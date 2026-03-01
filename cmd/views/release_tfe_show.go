// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"github.com/dustin/go-humanize"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// ReleaseTfeShowView handles rendering for release tfe show command
type ReleaseTfeShowView struct {
	*BaseView
}

func NewReleaseTfeShowView() *ReleaseTfeShowView {
	return &ReleaseTfeShowView{
		BaseView: NewBaseView(),
	}
}

// releaseTfeShowOutput is a JSON-safe representation of a TFE release for show views
type releaseTfeShowOutput struct {
	Tag     string                    `json:"tag"`
	Digest  string                    `json:"digest"`
	Created string                    `json:"created"`
	OS      string                    `json:"os"`
	Arch    string                    `json:"arch"`
	Size    int64                     `json:"size"`
	Layers  []releaseTfeLayerOutput   `json:"layers,omitempty"`
	History []releaseTfeHistoryOutput `json:"history,omitempty"`
}

type releaseTfeLayerOutput struct {
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	MediaType string `json:"mediaType"`
}

type releaseTfeHistoryOutput struct {
	Created    string `json:"created,omitempty"`
	CreatedBy  string `json:"createdBy,omitempty"`
	Author     string `json:"author,omitempty"`
	Comment    string `json:"comment,omitempty"`
	EmptyLayer bool   `json:"emptyLayer,omitempty"`
}

// Render renders TFE release details
func (v *ReleaseTfeShowView) Render(release map[string]interface{}) error {
	if v.IsJSON() {
		// Convert layers and history to JSON-safe representations
		var layers []releaseTfeLayerOutput
		if layersList, ok := release["Layers"].([]v1.Layer); ok {
			for _, layer := range layersList {
				digest, _ := layer.Digest()
				size, _ := layer.Size()
				mediaType, _ := layer.MediaType()
				layers = append(layers, releaseTfeLayerOutput{
					Digest:    digest.String(),
					Size:      size,
					MediaType: string(mediaType),
				})
			}
		}

		var history []releaseTfeHistoryOutput
		if historyList, ok := release["History"].([]v1.History); ok {
			for _, h := range historyList {
				history = append(history, releaseTfeHistoryOutput{
					Created:    h.Created.Format("01-02-2006 15:04:05"),
					CreatedBy:  h.CreatedBy,
					Author:     h.Author,
					Comment:    h.Comment,
					EmptyLayer: h.EmptyLayer,
				})
			}
		}

		output := releaseTfeShowOutput{
			Tag:     release["Tag"].(string),
			Digest:  release["Digest"].(string),
			Created: release["Created"].(string),
			OS:      release["OS"].(string),
			Arch:    release["Arch"].(string),
			Size:    release["Size"].(int64),
			Layers:  layers,
			History: history,
		}
		return v.Output().RenderJSON(output)
	}

	// Terminal mode: render as key-value properties
	properties := []PropertyPair{
		{Key: "Tag", Value: release["Tag"].(string)},
		{Key: "Digest", Value: release["Digest"].(string)},
		{Key: "Created", Value: release["Created"].(string)},
		{Key: "OS", Value: release["OS"].(string)},
		{Key: "Architecture", Value: release["Arch"].(string)},
		{Key: "Size", Value: humanize.Bytes(uint64(release["Size"].(int64)))},
	}

	if err := v.Output().RenderProperties(properties); err != nil {
		return err
	}

	return nil
}
