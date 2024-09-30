variable "kafka_zone" {
  type        = string
  description = "VM zone"
} 

variable "kafka_network_id" {
  type        = string
  description = "Network ID"
} 

variable "kafka_subnet_id" {
  type        = list(string)
  description = "Network interface subnet ID"
} 

variable "kafka_password" {
  type        = string
  description = "Kafka user password"
} 