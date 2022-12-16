package genericpam

// Run using "cd pam; go generate ./..."
// Uses forked oapi-codegen with fiber support, follow these instructions to run:
// 1. clone forked oapi-codegen: "git clone https://github.com/four-fingers/oapi-codegen.git"
// 2. build the command: "cd oapi-codegen; go build ./cmd/oapi-codegen"
// 3. go:generate should point to the locally built command binary
//go:generate ../../oapi-codegen/oapi-codegen --config model.cfg.yml ../../valkyrie/pam/pam_api.yml
//go:generate ../../oapi-codegen/oapi-codegen --config handlers.cfg.yml ../../valkyrie/pam/pam_api.yml
