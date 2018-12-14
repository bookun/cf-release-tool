.PHONY: build-install
build-install:
	go build -o ReleaseTool
	cf install-plugin -f ReleaseTool

.PHONY: reinstall
reinstall:
	cf install-plugin -f ReleaseTool
