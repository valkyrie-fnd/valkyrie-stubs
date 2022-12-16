# Generic PAM

A pam stub based on Valkyrie's [genericpam openapi spec](https://github.com/valkyrie-fnd/valkyrie/blob/main/pam/pam_api.yml).

The server and model implementations are generated using https://github.com/deepmap/oapi-codegen.

# Install oapi-codegen

```shell
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```

Make sure that `$(go env GOPATH)/bin` is added to `$PATH` in your shell. 

# Generate

```shell
go generate ./...
```
