# REST service

REST server integrated with MONGODB database, both services are contained within dockers. For simplicity of the example the ports are predefined, being **8080** for the REST application and **27017** for the database service. Both services listen on any of the network interfaces.

The REST application was developed in go 1.17v and the database service used the ready-made docker image found in the docker image repository.

## To start the services just run the command:

```
docker-compose up
```

## To monitor the logs of the REST application run:
```
docker-compose logs -f rest_app
```

If you wish to create the images and run outside docker compose follow these steps:

## Get **mongoDB** image:
```
docker pull mongo
```

## Create **rest_app** image:

```
docker build --tag rest_app .
```

## Create local network to connection between containers:
```
docker network create -d bridge containers_network
```

## Create mongo container:
```
docker run --rm --network containers_network -v mongo_volume --name mongodb -p 27017:27017 mongo
```

## Create rest_app container:
```
docker run --rm --network containers_network --name rest_app -p 8080:8080 rest_app
```

# Testing

For testing the endpoints cURL can be used, for example, below are the two test commands:

## Store data on database:
```
curl -X POST -d @${input} http://127.0.0.1:8080/input -v -H 'Content-Type: application/json'
```

Obs: In this command it is expected that the path to a json file is passed.

To facilitate testing an environment variable can be set with:
```
export input=/path/to/json/file
```

## Retrieve the last 10 records from the database:
```
curl -X GET http://127.0.0.1:8080/get
```
