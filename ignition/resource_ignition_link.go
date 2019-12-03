package ignition

import (
	"github.com/coreos/ignition/config/v2_1/types"
	// "github.com/coreos/ignition/config/v2_1/types/path"
	"github.com/coreos/ignition/config/shared/errors"
	"github.com/coreos/ignition/config/validate/report"
	"github.com/hashicorp/terraform/helper/schema"
	"path"
)

func dataSourceLink() *schema.Resource {
	return &schema.Resource{
		Exists: resourceLinkExists,
		Read:   resourceLinkRead,
		Schema: map[string]*schema.Schema{
			"filesystem": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hard": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLinkRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildLink(d, globalCache)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceLinkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildLink(d, globalCache)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildLink(d *schema.ResourceData, c *cache) (string, error) {
	link := &types.Link{}
	link.Filesystem = d.Get("filesystem").(string)
	link.Path = d.Get("path").(string)
	link.Target = d.Get("target").(string)
	link.Hard = d.Get("hard").(bool)

	uid := d.Get("uid").(int)
	if uid != 0 {
		link.User = types.NodeUser{ID: &uid}
	}

	gid := d.Get("gid").(int)
	if gid != 0 {
		link.Group = types.NodeGroup{ID: &gid}
	}

	return c.addLink(link), handleReport(ValidateTarget(link))
}

func ValidateTarget(s *types.Link) report.Report {
	r := report.Report{}
	if !s.Hard {
		err := validatePath(s.Target)
		if err != nil {
			r.Add(report.Entry{
				Message: err.Error(),
				Kind:    report.EntryError,
			})
		}
	}
	return r
}

func validatePath(p string) error {
	if !path.IsAbs(p) {
		return errors.ErrPathRelative
	}
	return nil
}
