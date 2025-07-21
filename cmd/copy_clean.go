package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ahmedYasserM/fo/internal/colors"
	"github.com/ahmedYasserM/fo/internal/utils"

	"github.com/spf13/cobra"
)

var copyCleanCmd = &cobra.Command{
	Use:   "copy-clean",
	Short: "Copies main.cpp content to clipboard after removing unused typedefs",
	Long: `Reads the content of main.cpp, analyzes which typedefs from the standard template
are actually used in the code, removes the unused ones, and then copies the
cleaned content to the system clipboard.`,
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

		cleanedContent := cleanTypedefs(content)

		err = utils.CopyToClipboard(cleanedContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Error copying to clipboard: %v%s\n", colors.RED, err, colors.RESET)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCleanCmd)
}

// cleanTypedefs processes the C++ code to remove unused typedefs.
func cleanTypedefs(code string) string {
	lines := strings.Split(code, "\n")
	typedefs := make(map[string]string) // alias -> full_typedef_line

	// Regex to capture typedef alias
	re := regexp.MustCompile(utils.TypedefRegex)

	// Identify typedef lines and their aliases
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			alias := matches[1]
			typedefs[alias] = line
		}
	}

	// Check usage of each typedef alias
	unusedAliases := make(map[string]bool)
	for alias := range typedefs {
		// Count occurrences of the alias in the entire code
		// Excluding its own typedef declaration line
		count := 0
		for _, line := range lines {
			if strings.Contains(line, typedefs[alias]) { // Skip the declaration line itself
				continue
			}
			// Use word boundaries to avoid matching substrings (e.g., 'i32' in 'myi32var')
			if regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(alias))).MatchString(line) {
				count++
			}
		}
		if count == 0 {
			unusedAliases[alias] = true
		}
	}

	// Reconstruct code without unused typedef lines
	var cleanedLines []string
	for _, line := range lines {
		isTypedefLine := false
		for alias, typedefLine := range typedefs {
			if line == typedefLine {
				if unusedAliases[alias] {
					isTypedefLine = true // This typedef line is unused, so skip it
					break
				}
			}
		}
		if !isTypedefLine {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}
