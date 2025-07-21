# fo

**fo** is a powerful command-line interface tool for competitive programming. It fetches sample cases from Codeforces, compiles C++ code, runs tests, and manages boilerplate for you.


## Configuration

`fo` looks for its configuration and templates in the standard config directory:

- **Config directory:**  
    **Linux**: `~/.config/fo` (e.g., `/home/ahmed/.config/fo`)

- **Config file:**  
  `config.yaml` — contains compiler commands, flags, source filename, executable name, etc.

- **C++ Template file:**  
  `template.cpp` — default C++ source template to use for new problems.


## Default Configuration Values

If no config file is present, `fo` will use these defaults internally:

```yaml
compiler:
  command: g++
  flags: "-Wall -Wextra -O2 -std=c++23"

sourceName: main.cpp
executableName: main
```


## Default C++ Template (`template.cpp`)

If no template file is found, `fo` will use the following built-in template:

```cpp
#include 
using namespace std;

typedef int i32;
typedef long long i64;
typedef unsigned int u32;
typedef unsigned long long u64;
typedef float f32;
typedef double f64;
typedef long double f80;
typedef vector vi;
typedef vector> vii;
typedef vector vl;
typedef vector> vll;
typedef pair pii;
typedef pair pll;
typedef pair psi;
typedef set si;
typedef map mii;
typedef unordered_map umii;

int main(void) {
  ios::sync_with_stdio(0);
  cin.tie(0);

  return 0;
}
```

## Installation

Make sure you have Go (1.18+) installed.

```bash
go install github.com/ahmedYasserM/fo@latest
```

`fo` is installed to your Go bin directory, add it to your PATH if necessary.

## Usage

```sh
fo [command]
```


## Available Commands

| Command | Description |
| :-- | :-- |
| `build` | Build your source (default `main.cpp`) using config settings  |
| `clean` | Removes generated files like `main` executable and `testcases.txt` |
| `copy` | Copies your source code (default: `main.cpp`) content to clipboard |
| `copy-clean` | Copies source code (default: `main.cpp`) content to clipboard after removing unused typedefs |
| `fetch` | Fetches sample test cases from a Codeforces problem URL |
| `run` | Builds (if needed) and runs the compiled program |
| `setup` | Sets up a new problem: fetches samples and creates the source file (default: `main.cpp`) if not exists |
| `test` | Run tests against sample inputs and outputs from `testcases.txt` |
| `completion` | Generate the autocompletion script for the specified shell |
| `help` | Help about any command |

## Flags

- `-h`, `--help`    Show help for `fo`


## Command-line Examples

### Build the C++ solution

```sh
fo build
```

**Quiet (suppress output):**

```sh
fo build --quiet
```


### Run the solution (auto-rebuilds if needed)

```sh
fo run
```

**Quiet (suppress build/run output):**

```sh
fo run --quiet
```


### Fetch sample test cases

```sh
fo fetch https://codeforces.com/contest/799/problem/A
```


### Test your solution against `testcases.txt`

```sh
fo test
```

**Quiet (suppress rebuild/test output):**

```sh
fo test --quiet
```


### Copy your solution to clipboard

```sh
fo copy
```


### Copy a cleaned solution (typeless) to clipboard

```sh
fo copy-clean
```


### Clean up generated files

```sh
fo clean
```


### Set up a new problem

```sh
fo setup https://codeforces.com/contest/799/problem/A
```


## Features

- Auto-detects if the source file (default: `main.cpp`) has changed and rebuilds automatically.
- Robust test parser for flexible `testcases.txt` format.
- Clipboard integration for code sharing.
- User-friendly colored output and error messages.
- `--quiet` flag for `build`, `run`, and `test` to suppress informational messages.
