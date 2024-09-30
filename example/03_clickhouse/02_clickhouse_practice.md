# И ещё немного практики — ClickHouse

Запустите ClickHouse-кластер локально в Docker — скачайте и распакуйте [архив](https://github.com/Yandex-Practicum/python_p2f_ugc/tree/058a2a4712f53b08804f26c2d9e1fd397ed55ad2/8th-sprint-refactor/03_clickhouse/src/clickhouse){target="_blank"}. Внутри — `docker-compose.yml` и папка `data`, в которой находятся конфиги для каждого контейнера из кластера.

В этой конфигурации кластер запускается из двух шардов по две реплики каждый, итого — четыре сервера. Для работы репликаций запускается ZooKeeper. Стандартный порт для ClickHouse — 8123, для ZooKeeper — 9000. В папке `data` есть папки для каждой ноды кластера с соответствующими файлами конфигурации.

Фактически в файлах конфигурации — информация о нодах кластера, и для каждой ноды есть два пользователя: `default` без пароля и `admin` с паролем `123`. Про все доступные конфигурационные параметры можно почитать [в официальной документации](https://clickhouse.tech/docs/ru/operations/server-configuration-parameters/settings/){target="_blank"}.

Для запуска кластера используйте уже привычную команду:

```bash
$ docker-compose -f docker-compose.yml up
```

С ClickHouse можно работать с помощью [клиента командной строки](https://clickhouse.tech/docs/ru/interfaces/cli/){target="_blank"} `clickhouse-client`. Но вам нужно работать с ClickHouse из Python. Для этого откройте [страницу со списком клиентов](https://clickhouse.tech/docs/ru/interfaces/third-party/client-libraries/){target="_blank"} и выберите подходящий вам.

Советуем выбрать библиотеку [clickhouse-driver](https://clickhouse-driver.readthedocs.io){target="_blank"}, потому что она: 

* **Не асинхронная**. Она нужна вам не для разработки, а чтобы сделать запросы с ClickHouse и научиться с ним работать. Асинхронный интерфейс для таких задач неудобен.
* **Не ORM**. Вам нужно научиться делать запросы, а не писать модели. ORM будет только усложнять весь процесс.
* **Имеет документацию**. Не у всех библиотек, представленных на сайте, есть документация.
* **Имеет больше всего звёзд на GitHub**. Это значит, что библиотекой пользуются и она больше всего отвечает потребностям разработчиков.

Установите библиотеку.

```bash
$ pip install clickhouse-driver
```

Далее запустите Python Shell или Jupyter Notebook, в котором вы будете вводить команды для работы с ClickHouse.

```bash
$ python
```

Создайте клиент, через который будут выполняться все запросы:

```python
from clickhouse_driver import Client

client = Client(host='localhost')
```

Если интерпретатор запущен не в Docker, а локально, то доступ к кластеру доступен на localhost.

Итак, первый запрос с ClickHouse. Проверьте, сколько у вас баз данных:

```python
client.execute('SHOW DATABASES')
# [('INFORMATION_SCHEMA',), ('default',), ('information_schema',), ('system',)]
```

Создайте базу данных, в которой будете работать, запросом `CREATE DATABASE`. Если нужны подробности, они есть [в документации](https://clickhouse.tech/docs/ru/sql-reference/statements/create/database/){target="_blank"}.

```sql
client.execute('CREATE DATABASE IF NOT EXISTS example ON CLUSTER company_cluster')
# [('clickhouse-node1', 9000, 0, '', 3, 0), ('clickhouse-node3', 9000, 0, '', 2, 0), ('clickhouse-node4', 9000, 0, '', 1, 0), ('clickhouse-node2', 9000, 0, '', 0, 0)]
```

Если внимательно посмотреть на конфиг любой ноды, то там можно найти название кластера `company_cluster`. Именно его вы использовали в запросе.

Теперь создайте таблицу `regular_table` [запросом `CREATE TABLE`](https://clickhouse.tech/docs/ru/sql-reference/statements/create/table/){target="_blank"}.

```python
client.execute('CREATE TABLE example.regular_table ON CLUSTER company_cluster (id Int64, x Int32) Engine=MergeTree() ORDER BY id')
[('clickhouse-node1', 9000, 0, '', 3, 0), ('clickhouse-node3', 9000, 0, '', 2, 0), ('clickhouse-node4', 9000, 0, '', 1, 0), ('clickhouse-node2', 9000, 0, '', 0, 0)]
```

Благодаря `ON CLUSTER company_cluster` таблица создалась на каждой ноде.

Запишите в таблицу одну строку [запросом `INSERT INTO`](https://clickhouse.tech/docs/ru/sql-reference/statements/insert-into/){target="_blank"}.

```python
client.execute('INSERT INTO example.regular_table (id, x) VALUES (1, 10), (2, 20)')
```

И сразу проверьте, что данные добавились в таблицу, [запросом `SELECT`](https://clickhouse.tech/docs/ru/sql-reference/statements/select/){target="_blank"}.

```python
client.execute('SELECT * FROM example.regular_table')
# [(1, 10), (2, 20)]
```

Вы только что попробовали добавить данные в простую таблицу. Теперь попробуйте создать дистрибутивную реплицированную таблицу. Для этого нужно создать реплицированную шардированную таблицу на нескольких нодах. В вашем случае — две реплики, поэтому таблица будет создана на двух нодах.

Нужно создавать таблицу на каждом шарде отдельным запросом, потому что использовать распределённые DDL-запросы (`ON CLUSTER`) не получится. Для этого нужно подключиться к каждому шарду отдельно. Самый простой способ — зайти в контейнер с базой и использовать [клиент командной строки](https://clickhouse.tech/docs/ru/interfaces/cli/){target="_blank"} `clickhouse-client`.

Зайдите в контейнер первого шарда и запустите клиент:

```bash
$ docker exec -it clickhouse-node1 bash
$ clickhouse-client
```

Настроим первую ноду, которая по совместительству является первым шардом. После каждого запроса должно быть сообщение об успехе `Ok.`.

```sql
CREATE DATABASE shard;
// Ok.
// 0 rows in set. Elapsed: 0.009 sec. 

CREATE TABLE shard.test (id Int64, event_time DateTime) Engine=ReplicatedMergeTree('/clickhouse/tables/shard1/test', 'replica_1') PARTITION BY toYYYYMMDD(event_time) ORDER BY id;
// Ok.
// 0 rows in set. Elapsed: 0.062 sec. 

CREATE TABLE default.test (id Int64, event_time DateTime) ENGINE = Distributed('company_cluster', '', test, rand());
// Ok.
// 0 rows in set. Elapsed: 0.112 sec. 

```

Первым запросом мы создаём базу данных для будущей дистрибутивной таблицы. Второй запрос создаст реплицированную таблицы на шарде. Для этого мы используем специальный [движок `ReplicatedMergeTree`](https://clickhouse.tech/docs/ru/engines/table-engines/mergetree-family/replication/#creating-replicated-tables){target="_blank"}). Он используется для реплицированных таблиц и принимает параметры `zookeeper_path` и `replica_name`. Таблицы с одинаковым `zookeeper_path` будут синхронизироваться.

Последний запрос `CREATE TABLE` создаёт дистрибутивную таблицу, используя [движок `Distributed`](https://clickhouse.tech/docs/ru/engines/table-engines/special/distributed/){target="_blank"}). Первым параметром он принимает название кластера — оно есть в файле конфигурации. Второй параметр — это название базы данных. Если его не указывать, будет использоваться дефолтная база. Дефолтную БД у каждой ноды можно найти в файле конфига, параметр `default_database`. У первой ноды он равен `<default_database>shard</default_database>`. Третий параметр — это локальная таблица, а четвёртый — это ключ шардирования данных. Вы указали созданную таблицу `test` из базы `shard` с рандомным распределением данных по шардам.

Далее подготовим вторую ноду — она будет являться репликой шарда на первой. Выйдите из клиента и контейнера, наберите дважды `exit` или `Ctrl + D`, а после зайдите на второй шард.

```bash
docker exec -it clickhouse-node2 bash
clickhouse-client
```

```sql
CREATE DATABASE replica;
// Ok.
// 0 rows in set. Elapsed: 0.009 sec. 

CREATE TABLE replica.test (id Int64, event_time DateTime) Engine=ReplicatedMergeTree('/clickhouse/tables/shard1/test', 'replica_2') PARTITION BY toYYYYMMDD(event_time) ORDER BY id;
// Ok.
// 0 rows in set. Elapsed: 0.062 sec. 
```

Обратите внимание, что значение параметра `zookeeper_path` у только что созданной реплицированной таблицы и таблицы на первой ноде совпадают. Таким образом, мы закончили настройку первого шарда и его реплиики. По аналогии настроим оставшиеся две ноды.

```bash
$ docker exec -it clickhouse-node3 bash
$ clickhouse-client
```

```sql
CREATE DATABASE shard;
// Ok.
// 0 rows in set. Elapsed: 0.009 sec. 

CREATE TABLE shard.test (id Int64, event_time DateTime) Engine=ReplicatedMergeTree('/clickhouse/tables/shard2/test', 'replica_1') PARTITION BY toYYYYMMDD(event_time) ORDER BY id;
// Ok.
// 0 rows in set. Elapsed: 0.062 sec. 

CREATE TABLE default.test (id Int64, event_time DateTime) ENGINE = Distributed('company_cluster', '', test, rand());
// Ok.
// 0 rows in set. Elapsed: 0.112 sec. 


```

```bash
$ docker exec -it clickhouse-node4 bash
$ clickhouse-client
```

```sql
CREATE DATABASE replica;
// Ok.
// 0 rows in set. Elapsed: 0.009 sec. 

CREATE TABLE replica.test (id Int64, event_time DateTime) Engine=ReplicatedMergeTree('/clickhouse/tables/shard2/test', 'replica_2') PARTITION BY toYYYYMMDD(event_time) ORDER BY id;
// Ok.
// 0 rows in set. Elapsed: 0.062 sec. 
```

 
После выполнения всех запросов по созданию таблицы можно начинать пользоваться дистрибутивной таблицей с помощью Python-клиента...

```python
client.execute('INSERT INTO default.test (id, event_time) VALUES (1, today()), (2, today()), (3, now())')
```

...и выполнить `SELECT`-запрос:

```python
client.execute('SELECT * FROM default.test')
```

## Асинхронный ClickHouse

Для работы с ClickHouse асинхронно можно использовать библиотеку `aiochclient`. Это асинхронный клиент ClickHouse для Python, использующий `asyncio` и базирующийся на `aiohttp`.

```bash
$ pip install aiochclient[aiohttp]
```

Для подключения к ClickHouse `aiochclient` требуется `aiohttp.ClientSession` или `httpx.AsyncClient`:

```python
from aiochclient import ChClient
from aiohttp import ClientSession

async def main():
    async with ClientSession() as s:
        client = ChClient(s)
        assert await client.is_alive()  # возвращает True, если соединение успешно
```

Создание базы:

```python
result = await client.execute('CREATE DATABASE IF NOT EXISTS example ON CLUSTER company_cluster')
[('clickhouse-node1', 9000, 0, '', 3, 0), ('clickhouse-node3', 9000, 0, '', 2, 0), ('clickhouse-node4', 9000, 0, '', 1, 0), ('clickhouse-node2', 9000, 0, '', 0, 0)]
```

Создание таблицы:

```python
result = await client.execute('CREATE TABLE example.regular_table ON CLUSTER company_cluster (id Int64, x Int32) Engine=MergeTree() ORDER BY id')
[('clickhouse-node1', 9000, 0, '', 3, 0), ('clickhouse-node3', 9000, 0, '', 2, 0), ('clickhouse-node4', 9000, 0, '', 1, 0), ('clickhouse-node2', 9000, 0, '', 0, 0)]
```

Вставить данные в таблицу:

```python
await client.execute('INSERT INTO example.regular_table (id, x) VALUES', (1, 10), (2, 20))
```

Выбрать данные из таблицы:

```python
rows = await client.fetch('SELECT * FROM example.regular_table')
[(1, 10), (2, 20)]
```

Этого хватит, чтобы начать работать с ClickHouse и понимать, для чего его можно использовать, а самое главное — как это делать правильно.

### Квиз

1. Как ClickHouse хранит данные в файлах?
   * [x] По колонкам.
 > Правильно! Это позволяет загружать только нужные колонки при обработке запросов.
   * [ ] По строкам.
 > Нет, ClickHouse хранит данные колонками, что позволяет загружать только нужные колонки.
   * [ ] По таблицам. 
 > Нет, ClickHouse хранит данные колонками, что позволяет загружать только нужные колонки.
   * [ ] Можно использовать любой из способов выше.
 >Неверно, ClickHouse хранит данные колонками, что позволяет загружать только нужные колонки.

2. Как будет выполняться запрос `SELECT * FROM db.table WHERE id = 100;`?
   * [ ] ClickHouse загрузит все данные в память, а потом простым перебором найдёт нужную строку.
   > Нет, ClickHouse найдёт нужную гранулу по файлам засечек, загрузит её в память и найдёт в ней нужную строку.
   * [x] ClickHouse найдёт нужную гранулу по файлам засечек, загрузит её в память и найдёт в ней нужную строку.
   > Верно! ClickHouse старается минимизировать количество данных, загружаемых в память. Для этого используются файлы засечек.
   * [ ] ClickHouse загрузит нужные колонки в память, найдёт нужную гранулу по файлам засечек и найдёт в ней нужную строку
   > Нет, ClickHouse найдёт нужную гранулу по файлам засечек, загрузит её в память и найдёт в ней нужную строку.
   * [ ] ClickHouse найдёт нужную гранулу по файлам засечек, найдёт в ней нужную строку и загрузит её в память
   > Нет, ClickHouse найдёт нужную гранулу по файлам засечек, загрузит её в память и найдёт в ней нужную строку.

3. Как устроен индекс в ClickHouse?
   * [ ] Индекс строится так же, как и в PostgreSQL.
   > Индекс строится по блокам данных по каждой 8192-й строке. Это допустимо, потому что данные хранятся упорядоченно. Индекс хранит данные в дополнительных файлах засечек.
   * [ ] Индекс строится по файлам засечек, которые есть у исходных данных
   > Индекс строится по блокам данных по каждой 8192-й строке. Это допустимо, потому что данные хранятся упорядоченно. Индекс хранит данные в дополнительных файлах засечек.
   * [ ] Индекс строится по блокам данных по каждой 8192-й строке. Это допустимо, так как данные хранятся упорядоченно. Индекс хранит данные в RAM.
   > Индекс строится по блокам данных по каждой 8192-й строке. Это допустимо, потому что данные хранятся упорядоченно. Индекс хранит данные в дополнительных файлах засечек.
   * [x] Индекс строится по блокам данных по каждой 8192-й строке. Это допустимо, потому что данные хранятся упорядоченно. Индекс хранит данные в дополнительных файлах засечек
   > Верно! Индекс работает по блокам.

4. Какие основные операции лежат в основе актуализации данных в реплицированных таблицах?
   * [ ] `INSERT`, `UPDATE`, `DELETE`
   > `INSERT` — добавление новых данных, `FETCH` — подгрузка недостающих кусочков у реплики, `MERGE` — объединение блоков `MergeTree`.
   * [x] `INSERT`, `FETCH`, `MERGE`
   > Верно! Добавление новых данных, подгрузка недостающих кусочков у реплики и объединение блоков `MergeTree`.
   * [ ] `APPEND`, `REQUEST`, `REDUCE`
   > `INSERT` — добавление новых данных, `FETCH` — подгрузка недостающих кусочков у реплики, `MERGE` — объединение блоков `MergeTree`.
   * [ ] `INSERT`, `FETCH`, `UNION`
   > `INSERT` — добавление новых данных, `FETCH` — подгрузка недостающих кусочков у реплики, `MERGE` — объединение блоков `MergeTree`.

5. Как хранятся данные в дистрибутивной таблице?
   * [ ] Данные реплицируются на несколько серверов.
   > Дистрибутивная таблица хранит данные в локальных таблицах на дисках нескольких серверов.
   * [ ] Дистрибутивная таблица хранит данные в RAM для более быстрой работы. Чем больше серверов, тем быстрее будет работать.
   > Дистрибутивная таблица хранит данные в локальных таблицах на дисках нескольких серверов.
   * [x] Дистрибутивная таблица хранит данные в локальных таблицах на дисках нескольких серверов.
   > Верно! Данные делятся между серверами, и большинство вычислений производится на серверах на кусочках данных. Это распараллеливание увеличивает время обработки данных.
   * [ ] Дистрибутивная таблица хранит данные в локальных таблицах в RAM нескольких серверов.
   > Дистрибутивная таблица хранит данные в локальных таблицах на дисках нескольких серверов.

## Полезные ссылки

* [Документация ClickHouse](https://clickhouse.tech/docs/ru/){target="_blank"}
* [Статья, описывающая архитектуру дистрибутивных реплицированных таблиц и работу с ними](https://habr.com/ru/post/509540){target="_blank"}
* [Статья Яндекса про запуск ClickHouse в open source](https://habr.com/ru/company/yandex/blog/303282){target="_blank"}
* [Простой и короткий туториал по созданию дистрибутивной реплицированной таблицы на английском](https://medium.com/@merticariug/distributed-clickhouse-configuration-d412c211687c){target="_blank"}
* [Документация по Python-клиенту clickhouse-driver на английском](https://clickhouse-driver.readthedocs.io/en/latest/quickstart.html){target="_blank"}
