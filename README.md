# Deploy

A lightweight, configurable CLI tool for building, organizing, and deploying local binaries.

## Overview

**Deploy** is a simple automation tool that compiles source code and moves the resulting binary (and optionally the source) into structured directories based on a configuration file.

It's designed to replace repetitive shell scripts with a consistent, reusable workflow.

## Features

* Compile projects using a custom command
* Automatically move or copy binaries after build
* Support for global and local binary directories
* Optional organization of source files
* Config-driven behavior via a `.deployfile`

## How It Works

1. Reads configuration from `.deployfile`
2. Executes the specified compilation command
3. Copies or moves the resulting binary
4. Optionally organizes source files

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/Eth3rna1/deploy.git
cd deploy
go build
```

Optionally move it to a directory in your `$PATH`.

## Usage

```bash
deploy <source>
```

Example:

```bash
deploy main.go
```

## Configuration

Deploy uses a `.deployfile` in the project root.

If one does not exist, it will be automatically created and initialized.

### Importing Another `.deployfile`

You can import variables from another `.deployfile` using:
```text
DEPLOYFILE_LOC=<path-to-other-deployfile>
```

If `DEPLOYFILE_LOC` is defined, all variables from the referenced file will be loaded and **override any previously defined variables**.

This allows reuse of shared configurations across multiple projects.

**Important:**
- Variables defined in the imported file take precedence
- Any variables defined *after* the import may be overwritten

### Boilerplate `.deployfile`

```text
# Use this variable to import another deployfile
DEPLOYFILE_LOC=

# if 'DEPLOYFILE_LOC' is defined, then all defined variables
# within such deployfile will get imported and overwrite any
# previously defined variables if defined in such deployfile.
# Everything defined past this comment will get overwritten
# or defined.

# Location of the global binary directory
GLOBAL_BIN_DIR=

# Location of the local binary directory
LOCAL_BIN_DIR=

# Location of the local scripts directory
SCRIPTS_DIR=

# Location of the base directory of the project
BASE_DIR=

# The CMD command to compile the project
COMPILATION_CMD=

# Tells where the binary is located after compiling
BINARY_LOC=
```

### Required Variables

| Key               | Description                             |
| ----------------- | --------------------------------------- |
| `BASE_DIR`        | Directory where compilation is executed |
| `COMPILATION_CMD` | Command used to build the project       |
| `BINARY_LOC`      | Path to the compiled binary             |

### Optional Variables

| Key              | Description                            |
| ---------------- | -------------------------------------- |
| `GLOBAL_BIN_DIR` | Directory to copy the binary into      |
| `LOCAL_BIN_DIR`  | Directory to move the binary into      |
| `SCRIPTS_DIR`    | Directory to move the source file into |

### Example `.deployfile`

```
BASE_DIR=./
COMPILATION_CMD=go build -o myapp
BINARY_LOC=./myapp
GLOBAL_BIN_DIR=/usr/local/bin
LOCAL_BIN_DIR=./bin
SCRIPTS_DIR=./scripts
```

## Example Workflow

```bash
deploy main.go
```

* Compiles the project
* Copies the binary to `/usr/local/bin`
* Moves the binary to `./bin`
* Moves `main.go` to `./scripts`

## Limitations

* No dependency management
* No versioning support
* No remote package sources

This tool is focused purely on local build and deployment workflows.

## Future Improvements

* Better validation for `.deployfile`
* Logging and verbose/debug modes

## Contributing

Contributions are welcome. Feel free to open issues or submit pull requests.
