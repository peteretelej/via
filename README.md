# via - URL Resolver

- Expands shortened URLs
- Get final URL in a link with redirects
- etc

Demo: [via.etelej.com](https://via.etelej.com) 

## Usage

Installing
``` bash
go get -u github.com/peteretelej/via
```

Resolving a URL from the command line
``` bash
via goo.gl/OZGX9M
```

Running your own instance of the web server
``` bash
via -server 
# launches at localhost:8080
```

You can then reverse proxy the above instance with your favourite web server. (nginx, caddy..)

Or run it live 
``` bash
via -server -listen :8080
# launches on 0.0.0.0:8080 (public)
``` 

