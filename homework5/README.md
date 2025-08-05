## Запрос с регулярным выражением

Получаем список всех пациентов, фамилия которых начинается на 'Пе' или 'По':
```
SELECT * FROM doctors WHERE surname SIMILAR TO 'П(е|о)%';
```
## Запросы с использованием LEFT JOIN и INNER JOIN

Получаем строки, которые есть в обеих таблицах по конкретному доктору:
```
SELECT * FROM bookings as b 
INNER JOIN doctors as p ON (p.id = b.doctors_id);
```  

Получаем строки только из таблицы bookings, doctors_id которыех совпадает с конкретным доктором:
```
SELECT * FROM bookings as b
LEFT JOIN doctors as p ON (p.id = b.doctors_id);
```

## Запрос на добавление данных с выводом информации о добавленных строках

```
INSERT INTO clinics (city, address)
VALUES ('Москва','Моховой переулок, дом 5')
RETURNING id;
```

## Запрос с обновлением данных, используя UPDATE FROM

```
UPDATE bookings
SET appointment_time = c.call_time
FROM patients as p LEFT JOIN calls as c ON (p.phone = c.phone)
WHERE p.name || ' ' || p.surname = 'Олег Петров';
```

## Запрос для удаления данных, используя JOIN с помощью USING

```
DELETE FROM bookings as b
USING doctors as d
WHERE d.id = b.doctor_id
AND NOT d.is_active;
```

## Пример использования утилиты COPY

Отчет оконченных приемов за последний месяц в csv файл
```
COPY (
	SELECT * FROM bookings as b
	INNER JOIN doctors as d ON d.id = b.doctor_id
	WHERE b.is_finished
	AND b.appointment_date > current_date - interval '30' day
	)
TO '/path/to/file.csv' WITH CSV;
```
