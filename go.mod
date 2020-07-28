module github.com/whywaita/terraform-provider-lovi

go 1.14

require (
	github.com/hashicorp/terraform-plugin-sdk v1.14.0
	github.com/whywaita/satelit v0.0.0-20200709101056-74c7b1556d9b
	google.golang.org/grpc v1.30.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
