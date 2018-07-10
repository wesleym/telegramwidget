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

type pair struct {
	key   string
	value string
}

func estimateSize(ps []pair) int {
	l := len(ps)
	if l <= 0 {
		return 0
	}
	s := 2*l - 1
	for _, p := range ps {
		s += len(p.key) + len(p.value)
	}
	return s
}
