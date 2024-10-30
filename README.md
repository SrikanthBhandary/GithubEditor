# GitHub Editor

**GitHub Editor** is a command-line tool written in Go that allows users to clone GitHub repositories, modify files using regex, and commit and push changes back to the repository. This tool simplifies the process of making quick edits to files in GitHub projects.

## Features

- Clone GitHub repositories using a personal access token.
- Checkout specific branches.
- Modify files using regex replacements.
- Commit and push changes to the repository.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/github-editor.git
   cd github-editor
   ```

2. Build the project:
   ```bash
   go build -o ghedit
   ```

## Usage

### Command Structure

```bash
ghedit --repo <repository-url> --branch <branch-name> --file <file-to-modify> --regEx <regex-pattern> --val <replacement-value> --token <github-token> --username <github-username>
```

### Flags

- `--repo`, `-r`: **(Required)** Repository to clone (e.g., `https://github.com/user/repo.git`).
- `--branch`, `-b`: **(Required)** Branch to checkout (e.g., `main`).
- `--file`, `-f`: **(Required)** File to modify (e.g., `path/to/file.go`).
- `--regEx`, `-e`: Regex pattern to find and replace (optional).
- `--val`, `-v`: Value to update in the file (optional).
- `--token`, `-t`: GitHub personal access token (required for authentication).
- `--username`, `-u`: GitHub username (required for authentication).

### Example

```bash
ghedit -r https://github.com/user/repo.git -b main -f path/to/file.go -e "// linter:\\d+" -v "// linter:12345" -t your-token -u your-username
```

This command will:
1. Clone the specified repository.
2. Checkout the `main` branch.
3. Replace all occurrences of comments matching the regex `// linter:\d+` with `// linter:12345` in `file.go`.
4. Commit the changes with a default message and push to the repository.

## Requirements

- Go 1.15 or higher.
- A valid GitHub personal access token with repo access.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for discussion.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Feel free to replace any placeholders (like `yourusername` and `your-token`) with your actual details. You may also want to customize sections based on any additional features or requirements specific to your project.
