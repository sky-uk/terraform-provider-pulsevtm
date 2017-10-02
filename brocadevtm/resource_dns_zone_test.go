package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/go-brocade-vtm/api/dns_zone"
	"github.com/sky-uk/go-rest-api"
	"regexp"
	"testing"
)

func TestAccBrocadeVTMDNSZoneBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	dnsZoneName := fmt.Sprintf("acctest_brocadevtm_dns_zone-%d", randomInt)
	dnsZoneResourceName := "brocadevtm_dns_zone.acctest"
	fmt.Printf("\n\nDNS zone is %s.\n\n", dnsZoneName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccBrocadeVTMDNSZoneCheckDestroy(state, dnsZoneName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccBrocadeVTMDNSZoneNoNameTemplate(),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccBrocadeDNSZoneCreateTemplate(dnsZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMDNSZoneExists(dnsZoneName, dnsZoneResourceName),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "name", dnsZoneName),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "origin", "example.com"),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "zone_file", "example.com.db"),
				),
			},
			{
				Config: testAccBrocadeDNSZoneUpdateTemplate(dnsZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccBrocadeVTMDNSZoneExists(dnsZoneName, dnsZoneResourceName),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "name", dnsZoneName),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "origin", "updated-example.com"),
					resource.TestCheckResourceAttr(dnsZoneResourceName, "zone_file", "updated-example.com.db"),
				),
			},
		},
	})
}

func testAccBrocadeVTMDNSZoneCheckDestroy(state *terraform.State, name string) error {

	vtmClient := testAccProvider.Meta().(*rest.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "brocadevtm_dns_zone" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		api := dnsZone.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Brocade vTM DNS zone - error occurred whilst retrieving a list of all DNS zones")
		}
		for _, dnsZone := range api.ResponseObject().(*dnsZone.DNSZones).Children {
			if dnsZone.Name == name {
				return fmt.Errorf("Brocade vTM DNS zone %s still exists", name)
			}
		}
	}
	return nil
}

func testAccBrocadeVTMDNSZoneExists(dnsZoneName, dnsZoneResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[dnsZoneResourceName]
		if !ok {
			return fmt.Errorf("\nBrocade vTM DNS zone %s wasn't found in resources", dnsZoneName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("\nBrocade vTM DNS zone ID not set for %s in resources", dnsZoneName)
		}

		vtmClient := testAccProvider.Meta().(*rest.Client)
		api := dnsZone.NewGetAll()
		err := vtmClient.Do(api)
		if err != nil {
			return fmt.Errorf("Error: %+v", err)
		}
		for _, dnsZone := range api.ResponseObject().(*dnsZone.DNSZones).Children {
			if dnsZone.Name == dnsZoneName {
				return nil
			}
		}
		return fmt.Errorf("Brocade vTM DNS zone %s not found on remote vTM", dnsZoneName)
	}
}

func testAccBrocadeVTMDNSZoneNoNameTemplate() string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone" "acctest" {
  origin = "example.com"
  zone_file = "example.com.db"
}
`)
}

func testAccBrocadeDNSZoneCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone" "acctest" {
  name = "%s"
  origin = "example.com"
  zone_file = "example.com.db"
}
`, name)
}

func testAccBrocadeDNSZoneUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "brocadevtm_dns_zone" "acctest" {
  name = "%s"
  origin = "updated-example.com"
  zone_file = "updated-example.com.db"
}
`, name)
}
