package main

import (
	"github.com/frostgrove/vv/utils/vvgoose"

	"github.com/frostgrove/vv-template/src/config"
)

func main() {
	configuration, _ := config.MustLoad()
	vvgoose.Execute(&configuration.Infra.Db)
}
