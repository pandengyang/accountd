all: accountd

accountd:
	go build

.PHONY: fmt clean cleanall

fmt:
	find . -name "*.go" | xargs -I {} go fmt {}

clean:
	-rm accountd

cleanall:
	-rm -rf accountd
