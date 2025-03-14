package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/iwittkau/tadoauth/auth"
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	http.DefaultClient.Timeout = 10 * time.Second

	switch flag.Arg(0) {
	case "setup":
		if err := setup(); err != nil {
			log.Fatal(err)
		}
	case "token":
		deviceCode := flag.Arg(1)
		if deviceCode == "" {
			log.Fatal("no device code provided")
		}
		res, err := auth.NewClient().ExchangeDeviceCode(context.Background(), deviceCode)
		if err != nil {
			log.Fatal(err)
		}
		res.Print()
	case "refresh":
		refreshToken := flag.Arg(1)
		if refreshToken == "" {
			log.Fatal("no refresh token provided")
		}
		res, err := auth.NewClient().RefreshToken(context.Background(), refreshToken)
		if err != nil {
			log.Fatal(err)
		}
		res.Print()
	case "":
		log.Fatal("missing command")
	default:
		log.Fatal("unknown command")
	}
}

func setup() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := auth.NewClient()
	deviceAuthRes, err := client.RequestRegistration(ctx)
	if err != nil {
		return fmt.Errorf("requesting registration: %w", err)
	}
	log.Print("Complete the registration here:", deviceAuthRes.VerificationURIComplete)
	log.Printf("This link will expire at %s.", deviceAuthRes.ExpiresAt().Local())
	log.Print("Press enter after you've completed the registration, Ctrl+C to abort.")
	if err := awaitUserConfirmation(); err != nil {
		return err
	}

	tokenRes, err := client.ExchangeDeviceCode(ctx, deviceAuthRes.DeviceCode)
	if err != nil {
		return err
	}
	tokenRes.Print()
	log.Print("Press enter if you want to test token refreshing, Ctrl+C to abort.")
	if err := awaitUserConfirmation(); err != nil {
		return err
	}

	tokenRes, err = client.RefreshToken(ctx, tokenRes.RefreshToken)
	if err != nil {
		return err
	}
	tokenRes.Print()

	return nil
}

func awaitUserConfirmation() error {
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("user aborted")
		}
		return fmt.Errorf("reading os.Stdin: %w", err)
	}

	return nil
}
