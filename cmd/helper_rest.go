// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/go-slug"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/viper"
	"github.com/straubt1/tfx/client"
)

func DownloadModule(c *client.TfxClient, moduleName string,
	providerName string, moduleVersion string, directory string) (string, error) {

	pmr, err := c.Client.RegistryModules.Read(c.Context, tfe.RegistryModuleID{
		Organization: c.OrganizationName,
		Name:         moduleName,
		Provider:     providerName,
		Namespace:    c.OrganizationName,
		RegistryName: tfe.PrivateRegistry,
	})
	if err != nil || pmr == nil {
		return "", errors.New("can't find module")
	}

	// create url
	url := fmt.Sprintf(
		"https://%s/api/registry/v1/modules/%s/%s/%s/%s/download",
		c.Hostname,
		c.OrganizationName,
		moduleName,
		providerName,
		moduleVersion,
	)
	// create http Client to make calls
	httpClient := &http.Client{}

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// add headers
	req.Header.Set("Authorization", "Bearer "+viper.GetString("tfeToken"))
	req.Header.Set("Accept", "application/vnd.api+json")

	// make request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	// wait for complete
	defer resp.Body.Close()

	downloadURL := resp.Header["X-Terraform-Get"][0]
	if downloadURL == "" {
		return "", errors.New("did not get a download Link")
	}

	httpClient2 := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = downloadURL
			return nil
		},
	}
	// Put content on file
	resp, err = httpClient2.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := slug.Unpack(resp.Body, directory); err != nil {
		return "", err
	}

	return directory, nil
}

func UploadBinary(uploadURL string, path string) error {
	data, err := os.Open(path)
	if err != nil {
		return err
	}
	defer data.Close()
	req, err := http.NewRequest("PUT", uploadURL, data)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func DownloadTextFile(downloadURL string) (string, error) {

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(downloadURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}

	return string(b), nil
}
