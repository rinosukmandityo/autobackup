# Auto Backup
To Backup MongoDB Database and Files into AWS S3 Bucket

How to use
---
#### Backup Database
It will dump your MongoDB database with archive and copy it into your local folder and S3 Bucket  
Here is step by step:
1. Change config file inside configs/configs.json with appropriate information about your database and S3 Bucket
2. If you want to upload into your S3 Bucket folder just add your folder name into `folder` key with structure `folder` or `folder/subfolder` etc, don't worry if you write it as `/folder` or `folder/` or `/folder/` this apps will trim it
3. `retentionday` key use to set your object retention in S3 bucket for N days
4. If you use retention feature please put it into separate folder because it will delete all data inside that folder if last modified object meets `retentionday` value  
<b>For example</b>  
if this apps run in 15 July 2019 and `retentionday` value is 7, it will keep object between 8 July 2019 - 15 July 2019 in particular folder, and all objects before 8 July 2019 will be deleted
5. If you don't want to use retention feature just set the `retentionday: 0`
6. To test it just run `go run main.go`.
7. To run it in any scheduler just build it and add parameter `-config=yourconfigfilelocation`.

#### Backup Files
It will copy your files into S3 Bucket  
Here is step by step:
1. Change config file inside configs/configs.json with appropriate information about your files directory location and S3 Bucket
2. If you want to upload into your S3 Bucket folder just add your folder name into `folder` key with structure `folder` or `folder/subfolder` etc
3. When you set `initialrun: true`, it will copy all files inside particular directory, after that it will set `initialrun: false` automatically
4. When you set `initialrun: false`, it will only copy files in the same day when this apps run
5. `backuptype` key use to determine whether you want to copy files or folders set
6. If you want to copy folders you can also determine archiving method for this particular directory on `archivemethod` key, the options are `zip`, `tar` or `tar.gz` 
7. To test it just run `go run main.go`.
8. To run it in any scheduler just build it and add parameter `-config=yourconfigfilelocation`.

