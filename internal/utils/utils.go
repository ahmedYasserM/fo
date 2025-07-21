package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ahmedYasserM/fo/internal/colors"
)

// ExecuteCmd runs a shell command and prints its output.
// It now connects stdin, stdout, and stderr for interactive use.
func ExecuteCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin   // Connect child's stdin to parent's stdin
	cmd.Stdout = os.Stdout // Connect child's stdout to parent's stdout
	cmd.Stderr = os.Stderr // Connect child's stderr to parent's stderr
	return cmd.Run()
}

// ExecuteCmdWithInput runs a shell command with provided input and returns output.
func ExecuteCmdWithInput(input string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(input)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %s %v, stderr: %s, stdout: %s, error: %w", name, args, errb.String(), outb.String(), err)
	}
	return outb.String(), nil
}

// CopyToClipboard copies the given content to the system clipboard.
func CopyToClipboard(content string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try wl-copy first for Wayland, then xclip for X11, then fallback to xsel
		_, errWl := exec.LookPath("wl-copy")
		_, errXclip := exec.LookPath("xclip")
		_, errXsel := exec.LookPath("xsel")

		if errWl == nil {
			cmd = exec.Command("wl-copy")
		} else if errXclip == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if errXsel == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard utility found (wl-copy, xclip, or xsel)")
		}
	default:
		return fmt.Errorf("unsupported operating system for clipboard operation: %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(content)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}
	fmt.Println(colors.GREEN + "✅ Content copied to clipboard!" + colors.RESET)
	return nil
}

// ReadFileToString reads a file's content into a string.
func ReadFileToString(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filepath, err)
	}
	return string(content), nil
}

// WriteStringToFile writes a string to a file.
func WriteStringToFile(filepath, content string) error {
	return os.WriteFile(filepath, []byte(content), 0o644)
}

// PathExists checks if a file or directory exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CppTemplate is the embedded C++ template
const CppTemplate = `#include <bits/stdc++.h>
using namespace std;

typedef int i32;
typedef long long i64;
typedef unsigned int u32;
typedef unsigned long long u64;
typedef float f32;
typedef double f64;
typedef long double f80;
typedef vector<int> vi;
typedef vector<vector<int>> vii;
typedef vector<long long> vl;
typedef vector<vector<long long>> vll;
typedef pair<int, int> pii;
typedef pair<long long, long long> pll;
typedef pair<string, int> psi;
typedef set<int> si;
typedef map<int, int> mii;
typedef unordered_map<int, int> umii;

int main(void) {
  ios::sync_with_stdio(0);
  cin.tie(0);

  return 0;
}
`

// TypedefRegex is used for identifying typedef lines for clip-clean.
// Note: This regex is specific to the example C++ template's typedefs.
const TypedefRegex = `(?m)^typedef\s+.*\s+([a-zA-Z0-9_]+);$`
