all::
	go build -o rstore main.go

clean::
	rm -rf bin/*
	rm -rf rstore
	rm -rf package

fmt::
	go fmt .

package::
	./PACKAGE.sh

build::
	./PACKAGE.sh
