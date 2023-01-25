package genericpam

// Run using "cd pam; go generate ./..."
// Uses forked oapi-codegen with fiber support
//go:generate go run github.com/four-fingers/oapi-codegen/cmd/oapi-codegen@latest --config model.cfg.yml ../../valkyrie/pam/pam_api.yml
//go:generate go run github.com/four-fingers/oapi-codegen/cmd/oapi-codegen@latest --config handlers.cfg.yml ../../valkyrie/pam/pam_api.yml
