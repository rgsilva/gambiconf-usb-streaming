build: timestampfs corruptfs streamfs

timestampfs:
	go build -o bin/timestampfs cmd/timestampfs.go

corruptfs:
	go build -o bin/corruptfs cmd/corruptfs.go

streamfs:
	go build -o bin/streamfs cmd/streamfs.go
