build:
	go build 

pi:
	GOOS=linux GOARCH=arm GOARM=6 go build -o co2Plotter main.go
