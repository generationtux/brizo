# Brizo

PAAS that runs on top of Kubernetes.

## Development

#### Autoreload

You can use [Gin](https://github.com/codegangsta/gin) to serve and autoreload the service during development.

Install Gin
```sh
$ go get github.com/codegangsta/gin
```

Then run the app
```sh
$ gin -a 8080
```
*8080 is the port this app is configured to listen on. If you've customized this you should use the appropriate port*

#### Dependency management

Dependencies are managed using [Govendor](https://github.com/kardianos/govendor). For simpliciy in version management
dependencies are committed with the repo.
