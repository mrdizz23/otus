1. Запускаю контейнер mongodb в Docker и подключаюсь к нему

```
(venv) [dizz@MUR-PC-3009-B2C ~]$ docker run --rm --name mongo -d -p 27017:27017 mongo
Unable to find image 'mongo:latest' locally
latest: Pulling from library/mongo
dcec2d403c4e: Pull complete
88373bcb58c1: Pull complete
1732f2a4d259: Pull complete
a89cc6fb1ead: Pull complete
6dc84fd2f3ac: Pull complete
20043066d3d5: Pull complete
7230971360e4: Pull complete
8a63055b2837: Pull complete
Digest: sha256:7245ffb851d149dbfac67397caf91bae4974d899972f9fd1d8985fc6eea3c13d
Status: Downloaded newer image for mongo:latest
6895524b7d8023f53bedecfefae24e3a2fe9cb8481b52d9d1e6ce7b83bc717f8
[dizz@MUR-PC-3009-B2C ~]$ docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED         STATUS         PORTS                                             NAMES
6895524b7d80   mongo     "docker-entrypoint.s…"   7 seconds ago   Up 4 seconds   0.0.0.0:27017->27017/tcp, [::]:27017->27017/tcp   mongo

(venv) [dizz@MUR-PC-3009-B2C ~]$ docker exec -it mongo mongosh --port 27017
Current Mongosh Log ID: 692eac869e90d78fff9dc29c
Connecting to:          mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.9
Using MongoDB:          8.2.2
Using Mongosh:          2.5.9

For mongosh info see: https://www.mongodb.com/docs/mongodb-shell/

------
   The server generated these startup warnings when booting
   2025-12-02T09:00:51.471+00:00: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine. See http://dochub.mongodb.org/core/prodnotes-filesystem
   2025-12-02T09:00:51.973+00:00: Access control is not enabled for the database. Read and write access to data and configuration is unrestricted
   2025-12-02T09:00:51.974+00:00: For customers running the current memory allocator, we suggest changing the contents of the following sysfsFile
   2025-12-02T09:00:51.974+00:00: For customers running the current memory allocator, we suggest changing the contents of the following sysfsFile
   2025-12-02T09:00:51.974+00:00: We suggest setting the contents of sysfsFile to 0.
   2025-12-02T09:00:51.974+00:00: vm.max_map_count is too low
   2025-12-02T09:00:51.974+00:00: We suggest setting swappiness to 0 or 1, as swapping can cause performance problems.
------

test> show databases
admin   40.00 KiB
config  12.00 KiB
local   40.00 KiB
test>
```

2. С помощью скрипта на python генерирую новую БД, коллекцию и генерирую 100 документов

```
(venv) [dizz@MUR-PC-3009-B2C ~]$ cat generate_data.py
from pymongo import MongoClient
from faker import Faker
import random

client = MongoClient('mongodb://localhost:27017/')
db = client['mydatabase']
collection = db['users']

fake = Faker()
Faker.seed(random.randint(0, 100))

for _ in range(90):
    user_data = {
        'name': fake.name(),
        'email': fake.email(),
        'age': random.randint(18, 65),
        'address': fake.address().replace("\\n", ", ")
    }
    collection.insert_one(user_data)
```

3. Подключаюсь к инстансу и проверяю результаты генерации

```
(venv) [dizz@MUR-PC-3009-B2C ~]$ docker exec -it mongo mongosh --port 27017 mydatabase
Current Mongosh Log ID: 692eae4c0c4f8733949dc29c
Connecting to:          mongodb://127.0.0.1:27017/mydatabase?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.9
Using MongoDB:          8.2.2
Using Mongosh:          2.5.9

For mongosh info see: https://www.mongodb.com/docs/mongodb-shell/

------
   The server generated these startup warnings when booting
   2025-12-02T09:00:51.471+00:00: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine. See http://dochub.mongodb.org/core/prodnotes-filesystem
   2025-12-02T09:00:51.973+00:00: Access control is not enabled for the database. Read and write access to data and configuration is unrestricted
   2025-12-02T09:00:51.974+00:00: For customers running the current memory allocator, we suggest changing the contents of the following sysfsFile
   2025-12-02T09:00:51.974+00:00: For customers running the current memory allocator, we suggest changing the contents of the following sysfsFile
   2025-12-02T09:00:51.974+00:00: We suggest setting the contents of sysfsFile to 0.
   2025-12-02T09:00:51.974+00:00: vm.max_map_count is too low
   2025-12-02T09:00:51.974+00:00: We suggest setting swappiness to 0 or 1, as swapping can cause performance problems.
------

mydatabase> show collections
users
mydatabase> db.users.countDocuments({})
100
```

4. Пара простых запросов на получение данных

```
(venv) [dizz@MUR-PC-3009-B2C ~]$ docker exec -it mongo mongosh --port 27017 mydatabase
Current Mongosh Log ID: 692eaeb448277b66c19dc29c
Connecting to:          mongodb://127.0.0.1:27017/mydatabase?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.5.9
Using MongoDB:          8.2.2
Using Mongosh:          2.5.9

For mongosh info see: https://www.mongodb.com/docs/mongodb-shell/

------
   The server generated these startup warnings when booting
   2025-12-02T09:00:51.471+00:00: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine. See http://dochub.mongodb.org/core/prodnotes-filesystem
   2025-12-02T09:00:51.973+00:00: Access control is not enabled for the database. Read and write access to data and configuration is unrestricted
   2025-12-02T09:00:51.974+00:00: For customers running the current memory allocator, we suggest changing the contents of the following sysfsFile
   2025-12-02T09:00:51.974+00:00: For customers running the current memory allocator, we suggest changing the contents of the following sysfsFile
   2025-12-02T09:00:51.974+00:00: We suggest setting the contents of sysfsFile to 0.
   2025-12-02T09:00:51.974+00:00: vm.max_map_count is too low
   2025-12-02T09:00:51.974+00:00: We suggest setting swappiness to 0 or 1, as swapping can cause performance problems.
------

mydatabase> db.users.find().limit(1)
[
  {
    _id: ObjectId('692eaca1a6d357a92743c547'),
    name: 'Jason Rodriguez',
    email: 'frazierkatherine@example.com',
    age: 64,
    address: '883 Berry Locks\nEast Jessicaville, VA 49792'
  }
]
mydatabase> db.users.find({email: {$regex: "^f"}})
[
  {
    _id: ObjectId('692eaca1a6d357a92743c547'),
    name: 'Jason Rodriguez',
    email: 'frazierkatherine@example.com',
    age: 64,
    address: '883 Berry Locks\nEast Jessicaville, VA 49792'
  },
  {
    _id: ObjectId('692eaca1a6d357a92743c54c'),
    name: 'Alejandra Johnson',
    email: 'freemanjordan@example.net',
    age: 48,
    address: '841 Robinson Walks Apt. 792\nRichardview, IN 22323'
  },
  {
    _id: ObjectId('692eae482b33900fb3073eba'),
    name: 'Sarah Poole',
    email: 'francisberry@example.org',
    age: 40,
    address: '0801 Faith Greens\nMarshallchester, MO 62128'
  },
  {
    _id: ObjectId('692eae482b33900fb3073eca'),
    name: 'Joshua Patton',
    email: 'fcochran@example.com',
    age: 61,
    address: '675 Taylor Drives Suite 003\nWest Christopherfort, NM 71216'
  },
  {
    _id: ObjectId('692eae482b33900fb3073efc'),
    name: 'Charlotte Watson',
    email: 'florestyler@example.com',
    age: 41,
    address: '87162 Price Hollow Apt. 571\nWest Josephville, AK 16632'
  },
  {
    _id: ObjectId('692eae482b33900fb3073f02'),
    name: 'Veronica Rojas',
    email: 'florestoni@example.org',
    age: 24,
    address: '00714 Richard Wells\nPughview, NM 91354'
  }
]
```

5. Обновление данных

```
mydatabase> db.users.find({name: "Sarah Poole"})
[
  {
    _id: ObjectId('692eae482b33900fb3073eba'),
    name: 'Sarah Poole',
    email: 'francisberry@example.org',
    age: 40,
    address: '0801 Faith Greens\nMarshallchester, MO 62128'
  }
]
mydatabase> db.users.update({name: "Sarah Poole"}, {$set: {age:404}})
DeprecationWarning: Collection.update() is deprecated. Use updateOne, updateMany, or bulkWrite.
{
  acknowledged: true,
  insertedId: null,
  matchedCount: 1,
  modifiedCount: 1,
  upsertedCount: 0
}
mydatabase> db.users.find({name: "Sarah Poole"})
[
  {
    _id: ObjectId('692eae482b33900fb3073eba'),
    name: 'Sarah Poole',
    email: 'francisberry@example.org',
    age: 404,
    address: '0801 Faith Greens\nMarshallchester, MO 62128'
  }
]
```

6. Создание индекса и проверка его работы

```
mydatabase> db.users.find({name: "Sarah Poole"}).explain()
{
  explainVersion: '1',
  queryPlanner: {
    namespace: 'mydatabase.users',
    parsedQuery: { name: { '$eq': 'Sarah Poole' } },
    indexFilterSet: false,
    queryHash: 'F4DDDCDC',
    planCacheShapeHash: 'F4DDDCDC',
    planCacheKey: 'E45FBFA1',
    optimizationTimeMillis: 0,
    maxIndexedOrSolutionsReached: false,
    maxIndexedAndSolutionsReached: false,
    maxScansToExplodeReached: false,
    prunedSimilarIndexes: false,
    winningPlan: {
      isCached: false,
      stage: 'COLLSCAN',
      filter: { name: { '$eq': 'Sarah Poole' } },
      direction: 'forward'
    },
    rejectedPlans: []
  },
  queryShapeHash: 'F969745E43D4E1C10C6940F2A8A0D59738840C4CF7251E6238A0A5270A3C8609',
  command: {
    find: 'users',
    filter: { name: 'Sarah Poole' },
    '$db': 'mydatabase'
  },
  serverInfo: {
    host: '6895524b7d80',
    port: 27017,
    version: '8.2.2',
    gitVersion: '594f839ceec1f4385be9a690131412d67b249a0a'
  },
  serverParameters: {
    internalQueryFacetBufferSizeBytes: 104857600,
    internalQueryFacetMaxOutputDocSizeBytes: 104857600,
    internalLookupStageIntermediateDocumentMaxSizeBytes: 104857600,
    internalDocumentSourceGroupMaxMemoryBytes: 104857600,
    internalQueryMaxBlockingSortMemoryUsageBytes: 104857600,
    internalQueryProhibitBlockingMergeOnMongoS: 0,
    internalQueryMaxAddToSetBytes: 104857600,
    internalDocumentSourceSetWindowFieldsMaxMemoryBytes: 104857600,
    internalQueryFrameworkControl: 'trySbeRestricted',
    internalQueryPlannerIgnoreIndexWithCollationForRegex: 1
  },
  ok: 1
}
mydatabase> db.users.createIndex({name: 1})
name_1
mydatabase> db.users.find({name: "Sarah Poole"}).explain()
{
  explainVersion: '1',
  queryPlanner: {
    namespace: 'mydatabase.users',
    parsedQuery: { name: { '$eq': 'Sarah Poole' } },
    indexFilterSet: false,
    queryHash: 'F4DDDCDC',
    planCacheShapeHash: 'F4DDDCDC',
    planCacheKey: '34C29116',
    optimizationTimeMillis: 0,
    maxIndexedOrSolutionsReached: false,
    maxIndexedAndSolutionsReached: false,
    maxScansToExplodeReached: false,
    prunedSimilarIndexes: false,
    winningPlan: {
      isCached: false,
      stage: 'FETCH',
      inputStage: {
        stage: 'IXSCAN',
        keyPattern: { name: 1 },
        indexName: 'name_1',
        isMultiKey: false,
        multiKeyPaths: { name: [] },
        isUnique: false,
        isSparse: false,
        isPartial: false,
        indexVersion: 2,
        direction: 'forward',
        indexBounds: { name: [ '["Sarah Poole", "Sarah Poole"]' ] }
      }
    },
    rejectedPlans: []
  },
  queryShapeHash: 'F969745E43D4E1C10C6940F2A8A0D59738840C4CF7251E6238A0A5270A3C8609',
  command: {
    find: 'users',
    filter: { name: 'Sarah Poole' },
    '$db': 'mydatabase'
  },
  serverInfo: {
    host: '6895524b7d80',
    port: 27017,
    version: '8.2.2',
    gitVersion: '594f839ceec1f4385be9a690131412d67b249a0a'
  },
  serverParameters: {
    internalQueryFacetBufferSizeBytes: 104857600,
    internalQueryFacetMaxOutputDocSizeBytes: 104857600,
    internalLookupStageIntermediateDocumentMaxSizeBytes: 104857600,
    internalDocumentSourceGroupMaxMemoryBytes: 104857600,
    internalQueryMaxBlockingSortMemoryUsageBytes: 104857600,
    internalQueryProhibitBlockingMergeOnMongoS: 0,
    internalQueryMaxAddToSetBytes: 104857600,
    internalDocumentSourceSetWindowFieldsMaxMemoryBytes: 104857600,
    internalQueryFrameworkControl: 'trySbeRestricted',
    internalQueryPlannerIgnoreIndexWithCollationForRegex: 1
  },
  ok: 1
}
```

<img width="1332" height="824" alt="image" src="https://github.com/user-attachments/assets/3a3bb408-31b1-4e30-a764-c807ac31a9eb" />
