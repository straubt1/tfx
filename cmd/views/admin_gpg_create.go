// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

// AdminGPGCreateView handles rendering for admin gpg create command
type AdminGPGCreateView struct {
	*BaseView
}

func NewAdminGPGCreateView() *AdminGPGCreateView {
	return &AdminGPGCreateView{
		BaseView: NewBaseView(),
	}
}

// Render renders the created GPG key
func (v *AdminGPGCreateView) Render(key *tfe.GPGKey) error {
	if v.IsJSON() {
		output := gpgKeyShowOutput{
			KeyID:      key.KeyID,
			Namespace:  key.Namespace,
			AsciiArmor: key.AsciiArmor,
			CreatedAt:  key.CreatedAt,
			UpdatedAt:  key.UpdatedAt,
		}
		return v.Output().RenderJSON(output)
	}

	// Terminal mode: render key fields
	properties := []PropertyPair{
		{Key: "Key ID", Value: key.KeyID},
		{Key: "Created", Value: key.CreatedAt.Format(time.RFC3339)},
		{Key: "Updated", Value: key.UpdatedAt.Format(time.RFC3339)},
		{Key: "ASCII Armor", Value: "\n" + key.AsciiArmor},
	}

	return v.Output().RenderProperties(properties)
}
