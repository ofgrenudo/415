// Package version holds build metadata injected at compile time via
// -ldflags "-X ...". Defaults below apply to `go run`/`go build` without
// the Makefile.
package version

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	Author  = "unknown"
)

func String(program string) string {
	return fmt.Sprintf("%s %s (commit %s, built %s) by %s", program, Version, Commit, Date, Author)
}
