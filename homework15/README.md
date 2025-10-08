### 1. Скрипт генерации магазинов и продаж (python)

```
import random
from datetime import timedelta, date, datetime

def generate_random_datetime(start_year=2023):
    now = datetime.now()
    start = datetime(start_year, 1, 1)
    delta = now - start
    total_seconds = (delta.days * 24 * 60 * 60) + delta.seconds
    random_second = random.randrange(total_seconds)
    return start + timedelta(seconds=random_second)

cursor = conn.cursor()

# Генерация 10 магазинов
for store in range(1, 11):

    address = f"Store #{store}"
    sql_query = """
    INSERT INTO stores (address) VALUES (%s)
    """
    cursor.execute(sql_query, (address))

# Генерация 100000 продаж 
for sale in range(1, 100000 + 1):

    date = generate_random_datetime()
    store_id = 1 if random.random() < 0.75 else random.randint(2, 10)
    sale_amount = round(random.uniform(1,10000),2)
    sql_product_query = """
    INSERT INTO sales (date, store_id, sale_amount)
    VALUES (%s, %s, %s)
    """
    values = (date, store_id, sale_amount)
    cursor.execute(sql_product_query, values)

conn.commit()
```


### 2. Запрос, который выведет нарастающий итог продаж по каждому магазину с группировкой по месяцам
 
```
SELECT
  st.address,
  sa.sale_amount,
  DATE_FORMAT(sa.date, '%Y-%M') AS year_month,
  SUM(sa.sale_amount) OVER(PARTITION BY DATE_FORMAT(sa.date, '%Y-%M') ORDER BY sa.sale_amount, sa.sale_id) AS total_price
FROM sales sa JOIN stores st USING (store_id);
```

### 3. Запрос, который выведет 7-дневное скользящее среднее за последний месяц по самому плодовитому магазину

```
WITH effective_store AS(
  SELECT store_id, SUM(sale_amount) AS sum_amount
  FROM sales GROUP BY store_id ORDER BY sum_amount DESC LIMIT 1
)
SELECT sa.* FROM sales sa
JOIN effective_store es ON sa.store_id = es.store_id
WHERE sa.date BETWEEN (CURRENT_DATE() - INTERVAL 1 MONTH) AND CURRENT_DATE()
ORDER BY sa.date;
```

### 4. Граничные случаи

По нарастающему итогу учет возможную одинаковую цену, включив первичный ключ в сортировкку оконной функции