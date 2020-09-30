resource "aws_cloudwatch_log_group" "hpschd-release" {
  name              = "${var.app}_${var.release}"
  retention_in_days = 0

  tags = {
    Name = "${var.app}_${var.release}"
  }
}
