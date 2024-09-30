# Как должен выглядеть результат
# Развёрнут один инструмент (NiFi, Hadoop).

cd task2
docker compose up -d

http://localhost:8080/nifi


# Настроено взаимодействие между инструментом и Kafka-кластером.
# Выполненное задание 2 должно включать:
# Скриншот из консоли с запущенными сервисами.

docker ps
CONTAINER ID   IMAGE                COMMAND                 CREATED          STATUS          PORTS                                                   NAMES
7af0db179034   apache/nifi:1.21.0   "../scripts/start.sh"   21 seconds ago   Up 20 seconds   8000/tcp, 8443/tcp, 10000/tcp, 0.0.0.0:8080->8080/tcp   nifi

/opt/nifi/nifi-current/cert/truststore.jks

# Конфигурационные файлы запуска.
./task2/docker-compose.yaml

# Код продюсера и консьюмера.

Код продюсера и консьюмера nifi сохранен в шаблоне
./task2/publish_and_consume.xml

# Логи успешной передачи данных:

Можно посомтреть в NiFi Data Provenance

# Kafka-топик с поступающими данными (kafka-console-consumer.sh);


kcat -C \
         -b rc1a-v868adntv9bpbg1i.mdb.yandexcloud.net:9091 \
         -t task2 \
         -X security.protocol=SASL_SSL \
         -X sasl.mechanism=SCRAM-SHA-512 \
         -X sasl.username="user" \
         -X sasl.password="user1234" \
         -X ssl.ca.location=./YandexInternalRootCA.crt -Z -K:
NULL:id,name,age,email 1,Екатерина,21,kate@yandex.ru 2,Никита,26,nikita@yandex.ru 3,Майя,21,maya@yandex.ru 4,Алексей,26,alex@yandex.ru 5,Илья,25,ilya@yandex.ru 6,Виктория,22,vika@yandex.ru 
NULL:id,name,age,email 1,Екатерина,21,kate@yandex.ru 2,Никита,26,nikita@yandex.ru 3,Майя,21,maya@yandex.ru 4,Алексей,26,alex@yandex.ru 5,Илья,25,ilya@yandex.ru 6,Виктория,22,vika@yandex.ru 
NULL:id,name,age,email 1,Екатерина,21,kate@yandex.ru 2,Никита,26,nikita@yandex.ru 3,Майя,21,maya@yandex.ru 4,Алексей,26,alex@yandex.ru 5,Илья,25,ilya@yandex.ru 6,Виктория,22,vika@yandex.ru 
NULL:id,name,age,email 1,Екатерина,21,kate@yandex.ru 2,Никита,26,nikita@yandex.ru 3,Майя,21,maya@yandex.ru 4,Алексей,26,alex@yandex.ru 5,Илья,25,ilya@yandex.ru 6,Виктория,22,vika@yandex.ru 
NULL:id,name,age,email 1,Екатерина,21,kate@yandex.ru 2,Никита,26,nikita@yandex.ru 3,Майя,21,maya@yandex.ru 4,Алексей,26,alex@yandex.ru 5,Илья,25,ilya@yandex.ru 6,Виктория,22,vika@yandex.ru 
NULL:id,name,age,email 1,Екатерина,21,kate@yandex.ru 2,Никита,26,nikita@yandex.ru 3,Майя,21,maya@yandex.ru 4,Алексей,26,alex@yandex.ru 5,Илья,25,ilya@yandex.ru 6,Виктория,22,vika@yandex.ru 
...
% Reached end of topic task2 [0] at offset 667
% Reached end of topic task2 [2] at offset 667
% Reached end of topic task2 [1] at offset 668


yc managed-kafka topic list --cluster-name kafka585
+-------+------------------+--------------------+
| NAME  | PARTITIONS COUNT | REPLICATION FACTOR |
+-------+------------------+--------------------+
| task2 |                3 |                  3 |
| topic |                3 |                  3 |
+-------+------------------+--------------------+


yc managed-kafka topic get topic --cluster-name kafka585
name: topic
cluster_id: c9qhd09k2eiqf6eq5hp5
partitions: "3"
replication_factor: "3"
topic_config_3:
  cleanup_policy: CLEANUP_POLICY_DELETE
  retention_bytes: "10000"
  retention_ms: "10000"

yc managed-kafka topic get task2 --cluster-name kafka585
name: task2
cluster_id: c9qhd09k2eiqf6eq5hp5
partitions: "3"
replication_factor: "3"
topic_config_3: {}


# логи успешной работы NiFi или Hadoop (вывод в консоли);

2025-05-01 18:04:11,729 INFO [Write-Ahead Local State Provider Maintenance] org.wali.MinimalLockingWriteAheadLog org.wali.MinimalLockingWriteAheadLog@419ae488 checkpointed with 3 Records and 0 Swap Files in 3 milliseconds (Stop-the-world time = 0 milliseconds, Clear Edit Logs time = 0 millis), max Transaction ID 4
2025-05-01 18:04:13,554 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:04:13,561 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 6 milliseconds (Stop-the-world time = 2 milliseconds), max Transaction ID 176
2025-05-01 18:04:13,561 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 6 milliseconds
2025-05-01 18:04:33,563 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:04:33,577 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 180
2025-05-01 18:04:33,577 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 13 milliseconds
2025-05-01 18:04:53,580 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:04:53,594 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 14 milliseconds (Stop-the-world time = 10 milliseconds), max Transaction ID 184
2025-05-01 18:04:53,595 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 14 milliseconds
2025-05-01 18:05:11,668 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Successfully deleted 0 files (0 bytes) from archive
2025-05-01 18:05:11,669 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Archive cleanup completed for container default; will now allow writing to this container. Bytes used = 93.2 GB, bytes free = 913.65 GB, capacity = 1,006.85 GB
2025-05-01 18:05:13,596 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:05:13,611 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 15 milliseconds (Stop-the-world time = 10 milliseconds), max Transaction ID 188
2025-05-01 18:05:13,611 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 15 milliseconds
2025-05-01 18:05:33,614 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:05:33,619 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 4 milliseconds (Stop-the-world time = 0 milliseconds), max Transaction ID 192
2025-05-01 18:05:33,619 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 4 milliseconds
2025-05-01 18:05:53,620 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:05:53,636 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 15 milliseconds (Stop-the-world time = 10 milliseconds), max Transaction ID 196
2025-05-01 18:05:53,636 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 16 milliseconds
2025-05-01 18:06:11,674 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Successfully deleted 0 files (0 bytes) from archive
2025-05-01 18:06:11,674 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Archive cleanup completed for container default; will now allow writing to this container. Bytes used = 93.2 GB, bytes free = 913.65 GB, capacity = 1,006.85 GB
2025-05-01 18:06:11,749 INFO [Write-Ahead Local State Provider Maintenance] org.wali.MinimalLockingWriteAheadLog org.wali.MinimalLockingWriteAheadLog@419ae488 checkpointed with 3 Records and 0 Swap Files in 14 milliseconds (Stop-the-world time = 1 milliseconds, Clear Edit Logs time = 2 millis), max Transaction ID 4
2025-05-01 18:06:13,637 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:06:13,652 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 14 milliseconds (Stop-the-world time = 9 milliseconds), max Transaction ID 200
2025-05-01 18:06:13,652 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 14 milliseconds
2025-05-01 18:06:33,655 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:06:33,669 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 204
2025-05-01 18:06:33,669 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 14 milliseconds
2025-05-01 18:06:53,670 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:06:53,684 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 9 milliseconds), max Transaction ID 208
2025-05-01 18:06:53,685 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 14 milliseconds
2025-05-01 18:07:11,681 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Successfully deleted 0 files (0 bytes) from archive
2025-05-01 18:07:11,682 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Archive cleanup completed for container default; will now allow writing to this container. Bytes used = 93.2 GB, bytes free = 913.65 GB, capacity = 1,006.85 GB
2025-05-01 18:07:13,685 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:07:13,692 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 6 milliseconds (Stop-the-world time = 2 milliseconds), max Transaction ID 212
2025-05-01 18:07:13,692 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 6 milliseconds
2025-05-01 18:07:33,693 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:07:33,699 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 5 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 216
2025-05-01 18:07:33,699 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 5 milliseconds
2025-05-01 18:07:53,701 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:07:53,711 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 10 milliseconds (Stop-the-world time = 0 milliseconds), max Transaction ID 220
2025-05-01 18:07:53,711 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 10 milliseconds
2025-05-01 18:08:11,688 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Successfully deleted 0 files (0 bytes) from archive
2025-05-01 18:08:11,688 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Archive cleanup completed for container default; will now allow writing to this container. Bytes used = 93.2 GB, bytes free = 913.65 GB, capacity = 1,006.85 GB
2025-05-01 18:08:11,772 INFO [Write-Ahead Local State Provider Maintenance] org.wali.MinimalLockingWriteAheadLog org.wali.MinimalLockingWriteAheadLog@419ae488 checkpointed with 3 Records and 0 Swap Files in 17 milliseconds (Stop-the-world time = 2 milliseconds, Clear Edit Logs time = 1 millis), max Transaction ID 4
2025-05-01 18:08:13,712 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:08:13,724 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 11 milliseconds (Stop-the-world time = 9 milliseconds), max Transaction ID 224
2025-05-01 18:08:13,724 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 12 milliseconds
2025-05-01 18:08:33,725 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:08:33,740 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 9 milliseconds), max Transaction ID 228
2025-05-01 18:08:33,740 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 14 milliseconds
2025-05-01 18:08:53,741 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:08:53,754 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 12 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 232
2025-05-01 18:08:53,755 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 12 milliseconds
2025-05-01 18:09:11,696 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Successfully deleted 0 files (0 bytes) from archive
2025-05-01 18:09:11,697 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Archive cleanup completed for container default; will now allow writing to this container. Bytes used = 93.2 GB, bytes free = 913.65 GB, capacity = 1,006.85 GB
2025-05-01 18:09:13,755 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:09:13,760 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 4 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 236
2025-05-01 18:09:13,760 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 4 milliseconds
2025-05-01 18:09:33,762 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:09:33,775 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 9 milliseconds), max Transaction ID 240
2025-05-01 18:09:33,775 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 13 milliseconds
2025-05-01 18:09:53,777 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:09:53,791 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 14 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 244
2025-05-01 18:09:53,791 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 14 milliseconds
2025-05-01 18:10:11,705 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Successfully deleted 0 files (0 bytes) from archive
2025-05-01 18:10:11,705 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Archive cleanup completed for container default; will now allow writing to this container. Bytes used = 93.2 GB, bytes free = 913.65 GB, capacity = 1,006.85 GB
2025-05-01 18:10:11,793 INFO [Write-Ahead Local State Provider Maintenance] org.wali.MinimalLockingWriteAheadLog org.wali.MinimalLockingWriteAheadLog@419ae488 checkpointed with 3 Records and 0 Swap Files in 15 milliseconds (Stop-the-world time = 1 milliseconds, Clear Edit Logs time = 1 millis), max Transaction ID 4
2025-05-01 18:10:13,792 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:10:13,806 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 14 milliseconds (Stop-the-world time = 9 milliseconds), max Transaction ID 248
2025-05-01 18:10:13,807 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 14 milliseconds
2025-05-01 18:10:33,808 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:10:33,820 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 12 milliseconds (Stop-the-world time = 8 milliseconds), max Transaction ID 252
2025-05-01 18:10:33,820 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 12 milliseconds
2025-05-01 18:10:48,301 INFO [Timer-Driven Process Thread-7] o.a.n.p.store.WriteAheadStorePartition Successfully rolled over Event Writer for Provenance Event Store Partition[directory=./provenance_repository] due to MAX_TIME_REACHED. Event File was 89.96 KB and contained 240 events.
2025-05-01 18:10:53,822 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:10:53,836 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 256
2025-05-01 18:10:53,836 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 13 milliseconds
2025-05-01 18:11:11,712 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Successfully deleted 0 files (0 bytes) from archive
2025-05-01 18:11:11,713 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Archive cleanup completed for container default; will now allow writing to this container. Bytes used = 93.2 GB, bytes free = 913.65 GB, capacity = 1,006.85 GB
2025-05-01 18:11:13,837 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:11:13,851 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 260
2025-05-01 18:11:13,851 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 13 milliseconds
2025-05-01 18:11:33,853 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:11:33,858 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 4 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 264
2025-05-01 18:11:33,858 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 4 milliseconds
2025-05-01 18:11:53,859 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:11:53,873 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 9 milliseconds), max Transaction ID 268
2025-05-01 18:11:53,874 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 14 milliseconds
2025-05-01 18:12:11,720 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Successfully deleted 0 files (0 bytes) from archive
2025-05-01 18:12:11,720 INFO [Cleanup Archive for default] o.a.n.c.repository.FileSystemRepository Archive cleanup completed for container default; will now allow writing to this container. Bytes used = 93.2 GB, bytes free = 913.65 GB, capacity = 1,006.85 GB
2025-05-01 18:12:11,815 INFO [Write-Ahead Local State Provider Maintenance] org.wali.MinimalLockingWriteAheadLog org.wali.MinimalLockingWriteAheadLog@419ae488 checkpointed with 3 Records and 0 Swap Files in 16 milliseconds (Stop-the-world time = 1 milliseconds, Clear Edit Logs time = 1 millis), max Transaction ID 4
2025-05-01 18:12:13,874 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Initiating checkpoint of FlowFile Repository
2025-05-01 18:12:13,888 INFO [pool-7-thread-1] o.a.n.wali.SequentialAccessWriteAheadLog Checkpointed Write-Ahead Log with 0 Records and 0 Swap Files in 13 milliseconds (Stop-the-world time = 1 milliseconds), max Transaction ID 272
2025-05-01 18:12:13,888 INFO [pool-7-thread-1] o.a.n.c.r.WriteAheadFlowFileRepository Successfully checkpointed FlowFile Repository with 0 records in 13 milliseconds

# подтверждение записи данных (если используется HDFS или БД).

Не используется HDFS или БД