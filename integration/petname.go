// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

//go:build integration
// +build integration

package integration

import (
	"fmt"
	"math/rand"
	"time"
)

var petAdjectives = []string{
	"amber", "brave", "calm", "daring", "eager", "fancy", "gentle", "happy",
	"iconic", "jolly", "keen", "lively", "merry", "noble", "quick", "rusty",
}

var petNouns = []string{
	"badger", "condor", "dragon", "eagle", "falcon", "gecko", "heron", "iguana",
	"jaguar", "koala", "lemur", "marten", "newt", "otter", "panda", "quail",
}

// uniquePetName returns a random adjective-noun suffix for unique resource names.
func uniquePetName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%s", petAdjectives[r.Intn(len(petAdjectives))], petNouns[r.Intn(len(petNouns))])
}
