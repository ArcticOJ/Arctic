//go:build ui

package main

import (
	"embed"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/ArcticOJ/wrapper/v0"
	"io/fs"
)

//go:embed all:avalanche/out/*
var AvalancheBundle embed.FS

func init() {
	sub, e := fs.Sub(AvalancheBundle, "avalanche/out")
	logger.Panic(e, "failed to open embedded avalanche bundle")
	wrapper.Register(Router, sub)
}
