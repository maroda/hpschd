/*

  VPC Networking

*/

data "aws_availability_zones" "aws-az" {
  state = "available"
}

# VPC
#
resource "aws_vpc" "hpschd" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true

  tags = {
    Name = "ecs-hpschd-${var.epithet}"
  }
}

# IGW
#
resource "aws_internet_gateway" "hpschd" {
  vpc_id = aws_vpc.hpschd.id

  tags = {
    Name = "ecs-hpschd-${var.epithet}"
  }
}

# Subnet
#
resource aws_subnet "hpschd" {
  count  = length(data.aws_availability_zones.aws-az.names)
  vpc_id = aws_vpc.hpschd.id

  # cidr_block              = var.subnet_cidr
  cidr_block              = cidrsubnet(aws_vpc.hpschd.cidr_block, 8, count.index + 1)
  availability_zone       = data.aws_availability_zones.aws-az.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "ecs-hpschd-${var.epithet}"
  }
}

# Subnet Routing
#
resource "aws_route_table" "hpschd" {
  vpc_id = aws_vpc.hpschd.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.hpschd.id
  }

  tags = {
    Name = "ecs-hpschd-${var.epithet}"
  }
}

resource "aws_main_route_table_association" "hpschd" {
  vpc_id         = aws_vpc.hpschd.id
  route_table_id = aws_route_table.hpschd.id
}
