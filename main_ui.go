//go:build ui

package main

import (
	"blizzard/logger"
	"embed"
	"io/fs"
	"wrapper"
)

//go:embed all:avalanche/out/*
var AvalancheBundle embed.FS

func init() {
	sub, e := fs.Sub(AvalancheBundle, "avalanche/out")
	logger.Panic(e, "failed to open embedded avalanche bundle")
	wrapper.Register(Router, sub)
}
