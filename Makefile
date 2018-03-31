VERSION :=	0.0.2dev

.PHONY:	dist

all:
	go build

clean:
	find . -name \*~ -print0 | xargs -0 rm -f
	rm -f cryptotrader
	rm -rf dist

dist: LDFLAGS = -w -s
dist:
	rm -rf dist

	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" \
		-o dist/cryptotrader-$(VERSION)-linux-x64/cryptotrader

	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" \
		-o dist/cryptotrader-$(VERSION)-windows-x64/cryptotrader.exe

	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" \
		-o dist/cryptotrader-$(VERSION)-macos-x64/cryptotrader

	cd dist && \
		for d in *; do \
			cp ../LICENSE.txt $$d; \
			cp ../README.md $$d; \
			zip -r $$d.zip $$d; \
		done
