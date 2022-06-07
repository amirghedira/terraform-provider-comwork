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
	httpClient *http.Client
}

type ErrorResponse struct {
	Error string `json:"error"`
}


type Project struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
	Region string `json:"region"`
	CreatedAt string `json:"created_at"`

}



type Instance struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Environment string `json:"environment"`
	Instance_type string `json:"instance_type"`
	Status string `json:"status"`
	Project int `json:"project_id"`
	Region string `json:"region"`

}

func NewClient(region string, token string) *Client {
	return &Client{
		region:       region,
		authToken:  token,
		httpClient: &http.Client{},
	}
}


func (c *Client) GetInstance(instanceId string) (*Instance, error) {
	body, err := c.httpRequest(fmt.Sprintf("/instance/%s/%v",c.region, instanceId), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	instance := &Instance{}
	err = json.NewDecoder(body).Decode(instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *Client) AddInstance(instance *Instance) (*Instance, error) {
	buf := bytes.Buffer{}
	instance.Region = c.region
	err := json.NewEncoder(&buf).Encode(instance)
	if err != nil {
		return nil, err
	}
	respBody, err := c.httpRequest(fmt.Sprintf("/instance/%s/provision/%s",c.region, instance.Environment), "POST", buf)
	if err != nil {
		return nil, err
	}
	created_instance := &Instance{}
	err = json.NewDecoder(respBody).Decode(created_instance)
	if err != nil {
		return nil, err
	}
	return created_instance, nil
}



func (c *Client) UpdateInstance(instance *Instance) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(instance)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("/instance/%s/%v", c.region,instance.Id), "PATCH", buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteInstance(instanceId string) error {
	_, err := c.httpRequest(fmt.Sprintf("/instance/%s/%s", c.region,instanceId), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}


func (c *Client) GetProject(projectId string) (*Project, error) {
	body, err := c.httpRequest(fmt.Sprintf("/project/%s/%s",c.region, projectId), "GET", bytes.Buffer{})
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

func (c *Client) AddProject(project *Project) (*Project, error) {
	buf := bytes.Buffer{}
	project.Region = c.region
	err := json.NewEncoder(&buf).Encode(project)
	if err != nil {
		return nil, err
	}
	respBody, err := c.httpRequest(fmt.Sprintf("/project/%s",c.region), "POST", buf)
	if err != nil {
		return nil, err
	}
	created_project := &Project{}
	err = json.NewDecoder(respBody).Decode(created_project)
	if err != nil {
		return nil, err
	}
	return created_project, nil
}

func (c *Client) DeleteProject(projectId string) error {
	_, err := c.httpRequest(fmt.Sprintf("/project/%s/%s", c.region,projectId), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}


func (c *Client) httpRequest(path, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)

	req.Header.Set("X-User-Token", c.authToken)


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

	if resp.StatusCode != http.StatusOK &&resp.StatusCode != http.StatusCreated {
		errorBody := &ErrorResponse{}
		json.NewDecoder(resp.Body).Decode(errorBody)
		return nil, fmt.Errorf("%s", errorBody.Error)
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	hostname := "https://cloud-api.comwork.io/v1"
	return fmt.Sprintf("%s%s", hostname, path)
}
