### Запуск продюсера (producer_app.py)
Продюсер непрерывно отправляет 10 сообщений в Kafka топик (раз в 600 
миллисекунд).

### Запуск Push Consumer (push_consumer_app.py)
Push Consumer слушает новые сообщения и обрабатывает их в реальном времени 
(реализация через очень частый запрос команды consume с низким timeout).

###  Pull Consumer (pull_consumer_app.py)
Pull Consumer периодически проверяет Kafka топик на наличие новых сообщений 
и обрабатывает их. Этот консьюмер будет опрашивать Kafka топик на наличие 
новых сообщений через заданные интервалы времени (10 секунд в данном случае).

### Проверка работы приложения
1. Запустить kafka кластер : ```docker-compose up -d``` (на macOS: ```docker 
compose up -d```)
2. Создать топик: ```docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic test-topic --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2```
3. Запустить `push_consumer_app.py` и `pull_consumer_app.py` в разных окнах
   терминала
4. Запустить `producer_app.py` в 3-ем окне терминала
5. Продюсер отправляет сообщения: в окне терминала продюсера будет 
   напечатано, что продюсер отправил сообщение, и информация об отправке
6. В терминале `push_consumer_app` проверьте, что консьюмер получает 
   сообщения в реальном времени и выводит информацию в лог
7. В терминале `pull_consumer_app` проверьте, что консьюмер получает 
   сообщения с периодичностью и батчами (пачками)
8. Проверка логов Kafka: вы можете проверить логи Kafka для отладки или для 
   уверенности, что сообщения успешно записываются и читаются
