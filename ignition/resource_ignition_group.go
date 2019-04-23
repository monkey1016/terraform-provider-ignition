package ignition

import (
	"github.com/coreos/ignition/config/v2_1/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Exists: resourceGroupExists,
		Read:   resourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"password_hash": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildGroup(d, globalCache)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildGroup(d, globalCache)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildGroup(d *schema.ResourceData, c *cache) (string, error) {
	group := &types.PasswdGroup{
		Name:         d.Get("name").(string),
		PasswordHash: d.Get("password_hash").(string),
		Gid:          getInt(d, "gid"),
	}

	return c.addGroup(group), nil
}
