variable "zoneapex" {
  default = "hpschd.xyz"
}

variable "aws_region" {
  default = "us-west-2"
}

variable "epithet" {
  description = "ECS Task name for HPSCHD"
  default     = "mesostic-api"
}

variable "vpc_cidr" {
  default = "10.0.0.0/16"
}

variable "subnet_cidr" {
  default = "10.0.1.0/24"
}

variable "app" {
  default = "hpschd"
}

variable "release" {
  default = "v1.3.0"
}

variable "repository" {
  default = "docker.io/maroda/chaquo"
}

variable "tcount" {
  default = 1
}
