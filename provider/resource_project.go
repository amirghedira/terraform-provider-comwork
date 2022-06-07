package provider

import (
	"fmt"
	"strings"

	"strconv"

	"github.com/comwork/comwork-provider/api/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the project",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"url": {
				Type:         schema.TypeString,
				Computed: true,
				Description:  "The url of the project",
			},
			"created_at": {
				Type:         schema.TypeString,
				Computed: true,
				Description:  "The creation date of the project",
			},
		},
		Create: projectCreateItem,
		Read:   projectReadItem,
		Delete: projectDeleteItem,
		Exists: projectExistsItem,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func projectCreateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	project := client.Project{
		Name: d.Get("name").(string),
	}
	created_project ,err := apiClient.AddProject(&project)

	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(created_project.Id))
	d.Set("url", created_project.Url)
	d.Set("created_at", created_project.CreatedAt)
	return nil
}

func projectReadItem(d *schema.ResourceData, m interface{}) error {
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

	d.SetId(strconv.Itoa(item.Id))

	d.Set("name", item.Name)
	d.Set("url", item.Url)
	d.Set("created_at", item.CreatedAt)

	return nil
}


func projectDeleteItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	projectId := d.Id()

	err := apiClient.DeleteProject(projectId)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func projectExistsItem(d *schema.ResourceData, m interface{}) (bool, error) {
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
