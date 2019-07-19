# Auto Backup
To Backup MongoDB Database and Files into AWS S3 Bucket

How to use
---
#### Backup Database
It will dump your MongoDB database with archive and copy it into your local folder and S3 Bucket  
It's located in `database` directory and following is brief description about the config file:
1. Change config file inside `configs/configs.json` with appropriate information about your database and S3 Bucket
2.	`uri` is connection string used to connect into your database
3.	`archivename` is your archive name that contains all dump files from your database
4.	`archivesuffix_dateformat` used to add suffix into archive name. It can be used as a flag to determine your dump time. This format 20060102 is similar to yyyyMMdd format and 150405 is similar to HHmmss
5.	`destpath` is destination directory for your dump files as archive file
6.	`retentionday` key use to set your object retention in S3 bucket for N days
For example: 
If this apps run in 15 July 2019 and `retentionday: 7`, it will keep object between 8 July 2019 - 15 July 2019 in particular folder, and all objects before 8 July 2019 will be deleted
7.	If you don't want to use retention feature just set the `retentionday: 0`
8.	`region` is your S3 Bucket region
9.	`bucket` is your S3 Bucket name
10.	If you want to upload into your S3 Bucket folder just add your folder name into `folder` key with structure `folder` or `folder/subfolder` etc, don't worry if you write it as `/folder` or `folder/` or `/folder/` this apps will trim it
11. To test it just run `go run main.go`.
12.	To run it in scheduler just build it and add parameter `-config=yourconfigfilelocation`.

#### Backup Files
It will copy your files into S3 Bucket  
It's located in `file` directory and following is brief description about the config file:
1. Change config file inside `configs/configs.json` with appropriate information about your files directory location and S3 Bucket
2. `dirpath` is your upload directory location
3.	When you set `initialrun: true`, it will copy all files inside particular directory, after that it will set `initialrun: false` automatically
4.	When you set `initialrun: false`, it will copy files in the same day when this apps run
5.	`backuptype`  key use to determine whether you want to copy files or folders set
6.	If you want to copy folders you can also determine archiving method for this particular directory on `archivemethod`  key, the options are `zip`, `tar`  or `tar.gz`  
7.	`region` is your S3 Bucket region
8.	`bucket` is your S3 Bucket name
9.	If you want to upload into your S3 Bucket folder just add your folder name into `folder` key with structure `folder` or `folder/subfolder` etc, don't worry if you write it as `/folder` or `folder/` or `/folder/` this apps will trim it
10. To test it just run `go run main.go`.
11.	To run it in scheduler just build it and add parameter `-config=yourconfigfilelocation`.


