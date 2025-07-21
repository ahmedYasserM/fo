package cmd

import (
	"fmt"
	"os"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/ahmedYasserM/fo/internal/utils"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copies source file (default: main.cpp) content to clipboard",
	Long:  `Reads the content of the source file (default: main.cpp) and copies it to the system clipboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.LoadConfigOnce(true); err != nil {
			fmt.Fprintf(os.Stderr, "%s❌  %v%s\n", colors.RED, err, colors.RESET)
			os.Exit(1)
		}

		if !utils.PathExists(utils.CmdConfig.SourceName) {
			fmt.Fprintf(os.Stderr, "%s❌ Error: %s not found.%s\n", colors.RED, utils.CmdConfig.SourceName, colors.RESET)
			os.Exit(1)
		}

		content, err := utils.ReadFileToString(utils.CmdConfig.SourceName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Error reading %s: %v%s\n", colors.RED, utils.CmdConfig.SourceName, err, colors.RESET)
			os.Exit(1)
		}

		err = utils.CopyToClipboard(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Error copying to clipboard: %v%s\n", colors.RED, err, colors.RESET)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
