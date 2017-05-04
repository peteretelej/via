package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"
)

var listen = flag.String("listen", "127.0.0.1:8080", "listen address for http server")

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Serve(*listen)
}

// ResolveURL returns the final URL address on visiting the provided url
func ResolveURL(u string) (string, error) {
	resp, err := http.Get(u)
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
	}
	res, err := ResolveURL(d.U)
	if err != nil {
		d.Err = "resolution failed"
		renderIndex(w, d)
	}
	d.Result = res
	renderIndex(w, d)
}
