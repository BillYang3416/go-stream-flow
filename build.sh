# swagger
cd go-file-gate
swag init -g cmd/go-file-gate/main.go -o ./docs

# Run Docker commands
docker-compose down
docker volume rm go-stream-flow_pgdata
docker-compose up -d --build 


