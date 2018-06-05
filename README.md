# Jabba
Fast and simple link shortener with cool features!

## Development
Some commands need to be installed to work on Jabba.
```bash
# Build system
$ go get -u -d github.com/magefile/mage
$ cd $GOPATH/src/github.com/magefile/mage
$ go run bootstrap.go

# Handling static files
$ go get -u github.com/gobuffalo/packr/packr
```

The development server can then be started. Static files can be reloaded without
restarting the server.
```bash
$ mage development
```

## Production
Jabba is packaged into a single binary, which can simply be copied onto the
server. The `.env` file is used for configuration.
```bash
$ mage production
```