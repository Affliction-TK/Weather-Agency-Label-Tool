# Contributing to Weather Agency Label Tool

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- Clear title and description
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version, Node version, MySQL version)
- Screenshots if applicable

### Suggesting Features

Feature requests are welcome! Please:
- Check if the feature already exists or is planned
- Describe the feature and its use case
- Explain why it would be useful
- Provide examples if possible

### Code Contributions

#### Setup Development Environment

1. Fork the repository
2. Clone your fork
3. Follow the [Quick Start Guide](QUICKSTART.md)
4. Create a new branch: `git checkout -b feature/your-feature-name`

#### Making Changes

**Backend (Go):**
- Follow Go conventions and best practices
- Use `go fmt` to format code
- Add tests for new functionality
- Update API documentation if adding/changing endpoints

**Frontend (Svelte):**
- Follow Svelte best practices
- Keep components modular and reusable
- Ensure responsive design
- Test in multiple browsers if possible

**Database:**
- Use migrations for schema changes
- Maintain backward compatibility when possible
- Add indexes for new queries
- Document schema changes

#### Testing

```bash
# Run Go tests
go test -v ./...

# Build frontend
cd frontend && npm run build

# Manual testing
./server
```

#### Commit Messages

Use clear, descriptive commit messages:
```
feat: Add image batch upload functionality
fix: Correct nearest station calculation for edge cases
docs: Update API documentation for new endpoint
refactor: Simplify annotation form validation
test: Add unit tests for haversine distance
```

Prefix conventions:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation only
- `style:` Code style (formatting, semicolons, etc.)
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance tasks

#### Pull Requests

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG if present
5. Create pull request with clear description:
   - What changes were made
   - Why they were made
   - How to test them
   - Screenshots for UI changes

### Code Style

**Go:**
- Use `gofmt` and `golint`
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Add comments for exported functions
- Keep functions focused and small

**JavaScript/Svelte:**
- Use 2 spaces for indentation
- Use semicolons
- Prefer `const` over `let`, avoid `var`
- Use meaningful variable names
- Add JSDoc comments for complex functions

**SQL:**
- Use uppercase for SQL keywords
- Indent nested queries
- Add comments for complex queries
- Use prepared statements always

### Project Structure

```
.
├── main.go              # Backend server
├── main_test.go         # Backend tests
├── go.mod/go.sum        # Go dependencies
├── schema.sql           # Database schema
├── uploads/             # Uploaded images
├── frontend/            # Svelte frontend
│   ├── src/
│   │   ├── App.svelte
│   │   ├── lib/         # Components
│   │   └── main.js
│   └── dist/            # Built frontend
├── docs/                # Additional documentation
└── scripts/             # Helper scripts
```

### Adding New Features

#### Adding a New API Endpoint

1. Add handler function in `main.go`
2. Register route in `main()` function
3. Add tests in `main_test.go`
4. Update `API.md` documentation
5. Update frontend if needed

Example:
```go
func getStatistics(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// In main()
api.HandleFunc("/statistics", getStatistics).Methods("GET")
```

#### Adding a New Frontend Component

1. Create component in `frontend/src/lib/`
2. Import and use in parent component
3. Follow existing style and patterns
4. Test in browser

Example:
```svelte
<!-- frontend/src/lib/Statistics.svelte -->
<script>
  export let data;
</script>

<div class="statistics">
  <!-- Component content -->
</div>

<style>
  .statistics { /* styles */ }
</style>
```

#### Adding a New Database Table

1. Update `schema.sql`
2. Add corresponding Go struct
3. Create database functions
4. Add API endpoints if needed
5. Update frontend if needed

### Security

- Never commit secrets or credentials
- Validate all user inputs
- Use prepared statements for SQL
- Sanitize file uploads
- Follow OWASP guidelines
- Report security issues privately

### Documentation

When adding features, update:
- README.md (if changing core functionality)
- API.md (if adding/changing API)
- QUICKSTART.md (if affecting setup)
- DEPLOYMENT.md (if affecting deployment)
- Code comments (for complex logic)

### Questions?

- Check existing issues and documentation
- Ask in GitHub issues
- Be respectful and patient

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Keep discussions on-topic
- Report inappropriate behavior

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT).

---

Thank you for contributing to the Weather Agency Label Tool! 🌤️
