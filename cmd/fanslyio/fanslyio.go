package main

import (
	"flag"
	"fmt"

	"github.com/spaam/fanslyio/internal/config"
	"github.com/spaam/fanslyio/internal/context"
	"github.com/spaam/fanslyio/internal/crawl"
	"github.com/spaam/fanslyio/internal/download"
)

var version = "dev"

func main() {
	var versionflag bool
	flag.BoolVar(&versionflag, "version", false, "display version number")
	flag.Parse()
	if versionflag {
		fmt.Printf("fanslyio v%s\n", version)
		return
	}
	file := config.Configfile()
	cc, err := context.ParseConfig(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	content, err := crawl.Crawl(&cc)
	if err != nil {
		fmt.Println(err)
		return
	}
	download.Download(&cc, content)
}
