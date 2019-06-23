[![](https://images.microbadger.com/badges/version/copybird/copybird.svg)](https://microbadger.com/images/copybird/copybird)
[![](https://images.microbadger.com/badges/image/copybird/copybird.svg)](https://microbadger.com/images/copybird/copybird)
[![](https://godoc.org/github.com/copybird/copybird?status.svg)](http://godoc.org/github.com/copybird/copybird)
[![GitHub release](https://img.shields.io/github/release/copybird/copybird/all.svg?style=flat-square)](https://github.com/copybird/copybird/releases)

![](https://travis-ci.org/copybird/copybird.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/copybird/copybird/badge.svg)](https://coveralls.io/github/copybird/copybird)
[![Go Report Card](https://goreportcard.com/badge/github.com/copybird/copybird)](https://goreportcard.com/report/github.com/copybird/copybird)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


# Copybird

## About

Copybird is open-source **cloud-native** universal backup tool for databases and files.

It allows you to:
1. Create database backup
2. Compress backup file
3. Encrypt backup file
4. Send it to various destinations fast and secure
5. Get notification about backup status in messagers and notification services
6. Enjoy simple backup as a service with k8s backup controller

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

## Backup as a Service (BAAS)
Run custom K8s controller with Backup custom resources

## Install & Run
Choose how to run the tool:

1. Run as a CLI tool with
```
go get -u github.com/copybird/copybird
```
2. Run with Docker
```
docker run copybird/copybird
```
3. Use k8s custom controller
```
kubectl apply -f your-backup-manifest.yaml
```

<div align="center">
    <h3>Built with Mad Devs support for the community</h3>
    <a href="https://maddevs.io"><img style="width: 100px" src ="docs/md-logo.png" /></a>
</div>

