/*

  Load Balancers

*/

#
# hpschd.xyz
#
resource "aws_lb" "hpschd" {
  name               = "hpschd-${var.epithet}"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.hpschd.id]
  subnets            = aws_subnet.hpschd.*.id

  tags = {
    Name = "hpschd-${var.epithet}"
  }
}

resource "aws_lb_listener" "hpschd" {
  load_balancer_arn = aws_lb.hpschd.arn
  port              = 9999
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.hpschd.arn
  }
}

resource "aws_lb_target_group" "hpschd" {
  name        = "hpschd-${var.epithet}"
  port        = 9999
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.hpschd.id

  health_check {
    healthy_threshold   = "3"
    interval            = "30"
    protocol            = "HTTP"
    matcher             = "200"
    timeout             = "3"
    unhealthy_threshold = "2"
    path                = "/ping"
    port                = 9999
  }
}
