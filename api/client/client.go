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


type Instance struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Environment string `json:"environment"`
	Instance_type string `json:"instance_type"`
	Status string `json:"status"`
	Project string `json:"project_id"`
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
	respBody, err := c.httpRequest(fmt.Sprintf("instance/%s/provision/%s",c.region, instance.Environment), "POST", buf)
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
	_, err = c.httpRequest(fmt.Sprintf("/instance/%s/%s", c.region,instance.Id), "PATCH", buf)
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
