### 1. Иземение данных в нескольких таблицах с помощью хранимой процедуры

Если брать в расчет базу данных medical_center, схему которой я прикладывал на начальных домашних заданиях (все еще не понимаю, является ли она проектом), то возможная бизнес-логика изменения данных одновременно в нескольких таблицах мне видется как заполнение таблиц **`patinets`** и **`bookings`** предполагая, что таблицы **`clinics`** и **`doctors`** уже заполнены.

Поскольку навыком написания процедур на mysql я не владею, написал скрипт на python

```
import random
import string
from datetime import timedelta, date, datetime

def generate_random_date(start_year=1950, end_year=2025):
    start = date(start_year, 1, 1)
    years = end_year - start_year + 1
    max_days = years * 365 + years // 4
    return str(start + timedelta(days=random.randrange(max_days)))

def generate_phone_number():
    first_part = f'7{random.randint(900, 999)}'
    second_part = random.randint(100, 999)
    third_part = random.randint(1000, 9999)
    return f'+{first_part}-{second_part}-{third_part}'

def generate_random_datetime(start_year=2025):
    now = datetime.now()
    start = datetime(start_year, 1, 1)
    delta = now - start
    total_seconds = (delta.days * 24 * 60 * 60) + delta.seconds
    random_second = random.randrange(total_seconds)
    dt = start + timedelta(seconds=random_second)
    rounded_minutes = (dt.minute // 10) * 10
    result = dt.replace(minute=rounded_minutes, second=0, microsecond=0)
    return result

def generate_email(domains=['example.com', 'test.ru']):
    letters = string.ascii_lowercase[:12]
    username_length = random.randint(5, 10)
    domain = random.choice(domains)
    username = ''.join(random.choice(letters) for _ in range(username_length))
    return f'{username}@{domain}'

cursor = conn.cursor()

for patinet in range(1, 11):
    # Генерирую каждого пациента

    name = f"Patinet name #{patinet}"
    surname = f"Patinet surname #{patinet}"
    sex = "Male" if random.random() < 0.6 else "Female"
    birthdate = generate_random_date()
    phone = generate_phone_number()
    email = generate_email()
    sql_query = """
    INSERT INTO categories (name, surname, sex, birthdate, phone, email) VALUES (%s, %s, %s, %s, %s, %s)
    """
    cursor.execute(sql_query, (name, surname, sex, birthdate, phone, email))
    
    # Получаю LAST_INSERT_ID
    last_category_id = cursor.lastrowid

    # Генерирую по 2 записи в таблицу bookings для каждого пациента
    for booking in range(1, 3):
        # Выбираю рандомные clinic_id и doctor_id
        clinic_id = random.randint(1, 10)
        patinent_id = last_category_id
        doctor_id = random.randint(1, 100)
        appointment_time = generate_random_datetime()
        booking_time = generate_random_datetime()
        is_planned = True if random.random() < 0.5 else False
        is_approved = True if random.random() < 0.5 else False
        is_finished = True if random.random() < 0.5 else False

        sql_product_query = """
        INSERT INTO products (clinic_id, patinent_id, doctor_id, appointment_time, is_planned, is_approved, is_finished)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        """
        values = (clinic_id, patinent_id, doctor_id, appointment_time, is_planned, is_approved, is_finished)
        cursor.execute(sql_product_query, values)

conn.commit()
```

### 2. Загрузка данных в базу из csv
 
Судя по заданию, в материалах должен быть csv-файл, но я его там не нашел. Потому буду загружать аналогичный csv, который был описан в листинге занятия

#### Запускаю контейнер mysql в докере, подключаюсь к нему и создаю /var/lib/mysql-files/users.csv

<img width="966" height="219" alt="image" src="https://github.com/user-attachments/assets/f61ac4f0-c4fd-426c-9ba9-debac4560f4a" />


#### Подключаюсь к инстансу, создаю БД, таблицу users и загружаю в нее данные

<img width="1045" height="749" alt="image" src="https://github.com/user-attachments/assets/fc3eb1d9-2605-475f-a1af-95bf9a35834f" />


> По заданибю со * утилиты mysqlimport в рекомендованном образе я не нашел, но если бы она была, импорт из того же csv делал бы так

```
docker exec -it mysql mysqlimport -u root -p123456 otus /var/lib/mysql-files/users.csv
```
