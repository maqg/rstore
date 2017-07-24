all:: fmt
	go build -o ../bin/rstore main.go
	go build -o ../bin/rstcli rstcli.go


clean::
	rm -rf rstore
	rm -rf rstcli


fmt::
	go fmt .
