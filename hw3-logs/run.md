**Example to run program**

**http path**

go run main.go -p https://raw.githubusercontent.com/elastic/examples/master/Common%20Data%20Formats/nginx_logs/nginx_logs --f json --o ../../remote_logs.json

**local path**

go run main.go -p '../../logs/2025/*' --f json --o ../../log.json

