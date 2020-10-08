// httpserver
package webserver

import (
	"crypto/tls"
	"crypto/x509"
	"gouploader/runconfig"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
	logrus "github.com/sirupsen/logrus"
)

type WebServer struct {
	srv    http.Server
	Port   int
	Router *mux.Router
}

func (w WebServer) StartServer() {

	log := logrus.WithFields(logrus.Fields{
		"file":   "httpserver.go",
		"method": "StartServer",
	})

	handler := cors.Default().Handler(w.Router)

	w.srv.Addr = ":" + strconv.Itoa(w.Port)
	w.srv.Handler = handler

	if runconfig.Config.GetString("check_client_certificat") == "check" {

		caCert, err := ioutil.ReadFile(runconfig.Config.GetString("client_certificat_file"))
		if err != nil {
			log.Error(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		cfg := &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caCertPool,
		}
		w.srv.TLSConfig = cfg
	} else {
		cfg := &tls.Config{}
		w.srv.TLSConfig = cfg
	}

	var err error
	if runconfig.Config.GetString("http_server_mode") == "https" {
		//log.Trace("httpS YES")
		err = w.srv.ListenAndServeTLS(runconfig.Config.GetString("ServerSrtFile"), runconfig.Config.GetString("ServerKeyFile"))
	} else {
		//log.Trace("httpS NO")
		err = w.srv.ListenAndServe()
	}

	if err == nil {
		log.Trace("Started http server on port: ", w.Port, " mode: ", runconfig.Config.GetString("http_server_mode"))
		log.Info("Started http server on port: ", w.Port, " mode: ", runconfig.Config.GetString("http_server_mode"))

	} else {
		log.Error("http server error: ", err)
	}

}

//not works :(
// func (w WebServer) StopServer(ctx context.Context) {
// 	log := logrus.WithFields(logrus.Fields{
// 		"file":   "httpserver.go",
// 		"method": "StopServer",
// 	})

// 	if err := w.srv.Shutdown(ctx); err != nil {
// 		log.Trace("WebServer not stoped with error: ", err)
// 		log.Error("WebServer not stoped with error: ", err)
// 	} else {
// 		log.Trace("WebServer stoped.")
// 		log.Info("WebServer stoped.")
// 	}
// }
