## Задание в сфере тестирования

### 1. Корректировка типов полей

> Отметил комментариями исправленные поля

```
CREATE TABLE categories IF NOT EXISTS (
    category_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(32) UNIQUE NOT NULL -- хочется уникальное поле, а то смысл
);

CREATE TABLE products IF NOT EXISTS (
    product_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(32) NOT NULL,
    category_id VARCHAR(32) REFERENCES categories (category_id),
    price INT UNSIGNED NOT NULL, -- цена не может быть отрицаиельной
    rating INT,
    status VARCHAR(32) NOT NULL,
    CONSTRAINT un__line_items__category_id__title
        UNIQUE (category_id, title) -- уникальная связка категории и товара
);
```



### 2. Скрипт генерации (python)

```
cursor = conn.cursor()

for cat_id in range(1, 21):
    # Генерирую каждую категорию
    category_title = f"Категория {cat_id}"
    sql_query = """
    INSERT INTO categories (title) VALUES (%s)
    """
    cursor.execute(sql_query, (category_title,))
    
    # Получаю LAST_INSERT_ID
    last_category_id = cursor.lastrowid

    # Генерирую цены для товаров внутри каждой категории исходя из требования уникальности
    start_price = (cat_id - 1) * 10000 + 1
    end_price = start_price + 9999
    prices = list(range(start_price, end_price + 1))
    
    # Генерирую товары
    for prod_id in range(1, 10001):
        title = f"Товар {prod_id} категории {cat_id}"
        price = prices.pop(0)
        rating = random.randint(1, 5)
        status = "В наличии" if random.random() < 0.8 else "Распродан"

        sql_product_query = """
        INSERT INTO products (title, category_id, price, rating, status)
        VALUES (%s, %s, %s, %s, %s)
        """
        values = (title, last_category_id, price, rating, status)
        cursor.execute(sql_product_query, values)

conn.commit()
```

### 3. Запрос на получение продуктов с сортировками

> Примечание  

Запрос выдаст только первую страницу - чтобы получить любую другую, необходимо в OFFSET передать номер страницы, умноженный на 50.  
Будет работать при условии только статусов "В наличии" и "Распродан" - при сортировке по status как раз и выстраивается в нужном порядке. В противном случае пишем запрос на конкретные статусы

```
SELECT *
FROM products
ORDER BY status, price
LIMIT 50 OFFSET 50;
```
