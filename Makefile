build:
	clean
	echo "Building for Linux"
	GOOS=linux GOARCH=386 go build -o bin/NekoGo-linux *.go
	echo Building for Windows
	GOOS=windows GOARCH=386 go build -o bin/NekoGo-windows.exe *.go
	echo "Build successful"

run:
	go run *.go

clean:
	go clean
