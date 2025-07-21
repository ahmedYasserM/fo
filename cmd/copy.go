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
	Short: "Copies main.cpp content to clipboard",
	Long:  `Reads the content of main.cpp and copies it to the system clipboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.PathExists("main.cpp") {
			fmt.Fprintf(os.Stderr, "%s❌ Error: main.cpp not found.%s\n", colors.RED, colors.RESET)
			os.Exit(1)
		}

		content, err := utils.ReadFileToString("main.cpp")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Error reading main.cpp: %v%s\n", colors.RED, err, colors.RESET)
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
