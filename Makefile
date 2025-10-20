.PHONY: test

test:
	@go test -coverprofile=./test/cover/lexer.out ./internal/lexer
	@go test -coverprofile=./test/cover/parser.out ./internal/parser
	@go tool cover -html=./test/cover/lexer.out -o ./test/cover/lexer.html
	@go tool cover -html=./test/cover/parser.out -o ./test/cover/parser.html
	@rm ./test/cover/*.out
