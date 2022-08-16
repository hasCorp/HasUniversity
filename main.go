package main

import (
	"flag"

	"github.com/hascorp/hasuniversity/cmd/hasuniversity"
)

var port = flag.Int("port", 8200, "port that the HTTP server listens on")
var cassandra = flag.String("cassandra", "cassandra:9042", "URL that points to the Cassandra cluster")

func main() {
	flag.Parse()

	hasuniversity.Start(*port, *cassandra)
}
