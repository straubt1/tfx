// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

// AdminGPGListView handles rendering for admin gpg list command
type AdminGPGListView struct {
	*BaseView
}

func NewAdminGPGListView() *AdminGPGListView {
	return &AdminGPGListView{
		BaseView: NewBaseView(),
	}
}

// gpgKeyListOutput is a JSON-safe representation of a GPG key for list views
type gpgKeyListOutput struct {
	KeyID     string    `json:"keyId"`
	Namespace string    `json:"namespace"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// Render renders GPG keys list
func (v *AdminGPGListView) Render(keys []*tfe.GPGKey) error {
	if v.IsJSON() {
		output := make([]gpgKeyListOutput, len(keys))
		for i, k := range keys {
			output[i] = gpgKeyListOutput{
				KeyID:     k.KeyID,
				Namespace: k.Namespace,
				UpdatedAt: k.UpdatedAt,
				CreatedAt: k.CreatedAt,
			}
		}
		return v.Output().RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Key ID", "Namespace", "Updated At", "Created At"}
	rows := make([][]interface{}, len(keys))
	for i, k := range keys {
		rows[i] = []interface{}{
			k.KeyID,
			k.Namespace,
			k.UpdatedAt.Format(time.RFC3339),
			k.CreatedAt.Format(time.RFC3339),
		}
	}

	return v.Output().RenderTable(headers, rows)
}
