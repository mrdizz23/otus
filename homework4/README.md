## Создание пользователя и базы данных

```CREATE USER medical_owner with password 'password'```  
```CREATE DATABASE medical_center_owner;```

<img width="771" height="501" alt="image" src="https://github.com/user-attachments/assets/3b29a6f5-48bf-401d-8326-42ca82630593" />

__В текущих реалиях дешевезны SSD дисков не вижу практического смысла для выделения каких-либо таблиц в отдельные табличные пространства, потому чисто для практики запросы на создание и выделение в него таблиц выглядели бы так:__

```
CREATE TABLESPACE spase1;
CREATE TABLE IF NOT EXISTS "doctors" (
	"id" serial NOT NULL UNIQUE,
	"name" varchar(255) NOT NULL,
	"surname" varchar(255) NOT NULL,
	"patronymic" varchar(255),
	"sex" varchar(255) NOT NULL,
	"speciality" varchar(255) NOT NULL,
	"is_active" boolean NOT NULL,
	PRIMARY KEY ("id")
) TABLESPACE space1;
```

## Создание схемы

```CREATE SCHEMA version1```

<img width="811" height="338" alt="image" src="https://github.com/user-attachments/assets/8279b505-9935-4be0-92ca-e3a15d603ed3" />

## Создание таблиц

__В виду большого количества таблиц, запросы на их создание я вынес в отдельный файл (аналогичен medical_center.sql из прошлых ДЗ) из применил через psql__

<img width="857" height="445" alt="image" src="https://github.com/user-attachments/assets/d17045ee-2412-4198-bcab-49eb3979fe65" />
