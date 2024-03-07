//go:build ui

package main

import (
	"github.com/ArcticOJ/wrapper/v0"
	rice "github.com/GeertJohan/go.rice"
	"os"
)

const DefaultAvalancheOut = "avalanche/out"

func init() {
	PostInit = append(PostInit, func() {
		out := os.Getenv("AVALANCHE_OUT")
		if out == "" {
			out = DefaultAvalancheOut
		}
		box := rice.MustFindBox(out)
		wrapper.Register(Router, box)
	})
}
