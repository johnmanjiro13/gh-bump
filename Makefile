.PHONY: test mockgen

test:
	go test -cover ./...

mockgen:
	mockgen -destination mock/mock_bump.go -package mock github.com/johnmanjiro13/gh-bump/bump Gh
