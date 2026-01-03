//go:build generate

package goprod

// Generate OpenAPI Stubs
//go:generate go tool oapi-codegen --config=api/openapi/openapi-config.yaml api/openapi/openapi.yaml
