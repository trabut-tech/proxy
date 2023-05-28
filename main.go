package main

import (
	"flag"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

var (
	flagListenAddr = flag.String("listen", ":33080", "the http address to start the proxy server on")
	flagVerbose    = flag.Bool("verbose", true, "whether to be verbose or not")
)

func main() {
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *flagVerbose

	log.Fatal(http.ListenAndServe(*flagListenAddr, proxy))
}
