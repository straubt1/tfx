// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
)

func TestRequireGPGKeyID(t *testing.T) {
	if err := requireGPGKeyID(&tfe.GPGKey{KeyID: "34365D9472D7468F"}, "34365D9472D7468F"); err != nil {
		t.Fatal(err)
	}
	if err := requireGPGKeyID(&tfe.GPGKey{KeyID: "34365d9472d7468f"}, "34365D9472D7468F"); err != nil {
		t.Fatal(err)
	}
	err := requireGPGKeyID(&tfe.GPGKey{KeyID: "AAAAAAAAAAAAAAAA"}, "34365D9472D7468F")
	if err == nil || !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("expected mismatch, got %v", err)
	}
	if err := requireGPGKeyID(nil, "34365D9472D7468F"); err == nil {
		t.Fatal("expected error for nil key")
	}
}
