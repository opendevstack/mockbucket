# MockBucket

The mock service of ODS

## Contents

This project provides a docker image with the following mocked services for ODS:

* GIT http server 
* APIs for creating:
  * Project
  * Repository
  * Webhooks   
  
**DISCLAIMER**: All mocked services are intended to support tests of ODS projects, the APIs are not fully functional  

## APIs
List of available APIs

### GIT

* `http://<host>:8080/scm/{projectKey}/{repositorySlug}.git`: git clone URL

### Rest

* `/projects [POST]`: Mimics the create project API 
* `/projects/{projectKey}/repos [POST]`: Mimics the create repository API
* `/projects/{projectKey}/repos/{repositorySlug}/webhooks [POST]`: Mimics the create webhook API **[INCOMPLETE]**



## Building and Running

```shell script
docker-compose up --build
```
