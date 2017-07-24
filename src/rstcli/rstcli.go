package main

import (
	"flag"
	"fmt"
	"os"
	"octlink/rstore/src/utils"
)

var (
	url  string
	cert string
	cmd  string
	imageid string
	output string
)

func init() {
	flag.StringVar(&cmd, "cmd", "", "cmd [pull,push,add,export,import,images,search]")
	flag.StringVar(&url, "url", "", "image url")
	flag.StringVar(&cert, "cert", "./cert/cert.crt", "cert file path")
	flag.StringVar(&imageid, "imageid", "", "image uuid like xxxxxxxxxxxxxxxxxxxxxx")
	flag.StringVar(&output, "output", "root.qcow2", "image's filepath output to local")
}

func usage() {
	fmt.Println("  RVM Store Client of V" + utils.Version() + "\n")
	fmt.Println("  ./rstcli -cmd pull -out xxx.qcow2 -imageid xxxxxxxxxxxxx")
	fmt.Println("  ./rstcli -cmd search -image ssss\n")
	flag.PrintDefaults()
}

func pull_image() {
	fmt.Printf("Running in image pulling function %s\n", imageid)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if len(os.Args) <= 2 || cmd == "" {
		usage()
		os.Exit(0)
	}

	/*fmt.Printf("cmd is %s\n", cmd)
	fmt.Printf("url is %s\n", url)
	fmt.Printf("cert is %s\n", cert)
	*/

	if cmd == "pull" {
		if imageid == "" {
			fmt.Printf("image id must specified for PULL action\n")
			os.Exit(1)
		}
		pull_image()
	}
}
