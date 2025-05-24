src = $(shell find . -name '*.go')

journal: $(src)
	go build -o ./journal ./journal.go
