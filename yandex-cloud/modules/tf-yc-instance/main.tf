resource "yandex_compute_instance" "vm-1" {
    name        = var.name
    platform_id = var.platform_id
    zone        = var.zone
    allow_stopping_for_update = true


    scheduling_policy {
        preemptible = var.preemptible
    }

    # Конфигурация ресурсов:
    # количество процессоров и оперативной памяти
    resources {
        cores  = var.resources_cores
        memory = var.resources_memory
    }

    # Загрузочный диск:
    # здесь указывается образ операционной системы
    # для новой виртуальной машины
    boot_disk {
        initialize_params {
            image_id = var.boot_disk_image_id
            size     = var.boot_disk_size
        }
    }

    # Сетевой интерфейс:
    # нужно указать идентификатор подсети, к которой будет подключена ВМ
    network_interface {
        subnet_id = var.network_interface_subnet_id
        nat       = var.network_interface_nat
    }

    # Метаданные машины:
    # здесь можно указать скрипт, который запустится при создании ВМ
    # или список SSH-ключей для доступа на ВМ
    metadata = {
        user-data = "${file(var.cloud_config_script_path)}"
    }
} 
