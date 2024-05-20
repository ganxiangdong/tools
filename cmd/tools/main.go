package main

import (
	"codeup.aliyun.com/5f9118049cffa29cfdd3be1c/tools/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A generator for Cobra based Applications",
	Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func main() {
	rootCmd.AddCommand(internal.ModelCmd)
	rootCmd.AddCommand(internal.WireCmd)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
