# Contributing to Runware Go SDK

Thank you for your interest in contributing to the Runware Go SDK! We welcome contributions from the community.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/sdk-go.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Run tests: `go test -v ./...`
6. Commit your changes: `git commit -am 'Add new feature'`
7. Push to the branch: `git push origin feature/your-feature-name`
8. Submit a pull request

## Development Guidelines

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code
- Use `golint` and `go vet` to check for common issues
- Write clear, descriptive comments for exported functions and types

### Testing

- Write unit tests for new functionality
- Ensure all tests pass before submitting a PR
- Aim for high test coverage (>80%)
- Use table-driven tests where appropriate

Example:

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "result", false},
        {"invalid input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Feature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Feature() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Feature() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Documentation

- Update the README.md if you add new features
- Add examples for new functionality in the `examples/` directory
- Document all exported functions, types, and constants
- Use godoc-style comments

### Commit Messages

- Use clear and descriptive commit messages
- Start with a verb in the present tense (e.g., "Add", "Fix", "Update")
- Reference issue numbers when applicable

Good examples:
- `Add support for video inference`
- `Fix reconnection logic in WebSocket client`
- `Update documentation for ControlNet usage`

### Pull Requests

- Provide a clear description of the changes
- Link to any related issues
- Ensure CI checks pass
- Be responsive to feedback and review comments

## Areas for Contribution

We welcome contributions in the following areas:

1. **New Features**: Implement support for new API endpoints or features
2. **Bug Fixes**: Fix reported bugs or issues
3. **Documentation**: Improve or expand documentation
4. **Examples**: Add new examples demonstrating SDK usage
5. **Tests**: Improve test coverage
6. **Performance**: Optimize performance-critical code
7. **Error Handling**: Improve error messages and handling

## Questions?

If you have questions about contributing, feel free to:

- Open an issue for discussion
- Check existing issues and pull requests
- Review the [documentation](https://runware.ai/docs)

## Code of Conduct

Be respectful and constructive in all interactions with the community.