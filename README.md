# maclnr

`maclnr` is a command-line tool written in Go that helps you manage directories on macOS by listing and cleaning files. It includes functionality to list files by size and remove `.DS_Store` files.

## Installation

To install `maclnr`, clone the repository and build the project using Go:

```sh
git clone https://github.com/yourusername/maclnr.git
cd maclnr
go build -o maclnr
```

## Usage

### List Files

The `list` command lists all files in a directory recursively, ordered by size from largest to smallest. You can also filter files based on a minimum size.

#### List all files:

```sh
./maclnr list --dir /path/to/directory
```

#### List files larger than a specified size (in bytes):

```sh
./maclnr list --dir /path/to/directory --min-size 1048576
```

### Clean Directory

The `clean` command removes `.DS_Store` files and files larger than a specified size from a directory. It includes options for a dry run, verbose output, and confirmation prompts.

#### Clean `.DS_Store` files:

```sh
./maclnr clean --dir /path/to/directory --ds-store
```

#### Clean files larger than a specified size (in bytes):

```sh
./maclnr clean --dir /path/to/directory --min-size 1048576
```

#### Clean files with verbose output:

```sh
./maclnr clean --dir /path/to/directory --verbose
```

#### Perform a dry run (no files will be deleted):

```sh
./maclnr clean --dir /path/to/directory --dry-run
```

#### Clean files without confirmation:

```sh
./maclnr clean --dir /path/to/directory --confirm
```

## Flags

- `--dir` (`-d`): The directory to list or clean files.
- `--min-size`: Minimum file size in bytes (for `list` and `clean` commands).
- `--dry-run`: Perform a dry run without deleting files (for `clean` command).
- `--verbose`: Enable verbose output (for `clean` command).
- `--confirm`: Skip confirmation prompt (for `clean` command).
- `--ds-store`: Only remove `.DS_Store` files (for `clean` command).

## Examples

### List all files in a directory ordered by size

```sh
./maclnr list --dir /Users/username/Documents
```

### List files larger than 1 MB

```sh
./maclnr list --dir /Users/username/Documents --min-size 1048576
```

### Clean `.DS_Store` files from a directory

```sh
./maclnr clean --dir /Users/username/Documents --ds-store
```

### Clean files larger than 1 MB

```sh
./maclnr clean --dir /Users/username/Documents --min-size 1048576
```

### Clean files with verbose output

```sh
./maclnr clean --dir /Users/username/Documents --verbose
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

This updated `README.md` file provides installation instructions, usage examples, and explanations for each command and flag available in the `maclnr` tool, including the new functionality to clean files based on a minimum size recursively.