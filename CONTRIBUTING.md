# Contributing to Vortex

Thank you for your interest in contributing to Vortex! This document provides guidelines and information for contributors.

## Table of Contents

- [Project Overview](#project-overview)
- [Development Setup](#development-setup)
- [Code Style and Standards](#code-style-and-standards)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Issue Reporting](#issue-reporting)
- [Pull Request Guidelines](#pull-request-guidelines)
- [Release Process](#release-process)

## Project Overview

Vortex is a modern frontend web framework for Go that compiles to WebAssembly. It provides:

- **Virtual DOM**: Efficient rendering and diffing
- **Component System**: Reusable UI components with lifecycle management
- **Type-Safe Styling**: CSS-in-Go with full type safety
- **Animation Engine**: Advanced animations, particles, and timeline effects
- **CLI Tooling**: Project initialization, building, and development server
- **Async Utilities**: Promise-like patterns for WebAssembly

## Development Setup

### Prerequisites

- **Go 1.25+** (required by the project)
- **Git**

### Getting Started

1. **Fork and clone the repository:**

   ```bash
   git clone https://github.com/YOUR_USERNAME/vortex.git
   cd vortex
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Install the CLI locally:**

   ```bash
   go install ./cmd/vortex
   ```

4. **Verify installation:**

   ```bash
   vortex --help
   ```

5. **Run tests:**
   ```bash
   go test ./...
   ```

## Code Style and Standards

### Go Code Standards

- **Follow Go conventions**: Use `gofmt`, `go vet`, and `golint`
- **File naming**: Use lowercase with underscores (e.g., `component_base.go`, not `ComponentBase.go`)
- **Package naming**: Short, concise, lowercase names
- **Error handling**: Always handle errors explicitly
- **Documentation**: Public functions and types must have godoc comments

### WebAssembly Specific

- **Build constraints**: Use `//go:build js && wasm` for WebAssembly-specific code
- **JS interop**: Use `syscall/js` package responsibly, avoid memory leaks
- **Performance**: Be mindful of WebAssembly performance characteristics

### Styling Framework

- **Type safety**: All CSS properties should have typed constants
- **Functional approach**: Use functional options pattern for styles
- **Consistency**: Follow existing naming conventions for CSS properties

### Example Code Style

```go
// Good: Proper godoc, typed constants, functional options
package style

// DisplayValue represents CSS display property values
type DisplayValue string

const (
    DisplayBlock DisplayValue = "block"
    DisplayFlex  DisplayValue = "flex"
)

// Display sets the display CSS property
func Display(value DisplayValue) StyleOption {
    return func(s *Style) {
        s.Base["display"] = string(value)
    }
}
```

## Project Structure

```
vortex/
├── cmd/vortex/           # CLI application
│   ├── main.go          # Entry point
│   ├── root.go          # Root command
│   ├── init.go          # Project initialization
│   ├── build.go         # Build command
│   └── dev.go           # Development server
├── pkg/                 # Public library code
│   ├── vdom/            # Virtual DOM implementation
│   ├── component/       # UI components
│   ├── style/           # Styling system
│   ├── animation/       # Animation engine
│   ├── async/           # Async utilities
│   └── renderer/        # Rendering engine
├── internal/            # Internal packages (if needed)
├── examples/            # Usage examples
├── docs/                # Documentation
└── testdata/           # Test fixtures
```

### Package Guidelines

- **`cmd/`**: Contains executable commands
- **`pkg/`**: Public API, importable by other projects
- **`internal/`**: Private code, not importable
- **Keep packages focused**: Each package should have a single responsibility

## Testing

### Test Requirements

- **Unit tests**: All public functions should have unit tests
- **Integration tests**: Test CLI commands end-to-end
- **WebAssembly tests**: Test WASM-specific functionality where possible
- **Examples**: Ensure all examples compile and work

### Test Organization

```go
// Good test structure
func TestStyleDisplay(t *testing.T) {
    tests := []struct {
        name     string
        value    DisplayValue
        expected string
    }{
        {"block display", DisplayBlock, "block"},
        {"flex display", DisplayFlex, "flex"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            style := New(Display(tt.value))
            if got := style.Base["display"]; got != tt.expected {
                t.Errorf("Display() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/style/

# Run with verbose output
go test -v ./...
```

## Submitting Changes

### Workflow

1. **Create a branch** from `main`:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the style guidelines

3. **Write or update tests** for your changes

4. **Run the full test suite**:

   ```bash
   go test ./...
   go vet ./...
   gofmt -s -w .
   ```

5. **Update documentation** if necessary

6. **Commit with descriptive messages**:

   ```bash
   git commit -m "feat(component): add Button hover animations

   - Add OnHover styling support to Button component
   - Include transition animations for smooth effects
   - Update examples to demonstrate hover states"
   ```

7. **Push and create a pull request**

### Commit Message Format

Use conventional commits format:

```
type(scope): short description

Longer description if necessary

- Bullet points for changes
- Reference issues with #123
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

## Issue Reporting

### Before Creating an Issue

1. **Search existing issues** to avoid duplicates
2. **Check the latest version** - issue might be fixed
3. **Prepare minimal reproduction** if reporting a bug

### Bug Reports

Include:

- **Vortex version**: `vortex --version`
- **Go version**: `go version`
- **Operating system**: Windows, macOS, Linux
- **Browser** (for WebAssembly issues): Chrome, Firefox, Safari
- **Minimal code example** that reproduces the issue
- **Expected vs actual behavior**
- **Error messages** (full stack traces)

### Feature Requests

Include:

- **Clear use case**: Why is this needed?
- **Proposed API**: How should it work?
- **Alternative solutions**: What workarounds exist?
- **Breaking changes**: Would this affect existing code?

## Pull Request Guidelines

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] Tests pass: `go test ./...`
- [ ] Code is properly formatted: `gofmt -s -w .`
- [ ] Documentation is updated
- [ ] Examples work (if applicable)
- [ ] No breaking changes (or clearly documented)

### PR Description Template

```markdown
## Description

Brief description of changes

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing

- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing performed

## Checklist

- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Examples updated (if applicable)
```

### Review Process

1. **Automated checks** must pass (tests, linting)
2. **Code review** by maintainers
3. **Testing** of new functionality
4. **Approval** and merge

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Workflow

1. **Create release branch**: `release/v1.2.3`
2. **Update version numbers** and documentation
3. **Create tag**: `git tag v1.2.3`
4. **Build and test** release artifacts
5. **Create GitHub release** with changelog
6. **Publish** to package repositories

## Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and community discussion
- **Documentation**: Check the `docs/` directory
- **Examples**: Look at the `examples/` directory

## Code of Conduct

We are committed to providing a welcoming and inclusive experience for everyone. We expect all contributors to:

- Be respectful and inclusive in interactions
- Focus on what is best for the community
- Show empathy towards other community members
- Accept responsibility for mistakes and learn from them

## Recognition

Contributors will be recognized in:

- **GitHub contributors list**
- **Release notes** for significant contributions
- **Project documentation** for major features

Thank you for contributing to Vortex! 🌪️


