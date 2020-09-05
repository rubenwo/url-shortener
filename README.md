# url-shortener

## Description
This is a simple url shortener written in Go. Deploy with docker-compose. 

The docker-compose.yml file describes the server service and a redis service which uses a persistent volume.

The certs folder contains self-generated certificates for use with localhost. Change the folder defined in the docker-compose.yml file to production certificates when deploying.