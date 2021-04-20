
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags " -s -w -linkmode 'external' -extldflags '-static'" -o porncrawler.exe ./main
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags " -s -w -linkmode 'external' -extldflags '-static'" -o porncrawler ./main
