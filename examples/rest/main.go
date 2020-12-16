package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/rest"
	"github.com/patractlabs/go-patract/utils"
)

var (
	path      = flag.String("path", ":8899", "path to url")
	isOffline = flag.Bool("offline", false, "is use offline mode")
)

func main() {
	flag.Parse()

	router := gin.Default()
	r := rest.NewRouter(router)

	r.WithRuntimeMetadata(utils.LoadRuntimeMetadata("./test/metadata.json"))

	erc20, err := metadata.NewFromFile("./test/contracts/ink/erc20.json")
	if err != nil {
		panic(err)
	}

	erc721, err := metadata.NewFromFile("./test/contracts/ink/erc721.json")
	if err != nil {
		panic(err)
	}

	flipper, err := metadata.NewFromFile("./test/contracts/ink/flipper.json")
	if err != nil {
		panic(err)
	}

	r.WithMetaData(erc20)
	r.WithMetaData(erc721)
	r.WithMetaData(flipper)

	r.Init()

	rest.NewRouter(router).Run(*path)
}
