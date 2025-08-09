## Индекс по полю speciality таблицы doctors

До создания индекса  

<img width="710" height="228" alt="image" src="https://github.com/user-attachments/assets/61daaa9e-a845-429a-93d2-1b79fc7d9da7" />

Неблокирующий запрос на создание индекса
```postgresql
CREATE INDEX CONCURRENTLY idx_doctors_speciality ON doctors (speciality);
```

После создания индекса  

<img width="705" height="228" alt="image" src="https://github.com/user-attachments/assets/967df688-7616-4def-a237-e9c38f358758" />


## Индекс для полнотекстового поиска

Для поиска по ключевым словам в диагнозе
```
CREATE INDEX CONCURRENTLY idx_medical_history_diagnostics on medical_history USING gin (to_tsvector('russian', diagnostics));
```  

## Индекс на поле с функцией

Для поиска пациентов по дате рождения
```
CREATE INDEX CONCURRENTLY idx_patients_age on patients (date_part('year', birthdate::date));
```

## Индекс на несколько полей

Для составления отчета при соответствии клиник и пациентов
```
CREATE INDEX CONCURRENTLY idx_bookings_clinic_patinent on bookings (clinic_id, patinent_id);
```

## Проблемы
Хотел реализовать индекс для поиска по количеству лет, реализовав его через функцию age(), и, как оказалось, нельзя этого сделать из-за того, что функции преобразования даты относятся к изменчивому типу - такие функции в индексе недопустимы.