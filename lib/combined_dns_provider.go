package mycaddy

import (
	"context"
	"fmt"
	"strings"

	myaddr_dns_provider "github.com/shadow750d6/myaddr-dns-provider/lib"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/duckdns"
	"github.com/libdns/libdns"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct {
	myaddr  *myaddr_dns_provider.Provider
	duckdns *duckdns.Provider
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "dns.providers.combined",
		New: func() caddy.Module {
			return &Provider{
				myaddr:  &myaddr_dns_provider.Provider{},
				duckdns: &duckdns.Provider{},
			}
		},
	}
}

// Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.duckdns.APIToken = repl.ReplaceAll(p.duckdns.APIToken, "")
	p.myaddr.Key = repl.ReplaceAll(p.myaddr.Key, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens.
// Syntax:
//
//	acme_dns combined {
//	  duckdns_token TOKEN
//	  myaddr_key KEY
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// No extra args after provider name
		if d.NextArg() {
			return d.ArgErr()
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "duckdns_token":
				if p.duckdns.APIToken != "" {
					return d.Err("duckdns_token already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.duckdns.APIToken = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "myaddr_key":
				if p.myaddr.Key != "" {
					return d.Err("myaddr_key already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.myaddr.Key = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized option '%s'", d.Val())
			}
		}
	}

	if p.duckdns.APIToken == "" {
		return d.Err("missing duckdns_token")
	}
	if p.myaddr.Key == "" {
		return d.Err("missing myaddr_key")
	}

	return nil
}

// AppendRecords adds records to a zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	switch strings.Trim(zone, ".") {
	case "duckdns.org":
		return p.duckdns.AppendRecords(ctx, zone, records)
	case "myaddr.io":
	case "myaddr.dev":
	case "myaddr.tools":
		return p.myaddr.AppendRecords(ctx, zone, records)
	}
	return nil, fmt.Errorf("unsupported zone %s", zone)
}

// DeleteRecords deletes records from a zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	switch strings.Trim(zone, ".") {
	case "duckdns.org":
		return p.duckdns.DeleteRecords(ctx, zone, records)
	case "myaddr.io":
	case "myaddr.dev":
	case "myaddr.tools":
		return p.myaddr.DeleteRecords(ctx, zone, records)
	}
	return nil, fmt.Errorf("unsupported zone %s", zone)
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
