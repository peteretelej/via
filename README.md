# via - URL Resolver / URL Expander

- Expands shortened URLs
- Get final URL from a link with redirects (or shortened link)

## Usage

Installing
``` bash
go get -u github.com/peteretelej/via
```

Resolving a URL from the command line
``` bash
via bit.ly/3jHZKEC
```

Running your own instance of the web server
``` bash
via -server 

via -server -log
# launcher server and logs all resolution requests  (debug)
```
Server launches at http://localhost:8080 with a Web UI for expanding URLs


You can then reverse proxy the above instance with your favourite web server. (nginx, caddy..)

Or run it live 
``` bash
via -server -listen :8080
# launches on 0.0.0.0:8080 (public)
``` 

Web UI for the launched `-server`

<img src="https://user-images.githubusercontent.com/2271973/128786970-42d0618c-2f6b-4af3-950f-c7595a5a5455.png" width="400" height="400">

License: **MIT**
