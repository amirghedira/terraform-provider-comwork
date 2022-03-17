package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/comwork/comwork-provider/api/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func validateInstanceType(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected instance type to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	instanceAllowedTypes := map[string]bool {
		"DEV1-S": true,
		"DEV1-M": true,
		"DEV1-L": true,
		"DEV1-XL": true,
	}
	if !instanceAllowedTypes[value]{
		errs = append(errs, fmt.Errorf("no instance type with that name. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_name": {
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
			"project_type": {
				Type:        schema.TypeString,
				Required:     true,
				Description: "Type of the project",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the instance",
				ValidateFunc: validateInstanceType,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "status of the instance (poweroff,poweron)",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email attached to this resource",
			},
		},
		Create: instanceCreateItem,
		Read:   instanceReadItem,
		Update: instanceUpdateItem,
		Delete: instanceDeleteItem,
		Exists: instanceExistsItem,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func instanceCreateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	project := client.Project{
		Project_name: d.Get("project_name").(string),
		Stack_name: d.Get("stack_name").(string),
		Project_type: d.Get("project_type").(string),
		Instance_type: d.Get("instance_type").(string),
		Status: d.Get("status").(string),
		Email: d.Get("email").(string),
	}
	err := apiClient.AddProject(&project)

	if err != nil {
		return err
	}
	d.SetId(project.Project_name)
	return nil
}

func instanceReadItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	itemId := d.Id()
	item, err := apiClient.GetProject(itemId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Item with ID %s", itemId)
		}
	}

	d.SetId(item.Project_name)
	d.Set("project_name", item.Project_name)
	d.Set("instance_type", item.Instance_type)
	d.Set("project_type", item.Project_type)
	d.Set("status", item.Status)
	d.Set("email", item.Email)

	return nil
}

func instanceUpdateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	project := client.Project{
		Project_name: d.Get("project_name").(string),
		Stack_name: d.Get("project_name").(string),
		Project_type: d.Get("project_type").(string),
		Instance_type: d.Get("instance_type").(string),
		Status: d.Get("status").(string),
		Email: d.Get("email").(string),
	}

	err := apiClient.UpdateProject(&project)
	if err != nil {
		return err
	}
	return nil
}

func instanceDeleteItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	projectId := d.Id()

	err := apiClient.DeleteProject(projectId)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func instanceExistsItem(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	projectId := d.Id()
	_, err := apiClient.GetProject(projectId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
