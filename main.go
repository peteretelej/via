package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

//ErrInvalidURL is returned when an invalid url is provided
var ErrInvalidURL = errors.New("invalid url submitted")

var (
	listen = flag.String("listen", "127.0.0.1:8080", "listen address for http server")
	server = flag.Bool("server", false, "launches web server")
)

var usage = func() {
	fmt.Fprintf(os.Stderr, `via resolves URLs

Usage: 
	via [URL to resolve]
	via [flags] 

Examples:
	via goo.gl/OZGX9M	Resolves the URL 	
	via -server		Launches a webserver at localhost:8080
	via -listen :9000	Launches a webserver at 0.0.0.0:9000

Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if *server || *listen != "127.0.0.1:8080" {
		Serve(*listen)
		return
	}

	if len(os.Args) < 2 {
		fmt.Print("A URL is required, see -help")
		os.Exit(1)
	}
	res, err := ResolveURL(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

// ResolveURL returns the final URL address on visiting the provided url
func ResolveURL(u string) (string, error) {
	theurl, err := url.Parse(u)
	if err != nil {
		log.Print(err)
		return "", ErrInvalidURL
	}
	if theurl.Scheme == "" {
		theurl.Scheme = "http"
	}
	resp, err := http.Get(theurl.String())
	if err != nil {
		return "", err
	}

	return resp.Request.URL.String(), nil
}

// Serve launches a http server on the specified listen address and serves an html page for resolving urls
func Serve(listenAddr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	svr := &http.Server{
		Addr:           listenAddr,
		Handler:        mux,
		ReadTimeout:    time.Minute,
		WriteTimeout:   2 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Launching http server on %s", listenAddr)
	log.Fatal(svr.ListenAndServe())
}

var indexTmpl = template.Must(template.New("index").Parse(indexHTML))

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<title>via</title>
<meta name="viewport" content="width=device-width,initial-scale=1" />
<style>
</style>
</head>
<body>
<h3>via - URL Resolver</h3>
<form action="/" method="get">
	<label>URL
	<input name="url" value="{{.U}}" required />
	</label>
	<input type="submit" class="button" value="Get Final URL" />
</form>

{{if .Err}}
<p class="background-color:rgb(250,150,150)"> Unable to resolve URL( {{.U}} ): {{.Err}}.</p>
{{end}}
{{if .Result}}
<div class="callout">
<p style="display:inline-block;">{{.Result}}</p>
<div>
{{end}}
</body>
</html>
`

func renderIndex(w http.ResponseWriter, data interface{}) {
	if err := indexTmpl.Execute(w, data); err != nil {
		log.Print(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	d := struct{ U, Result, Err string }{}
	d.U = r.URL.Query().Get("url")
	if d.U == "" {
		renderIndex(w, d)
		return
	}
	res, err := ResolveURL(d.U)
	if err != nil {
		d.Err = "resolution failed"
		if err == ErrInvalidURL {
			d.Err = err.Error()
		}
		renderIndex(w, d)
		return
	}
	d.Result = res
	renderIndex(w, d)
}
