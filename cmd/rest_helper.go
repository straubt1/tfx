/*
Copyright © 2021 Tom Straub <tstraub@hashicorp.com>

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
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	// Convert response body to Todo struct
	var rmv RegistryModuleVersion
	err = json.Unmarshal(bodyBytes, &rmv)
	if err != nil {
		return nil, err
	}
	fmt.Println("Create Module Version:", version)
	fmt.Println("Upload Link:", rmv.Data.Links.Upload)

	return &rmv.Data.Links.Upload, nil
}

func RegistryModulesUpload(token string, url string) error {
	var err error

	// create http Client to make calls
	client := &http.Client{}

	file, err := os.Open("/Users/tstraub/tfx/module/module.tar.gz")
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("filetype", filepath.Base(file.Name()))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	// create request
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}

	// add headers
	req.Header.Set("Authorization", "Bearer "+token)
	// req.Header.Set("Accept", "application/vnd.api+json")
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

	// // read all bytes, convert to object
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// // Convert response body to Todo struct
	// var rmv RegistryModuleVersion
	// json.Unmarshal(bodyBytes, &rmv)
	// fmt.Println("Create Module Version:", version)
	// fmt.Println("Upload Link:", rmv.Data.Links.Upload)

	// return &rmv.Data.Links.Upload, nil
}
