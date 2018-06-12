package jabba

import (
	fswatch "github.com/andreaskoch/go-fswatch"
)

// Since vgo does not look at all build tags, the magefile dependencies will
// not be vendored. Thus, this file is used to import all the dependencies of
// magefile.go.

func init() {
	fswatch.NewFileWatcher("", 0)
}
