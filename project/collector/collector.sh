#!/bin/bash

# строка для подключения к инстансу для отбора статистики
connection_string="psql -qtAX postgresql://stat_monitor@localhost:5432/postgres"

# строка для подключения к служебной БД
stats_connect="psql -qtAX postgresql://statements_owner@statements:5432/statements"

# формирование инсертов отобранной статистики во временный файл
stats=$($connection_string -f "/app/get_stats.sql" > /tmp/stat_statements.tmp)

# сброс статистики после отбора
reset=$($connection_string -c "SELECT pg_stat_statements_reset();")

# пуш отобранной статистики в служебную БД
result=$($stats_connect -f "/tmp/stat_statements.tmp")

echo 0
