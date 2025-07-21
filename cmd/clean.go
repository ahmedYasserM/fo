package cmd

import (
	"fmt"
	"os"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/ahmedYasserM/fo/internal/utils"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes generated files like 'main' executable and 'testcases.txt'",
	Long:  `Removes the 'main' executable and the 'testcases.txt' file if they exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Remove main executable
		if utils.PathExists("main") {
			err := os.Remove("main")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s❌ Error removing 'main' executable: %v%s\n", colors.RED, err, colors.RESET)
				os.Exit(1)
			}
			fmt.Printf("%s✅ Removed 'main' executable.%s\n", colors.GREEN, colors.RESET)
		} else {
			fmt.Printf("%s'main' executable%s not found, nothing to clean.%s\n", colors.YELLOW, colors.RESET, colors.RESET)
		}

		// Remove testcases.txt
		if utils.PathExists("testcases.txt") {
			err := os.Remove("testcases.txt")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s❌ Error removing 'testcases.txt': %v%s\n", colors.RED, err, colors.RESET)
				os.Exit(1)
			}
			fmt.Printf("%s✅ Removed 'testcases.txt'.%s\n", colors.GREEN, colors.RESET)
		} else {
			fmt.Printf("%s'testcases.txt'%s not found, nothing to clean.%s\n", colors.YELLOW, colors.RESET, colors.RESET)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
