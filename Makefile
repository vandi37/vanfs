.PHONY: install uninstall
install:
	go build 
	sudo mv vanfs /usr/bin 
uninstall:
	sudo rm -rf /usr/bin/vanfs