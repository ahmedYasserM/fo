package utils

import (
	"fmt"
	"strings"

	"github.com/ahmedYasserM/fo/internal/colors"
)

// buildExecutable encapsulates the C++ build logic.
// It returns an error if the build fails.
func BuildExecutable(quiet bool) error {
	args := strings.Fields(CmdConfig.Compiler.Flags)
	args = append(args, CmdConfig.SourceName, "-o", CmdConfig.ExecutableName)

	if !PathExists(CmdConfig.SourceName) {
		return fmt.Errorf("%s not found. Cannot compile.", CmdConfig.SourceName)
	}

	if !quiet {
		fmt.Printf("Compiling %s%s%s...\n", colors.CYAN, CmdConfig.SourceName, colors.RESET)
	}
	err := ExecuteCmd(CmdConfig.Compiler.Command, args...)
	if err != nil {
		return fmt.Errorf("%s command failed: %w", CmdConfig.Compiler.Command, err)
	}
	return nil
}
