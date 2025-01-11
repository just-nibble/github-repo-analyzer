# GitHub Repository Analyzer

A command-line tool that analyzes the file structure and sizes of GitHub repositories.

## Features

- Clones public GitHub repositories
- Analyzes file and directory structure
- Calculates file sizes in MB
- Outputs analysis in JSON format
- Handles nested folder structures
- Includes comprehensive error handling
- Cleans up temporary files automatically

## Requirements

- Go 1.16 or higher
- Git installed on your system

## Installation

  1. Clone this repository:

```bash
git clone https://github.com/just-nibble/github-repo-analyzer.git
```

2. Install dependencies:

```bash
go mod init github-repo-analyzer
go get github.com/go-git/go-git/v5
```

## Usage

Run the analyzer by providing a GitHub repository URL:

```bash
go run main.go https://github.com/username/repository.git
```

The tool will output JSON-formatted analysis including:

- Repository clone URL
- Total repository size
- Folder structure with file sizes

Example output:

```json
{
  "clone_url": "https://github.com/example/example.git",
  "size": 20.45,
  "folders": [
    {
      "name": "src",
      "files": [
        {
          "name": "main.go",
          "size": 0.15
        }
      ]
    }
  ]
}
```

## Running Tests

To run the test suite:

```bash
go test -v ./...
```

## Error Handling

The tool handles various error cases:

- Invalid GitHub URLs
- Repository cloning failures
- File system access issues
- JSON formatting errors

## Performance Considerations

- Uses efficient file system traversal
- Minimizes memory usage during analysis
- Cleans up temporary files automatically
- Handles large repositories gracefully

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
