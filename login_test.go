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
	"testing"
)

var testBotToken = "123456789:abcdefGHIJKLmnopqrSTUVWXyz123456789"

func TestConstructString_WithEmptyPairs(t *testing.T) {
	if s := constructCheckString([]pair{}); s != "" {
		t.Errorf("result should be empty string, but was %v", s)
	}
}

func TestConstructString_WithOnlyID(t *testing.T) {
	if s := constructCheckString([]pair{
		{"id", "12345678"},
	}); s != "id=12345678" {
		t.Errorf("result should be 'id=12341234', but was %v", s)
	}
}

func TestConstructString_WithFullMap(t *testing.T) {
	if s := constructCheckString([]pair{
		{"id", "12345678"},
		{"first_name", "John 🕶"},
		{"last_name", "Smith"},
		{"username", "jsmith"},
		{"photo_url", "https://t.me/i/userpic/320/jsmith.jpg"},
		{"auth_date", "1512345678"},
	}); s != "auth_date=1512345678\nfirst_name=John 🕶\nid=12345678\nlast_name=Smith\nphoto_url=https://t.me/i/userpic/320/jsmith.jpg\nusername=jsmith" {
		t.Errorf("result should be 'auth_date=1512345678\nfirst_name=John 🕶\nid=12345678\nlast_name=Smith\nphoto_url=https://t.me/i/userpic/320/jsmith.jpg\nusername=jsmith', but was %v", s)
	}
}
