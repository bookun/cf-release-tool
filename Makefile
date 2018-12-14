.PHONY: build-install
build-install:
	GOOS=linux GOARCH=amd64 go build -o ReleaseTool
	cf install-plugin -f ReleaseTool

.PHONY: reinstall
reinstall:
	cf install-plugin -f ReleaseTool
