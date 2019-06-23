<div style="display: flex; align-items: center;" align="center">
<a href="https://copybird.org"><img width="100px" src="https://raw.githubusercontent.com/copybird/copybird/master/docs/logo.svg?sanitize=true" alt="Copybird"></a>
</div>

# Copybird

[![](https://images.microbadger.com/badges/version/copybird/copybird.svg)](https://microbadger.com/images/copybird/copybird)
[![](https://images.microbadger.com/badges/image/copybird/copybird.svg)](https://microbadger.com/images/copybird/copybird)
[![](https://godoc.org/github.com/copybird/copybird?status.svg)](http://godoc.org/github.com/copybird/copybird)
[![GitHub release](https://img.shields.io/github/release/copybird/copybird/all.svg?style=flat-square)](https://github.com/copybird/copybird/releases)

![](https://travis-ci.org/copybird/copybird.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/copybird/copybird/badge.svg)](https://coveralls.io/github/copybird/copybird)
[![Go Report Card](https://goreportcard.com/badge/github.com/copybird/copybird)](https://goreportcard.com/report/github.com/copybird/copybird)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## About

Copybird is open-source **cloud-native** universal backup tool for databases and files.

It allows you to:
1. Create database backup
2. Compress backup file
3. Encrypt backup file
4. Send it to various destinations fast and secure
5. Get notification about backup status in messagers and notification services
6. Enjoy simple backup as a service with k8s backup controller

learn more at [copybird.org](https://copybird.org)

## Databases
Currently Copybird supports the following databases:
- MySQL
- Postgres
- MongoDB
- Etcd (v2 and v3 API)

## Compression
Copybird compresses with the following tools:
- gzip
- lz4

## Encryption
Copybird uses AES-GCM for Efficient Authenticated Encryption

## Output
Copybird can deliver encrypted compressed backup to the following destinations:
- store the file locally
- save it on [GCP](https://cloud.google.com/‎)
- save it on [S3](https://aws.amazon.com/s3/)
- send over HTTP
- send over SCP

## Notification services
Copybird currently supports the following notification services: 

- Slack
- Telegram
- AWS SES
- AWS SQS
- get notificatoin on email
- Kafka
- Nats
- Create issue in PagerDuty
- Pushbullet
- RabbitMQ
- Twilio
- Webcallback

If you would like to add additional service, please submit an issue with feature request or add it yourself and send a Pull Request.

## How to Run the tool
There are different ways you can use this tool: 

### Run locally
First get the source code on your machine
```
go get -u github.com/copybird/copybird
```
Then run it with `go run main.go` to see helpers for various optional parameters

### Run with Docker
Run `docker run copybird/copybird` to see the available optional parameters

### Use Backup Custom Controller/Operator for k8s

First create custom resource definition in your cluster: 
```
kubectl apply -f operator/crd/crd.yaml
```

To run the controller:
``` 
go run main.go operator
```

And then in a separate shell, create custom resource:
```
kubectl create -f operator/example/backup-example.yaml
```
As output you get the following logs when creating, updating or deleting custom resource:
```
INFO[0000] Successfully constructed k8s client          
INFO[0000] Starting Foo controller                      
INFO[0000] Waiting for informer caches to sync          
INFO[0001] Starting workers                             
INFO[0001] Started workers               
```
You can modify example file as you wish to get proper configuration for your jobs

## Contributing
Pull requests are more than welcomed. For major changes, please open an issue first to discuss what you would like to change. 

Before submission of pull request make sure you pulled recent updates, included tests for your code that covers at least the core functionality and you submitted a desciptive issue that will be fixed with your pull request. Do not forget to mention the issue in the pull request. 


<div align="center">
    <h3>Built with Mad Devs support for the community</h3>
    <a href="https://maddevs.io"><img height="100px" src ="docs/md-logo.png" /></a>
</div>

