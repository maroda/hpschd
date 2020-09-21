/*

    IAM

*/

# ECS Access
#
data "aws_iam_policy_document" "ecs-policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs-task-exec" {
  name               = "ecs-task-exec-${var.epithet}"
  assume_role_policy = data.aws_iam_policy_document.ecs-policy.json

  tags = {
    Name = "ecs-task-exec-${var.epithet}"
  }
}

resource "aws_iam_role" "ecs-task" {
  name               = "ecs-task-${var.epithet}"
  assume_role_policy = data.aws_iam_policy_document.ecs-policy.json

  tags = {
    Name = "ecs-task-${var.epithet}"
  }
}

resource "aws_iam_policy_attachment" "ecs-task-exec" {
  name       = "ecs-task-exec-${var.epithet}"
  roles      = [aws_iam_role.ecs-task-exec.name]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_policy_attachment" "ecs-task" {
  name       = "ecs-task-${var.epithet}"
  roles      = [aws_iam_role.ecs-task.name]
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}
