src = $(shell find . -name '*.go')

journal: $(src)
	go build -o ./journal ./main.go
