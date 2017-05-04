package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

//ErrInvalidURL is returned when an invalid url is provided
var ErrInvalidURL = errors.New("invalid url submitted")

var (
	listen      = flag.String("listen", "127.0.0.1:8080", "listen address for http server")
	server      = flag.Bool("server", false, "launches web server")
	logRequests = flag.Bool("log", false, "logs all requests for resolution from server")
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

var cli bool

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
	cli = true
	res, err := ResolveURL(os.Args[1], "GET", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

// ResolveURL returns the final URL address on visiting the provided url
func ResolveURL(u, method string, headers map[string]string) (string, error) {
	theurl, err := url.Parse(u)
	if err != nil {
		log.Print(err)
		return "", ErrInvalidURL
	}
	if theurl.Scheme == "" {
		theurl.Scheme = "http"
	}
	cl := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, theurl.String(), nil)
	if err != nil {
		return "", err
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	if *logRequests && !cli {
		defer logRequest(os.Stdout, req)
	}
	resp, err := cl.Do(req)
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
body{ display:block; width:80%; margin:20px  auto;
}
label{ display:inline-block; }
input,textarea { display:block; margin-bottom:20px;}
</style>
</head>
<body>
<h3>via - URL Resolver</h3>
<form action="/" method="post">
	<label for="url">URL</label>
	<input name="url" id="url"  style="width:300px" value="{{.U}}" required />
	<label for="method">Request Method</label>
	<input type="text" name="method" value="{{.Method}}" />
	<label for="headers">HTTP Headers (optional)</label>
	<textarea name="headers" style="width:500px;height:70px"
placeholder="Accept-Encoding:gzip
Accept-Language: en-US,en;q=0.8">{{with .Headers}}{{.}}{{end}}</textarea>
	<input type="submit" class="button" value="Get Final URL" />
	<a href="/">Clear Query</a>
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
	d := struct {
		U, Result, Err string
		Headers        string
		Method         string
	}{Method: "GET"}
	if r.Method != "POST" {
		renderIndex(w, d)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}
	m := r.Form.Get("method")
	if m != "" {
		d.Method = m
	}
	d.U = r.Form.Get("url")
	d.Headers = r.Form.Get("headers")

	headers := retrieveHeaders(d.Headers)
	res, err := ResolveURL(d.U, d.Method, headers)
	if err != nil {
		if err == ErrInvalidURL {
			d.Err = err.Error()
		}
		d.Err = "something strange happened"
		renderIndex(w, d)
		return
	}
	d.Result = res
	renderIndex(w, d)
}

func retrieveHeaders(s string) map[string]string {
	s = strings.Replace(s, "\r", "", -1)
	ss := strings.Split(s, "\n")
	h := make(map[string]string)
	for _, val := range ss {
		if strings.Count(val, ":") != 1 {
			continue
		}
		ind := strings.Index(val, ":")
		if ind < 1 {
			continue
		}
		if len(val) <= ind {
			continue
		}
		h[val[:ind]] = val[ind+1:]
	}
	return h
}

func logRequest(wr io.Writer, req *http.Request) {
	dmp, err := httputil.DumpRequest(req, true)
	if err != nil {
		return
	}
	fmt.Fprintf(wr, "%q\n", dmp)
}
