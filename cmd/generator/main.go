package main

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "generator",
		Short: "generate a blizzard key",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(paseto.NewV4SymmetricKey().ExportHex())
		},
	}
	_ = root.Execute()
}
