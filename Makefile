all: accountd

accountd:
	go build

.PHONY: fmt clean cleanall dist deploy

fmt:
	find . -name "*.go" | xargs -I {} go fmt {}

clean:
	-rm accountd

cleanall:
	-rm -rf accountd dist dist.tar.gz

dist:
	-mkdir -p dist/web/views/CollectionJSON/v1/
	cp README.md dist
	cp accountd dist
	cp -r web/views/CollectionJSON/v1/* dist/web/views/CollectionJSON/v1/
	mkdir -p dist/env
	cp -r env/production dist/env/
	rm env/working
	cd env && ln -s production working
	tar czvf dist.tar.gz ./dist

deploy:
	scp -r dist/* zack@172.16.140.17:/var/www/accountd
