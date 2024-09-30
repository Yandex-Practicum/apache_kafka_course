## Параметры модуля

*  "name" (
  по умолчанию     = "chapter5-std-031-31-vm",
  тип        = string
). Имя ВМ

*  "platform_id" (
  по умолчанию     = "standard-v1",
  тип        = string
). platform_id ВМ

*  "zone" (
  по умолчанию     = "ru-central1-a",
  тип        = string
). Зона ВМ

*  "preemptible" (
  по умолчанию     = false,
  тип        = bool
). Преемственность ВМ

*  "resources_cores" (
  по умолчанию     = 2,
  тип        = number
). Число ядер ВМ

*  "resources_memory" (
  по умолчанию     = 2,
  тип        = number
). Количество памяти ВМ в Гб 

*  "boot_disk_image_id" (
  по умолчанию     = "fd80qm01ah03dkqb14lc",
  тип        = string
). ID исходного образа диска ВМ

*  "boot_disk_size" (
  по умолчанию     = 25,
  тип        = number
). Размер загрузочного диска ВМ

*  "network_interface_subnet_id" (
  по умолчанию     = "e9btneeb3j253qf6156g",
  тип        = string
). ID подсети ВМ

*  "network_interface_nat" (
  по умолчанию     = true,
  тип        = bool
). Включать ли NAT на сетевом интерфейсе 

*  "cloud_config_script_path" (
  по умолчанию     = "scripts/add-ssh-web-app.yaml",
  тип        = string
). Путь до файла cloud-config 

## Вывод модуля

* ip_address -- внутренний IP-адрес созданной ВМ