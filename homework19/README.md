> Реализация домашнего задания произведения с помощью Docker

0. Создаем сеть в режиме bridge для общения контейнеров друг с другом по DNS
1. Создаем 2 контейнера из Docker-образов mysql для мастера и реплики в той же сети на разных портах для наглядности

```
[dizz@MUR-PC-3009-B2C ~]$ docker network create -d bridge mysql
ecedd6dc2e7332ceb6d536d38f14efccfcda0f2f34c8d3d7e4761ae853434898
[dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mysql-master -p 3306:3306 --network=mysql -e MYSQL_ROOT_PASSWORD=123456 -e LANG=C.UTF-8 -d mysql:8.4.4
4c77a114892f32d6abbad9efdf298e310af3219bd14ea97350511578c4f35542
[dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mysql-replica -p 3307:3306 --network=mysql -e MYSQL_ROOT_PASSWORD=123456 -e LANG=C.UTF-8 -d mysql:8.4.4
b9b7ecb6eb22b342c6409a6d9c28b463388834ce3cc28f7ed62afb7ccb402098
[dizz@MUR-PC-3009-B2C ~]$ docker ps
CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS          PORTS                                         NAMES
b9b7ecb6eb22   mysql:8.4.4   "docker-entrypoint.s…"   6 seconds ago    Up 5 seconds    0.0.0.0:3307->3306/tcp, [::]:3307->3306/tcp   mysql-replica
4c77a114892f   mysql:8.4.4   "docker-entrypoint.s…"   17 seconds ago   Up 17 seconds   0.0.0.0:3306->3306/tcp, [::]:3306->3306/tcp   mysql-master
```
2. Создаю пользователя для репликации на мастере и даю ему привилегии

```
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-master mysql -u root -p123456
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 8
Server version: 8.4.4 MySQL Community Server - GPL

Copyright (c) 2000, 2025, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> create user repl@'%' identified with 'caching_sha2_password' by 'replpass';
Query OK, 0 rows affected (0.02 sec)

mysql> grant replication slave on *.* to repl@'%';
Query OK, 0 rows affected (0.01 sec)
```

3.
