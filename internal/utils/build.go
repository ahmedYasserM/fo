package utils

import (
	"fmt"

	"github.com/ahmedYasserM/fo/internal/colors"
)

// buildExecutable encapsulates the C++ build logic.
// It returns an error if the build fails.
func BuildExecutable(quiet bool) error {
	mainCppPath := "main.cpp"
	mainExecPath := "main"

	if !PathExists(mainCppPath) {
		return fmt.Errorf("%s not found. Cannot compile.", mainCppPath)
	}

	if !quiet {
		fmt.Printf("Compiling %s%s%s...\n", colors.CYAN, mainCppPath, colors.RESET)
	}
	err := ExecuteCmd("g++", "-Wall", "-Wextra", "-O2", "-std=c++23", "-o", mainExecPath, mainCppPath)
	if err != nil {
		return fmt.Errorf("g++ command failed: %w", err)
	}
	return nil
}
