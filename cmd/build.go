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
	Short: "Builds the C++ source file (default: main.cpp)",
	Long: `Builds the source file (default: main.cpp) using the compiler command specified in the configuration.
By default, it uses 'g++' with standard compilation flags.
The resulting executable will be named 'main'.

You can customize the compiler command and flags in your configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		if err := utils.LoadConfigOnce(false); err != nil {
			return err
		}

		err := utils.BuildExecutable(buildQuiet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ %v%s\n", colors.RED, err, colors.RESET)
			return err
		}
		fmt.Printf("%s✅ Build successful! Executable: %s%s\n", colors.GREEN, utils.CmdConfig.ExecutableName, colors.RESET)
		return nil
	},
}

func init() {
	buildCmd.Flags().BoolVarP(&buildQuiet, "quiet", "q", false, "Suppress build output")
	rootCmd.AddCommand(buildCmd)
}
