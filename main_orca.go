//go:build orca && headless

package main

import "github.com/ArcticOJ/orca/v0"

func init() {
	LateInit = append(LateInit, orca.Init)
}
