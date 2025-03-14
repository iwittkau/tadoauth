package server

import "context"

type Config struct {
	SessionPath string
	Address     string
}

type Server struct {
	config  *Config
	session *session
}

func New(c *Config) *Server {
	return &Server{config: c}
}

func (s *Server) Serve(ctx context.Context) error {
	var err error
	if s.session, err = loadSession(s.config.SessionPath); err != nil {
		s.session = &session{path: s.config.SessionPath}
	}

	return nil
}
