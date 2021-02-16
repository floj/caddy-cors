package caddycors

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("cors", parseCORS)
}

// parseCORS parses the Caddyfile tokens for the webdav directive.
func parseCORS(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	m := &CORSMiddleware{}
	err := m.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalCaddyfile sets up the handler from Caddyfile tokens.
//
//    cors [<matcher>] {
//        origin          <origin>
//        allowed_methods <methods>...
//        allowed_headers <headers>...
//        exposed_headers <headers>...
//        max_age         <seconds>
//        allow_credentials
//    }
//
func (m *CORSMiddleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for d.NextBlock(0) {
			switch d.Val() {
			case "origin":
				if !d.NextArg() {
					return d.ArgErr()
				}
				m.Origins = append(m.Origins, d.Val())
			case "allowed_methods":
				if m.AllowedMethods != nil {
					return d.Err("allowed_methods already specified")
				}
				m.AllowedMethods = append(m.AllowedMethods, d.RemainingArgs()...)
			case "allowed_headers":
				if m.AllowedHeaders != nil {
					return d.Err("allowed_headers already specified")
				}
				m.AllowedHeaders = append(m.AllowedHeaders, d.RemainingArgs()...)
			case "exposed_headers":
				if m.ExposedHeaders != nil {
					return d.Err("exposed_headers already specified")
				}
				m.ExposedHeaders = append(m.ExposedHeaders, d.RemainingArgs()...)
			case "allow_credentials":
				if d.NextArg() {
					return d.ArgErr()
				}
				m.AllowCredentials = true
			case "max_age":
				if !d.NextArg() {
					return d.ArgErr()
				}
				m.MaxAge = 1
			default:
				return d.Errf("unrecognized subdirective: %s", d.Val())
			}
		}
	}
	return nil
}
