# Auto Backup
To Backup MongoDB Database and Files into AWS S3 Bucket

How to use
---
#### Backup Database
It will dump your MongoDB database with archive and copy it into your local folder and S3 Bucket  
Here is step by step:
1. Change config file inside configs/configs.json with appropriate information about your database and S3 Buecket
2. If you want to upload into your S3 Bucket folder just add `/yourfoldername/` into `bucket` key
3. `retentionday` key used to set your object retention in S3 bucket for N days
3. To test it just run `go run main.go`.
4. To run it in any scheduler just build it and add parameter `-config=yourconfigfilelocation`.

#### Backup Files
It will copy your files into S3 Bucket  
Here is step by step:
1. Change config file inside configs/configs.json with appropriate information about your files directory location and S3 Bucket
2. It will only copy files in the same day when this apps run
3. To test it just run `go run main.go`.
4. To run it in any scheduler just build it and add parameter `-config=yourconfigfilelocation`.

