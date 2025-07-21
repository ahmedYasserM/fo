package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/ahmedYasserM/fo/internal/utils"

	"github.com/spf13/cobra"
)

type Testcase struct {
	Input    string
	Expected string
}

var testQuiet bool

// parseTestcases reads testcases.txt and extracts input/output samples.
// It tolerates blank lines and flexible formatting.
func parseTestcases(filename string) ([]Testcase, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}
	defer file.Close()

	var tests []Testcase
	scanner := bufio.NewScanner(file)
	var state string // "input", "output", or ""
	var inputLines []string
	var outputLines []string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "--- Sample") && strings.Contains(trimmed, "Input ---") {
			if state != "" {
				// Save previous test before switching
				tests = append(tests, Testcase{
					Input:    strings.Join(inputLines, "\n"),
					Expected: strings.Join(outputLines, "\n"),
				})
				inputLines = nil
				outputLines = nil
			}
			state = "input"
			continue
		}
		if strings.HasPrefix(trimmed, "--- Sample") && strings.Contains(trimmed, "Output ---") {
			state = "output"
			continue
		}

		// Accumulate lines according to current state
		switch state {
		case "input":
			inputLines = append(inputLines, line)
		case "output":
			outputLines = append(outputLines, line)
		default:
			// Outside a sample block, ignore lines
		}
	}
	// Add last sample after EOF
	if len(inputLines) > 0 || len(outputLines) > 0 {
		tests = append(tests, Testcase{
			Input:    strings.Join(inputLines, "\n"),
			Expected: strings.Join(outputLines, "\n"),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %w", filename, err)
	}

	return tests, nil
}

// ensureBuilt recompiles main.cpp if missing or outdated
func ensureBuilt(quiet bool) error {
	mainCpp := "main.cpp"
	mainExe := "main"

	if !utils.PathExists(mainCpp) {
		return fmt.Errorf("%smain.cpp not found.%s", colors.RED, colors.RESET)
	}

	needsBuild := false

	cppInfo, err := os.Stat(mainCpp)
	if err != nil {
		return err
	}

	exeInfo, err := os.Stat(mainExe)
	if os.IsNotExist(err) {
		needsBuild = true
	} else if err != nil {
		return err
	} else if cppInfo.ModTime().After(exeInfo.ModTime()) {
		needsBuild = true
	}

	if needsBuild {
		if !quiet {
			fmt.Printf("%smain.cpp changed or executable missing. Rebuilding...%s\n", colors.YELLOW, colors.RESET)
		}
		if err := utils.BuildExecutable(testQuiet); err != nil {
			return fmt.Errorf("%sBuild failed: %w%s", colors.RED, err, colors.RESET)
		}
		if !quiet {
			fmt.Printf("%sBuild succeeded.%s\n", colors.GREEN, colors.RESET)
		}
	} else {
		if !quiet {
			fmt.Printf("%sExecutable up-to-date. Skipping build.%s\n", colors.CYAN, colors.RESET)
		}
	}
	return nil
}

// testCmd runs tests from testcases.txt by feeding inputs to the program and comparing output
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests against sample inputs and outputs from testcases.txt",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Step 1. Ensure executable is up to date
		if err := ensureBuilt(testQuiet); err != nil {
			return err
		}

		// Step 2. Parse test cases
		tests, err := parseTestcases("testcases.txt")
		if err != nil {
			return fmt.Errorf("error parsing testcases.txt: %w", err)
		}
		if len(tests) == 0 {
			return fmt.Errorf("no tests found in testcases.txt")
		}

		// Step 3. Run each test
		fmt.Printf("%sRunning tests...%s\n", colors.CYAN, colors.RESET)

		passed := 0
		for i, test := range tests {
			actual, err := executeWithInput("./main", test.Input)
			if err != nil {
				fmt.Printf("%sTest #%d execution error: %v%s\n", colors.RED, i+1, err, colors.RESET)
				continue
			}
			actual = strings.TrimSpace(actual)

			if actual == strings.TrimSpace(test.Expected) {
				fmt.Printf("%s=== Test %d === %s[OK]%s\n", colors.BOLD, i+1, colors.GREEN, colors.RESET)
				passed++
			} else {
				fmt.Printf("%s=== Test %d === %s[FAIL]%s\n", colors.BOLD, i+1, colors.RED, colors.RESET)
				fmt.Printf("%sInput:%s\n%s\n", colors.YELLOW, colors.RESET, strings.TrimSpace(test.Input))
				fmt.Printf("%sYour output:%s\n%s\n", colors.YELLOW, colors.RESET, strings.TrimSpace(actual))
				fmt.Printf("%sExpected:%s\n%s\n\n", colors.YELLOW, colors.RESET, strings.TrimSpace(test.Expected))
			}

		}

		if passed == len(tests) {
			fmt.Printf("%s✅ Test summary: Passed %d out of %d tests.%s\n", colors.BOLD+colors.CYAN, passed, len(tests), colors.RESET)
		} else {
			fmt.Printf("%s❌ Test summary: Passed %d out of %d tests.%s\n", colors.BOLD+colors.CYAN, passed, len(tests), colors.RESET)
		}
		return nil
	},
}

// executeWithInput runs a command feeding its stdin and returns the stdout output (or error)
func executeWithInput(command, input string) (string, error) {
	cmd := exec.Command(command)
	cmd.Stdin = strings.NewReader(input)

	outBytes, err := cmd.Output()
	if err != nil {
		// Try to get stderr text in case of exec error
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("program error: %s", string(exitErr.Stderr))
		}
		return "", err
	}
	return string(outBytes), nil
}

func init() {
	testCmd.Flags().BoolVarP(&testQuiet, "quiet", "q", false, "Suppress build output during tests")

	rootCmd.AddCommand(testCmd)
}
