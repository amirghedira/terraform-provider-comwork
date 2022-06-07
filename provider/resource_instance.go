package provider

import (
	"fmt"
	"regexp"
	"strings"

	"strconv"

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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the instance",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"environment": {
				Type:        schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description: "Type of the project",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:     true,
				Description: "Type of the instance",
				ValidateFunc: validateInstanceType,
			},
			"attach": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default: false,
				Description:  "Whether the instance will be attached to a project or created from scratch",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "status of the instance (poweroff,poweron)",
			},
			
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:     true,
				Description: "project attached to this resource",
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
	instance := client.Instance{
		Name: d.Get("name").(string),
		Environment: d.Get("environment").(string),
		Instance_type: d.Get("instance_type").(string),
		Status: d.Get("status").(string),
		Project: d.Get("project_id").(int),
		Attach: d.Get("attach").(bool),
	}
	if instance.Attach{
		created_instance ,err := apiClient.AttachInstance(&instance)
		if err != nil {
			return err
		}
		d.SetId(strconv.Itoa(created_instance.Id))
		return nil
	} else {
		created_instance ,err := apiClient.AddInstance(&instance)
		if err != nil {
			return err
		}
		d.SetId(strconv.Itoa(created_instance.Id))
		return nil
	}
}

func instanceReadItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	itemId := d.Id()
	item, err := apiClient.GetInstance(itemId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Item with ID %s", itemId)
		}
	}

	d.SetId(strconv.Itoa(item.Id))

	d.Set("name", item.Name)
	d.Set("instance_type", item.Instance_type)
	d.Set("environment", item.Environment)
	d.Set("status", item.Status)
	d.Set("status", item.Status)
	d.Set("project_id", item.Project)

	return nil
}

func instanceUpdateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	instance_id ,_ :=strconv.Atoi(d.Id())
	instance := client.Instance{
		Id: instance_id,
		Name: d.Get("name").(string),
		Environment: d.Get("environment").(string),
		Instance_type: d.Get("instance_type").(string),
		Status: d.Get("status").(string),
		Project: d.Get("project_id").(int),
	}

	err := apiClient.UpdateInstance(&instance)
	if err != nil {
		return err
	}
	return nil
}

func instanceDeleteItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	instanceId := d.Id()

	err := apiClient.DeleteInstance(instanceId)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func instanceExistsItem(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	instanceId := d.Id()
	_, err := apiClient.GetInstance(instanceId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
