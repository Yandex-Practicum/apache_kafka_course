resource "yandex_mdb_kafka_cluster" "kafka-cluster" {
  name        = "kafka-cluster"
  environment = "PRODUCTION"
  network_id  = "${var.kafka_network_id}"
  subnet_ids  = "${var.kafka_subnet_id}"

  config {
    version          = "3.6"
    brokers_count    = 3
    zones            = ["${var.kafka_zone}"]
    assign_public_ip = true
    schema_registry  = true
    rest_api {
      enabled = true
    }
    kafka {
      resources {
        resource_preset_id = "s2.micro"
        disk_type_id       = "network-hdd"
        disk_size          = 50
      }
      kafka_config {}
    }
    kraft {
      resources {
        resource_preset_id = "s2.micro"
        disk_type_id       = "network-ssd"
        disk_size          = 50
      }
    }
  }
}

resource "yandex_mdb_kafka_topic" "kafka-topic" {
  cluster_id         = "${yandex_mdb_kafka_cluster.kafka-cluster.id}"
  name               = "kafka-topic"
  partitions         = 3
  replication_factor = 3
  topic_config {
    cleanup_policy        = "CLEANUP_POLICY_DELETE"
    retention_bytes       = 1073741824
    retention_ms          = 86400000
  }
}

resource "yandex_mdb_kafka_user" "kafka-user" {
  cluster_id = "${yandex_mdb_kafka_cluster.kafka-cluster.id}"
  name       = "kafka-user"
  password   = "${var.kafka_password}"

  permission {
    topic_name  = "*"
    role        = "ACCESS_ROLE_ADMIN"
  }
}