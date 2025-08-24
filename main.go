package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/mholt/caddy-webdav"
	_ "github.com/shadow750d6/caddy-with-custom/lib"
)

func main() {
	caddycmd.Main()
}
