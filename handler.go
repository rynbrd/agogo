package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

var (
	handlerClient   *http.Client   = &http.Client{}
	handlerSplitter *regexp.Regexp = regexp.MustCompile(`[ \t]+`)
)

type Handler struct {
	links     *url.URL
	reload    bool
	allow     []net.IP
	redirects map[string]string
	mutex     sync.Mutex
}

// Create a new handler which loads its redirects from `links`.
func NewHandler(links *url.URL, reload bool, allow []net.IP) *Handler {
	h := &Handler{
		links:     links,
		reload:    reload,
		allow:     allow,
		redirects: map[string]string{},
	}
	h.load()
	return h
}

func (h *Handler) allowed(req *http.Request) bool {
	if len(h.allow) == 0 {
		return true
	}
	parts := strings.SplitN(req.RemoteAddr, ":", 2)
	if len(parts) < 1 {
		return false
	}
	ip := net.ParseIP(parts[0])
	if ip == nil {
		return false
	}
	for _, allowIP := range h.allow {
		if ip.Equal(allowIP) {
			return true
		}
	}
	return false
}

func (h *Handler) load() error {
	res, err := handlerClient.Get(h.links.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	n := 0
	sc := bufio.NewScanner(res.Body)
	redirects := map[string]string{}
	for sc.Scan() {
		n++
		line := strings.TrimSpace(sc.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := handlerSplitter.Split(line, 2)
		if len(parts) != 2 {
			log.Printf("ignoring invalid line %d", n)
			continue
		}
		redirects[CleanPath(parts[0])] = parts[1]
	}
	h.redirects = redirects
	log.Printf("loaded %d redirects", n)
	return nil
}

func (h *Handler) Load() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.load()
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	cleanPath := CleanPath(req.URL.Path)
	if h.reload && cleanPath == "_reload" {
		if !h.allowed(req) {
			http.Error(res, "permission denied", 403)
		} else {
			if err := h.load(); err == nil {
				res.WriteHeader(200)
				res.Write([]byte("links reloaded"))
			} else {
				http.Error(res, "internal server error", 500)
			}
		}
	} else {
		if location, ok := h.redirects[cleanPath]; ok {
			res.Header().Set("Location", location)
			res.WriteHeader(307)
		} else {
			http.Error(res, "not found", 404)
		}
	}
}

// Remove leading and trailing slashes from a path.
func CleanPath(path string) string {
	return strings.Trim(path, "/")
}
