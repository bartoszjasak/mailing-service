# Mailing service

## Deployment

### Prerequisities:
1. minikube installed
2. skaffold installed (https://skaffold.dev/docs/install/)
3. goose installed ("go install github.com/pressly/goose/v3/cmd/goose@latest") 

### Deployment:
1. minikube start
2. skaffold dev 
3. cd db\sqlc\migrations
4. goose postgres "postgres://postgres:postgres@localhost:5432/postgres" up
5. App will be running on localhost:8080