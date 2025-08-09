## Индекс по полю speciality таблицы doctors

До создания индекса

<img width="710" height="228" alt="image" src="https://github.com/user-attachments/assets/61daaa9e-a845-429a-93d2-1b79fc7d9da7" />


```postgresql
CREATE INDEX CONCURRENTLY idx_doctors_speciality ON doctors (speciality);
```
После создания индекса

<img width="705" height="228" alt="image" src="https://github.com/user-attachments/assets/967df688-7616-4def-a237-e9c38f358758" />


## Индекс для полнотекстового поиска

```

```  

## Индекс на поле с функцией

```

```

## Индекс на несколько полей

