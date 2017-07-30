all::
	go build -o bin/rstore src/main.go
	go build -o bin/rstcli src/rstcli/rstcli.go


prepare::


clean::
	rm -rf bin/*
	rm -rf rstore
	rm -rf rstcli


fmt::
	go fmt .
