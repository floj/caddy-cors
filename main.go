package caddycors

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/floj/caddy-cors/cors"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(&CORSMiddleware{})
}

// CORSMiddleware implements an HTTP handler that adds CORS headers
type CORSMiddleware struct {
	Origins          []string `json:"allowed_origins,omitempty"`
	AllowedMethods   []string `json:"allowed_methods,omitempty"`
	AllowedHeaders   []string `json:"allowed_headers,omitempty"`
	ExposedHeaders   []string `json:"exposed_headers,omitempty"`
	AllowCredentials bool     `json:"allow_credentials,omitempty"`
	MaxAge           int      `json:"max_age,omitempty"`
	OptionStatusCode int      `json:"option_status_code,omitempty"`

	c   *cors.Cors
	log *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (*CORSMiddleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.cors",
		New: func() caddy.Module { return &CORSMiddleware{} },
	}
}

// Provision implements caddy.Provisioner.
func (m *CORSMiddleware) Provision(ctx caddy.Context) error {
	m.log = ctx.Logger(m)

	if m.OptionStatusCode == 0 {
		m.OptionStatusCode = http.StatusNoContent
	}

	opts := []cors.CORSOption{
		cors.OptionStatusCode(m.OptionStatusCode),
	}

	if len(m.Origins) > 0 {
		opts = append(opts, cors.AllowedOrigins(m.Origins))
	}
	if len(m.AllowedMethods) > 0 {
		opts = append(opts, cors.AllowedMethods(m.AllowedMethods))
	}
	if len(m.AllowedHeaders) > 0 {
		opts = append(opts, cors.AllowedHeaders(m.AllowedHeaders))
	}
	if len(m.ExposedHeaders) > 0 {
		opts = append(opts, cors.ExposedHeaders(m.ExposedHeaders))
	}
	if m.AllowCredentials {
		opts = append(opts, cors.AllowCredentials())
	}
	if m.MaxAge > 0 {
		opts = append(opts, cors.MaxAge(m.MaxAge))
	}

	m.c = cors.CORS(opts...)
	return nil
}

// Validate implements caddy.Validator.
func (m *CORSMiddleware) Validate() error {
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m *CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	var err error
	m.c.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
		err = next.ServeHTTP(w, r)
	})
	return err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*CORSMiddleware)(nil)
	_ caddy.Validator             = (*CORSMiddleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*CORSMiddleware)(nil)
	_ caddyfile.Unmarshaler       = (*CORSMiddleware)(nil)
)
