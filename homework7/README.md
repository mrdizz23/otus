## Запрос суммы очков с группировкой и сортировкой по годам

```
SELECT year_game, SUM(points) FROM statistic GROUP BY year_game ORDER BY year_game;
```

<img width="729" height="284" alt="image" src="https://github.com/user-attachments/assets/241dc241-7c1d-4cc7-8d51-00687d2a538e" />


## Написать cte показывающее тоже самое

```
WITH stats as(
  SELECT year_game, SUM(points)
    FROM statistic GROUP BY year_game
)
SELECT * from stats ORDER by year_game;
```  

<img width="681" height="358" alt="image" src="https://github.com/user-attachments/assets/04628958-e998-43b1-a4ac-d19a10a5402e" />


## Используя функцию LAG вывести кол-во очков по всем игрокам за текущий код и за предыдущий

```
WITH stats AS(
  SELECT year_game, SUM(points) AS sum_points_current_year
    FROM statistic GROUP BY year_game
)
SELECT *, LAG(sum_points_current_year)
OVER (ORDER by year_game) as sum_points_pre_year FROM stats;
```

<img width="881" height="382" alt="image" src="https://github.com/user-attachments/assets/aeacf346-a1f9-411d-a44b-a8e1331e86a2" />

