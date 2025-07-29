package cmd

import (
	"fmt"
	"os"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/ahmedYasserM/fo/internal/utils"

	"github.com/spf13/cobra"
)

func preload() error {

	// Load config file
	if err := utils.LoadConfigOnce(false); err != nil {
		return fmt.Errorf("%s❌ %v%s\n", colors.RED, err, colors.RESET)
	}

	// Load C++ template
	if err := utils.LoadTemplateOnce(); err != nil {
		return fmt.Errorf("%s❌  %v%s\n", colors.RED, err, colors.RESET)
	}

	return nil
}

var setupCmd = &cobra.Command{
	Use:   "setup [URL]",
	Short: "Sets up a new problem: fetches samples and creates source file if not exists",
	Long: `This command streamlines the setup for a new competitive programming problem.
It fetches sample test cases from the provided Codeforces URL using 'fo fetch'.
Then, if the source file does not already exist in the current directory, it creates it
using the C++ template located in the configuration directory, with the filename
determined by your configuration (default is 'main.cpp').

You can customize the executable/file name in your config file, which affects which
source file is created and used.

This helps keep your workflow flexible and consistent across projects.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if err := preload(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fetchURL := args[0]
		fmt.Printf("Running %sfetch %s%s...\n", colors.CYAN, fetchURL, colors.RESET)
		if err := fetchCmd.RunE(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ %v%s\n", colors.RED, err, colors.RESET)
			os.Exit(1)
		}

		fmt.Printf("Creating template %s%s%s...\n", colors.CYAN, utils.CmdConfig.SourceName, colors.RESET)
		err := utils.WriteStringToFile(utils.CmdConfig.SourceName, utils.CmdTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Failed to create %s: %v%s\n", colors.RED, utils.CmdConfig.SourceName, err, colors.RESET)
			os.Exit(1)
		}
		fmt.Printf("%s✅ Created template %s%s\n", colors.GREEN, utils.CmdConfig.SourceName, colors.RESET)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
