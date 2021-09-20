package aws

import (
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccAWSDefaultVpcDhcpOptions_basic(t *testing.T) {
	var d ec2.DhcpOptions
	resourceName := "aws_default_vpc_dhcp_options.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSDefaultVpcDhcpOptionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDefaultVpcDhcpOptionsConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDHCPOptionsExists(resourceName, &d),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`dhcp-options/dopt-.+`)),
					resource.TestCheckResourceAttr(resourceName, "domain_name", regionalPrivateDNSSuffix(acctest.Region())),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers", "AmazonProvidedDNS"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Default DHCP Option Set"),
					acctest.CheckResourceAttrAccountID(resourceName, "owner_id"),
				),
			},
		},
	})
}

func TestAccAWSDefaultVpcDhcpOptions_owner(t *testing.T) {
	var d ec2.DhcpOptions
	resourceName := "aws_default_vpc_dhcp_options.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSDefaultVpcDhcpOptionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDefaultVpcDhcpOptionsConfigOwner,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDHCPOptionsExists(resourceName, &d),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`dhcp-options/dopt-.+`)),
					resource.TestCheckResourceAttr(resourceName, "domain_name", regionalPrivateDNSSuffix(acctest.Region())),
					resource.TestCheckResourceAttr(resourceName, "domain_name_servers", "AmazonProvidedDNS"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "Default DHCP Option Set"),
					acctest.CheckResourceAttrAccountID(resourceName, "owner_id"),
				),
			},
		},
	})
}

func testAccCheckAWSDefaultVpcDhcpOptionsDestroy(s *terraform.State) error {
	// We expect DHCP Options Set to still exist
	return nil
}

const testAccAWSDefaultVpcDhcpOptionsConfigBasic = `
resource "aws_default_vpc_dhcp_options" "test" {
  tags = {
    Name = "Default DHCP Option Set"
  }
}
`

const testAccAWSDefaultVpcDhcpOptionsConfigOwner = `
data "aws_caller_identity" "current" {}

resource "aws_default_vpc_dhcp_options" "test" {
  owner_id = data.aws_caller_identity.current.account_id

  tags = {
    Name = "Default DHCP Option Set"
  }
}
`
