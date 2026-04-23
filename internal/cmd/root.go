package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goploy",
	Short: "Heroku-lite deployment CLI",
	Long:  `Goploy is a CLI tool to deploy and manage apps on your VPS easily.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
