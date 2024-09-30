## Доработки
* Клиент HDFS создаётся в [consumer/src/main.py](consumer/src/main.py), строка 65
* Запись файла в HDFS осуществляется в [consumer/src/consumer/service.py](consumer/src/consumer/service.py), строка 28

## Подготовка

* Скачать сертификат yandex cloud `wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" --output-document yandex-kafka-production/hadoop-consumer/ca.crt`
* Посмотреть в консоли Yandex Cloud хост для подключения к кластеру Kafka и создать файл `yandex-kafka-production/hadoop-consumer/.env` со следующим содержимым (заменив имя хоста и пароль на свои):
```
CONSUMER_KAFKA_BOOTSTRAP_SERVERS=rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:9091
CONSUMER_SCHEMA_REGISTRY_SERVER=https://kafka-user:...@rc1a-96eh53qo92ss79ic.mdb.yandexcloud.net:443
CONSUMER_KAFKA_PASSWORD=...
```
* Выполнить `cd yandex-kafka-production/hadoop-consumer`
* Выполнить `docker compose up -d`

## Тестирование

* Проверим, что в логах консьюмера нет ошибок (`docker compose logs -f`)
```
hadoop-consumer-1  | 2025-06-01 19:19:54.020 | INFO     | src.consumer.service:process_message:18 - {"topic": "kafka-topic", "partition": 0, "offset": 102, "key": null, "message": {"text": "message"}, "timestamp": 1748804873253}
hadoop-consumer-1  | 2025-06-01 19:19:54.042 | INFO     | src.consumer.service:process_message:30 - {'text': 'message'} written to HDFS as data/message_8e2a3e70-703d-4060-9ba0-ef64a53e06b9
hadoop-consumer-1  | 2025-06-01 19:20:04.003 | INFO     | src.consumer.service:process_message:18 - {"topic": "kafka-topic", "partition": 0, "offset": 103, "key": null, "message": {"text": "message"}, "timestamp": 1748804873293}
hadoop-consumer-1  | 2025-06-01 19:20:04.013 | INFO     | src.consumer.service:process_message:30 - {'text': 'message'} written to HDFS as data/message_8d7a62c7-fe78-44a6-9ede-b3b601245f51
```

* Проверим, что в дашборде Hadoop появились файлы

![files](README_img/scr_files.png)

* Проверим содержимое файла, для этого скачаем один из них

![file](README_img/scr_file.png)

* Можно также посмотреть логи Hadoop на отсутствие ошибок 

```
hadoop-namenode    | 2025-06-01 19:30:54 INFO  StateChange:801 - BLOCK* allocate blk_1073741891_1067, replicas=172.20.0.3:9972, 172.20.0.4:9970, 172.20.0.5:9971 for /user/root/data/message_b0ec6a10-76af-492a-8a55-b8b25ed71b8b
hadoop-datanode-3  | 2025-06-01 19:30:54 INFO  DataNode:759 - Receiving BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067 src: /172.20.0.3:33356 dest: /172.20.0.3:9972
hadoop-datanode-1  | 2025-06-01 19:30:54 INFO  DataNode:759 - Receiving BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067 src: /172.20.0.3:39362 dest: /172.20.0.4:9970
hadoop-datanode-2  | 2025-06-01 19:30:54 INFO  DataNode:759 - Receiving BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067 src: /172.20.0.4:35490 dest: /172.20.0.5:9971
hadoop-datanode-2  | 2025-06-01 19:30:54 INFO  clienttrace:1576 - src: /172.20.0.4:35490, dest: /172.20.0.5:9971, volume: /usr/local/hadoop/hdfs/datanode, bytes: 20, op: HDFS_WRITE, cliID: DFSClient_NONMAPREDUCE_-971631668_65, offset: 0, srvID: d22c0394-9547-4df3-937f-6b141e771f7e, blockid: BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067, duration(ns): 10095531
hadoop-datanode-2  | 2025-06-01 19:30:54 INFO  DataNode:1549 - PacketResponder: BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067, type=LAST_IN_PIPELINE terminating
hadoop-namenode    | 2025-06-01 19:30:54 INFO  BlockStateChange:3777 - BLOCK* addStoredBlock: 172.20.0.5:9971 is added to blk_1073741891_1067 (size=20)
hadoop-datanode-1  | 2025-06-01 19:30:54 INFO  clienttrace:1576 - src: /172.20.0.3:39362, dest: /172.20.0.4:9970, volume: /usr/local/hadoop/hdfs/datanode, bytes: 20, op: HDFS_WRITE, cliID: DFSClient_NONMAPREDUCE_-971631668_65, offset: 0, srvID: acc73f41-048e-48a9-ad8c-c21220544c12, blockid: BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067, duration(ns): 11238194
hadoop-namenode    | 2025-06-01 19:30:54 INFO  BlockStateChange:3777 - BLOCK* addStoredBlock: 172.20.0.4:9970 is added to blk_1073741891_1067 (size=20)
hadoop-datanode-1  | 2025-06-01 19:30:54 INFO  DataNode:1549 - PacketResponder: BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067, type=HAS_DOWNSTREAM_IN_PIPELINE, downstreams=1:[172.20.0.5:9971] terminating
hadoop-namenode    | 2025-06-01 19:30:54 INFO  BlockStateChange:3777 - BLOCK* addStoredBlock: 172.20.0.3:9972 is added to blk_1073741891_1067 (size=20)
hadoop-datanode-3  | 2025-06-01 19:30:54 INFO  clienttrace:1576 - src: /172.20.0.3:33356, dest: /172.20.0.3:9972, volume: /usr/local/hadoop/hdfs/datanode, bytes: 20, op: HDFS_WRITE, cliID: DFSClient_NONMAPREDUCE_-971631668_65, offset: 0, srvID: 463d9bea-f8cf-4646-882a-af6de74fe3b9, blockid: BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067, duration(ns): 17312899
hadoop-datanode-3  | 2025-06-01 19:30:54 INFO  DataNode:1549 - PacketResponder: BP-867886778-172.20.0.2-1748805326388:blk_1073741891_1067, type=HAS_DOWNSTREAM_IN_PIPELINE, downstreams=2:[172.20.0.4:9970, 172.20.0.5:9971] terminating
hadoop-namenode    | 2025-06-01 19:30:54 INFO  StateChange:3239 - DIR* completeFile: /user/root/data/message_b0ec6a10-76af-492a-8a55-b8b25ed71b8b is closed by DFSClient_NONMAPREDUCE_-971631668_65
```