package provider

import (
	"fmt"
	"github.com/amirghedira/terraform-provider/api/client"
	"regexp"
	"strings"
	"github.com/hashicorp/terraform/helper/schema"
)

func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the project",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"stack_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the stack",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the environement",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the instance",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email attached to this resource",
			},
			// "tags": {
			// 	Type:        schema.TypeSet,
			// 	Optional:    true,
			// 	Description: "An optional list of tags, represented as a key, value pair",
			// 	Elem:        &schema.Schema{Type: schema.TypeString},
			// },
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceUpdateItem,
		Delete: resourceDeleteItem,
		Exists: resourceExistsItem,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	item := client.Project{
		project_name: d.Get("name").(string),
		stack_name: d.Get("stack_name").(string),
		project_type: d.Get("project_type").(string),
		instance_type: d.Get("instance_type").(string),
		status: d.Get("status").(string),
		email: d.Get("email").(string),
	}
	err := apiClient.NewItem(&item)

	if err != nil {
		return err
	}
	d.SetId(item.project_name)
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	itemId := d.Id()
	item, err := apiClient.GetItem(itemId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Item with ID %s", itemId)
		}
	}

	d.SetId(item.project_name)
	d.Set("project_name", item.project_name)
	d.Set("instance_type", item.instance_type)
	d.Set("project_type", item.project_type)
	d.Set("status", item.status)
	d.Set("email", item.email)

	return nil
}

func resourceUpdateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	item := client.Project{
		project_name: d.Get("project_name").(string),
		stack_name: d.Get("project_name").(string),
		project_type: d.Get("project_type").(string),
		instance_type: d.Get("instance_type").(string),
		status: d.Get("status").(string),
		email: d.Get("email").(string),
	}

	err := apiClient.UpdateItem(&item)
	if err != nil {
		return err
	}
	return nil
}

func resourceDeleteItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	itemId := d.Id()

	err := apiClient.DeleteItem(itemId)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceExistsItem(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	itemId := d.Id()
	_, err := apiClient.GetItem(itemId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
