// Copyright © 2021 Tom Straub <github.com/straubt1>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-slug"
	"github.com/hashicorp/go-tfe"
)

func DownloadModule(token string, tfeHostname string, tfeOrganization string, moduleName string,
	providerName string, moduleVersion string, directory string) (string, error) {

	tfeClient, ctx := getClientContext()
	pmr, err := tfeClient.RegistryModules.Read(ctx, tfe.RegistryModuleID{
		Organization: tfeOrganization,
		Name:         moduleName,
		Provider:     providerName,
	})
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

	downloadURL := resp.Header["X-Terraform-Get"][0]
	if downloadURL == "" {
		return "", errors.New("did not get a download Link")
	}

	client2 := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = downloadURL
			return nil
		},
	}
	// Put content on file
	resp, err = client2.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := slug.Unpack(resp.Body, directory); err != nil {
		return "", err
	}

	return directory, nil
}
