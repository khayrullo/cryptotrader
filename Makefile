VERSION :=	0.0.1

.PHONY:	dist

all:
	go build

clean:
	find . -name \*~ -print0 | xargs -0 rm -f
	rm -f ctrader
	rm -rf dist

dist: LDFLAGS = -w -s
dist:
	rm -rf dist

	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" \
		-o dist/ctrader-$(VERSION)-linux-x64/ctrader

	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" \
		-o dist/ctrader-$(VERSION)-windows-x64/ctrader.exe

	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" \
		-o dist/ctrader-$(VERSION)-macos-x64/ctrader

	cd dist && \
		for d in *; do \
			cp ../LICENSE.txt $$d; \
			cp ../README.md $$d; \
			zip -r $$d.zip $$d; \
		done
