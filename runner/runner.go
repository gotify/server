package runner

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/config"
	"golang.org/x/crypto/acme/autocert"
)

// Run starts the http server and if configured a https server.
func Run(engine *gin.Engine, conf *config.Configuration) {
	var httpHandler http.Handler = engine

	if *conf.Server.SSL.Enabled {
		fmt.Println(*conf.Server.SSL.RedirectToHTTPS)
		if *conf.Server.SSL.RedirectToHTTPS {
			httpHandler = redirectToHTTPS(string(conf.Server.SSL.Port))
		}

		s := &http.Server{
			Addr:    fmt.Sprintf(":%d", conf.Server.SSL.Port),
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
		fmt.Println("Started Listening on port", conf.Server.SSL.Port)
		go log.Fatal(s.ListenAndServeTLS(conf.Server.SSL.CertFile, conf.Server.SSL.CertKey))
	}
	fmt.Println("Started Listening on port", conf.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Server.Port), httpHandler))
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
