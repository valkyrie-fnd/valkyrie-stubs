FROM golang:1.20-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -ldflags="-w -s" .

FROM scratch
WORKDIR /app
COPY --from=builder /app/datastore.config.yaml /app/
COPY --from=builder /app/valkyrie-stubs /usr/bin/
ENTRYPOINT ["valkyrie-stubs"]
