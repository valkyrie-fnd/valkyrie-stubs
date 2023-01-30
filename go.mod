module github.com/valkyrie-fnd/valkyrie-stubs

go 1.19

require (
	github.com/four-fingers/oapi-codegen v0.0.0-20221219135408-9237c9743c67
	github.com/gofiber/fiber/v2 v2.38.1
	github.com/joho/godotenv v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.28.0
	github.com/stretchr/testify v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/four-fingers/oapi-codegen-runtime v0.0.0-20230125082134-9d9fdf1239ab // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.40.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20221010170243-090e33056c14 // indirect
)

// avoids depending on all of oapi-codegen's dependencies
replace github.com/four-fingers/oapi-codegen => github.com/four-fingers/oapi-codegen-runtime v0.1.0
