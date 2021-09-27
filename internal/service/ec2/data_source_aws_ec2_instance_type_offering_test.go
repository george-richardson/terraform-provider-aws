package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccAWSEc2InstanceTypeOfferingDataSource_Filter(t *testing.T) {
	dataSourceName := "data.aws_ec2_instance_type_offering.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEc2InstanceTypeOffering(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEc2InstanceTypeOfferingDataSourceConfigFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_type"),
				),
			},
		},
	})
}

func TestAccAWSEc2InstanceTypeOfferingDataSource_LocationType(t *testing.T) {
	dataSourceName := "data.aws_ec2_instance_type_offering.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEc2InstanceTypeOffering(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEc2InstanceTypeOfferingDataSourceConfigLocationType(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_type"),
				),
			},
		},
	})
}

func TestAccAWSEc2InstanceTypeOfferingDataSource_PreferredInstanceTypes(t *testing.T) {
	dataSourceName := "data.aws_ec2_instance_type_offering.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSEc2InstanceTypeOffering(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEc2InstanceTypeOfferingDataSourceConfigPreferredInstanceTypes(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "instance_type", "t3.micro"),
				),
			},
		},
	})
}

func testAccPreCheckAWSEc2InstanceTypeOffering(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	input := &ec2.DescribeInstanceTypeOfferingsInput{
		MaxResults: aws.Int64(5),
	}

	_, err := conn.DescribeInstanceTypeOfferings(input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccAWSEc2InstanceTypeOfferingDataSourceConfigFilter() string {
	return `
# Rather than hardcode an instance type in the testing,
# use the first result from all available offerings.
data "aws_ec2_instance_type_offerings" "test" {}

data "aws_ec2_instance_type_offering" "test" {
  filter {
    name   = "instance-type"
    values = [tolist(data.aws_ec2_instance_type_offerings.test.instance_types)[0]]
  }
}
`
}

func testAccAWSEc2InstanceTypeOfferingDataSourceConfigLocationType() string {
	return acctest.ConfigAvailableAZsNoOptIn() + `
# Rather than hardcode an instance type in the testing,
# use the first result from all available offerings.
data "aws_ec2_instance_type_offerings" "test" {
  filter {
    name   = "location"
    values = [data.aws_availability_zones.available.names[0]]
  }

  location_type = "availability-zone"
}

data "aws_ec2_instance_type_offering" "test" {
  filter {
    name   = "instance-type"
    values = [tolist(data.aws_ec2_instance_type_offerings.test.instance_types)[0]]
  }

  filter {
    name   = "location"
    values = [data.aws_availability_zones.available.names[0]]
  }

  location_type = "availability-zone"
}
`
}

func testAccAWSEc2InstanceTypeOfferingDataSourceConfigPreferredInstanceTypes() string {
	return `
data "aws_ec2_instance_type_offering" "test" {
  filter {
    name   = "instance-type"
    values = ["t1.micro", "t2.micro", "t3.micro"]
  }

  preferred_instance_types = ["t3.micro", "t2.micro", "t1.micro"]
}
`
}