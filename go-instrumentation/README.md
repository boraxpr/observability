# GO Instrumentation

## How to initial and run

```bash
# Delete all go module files
rm -rf go.*

# Initial go module and download dependencies
go mod init go-instrumentation
go get ./...

# If you want to update dependencies
go get -u all ./...

# Run main.go in sample/ directory
go run ./sample
```

## Sample Docker Compose

```yaml
services:
  go-sample:
    container_name: go-sample
    build:
      context: .
      args:
        APP_NAME: sample
        APP_PORT: 8081
    environment:
      APP_PORT: 8081
    ports:
      - 8081:8081
  go-practice:
    container_name: go-practice
    build:
      context: .
      args:
        APP_NAME: practice
        APP_PORT: 8082
    environment:
      APP_PORT: 8082
    ports:
      - 8082:8082
```
