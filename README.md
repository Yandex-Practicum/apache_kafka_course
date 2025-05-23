# Kafka connect (module 3)

## Задание 1. Оптимизация параметров для повышения пропускной способности JDBC Source Connector

Небольшого приростра скорости записи удалось добиться путем изменения параметра batch.size, linger.ms и включения сжатия.
```
batch.size = batch.max.rows * record_size_average_in_bytes = 100 * 140b = 14000
```

| Эксперимент | batch.size | linger.ms | compression.type | buffer.memory | Record Write Rate |
|-------------|------------|-----------|------------------|---------------|-------------------|
| 1           | 100        | 0         | none             | 33554432      | 171               |
| 2           | 14000      | 0         | snappy           | 33554432      | 175               |
| 3           | 14000      | 200       | snappy           | 33554432      | 179               |

## Задание 2. Создание собственного коннектора для переноса данных из Apache Kafka в Prometheus

Пример .env файла:
```
KAFKA_BOOTSTRAP_SERVERS=127.0.0.1:9094
TOPIC=metrics
KAFKA_CONSUMER_PULL_TIMEOUT_MS=100
KAFKA_CUSTOMER_GROUP_1=consumer_group_1
KAFKA_SESSION_TIMEOUT_MS=6000
SCHEMA_REGISTRY_URL=http://127.0.0.1:8081
REDIS_OPT=redis://default:my-password@localhost:26379/0
```

Запуск окружения:
```
make up
```

Создание топика:
```
make create-topic
```

Создание коннектора:
```
make prometheus
```

## Задание 3. Получение лога вывода Debezium PostgresConnector

```
[2025-01-27 10:39:50,865] INFO Starting PostgresConnectorTask with configuration:
2025-01-27T10:39:50.865681091Z    connector.class = io.debezium.connector.postgresql.PostgresConnector
2025-01-27T10:39:50.865682341Z    database.user = postgres-user
2025-01-27T10:39:50.865683299Z    database.dbname = customers
2025-01-27T10:39:50.865684091Z    transforms.unwrap.delete.handling.mode = rewrite
2025-01-27T10:39:50.865685007Z    topic.creation.default.partitions = -1
2025-01-27T10:39:50.865685799Z    transforms = unwrap
2025-01-27T10:39:50.865686674Z    database.server.name = customers
2025-01-27T10:39:50.865693466Z    database.port = 5432
2025-01-27T10:39:50.865694466Z    table.whitelist = public.customers
2025-01-27T10:39:50.865695257Z    topic.creation.enable = true
2025-01-27T10:39:50.865696257Z    topic.prefix = customers
2025-01-27T10:39:50.865697132Z    task.class = io.debezium.connector.postgresql.PostgresConnectorTask
2025-01-27T10:39:50.865697965Z    database.hostname = postgres
2025-01-27T10:39:50.865698757Z    database.password = ********
2025-01-27T10:39:50.865699549Z    transforms.unwrap.drop.tombstones = false
2025-01-27T10:39:50.865700590Z    name = debezium-source
2025-01-27T10:39:50.865701340Z    topic.creation.default.replication.factor = -1
2025-01-27T10:39:50.865702174Z    transforms.unwrap.type = io.debezium.transforms.ExtractNewRecordState
2025-01-27T10:39:50.865703257Z    skipped.operations = none
2025-01-27T10:39:50.865704007Z    value.converter = org.apache.kafka.connect.storage.StringConverter
2025-01-27T10:39:50.865704840Z    key.converter = org.apache.kafka.connect.storage.StringConverter
2025-01-27T10:39:50.865705715Z  (io.debezium.connector.common.BaseSourceTask)
2025-01-27T10:39:50.865792923Z [2025-01-27 10:39:50,865] INFO Loading the custom source info struct maker plugin: io.debezium.connector.postgresql.PostgresSourceInfoStructMaker (io.debezium.config.CommonConnectorConfig)
2025-01-27T10:39:50.866080170Z [2025-01-27 10:39:50,866] INFO Loading the custom topic naming strategy plugin: io.debezium.schema.SchemaTopicNamingStrategy (io.debezium.config.CommonConnectorConfig)
2025-01-27T10:39:50.871267702Z [2025-01-27 10:39:50,871] INFO Connection gracefully closed (io.debezium.jdbc.JdbcConnection)
2025-01-27T10:39:50.888306034Z [2025-01-27 10:39:50,886] WARN Type [oid:13529, name:_pg_user_mappings] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.888322117Z [2025-01-27 10:39:50,886] WARN Type [oid:13203, name:cardinal_number] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.888323784Z [2025-01-27 10:39:50,886] WARN Type [oid:13206, name:character_data] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.888324825Z [2025-01-27 10:39:50,886] WARN Type [oid:13208, name:sql_identifier] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.888325784Z [2025-01-27 10:39:50,886] WARN Type [oid:13214, name:time_stamp] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.888326700Z [2025-01-27 10:39:50,886] WARN Type [oid:13216, name:yes_or_no] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.902906639Z [2025-01-27 10:39:50,902] INFO No previous offsets found (io.debezium.connector.common.BaseSourceTask)
2025-01-27T10:39:50.911843801Z [2025-01-27 10:39:50,911] WARN Type [oid:13529, name:_pg_user_mappings] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.911853718Z [2025-01-27 10:39:50,911] WARN Type [oid:13203, name:cardinal_number] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.911880176Z [2025-01-27 10:39:50,911] WARN Type [oid:13206, name:character_data] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.911979133Z [2025-01-27 10:39:50,911] WARN Type [oid:13208, name:sql_identifier] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.911980675Z [2025-01-27 10:39:50,911] WARN Type [oid:13214, name:time_stamp] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.911981841Z [2025-01-27 10:39:50,911] WARN Type [oid:13216, name:yes_or_no] is already mapped (io.debezium.connector.postgresql.TypeRegistry)
2025-01-27T10:39:50.917818784Z [2025-01-27 10:39:50,917] INFO Connector started for the first time. (io.debezium.connector.common.BaseSourceTask)
2025-01-27T10:39:50.918187197Z [2025-01-27 10:39:50,918] INFO No previous offset found (io.debezium.connector.postgresql.PostgresConnectorTask)
2025-01-27T10:39:50.920577173Z [2025-01-27 10:39:50,920] INFO user 'postgres-user' connected to database 'customers' on PostgreSQL 16.4 (Debian 16.4-1.pgdg110+2) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 10.2.1-6) 10.2.1 20210110, 64-bit with roles:
2025-01-27T10:39:50.920589798Z 	role 'pg_read_all_settings' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920591590Z 	role 'pg_database_owner' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920592756Z 	role 'pg_stat_scan_tables' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920593756Z 	role 'pg_checkpoint' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920594965Z 	role 'pg_write_server_files' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920595923Z 	role 'pg_use_reserved_connections' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920596881Z 	role 'postgres-user' [superuser: true, replication: true, inherit: true, create role: true, create db: true, can log in: true]
2025-01-27T10:39:50.920598006Z 	role 'pg_read_all_data' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920599048Z 	role 'pg_write_all_data' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920599923Z 	role 'pg_monitor' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920600840Z 	role 'pg_read_server_files' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920601756Z 	role 'pg_create_subscription' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920609256Z 	role 'pg_execute_server_program' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920610756Z 	role 'pg_read_all_stats' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false]
2025-01-27T10:39:50.920611756Z 	role 'pg_signal_backend' [superuser: false, replication: false, inherit: true, create role: false, create db: false, can log in: false] (io.debezium.connector.postgresql.PostgresConnectorTask)
2025-01-27T10:39:50.922985483Z [2025-01-27 10:39:50,922] INFO Obtained valid replication slot ReplicationSlot [active=false, latestFlushedLsn=null, catalogXmin=null] (io.debezium.connector.postgresql.connection.PostgresConnection)
2025-01-27T10:39:50.928790009Z [2025-01-27 10:39:50,928] INFO Creating replication slot with command CREATE_REPLICATION_SLOT "debezium"  LOGICAL decoderbufs  (io.debezium.connector.postgresql.connection.PostgresReplicationConnection)
2025-01-27T10:39:50.937827961Z [2025-01-27 10:39:50,937] INFO Requested thread factory for component PostgresConnector, id = customers named = SignalProcessor (io.debezium.util.Threads)
2025-01-27T10:39:50.942779454Z [2025-01-27 10:39:50,942] INFO Requested thread factory for component PostgresConnector, id = customers named = change-event-source-coordinator (io.debezium.util.Threads)
2025-01-27T10:39:50.942787537Z [2025-01-27 10:39:50,942] INFO Requested thread factory for component PostgresConnector, id = customers named = blocking-snapshot (io.debezium.util.Threads)
2025-01-27T10:39:50.944478937Z [2025-01-27 10:39:50,944] INFO Creating thread debezium-postgresconnector-customers-change-event-source-coordinator (io.debezium.util.Threads)
2025-01-27T10:39:50.945837840Z [2025-01-27 10:39:50,944] INFO WorkerSourceTask{id=debezium-source-0} Source task finished initialization and start (org.apache.kafka.connect.runtime.AbstractWorkerSourceTask)
2025-01-27T10:39:50.945852423Z [2025-01-27 10:39:50,945] INFO Metrics registered (io.debezium.pipeline.ChangeEventSourceCoordinator)
2025-01-27T10:39:50.946124671Z [2025-01-27 10:39:50,946] INFO Context created (io.debezium.pipeline.ChangeEventSourceCoordinator)
2025-01-27T10:39:50.947483532Z [2025-01-27 10:39:50,947] INFO According to the connector configuration data will be snapshotted (io.debezium.connector.postgresql.PostgresSnapshotChangeEventSource)
2025-01-27T10:39:50.948414148Z [2025-01-27 10:39:50,948] INFO Snapshot step 1 - Preparing (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:50.948921101Z [2025-01-27 10:39:50,948] INFO Setting isolation level (io.debezium.connector.postgresql.PostgresSnapshotChangeEventSource)
2025-01-27T10:39:50.949150308Z [2025-01-27 10:39:50,949] INFO Opening transaction with statement SET TRANSACTION ISOLATION LEVEL REPEATABLE READ; 
2025-01-27T10:39:50.949153516Z SET TRANSACTION SNAPSHOT '00000006-00000002-1'; (io.debezium.connector.postgresql.PostgresSnapshotChangeEventSource)
2025-01-27T10:39:50.977740192Z [2025-01-27 10:39:50,977] INFO Snapshot step 2 - Determining captured tables (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:50.978803056Z [2025-01-27 10:39:50,978] INFO Adding table public.users to the list of capture schema tables (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:50.980192876Z [2025-01-27 10:39:50,980] INFO Created connection pool with 1 threads (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:50.980198251Z [2025-01-27 10:39:50,980] INFO Snapshot step 3 - Locking captured tables [public.users] (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:50.981936275Z [2025-01-27 10:39:50,981] INFO Snapshot step 4 - Determining snapshot offset (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:50.982389312Z [2025-01-27 10:39:50,982] INFO Creating initial offset context (io.debezium.connector.postgresql.PostgresOffsetContext)
2025-01-27T10:39:50.983488010Z [2025-01-27 10:39:50,983] INFO Read xlogStart at 'LSN{0/1BDF3C0}' from transaction '747' (io.debezium.connector.postgresql.PostgresOffsetContext)
2025-01-27T10:39:50.986212900Z [2025-01-27 10:39:50,986] INFO Read xlogStart at 'LSN{0/1BDF3C0}' from transaction '747' (io.debezium.connector.postgresql.PostgresSnapshotChangeEventSource)
2025-01-27T10:39:50.986218150Z [2025-01-27 10:39:50,986] INFO Snapshot step 5 - Reading structure of captured tables (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:50.986339982Z [2025-01-27 10:39:50,986] INFO Reading structure of schema 'public' of catalog 'customers' (io.debezium.connector.postgresql.PostgresSnapshotChangeEventSource)
2025-01-27T10:39:51.000139554Z [2025-01-27 10:39:51,000] INFO Snapshot step 6 - Persisting schema history (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:51.000231261Z [2025-01-27 10:39:51,000] INFO Snapshot step 7 - Snapshotting data (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:51.000586299Z [2025-01-27 10:39:51,000] INFO Creating snapshot worker pool with 1 worker thread(s) (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:51.001740288Z [2025-01-27 10:39:51,001] INFO For table 'public.users' using select statement: 'SELECT "id", "name", "updated_at", "private_info" FROM "public"."users"' (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:51.003625061Z [2025-01-27 10:39:51,003] INFO Exporting data from table 'public.users' (1 of 1 tables) (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:51.082878485Z [2025-01-27 10:39:51,082] INFO The task will send records to topic 'customers.public.users' for the first time. Checking whether topic exists (org.apache.kafka.connect.runtime.AbstractWorkerSourceTask)
2025-01-27T10:39:51.086463492Z [2025-01-27 10:39:51,086] INFO Creating topic 'customers.public.users' (org.apache.kafka.connect.runtime.AbstractWorkerSourceTask)
2025-01-27T10:39:51.131892709Z [2025-01-27 10:39:51,131] INFO Created topic (name=customers.public.users, numPartitions=-1, replicationFactor=-1, replicasAssignments=null, configs={}) on brokers at kafka-0:9092 (org.apache.kafka.connect.util.TopicAdmin)
2025-01-27T10:39:51.131912334Z [2025-01-27 10:39:51,131] INFO Created topic '(name=customers.public.users, numPartitions=-1, replicationFactor=-1, replicasAssignments=null, configs={})' using creation group TopicCreationGroup{name='default', inclusionPattern=.*, exclusionPattern=, numPartitions=-1, replicationFactor=-1, otherConfigs={}} (org.apache.kafka.connect.runtime.AbstractWorkerSourceTask)
2025-01-27T10:39:51.133884398Z [2025-01-27 10:39:51,133] WARN [Producer clientId=connector-producer-debezium-source-0] Error while fetching metadata with correlation id 4 : {customers.public.users=UNKNOWN_TOPIC_OR_PARTITION} (org.apache.kafka.clients.NetworkClient)
2025-01-27T10:39:51.134064604Z [2025-01-27 10:39:51,134] INFO 	 Finished exporting 10000 records for table 'public.users' (1 of 1 tables); total duration '00:00:00.13' (io.debezium.relational.RelationalSnapshotChangeEventSource)
2025-01-27T10:39:51.135276509Z [2025-01-27 10:39:51,135] INFO Snapshot - Final stage (io.debezium.pipeline.source.AbstractSnapshotChangeEventSource)
2025-01-27T10:39:51.135292592Z [2025-01-27 10:39:51,135] INFO Snapshot completed (io.debezium.pipeline.source.AbstractSnapshotChangeEventSource)
2025-01-27T10:39:51.138942514Z [2025-01-27 10:39:51,138] INFO Snapshot ended with SnapshotResult [status=COMPLETED, offset=PostgresOffsetContext [sourceInfoSchema=Schema{io.debezium.connector.postgresql.Source:STRUCT}, sourceInfo=source_info[server='customers'db='customers', lsn=LSN{0/1BDF3C0}, txId=747, timestamp=2025-01-27T10:39:50.950686Z, snapshot=FALSE, schema=public, table=users], lastSnapshotRecord=true, lastCompletelyProcessedLsn=null, lastCommitLsn=null, streamingStoppingLsn=null, transactionContext=TransactionContext [currentTransactionId=null, perTableEventCount={}, totalEventCount=0], incrementalSnapshotContext=IncrementalSnapshotContext [windowOpened=false, chunkEndPosition=null, dataCollectionsToSnapshot=[], lastEventKeySent=null, maximumKey=null]]] (io.debezium.pipeline.ChangeEventSourceCoordinator)
2025-01-27T10:39:51.140215252Z [2025-01-27 10:39:51,140] INFO Connected metrics set to 'true' (io.debezium.pipeline.ChangeEventSourceCoordinator)
2025-01-27T10:39:51.151591931Z [2025-01-27 10:39:51,151] INFO REPLICA IDENTITY for 'public.users' is 'DEFAULT'; UPDATE and DELETE events will contain previous values only for PK columns (io.debezium.connector.postgresql.PostgresSchema)
2025-01-27T10:39:51.154801274Z [2025-01-27 10:39:51,154] INFO SignalProcessor started. Scheduling it every 5000ms (io.debezium.pipeline.signal.SignalProcessor)
2025-01-27T10:39:51.154956689Z [2025-01-27 10:39:51,154] INFO Creating thread debezium-postgresconnector-customers-SignalProcessor (io.debezium.util.Threads)
2025-01-27T10:39:51.155111729Z [2025-01-27 10:39:51,155] INFO Starting streaming (io.debezium.pipeline.ChangeEventSourceCoordinator)
2025-01-27T10:39:51.155117188Z [2025-01-27 10:39:51,155] INFO Retrieved latest position from stored offset 'LSN{0/1BDF3C0}' (io.debezium.connector.postgresql.PostgresStreamingChangeEventSource)
2025-01-27T10:39:51.155353769Z [2025-01-27 10:39:51,155] INFO Looking for WAL restart position for last commit LSN 'null' and last change LSN 'LSN{0/1BDF3C0}' (io.debezium.connector.postgresql.connection.WalPositionLocator)
2025-01-27T10:39:51.160115138Z [2025-01-27 10:39:51,160] INFO Obtained valid replication slot ReplicationSlot [active=false, latestFlushedLsn=LSN{0/1BDF3C0}, catalogXmin=747] (io.debezium.connector.postgresql.connection.PostgresConnection)
2025-01-27T10:39:51.160505134Z [2025-01-27 10:39:51,160] INFO Connection gracefully closed (io.debezium.jdbc.JdbcConnection)
2025-01-27T10:39:51.173873877Z [2025-01-27 10:39:51,173] INFO Requested thread factory for component PostgresConnector, id = customers named = keep-alive (io.debezium.util.Threads)
2025-01-27T10:39:51.174066458Z [2025-01-27 10:39:51,173] INFO Creating thread debezium-postgresconnector-customers-keep-alive (io.debezium.util.Threads)
2025-01-27T10:39:51.184404731Z [2025-01-27 10:39:51,184] INFO REPLICA IDENTITY for 'public.users' is 'DEFAULT'; UPDATE and DELETE events will contain previous values only for PK columns (io.debezium.connector.postgresql.PostgresSchema)
2025-01-27T10:39:51.184839852Z [2025-01-27 10:39:51,184] INFO Processing messages (io.debezium.connector.postgresql.PostgresStreamingChangeEventSource)
```