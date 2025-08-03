## Создание пользователя и базы данных

```CREATE USER medical_owner with password 'password'```
```CREATE DATABASE medical_center_owner;```

<img width="907" height="722" alt="image" src="https://github.com/user-attachments/assets/08fc9bba-cec3-4359-bea4-0b3f2119e76a" />

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

<img width="907" height="722" alt="image" src="https://github.com/user-attachments/assets/08fc9bba-cec3-4359-bea4-0b3f2119e76a" />

## Создание таблиц

__В виду большого количества таблиц, запросы на их создание я вынес в отдельный файл (аналогичен medical_center.sql из прошлых ДЗ) из применил через psql__

<img width="907" height="722" alt="image" src="https://github.com/user-attachments/assets/08fc9bba-cec3-4359-bea4-0b3f2119e76a" />