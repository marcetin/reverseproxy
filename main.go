package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/comhttp/jorm/pkg/utl"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type (
	JORM struct {
		WWW    *http.Server
		config map[string]string
	}
)

func main() {
	// Get cmd line parameters
	// service := flag.String("srv", "", "Service")
	// path := flag.String("path", "reverseproxy.json", "Path")
	port := flag.String("port", "80", "Port")
	loglevel := flag.String("loglevel", "debug", "Logging level (debug, info, warn, error)")
	flag.Parse()

	//j.Log.SetLevel(parseLogLevel(*loglevel))
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Default level for this example is info, unless debug flag is present

	switch *loglevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	j := new(JORM)

	j.WWW = &http.Server{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("reverse proxy")
	h := &baseHandle{}
	http.Handle("/", h)
	j.WWW.Handler = h
	j.WWW.Addr = ":11444"
	log.Fatal().Err(j.WWW.ListenAndServe())

	log.Info().Msg("Port: " + *port)

}

// var (
// 	hostTarget = map[string]string{
// 		"okno.rs":                    "http://127.0.0.1:4433",
// 		"parallelcoin.info":          "http://127.0.0.1:4433",
// 		"explorer.parallelcoin.info": "http://127.0.0.1:4433",
// 		"jorm.okno.rs":               "http://127.0.0.1:14411",
// 		"our.okno.rs":                "http://127.0.0.1:14422",
// 		"enso.okno.rs":               "http://127.0.0.1:14433",
// 		"p9c.okno.rs":                "http://127.0.0.1:1337",
// 		"admin.parallelcoin.io":      "http://127.0.0.1:11122",
// 		"api.parallelcoin.io":        "http://127.0.0.1:11123",
// 	}
// )

type baseHandle struct{}

func (h *baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	log.Print("target hosthosthost fail:", host)

	c, err := Read("reverseproxy.json")
	utl.ErrorLog(err)

	if target, ok := c[host]; ok {
		reverseproxy(w, r, target)
	} else {
		reverseproxy(w, r, "http://localhost:14444")
	}
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
}

func reverseproxy(w http.ResponseWriter, r *http.Request, target string) {
	remoteUrl, err := url.Parse(target)
	if err != nil {
		log.Print("target parse fail:", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	// w.Header().Set("AMP-Access-Control-Allow-Source-Origin", "*")
	// w.Header().Set("Access-Control-Expose-Headers", "AMP-Access-Control-Allow-Source-Origin")

	// w.Header().Set("Access-Control-Allow-Credentials","true")

	// w.Header().Set("AMP-Same-Origin","true")

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods","POST, GET, OPTIONS")

	// w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	proxy.ServeHTTP(w, r)
	return
}

//
//func (j *JORM) ReverseProxySRV() {
//	h := &BaseHandle{}
//	http.Handle("/", h)
//	server := &http.Server{
//		Addr:    ":80",
//		Handler: h,
//	}
//	log.Fatal(server.ListenAndServe())
//}
func status(w http.ResponseWriter, r *http.Request) {
	// Handles top-level page.
	fmt.Fprintf(w, "You are on the status home page")
}

// Read a record from the database
func Read(path string) (map[string]string, error) {
	// read record from database
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := make(map[string]string)
	// unmarshal data
	json.Unmarshal(b, &conf)
	return conf, err
}
