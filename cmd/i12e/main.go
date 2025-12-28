//
// build:
//   $ go build -buildvcs=false ./cmd/i12e
//

package main

import "github.com/sfmunoz/logit"

var log = logit.Logit().
	WithLevel(logit.LevelNotice).
	With("app", "i12e")

func main() {
	log.Error("i12e not implemented yet")
}
