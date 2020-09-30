variable "zoneapex" {
  default = "hpschd.xyz"
}

variable "ecs" {
  description = "The ECS ALB endpoint for DNS CNAME"
  default     = "hpschd-mesostic-api-1116213348.us-west-2.elb.amazonaws.com"
}
