# Jabba
Fast and simple link shortener.

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

The development server can then be started. The server will be restarted 
automatically when changes are made.
```bash
$ mage development
```

## Production
Jabba is packaged into a single binary, which can simply be copied onto the
server. The `.env` file is used for configuration.
```bash
$ mage production
```
