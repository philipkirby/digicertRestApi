# Welcome to My Dockerized REST API with MongoDB

## Greetings  👋
Thank you for exploring this technical interview project showcasing a Dockerized REST API using MongoDB.
This project is designed as a response to the technical interview with DigiCert.


## Key Features
* RESTful API: Utilising Go to create a RESTful API that performs basic CRUD operations as well as a full library list request.
* MongoDB: The application utilises MongoDB as its database; if any other database is to be used, it can be easily implemented using the db interface in the db directory.
* Docker Compose: For a prebuilt environment, see [docker-compose.yaml](docker-compose.yaml). This will set up the entire environment for you with `docker compose up -d`.



## Manual rebuilding
To build and run the Docker environment manually, you can use the [DockerFile](Dockerfile).
This will build the go programme in a golang:alpine container, and then the binary will be moved to a minimal image for execution.
To build the Docker container, execute:

`docker build -t pdkirby/golangrestapi:yourtag .`

To create the network, execute:

`docker network create restNet`

To run the MongoDB instance, execute:

`docker run -d -v mongodb_data:/data/db --network="restNet" --name dockerMongoDB -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=pdkirby -e MONGO_INITDB_ROOT_PASSWORD=DontLookAtMyPassword mongo:latest`

To run the rest api, execute:

`docker run -d --network="restNet" --name restApi -p 8081:8081 -e mongoDSN=mongodb://pdkirby:DontLookAtMyPassword@dockerMongoDB:27017 -e restPort=8081 pdkirby/golangrestapi:yourtag`

## Misc
### Postman
If you would like to test the rest api with Postman, you can import my postman calls found in [rest.postman_collection_for_testing.json](rest.postman_collection_for_testing.json).

### Local library
In a larger system, the DB interface and lib would be in a separate library, as the odds are not just that the restapi would want access to the book library.
But for the sake of this technical interview, both will be stored in this project directory.


## RestApi Data
See below for examples of rest calls:

#### List
GET: `http://localhost:8081/api/library/getlist?Content-Type=application/json`

#### Create
PUT: `http://localhost:8081/api/library/create?Content-Type=application/json`

body:
{
"name": "harry potter 2",
"author": "JKR",
"contents": "A boy once nearly died."
}

#### Retrieve
GET : `http://localhost:8081/api/library/get/harry potter 2/JKR?Content-Type=application/json`

#### Update
PUT: `http://localhost:8081/api/library/update?Content-Type=application/json`

body:
{
"name": "harry potter 2",
"author": "JKR",
"contents": "A boy once died."
}

#### Delete
PUT: `http://localhost:8081/api/library/delete/harry potter 2/JKR?Content-Type=application/json`

