package server

import (
	"crypto/tls"
	"github.com/Jamshid90/go-clean-architecture/pkg/config"
	"net/http"
	"time"
)

type Server struct {
	config  *config.Config
	handler http.Handler
}

func NewServer(config *config.Config, handler http.Handler) *Server {
	return &Server{
		config:  config,
		handler: handler,
	}
}

func (s *Server) GetServerAddr() string {
	return s.config.Server.Host + ":" + s.config.Server.Port
}

func (s *Server) tlsConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256, tls.X25519},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

func (s *Server) Run() error {

	server := &http.Server{
		Addr:         s.GetServerAddr(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      s.handler,
		TLSConfig:    s.tlsConfig(),
	}

	if s.config.Server.Protocol == "https" {
		return server.ListenAndServeTLS(s.config.Server.SSLCert, s.config.Server.SSLPrivKey)
	}

	return server.ListenAndServe()
}
