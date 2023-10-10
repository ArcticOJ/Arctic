//go:build orca && headless

package main

import "github.com/ArcticOJ/orca/v0"

func init() {
	OnInit = append(OnInit, orca.Init)
}
