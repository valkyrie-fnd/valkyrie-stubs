module github.com/valkyrie-fnd/valkyrie-stubs

go 1.20

require (
	github.com/four-fingers/oapi-codegen v0.0.0-20221219135408-9237c9743c67
	github.com/gofiber/fiber/v2 v2.42.0
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.29.0
	github.com/stretchr/testify v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20230203094617-bcbc01813b4f // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/four-fingers/oapi-codegen-runtime v0.0.0-20230125082134-9d9fdf1239ab // indirect
	github.com/gofiber/template v1.7.5
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.15.15 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.44.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.5.0 // indirect

)

// avoids depending on all of oapi-codegen's dependencies
replace github.com/four-fingers/oapi-codegen => github.com/four-fingers/oapi-codegen-runtime v0.1.0
