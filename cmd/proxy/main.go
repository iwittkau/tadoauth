package main

import (
	"cmp"
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	conf := configFromEnv()
	grp, _ := errgroup.WithContext(ctx)
	grp.SetLimit(-1)
	grp.Go(func() error { return runProxy(ctx, conf) })
	grp.Go(func() error { return runWebUI(ctx, conf) })
	switch err := grp.Wait(); {
	case errors.Is(err, net.ErrClosed):
		return nil
	default:
		return err
	}
}

type config struct {
	proxyAddress string
	webUIAddress string
	sessionPath  string
}

func configFromEnv() *config {
	var c config
	c.proxyAddress = cmp.Or(os.Getenv("PROXY_ADDRESS"), ":8080")
	c.webUIAddress = cmp.Or(os.Getenv("WEBUI_ADDRESS"), ":8090")
	c.sessionPath = cmp.Or(os.Getenv("SESSION_PATH"), "session.json")
	return &c
}

func listenContext(ctx context.Context, addr string) (net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		lis.Close()
	}()
	return lis, nil
}

func runProxy(ctx context.Context, conf *config) error {
	http.HandleFunc("/", handleProxy)
	lis, err := listenContext(ctx, conf.proxyAddress)
	if err != nil {
		return err
	}
	log.Println("starting proxy", lis.Addr())
	return http.Serve(lis, nil)
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, "https://my.tado.com"+r.RequestURI, r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = r.Header.Clone()
	proxyReq.Header.Set("Authorization", "bearer "+os.Getenv("ACCESS_TOKEN"))
	log.Println(proxyReq)
	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, err = io.Copy(w, resp.Body); err != nil {
		log.Println(err)
		http.Error(w, "failed to copy response body", http.StatusInternalServerError)
		return
	}
}

func runWebUI(ctx context.Context, conf *config) error {
	lis, err := listenContext(ctx, conf.webUIAddress)
	if err != nil {
		return err
	}
	log.Println("starting web UI", lis.Addr())
	return http.Serve(lis, nil)
}
