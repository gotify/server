package runner

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gotify/server/v2/config"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// Run starts the http server and if configured a https server.
func Run(router http.Handler, conf *config.Configuration) error {
	shutdown := make(chan error)
	go doShutdownOnSignal(shutdown)

	httpListener, err := startListening("plain connection", conf.Server.ListenAddr, conf.Server.Port, conf.Server.KeepAlivePeriodSeconds)
	if err != nil {
		return err
	}
	defer httpListener.Close()

	s := &http.Server{Handler: router}
	if conf.Server.SSL.Enabled {
		if conf.Server.SSL.LetsEncrypt.Enabled {
			applyLetsEncrypt(s, conf)
		} else if conf.Server.SSL.CertFile != "" && conf.Server.SSL.CertKey != "" {
			log.Fatalln("CertFile and CertKey must be set to use HTTPS when LetsEncrypt is disabled, please set GOTIFY_SERVER_SSL_CERTFILE and GOTIFY_SERVER_SSL_CERTKEY")
		}

		httpsListener, err := startListening("TLS connection", conf.Server.SSL.ListenAddr, conf.Server.SSL.Port, conf.Server.KeepAlivePeriodSeconds)
		if err != nil {
			return err
		}
		defer httpsListener.Close()

		go func() {
			err := s.ServeTLS(httpsListener, conf.Server.SSL.CertFile, conf.Server.SSL.CertKey)
			doShutdown(shutdown, err)
		}()
	}
	go func() {
		err := s.Serve(httpListener)
		doShutdown(shutdown, err)
	}()

	err = <-shutdown
	fmt.Println("Shutting down:", err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Shutdown(ctx)
}

func doShutdownOnSignal(shutdown chan<- error) {
	onSignal := make(chan os.Signal, 1)
	signal.Notify(onSignal, os.Interrupt, syscall.SIGTERM)
	sig := <-onSignal
	doShutdown(shutdown, fmt.Errorf("received signal %s", sig))
}

func doShutdown(shutdown chan<- error, err error) {
	select {
	case shutdown <- err:
	default:
		// If there is no one listening on the shutdown channel, then the
		// shutdown is already initiated and we can ignore these errors.
	}
}

func startListening(connectionType, listenAddr string, port, keepAlive int) (net.Listener, error) {
	network, addr := getNetworkAndAddr(listenAddr, port)
	lc := net.ListenConfig{KeepAlive: time.Duration(keepAlive) * time.Second}

	oldMask := umask(0)
	defer umask(oldMask)

	l, err := lc.Listen(context.Background(), network, addr)
	if err == nil {
		fmt.Println("Started listening for", connectionType, "on", l.Addr().Network(), l.Addr().String())
	}
	return l, err
}

func getNetworkAndAddr(listenAddr string, port int) (string, string) {
	if strings.HasPrefix(listenAddr, "unix:") {
		return "unix", strings.TrimPrefix(listenAddr, "unix:")
	}
	return "tcp", fmt.Sprintf("%s:%d", listenAddr, port)
}

type LoggingRoundTripper struct {
	Name         string
	RoundTripper http.RoundTripper
}

func (l *LoggingRoundTripper) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	resp, err = l.RoundTripper.RoundTrip(r)
	if resp.StatusCode == 429 {
		log.Printf("%s Rate Limited: Retry-After %s on %s %s\n", l.Name, resp.Header.Get("Retry-After"), r.Method, r.URL.String())
	} else if resp.StatusCode >= 400 {
		log.Printf("%s Request Failed: Unexpected status code %d on %s %s\n", l.Name, resp.StatusCode, r.Method, r.URL.String())
	} else if err != nil {
		log.Printf("%s Request Failed: %s on %s %s\n", l.Name, err.Error(), r.Method, r.URL.String())
	}
	return
}

func applyLetsEncrypt(s *http.Server, conf *config.Configuration) {
	httpClient := &http.Client{
		Transport: &LoggingRoundTripper{Name: "Let's Encrypt", RoundTripper: http.DefaultTransport},
		Timeout:   60 * time.Second,
	}

	acmeClient := &acme.Client{
		HTTPClient:   httpClient,
		DirectoryURL: conf.Server.SSL.LetsEncrypt.DirectoryURL,
	}
	certManager := autocert.Manager{
		Client: acmeClient,
		Prompt: func(tosURL string) bool {
			if !conf.Server.SSL.LetsEncrypt.AcceptTOS {
				log.Fatalf("Let's Encrypt TOS must be accepted to use Let's Encrypt, please acknowledge TOS at %s and set GOTIFY_SERVER_SSL_LETSENCRYPT_ACCEPTTOS=true\n", tosURL)
			}
			return true
		},
		HostPolicy: autocert.HostWhitelist(conf.Server.SSL.LetsEncrypt.Hosts...),
		Cache:      autocert.DirCache(conf.Server.SSL.LetsEncrypt.Cache),
	}
	s.Handler = certManager.HTTPHandler(s.Handler)
	s.TLSConfig = certManager.TLSConfig()
}
