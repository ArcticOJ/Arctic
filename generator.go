//go:build headless

package main

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate a secret key for client sessions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(paseto.NewV4SymmetricKey().ExportHex())
	},
}
