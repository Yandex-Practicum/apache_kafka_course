# Что такое Kafka

Kafka — это распределённый брокер сообщений для стриминга событий с открытым исходным кодом.

Стриминг событий — это цифровой эквивалент центральной нервной системы человека. Говоря техническим языком, это практика снятия данных в режиме реального времени с различных источников — баз данных, датчиков, мобильных устройств, облачных сервисов и программных приложений — в форме потока событий и сохранение этого потока для дальнейшей обработки, изменения или реакции в режиме реального времени или отложенно.

Событие обозначает факт, что в информационной системе что-то произошло — в бизнес-терминах приложения. Когда вы пишете в Kafka или читаете из неё, вы делаете это в форме событий. Концептуально каждое событие обычно имеет имя или ключ, значение, дату и время, когда оно произошло, а также опциональные заголовки метаданных. Каждое событие иммутабельно (его нельзя редактировать). Пример события о прогрессе просмотра фильма пользователем может выглядеть так:

- Ключ события: `Иван Иванов`.
- Значение события: `Просмотрел Иронию судьбы с 00:53:24 до 00:53:25 (1 секунду)`.
- Дата события: `14 января 2021 года в 09:37`.

В контексте Kafka, эти «события» обычно называются сообщениями. Сообщение в Kafka — это просто массив байтов, поэтому вы можете сохранить в нём любые данные в любом формате. Однако на практике сообщения часто представляют собой данные в форматах типа JSON, Avro и Protobuf. Каждое сообщение также имеет ключ, который представляет собой массив байтов и может быть использован для определения партиции, в которую будет отправлено сообщение.

Стриминг событий обеспечивает непрерывный поток и интерпретацию данных так, что нужная информация оказывается в нужном месте и в нужное время.

Kafka объединяет три ключевые возможности, которые позволяют использовать её для реализации стриминга событий онлайн-кинотеатра:

- Публиковать или добавлять новые события в поток и подписываться на события из потока или считывать их.
- Сохранять потоки событий надёжно и долговечно, как это делают РСУБД.
- Реагировать на события, как только они происходят либо позже (ретроспективно).

Писатели (Producers) — это клиентские приложения, которые публикуют или добавляют события в Kafka, а читатели (Consumers) — это приложения, которые подписываются на интересующие их события, читают и реагируют на них. Важно понимать, что читатели сами решают, когда и сколько событий им нужно нужно получить (Pull-модель).

Kafka обеспечивает различные гарантии доставки сообщений.  At most once и At least once полностью поддерживаются. С Exactly once сложнее. Это более сложная гарантия доставки, и Kafka обеспечивает её в основном в контексте обмена сообщениями внутри кластера. Во внешнем взаимодействии с клиентскими приложениями выполнение этой гарантии в значительной степени зависит от специфики реализации клиентского приложения, а также от обработки ошибок и сбоев.

## Топики

Топики (англ. topics) в Kafka представляют собой логические контейнеры, в которых организуются и надёжно хранятся события. Если провести аналогию, то топик можно представить как папку в файловой системе. Топики содержат партиции, которые представляют собой append-only-файлы. В этих партициях (файлах) хранятся события, а каждое событие записывается в конец файла, отсюда и “append-only”.

Топики в Kafka могут иметь один или несколько писателей или не иметь их вообще. То же самое справедливо и для читателей. Таким образом, читатели и писатели полностью независимы друг от друга — это один из ключевых элементов дизайна Kafka для широких возможностей масштабирования.

События из топиков могут быть прочитаны бесконечное количество раз. В отличие от традиционных систем доставки сообщений вроде [RabbitMQ](https://www.rabbitmq.com/){target="_blank"}, события не удаляются после прочтения. Вместо этого вы можете настроить, как долго сообщения будут храниться в каждом топике. Производительность Kafka относительно стабильна в зависимости от размера данных, поэтому данные можно накапливать столько, сколько нужно.

## Партиции

Топики в Kafka делятся на **партиции**. Это означает, что данные в топиках распределены между набором «корзин», которые размещены на разных брокерах. Такое распределение данных очень важно для масштабирования, поскольку это позволяет клиентским приложениям читать и записывать данные параллельно в разные брокеры. Когда новое событие публикуется в топик, оно добавляется в одну из его партиций.

Механизм выбора партиции зависит от писателя и настроек. Если при записи сообщения указан ключ, то Kafka использует его для определения, в какую партицию пойдет сообщение. В этом случае, все события с одним и тем же ключом, например, ID пользователя, всегда будут добавлены в одну и ту же партицию. Если ключ не указан, то сообщения распределяются по всем партициям топика равномерно.

Kafka гарантирует, что любой потребитель для конкретного топика и партиции будет всегда считывать события этой партиции в том же порядке, в котором они были туда записаны.

Когда новое событие публикуется в топик, оно добавляется в одну из его партиций. Kafka выбирает подходящую партицию с помощью ключа партиции. Ключ партиции (англ. partition key) — это опциональный компонент сообщения, который используется для определения того, в какую партицию будет отправлено сообщение. Если ключ партиции указан, Kafka гарантирует, что все сообщения с одним и тем же ключом будут направлены в одну и ту же партицию. За это отвечает алгоритм хэширования Murmur2. Его работу можно писать так:

1. Kafka получает ключ сообщения, который может быть любого типа. Для облегчения процесса хеширования ключ обычно сериализуется в байтовый массив.
2. Kafka использует алгоритм MurmurHash2 для преобразования сериализованного ключа в 32-битное целое число, известное как хеш-код.
3. Хеш-код делится по модулю на количество партиций в топике. Результат этого деления определяет номер партиции, в которую будет отправлено сообщение.

Например, предположим, что у нас есть топик с тремя партициями и ключ сообщения `example`. Kafka преобразует `example` в байтовый массив, затем использует MurmurHash2 для получения хеш-кода. Предположим, что полученный хеш-код равен 250. Тогда 250 делится по модулю на три (количество партиций), и остаток от деления (единица) становится номером партиции, в которую будет записано сообщение.

Если ключ сообщения равен `null`, Kafka будет распределять сообщения по партициям в режиме round-robin. Однако, следует учесть, что распределение round-robin в Kafka не работает на уровне отдельных сообщений. Вместо этого он сначала отправляет пачку сообщений в одну партицию, затем переключается и отправляет следующую пачку сообщений в другую партицию, и так далее. Это означает, что сообщения будут распределены по всем доступным партициям не строго поочередно, а пакетами.

Использование ключа партиции обеспечивает сохранение порядка сообщений с одним и тем же ключом. Kafka гарантирует, что любой потребитель для конкретного топика и партиции будет всегда считывать события этой партиции в том же порядке, в котором они были туда записаны.

## Кластер Kafka

Kafka работает на уровне кластера. Кластер Kafka состоит из одного или нескольких серверов, называемых брокерами, которые могут хранить, обрабатывать и передавать данные. Каждый брокер хранит определённые топики и партиции, распределённые между всеми брокерами кластера. Это распределение данных помогает достигнуть отказоустойчивости и масштабируемости, так как даже при отказе одного брокера, другие брокеры продолжают работать и обрабатывать данные.

Какое-то время назад  Kafka активно использовала Apache Zookeeper для своих внутренних нужд. Zookeeper — это сервис, который используется для обеспечения согласованности в распределённых системах. Kafka использовала его для управления и синхронизации состояния брокеров и топиков в кластере.

С помощью Zookeeper контроллер Kafka мог управлять многими аспектами работы кластера: отслеживать состояние брокеров, координировать выборы лидеров для партиций, следить за конфигурацией топиков и многое другое

И возникает разумный вопрос — зачем использовать дополнительное хранилище метаданных, когда сама Kafka поддерживает репликацию и может содержать и распространять по своим узлам конфигацию. Разработчики думали также — так появился KRaft.

KRaft, внедрённый в Kafka, использует концепцию контроллеров кворума (англ. Quorum Controllers) для управления конфигурацией. Эти контроллеры обмениваются данными и синхронизируются друг с другом. Изменения в конфигурации считаются успешно выполненными, когда они достигают большинства, то есть (n+1)/2 узлов. Любой из этих узлов впоследствии может предоставить подтверждённую конфигурацию для всего остального кластера.

В каждом кворуме выбирается лидер, задача которого — координировать обновления конфигурации на других серверах. Каждое обновление конфигурации становится новым элементом в логе.

Одним из ключевых преимуществ этой модели является быстрое восстановление работы кластера в случае потери соединения с лидером. Новый лидер выбирается из узлов, которые уже имеют полный лог с метаданными. Это обеспечивает быстрое восстановление, поскольку нет необходимости в медленных операциях обмена данными с внешним хранилищем, таким как Zookeeper, особенно при наличии большого количества топиков и партиций.

В KRaft режиме, Kafka использует Raft-протокол для обеспечения согласованности между брокерами в кластере. Это позволяет обойтись без внешнего координатора состояний, вместо этого все операции управления состоянием выполняются непосредственно в рамках Kafka.

Режим KRaft предлагает ряд преимуществ по сравнению с использованием Zookeeper, включая упрощение операций и улучшение общей производительности, так как устраняется потребность в двойной записи (сначала в Zookeeper, затем в Kafka).

Не все узлы в Kafka по умолчанию превращаются в контроллеры. Задача назначения соответствующих ролей лежит на плечах разработчика. Один из узлов выбирается как лидер, другие узлы, входящие в кворум, становятся участниками голосования (англ. voters), а те, кто не входят в кворум, становятся наблюдателями (англ. observers). Остальные узлы остаются в роли брокеров и не участвуют в процессе хранения метаданных.


## Механизмы репликации партиций

Для обеспечения отказоустойчивости и сохранности данных, Kafka использует механизм репликации. Каждая партиция в Kafka может иметь несколько реплик, распределённых по различным брокерам. Однако не все реплики равны. Одна из реплик выбирается в качестве лидера (англ. leader), а остальные становятся последователями (англ. followers).

Лидер ответственен за все чтения и записи для конкретной партиции, и каждое сообщение, записанное в партицию, сначала записывается в лидера. Затем последователи копируют сообщения из лидера и синхронизируют свое состояние.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/a5c0bee89df29f80a4aa62fea050f3790ac09664/8th-sprint-refactor/02_brokers/images/replica-example_fin.png)


На примере изображен кластер кафка с тремя брокерами и одним топиком. Топик имеет четыре партиции. Каждя партиция имеет три копии. Одна из копий является лидером, остальные — последователями.

In-Sync Replicas (ISR) — это набор реплик партиции, которые находятся в синхронизированном состоянии с лидером. Только реплики, находящиеся в ISR, могут быть выбраны в качестве нового лидера, если текущий лидер не справляется со своими обязанностями (например, если брокер с лидером отключается). Если реплика не синхронизирована достаточно долго (то есть, не может получать данные от лидера), она исключается из ISR.

Использование ISR обеспечивает, что данные не теряются даже в случае отказа брокера, и что чтение всегда возвращает последнее написанное сообщение.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/a5c0bee89df29f80a4aa62fea050f3790ac09664/8th-sprint-refactor/02_brokers/images/in-sync-1_fin.png)

На примере изображен всё тот же кластер. Сейчас все партиции синхронизованы и работают. Представим, что третий брокер вышел из строя.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/a5c0bee89df29f80a4aa62fea050f3790ac09664/8th-sprint-refactor/02_brokers/images/in-sync-2_fin.png)

Теперь, все партиции третьего брокера больше недоступны. Лидер-реплики не могут достучаться до последователей на этом брокере, и через какое-то время они перестают входить в набор синхронизованных реплик. Также произошла автоматическая смена лидера у партиции 2, новый лидером теперь стала партиция на брокере 2.

## Группы потребителей 

В контексте партиционирования топиков становится очень важным понятие групп потребителей (англ. Consumer group) . Потребители в Kafka объединяются в группы для параллельного чтения данных. Каждый потребитель в группе читает данные из одной или нескольких уникальных партиций, тем самым гарантируя, что каждое сообщение будет прочитано одним из потребителей в группе. Группы потребителей позволяют Kafka эффективно распределять обработку данных между несколькими потребителями и обеспечивают масштабируемость обработки данных по мере увеличения количества потребителей в группе. Если один из потребителей выходит из строя, его партиции автоматически перераспределяются между другими потребителями в группе.

![](https://github.com/Yandex-Practicum/python_p2f_ugc/blob/c6e282cb1ce25368e582dac4393379f29886abc2/8th-sprint-refactor/02_brokers/images/UGC_kafka.png)

На рисунке топик состоит из четырёх партиций P1–P4. Два разных клиентских приложения независимо друг от друга через сеть добавляют новые события в партиции топика. События с одним и тем же ключом (на рисунке обозначены одним цветом) добавляются в одну и ту же партицию. Оба клиента также могут добавлять события в одну партицию, если нужно.

Чтобы уберечь данные от потери и сделать их высокодоступными, топики могут реплицироваться. При этом репликацию между разными дата-центрами обычно не делают: это приводит к разрушению кластера. Чтобы реплицировать в другой ЦОД, используют [зеркалирование](https://cwiki.apache.org/confluence/pages/viewpage.action?pageId=27846330){target="_blank"}. Обычно используется фактор репликации 3, то есть всегда будут доступны три копии данных.


### Квиз

Какие из приведённых характеристик брокеров сообщений можно отнести к Kafka?

- [x] Kafka использует pull-based модель чтения сообщений.
 > В отличие от push-based систем, где брокер принимает решение о том, когда доставить сообщение, в Kafka клиенты сами решают, когда и сколько сообщений они хотят извлечь.
- [ ] Kafka автоматически удаляет сообщения после их чтения.
- [ ] Из-за технических особенностей реализации топиков, Kafka не поддерживает гарантию доставки сообщений Exactly once.
- [ ] Kafka не подходит для аналитических задач.

