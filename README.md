# Function 
1. backup 
2. restore
3. create table
4. insert data?
5. check postgresql version 


read config file to handle the postgres settings.
dbtool backup -n local -c config.yaml


dbtool restore -n remote -c config.yaml -f temp.dump


dbtool build -n remote -c config.yaml -d sql/create
