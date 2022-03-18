# GitNoter Backend Module
This is the backend service of gitnoter application.

## Start the backend service
There are several ways to start backend service. Refer any one of the following method
### Run application locally with make & go
#### Setup & start database
Make sure docker is installed on the system since below make commands use docker to start the database container
```
make network
make postgres
make createdb
```

#### Create configuration file from template
```shell
cp gitnoter.yaml .gitnoter.yaml
```

#### Start the server
Make sure that the `.gitnoter.yaml` file is configured correctly & database is up.
```shell
go run main.go migrateup
go run main.go serve
```

### Build & run backend service locally using executable
#### Build the backend service
```shell
go build
```
This will build the service & generate executable file.

#### Create configuration file from template
```shell
cp gitnoter.yaml .gitnoter.yaml
```

#### Start http server
Make sure that the `.gitnoter.yaml` file is configured correctly & database is up.
```shell
./gitnoter migrateup
./gitnoter serve
```

### Build docker image & run the application with docker locally

#### Build docker image locally
```
docker build -t gitnoter-backend:latest .
```

#### Start containers with required configuration
```
# setup pre-requisite
docker network create gn-network
docker run --name postgres12 --network gn-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
docker exec -it postgres12 createdb --username=root --owner=root gn_db

# inspect the ip-address of postgres container
docker inspect <postgres-container-id>

# run migration (replace '<postgres-container-ip>' with the actual ip-address)
docker run --network gn-network -it \
-e DATABASE_HOST='<postgres-container-ip>' \
-e DATABASE_PORT='5432' \
-e DATABASE_DBNAME='gn_db' \
-e DATABASE_USERNAME='root' \
-e DATABASE_PASSWORD='secret' \
-e DATABASE_DEBUG='true' \
gitnoter-backend /app/gitnoter migrateup

# start server (replace '<postgres-container-ip>' with the actual ip-address)
docker run --network gn-network -it \
-e SECRETKEY='secret' \
-e DATABASE_HOST='<postgres-container-ip>' \
-e DATABASE_PORT='5432' \
-e DATABASE_DBNAME='gn_db' \
-e DATABASE_USERNAME='root' \
-e DATABASE_PASSWORD='secret' \
-e DATABASE_DEBUG='true' \
-e HTTPSERVER_HOST='0.0.0.0' \
-e HTTPSERVER_PORT='8080' \
-e HTTPSERVER_DEBUG='true' \
-e OAUTH2_GITHUB_CLIENTID='<github_client_id>' \
-e OAUTH2_GITHUB_CLIENTSECRET='<github_client_secret>' \
-e OAUTH2_GITHUB_REDIRECTURL='http://localhost:8080/api/v1/oauth2/github/callback' \
-p 8080:8080 \
gitnoter-backend
```

This will start the server on the specified port
