package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ahmedYasserM/fo/internal/colors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Compiler struct {
		Command string `yaml:"command"`
		Flags   string `yaml:"flags"`
	} `yaml:"compiler"`
	SourceName     string `yaml:"source_name"`
	ExecutableName string `yaml:"executable_name"`
}

var (
	configDir      string
	CmdConfig      *Config
	CmdTemplate    string
	templateLoaded bool
	defaultConfig  = Config{
		Compiler: struct {
			Command string `yaml:"command"`
			Flags   string `yaml:"flags"`
		}{
			Command: "g++",
			Flags:   "-Wall -Wextra -O2 -std=c++23",
		},
		SourceName:     "main.cpp",
		ExecutableName: "main",
	}
)

const defaultCppTemplate = `#include <bits/stdc++.h>
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
  ios::sync_with_stdio(false);
  cin.tie(nullptr);

  return 0;
}
`

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	configDir = filepath.Join(homeDir, ".config", "fo")
}

func parseConfig(quiet bool) error {
	configPath := filepath.Join(configDir, "config.yaml")

	if !PathExists(configPath) {
		fmt.Fprintf(os.Stderr, "%s⚠️ Config file not found. Using defaults.%s\n", colors.YELLOW, colors.RESET)
		CmdConfig = &defaultConfig
		return nil
	}

	data, err := ReadFileToBytes(configPath)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, &CmdConfig); err != nil {
		return err
	}

	if !quiet {
		fmt.Printf("%s✅ Config loaded successfully! %s\n", colors.GREEN, colors.RESET)
	}

	return nil
}

func parseTemplate() error {
	templatePath := filepath.Join(configDir, "template.cpp")

	if !PathExists(templatePath) {
		fmt.Fprintf(os.Stderr, "%s⚠️ Template file not found. Using default template.%s\n", colors.YELLOW, colors.RESET)
		CmdTemplate = defaultCppTemplate
		return nil
	}

	content, err := ReadFileToString(templatePath)
	if err != nil {
		return err
	}

	CmdTemplate = content
	templateLoaded = true
	fmt.Printf("%s✅ C++ template loaded successfully! %s\n", colors.GREEN, colors.RESET)

	return nil
}

func LoadConfigOnce(quiet bool) error {
	if CmdConfig != nil {
		return nil
	}

	return parseConfig(quiet)
}

func LoadTemplateOnce() error {
	if templateLoaded {
		return nil
	}

	return parseTemplate()
}
