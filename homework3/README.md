# Установка и запуск СУБД PostgreSQL

Инстанс БД запущен в docker.

![image](https://github.com/user-attachments/assets/607f7c56-dc6a-4f90-965d-2faec06379d4)


## Описание команд
> Предварительно на машине был установлен клиент и стянут образ postgresql для docker

```docker ps -a```  
Получение списка запущенных контейнеров.

```docker run --rm --name postgres -e POSTGRES_PASSWORD=pg_pass -d -p 5432:5432 postgres```  
Запуск контейнера

```ss -tlpn```  
Проверка доступности порта 5432

```psql -U postgres -h 127.0.0.1 -p 5432```  
Подключение к инстансу с вводом пароля

```create role test_role;```  
```create database test_db owner test_role;```  
Создание пользователя и базы данных
