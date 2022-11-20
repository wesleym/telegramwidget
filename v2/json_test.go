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
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestConvertAndVerifyJSON_WithValidCredentials(t *testing.T) {
	u, err := ConvertAndVerifyJSON(strings.NewReader(`{
		"auth_date": 1512345678,
		"first_name": "John 🕶",
		"hash": "25409759c10beb29bd3f3fe1d16ee0605ac82eb2907d886e196d481371b91501",
		"id": 12345678,
		"last_name": "Smith",
		"photo_url": "https://t.me/i/userpic/320/jsmith.jpg",
		"username": "jsmith"
	}`), testBotTokenHash)
	if err != nil {
		t.Fatalf("failed to convert and verify: %v", err)
	}
	if !time.Date(2017, time.December, 4, 0, 1, 18, 0, time.UTC).Equal(u.AuthDate) {
		t.Errorf("auth date should be 2017-12-04T00:01:18Z, but was %v", u.AuthDate)
	}

	if u.FirstName != "John 🕶" {
		t.Errorf("first name should be John 🕶, but was %v", u.FirstName)
	}

	if u.ID != 12345678 {
		t.Errorf("ID should be 12345678, but was %d", u.ID)
	}

	if u.LastName != "Smith" {
		t.Errorf("last name should be Smith, but was %v", u.LastName)
	}

	p := url.URL{Scheme: "https", Host: "t.me", Path: "/i/userpic/320/jsmith.jpg"}
	if *u.PhotoURL != p {
		t.Errorf("photo URL should be https://t.me/i/userpic/320/jsmith.jpg but was %v", u.PhotoURL)
	}

	if u.Username != "jsmith" {
		t.Errorf("username should be jsmith, but was %s", u.Username)
	}
}

func TestConvertAndVerifyJSON_WithoutHash(t *testing.T) {
	_, err := ConvertAndVerifyJSON(strings.NewReader(`{
		"auth_date": 1512345678,
		"first_name": "John 🕶",
		"id": 12345678,
		"last_name": "Smith",
		"photo_url": "https://t.me/i/userpic/320/jsmith.jpg",
		"username": "jsmith"
	}`), testBotTokenHash)
	if err != ErrInvalidHash {
		t.Errorf("expected ErrInvalidHash, but was %v", err)
	}
}

func TestConvertAndVerifyJSON_WithIncorrectHash(t *testing.T) {
	_, err := ConvertAndVerifyJSON(strings.NewReader(`{
		"auth_date": 1512345678,
		"first_name": "John 🕶",
		"hash": "0000000000000000000000000000000000000000000000000000000000000000",
		"id": 12345678,
		"last_name": "Smith",
		"photo_url": "https://t.me/i/userpic/320/jsmith.jpg",
		"username": "jsmith"
	}`), testBotTokenHash)
	if err != ErrInvalidHash {
		t.Errorf("expected ErrInvalidHash, but was %v", err)
	}
}

func TestConvertAndVerifyJSON_WithMalformedJSON(t *testing.T) {
	_, err := ConvertAndVerifyJSON(strings.NewReader("food"), testBotTokenHash)
	if err == nil {
		t.Errorf("should have returned error, but was nil")
	}
}

func TestConvertAndVerifyJSON_WithWrongJSONValueType(t *testing.T) {
	_, err := ConvertAndVerifyJSON(strings.NewReader(`"food"`), testBotTokenHash)
	if err == nil {
		t.Errorf("should have returned error, but was nil")
	}
}

func TestConvertAndVerifyJSON_WithNestedObject(t *testing.T) {
	_, err := ConvertAndVerifyJSON(strings.NewReader(`{
		"auth_date": {},
	}`), testBotTokenHash)
	if err == nil {
		t.Errorf("should have returned error, but was nil")
	}
}

func TestConvertAndVerifyJSON_WithNestedArray(t *testing.T) {
	_, err := ConvertAndVerifyJSON(strings.NewReader(`{
		"auth_date": [],
	}`), testBotTokenHash)
	if err == nil {
		t.Errorf("should have returned error, but was nil")
	}
}

func TestConvertAndVerifyJSON_WithEmptyString(t *testing.T) {
	_, err := ConvertAndVerifyJSON(strings.NewReader(""), testBotTokenHash)
	if err == nil {
		t.Errorf("should have returned error, but was nil")
	}
}

func TestConvertAndVerifyJSON_MarksMissingFields(t *testing.T) {
	u, err := ConvertAndVerifyJSON(strings.NewReader(`{
		"auth_date": 1512345678,
		"id": 12345678,
		"hash": "180f7d26839de06e6ecb26148f181553d24e1c62153400da55ae31483ee62ad3"
	}`), testBotTokenHash)
	if err != nil {
		t.Errorf("failed to convert and verify: %v", err)
	}
	if u.FirstName != "" {
		t.Errorf("first name should be absent, but was %s", u.FirstName)
	}
	if u.LastName != "" {
		t.Errorf("last name should be absent, but was %s", u.LastName)
	}
	if u.PhotoURL != nil {
		t.Errorf("photo URL should be absent, but was %s", u.PhotoURL)
	}
	if u.Username != "" {
		t.Errorf("username should be absent, but was %s", u.Username)
	}
}
