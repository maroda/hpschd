/*

    Security Groups

*/

# hpschd PUBLIC
#
resource "aws_security_group" "hpschd" {
  vpc_id = aws_vpc.hpschd.id
  name   = "hpschd-${var.epithet}"

  # Mesostic API Access
  ingress {
    protocol    = "tcp"
    from_port   = 9999
    to_port     = 9999
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "hpschd-${var.epithet}"
  }
}

# hpschd task PRIVATE
#
resource "aws_security_group" "hpschd-task" {
  vpc_id = aws_vpc.hpschd.id
  name   = "hpschd-task-${var.epithet}"

  # Mesostic API Access
  ingress {
    protocol        = "tcp"
    from_port       = 9999
    to_port         = 9999
    security_groups = [aws_security_group.hpschd.id]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "hpschd-task-${var.epithet}"
  }
}
