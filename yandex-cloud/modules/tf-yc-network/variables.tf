variable "network_zones" {
  default     = ["ru-central1-a", "ru-central1-b", "ru-central1-d"]
  type        = set(string)
  description = "Network zones"
} 