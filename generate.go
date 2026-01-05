//go:build generate

package goprod

//go:generate go tool oapi-codegen --config=api/rest/openapi-config.yaml api/rest/openapi.yaml
//go:generate go tool mockery
