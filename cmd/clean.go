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
	Long:  `Removes the executable file and the 'testcases.txt' if they exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.LoadConfigOnce(true); err != nil {
			fmt.Fprintf(os.Stderr, "%s❌  %v%s\n", colors.RED, err, colors.RESET)
			os.Exit(1)
		}

		// Remove executable file
		if utils.PathExists(utils.CmdConfig.ExecutableName) {
			err := os.Remove(utils.CmdConfig.ExecutableName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s❌ Error removing '%s' executable: %v%s\n", colors.RED, utils.CmdConfig.ExecutableName, err, colors.RESET)
				os.Exit(1)
			}
			fmt.Printf("%s✅ Removed '%s' executable.%s\n", colors.GREEN, utils.CmdConfig.ExecutableName, colors.RESET)
		} else {
			fmt.Printf("%s'%s' executable%s not found, nothing to clean.%s\n", colors.YELLOW, utils.CmdConfig.ExecutableName, colors.RESET, colors.RESET)
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
