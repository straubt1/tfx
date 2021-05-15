/*
Copyright © 2021 Tom Straub <github.com/straubt1>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-slug"
)

type RegistryModuleCreateVersionOptions struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Version string `json:"version"`
		} `json:"attributes"`
	} `json:"data"`
}

type RegistryModuleVersion struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Source    string    `json:"source"`
			Status    string    `json:"status"`
			Version   string    `json:"version"`
			CreatedAt time.Time `json:"created-at"`
			UpdatedAt time.Time `json:"updated-at"`
		} `json:"attributes"`
		// Relationships struct {
		// 	RegistryModule struct {
		// 		Data struct {
		// 			ID   string `json:"id"`
		// 			Type string `json:"type"`
		// 		} `json:"data"`
		// 	} `json:"registry-module"`
		// } `json:"relationships"`
		Links struct {
			Upload string `json:"upload"`
		} `json:"links"`
	} `json:"data"`
}

func RegistryModulesCreateVersion(token string, tfeHostname string, tfeOrganization string,
	moduleName string, providerName string, version string) (*string, error) {
	var err error

	// create url
	url := fmt.Sprintf(
		"https://%s/api/v2/registry-modules/%s/%s/%s/versions",
		tfeHostname,
		tfeOrganization,
		moduleName,
		providerName,
	)

	// create http Client to make calls
	client := &http.Client{}
	postBody := fmt.Sprintf(`{
		"data": {
			"type": "registry-module-versions",
			"attributes": {
				"version": "%s"
			}
		}
	}`, version)

	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// wait for complete
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, errors.New("Non-OK HTTP status:" + string(resp.StatusCode))
	}

	// read all bytes, convert to object
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// bodyString := string(bodyBytes)
	// fmt.Println("API Response as String:\n" + bodyString)

	// Convert response body to Todo struct
	var rmv RegistryModuleVersion
	err = json.Unmarshal(bodyBytes, &rmv)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Create Module Version:", version)
	// fmt.Println("Upload Link:", rmv.Data.Links.Upload)

	return &rmv.Data.Links.Upload, nil
}

func RegistryModulesUpload(token string, url *string, directory string) error {
	var err error

	// TODO: verify directory exists

	// create http Client to make calls
	client := &http.Client{}

	// create body to store tar
	body := bytes.NewBuffer(nil)

	// pack directory into a slug
	if _, err := slug.Pack(directory, body, true); err != nil {
		return err
	}

	// create request
	req, err := http.NewRequest("PUT", *url, body)
	if err != nil {
		return err
	}

	// add headers - no auth needed, its baked into the url
	req.Header.Set("Content-Type", "application/octet-stream")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// wait for complete
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return errors.New("Non-OK HTTP status:" + string(resp.StatusCode))
	}
	return nil
}

type PMRList struct {
	Meta struct {
		Limit         int `json:"limit"`
		CurrentOffset int `json:"current_offset"`
	} `json:"meta"`
	Modules []struct {
		ID          string    `json:"id"`
		Owner       string    `json:"owner"`
		Namespace   string    `json:"namespace"`
		Name        string    `json:"name"`
		Version     string    `json:"version"`
		Provider    string    `json:"provider"`
		Description string    `json:"description"`
		Source      string    `json:"source"`
		Tag         string    `json:"tag"`
		PublishedAt time.Time `json:"published_at"`
		Downloads   int       `json:"downloads"`
		Verified    bool      `json:"verified"`
	} `json:"modules"`
}

func GetAllPMRModules(token string, tfeHostname string, tfeOrganization string) (*PMRList, error) {
	// create url
	url := fmt.Sprintf(
		"https://%s/api/registry/v1/modules/%s",
		tfeHostname,
		tfeOrganization,
	)
	// create http Client to make calls
	client := &http.Client{}

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.api+json")
	// req.Header.Set("Content-Type", "application/vnd.api+json")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// wait for complete
	defer resp.Body.Close()

	// read all bytes, convert to object
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// bodyString := string(bodyBytes)
	// fmt.Println("API Response as String:\n" + bodyString)

	// Convert response body to Todo struct
	var pmr *PMRList
	err = json.Unmarshal(bodyBytes, &pmr)
	if err != nil {
		return nil, err
	}

	return pmr, nil
}

func DownloadModule(token string, tfeHostname string, tfeOrganization string, moduleName string,
	providerName string, moduleVersion string) (string, error) {

	tfeClient, ctx := getClientContext()
	pmr, err := tfeClient.RegistryModules.Read(ctx, tfeOrganization, "moduleName", providerName)
	if err != nil || pmr == nil {
		return "", errors.New("can't find module")
	}

	// create url
	url := fmt.Sprintf(
		"https://%s/api/registry/v1/modules/%s/%s/%s/%s/download",
		tfeHostname,
		tfeOrganization,
		moduleName,
		providerName,
		moduleVersion,
	)
	// create http Client to make calls
	client := &http.Client{}

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// add headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.api+json")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// wait for complete
	defer resp.Body.Close()

	downloadUrl := resp.Header["X-Terraform-Get"][0]
	if downloadUrl == "" {
		return "", errors.New("did not get a download Link")
	}

	client2 := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = downloadUrl
			return nil
		},
	}
	// Put content on file
	resp, err = client2.Get(downloadUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create a directory to unpack the slug contents into.
	dst, err := ioutil.TempDir("", "slug")
	if err != nil {
		return "", err
	}

	if err := slug.Unpack(resp.Body, dst); err != nil {
		return "", err
	}

	return dst, nil
}
