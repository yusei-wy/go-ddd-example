.PHYNMY: format
format:
	gofumpt -e -d -l -w .
	golangci-lint run --fix
