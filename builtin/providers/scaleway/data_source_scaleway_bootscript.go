package scaleway

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

func dataSourceScalewayBootscript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScalewayBootscriptRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"architecture": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// Computed values.
			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"boot_cmd_args": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dtb": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"initrd": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"kernel": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func bootscriptDescriptionAttributes(d *schema.ResourceData, script api.ScalewayBootscript) error {
	d.Set("architecture", script.Arch)
	d.Set("organization", script.Organization)
	d.Set("public", script.Public)
	d.Set("boot_cmd_args", script.Bootcmdargs)
	d.Set("dtb", script.Dtb)
	d.Set("initrd", script.Initrd)
	d.Set("kernel", script.Kernel)
	d.SetId(script.Identifier)

	return nil
}

func dataSourceScalewayBootscriptRead(d *schema.ResourceData, meta interface{}) error {
	scaleway := meta.(*Client).scaleway

	scripts, err := scaleway.GetBootscripts()
	if err != nil {
		return err
	}

	var isMatch func(api.ScalewayBootscript) bool

	if name, ok := d.GetOk("name"); ok {
		isMatch = func(s api.ScalewayBootscript) bool {
			return s.Title == name.(string)
		}
	} else if nameFilter, ok := d.GetOk("name_filter"); ok {
		architecture := d.Get("architecture")
		exp, err := regexp.Compile(nameFilter.(string))
		if err != nil {
			return err
		}

		isMatch = func(s api.ScalewayBootscript) bool {
			nameMatch := exp.MatchString(s.Title)
			architectureMatch := true
			if architecture != "" {
				architectureMatch = architecture == s.Arch
			}
			return nameMatch && architectureMatch
		}
	}

	var matches []api.ScalewayBootscript
	for _, script := range *scripts {
		if isMatch(script) {
			matches = append(matches, script)
		}
	}

	if len(matches) > 1 {
		return fmt.Errorf("The query returned more than one result. Please refine your query.")
	}
	if len(matches) == 0 {
		return fmt.Errorf("The query returned no result. Please refine your query.")
	}

	return bootscriptDescriptionAttributes(d, matches[0])
}
