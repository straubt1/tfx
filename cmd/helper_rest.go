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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	b64 "encoding/base64"

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

type GPGList struct {
	Meta struct {
		Pagination struct {
			PageSize    int         `json:"page-size"`
			CurrentPage int         `json:"current-page"`
			NextPage    interface{} `json:"next-page"`
			PrevPage    interface{} `json:"prev-page"`
			TotalPages  int         `json:"total-pages"`
			TotalCount  int         `json:"total-count"`
		} `json:"pagination"`
	} `json:"meta"`
	Keys []struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			ASCIIArmor     string      `json:"ascii-armor"`
			CreatedAt      time.Time   `json:"created-at"`
			KeyID          string      `json:"key-id"`
			Namespace      string      `json:"namespace"`
			Source         string      `json:"source"`
			SourceURL      interface{} `json:"source-url"`
			TrustSignature string      `json:"trust-signature"`
			UpdatedAt      time.Time   `json:"updated-at"`
		} `json:"attributes"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
	// Links struct {
	// 	First string      `json:"first"`
	// 	Last  string      `json:"last"`
	// 	Next  interface{} `json:"next"`
	// 	Prev  interface{} `json:"prev"`
	// } `json:"links"`
}

func ListGPGKeys(token string, tfeHostname string, tfeOrganization string) (*GPGList, error) {
	// create url "https://${HOST}/api/registry/private/v2/gpg-keys?filter%5Bnamespace%5D=${provider_namespace}"
	url := fmt.Sprintf(
		"https://%s/api/registry/private/v2/gpg-keys?filter[namespace]=%s",
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
	var keys *GPGList
	err = json.Unmarshal(bodyBytes, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

type TFEBinaries struct {
	Releases []TFERelease
}

type TFERelease struct {
	ReleaseSequence      int       `json:"release_sequence"`
	PatchReleaseSequence int       `json:"patch_release_sequence"`
	ReleaseDate          time.Time `json:"release_date"`
	Required             bool      `json:"required"`
	Label                string    `json:"label"`
	ReleaseNotes         string    `json:"release_notes"`
	DownloadLink         string    `json:"download_link"`
	BuildStatus          string    `json:"build_status"`
	LastUpdateTime       time.Time `json:"last_update_time"`
	Checksum             string    `json:"checksum"`
	ImagelessChecksum    string    `json:"imageless_checksum"`
}

func ListTFEBinaries(password string, licenseId string) (*TFEBinaries, error) {
	passwordB64 := b64.URLEncoding.EncodeToString([]byte(password))
	// create url "https://api.replicated.com/market/v1/airgap/releases?license_id=${LICENSE_ID}"
	url := fmt.Sprintf(
		"https://api.replicated.com/market/v1/airgap/releases?license_id=%s",
		licenseId,
	)
	// create http Client to make calls
	client := &http.Client{}

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Set("Authorization", "Basic "+passwordB64)
	req.Header.Set("Accept", "application/json")

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
	var tfeBinaries *TFEBinaries
	err = json.Unmarshal(bodyBytes, &tfeBinaries)
	if err != nil {
		return nil, err
	}

	return tfeBinaries, nil
}

type TfeReleaseUrl struct {
	URL          string `json:"url"`
	ImagelessURL string `json:"imageless_url"`
}

func GetTFEBinary(password string, licenseId string, releaseSequence string) (*TfeReleaseUrl, error) {
	passwordB64 := b64.URLEncoding.EncodeToString([]byte(password))
	// create url "https://api.replicated.com/market/v1/airgap/images/url?license_id=${LICENSE_ID}&sequence=${release_sequence}"
	url := fmt.Sprintf(
		"https://api.replicated.com/market/v1/airgap/images/url?license_id=%s&sequence=%s",
		licenseId,
		releaseSequence,
	)
	// create http Client to make calls
	client := &http.Client{}

	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Set("Authorization", "Basic "+passwordB64)
	req.Header.Set("Accept", "application/json")

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
	var tfeUrl *TfeReleaseUrl
	err = json.Unmarshal(bodyBytes, &tfeUrl)
	if err != nil {
		return nil, err
	}

	return tfeUrl, nil
}
