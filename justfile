run * ARGS:
  go run . --config mpcube.cfg {{ ARGS }}

build:
  go build .

fmt:
  go fmt

test:
  go test ./...
