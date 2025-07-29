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

source_name: main.cpp
executable_name: main
```


## Default C++ Template (`template.cpp`)

If no template file is found, `fo` will use the following built-in template:

```cpp
#include <bits/stdc++.h>
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
| `setup` | Sets up a new problem: fetches samples and creates the source file (default: `main.cpp`) if not exists |
| `test` | Run tests against sample inputs and outputs from `testcases.txt` |
| `copy-clean` | Copies source code (default: `main.cpp`) content to clipboard after removing unused typedefs |
| `copy` | Copies your source code (default: `main.cpp`) content to clipboard |
| `fetch` | Fetches sample test cases from a Codeforces problem URL |
| `build` | Build your source (default `main.cpp`) using config settings  |
| `run` | Builds (if needed) and runs the compiled program |
| `clean` | Removes generated files like `main` executable and `testcases.txt` |
| `completion` | Generate the autocompletion script for the specified shell |
| `help` | Help about any command |

## Flags

- `-h`, `--help`    Show help for `fo`


## Command-line Examples

### Set up a new problem

```sh
fo setup https://codeforces.com/contest/799/problem/A
```

### Test your solution against `testcases.txt`

```sh
fo test
```

**Quiet (suppress rebuild/test output):**

```sh
fo test --quiet
```

### Copy a cleaned solution (typeless) to clipboard

```sh
fo copy-clean
```

### Copy your solution to clipboard

```sh
fo copy
```

### Fetch sample test cases

```sh
fo fetch https://codeforces.com/contest/799/problem/A
```

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


### Clean up generated files

```sh
fo clean
```

## Features

- Auto-detects if the source file (default: `main.cpp` has changed and rebuilds automatically.
- Robust test parser for flexible `testcases.txt` format.
- Clipboard integration for code sharing.
- User-friendly colored output and error messages.
- `--quiet` flag for `build`, `run`, and `test` to suppress informational messages.
