# Upload to S3

To upload a file to S3 bucket first setup env variables 
```
export AWS_REGION=<name of your aws region>
export AWS_ACCESS_KEY_ID=<your aws access key id>
export AWS_SECRET_ACCESS_KEY=<your aws secret access key>
export AWS_BUCKET=<name of your bucket>
export FILE_PATH=<path to your file>
```
Then run 
```
go run s3.go
```