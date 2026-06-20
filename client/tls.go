// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package client

import (
	"crypto/tls"
	"net/http"
)

func baseTransport(sslSkipVerify bool) *http.Transport {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	if sslSkipVerify {
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		tr.TLSClientConfig.InsecureSkipVerify = true
	}
	return tr
}

// HTTPClientForTFE returns an http.Client configured for TFE API calls.
func HTTPClientForTFE(sslSkipVerify bool) *http.Client {
	if !sslSkipVerify {
		return &http.Client{}
	}
	return &http.Client{Transport: baseTransport(true)}
}
