all::
	go build -o rstore main.go


prepare::

clean::
	rm -rf bin/*
	rm -rf rstore

fmt::
	go fmt .
