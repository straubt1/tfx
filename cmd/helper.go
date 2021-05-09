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
	"context"
	"fmt"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

func createOrUpdateVariable(ctx context.Context, client *tfe.Client, workspaceId string, key string, value string) error {
	var err error
	var v *tfe.Variable
	var allV *tfe.VariableList

	// Read all variables and search
	// TODO: is there a better way? API doesnt expose a variable by name lookup
	allV, err = client.Variables.List(ctx, workspaceId, tfe.VariableListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(allV)

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
		v, err = client.Variables.Create(ctx, workspaceId, tfe.VariableCreateOptions{
			Key:         &key,
			Value:       &value,
			Description: tfe.String(fmt.Sprintf("Written by TFx at %s", timestamp)),
			Category:    tfe.Category("env"),
			Sensitive:   tfe.Bool(false),
		})
	} else {
		fmt.Println("Updating existing Variable", key)
		v, err = client.Variables.Update(ctx, workspaceId, found.ID, tfe.VariableUpdateOptions{
			Key:         &key,
			Value:       &value,
			Description: tfe.String(fmt.Sprintf("Written by TFx at %s", timestamp)),
			// Category:    tfe.Category("env"),
			Sensitive: tfe.Bool(false),
		})
	}
	_ = v // Golang...
	return err
}
