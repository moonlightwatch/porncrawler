
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags " -v -s -w -linkmode 'external' -extldflags '-static'" -o porncrawler ./main
