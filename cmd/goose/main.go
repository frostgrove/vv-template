package main

import (
	"github.com/frostgrove/vv-template/src/config"
	"github.com/frostgrove/vv/utils/vvcfg"
	"github.com/frostgrove/vv/utils/vvgoose"
)

func main() {
	cfg := vvcfg.MustLoad[config.Config]()
	vvgoose.Execute(cfg.Infra.Db)
}
