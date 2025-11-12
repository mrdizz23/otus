### Дополнительные индексы.

Анализ "проекта" и добавление необходимых индексов уже были проведены в ДЗ №2, потому дублирую информацию

```CREATE INDEX idx_patients_full_name ON patients (surname, name);```  
Составной индекс по Фамилии и имени пацинта (именно в таком порядке с уменьшением кардинальности) для быстрого поиска  

```CREATE INDEX idx_patients_phone ON patients (phone);```  
Индекс для быстрого поиска пациента по номеру телефона

```CREATE INDEX idx_doctors_full_name ON doctors (surname, name) where (is_active);```  
Составной индекс для быстрого поиска активного (находящегося не в отпуске, и не уволенного) доктора

```CREATE INDEX idx_calls_phone ON calls (phone);```
Индекс для получения отчета истории звонков по номеру телефона

### Полнотекстовый индекс

Для поиска по ключевым словам в поле диагноз таблицы **'medical_history'** (имитирую, что поле заполнено на английском языке)

```
mysql> select * from medical_history where match (diagnostics) against ('rupture');
ERROR 1191 (HY000): Can't find FULLTEXT index matching the column list

mysql> create fulltext index fulltext_idx_medical_history_diagnostics on medical_history (diagnostics);
Query OK, 0 rows affected, 1 warning (0.53 sec)
Records: 0  Duplicates: 0  Warnings: 1

mysql> select * from medical_history where match (diagnostics) against ('rupture');
Empty set (0.00 sec)

mysql> explain select * from medical_history where match (diagnostics) against ('rupture');
+----+-------------+-----------------+------------+----------+------------------------------------------+------------------------------------------+---------+-------+------+----------+-------------------------------+
| id | select_type | table           | partitions | type     | possible_keys                            | key                                      | key_len | ref   | rows | filtered | Extra                         |
+----+-------------+-----------------+------------+----------+------------------------------------------+------------------------------------------+---------+-------+------+----------+-------------------------------+
|  1 | SIMPLE      | medical_history | NULL       | fulltext | fulltext_idx_medical_history_diagnostics | fulltext_idx_medical_history_diagnostics | 0       | const |    1 |   100.00 | Using where; Ft_hints: sorted |
+----+-------------+-----------------+------------+----------+------------------------------------------+------------------------------------------+---------+-------+------+----------+-------------------------------+
1 row in set, 1 warning (0.00 sec)

mysql>
```
