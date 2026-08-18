// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"strings"
	"testing"
)

func TestUploadCommands(t *testing.T) {
	cmds := uploadCommands(&RegistryProviderDownloadResult{
		Directory:        "/tmp/providers/chainguard-dev/cosign/0.4.16",
		KeyID:            "5BBEE08F6BF07616",
		GPGPublicKeyPath: "/tmp/providers/chainguard-dev/cosign/0.4.16/5BBEE08F6BF07616.asc",
	})
	if len(cmds) != 1 {
		t.Fatalf("commands = %v", cmds)
	}
	if !strings.Contains(cmds[0], "tfx registry provider version create --directory /tmp/providers/chainguard-dev/cosign/0.4.16") {
		t.Fatalf("commands = %v", cmds)
	}
}
