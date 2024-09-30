## Задание 1

* В первую очередь была развёрнута инфраструктура в yandex cloud при помощи terraform, команды для запуска описаны в [README](yandex-cloud/README.md) папки "yandex-cloud", в которой хранится конфигурация terraform

* Далее были из [предыдущей практической работы](https://github.com/sakharovaan/yandex-kafka-security) перенесены consumer и producer в папку "app", подключены к облачной Kafka, действия для подключения и ход тестирования описаны в [README](app/README.md)

## Задание 2

* Был развёрнут кластер Hadoop на ВМ в Yandex Cloud по конфигурации из теоретической части (папка "hadoop"), описание в [README](hadoop/README.md)

* Был доработан consumer (папка "hadoop-consumer"), чтобы он не только выводил в консоль полученное сообщение, но и записывал его как файл в HDFS, продьюсер остался тем же. Детали и тестирование в [README](hadoop-consumer/README.md)