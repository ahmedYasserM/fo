package cmd

import (
	"fmt"
	"os"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/ahmedYasserM/fo/internal/utils"

	"github.com/spf13/cobra"
)

var runQuiet bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Builds (if needed) and runs the compiled program",
	RunE: func(cmd *cobra.Command, args []string) error {
		mainCppPath := "main.cpp"
		mainExecPath := "main"

		if !utils.PathExists(mainCppPath) {
			return fmt.Errorf("%sError: %s%s not found. Cannot compile or run.%s", colors.RED, colors.BOLD, mainCppPath, colors.RESET)
		}

		needsBuild := false
		mainCppInfo, err := os.Stat(mainCppPath)
		if err != nil {
			return fmt.Errorf("could not get info for %s: %w", mainCppPath, err)
		}

		mainExecInfo, err := os.Stat(mainExecPath)
		if os.IsNotExist(err) {
			needsBuild = true
		} else if err != nil {
			return fmt.Errorf("could not get info for %s: %w", mainExecPath, err)
		} else {
			if mainCppInfo.ModTime().After(mainExecInfo.ModTime()) {
				needsBuild = true
			}
		}

		if needsBuild {
			if !runQuiet {
				fmt.Printf("%sExecutable '%s' is missing or outdated. Building...%s\n", colors.YELLOW, mainExecPath, colors.RESET)
			}
			if err := utils.BuildExecutable(runQuiet); err != nil {
				return fmt.Errorf("%sBuild failed, cannot run:%s %w", colors.RED, colors.RESET, err)
			}
			fmt.Printf("%sBuild successful!%s\n", colors.GREEN, colors.RESET)
		} else {
			if !runQuiet {
				fmt.Printf("%s'%s' is up to date, skipping build.%s\n", colors.CYAN, mainExecPath, colors.RESET)
			}
		}

		fmt.Printf("%sRunning '%s'...%s\n", colors.CYAN, mainExecPath, colors.RESET)
		err = utils.ExecuteCmd(fmt.Sprintf("./%s", mainExecPath))
		if err != nil {
			return fmt.Errorf("%sProgram exited with error:%s %w", colors.RED, colors.RESET, err)
		}

		return nil
	},
}

func init() {
	runCmd.Flags().BoolVarP(&runQuiet, "quiet", "q", false, "Suppress build output")
	rootCmd.AddCommand(runCmd)
}
