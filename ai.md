# Project Context: Bookmark Manager (Go)

## Architecture & Transport
- **Layering**: Hexagonal / Clean Architecture.
- **Transport**: Located in `internal/infra/transport/rest`. 
- **API**: OpenAPI 3.1 driven. Use `oapi-codegen` for server stubs.
- **Package Naming**: Use `rest` for transport to avoid `net/http` collision.

## Tooling (Go 1.24+)
- **Tool Management**: Use the `tool` directive in `go.mod`. 
- **Generation**: Always run `go generate ./...` from the root.
- **Mocking**: Use `mockery` with the `EXPECT()` syntax (Expecter enabled).
- **Linting**: `golangci-lint` is configured in `.golangci.yml`.

## Coding Standards
- Interfaces live in `internal/domain`.
- Mocks live in `internal/mocks`.
- Use Table-Driven Tests for all handlers and services.
- Always use `context.Context` in service and repository methods.
