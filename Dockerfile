FROM golang:alpine as builder
ENV CGO_ENABLED=0
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -ldflags="-w -s" .

FROM scratch
WORKDIR /app
COPY --from=builder /app/pam-stub /usr/bin/
ENTRYPOINT ["pam-stub"]