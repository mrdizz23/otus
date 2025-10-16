> Реализация домашнего задания произведения с помощью Docker

0. Создаю сеть в режиме bridge для общения контейнеров друг с другом по DNS

```
[dizz@MUR-PC-3009-B2C ~]$ docker network create -d bridge mysql
ecedd6dc2e7332ceb6d536d38f14efccfcda0f2f34c8d3d7e4761ae853434898
```

1. Готовлю 2 конфига для запуска контейнеров для мастера и реплики

```
[dizz@MUR-PC-3009-B2C ~]$ cat my_master/config-file.cnf
[mysqld]
host_cache_size=0
skip-name-resolve
server_id = 1
gtid_mode = ON
enforce_gtid_consistency = ON

[dizz@MUR-PC-3009-B2C ~]$ cat my_replica/config-file.cnf
[mysqld]
host_cache_size=0
skip-name-resolve
server_id = 2
read_only = ON
gtid_mode = ON
enforce_gtid_consistency = ON
```

2. Создаю контейнера мастера в только что созданной сети и примапливаю ему подготовленный конфиг

```
[dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mysql-master -p 3306:3306 --network=mysql -v /home/dizz/my_master:/etc/my.cnf.d -e MYSQL_ROOT_PASSWORD=123456 -e LANG=C.UTF-8 -d percona/percona-server:8.0
4c77a114892f32d6abbad9efdf298e310af3219bd14ea97350511578c4f35542
```
2. Подключаюсь к инстансу мастера, смотрю GTID, создаю пользователя для репликации и даю ему нужные привелегии

```
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-master mysql -u root -p123456
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 12
Server version: 8.0.43-34 Percona Server (GPL), Release 34, Revision e2841f91

Copyright (c) 2009-2025 Percona LLC and/or its affiliates
Copyright (c) 2000, 2025, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show binary logs;
+---------------+-----------+-----------+
| Log_name      | File_size | Encrypted |
+---------------+-----------+-----------+
| binlog.000001 |       180 | No        |
| binlog.000002 |       157 | No        |
+---------------+-----------+-----------+
2 rows in set (0.00 sec)

mysql> show master status;
+---------------+----------+--------------+------------------+-------------------+
| File          | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+---------------+----------+--------------+------------------+-------------------+
| binlog.000002 |      157 |              |                  |                   |
+---------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

mysql> create user repl@'%' identified with 'caching_sha2_password' by 'replpass';
Query OK, 0 rows affected (0.02 sec)

mysql> grant replication slave on *.* to repl@'%';
Query OK, 0 rows affected (0.01 sec)

mysql> show master status;
+---------------+----------+--------------+------------------+------------------------------------------+
| File          | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set                        |
+---------------+----------+--------------+------------------+------------------------------------------+
| binlog.000002 |      693 |              |                  | ab752fe1-aa89-11f0-b453-ca57696a816d:1-2 |
+---------------+----------+--------------+------------------+------------------------------------------+
1 row in set (0.00 sec)
```

3. Аналогично стартую контейнер с репликой, подключаюсь к ней и стартую репликацию с мастера

```
[dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mysql-replica -p 3307:3306 --network=mysql -v /home/dizz/my_replica:/etc/my.cnf.d -e MYSQL_ROOT_PASSWORD=123456 -e LANG=C.UTF-8 -d percona/percona-server:8.0
d486111397efee5273f8a298f2171a30ecef6aca5818404700d2ae3eed73339d

[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-replica mysql -u root -p123456
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 9
Server version: 8.0.43-34 Percona Server (GPL), Release 34, Revision e2841f91

Copyright (c) 2009-2025 Percona LLC and/or its affiliates
Copyright (c) 2000, 2025, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> stop replica
    -> ;
Query OK, 0 rows affected, 1 warning (0.00 sec)

mysql> CHANGE REPLICATION SOURCE TO SOURCE_HOST='mysql-master', SOURCE_USER = 'repl', SOURCE_PASSWORD = 'replpass', SOURCE_AUTO_POSITION = 1, GET_SOURCE_PUBLIC_KEY = 1;
Query OK, 0 rows affected, 2 warnings (0.06 sec)

mysql> start replica;
Query OK, 0 rows affected (0.09 sec)
```

4. Смотрим статус репликации на реплике

```
mysql> show replica status\G
*************************** 1. row ***************************
             Replica_IO_State: Waiting for source to send event
                  Source_Host: mysql-master
                  Source_User: repl
                  Source_Port: 3306
                Connect_Retry: 60
              Source_Log_File: binlog.000002
          Read_Source_Log_Pos: 693
               Relay_Log_File: d486111397ef-relay-bin.000002
                Relay_Log_Pos: 903
        Relay_Source_Log_File: binlog.000002
           Replica_IO_Running: Yes
          Replica_SQL_Running: Yes
              Replicate_Do_DB:
          Replicate_Ignore_DB:
           Replicate_Do_Table:
       Replicate_Ignore_Table:
      Replicate_Wild_Do_Table:
  Replicate_Wild_Ignore_Table:
                   Last_Errno: 0
                   Last_Error:
                 Skip_Counter: 0
          Exec_Source_Log_Pos: 693
              Relay_Log_Space: 1120
              Until_Condition: None
               Until_Log_File:
                Until_Log_Pos: 0
           Source_SSL_Allowed: No
           Source_SSL_CA_File:
           Source_SSL_CA_Path:
              Source_SSL_Cert:
            Source_SSL_Cipher:
               Source_SSL_Key:
        Seconds_Behind_Source: 0
Source_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error:
               Last_SQL_Errno: 0
               Last_SQL_Error:
  Replicate_Ignore_Server_Ids:
             Source_Server_Id: 1
                  Source_UUID: ab752fe1-aa89-11f0-b453-ca57696a816d
             Source_Info_File: mysql.slave_master_info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
    Replica_SQL_Running_State: Replica has read all relay log; waiting for more updates
           Source_Retry_Count: 86400
                  Source_Bind:
      Last_IO_Error_Timestamp:
     Last_SQL_Error_Timestamp:
               Source_SSL_Crl:
           Source_SSL_Crlpath:
           Retrieved_Gtid_Set: ab752fe1-aa89-11f0-b453-ca57696a816d:1-2
            Executed_Gtid_Set: ab752fe1-aa89-11f0-b453-ca57696a816d:1-2
                Auto_Position: 1
         Replicate_Rewrite_DB:
                 Channel_Name:
           Source_TLS_Version:
       Source_public_key_path:
        Get_Source_public_key: 1
            Network_Namespace:
1 row in set (0.00 sec)
```
5. Проверяю работы репликации через создание базы данных (слева мастер, справа - реплика)

<img width="1332" height="824" alt="image" src="https://github.com/user-attachments/assets/3a3bb408-31b1-4e30-a764-c807ac31a9eb" />
