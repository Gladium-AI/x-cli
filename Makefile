BINARY  := x-cli
PREFIX  := /usr/local
BINDIR  := $(PREFIX)/bin

.PHONY: build install uninstall clean

build:
	go build -o $(BINARY) .

install: build
	install -d $(BINDIR)
	install -m 755 $(BINARY) $(BINDIR)/$(BINARY)

uninstall:
	rm -f $(BINDIR)/$(BINARY)

clean:
	rm -f $(BINARY)
