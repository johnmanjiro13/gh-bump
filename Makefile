.PHONY: test

test:
	go test -cover ./...

mockgen:
	mockgen -destination bump/mock_bump/mock_bump.go github.com/johnmanjiro13/gh-bump/bump Gh
