package cmd

import (
	"fmt"
	"os"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/ahmedYasserM/fo/internal/utils"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [URL]",
	Short: "Sets up a new problem: fetches samples and creates main.cpp if not exists",
	Long: `This command streamlines the setup for a new competitive programming problem.
It first fetches sample test cases from the provided Codeforces URL using 'fo fetch'.
Then, if 'main.cpp' does not already exist in the current directory, it creates it
with a standard C++ template.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Run fetch command
		fetchURL := args[0]
		fmt.Printf("Running %sfetch %s%s...\n", colors.CYAN, fetchURL, colors.RESET)
		if err := fetchCmd.RunE(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Setup failed during fetch: %v%s\n", colors.RED, err, colors.RESET)
			os.Exit(1)
		}

		// 2. Create main.cpp if it doesn't exist
		// if !utils.PathExists("main.cpp") {
		fmt.Printf("Creating template %smain.cpp%s...\n", colors.CYAN, colors.RESET)
		err := utils.WriteStringToFile("main.cpp", utils.CppTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Failed to create main.cpp: %v%s\n", colors.RED, err, colors.RESET)
			os.Exit(1)
		}
		fmt.Printf("%s✅ Created template main.cpp%s\n", colors.GREEN, colors.RESET)
		// } else {
		// fmt.Printf("%smain.cpp%s already exists, skipping template creation.\n", colors.CYAN, colors.RESET)
		// }
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
