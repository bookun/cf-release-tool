.PHONY: install
install:
	go build -o ReleaseTool
	cf uninstall-plugin ReleaseTool
	cf install-plugin -f ReleaseTool

.PHONY: plug-install
plug-install:
	cf uninstall-plugin ReleaseTool
	cf install-plugin -f ReleaseTool
