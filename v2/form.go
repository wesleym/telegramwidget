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

package telegramwidget

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
)

// ErrNotSingleValue indicates that the provided form has zero or more than one value for one of the parameters that
// contains user data.
var ErrNotSingleValue = errors.New("zero or multiple values for a key in form")

// ConvertAndVerifyForm accepts form encoded data from the provided form and parses it into the returned User. The hash
// property of the input form is used to validate the user data before it is returned.
func ConvertAndVerifyForm(f url.Values, tokenHash []byte) (User, error) {
	u, ps, expectedMAC, err := parseUserFromForm(f)
	if err != nil {
		return u, err
	}

	if !validate(ps, tokenHash, expectedMAC) {
		return u, ErrInvalidHash
	}

	return u, nil
}

func parseUserFromForm(f url.Values) (User, []pair, []byte, error) {
	var tu User
	// There are six supported properties.
	ps := make([]pair, 0, 6)
	expectedMAC := make([]byte, sha256.Size)

	for k, vs := range f {
		if len(vs) != 1 {
			return tu, nil, expectedMAC, ErrNotSingleValue
		}
		v := vs[0]

		switch k {
		case "id":
			ps = append(ps, pair{"id", v})
			var err error
			if tu.ID, err = strconv.ParseInt(v, 10, 64); err != nil {
				return tu, nil, expectedMAC, err
			}
		case "first_name":
			ps = append(ps, pair{"first_name", v})
			tu.FirstName = v
		case "last_name":
			ps = append(ps, pair{"last_name", v})
			tu.LastName = v
		case "username":
			ps = append(ps, pair{"username", v})
			tu.Username = v
		case "photo_url":
			ps = append(ps, pair{"photo_url", v})
			var err error
			if tu.PhotoURL, err = url.Parse(v); err != nil {
				return tu, nil, expectedMAC, err
			}
		case "auth_date":
			ps = append(ps, pair{"auth_date", v})
			// Fractional seconds are lost by this conversion.
			seconds, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return tu, nil, expectedMAC, err
			}
			tu.AuthDate = time.Unix(seconds, 0)
		case "hash":
			// This is only used to check validity, then is dropped.
			if hex.DecodedLen(len(v)) != sha256.Size {
				return tu, nil, expectedMAC, fmt.Errorf("hash must be 64 characters long, but wasn't")
			}
			if _, err := hex.Decode(expectedMAC, []byte(v)); err != nil {
				return tu, nil, expectedMAC, fmt.Errorf("failure to decode incoming hash: %v", err)
			}
		default:
			log.Printf("unexpected field in Telegram user: %s", k)
		}
	}

	return tu, ps, expectedMAC, nil
}
