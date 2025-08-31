> Реализация домашнего задания произведения с помощью Docker

## Физическая репликация

0. Создаем сеть в режиме bridge для общения контейнеров друг с другом по DNS
1. Создаем 2 контейнера из Docker-образов postgres для мастера и реплики в той же сети на разных портах для наглядности
2. Убеждаемся, что инстанс мастера доступен из контейнера реплики

<img width="1108" height="201" alt="image" src="https://github.com/user-attachments/assets/5f59cd79-39e8-4eeb-93a8-bc3c6a41f209" />

3. Создаем пустую директорию в контейнере реплики для "восстановления" из мастера
4. Назначаем ей владельца и группу postgres
5. Устанавливаем права на созданную директорию 0700

<img width="627" height="142" alt="image" src="https://github.com/user-attachments/assets/ae324645-25c0-411e-939b-a75112b86343" />

6. Проваливаемся в контейнер мастера и разрешаем подключаться пользователю postgres для репликации, добавляя правило в pg_hba.conf и перечитывая конфиг

<img width="1405" height="639" alt="image" src="https://github.com/user-attachments/assets/3b9f3b86-c33f-48ba-9605-2de6d3bdae02" />

7. В конрейнере реплики от пользователя postgres запускаем "восстановление" черезе pg_basebackup с инстанса мастера в созданную директорию

<img width="1167" height="238" alt="image" src="https://github.com/user-attachments/assets/aabc95ea-c4ee-4c87-8017-11adea3b8008" />

8. Устанавлиаем порт для реплики в 5433, добавляя строку с параметром в postgresql.conf и запускаем инстанс реплики

<img width="1584" height="276" alt="image" src="https://github.com/user-attachments/assets/c8f00f0e-7498-475a-b790-f42f95c33f99" />

9. Убеждаемся на инстансе реплики, что она "подключена" к мастеру

<img width="1906" height="520" alt="image" src="https://github.com/user-attachments/assets/35ec7f7d-a12b-47c9-a4f6-933dbc676168" />

10. На мастере устанавливаем слот физической репликации

<img width="614" height="219" alt="image" src="https://github.com/user-attachments/assets/49fc59fd-d29f-4073-8792-da5158e8f2e3" />
    
12. На реплике его же помещяем в переменную primary_slot_name, а также устанавливаем задержку репликации через установку параметра recovery_min_apply_delay, и перечитывает конфиг реплики

<img width="910" height="300" alt="image" src="https://github.com/user-attachments/assets/7b48ab93-d9a4-474b-b1ed-85fdedffa599" />
<img width="1899" height="309" alt="image" src="https://github.com/user-attachments/assets/df551681-0c1f-43aa-8bbc-62582c3bdc24" />

13. Произведем изменения на мастере и дождемся изменений на реплике

<img width="690" height="468" alt="image" src="https://github.com/user-attachments/assets/f16c1818-8ffc-492c-bfc1-ef529df2544f" />

Через 5 минут

<img width="694" height="230" alt="image" src="https://github.com/user-attachments/assets/9e25b681-7925-4236-8605-0366580cba86" />


## Логическая репликация

1. В той же сети создаем 2 контейнера - для мастера и для реплики (порт форвард в 5434) с wal-level=logical

<img width="1262" height="90" alt="image" src="https://github.com/user-attachments/assets/615d8567-326b-42a8-ba9b-5357a7474633" />

2. На мастере создаем базу данных test_master, таблицу, инсертим в нее данные и создаем публикацию

<img width="768" height="461" alt="image" src="https://github.com/user-attachments/assets/0d403bf1-6fea-4c1a-a5a6-27b1d58fe6dc" />

3. На реплике создаем базу данных test_replica и создаем аналогичную таблицу без данных

<img width="589" height="438" alt="image" src="https://github.com/user-attachments/assets/279ac0ed-1d02-4f42-ba43-eb05895be0a7" />

4. Создаем подписку и смотрим реплицированные данные

<img width="1351" height="161" alt="image" src="https://github.com/user-attachments/assets/a5cb853c-15e0-489a-9f55-f4c6183964b6" />

5. Для чистоты эксперимента добавим на мастере еще одну строку в эту таблицу и проверим репликацию

<img width="645" height="575" alt="image" src="https://github.com/user-attachments/assets/ede4ff96-edb9-4465-8326-63e90135ffa2" />








