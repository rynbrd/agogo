package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"strings"
)

const (
	DefaultListen   = ":http"
	DefaultNoReload = false
)

// An array of strings that can be used as an flag var.
type Strings []string

func (strs *Strings) String() string {
	return strings.Join(*strs, ",")
}

func (strs *Strings) Set(value string) error {
	*strs = append(*strs, strings.Split(value, ",")...)
	return nil
}

type Config struct {
	Links    *url.URL
	Listen   string
	Allow    []net.IP
	NoReload bool
}

// Configure builds a populated `Config` struct using the provided command line
// arguments. Should that fail an error will be returned.
func Configure(args []string) (*Config, error) {
	var links string = ""
	var allow Strings
	var err error
	cfg := &Config{}
	flags := flag.NewFlagSet(args[0], 0)
	flags.StringVar(&links, "links", "", "A URL pointing to the links database.")
	flags.StringVar(&cfg.Listen, "listen", DefaultListen, "The address to listen on.")
	flags.Var(&allow, "allow", "Only allow reload requests from this host.")
	flags.BoolVar(&cfg.NoReload, "no-reload", DefaultNoReload, "Disable the reload endpoint.")
	if err = flags.Parse(args[1:]); err != nil {
		return nil, err
	}
	if links == "" {
		return nil, fmt.Errorf("Must provide a value to -links.")
	}
	if cfg.Links, err = url.Parse(links); err != nil {
		return nil, fmt.Errorf("Invalid links url '%s'", links)
	}
	cfg.Allow = []net.IP{}
	for _, host := range allow {
		if addrs, err := net.LookupIP(host); err == nil {
			cfg.Allow = append(cfg.Allow, addrs...)
		} else {
			return nil, fmt.Errorf("Allowed host '%s' cannot be resolved.", host)
		}
	}
	return cfg, nil
}
