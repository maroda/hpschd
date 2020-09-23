variable "zoneapex" {
  default = "hpschd.xyz"
}

variable "aws_region" {
  default = "us-west-1"
}

variable "endpoint_ip" {
  default = "65.50.222.218"
}

variable "ecs" {
  default = "EC2Co-EcsEl-1HU46LZBDSTRC-1958930581.us-west-2.elb.amazonaws.com"
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
