package cmd

import (
	"fmt"
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fo",
	Short: "fo is a CLI tool for competitive programming workflows",
	Long: `fo is a powerful command-line interface tool designed to streamline
your competitive programming workflow. It handles fetching sample cases from
Codeforces, compiling C++ code, running tests, and managing boilerplate.`,
}

func Execute() {
	// Customize cobra output
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiYellow + cc.Bold,
		Commands:        cc.HiGreen + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.HiGreen + cc.Bold,
		Flags:           cc.HiBlue + cc.Bold,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
