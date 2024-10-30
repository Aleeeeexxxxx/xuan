binaries: | xuan

xuan:
	go build -o xuan src/cmd/xuan/main.go
	chmod +x xuan

xuan_windows:
	GOOS=windows GOARCH=amd64 go build -o xuan.exe src/cmd/xuan/main.go

clean:
	rm -f xuan xuan.exe