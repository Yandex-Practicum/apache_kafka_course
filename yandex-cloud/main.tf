module "yandex_cloud_network" {
    source = "./modules/tf-yc-network"
}

module "yandex_cloud_vm" {
    source = "./modules/tf-yc-instance"

    zone = var.yandex_cloud_zone
    network_interface_subnet_id = "${module.yandex_cloud_network.yandex_vpc_subnets[var.yandex_cloud_zone].id}"
}

module "yandex_cloud_kafka" {
    source = "./modules/tf-yc-kafka"

    kafka_zone = var.yandex_cloud_zone
    kafka_network_id = "${module.yandex_cloud_network.yandex_vpc_subnets[var.yandex_cloud_zone].network_id}"
    kafka_subnet_id = [for zone in module.yandex_cloud_network.yandex_vpc_subnets : zone.id]
    kafka_password = var.kafka_password
}
