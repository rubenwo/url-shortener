# url-shortener

## Description
This is a url shortener written in Go. 

We now support malicious browser checking using the Safe Browsing API from Google.
To run the app with the safebrowsing check, it is necessary to provide an API key that can be retrieved from google developper panel as stated here. This key has to be provisioned in a .env file or as an environment variable called API_KEY so the backend can make use of.

In case you don't want to make use of the url checker, a dummy interface is also provided which allows any url to be stored in the database.

### Deploy with docker-compose.
The docker-compose.yml file describes the server service and a redis service which uses a persistent volume.

Change the certs folder defined in the docker-compose.yml file to production certificates when deploying.
This folder needs to contain the following files:
- server.pem (certificate)
- server.key (key)