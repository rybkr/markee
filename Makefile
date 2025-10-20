.PHONY: test

test:
	@go test -cover ./internal/lexer
	@go test -cover ./internal/parser
