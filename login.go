// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package telegramwidget provides a data type to represent a Telegram user and utilities to parse and verify a Telegram
// user as returned from the Telegram login widget. This library currently supports version 4 of the login widget only.
//
// For more detail about the Telegram login widget, see https://core.telegram.org/widgets/login.
package telegramwidget

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"sort"
)

// ErrInvalidHash indicates that the data received could not be authenticated with its hash property and the provided
// token.
var ErrInvalidHash = errors.New("the hash is invalid")

// constructCheckString accepts key-value pairs of user data parameters and constructs a string that can be hashed to
// authenticate it. If passing in pairs from the Telegram login widget, make sure not to include the parameter named
// "hash", as the hash itself is not used in computing the hash.
func constructCheckString(r []pair) string {
	l := len(r)
	if l <= 0 {
		return ""
	}

	sort.Slice(r, func(i, j int) bool {
		return r[i].key < r[j].key
	})

	s := make([]byte, 0, estimateSize(r))

	s = append(s, r[0].key...)
	s = append(s, '=')
	s = append(s, r[0].value...)
	for _, p := range r[1:] {
		s = append(s, '\n')
		s = append(s, p.key...)
		s = append(s, '=')
		s = append(s, p.value...)
	}

	return string(s)
}

func validate(ps []pair, token string, expectedMAC []byte) bool {
	s := constructCheckString(ps)
	key := sha256.Sum256([]byte(token))
	mac := hmac.New(sha256.New, key[:])
	mac.Write([]byte(s))
	computedMAC := mac.Sum(nil)
	// This will fail anyway if the input JSON doesn't include a hash.
	return hmac.Equal(expectedMAC, computedMAC)
}
