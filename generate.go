//go:build generate

package goprod

//go:generate go tool oapi-codegen --config=api/openapi/openapi-config.yaml api/openapi/openapi.yaml
//go:generate go tool mockery
//go:generate go tool golangci-lint run ./...
