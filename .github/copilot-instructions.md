# GitHub Copilot Instructions

## Code Style & Quality Guidelines

When generating code for this project, please ensure compliance with the following standards:

### Go Code Standards
- **Follow Go conventions**: Use `gofmt`, `goimports`, and standard Go idioms
- **Naming conventions**: 
  - Use camelCase for variables and functions
  - Use PascalCase for exported types and functions
  - Use meaningful, descriptive names
  - Avoid abbreviations unless they're well-known (e.g., `HTTP`, `URL`, `ID`)

### Linting Compliance
- **golangci-lint**: All generated code must pass our golangci-lint configuration
- **Common lint rules to follow**:
  - No unused variables or imports
  - Proper error handling (never ignore errors)
  - Use `require` for error assertions in tests, `assert` for other validations
  - Avoid useless assertions (comparing variables to themselves)
  - Add proper context to error messages

### Testing Standards
- **Test naming**: Use `Test<FunctionName>` pattern
- **Test structure**: Follow Arrange-Act-Assert pattern
- **Parallel tests**: Use `t.Parallel()` for independent tests
- **Error testing**: Use `require.Error()` for error assertions, `assert` for other checks
- **Coverage**: Aim for meaningful test coverage, not just line coverage
- **Table-driven tests**: Use for multiple test cases with similar structure

### Package Organization
- **Internal packages** (`internal/`): Private application code
  - `config/`: Configuration management
  - `db/`: Database repositories and data access
  - `handlers/`: HTTP request handlers
  - `middleware/`: HTTP middleware components
  - `models/`: Domain models and data structures
  - `server/`: HTTP server setup and lifecycle
  - `utilities/`: Internal utility functions
- **Public packages** (`pkg/`): Reusable libraries
  - `logger/`: Structured logging utilities
  - `mongodb/`: MongoDB connection management

### Code Patterns to Follow

#### Error Handling
```go
// Good: Proper error wrapping
if err != nil {
    return fmt.Errorf("failed to process request: %w", err)
}

// Good: Context-aware operations
func (r *Repository) GetByID(ctx context.Context, id string) (*Model, error) {
    // Implementation with proper context usage
}
```

#### Testing Patterns
```go
// Good: Table-driven tests
func TestFunction(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // Test implementation
        })
    }
}
```

#### MongoDB Patterns
```go
// Good: Use functional options pattern
dbConnMgr, err := mongodb.NewMongoManager(
    svcEnv.DBHosts,
    svcEnv.DBName,
    credentials,
    mongodb.WithQueryLogging(true),
    mongodb.WithReplicaSet("rs0"),
)
```

#### Logging Patterns
```go
// Good: Structured logging with context
lgr.Info().
    Str("requestID", reqID).
    Str("operation", "create_order").
    Dur("elapsed", elapsed).
    Msg("operation completed successfully")
```

### Architecture Principles
- **Dependency Injection**: Use interfaces and inject dependencies
- **Single Responsibility**: Each function/type should have one clear purpose
- **Interface Segregation**: Keep interfaces small and focused
- **Error Boundaries**: Handle errors at appropriate levels
- **Context Propagation**: Pass context through all operations that might be cancelled

### Security Considerations
- **Input Validation**: Always validate and sanitize inputs
- **Secrets Management**: Never hardcode secrets, use environment variables or sidecar files
- **SQL Injection Prevention**: Use parameterized queries (though we use MongoDB)
- **OWASP Compliance**: Follow security headers and validation patterns

### Performance Guidelines
- **Connection Pooling**: Reuse database connections
- **Context Timeouts**: Set appropriate timeouts for operations
- **Graceful Shutdown**: Implement proper cleanup on application termination
- **Resource Management**: Always close/cleanup resources (defer patterns)

## Specific to This Project

### MongoDB Connection Management
- Use the `pkg/mongodb` package for all database connections
- Follow the functional options pattern for configuration
- Always use context-aware operations
- Implement proper connection cleanup

### Middleware Chain
- Keep middleware focused and composable
- Use proper request ID propagation
- Implement structured logging in middleware
- Handle panics gracefully

### Configuration Management
- Use the `internal/config` package for environment configuration
- Validate all required configuration at startup
- Provide sensible defaults where appropriate

When generating code, please ensure it follows these patterns and will pass both our linting rules and maintain consistency with the existing codebase architecture.