install:
	go install .
uninstall:
	rm -f $(GOPATH)/bin/$(shell basename $(PWD))
	test -d ~/.sourdough && rm -rf ~/.sourdough
