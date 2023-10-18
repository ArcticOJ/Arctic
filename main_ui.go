//go:build ui

package main

import (
	"embed"
	"github.com/ArcticOJ/wrapper/v0"
	"io/fs"
)

//go:embed all:avalanche/out/*
var AvalancheBundle embed.FS

func init() {
	sub, e := fs.Sub(AvalancheBundle, "avalanche/out")
	if e != nil {
		panic(e)
	}
	wrapper.Register(Router, sub)
}
