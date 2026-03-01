// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type StateVersionCreateView struct{ *BaseView }

func NewStateVersionCreateView() *StateVersionCreateView {
	return &StateVersionCreateView{NewBaseView()}
}

func (v *StateVersionCreateView) Render(sv *tfe.StateVersion) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(map[string]interface{}{
			"id":               sv.ID,
			"terraformVersion": sv.TerraformVersion,
			"serial":           sv.Serial,
		})
	}
	props := []PropertyPair{
		{Key: "Created Serial", Value: sv.Serial},
	}
	return v.Output().RenderProperties(props)
}
