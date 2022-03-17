package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type Client struct {
	region   string
	authToken  string
	ngx_username string
	ngx_password string
	httpClient *http.Client
}


type Project struct {
	Project_name string `json:"project_name"`
	Stack_name string `json:"stack_name"` 
	Project_type string `json:"project_type"`
	Instance_type string `json:"instance_type"`
	Status string `json:"status"`
	Email string `json:"email"`
	Region string `json:"region"`

}

func NewClient(region string, token string, nginx_username string, nginx_password string) *Client {
	return &Client{
		region:       region,
		authToken:  token,
		ngx_username: nginx_username,
		ngx_password: nginx_password,
		httpClient: &http.Client{},
	}
}

// GetAll Retrieves all of the Items from the server
// func (c *Client) GetAll() (*map[string]Project, error) {
// 	body, err := c.httpRequest("item", "GET", bytes.Buffer{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	items := map[string]Project{}
// 	err = json.NewDecoder(body).Decode(&items)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &items, nil
// }

func (c *Client) GetProject(name string) (*Project, error) {
	body, err := c.httpRequest(fmt.Sprintf("/instance/%v", name), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	project := &Project{}
	err = json.NewDecoder(body).Decode(project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (c *Client) AddProject(project *Project) error {
	buf := bytes.Buffer{}
	project.Region = c.region
	err := json.NewEncoder(&buf).Encode(project)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("/provision/%s", project.Project_type), "POST", buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateProject(project *Project) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(project)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("/instance/%s", project.Project_name), "PATCH", buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteProject(itemName string) error {
	_, err := c.httpRequest(fmt.Sprintf("/instance/%s", itemName), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) httpRequest(path, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)

	req.Header.Set("X-User-Token", c.authToken)
	req.SetBasicAuth(c.ngx_username,c.ngx_password)


	if err != nil {
		return nil, err
	}
	switch method {
	case "GET":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	hostname := "https://cloud-api.comwork.io/v1"
	return fmt.Sprintf("%s%s", hostname, path)
}
