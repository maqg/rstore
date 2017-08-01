all::
	go build -o bin/rstore main.go
	go build -o bin/rstcli rstcli/main.go


prepare::

clean::
	rm -rf bin/*

fmt::
	go fmt .
