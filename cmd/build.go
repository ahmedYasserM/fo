package cmd

import (
	"fmt"
	"os"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/ahmedYasserM/fo/internal/utils"

	"github.com/spf13/cobra"
)

var buildQuiet bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the C++ source file (main.cpp)",
	Long: `Builds the main.cpp file using g++.
The executable will be named 'main'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := utils.BuildExecutable(buildQuiet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Build failed: %v%s\n", colors.RED, err, colors.RESET)
			return err
		}
		fmt.Printf("%s✅ Build successful! Executable: main%s\n", colors.GREEN, colors.RESET)
		return nil
	},
}

func init() {
	buildCmd.Flags().BoolVarP(&buildQuiet, "quiet", "q", false, "Suppress build output")
	rootCmd.AddCommand(buildCmd)
}
