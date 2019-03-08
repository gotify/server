package runner

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/config"
	"golang.org/x/crypto/acme/autocert"
)

// Run starts the http server and if configured a https server.
func Run(engine *gin.Engine, conf *config.Configuration) {
	var httpHandler http.Handler = engine

	if *conf.Server.SSL.Enabled {
		if *conf.Server.SSL.RedirectToHTTPS {
			httpHandler = redirectToHTTPS(string(conf.Server.SSL.Port))
		}

		addr := fmt.Sprintf("%s:%d", conf.Server.SSL.ListenAddr, conf.Server.SSL.Port)
		s := &http.Server{
			Addr:    addr,
			Handler: engine,
		}

		if *conf.Server.SSL.LetsEncrypt.Enabled {
			certManager := autocert.Manager{
				Prompt:     func(tosURL string) bool { return *conf.Server.SSL.LetsEncrypt.AcceptTOS },
				HostPolicy: autocert.HostWhitelist(conf.Server.SSL.LetsEncrypt.Hosts...),
				Cache:      autocert.DirCache(conf.Server.SSL.LetsEncrypt.Cache),
			}
			httpHandler = certManager.HTTPHandler(httpHandler)
			s.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}
		}
		fmt.Println("Started Listening for TLS connection on " + addr)
		go func() {
			log.Fatal(s.ListenAndServeTLS(conf.Server.SSL.CertFile, conf.Server.SSL.CertKey))
		}()
	}
	addr := fmt.Sprintf("%s:%d", conf.Server.ListenAddr, conf.Server.Port)
	fmt.Println("Started Listening for plain HTTP connection on " + addr)
	log.Fatal(http.ListenAndServe(addr, httpHandler))
}

func redirectToHTTPS(port string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			http.Error(w, "Use HTTPS", http.StatusBadRequest)
			return
		}

		target := "https://" + changePort(r.Host, port) + r.URL.RequestURI()
		http.Redirect(w, r, target, http.StatusFound)
	}
}

func changePort(hostPort string, port string) string {
	host, _, err := net.SplitHostPort(hostPort)
	if err != nil {
		return hostPort
	}
	return net.JoinHostPort(host, port)
}
