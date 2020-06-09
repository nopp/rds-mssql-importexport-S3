# AWS RDS MSSQL Import/Export using S3

![Code scanning - action](https://github.com/nopp/rds-mssql-importexport-S3/workflows/Code%20scanning%20-%20action/badge.svg)

Requisites
==========
    yum install freetds-devel
    go get github.com/minus5/gofreetds

Compiling
=========
$ go build tool.go

Configuration
=============
Configure config.json file with correct informations:

    {
    "host": rdsMSSQLhost,
    "user": rdsMSSQLuser,
    "password": rdsMSSQLpassword
    }

Running (Obs:. config.json need to stay in the same root path of tool)
=======
$ ./tool

    tool status dbName    
    tool export dbName s3Name backupName
    tool import dbName s3Name restoreName
    
