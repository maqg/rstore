all::
	go build -o rstore main.go


prepare::
	go run tools/funcmap/create_funcmap.go > api/funcmap.go 

clean::
	rm -rf bin/*
	rm -rf rstore

fmt::
	go fmt .

package::
	./PACKAGE.sh

build::
	./PACKAGE.sh
