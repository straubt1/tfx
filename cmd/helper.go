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
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

// Create or read the CV to prepare for plan
func createOrReadConfigurationVersion(ctx context.Context, client *tfe.Client, workspaceId string, cvId string, tfDirectory string, speculative bool) (*tfe.ConfigurationVersion, error) {
	var err error
	var cv *tfe.ConfigurationVersion

	if cvId == "" { // config id was not give, create a new one
		fmt.Println("Creating new Config Version")
		cv, err = client.ConfigurationVersions.Create(ctx, workspaceId, tfe.ConfigurationVersionCreateOptions{
			AutoQueueRuns: tfe.Bool(false), // wait for upload
			Speculative:   tfe.Bool(speculative),
		})
		if err != nil {
			return nil, err
		}

		// upload code
		err = client.ConfigurationVersions.Upload(ctx, cv.UploadURL, tfDirectory)
		if err != nil {
			return nil, err
		}
	} else {
		// config id was given, read
		fmt.Println("Using existing Config Version", cvId)
		cv, err = client.ConfigurationVersions.Read(ctx, cvId)
		if err != nil {
			return nil, err
		}
		// TODO: check CV is uploaded and ready to go?
	}

	return cv, nil
}

//
func getRunLogs(ctx context.Context, client *tfe.Client, planId string) error {
	var err error
	var logs io.Reader
	logs, err = client.Plans.Logs(ctx, planId)
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------------------------------------")
	// mostly found from here: https://github.com/hashicorp/terraform/blob/89f986ded6fb07e7d5f27aaf340f69c353860c12/backend/remote/backend_plan.go#L332
	reader := bufio.NewReaderSize(logs, 64*1024)
	for next := true; next; {
		var l, line []byte

		for isPrefix := true; isPrefix; {
			l, isPrefix, err = reader.ReadLine()
			if err != nil {
				if err != io.EOF {
					return err
				}
				next = false
			}
			line = append(line, l...)
		}

		if next || len(line) > 0 {
			fmt.Println(string(line))
		}
	}

	fmt.Println("------------------------------------------------------------------------")
	return nil
}

//
func getApplyLogs(ctx context.Context, client *tfe.Client, applyId string) error {
	var err error
	var logs io.Reader
	logs, err = client.Applies.Logs(ctx, applyId)
	if err != nil {
		return err
	}

	// mostly found from here: https://github.com/hashicorp/terraform/blob/89f986ded6fb07e7d5f27aaf340f69c353860c12/backend/remote/backend_plan.go#L332
	reader := bufio.NewReaderSize(logs, 64*1024)
	for next := true; next; {
		var l, line []byte

		for isPrefix := true; isPrefix; {
			l, isPrefix, err = reader.ReadLine()
			if err != nil {
				if err != io.EOF {
					return err
				}
				next = false
			}
			line = append(line, l...)
		}

		if next || len(line) > 0 {
			fmt.Println(string(line))
		}
	}

	return nil
}

// Ensure variable is up to date (upsert)
func createOrUpdateVariable(ctx context.Context, client *tfe.Client, workspaceId string, key string, value string) error {
	var err error
	var allV *tfe.VariableList

	// Read all variables and search
	// TODO: is there a better way? API doesnt expose a variable by name lookup
	allV, err = client.Variables.List(ctx, workspaceId, tfe.VariableListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(allV)

	var found *tfe.Variable
	// see if variable exists
	for i := range allV.Items {
		if allV.Items[i].Key == key {
			found = allV.Items[i]
		}
	}

	timestamp := time.Now()
	if found == nil {
		fmt.Println("Creating new Variable", key)
		_, err = client.Variables.Create(ctx, workspaceId, tfe.VariableCreateOptions{
			Key:         &key,
			Value:       &value,
			Description: tfe.String(fmt.Sprintf("Written by TFx at %s", timestamp)),
			Category:    tfe.Category("env"),
			Sensitive:   tfe.Bool(false),
		})
	} else {
		fmt.Println("Updating existing Variable", key)
		_, err = client.Variables.Update(ctx, workspaceId, found.ID, tfe.VariableUpdateOptions{
			Key:         &key,
			Value:       &value,
			Description: tfe.String(fmt.Sprintf("Written by TFx at %s", timestamp)),
			// Category:    tfe.Category("env"),
			Sensitive: tfe.Bool(false),
		})
	}
	return err
}

// Determine if a run status can be applied
func runCanBeApplied(status string) bool {
	allowed := []string{"planned", "cost_estimated", "policy_checked"}

	for _, a := range allowed {
		if status == a {
			return true
		}
	}
	return false
}
