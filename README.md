# goprod
Production-Ready Go: Tools, Patterns and Techniques

## TODO
- Add table driven tests
- Add goimports
- Add postgres database and migrations
- Add docker-compose
- Add Tasks


## Run Container
podman run -d --name my-go-container -p 8080:8080 -v $(pwd):/app back-to-go

## Follow Logs
podman logs -f my-go-container
