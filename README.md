# AWS RDS MSSQL Import/Export using S3

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
