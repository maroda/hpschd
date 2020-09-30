resource "aws_ecs_cluster" "ecs-cluster" {
  name = "hpschd_${var.epithet}"
}

resource "aws_ecs_task_definition" "ecs-task" {
  family                   = var.epithet
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs-task-exec.arn
  container_definitions    = <<MESOTASK
[
  {
  "image": "${var.repository}:${var.app}_${var.release}",
  "name": "${var.epithet}",
  "essential": true,
        "logConfiguration": {
            "logDriver": "awslogs",
            "options": {
                "awslogs-group": "${var.app}_${var.release}",
                "awslogs-region": "${var.aws_region}",
                "awslogs-stream-prefix": "ecs"
              }
          },
      "portMappings": [
        {
          "hostPort": 9999,
          "protocol": "tcp",
          "containerPort": 9999
        }
      ]
    }
]
MESOTASK
}

resource "aws_ecs_service" "ecs-service" {
  name             = "hpschd_${var.epithet}"
  cluster          = aws_ecs_cluster.ecs-cluster.arn
  task_definition  = aws_ecs_task_definition.ecs-task.arn
  launch_type      = "FARGATE"
  platform_version = "LATEST"

  desired_count = var.tcount

  lifecycle {
    create_before_destroy = true
    ignore_changes        = [desired_count]
  }

  # when set to true will redeploy new containers, useful for version upgrades
  force_new_deployment = true

  # required for type awsvpc
  network_configuration {
    security_groups  = [aws_security_group.hpschd-task.id]
    subnets          = aws_subnet.hpschd.*.id
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.hpschd.arn
    container_name   = var.epithet
    container_port   = 9999
  }

  depends_on = [aws_lb_listener.hpschd]
}
