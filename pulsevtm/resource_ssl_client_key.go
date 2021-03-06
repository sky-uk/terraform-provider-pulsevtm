package pulsevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/terraform-provider-pulsevtm/pulsevtm/util"
)

func resourceSSLClientKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSLClientKeyCreate,
		Read:   resourceSSLClientKeyRead,
		Update: resourceSSLClientKeyUpdate,
		Delete: resourceSSLClientKeyDelete,

		Schema: util.SchemaSSLKey(),
	}
}

func resourceSSLClientKeyCreate(d *schema.ResourceData, m interface{}) error {
	err := util.SSLKeyCreate(d, m, "ssl/client_keys")
	if err != nil {
		return err
	}
	return resourceSSLClientKeyRead(d, m)
}

func resourceSSLClientKeyRead(d *schema.ResourceData, m interface{}) error {
	err := util.SSLKeyRead(d, m, "ssl/client_keys")
	if err != nil {
		return err
	}
	return nil
}

func resourceSSLClientKeyUpdate(d *schema.ResourceData, m interface{}) error {
	err := util.SSLKeyUpdate(d, m, "ssl/client_keys")
	if err != nil {
		return err
	}
	return resourceSSLClientKeyRead(d, m)
}

func resourceSSLClientKeyDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("ssl/client_keys", d, m)
}
