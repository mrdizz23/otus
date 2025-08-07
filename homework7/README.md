## Запрос суммы очков с группировкой и сортировкой по годам

```
SELECT year_game, SUM(points) FROM statistic GROUP BY year_game ORDER BY year_game;
```
## Написать cte показывающее тоже самое

```
WITH stats as(
  SELECT year_game, SUM(points)
    FROM statistic GROUP BY year_game
)
SELECT * from stats ORDER by year_game;
```  

## Используя функцию LAG вывести кол-во очков по всем игрокам за текущий код и за предыдущий

```
WITH stats AS(
  SELECT year_game, SUM(points) AS sum_points_current_year
    FROM statistic GROUP BY year_game
)
SELECT *, LAG(sum_points_current_year)
OVER (ORDER by year_game) as sum_points_pre_year FROM stats;
```

