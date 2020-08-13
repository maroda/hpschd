#
# hosted zone records
#
resource "aws_route53_record" "apex" {
  zone_id = aws_route53_zone.apex.zone_id
  name    = var.zoneapex
  type    = "A"

  alias {
    zone_id                = aws_route53_zone.apex.zone_id
    name                   = aws_route53_record.www.name
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "wildcard" {
  zone_id = aws_route53_zone.apex.zone_id
  name    = "*.${var.zoneapex}"
  type    = "A"
  ttl     = "300"
  records = [var.endpoint_ip]
}

resource "aws_route53_record" "www" {
  zone_id = aws_route53_zone.apex.zone_id
  name    = "www.${var.zoneapex}"
  type    = "A"
  ttl     = "300"
  records = [var.endpoint_ip]
}

