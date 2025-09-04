# License Manager Tests

This directory contains all tests for the License Manager application, organized following Go best practices.

## Structure

```
tests/
├── unit/              # Unit tests for individual components
│   ├── ssh_service_test.go
│   └── handlers_test.go
├── integration/       # Integration tests for API endpoints
│   └── api_test.go
├── fixtures/          # Test data and fixtures
│   └── test_data.go
└── README.md          # This file
```

## Running Tests

### Using Make (Recommended)
```bash
make test              # Run all tests
make test-unit         # Run unit tests only
make test-integration  # Run integration tests only
make test-verbose      # Run with verbose output
make test-coverage     # Run with coverage report
make test-benchmark    # Run benchmark tests
```

### Using Docker directly
```bash
# Run all tests
docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine go test ./tests/...

# Run unit tests only
docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine go test ./tests/unit/...

# Run with coverage
docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine go test -cover ./tests/...
```


## Test Types

### Unit Tests (`tests/unit/`)
- Test individual functions and methods in isolation
- Use mocks and stubs where appropriate
- Fast execution, no external dependencies
- Examples: SSH service methods, handler functions

### Integration Tests (`tests/integration/`)
- Test complete workflows and API endpoints
- Test interactions between components
- May use real external services (with test data)
- Examples: Full API request/response cycles

### Fixtures (`tests/fixtures/`)
- Test data and configuration
- Reusable test objects
- Mock data for consistent testing


## Best Practices

1. **Package Organization**: Tests are organized by type (unit, integration) rather than by source package
2. **Docker-Only**: All tests run in Docker containers to avoid local Go installation
3. **Comprehensive Coverage**: Tests cover happy paths, error cases, and edge cases
4. **Test Data**: Use fixtures for consistent test data
5. **Utilities**: Reuse common test utilities to reduce duplication

## Adding New Tests

1. **Unit Tests**: Add to `tests/unit/` with descriptive names
2. **Integration Tests**: Add to `tests/integration/` for API/component interactions
3. **Test Data**: Add reusable test data to `tests/fixtures/`

## Example Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
    }{
        {
            name:     "valid input",
            input:    fixtures.ValidInput,
            expected: fixtures.ExpectedOutput,
        },
        {
            name:     "invalid input",
            input:    fixtures.InvalidInput,
            expected: fixtures.ExpectedError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionToTest(tt.input)
            // assertions...
        })
    }
}
```
