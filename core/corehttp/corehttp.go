/*
Package corehttp provides utilities for the webui, gateways, and other
high-level HTTP interfaces to IPFS.
*/
package corehttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	config "github.com/ipfs/kubo/config"
	core "github.com/ipfs/kubo/core"
	"github.com/jbenet/goprocess"
	periodicproc "github.com/jbenet/goprocess/periodic"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

var log = logging.Logger("core/server")

// shutdownTimeout is the timeout after which we'll stop waiting for hung
// commands to return on shutdown.
const shutdownTimeout = 30 * time.Second

// ServeOption registers any HTTP handlers it provides on the given mux.
// It returns the mux to expose to future options, which may be a new mux if it
// is interested in mediating requests to future options, or the same mux
// initially passed in if not.
type ServeOption func(*core.IpfsNode, net.Listener, *http.ServeMux) (*http.ServeMux, error)

// MakeHandler turns a list of ServeOptions into a http.Handler that implements
// all of the given options, in order.
func MakeHandler(n *core.IpfsNode, l net.Listener, options ...ServeOption) (http.Handler, error) {
	topMux := http.NewServeMux()
	mux := topMux
	for _, option := range options {
		var err error
		mux, err = option(n, l, mux)
		if err != nil {
			return nil, err
		}
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ServeMux does not support requests with CONNECT method,
		// so we need to handle them separately
		// https://golang.org/src/net/http/request.go#L111
		if r.Method == http.MethodConnect {
			w.WriteHeader(http.StatusOK)
			return
		}
		topMux.ServeHTTP(w, r)
	})
	return handler, nil
}

// ListenAndServe runs an HTTP server listening at |listeningMultiAddr| with
// the given serve options. The address must be provided in multiaddr format.
//
// TODO intelligently parse address strings in other formats so long as they
// unambiguously map to a valid multiaddr. e.g. for convenience, ":8080" should
// map to "/ip4/0.0.0.0/tcp/8080".
func ListenAndServe(n *core.IpfsNode, listeningMultiAddr string, options ...ServeOption) error {
	addr, err := ma.NewMultiaddr(listeningMultiAddr)
	if err != nil {
		return err
	}

	list, err := manet.Listen(addr)
	if err != nil {
		return err
	}

	// we might have listened to /tcp/0 - let's see what we are listing on
	addr = list.Multiaddr()
	fmt.Printf("RPC API server listening on %s\n", addr)

	return Serve(n, manet.NetListener(list), options...)
}

// Serve accepts incoming HTTP connections on the listener and pass them
// to ServeOption handlers.
func Serve(node *core.IpfsNode, lis net.Listener, options ...ServeOption) error {
	// make sure we close this no matter what.
	defer lis.Close()

	handler, err := MakeHandler(node, lis, options...)
	if err != nil {
		return err
	}
	cfg, err := node.Repo.Config()
	if err != nil {
		return err
	}

	middlewareHandler := DedicatedGatewayMiddleware(handler, cfg)

	addr, err := manet.FromNetAddr(lis.Addr())
	if err != nil {
		return err
	}

	select {
	case <-node.Process.Closing():
		return fmt.Errorf("failed to start server, process closing")
	default:
	}

	server := &http.Server{
		Handler: middlewareHandler,
	}

	var serverError error
	serverProc := node.Process.Go(func(p goprocess.Process) {
		serverError = server.Serve(lis)
	})

	// wait for server to exit.
	select {
	case <-serverProc.Closed():
	// if node being closed before server exits, close server
	case <-node.Process.Closing():
		log.Infof("server at %s terminating...", addr)

		warnProc := periodicproc.Tick(5*time.Second, func(_ goprocess.Process) {
			log.Infof("waiting for server at %s to terminate...", addr)
		})

		// This timeout shouldn't be necessary if all of our commands
		// are obeying their contexts but we should have *some* timeout.
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctx)

		// Should have already closed but we still need to wait for it
		// to set the error.
		<-serverProc.Closed()
		serverError = err

		warnProc.Close()
	}

	log.Infof("server at %s terminated", addr)
	return serverError
}

var ipLimiters = make(map[string]*rate.Limiter)
var cidLimiters = make(map[string]*rate.Limiter)
var mtx sync.Mutex

func getLimiter(limit string, limitMap map[string]*rate.Limiter, rps float64) *rate.Limiter {
	mtx.Lock()
	defer mtx.Unlock()

	limiter, exists := limitMap[limit]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Minute), int(rps))
		limitMap[limit] = limiter
	}

	return limiter
}

func DedicatedGatewayMiddleware(next http.Handler, cfg *config.Config) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is follow the pattern /ipfs/<hash>
		if cfg.ConfigPinningService.DedicatedGateway && strings.HasPrefix(r.URL.Path, "/ipfs/") {
			// Get the hash from the request URL
			pathPattern := regexp.MustCompile(`/ipfs/([^/]+)`)

			matches := pathPattern.FindStringSubmatch(r.URL.Path)
			if matches == nil || len(matches) < 2 {
				http.Error(w, "Invalid path", http.StatusBadRequest)
				return
			}
			cid, err := cid.Parse(matches[1])
			if err != nil {
				http.Error(w, "Invalid hash", http.StatusBadRequest)
				return
			}

			status, err := checkDmca(cid.String(), cfg)
			if err != nil {
				http.Error(w, err.Error(), status)
				return
			}
			// Call the getDedicatedGatewayAccess function
			status, err = getDedicatedGatewayAccess(cid.Hash().HexString(), cfg)
			if err != nil {
				http.Error(w, err.Error(), status)
				return
			}
		} else if !cfg.ConfigPinningService.DedicatedGateway && strings.HasPrefix(r.URL.Path, "/ipfs/") {
			ipLimiter := getLimiter(r.RemoteAddr, ipLimiters, 100)
			if !ipLimiter.Allow() {
				http.Error(w, "Too many requests from this IP", http.StatusTooManyRequests)
				return
			}
			pathPattern := regexp.MustCompile(`/ipfs/([^/]+)`)

			matches := pathPattern.FindStringSubmatch(r.URL.Path)
			if matches == nil || len(matches) < 2 {
				http.Error(w, "Invalid path", http.StatusBadRequest)
				return
			}
			cid, err := cid.Parse(matches[1])
			if err != nil {
				http.Error(w, "Invalid hash", http.StatusBadRequest)
				return
			}

			cidLimiter := getLimiter(cid.String(), cidLimiters, 15)
			if !cidLimiter.Allow() {
				http.Error(w, "Too many requests for this CID", http.StatusTooManyRequests)
				return
			}

			status, err := checkDmca(cid.String(), cfg)
			if err != nil {
				http.Error(w, err.Error(), status)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func getDedicatedGatewayAccess(hash string, cfg *config.Config) (int, error) {
	apiUrl := fmt.Sprintf("%s/api/dedicatedGateways/%s", cfg.ConfigPinningService.PinningService, hash)
	req, err := http.NewRequest("GET", apiUrl, bytes.NewBuffer(nil))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("blockservice-API-Key", cfg.ConfigPinningService.BlockserviceApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error while calling dedicated gateway API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp.StatusCode, errors.New("No users have subscribed to this hash yet.")
	}
	return http.StatusOK, nil
}

func checkDmca(hash string, cfg *config.Config) (int, error) {
	apiUrl := fmt.Sprintf("%s/api/dmca/%s", cfg.ConfigPinningService.PinningService, hash)
	req, err := http.NewRequest("GET", apiUrl, bytes.NewBuffer(nil))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("blockservice-API-Key", cfg.ConfigPinningService.BlockserviceApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error while calling DMCA API")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusGone {
		return http.StatusGone, errors.New("The content that you requested has been blocked because of legal, abuse, malware or security reasons. Please contact support@w3ipfs.storage for more information")
	}

	if resp.StatusCode != 200 {
		return resp.StatusCode, errors.New("Something went wrong")
	}

	return http.StatusOK, nil
}
