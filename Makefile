.PHONY: install uninstall
install:
	sudo rm -rf /usr/bin/vanfs
	go build 
	sudo mv vanfs /usr/bin 
uninstall:
	sudo rm -rf /usr/bin/vanfs