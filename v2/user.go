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
	"time"
)

// A User is a Telegram user. All of the data returned from the Telegram login
// widget is represented in this type.
//
// Absent fields are parsed as their zero values. For example, when username is
// not provided, the Username field contains the empty string.
type User struct {
	AuthDate  time.Time
	FirstName string
	ID        int64
	LastName  string
	PhotoURL  *url.URL
	Username  string
}
