.PHONY: install
install:
	go build -o release
	cf uninstall-plugin ReleaseTool
	cf install-plugin -f release
