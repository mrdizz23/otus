## Работа с JSON в MySQL

> По следам веббинара с Дмитрием Кирилловым для успешной сдачи ДЗ рекомендовано показать работу с JSON в MySQL

### 1. Создание таблицы и заполнение таблицы с JSON

```
mysql> CREATE TABLE documents (
    ->     id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ->     title VARCHAR(64),
    ->     metadata JSON,
    ->     contents TEXT
    -> );
Query OK, 0 rows affected (0.01 sec)
```

```
mysql> INSERT INTO documents (title, metadata, contents)
    -> VALUES
    ->     ( 'Document 1',
    ->       '{"author": "John",  "tags": ["legal", "real estate"]}',
    ->       'This is a legal document about real estate.' ),
    ->     ( 'Document 2',
    ->       '{"author": "Jane",  "tags": ["finance", "legal"]}',
    ->       'Financial statements should be verified.' ),
    ->     ( 'Document 3',
    ->       '{"author": "Paul",  "tags": ["health", "nutrition"]}',
    ->       'Regular exercise promotes better health.' ),
    ->     ( 'Document 4',
    ->       '{"author": "Alice", "tags": ["travel", "adventure"]}',
    ->       'Mountaineering requires careful preparation.' ),
    ->     ( 'Document 5',
    ->       '{"author": "Bob",   "tags": ["legal", "contracts"]}',
    ->       'Contracts are binding legal documents.' ),
    ->     ( 'Document 6',
    ->        '{"author": "Eve",  "tags": ["legal", "family law"]}',
    ->        'Family law addresses diverse issues.' ),
    ->     ( 'Document 7',
    ->       '{"author": "John",  "tags": ["technology", "innovation"]}',
    ->       'Tech innovations are changing the world.' );
Query OK, 7 rows affected (0.01 sec)
Records: 7  Duplicates: 0  Warnings: 0
```

```
mysql> SELECT * FROM documents;
+----+------------+----------------------------------------------------------+----------------------------------------------+
| id | title      | metadata                                                 | contents                                     |
+----+------------+----------------------------------------------------------+----------------------------------------------+
|  1 | Document 1 | {"tags": ["legal", "real estate"], "author": "John"}     | This is a legal document about real estate.  |
|  2 | Document 2 | {"tags": ["finance", "legal"], "author": "Jane"}         | Financial statements should be verified.     |
|  3 | Document 3 | {"tags": ["health", "nutrition"], "author": "Paul"}      | Regular exercise promotes better health.     |
|  4 | Document 4 | {"tags": ["travel", "adventure"], "author": "Alice"}     | Mountaineering requires careful preparation. |
|  5 | Document 5 | {"tags": ["legal", "contracts"], "author": "Bob"}        | Contracts are binding legal documents.       |
|  6 | Document 6 | {"tags": ["legal", "family law"], "author": "Eve"}       | Family law addresses diverse issues.         |
|  7 | Document 7 | {"tags": ["technology", "innovation"], "author": "John"} | Tech innovations are changing the world.     |
+----+------------+----------------------------------------------------------+----------------------------------------------+
7 rows in set (0.01 sec)
```

### 2. Получение значение объектов JSON из таблицы

```
mysql> SELECT * FROM documents
    ->     WHERE metadata->>'$.author' = 'John';
+----+------------+----------------------------------------------------------+---------------------------------------------+
| id | title      | metadata                                                 | contents                                    |
+----+------------+----------------------------------------------------------+---------------------------------------------+
|  1 | Document 1 | {"tags": ["legal", "real estate"], "author": "John"}     | This is a legal document about real estate. |
|  7 | Document 7 | {"tags": ["technology", "innovation"], "author": "John"} | Tech innovations are changing the world.    |
+----+------------+----------------------------------------------------------+---------------------------------------------+
2 rows in set (0.00 sec)
```

### 3. Изменение объектов JSON

```
mysql> UPDATE documents
    -> SET metadata = JSON_SET(metadata, '$.author', 'Bob')
    -> WHERE id = 4;
Query OK, 1 row affected (0.01 sec)
Rows matched: 1  Changed: 1  Warnings: 0
```

```
mysql> SELECT * FROM documents WHERE id = 4;
+----+------------+----------------------------------------------------+----------------------------------------------+
| id | title      | metadata                                           | contents                                     |
+----+------------+----------------------------------------------------+----------------------------------------------+
|  4 | Document 4 | {"tags": ["travel", "adventure"], "author": "Bob"} | Mountaineering requires careful preparation. |
+----+------------+----------------------------------------------------+----------------------------------------------+
1 row in set (0.00 sec)
```
