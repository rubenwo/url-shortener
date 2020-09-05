# url-shortener

## Description
This is a simple url shortener written in Go. Deploy with docker-compose. 

The docker-compose.yml file describes the server service and a redis service which uses a persistent volume.

Change the certs folder defined in the docker-compose.yml file to production certificates when deploying.
This folder needs to contain the following files:
- server.pem (certificate)
- server.key (key)