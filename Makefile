.SILENT: 
build: clean
	echo "Building for Linux"
	GOOS=linux GOARCH=386 go build -o bin/NekoGo-linux *.go
	echo Building for Windows
	GOOS=windows GOARCH=386 go build -o bin/NekoGo-windows.exe *.go
	echo "Build successful"

.SILENT:
run: clean
	go run *.go

.SILENT:
clean:
	echo "Cleaning"
	go clean
	echo "Clean successful"
