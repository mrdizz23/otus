> Реализация домашнего задания произведения с помощью Docker

0. Создаю сеть в режиме bridge для общения контейнеров друг с другом по DNS

```
[dizz@MUR-PC-3009-B2C ~]$ docker network create -d bridge mysql
ecedd6dc2e7332ceb6d536d38f14efccfcda0f2f34c8d3d7e4761ae853434898
```

1. Готовлю конфиг для запуска контейнеров c разницей в server_id - 1,2,3 для каждого контейнера соответственно

```
[dizz@MUR-PC-3009-B2C ~]$ cat /home/dizz/idb_cluster_node1/config-file.cnf
[mysqld]
server_id = 1
binlog_transaction_dependency_tracking = WRITESET
gtid_mode = ON
enforce_gtid_consistency = ON
```

2. Создаю контейнеры в только что созданной сети и примапливаю им подготовленные конфиги

```
[dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mysql-node1 -p 3306:3306 --network=mysql -v /home/dizz/idb_cluster_node1:/etc/my.cnf.d -e MYSQL_ROOT_PASSWORD=123456 -e LANG=C.UTF-8 -d percona/percona-server:8.0
8eb4e367a2e2772845b173c19ec7a2289d2709ac231589b488a1e323cd0fc35d
[dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mysql-node2 -p 3307:3306 --network=mysql -v /home/dizz/idb_cluster_node2:/etc/my.cnf.d -e MYSQL_ROOT_PASSWORD=123456 -e LANG=C.UTF-8 -d percona/percona-server:8.0
efef5be9c4a0f2e730701bd5b9f4ee9e86b19cafecffe087fe227722de9b696e
[dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mysql-node3 -p 3308:3306 --network=mysql -v /home/dizz/idb_cluster_node3:/etc/my.cnf.d -e MYSQL_ROOT_PASSWORD=123456 -e LANG=C.UTF-8 -d percona/percona-server:8.0
df9dbc22e21bbd0ed4e65645b01af1df294d284096cd91a75283e056cce4e5e7
[dizz@MUR-PC-3009-B2C ~]$ docker ps
CONTAINER ID   IMAGE                        COMMAND                  CREATED          STATUS          PORTS                                         NAMES
df9dbc22e21b   percona/percona-server:8.0   "/docker-entrypoint.…"   3 seconds ago    Up 2 seconds    0.0.0.0:3308->3306/tcp, [::]:3308->3306/tcp   mysql-node3
efef5be9c4a0   percona/percona-server:8.0   "/docker-entrypoint.…"   25 seconds ago   Up 23 seconds   0.0.0.0:3307->3306/tcp, [::]:3307->3306/tcp   mysql-node2
8eb4e367a2e2   percona/percona-server:8.0   "/docker-entrypoint.…"   38 seconds ago   Up 37 seconds   0.0.0.0:3306->3306/tcp, [::]:3306->3306/tcp   mysql-node1
```
2. Настройку кластера буду производить с помощью mysqlsh на первой ноде

```
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-node1 mysqlsh
Cannot set LC_ALL to locale en_US.UTF-8: No such file or directory
MySQL Shell 8.0.43

Copyright (c) 2016, 2025, Oracle and/or its affiliates.
Oracle is a registered trademark of Oracle Corporation and/or its affiliates.
Other names may be trademarks of their respective owners.

Type '\help' or '\?' for help; '\quit' to exit.
 MySQL  JS > dba.configureInstance('root@mysql-node1',{clusterAdmin: "clusteradmin@'%'", clusterAdminPassword: '123456'});
Please provide the password for 'root@mysql-node1': ******
Save password for 'root@mysql-node1'? [Y]es/[N]o/Ne[v]er (default No): Y
Configuring local MySQL instance listening at port 3306 for use in an InnoDB cluster...

This instance reports its own address as 8eb4e367a2e2:3306
Clients and other cluster members will communicate with it through this address by default. If this is not correct, the report_host MySQL system variable should be changed.

applierWorkerThreads will be set to the default value of 4.

The instance '8eb4e367a2e2:3306' is valid to be used in an InnoDB cluster.

Creating user clusteradmin@%.
Account clusteradmin@% was successfully created.

The instance '8eb4e367a2e2:3306' is already ready to be used in an InnoDB cluster.

Successfully enabled parallel appliers.
 MySQL  JS > dba.configureInstance('root@mysql-node2',{clusterAdmin: "clusteradmin@'%'", clusterAdminPassword: '123456'});
Please provide the password for 'root@mysql-node2': ******
Save password for 'root@mysql-node2'? [Y]es/[N]o/Ne[v]er (default No): Y
Configuring MySQL instance at efef5be9c4a0:3306 for use in an InnoDB cluster...

This instance reports its own address as efef5be9c4a0:3306
Clients and other cluster members will communicate with it through this address by default. If this is not correct, the report_host MySQL system variable should be changed.

applierWorkerThreads will be set to the default value of 4.

The instance 'efef5be9c4a0:3306' is valid to be used in an InnoDB cluster.

Creating user clusteradmin@%.
Account clusteradmin@% was successfully created.

The instance 'efef5be9c4a0:3306' is already ready to be used in an InnoDB cluster.

Successfully enabled parallel appliers.
 MySQL  JS > dba.configureInstance('root@mysql-node3',{clusterAdmin: "clusteradmin@'%'", clusterAdminPassword: '123456'});
Please provide the password for 'root@mysql-node3': ******
Save password for 'root@mysql-node3'? [Y]es/[N]o/Ne[v]er (default No): Y
Configuring MySQL instance at df9dbc22e21b:3306 for use in an InnoDB cluster...

This instance reports its own address as df9dbc22e21b:3306
Clients and other cluster members will communicate with it through this address by default. If this is not correct, the report_host MySQL system variable should be changed.

applierWorkerThreads will be set to the default value of 4.

The instance 'df9dbc22e21b:3306' is valid to be used in an InnoDB cluster.

Creating user clusteradmin@%.
Account clusteradmin@% was successfully created.

The instance 'df9dbc22e21b:3306' is already ready to be used in an InnoDB cluster.

Successfully enabled parallel appliers.
 MySQL  JS >
```

3. Подключаюсь к инстансу на node1, создаю кластер и добавляю в него ноды 2 и 3 через [I]ncremental recovery, посколько инстанс у меня пустой, а также [C]lone рестартит инснтасы, что недопустимо в моей конфигурации через Docker

```
 MySQL  JS > shell.connect('clusteradmin@mysql-node1:3306');
Creating a session to 'clusteradmin@mysql-node1:3306'
Please provide the password for 'clusteradmin@mysql-node1:3306': ******
Save password for 'clusteradmin@mysql-node1:3306'? [Y]es/[N]o/Ne[v]er (default No): Y
Fetching schema names for auto-completion... Press ^C to stop.
Your MySQL connection id is 11
Server version: 8.0.43-34 Percona Server (GPL), Release 34, Revision e2841f91
No default schema selected; type \use <schema> to set one.
<ClassicSession:clusteradmin@mysql-node1:3306>
 MySQL  mysql-node1:3306 ssl  JS > cluster = dba.createCluster('my_first_cluster');
A new InnoDB Cluster will be created on instance '1ed9eef1f6e2:3306'.

Validating instance configuration at mysql-node1:3306...

This instance reports its own address as 1ed9eef1f6e2:3306

Instance configuration is suitable.
NOTE: Group Replication will communicate with other members using '1ed9eef1f6e2:3306'. Use the localAddress option to override.

* Checking connectivity and SSL configuration...

Creating InnoDB Cluster 'my_first_cluster' on '1ed9eef1f6e2:3306'...

Adding Seed Instance...
Cluster successfully created. Use Cluster.addInstance() to add MySQL instances.
At least 3 instances are needed for the cluster to be able to withstand up to
one server failure.

<Cluster:my_first_cluster>
 MySQL  mysql-node1:3306 ssl  JS > cluster.addInstance('clusteradmin@mysql-node2:3306');

NOTE: The target instance 'e518101423ab:3306' has not been pre-provisioned (GTID set is empty). The Shell is unable to decide whether incremental state recovery can correctly provision it.
The safest and most convenient way to provision a new instance is through automatic clone provisioning, which will completely overwrite the state of 'e518101423ab:3306' with a physical snapshot from an existing cluster member. To use this method by default, set the 'recoveryMethod' option to 'clone'.

The incremental state recovery may be safely used if you are sure all updates ever executed in the cluster were done with GTIDs enabled, there are no purged transactions and the new instance contains the same GTID set as the cluster or a subset of it. To use this method by default, set the 'recoveryMethod' option to 'incremental'.


Please select a recovery method [C]lone/[I]ncremental recovery/[A]bort (default Clone): I
Validating instance configuration at mysql-node2:3306...

This instance reports its own address as e518101423ab:3306

Instance configuration is suitable.
NOTE: Group Replication will communicate with other members using 'e518101423ab:3306'. Use the localAddress option to override.

* Checking connectivity and SSL configuration...
A new instance will be added to the InnoDB Cluster. Depending on the amount of
data on the cluster this might take from a few seconds to several hours.

Adding instance to the cluster...

Monitoring recovery process of the new cluster member. Press ^C to stop monitoring and let it continue in background.
Incremental state recovery is now in progress.

* Waiting for distributed recovery to finish...
NOTE: 'e518101423ab:3306' is being recovered from '1ed9eef1f6e2:3306'
* Distributed recovery has finished

The instance 'e518101423ab:3306' was successfully added to the cluster.

 MySQL  mysql-node1:3306 ssl  JS > cluster.addInstance('clusteradmin@mysql-node3:3306');

NOTE: The target instance '7c32a9d06978:3306' has not been pre-provisioned (GTID set is empty). The Shell is unable to decide whether incremental state recovery can correctly provision it.
The safest and most convenient way to provision a new instance is through automatic clone provisioning, which will completely overwrite the state of '7c32a9d06978:3306' with a physical snapshot from an existing cluster member. To use this method by default, set the 'recoveryMethod' option to 'clone'.

The incremental state recovery may be safely used if you are sure all updates ever executed in the cluster were done with GTIDs enabled, there are no purged transactions and the new instance contains the same GTID set as the cluster or a subset of it. To use this method by default, set the 'recoveryMethod' option to 'incremental'.


Please select a recovery method [C]lone/[I]ncremental recovery/[A]bort (default Clone): I
Validating instance configuration at mysql-node3:3306...

This instance reports its own address as 7c32a9d06978:3306

Instance configuration is suitable.
NOTE: Group Replication will communicate with other members using '7c32a9d06978:3306'. Use the localAddress option to override.

* Checking connectivity and SSL configuration...
A new instance will be added to the InnoDB Cluster. Depending on the amount of
data on the cluster this might take from a few seconds to several hours.

Adding instance to the cluster...

Monitoring recovery process of the new cluster member. Press ^C to stop monitoring and let it continue in background.
Incremental state recovery is now in progress.

* Waiting for distributed recovery to finish...
NOTE: '7c32a9d06978:3306' is being recovered from '1ed9eef1f6e2:3306'
* Distributed recovery has finished

The instance '7c32a9d06978:3306' was successfully added to the cluster.

 MySQL  mysql-node1:3306 ssl  JS > cluster.status();
{
    "clusterName": "my_first_cluster",
    "defaultReplicaSet": {
        "name": "default",
        "primary": "1ed9eef1f6e2:3306",
        "ssl": "REQUIRED",
        "status": "OK",
        "statusText": "Cluster is ONLINE and can tolerate up to ONE failure.",
        "topology": {
            "1ed9eef1f6e2:3306": {
                "address": "1ed9eef1f6e2:3306",
                "memberRole": "PRIMARY",
                "mode": "R/W",
                "readReplicas": {},
                "replicationLag": "applier_queue_applied",
                "role": "HA",
                "status": "ONLINE",
                "version": "8.0.43"
            },
            "7c32a9d06978:3306": {
                "address": "7c32a9d06978:3306",
                "memberRole": "SECONDARY",
                "mode": "R/O",
                "readReplicas": {},
                "replicationLag": "applier_queue_applied",
                "role": "HA",
                "status": "ONLINE",
                "version": "8.0.43"
            },
            "e518101423ab:3306": {
                "address": "e518101423ab:3306",
                "memberRole": "SECONDARY",
                "mode": "R/O",
                "readReplicas": {},
                "replicationLag": "applier_queue_applied",
                "role": "HA",
                "status": "ONLINE",
                "version": "8.0.43"
            }
        },
        "topologyMode": "Single-Primary"
    },
    "groupInformationSourceMember": "1ed9eef1f6e2:3306"
}
 MySQL  mysql-node1:3306 ssl  JS >

```

4. Создаю тестовую базу данных для проверки репликации, попутно выдавая права пользователю clusteradmin

```
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-node1 mysql -u root -p123456
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 66
Server version: 8.0.43-34 Percona Server (GPL), Release 34, Revision e2841f91

Copyright (c) 2009-2025 Percona LLC and/or its affiliates
Copyright (c) 2000, 2025, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show databases;
+-------------------------------+
| Database                      |
+-------------------------------+
| information_schema            |
| mysql                         |
| mysql_innodb_cluster_metadata |
| performance_schema            |
| sys                           |
+-------------------------------+
5 rows in set (0.00 sec)

mysql> create database testdb;
Query OK, 1 row affected (0.03 sec)

mysql> grant all on *.* to clusteradmin@'%';
Query OK, 0 rows affected (0.02 sec)

mysql>
Bye
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-node2 mysql -u root -p123456 -e "show databases;"
mysql: [Warning] Using a password on the command line interface can be insecure.
+-------------------------------+
| Database                      |
+-------------------------------+
| information_schema            |
| mysql                         |
| mysql_innodb_cluster_metadata |
| performance_schema            |
| sys                           |
| testdb                        |
+-------------------------------+
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-node3 mysql -u root -p123456 -e "show databases;"
mysql: [Warning] Using a password on the command line interface can be insecure.
+-------------------------------+
| Database                      |
+-------------------------------+
| information_schema            |
| mysql                         |
| mysql_innodb_cluster_metadata |
| performance_schema            |
| sys                           |
| testdb                        |
+-------------------------------+
```
5. Поднимаю еще один контейнер с mysql-router и проверяю работу балансировщика по разным портам

```
[dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mysql-router --network=mysql -p 6446:6446 -p 6447:6447 -e MYSQL_HOST=mysql-node1 -e MYSQL_PORT=3306 -e MYSQL_USER=root -e MYSQL_PASSWORD=123456 -e MYSQL_INNODB_CLUSTER_MEMBERS=3 -d mysql/mysql-router
30efd01fa1adcc0072d5066482d66a4e3b1ea1474e3874d72e8fb8b33142aac2
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-node1 bash
[mysql@1ed9eef1f6e2 /]$ mysql -u root -h mysql-router --port 6446 -p123456 -e "select @@hostname;"
mysql: [Warning] Using a password on the command line interface can be insecure.
+--------------+
| @@hostname   |
+--------------+
| 1ed9eef1f6e2 |
+--------------+
[mysql@1ed9eef1f6e2 /]$ mysql -u root -h mysql-router --port 6447 -p123456 -e "select @@hostname;"
mysql: [Warning] Using a password on the command line interface can be insecure.
+--------------+
| @@hostname   |
+--------------+
| 7c32a9d06978 |
+--------------+
[mysql@1ed9eef1f6e2 /]$ mysql -u root -h mysql-router --port 6447 -p123456 -e "select @@hostname;"
mysql: [Warning] Using a password on the command line interface can be insecure.
+--------------+
| @@hostname   |
+--------------+
| e518101423ab |
+--------------+
[mysql@1ed9eef1f6e2 /]$
```
6. Ручной фейловер, смотрю текйщий статус инстансов кластера, останавливаю контейнер mysql-node3 и снова получаю статус

```
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-node1 mysqlsh root@mysql-node1:3306
Cannot set LC_ALL to locale en_US.UTF-8: No such file or directory
MySQL Shell 8.0.43

Copyright (c) 2016, 2025, Oracle and/or its affiliates.
Oracle is a registered trademark of Oracle Corporation and/or its affiliates.
Other names may be trademarks of their respective owners.

Type '\help' or '\?' for help; '\quit' to exit.
Creating a session to 'root@mysql-node1:3306'
Fetching schema names for auto-completion... Press ^C to stop.
Your MySQL connection id is 1488
Server version: 8.0.43-34 Percona Server (GPL), Release 34, Revision e2841f91
No default schema selected; type \use <schema> to set one.
 MySQL  mysql-node1:3306 ssl  JS > cluster = dba.getCluster();
<Cluster:my_first_cluster>
 MySQL  mysql-node1:3306 ssl  JS > cluster.status();
{
    "clusterName": "my_first_cluster",
    "defaultReplicaSet": {
        "name": "default",
        "primary": "1ed9eef1f6e2:3306",
        "ssl": "REQUIRED",
        "status": "OK",
        "statusText": "Cluster is ONLINE and can tolerate up to ONE failure.",
        "topology": {
            "1ed9eef1f6e2:3306": {
                "address": "1ed9eef1f6e2:3306",
                "memberRole": "PRIMARY",
                "mode": "R/W",
                "readReplicas": {},
                "replicationLag": "applier_queue_applied",
                "role": "HA",
                "status": "ONLINE",
                "version": "8.0.43"
            },
            "7c32a9d06978:3306": {
                "address": "7c32a9d06978:3306",
                "memberRole": "SECONDARY",
                "mode": "R/O",
                "readReplicas": {},
                "replicationLag": "applier_queue_applied",
                "role": "HA",
                "status": "ONLINE",
                "version": "8.0.43"
            },
            "e518101423ab:3306": {
                "address": "e518101423ab:3306",
                "memberRole": "SECONDARY",
                "mode": "R/O",
                "readReplicas": {},
                "replicationLag": "applier_queue_applied",
                "role": "HA",
                "status": "ONLINE",
                "version": "8.0.43"
            }
        },
        "topologyMode": "Single-Primary"
    },
    "groupInformationSourceMember": "1ed9eef1f6e2:3306"
}
 MySQL  mysql-node1:3306 ssl  JS >
Bye!
[dizz@MUR-PC-3009-B2C ~]$ docker stop mysql-node3
mysql-node3
[dizz@MUR-PC-3009-B2C ~]$ docker exec -it mysql-node1 mysqlsh root@mysql-node1:3306
Cannot set LC_ALL to locale en_US.UTF-8: No such file or directory
MySQL Shell 8.0.43

Copyright (c) 2016, 2025, Oracle and/or its affiliates.
Oracle is a registered trademark of Oracle Corporation and/or its affiliates.
Other names may be trademarks of their respective owners.

Type '\help' or '\?' for help; '\quit' to exit.
Creating a session to 'root@mysql-node1:3306'
Fetching schema names for auto-completion... Press ^C to stop.
Your MySQL connection id is 1586
Server version: 8.0.43-34 Percona Server (GPL), Release 34, Revision e2841f91
No default schema selected; type \use <schema> to set one.
 MySQL  mysql-node1:3306 ssl  JS > cluster = dba.getCluster();
<Cluster:my_first_cluster>
 MySQL  mysql-node1:3306 ssl  JS > cluster.status();
{
    "clusterName": "my_first_cluster",
    "defaultReplicaSet": {
        "name": "default",
        "primary": "1ed9eef1f6e2:3306",
        "ssl": "REQUIRED",
        "status": "OK_NO_TOLERANCE_PARTIAL",
        "statusText": "Cluster is NOT tolerant to any failures. 1 member is not active.",
        "topology": {
            "1ed9eef1f6e2:3306": {
                "address": "1ed9eef1f6e2:3306",
                "memberRole": "PRIMARY",
                "mode": "R/W",
                "readReplicas": {},
                "replicationLag": "applier_queue_applied",
                "role": "HA",
                "status": "ONLINE",
                "version": "8.0.43"
            },
            "7c32a9d06978:3306": {
                "address": "7c32a9d06978:3306",
                "memberRole": "SECONDARY",
                "mode": "n/a",
                "readReplicas": {},
                "role": "HA",
                "shellConnectError": "MySQL Error 2005: Could not open connection to '7c32a9d06978:3306': Unknown MySQL server host '7c32a9d06978' (-2)",
                "status": "(MISSING)"
            },
            "e518101423ab:3306": {
                "address": "e518101423ab:3306",
                "memberRole": "SECONDARY",
                "mode": "R/O",
                "readReplicas": {},
                "replicationLag": "applier_queue_applied",
                "role": "HA",
                "status": "ONLINE",
                "version": "8.0.43"
            }
        },
        "topologyMode": "Single-Primary"
    },
    "groupInformationSourceMember": "1ed9eef1f6e2:3306"
}
 MySQL  mysql-node1:3306 ssl  JS >
```
