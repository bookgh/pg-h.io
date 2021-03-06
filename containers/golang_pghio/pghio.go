package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	pghioGetCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pghio_get_count",
		Help: "is the count of page requests to / since server started.",
	})
	pghioPrivacyGetCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pghio_privacy_get_count",
		Help: "is the count of page requests to /privacy.html since server started.",
	})
	pghioPostCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pghio_post_count",
		Help: "is the count of posts to the command prompt since server started.",
	})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(pghioGetCount)
	prometheus.MustRegister(pghioPrivacyGetCount)
	prometheus.MustRegister(pghioPostCount)

}

func main() {
	log.Println("Starting pg-h.io")
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/privacy.html", privacy)
	http.HandleFunc("/google776b578cc5a81cc0.html", google)
	http.HandleFunc("/sitemap.txt", sitemap)
	http.HandleFunc("/", pghio)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":80", nil)
}

func pghio(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		pghioGetCount.Inc()
		tmpl := template.New("index.html")
		tmpl, err := tmpl.ParseFiles("html/index.html")
		if err != nil {
			log.Println("ERROR: pghio tmpl.ParseFiles", err)
		}
		tmpl.Execute(w, "")
		logRequestInfo(r)
	}
	if r.Method == "POST" {
		pghioPostCount.Inc()
		r.ParseForm() // need to parse the form before interacting with form data
		text := (r.Form["text"])[0]
		switch text {
		case "pg --alertmanager":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://alertmanager.pg-h.io/", http.StatusFound)
		case "pg --alertmanager-metrics":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://alertmanager.pg-h.io/metrics", http.StatusFound)
		case "pg --blog":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://blog.pg-h.io/", http.StatusFound)
		case "pg --blog-metrics":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://blog.pg-h.io/metrics", http.StatusFound)
		case "pg --cadvisor-metrics":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://cadvisor.pg-h.io/metrics", http.StatusFound)
		case "pg --github":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "https://github.com/silentpete", http.StatusFound)
		case "pg --grafana":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://grafana.pg-h.io/", http.StatusFound)
		case "pg --grafana-metrics":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://grafana.pg-h.io/metrics", http.StatusFound)
		case "pg --help":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "/", http.StatusFound)
		case "pg -h":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "/", http.StatusFound)
		case "pg --influxdb-metrics":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://influxdb.pg-h.io/metrics", http.StatusFound)
		case "pg --node-exporter-metrics":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://node-exporter.pg-h.io/metrics", http.StatusFound)
		case "pg --privacy":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://pg-h.io/privacy.html", http.StatusFound)
		case "pg --prometheus":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://prometheus.pg-h.io/", http.StatusFound)
		case "pg --prometheus-metrics":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://prometheus.pg-h.io/metrics", http.StatusFound)
		case "pg --twitter":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "https://twitter.com/PeterGallerani", http.StatusFound)
		case "pg --resume":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "https://www.linkedin.com/in/petegallerani/", http.StatusFound)
		case "pg --site-metrics":
			log.Printf("success POST: \"%v\"\n", text)
			http.Redirect(w, r, "http://pg-h.io/metrics", http.StatusFound)
		default:
			log.Printf("failed POST: \"%v\"\n", text)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func logRequestInfo(r *http.Request) {
	// comes through the proxy, print IP proxied
	if len(r.Header["X-Real-Ip"]) > 0 {
		for _, ip := range r.Header["X-Real-Ip"] {
			log.Println(ip, "requested", r.RequestURI)
		}
	} else {
		// doesn't come through the proxy
		log.Println(r.RemoteAddr, "requested", r.RequestURI)
	}
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "imgs/favicon.ico")
	logRequestInfo(r)
}

// sitemap is the handler used for requests to /sitemap.txt
func sitemap(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "files/sitemap.txt")
	logRequestInfo(r)
}

func google(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "google-site-verification: google776b578cc5a81cc0.html")
	logRequestInfo(r)
}

func privacy(w http.ResponseWriter, r *http.Request) {
	pghioPrivacyGetCount.Inc()
	tmpl := template.New("privacy.html")
	tmpl, err := tmpl.ParseFiles("html/privacy.html")
	if err != nil {
		log.Println("ERROR: pghio privicy tmpl.ParseFiles", err)
	}
	tmpl.Execute(w, "")
	logRequestInfo(r)
}
