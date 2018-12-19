# via - URL Resolver
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fpeteretelej%2Fvia.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fpeteretelej%2Fvia?ref=badge_shield)


- Expands shortened URLs
- Get final URL in a link with redirects
- etc


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

via -server -log
# launcher server and logs all resolution requests  (debug)
```


You can then reverse proxy the above instance with your favourite web server. (nginx, caddy..)

Or run it live 
``` bash
via -server -listen :8080
# launches on 0.0.0.0:8080 (public)
``` 



## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fpeteretelej%2Fvia.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fpeteretelej%2Fvia?ref=badge_large)