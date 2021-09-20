package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/eks/token"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccAWSEksClusterAuthDataSource_basic(t *testing.T) {
	dataSourceResourceName := "data.aws_eks_cluster_auth.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, eks.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsEksClusterAuthConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceResourceName, "name", "foobar"),
					resource.TestCheckResourceAttrSet(dataSourceResourceName, "token"),
					testAccCheckAwsEksClusterAuthToken(dataSourceResourceName),
				),
			},
		},
	})
}

func testAccCheckAwsEksClusterAuthToken(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource ID is set")
		}

		name := rs.Primary.Attributes["name"]
		tok := rs.Primary.Attributes["token"]
		verifier := token.NewVerifier(name)
		identity, err := verifier.Verify(tok)
		if err != nil {
			return fmt.Errorf("Error verifying token for cluster %q: %v", name, err)
		}
		if identity.ARN == "" {
			return fmt.Errorf("Unexpected blank ARN for token identity")
		}

		return nil
	}
}

const testAccCheckAwsEksClusterAuthConfig_basic = `
data "aws_eks_cluster_auth" "test" {
  name = "foobar"
}
`
