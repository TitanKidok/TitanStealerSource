set GOARCH=386
go build -trimpath -ldflags "-s -H windowsgui -w" -gcflags=all="-l -B"