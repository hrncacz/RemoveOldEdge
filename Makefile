run:
	env GOOS=windows GOARCH=amd64 go build .
	cp remove-old-edge.exe /mnt/c/_install

