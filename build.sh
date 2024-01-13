# swagger
cd go-flow-gateway
swag init -g cmd/go-flow-gateway/main.go -o ./docs

# Run Docker commands
docker-compose down
docker volume rm go-stream-flow_pgdata
docker-compose up -d --build 


