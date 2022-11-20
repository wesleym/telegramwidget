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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"
)

// ConvertAndVerifyJSON accepts JSON from the provided reader and parses it into the returned User. The hash property of the
// input JSON is used to validate the user data before it is returned.
func ConvertAndVerifyJSON(r io.Reader, tokenHash []byte) (User, error) {
	u, ps, expectedMAC, err := parseUserFromJSON(r)
	if err != nil {
		return u, err
	}

	if !validate(ps, tokenHash, expectedMAC) {
		return u, ErrInvalidHash
	}

	return u, nil
}

func parseUserFromJSON(r io.Reader) (User, []pair, []byte, error) {
	d := json.NewDecoder(r)
	d.UseNumber()
	var tu User
	// There are six supported properties.
	ps := make([]pair, 0, 6)
	expectedMAC := make([]byte, sha256.Size)

	if t, err := d.Token(); err == io.EOF {
		return tu, nil, expectedMAC, fmt.Errorf("expected start of object, got EOF")
	} else if err != nil {
		return tu, nil, expectedMAC, fmt.Errorf("expected start of object, got error: %v", err)
	} else if d, ok := t.(json.Delim); !ok || d != '{' {
		return tu, nil, expectedMAC, fmt.Errorf("expected start of object, got token: %v", t)
	}

	for d.More() {
		k, err := d.Token()
		if err != nil {
			return tu, nil, expectedMAC, fmt.Errorf("expected key, got error: %v", err)
		} else if d, ok := k.(json.Delim); ok {
			// This case should be impossible in well-formed JSON. We're inside an object, so keys should always be
			// strings.
			return tu, nil, expectedMAC, fmt.Errorf("expected key, got delimeter: %v", d)
		}

		v, err := d.Token()
		if err != nil {
			return tu, nil, expectedMAC, fmt.Errorf("expected value, got error: %v", err)
		} else if _, ok := v.(json.Delim); ok {
			return tu, nil, expectedMAC, fmt.Errorf("expected value, got delimeter: %v", v)
		}

		switch k {
		case "id":
			id := v.(json.Number)
			ps = append(ps, pair{"id", id.String()})
			if tu.ID, err = id.Int64(); err != nil {
				return tu, nil, expectedMAC, err
			}
		case "first_name":
			firstName := v.(string)
			ps = append(ps, pair{"first_name", firstName})
			tu.FirstName = firstName
		case "last_name":
			lastName := v.(string)
			ps = append(ps, pair{"last_name", lastName})
			tu.LastName = lastName
		case "username":
			username := v.(string)
			ps = append(ps, pair{"username", username})
			tu.Username = username
		case "photo_url":
			photoURL := v.(string)
			ps = append(ps, pair{"photo_url", photoURL})
			if tu.PhotoURL, err = url.Parse(photoURL); err != nil {
				return tu, nil, expectedMAC, err
			}
		case "auth_date":
			authDate := v.(json.Number)
			ps = append(ps, pair{"auth_date", authDate.String()})
			// Fractional seconds are lost by this conversion.
			seconds, err := authDate.Int64()
			if err != nil {
				return tu, nil, expectedMAC, err
			}
			tu.AuthDate = time.Unix(seconds, 0)
		case "hash":
			// This is only used to check validity, then is dropped.
			hash := v.(string)
			if hex.DecodedLen(len(hash)) != sha256.Size {
				return tu, nil, expectedMAC, fmt.Errorf("hash must be 64 characters long, but wasn't")
			}
			if _, err := hex.Decode(expectedMAC, []byte(hash)); err != nil {
				return tu, nil, expectedMAC, fmt.Errorf("failure to decode incoming hash: %v", err)
			}
		default:
			log.Printf("unexpected field in Telegram user: %s", k)
		}
	}

	// The following error cases should only occur if the input JSON is not well-formed.
	// The only way for d.More to return false is at the end of an object, and the end of an object should only be
	// indicated by a '}' delimiter.
	if t, err := d.Token(); err == io.EOF {
		return tu, nil, expectedMAC, fmt.Errorf("expected end of object, got EOF")
	} else if err != nil {
		return tu, nil, expectedMAC, fmt.Errorf("expected end of object, got error: %v", err)
	} else if d, ok := t.(json.Delim); !ok || d != '}' {
		return tu, nil, expectedMAC, fmt.Errorf("expected end of object, got token: %v", t)
	}

	// A JSON document should only represent a single value, right? If so, the only possibility after closing the top
	// level object is EOF.
	if _, err := d.Token(); err == nil {
		return tu, nil, expectedMAC, fmt.Errorf("expected EOF, but got a token")
	} else if err != io.EOF {
		return tu, nil, expectedMAC, fmt.Errorf("expected EOF, but got a different error: %v", err)
	}

	return tu, ps, expectedMAC, nil
}
