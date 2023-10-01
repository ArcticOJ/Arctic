//go:build !(headless || ui)

package main

func main() {
	panic("nothing to run!")
}
