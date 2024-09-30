variable "name" {
  default     = "kafka-vm"
  type        = string
  description = "VM name"
} 

variable "platform_id" {
  default     = "standard-v1"
  type        = string
  description = "VM platform_id"
} 

variable "zone" {
  default     = "ru-central1-a"
  type        = string
  description = "VM zone"
} 

variable "preemptible" {
  default     = false
  type        = bool
  description = "Is VM preemptible or not"
} 

variable "resources_cores" {
  default     = 2
  type        = number
  description = "Number of VM cores"
}

variable "resources_memory" {
  default     = 16
  type        = number
  description = "Amount of VM memory in GB"
} 

variable "boot_disk_image_id" {
  default     = "fd80qm01ah03dkqb14lc"
  type        = string
  description = "Boot disk image ID"
} 

variable "boot_disk_size" {
  default     = 35
  type        = number
  description = "Boot disk size"
} 

variable "network_interface_subnet_id" {
  default     = "e9btneeb3j253qf6156g"
  type        = string
  description = "Network interface subnet ID"
} 

variable "network_interface_nat" {
  default     = true
  type        = bool
  description = "Network interface nat parameter"
} 

variable "cloud_config_script_path" {
  default     = "scripts/add-ssh-web-app.yaml"
  type        = string
  description = "Path to cloud-config file"
} 