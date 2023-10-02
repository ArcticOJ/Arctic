//go:build orca && headless

package main

import "orca"

func init() {
	OnInit = append(OnInit, orca.Init)
}
