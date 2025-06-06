src = $(shell find . -name '*.go')

journal: $(src)
	go mod tidy
	go build .

install:
	go install .
