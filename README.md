# fo

**fo** is a powerful command-line interface tool for competitive programming. It fetches sample cases from Codeforces, compiles C++ code, runs tests, and manages boilerplate for you.

## Usage

```sh
fo [command]
```


## Available Commands

| Command | Description |
| :-- | :-- |
| `build` | Builds the C++ source file (`main.cpp`) |
| `clean` | Removes generated files like `main` executable and `testcases.txt` |
| `copy` | Copies `main.cpp` content to clipboard |
| `copy-clean` | Copies `main.cpp` content to clipboard after removing unused typedefs |
| `fetch` | Fetches sample test cases from a Codeforces problem URL |
| `run` | Builds (if needed) and runs the compiled program |
| `setup` | Sets up a new problem: fetches samples and creates `main.cpp` if not exists |
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

- Auto-detects if `main.cpp` has changed and rebuilds automatically.
- Robust test parser for flexible `testcases.txt` format.
- Clipboard integration for code sharing.
- User-friendly colored output and error messages.
- `--quiet` flag for `build`, `run`, and `test` to suppress informational messages.
