package main

import (
	"flag"

	"uniport/api"
	"uniport/config"
	"uniport/utils"
)

func main() {
	var lDebug bool
	flag.BoolVar(&lDebug, "debug", false, "Debug flag forces no cache")
	flag.Parse()

	utils.Logger("")
	config.Init(nil) //Init Config.yaml
	config.Get().Debug = lDebug
	api.StartRouter()
}
