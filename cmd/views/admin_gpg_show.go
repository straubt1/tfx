// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

// AdminGPGShowView handles rendering for admin gpg show command
type AdminGPGShowView struct {
	*BaseView
}

func NewAdminGPGShowView() *AdminGPGShowView {
	return &AdminGPGShowView{
		BaseView: NewBaseView(),
	}
}

// gpgKeyShowOutput is a JSON-safe representation of a GPG key
type gpgKeyShowOutput struct {
	KeyID      string    `json:"keyId"`
	Namespace  string    `json:"namespace"`
	AsciiArmor string    `json:"asciiArmor"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Render renders a single GPG key's details
func (v *AdminGPGShowView) Render(key *tfe.GPGKey) error {
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
		{Key: "Namespace", Value: key.Namespace},
		{Key: "Created", Value: key.CreatedAt.Format(time.RFC3339)},
		{Key: "Updated", Value: key.UpdatedAt.Format(time.RFC3339)},
		{Key: "ASCII Armor", Value: "\n" + key.AsciiArmor},
	}

	return v.Output().RenderProperties(properties)
}
