# docker rest API

## About
This yada yada
docker network create restNet
docker run -d --network="restNet" --name dockerMongoDB -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=pdkirby -e MONGO_INITDB_ROOT_PASSWORD=DontLookAtMyPassword mongo:latest
docker run -d --network="restNet" --name restApi -p 8081:8081 -e mongoDSN=mongodb://pdkirby:DontLookAtMyPassword@dockerMongoDB:27017 pdkirby/golangrestapi:latest

for building a new rest docker image:

docker build -t pdkirby/golangrestapi:yourtag .

or 

docker compose up -d
