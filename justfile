run * ARGS:
  go run . {{ ARGS }}

build:
  go build .

fmt:
  go fmt

test:
  go test ./...
