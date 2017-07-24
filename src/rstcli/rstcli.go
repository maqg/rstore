package main

import (
	"flag"
	"fmt"
	"octlink/rstore/src/utils"
)

var (
	url  string
	cert string
	cmd  bool
)

func init() {
	flag.Var(&cmd, "cmd", "cmd [pull,push,add,export,import,images,search]")
	flag.StringVar(&url, "url", "", "image url")
	flag.StringVar(&cert, "cert", "./cert/cert.crt", "cert file path")
}

func usage() {
	fmt.Println("RVM Store Client of V" + utils.Version())
	flag.PrintDefaults()
}

func main() {

	flag.Usage = usage
	flag.Parse()

	fmt.Printf("cmd is %s\n", cmd)
	fmt.Printf("url is %s\n", url)
	fmt.Printf("cert is %s\n", cert)
}
