## Подготовка

* Установить утилиту yc `curl -sSL https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash`
* Настроить yc (получить токен, выбрать облако и каталог) в соответствии с [инструкцией](https://yandex.cloud/ru/docs/cli/quickstart#install)
* Создать файл terraform.tfvars с со следующим содержимым:
```
yandex_cloud_id=""
yandex_folder_id=""
yandex_cloud_zone="ru-central1-a"
yandex_cloud_token=""
kafka_password=""
```
* Заполнить в нём yandex_cloud_id и yandex_folder_id из вывода команды `yc init`, kafka_password - пароль для пользователя kafka
* Сгенерировать токен через `yc iam create-token` и скопировать в yandex_cloud_token
* Установить terraform по [инструкции](https://yandex.cloud/ru/docs/tutorials/infrastructure-management/terraform-quickstart)
* `cd yandex-cloud`
* Добавить свой ssh-ключ в `scripts/add-ssh-web-app.yaml`
* `terraform init`
* `terraform plan`
* `terraform apply`

## Итог

В облаке будут созданы:
* [Виртуальная машина "kafka-vm"](modules/tf-yc-instance/main.tf) с 2 vCPU/16 Gb RAM для запуска консьюмера и продьюсера из 1 задания и hadoop из 2 задания. Её внешний ip-адрес можно посмотреть в консоли yandex cloud, подключиться можно по ssh. Виртуальная машина будет использовать облачную сеть по умолчанию (default).
* [Кластер Managed Service for Kafka](modules/tf-yc-kafka/main.tf) версии 3.6 на 3 ноды с шаблоном конфигурации хоста s2.micro (2 CPU/8 Gb ram). Для упрощения развёртывания в этом задании кластер развёрнут в одной зоне доступности, чего нельзя делать в продакшен-среде: по рекомендациям, в каждой зоне доступности нужно разворачивать свой кластер и включать зеркалирование. Кластер развёрнут с использованием KRaft
* [Топик kafka-topic](modules/tf-yc-kafka/main.tf) с обозначенными в задании требованиями: 3 раздела, фактор репликации 3, политика очистки delete, максимальный размер раздела 1 Гб, время жизни сегмента 1 день
* [Пользователь kafka-user](modules/tf-yc-kafka/main.tf) с правами admin на все топики (для упрощения; для продакшена лучше сделать отдельных пользователей на конкретные топики)