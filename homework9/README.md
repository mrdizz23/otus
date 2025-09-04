## Физическая репликация

1. Правлю init.sql, вставляю в него БД sbtest для будущего проведения тестирования с помощью sysbench

```
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ ll
total 28
drwxrwxr-x  4 dizz dizz 4096 Sep  4 15:04 .
drwx------ 18 dizz dizz 4096 Sep  4 15:04 ..
drwxrwxr-x  2 dizz dizz 4096 Sep  4 13:44 custom.conf
-rw-rw-r--  1 dizz dizz  339 Sep  4 13:30 docker-compose.yml
drwxrwxr-x  8 dizz dizz 4096 Sep  4 13:30 .git
-rw-rw-r--  1 dizz dizz   36 Sep  4 15:04 init.sql
-rw-rw-r--  1 dizz dizz  416 Sep  4 13:30 README.md
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ cat init.sql
CREATE database sbtest;
USE sbtest;
```

2. Добавляю в custom.conf еще один файл с кастомной конфигурацией

```
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ ll custom.conf/
total 16
drwxrwxr-x 2 dizz dizz 4096 Sep  4 13:44 .
drwxrwxr-x 4 dizz dizz 4096 Sep  4 15:04 ..
-rw-rw-r-- 1 dizz dizz   43 Sep  4 13:44 inno.cnf
-rw-rw-r-- 1 dizz dizz   62 Sep  4 13:30 my.cnf
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ cat custom.conf/inno.cnf
[mysqld]
innodb-buffer-pool-size=536870912
```

3. Запускаю сервис и проверяю, что инстанс доступен, база создалась, параметр применился (дефолтное значение 128Мб)

```
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ docker compose up otusdb -d
[+] Running 1/1
 ✔ Container otus-mysql-docker-otusdb-1  Started                                                                                                                                                                                       0.9s
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ docker ps
CONTAINER ID   IMAGE          COMMAND                  CREATED         STATUS         PORTS                                         NAMES
c3052b64dfac   mysql:8.0.15   "docker-entrypoint.s…"   4 seconds ago   Up 4 seconds   0.0.0.0:3309->3306/tcp, [::]:3309->3306/tcp   otus-mysql-docker-otusdb-1
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ mysql -u root -p12345 -h 127.0.0.1 --port=3309 --protocol=tcp sbtest
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 9
Server version: 8.0.15 MySQL Community Server - GPL

Copyright (c) 2000, 2021, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sbtest             |
| sys                |
+--------------------+
5 rows in set (0.01 sec)

mysql> SELECT DATABASE();
+------------+
| DATABASE() |
+------------+
| sbtest     |
+------------+
1 row in set (0.00 sec)

mysql> select @@innodb_buffer_pool_size/1024/1024;
+-------------------------------------+
| @@innodb_buffer_pool_size/1024/1024 |
+-------------------------------------+
|                        512.00000000 |
+-------------------------------------+
1 row in set (0.00 sec)

mysql>
```

4. Проведение тестирования. Наполняю тестовую базу данными, предварительно проверив размер текущей пустой БД

```
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ docker compose exec otusdb du -sh /var/lib/mysql/sbtest/
4.0K    /var/lib/mysql/sbtest/
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ sysbench --db-driver=mysql --mysql-user=root --mysql-db=sbtest --mysql-host=127.0.0.1 --mysql-port=3309 --mysql-password=12345 --tables=16 --table-size=10000 /usr/share/sysbench/oltp_read_write.lua prepare
sysbench 1.0.20 (using system LuaJIT 2.1.0-beta3)

Creating table 'sbtest1'...
Inserting 10000 records into 'sbtest1'
Creating a secondary index on 'sbtest1'...
Creating table 'sbtest2'...
Inserting 10000 records into 'sbtest2'
Creating a secondary index on 'sbtest2'...
Creating table 'sbtest3'...
Inserting 10000 records into 'sbtest3'
Creating a secondary index on 'sbtest3'...
Creating table 'sbtest4'...
Inserting 10000 records into 'sbtest4'
Creating a secondary index on 'sbtest4'...
Creating table 'sbtest5'...
Inserting 10000 records into 'sbtest5'
Creating a secondary index on 'sbtest5'...
Creating table 'sbtest6'...
Inserting 10000 records into 'sbtest6'
Creating a secondary index on 'sbtest6'...
Creating table 'sbtest7'...
Inserting 10000 records into 'sbtest7'
Creating a secondary index on 'sbtest7'...
Creating table 'sbtest8'...
Inserting 10000 records into 'sbtest8'
Creating a secondary index on 'sbtest8'...
Creating table 'sbtest9'...
Inserting 10000 records into 'sbtest9'
Creating a secondary index on 'sbtest9'...
Creating table 'sbtest10'...
Inserting 10000 records into 'sbtest10'
Creating a secondary index on 'sbtest10'...
Creating table 'sbtest11'...
Inserting 10000 records into 'sbtest11'
Creating a secondary index on 'sbtest11'...
Creating table 'sbtest12'...
Inserting 10000 records into 'sbtest12'
Creating a secondary index on 'sbtest12'...
Creating table 'sbtest13'...
Inserting 10000 records into 'sbtest13'
Creating a secondary index on 'sbtest13'...
Creating table 'sbtest14'...
Inserting 10000 records into 'sbtest14'
Creating a secondary index on 'sbtest14'...
Creating table 'sbtest15'...
Inserting 10000 records into 'sbtest15'
Creating a secondary index on 'sbtest15'...
Creating table 'sbtest16'...
Inserting 10000 records into 'sbtest16'
Creating a secondary index on 'sbtest16'...
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ docker compose exec otusdb du -sh /var/lib/mysql/sbtest/
160M    /var/lib/mysql/sbtest/
```

5. Запуск и результаты теста

```
[dizz@MUR-PC-3009-B2C otus-mysql-docker]$ sysbench --db-driver=mysql --mysql-user=root --mysql-db=sbtest --mysql-host=127.0.0.1 --mysql-port=3309 --mysql-password=12345 --tables=16 --table-size=10000 --threads=4 --time=10 --events=0 --report-interval=1 /usr/share/sysbench/oltp_read_write.lua run
sysbench 1.0.20 (using system LuaJIT 2.1.0-beta3)

Running the test with following options:
Number of threads: 4
Report intermediate results every 1 second(s)
Initializing random number generator from current time


Initializing worker threads...

Threads started!

[ 1s ] thds: 4 tps: 56.61 qps: 1198.68 (r/w/o: 847.12/234.37/117.19) lat (ms,95%): 153.02 err/s: 0.00 reconn/s: 0.00
[ 2s ] thds: 4 tps: 64.36 qps: 1296.26 (r/w/o: 902.05/265.49/128.72) lat (ms,95%): 99.33 err/s: 0.00 reconn/s: 0.00
[ 3s ] thds: 4 tps: 117.99 qps: 2354.78 (r/w/o: 1650.84/467.96/235.98) lat (ms,95%): 61.08 err/s: 0.00 reconn/s: 0.00
[ 4s ] thds: 4 tps: 111.01 qps: 2206.11 (r/w/o: 1541.07/443.02/222.01) lat (ms,95%): 59.99 err/s: 0.00 reconn/s: 0.00
[ 5s ] thds: 4 tps: 131.96 qps: 2631.24 (r/w/o: 1846.46/520.85/263.92) lat (ms,95%): 55.82 err/s: 0.00 reconn/s: 0.00
[ 6s ] thds: 4 tps: 128.03 qps: 2573.52 (r/w/o: 1799.36/518.10/256.05) lat (ms,95%): 62.19 err/s: 0.00 reconn/s: 0.00
[ 7s ] thds: 4 tps: 117.00 qps: 2343.93 (r/w/o: 1643.95/465.99/233.99) lat (ms,95%): 57.87 err/s: 0.00 reconn/s: 0.00
[ 8s ] thds: 4 tps: 152.01 qps: 3044.25 (r/w/o: 2128.18/612.05/304.03) lat (ms,95%): 35.59 err/s: 0.00 reconn/s: 0.00
[ 9s ] thds: 4 tps: 125.99 qps: 2494.76 (r/w/o: 1746.83/495.95/251.98) lat (ms,95%): 59.99 err/s: 0.00 reconn/s: 0.00
[ 10s ] thds: 4 tps: 140.01 qps: 2782.25 (r/w/o: 1942.17/560.05/280.02) lat (ms,95%): 59.99 err/s: 0.00 reconn/s: 0.00
SQL statistics:
    queries performed:
        read:                            16086
        write:                           4596
        other:                           2298
        total:                           22980
    transactions:                        1149   (114.58 per sec.)
    queries:                             22980  (2291.55 per sec.)
    ignored errors:                      0      (0.00 per sec.)
    reconnects:                          0      (0.00 per sec.)

General statistics:
    total time:                          10.0271s
    total number of events:              1149

Latency (ms):
         min:                                   17.84
         avg:                                   34.87
         max:                                  163.76
         95th percentile:                       66.84
         sum:                                40067.01

Threads fairness:
    events (avg/stddev):           287.2500/1.09
    execution time (avg/stddev):   10.0168/0.01
```
