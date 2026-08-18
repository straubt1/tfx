// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"strings"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

func TestIsConflict(t *testing.T) {
	if isConflict(nil) {
		t.Fatal("nil")
	}
	if !isConflict(errors.New("Conflict\n\nVersion has already been taken")) {
		t.Fatal("already been taken")
	}
	if !isConflict(errors.New("resource already exists")) {
		t.Fatal("already exists")
	}
	if isConflict(errors.New("failed to hash file")) {
		t.Fatal("unrelated error should not be a conflict")
	}
}

func TestSkipExistingPlatform(t *testing.T) {
	sum := strings.Repeat("a", 64)
	skip, err := skipExistingPlatform(&tfe.RegistryProviderPlatform{
		OS:                     "linux",
		Arch:                   "amd64",
		Shasum:                 sum,
		ProviderBinaryUploaded: true,
	}, sum)
	if err != nil || !skip {
		t.Fatalf("skip=%v err=%v, want skip", skip, err)
	}

	skip, err = skipExistingPlatform(&tfe.RegistryProviderPlatform{
		OS:                     "linux",
		Arch:                   "amd64",
		Shasum:                 strings.Repeat("b", 64),
		ProviderBinaryUploaded: true,
	}, sum)
	if err == nil || skip {
		t.Fatal("expected shasum mismatch error")
	}

	skip, err = skipExistingPlatform(&tfe.RegistryProviderPlatform{
		OS:                     "linux",
		Arch:                   "amd64",
		ProviderBinaryUploaded: false,
	}, sum)
	if err != nil || skip {
		t.Fatalf("incomplete platform should not skip, skip=%v err=%v", skip, err)
	}
}
