set GOARCH=386
go build -trimpath -ldflags "-s -w" -gcflags=all="-l -B"