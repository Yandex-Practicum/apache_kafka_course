# Агрегация терабайта данных на лету

## История

Вам знаком сервис [Яндекс Метрика](https://metrika.yandex.ru){target="_blank"}? Он собирает аналитику посещения сайтов и позволяет строить произвольные отчёты по этой аналитике.

А слышали про [Google Analytics](https://analytics.google.com){target="_blank"}? Это примерно то же самое, но от Google. Для веб-мастера эти сервисы — главный источник информации о пользователях.

А знаком термин «таргетированная реклама»? Это реклама, которая показывается только аудитории, обладающей общими свойствами. Например, только жителям Москвы. Для создания таргетированной рекламы нужно знать портрет пользователя. В социальных сетях люди сами заполняют информацию о себе. А как быть с онлайн-кинотеатром? Вряд ли пользователи будут заполнять информацию о любимых актёрах. И не факт, что получится выяснить пол и возраст пользователя.

Несмотря на это, простая история путешествия по сайту уже многое говорит о пользователе. Яндекс Метрика установлена на многих российских сайтах. Пока вы путешествуете по интернету, Яндекс Метрика наблюдает за вами от сайта к сайту и постепенно выстраивает ваш портрет. Этот портрет Яндекс использует в таргетированной рекламе ([Яндекс Директ](https://direct.yandex.ru/){target="_blank"}).

Но зачем устанавливать Яндекс Метрику на сайт, если она собирает данные о пользователях в собственных целях? Не совсем так. Веб-мастеру тоже доступна эта аналитика о каждом клике пользователя на сайте. Таким образом, вы тоже можете формировать портрет именно вашего пользователя.

Проверять посещаемость в будни и выходные, с телефонов и компьютеров, проверять рост посещаемости в дни акций, изучать результаты {{A/B-тестирования}}[p2f_A/B].

А теперь представьте, что вы разработали Яндекс Метрику. 80% рунета её установили на своих веб-ресурсах. Сколько запросов с данными о пользователях будет приходить ежеминутно? «Миллиарды данных» о кликах, переходах и других действиях в интернете! Эти данные нужно где-то хранить. Помимо хранения всех событий сёрфа по интернету, нужно дать возможность веб-мастерам формировать отчёты по этим данным для своих ресурсов.

Для полноты картины попробуем оценить объём информации. Допустим, один клик занимает 1 байт и каждую минуту приходит 10 миллиардов новых кликов. Тогда примерный объём будет 1 × 10000000000 / 1024 / 1024 / 1024 ≈ 9,31 Гб. Посчитаем объём данных в сутки: 9,31 × 60 × 24 / 1024 ≈ 13,1 Тб. Внушительное число!

Для хранения и обработки таких данных одним сервером точно не обойтись. Можно использовать распределённые системы, например инструментарий экосистемы Apache Hadoop: Apache HDFS для хранения сырых данных, Apache Spark для расчёта пользовательских отчётов, Apache HBase для хранения уже вычисленных отчётов.

Но у такого подхода есть пара существенных недостатков:

* Пользователю приходится ждать формирования произвольных отчётов. Фактически ему нужно запустить вычисления на Spark-кластере. Если он хочет посмотреть отчёт за месяц, вычисления точно не будут быстрыми и могут занимать по несколько минут. Но если пользователь сможет формировать отчёты за секунды, работа будет намного продуктивней, эффективней и без отрыва от контекста между запросами.  
* Можно не давать строить произвольные отчёты, а предоставить набор шаблонных. По заранее известным отчётам вычисления можно производить сразу, когда приходят данные. Но, во-первых, это потеря гибкости, поэтому веб-мастера не смогут получить полную информацию о пользователе. А во-вторых, это много вычислений, которые могут и не понадобиться. Например, вычислить отчёт за месяц, который веб-мастер ни разу не откроет.

Именно с такими задачами столкнулась команда Яндекс Метрики в 2009 году и попробовала разработать собственное решение. В 2009 году появился прототип собственной СУБД, который умел очень мало. Но прототип постепенно развивался, закрывал всё больше потребностей, и к 2014 году на эту СУБД полностью перебрались все Яндекс Метрики. Называлась СУБД ClickHouse. Объём переводимых данных всех Яндекс Метрик был 2 Пб. Флешкой такие данные не перенести.

ClickHouse быстро обрёл популярность внутри компании, поэтому многие другие продукты начали его использовать в своих целях. В 2016 году ClickHouse выложили в open source под лицензию Apache 2.0. Об этом есть [хабровская статья Яндекса](https://habr.com/ru/company/yandex/blog/303282/){target="_blank"} и [видео доклада](https://www.youtube.com/watch?v=Ho4_dQk7dAg&ab_channel=%D0%A0%D0%B0%D0%B7%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%BA%D0%B0){target="_blank"}. А 2021 году ClickHose вообще отделился от основного Яндекса в отедльную компанию.

Сегодня ClickHouse используют многие компании не только в России, но и за рубежом.

## Как выглядит ClickHouse

Общая задача и область применения ClickHouse — это сохранение каких-то событий и последующий их анализ. Те же клики пользователя по сайту.

ClickHouse — распределённая СУБД колоночного типа с поддержкой SQL-диалекта c расширениями (массивы, кортежи). В отличие от GFS (Google File System), в котором есть два типа серверов — MasterServer и ChunkServer, — в ClickHouse можно всё уместить на одном сервере. Данные в ClickHouse хранятся на диске, что значительно снижает стоимость хранения данных. ClickHouse является колоночной СУБД и данные хранит тоже колонками. При выполнении запроса это позволяет загружать только те колонки, которые нужны. И именно поэтому в ClickHouse лучше хранить широкие таблицы.

ClickHouse хранит блоки данных и работает с ними, а не с каждой строкой по отдельности. Минимальные неделимые блоки данных называются «гранулами».

Как и у всех СУБД, у ClickHouse есть индекс. Во многом он похож на индексы в других СУБД, но есть кардинальное отличие — разреженность. Он индексирует не каждую запись, а каждую десятитысячную. Фактически это значение ключа индекса для каждой 8192-й строки.

Например, есть два индекса: по дате и некоему счётчику (id) события — и запрос.

```sql
SELECT Name, count(*) FROM events WHERE CounterID = 1234 AND Date >= today() - 7;
```

Допустим, сегодня 31 января. В индексе с датой последняя запись — от 15 мая. В индексе со счётчиком есть запись 1235. СУБД прочитает данные с диска от записи 15 мая до записи со счётчиком 1235. С диска будет прочитано чуть больше данных, чем нужно: от 1 до 10 000 записей. Но поиск в массиве из 10 000 записей — тривиальная и быстрая задача. Также обратите внимание, что будут прочитаны только три столбца: `Name`, `CounterID` и `Date`.

Столбцы, в которых хранятся данные, упорядочены по первичному ключу. Файлы, которые содержат значения первичного ключа в индексах и указатели на гранулы, называются «файлы засечек». При обработке запроса первым делом по файлам засечек ищутся гранулы, после чего они считываются с диска и собираются в блок, который отправляется на обработку.

![image](https://pictures.s3.yandex.net/resources/UGC_aggregation_1_1630518480.png)

Файлы засечек и гранулы — основа хранения данных в ClickHouse. Файлы засечек создаются не только для индекса, но и для столбцов, независимо от того, есть индекс или нет.

Порядок в гранулах устанавливается ещё на этапе добавления данных в базу.

Основной движок хранения данных таблиц в ClickHouse называется [MergeTree](https://clickhouse.tech/docs/ru/engines/table-engines/mergetree-family/mergetree/){target="_blank"}. Идея заключается в хранении небольшого количества упорядоченных кусочков, а если их становится много, то MergeTree объединяет несколько кусочков в один. Тем самым алгоритм постоянно поддерживает небольшое количество кусочков данных.

![image](https://pictures.s3.yandex.net/resources/UGC_aggregation_2__1630518484.png)

Чёрные точки — это данные на диске, а синие — это те, что предстоит записать. ClickHouse сначала сортирует новые данные по первичному ключу, а потом записывает их на диск.

![image](https://pictures.s3.yandex.net/resources/UGC_aggregation_3__1630518488.png)

Из-за записи на диск очень важно следить за общим количеством операций записи. Чтобы уменьшить его, лучше вставлять данные пачками. Создание записей по одной сильно снизит эффективность и скорость добавления данных.

Когда кусочков данных, записанных на диск, становится много, MergeTree уменьшает их количество, сливая кусочки. Это происходит в фоновом режиме. Алгоритм сливает только кусочки, созданные друг за другом.

![image](https://pictures.s3.yandex.net/resources/UGC_aggregation_4_1630518494.png)

На примере выше MergeTree слил два кусочка: [M, N] и кусочек, который идёт после N и называется N+1. Получился один новый кусочек [M, N+1].

## OTLP vs OLAP

Теперь мы знаем, что ClickHouse работает с блоками данных, а не конкретными строками. Для таких систем даже есть специальное название — OLAP. Online Analytical Processing — это подход к обработке данных, который поддерживает многомерный анализ данных. OLAP предназначен для обработки долгих транзакций и сложных запросов с большим объёмом данных. Примеры OLAP-систем включают в себя системы, основанные на концепции Data Warehouse, вроде ClickHouse.

Альтернативой OLAP явялется OTLP (Online Transaction Processing) — это подход к обработке данных, который управляет транзакционными приложениями в режиме реального времени. Системы OLTP обычно предназначены для работы с короткими транзакциями, которые требуют быстрого времени отклика. Они в основном поддерживают операции чтения и записи. Примеры OLTP-систем включают в себя реляционные базы данных, такие как PostgreSQL, MySQL.

В целом, разница между OLTP и OLAP сводится к тому, как они используют и обрабатывают данные. OLTP обычно используется для обработки и хранения транзакционных данных, в то время как OLAP используется для анализа данных. По этой причине, ClickHouse, как OLAP-система, оптимизирован для чтения больших объёмов данных и предназначен для работы с блоками данных, что делает его идеально подходящим для выполнения сложных аналитических запросов.


## Распределённое хранение данных

Один из любимых примеров преподавателей Python о модуле `threading`: если вы решили распараллелить программу на Python с помощью потоков, насколько  быстрее она станет выполняться? Правильный ответ: она будет выполняться дольше однопоточной из-за GIL. В случае с ClickHouse всё ровно наоборот: чем больше серверов, тем быстрее будут выполняться запросы. Такой эффект дают дистрибутивные таблицы (англ. distributed table).

Для работы дистрибутивных таблиц в первую очередь нужно {{шардировать данные}}[p2f_shard]. Данные шардированной таблицы на одном шарде называются локальной таблицей (local table).

![image](https://pictures.s3.yandex.net/resources/UGC_aggregation_8_1630518500.png)

Дистрибутивная таблица не хранит сами данные, фактически это view над локальными таблицами. Когда пользователь делает запрос в дистрибутивную таблицу, она формирует и делает запросы в локальные таблицы, объединяет полученные данные и возвращает пользователю.

![image](https://pictures.s3.yandex.net/resources/UGC_aggregation_6_1630518512.png)

Локальные таблицы вычисляют запрос, насколько это возможно, чтобы передавать по сети минимум данных. Тем самым вычисления агрегаций, которые требуют ресурса CPU, выполняются параллельно и изолированно в каждой локальной таблице. Дистрибутивная таблица, отправив все запросы в локальные таблицы, получит от них частично или полностью агрегированный результат, объединит и вернёт пользователю.

Получается, что шардирование данных по серверам — это ещё и распараллеливание агрегаций этих данных. Но за это нужно платить удобством:

1. Дистрибутивную таблицу можно создавать на любом сервере — даже на сервере без данных. Ведь дистрибутивная таблица — это функция, которая сама ходит на серверы с данными. Но обращаться к ней можно только на том сервере, на котором она создана. Поэтому лучше её создавать на каждом сервере: при обращении на любой сервер запрос всегда выполнится по всем данным.
2. Записывать данные в дистрибутивную таблицу можно двумя способами: запросом `INSERT` в дистрибутивную таблицу или самостоятельно раскладывать по локальным таблицам.

Запись в дистрибутивную таблицу через `INSERT` — асинхронная. Таблица кладёт данные во временную папку, опрашивает все шарды, вычисляет с помощью функции ключа шардирования, куда записать данные, и только потом будет записывать. Сам процесс несложный и даже недолгий, но при этом данные будут разложены по шардам без учёта самих данных. Это может сыграть против вас, когда начнётся работа с joins.

Если данные для запроса будут находиться на разных шардах, а именно данные, добавленные с помощью `JOIN`, то вся магия параллельных вычислений пропадёт. Локальные таблицы без необходимых `JOIN`-данных на их шардах вернут прочитанные с диска данные практически без вычислений агрегаций, а сервер с дистрибутивной таблицей примет все данные и будет вычислять их у себя.

Помимо того, что сам по себе процесс вычисления агрегаций с большим набором данных на одном процессоре не самый быстрый, есть и другие нюансы, которые замедляют выполнение запроса:

1. Из-за блочного хранения данных их часто читается немного больше, чем необходимо пользователю. Если у вас 1000+ шардов и несколько сотен терабайт данных, передача по сети каждой лишней гранулы приводит к существенным затратам времени.
2. Из-за невыполненных агрегаций данные не схлопываются и передаются по сети в том объёме, в котором хранятся. Это может оказаться непосильной операцией для сети кластера.

Именно поэтому рекомендуется добавлять данные в дистрибутивную таблицу, самостоятельно раскладывая по локальным таблицам, контролируя и оптимизируя расположение данных для их дальнейшего использования.

## Безопасность данных

Ещё одно важное требование к хранилищам данных — их доступность. В ClickHouse поддерживаются реплицированные таблицы (англ. replicated table), за это отвечает алгоритм [ReplicatedMergeTree](https://clickhouse.tech/docs/ru/engines/table-engines/mergetree-family/replication/){target="_blank"} из семейства [MergeTree](https://clickhouse.tech/docs/ru/engines/table-engines/mergetree-family/mergetree/){target="_blank"}. Идея заключается в выполнении трёх операций и в согласованной работе всех реплик. Согласованность обеспечивает ZooKeeper, а операции достаточно простые:

* `INSERT` — вставить данные в реплику;
* `FETCH` — одна реплика скачивает кусок данных у другой;
* `MERGE` — реплика сливает несколько кусков в один.

Вставка происходит на любую реплику с помощью `INSERT`, и информация о вставке записывается в ZooKeeper. С помощью ZooKeeper все реплики следят за одним набором данных, и когда видят у себя недостающий кусок, то скачивают его с помощью `FETCH`.

Основной алгоритм `MergeTree` сливает кусочки, а значит, их нужно сливать на всех репликах. Чтобы данные не разошлись, сливать кусочки тоже нужно согласованно. Когда приходит такой момент, одна из реплик предупреждает остальные о своих намерениях. Записывает в ZooKeeper информацию о слиянии кусков и выполняет `MERGE`. Благодаря ZooKeeper другие реплики сделают то же самое, и расхождений в данных не будет.

На одном сервере у вас могут быть как реплицированные таблицы, так и простые.

## А вместе можно?

Хотя дистрибутивные и реплицированные таблицы работают независимо, использовать их вместе — лучший вариант. Они прекрасно работают друг с другом. 

Нормальный быстрый и отказоустойчивый кластер ClickHouse — это N шардов по три реплики каждого. Дистрибутивные таблицы умеют понимать, что обращаются в реплицированную, и если она недоступна, то пойдут в другую.

![image](https://pictures.s3.yandex.net/resources/UGC_aggregation_7_1630518517.png)