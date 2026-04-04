// Copyright 2006-2019 WebPKI.org (http://webpki.org).
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

package jsoncanonicalizer

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const invalidPattern uint64 = 0x7ff0000000000000

func NumberToJSON(ieeeF64 float64) (res string, err error) {
	ieeeU64 := math.Float64bits(ieeeF64)

	if (ieeeU64 & invalidPattern) == invalidPattern {
		return "null", errors.New("Invalid JSON number: " + strconv.FormatUint(ieeeU64, 16))
	}

	if ieeeF64 == 0 {
		return "0", nil
	}

	sign := ""
	if ieeeF64 < 0 {
		ieeeF64 = -ieeeF64
		sign = "-"
	}

	format := byte('e')
	if ieeeF64 < 1e+21 && ieeeF64 >= 1e-6 {
		format = 'f'
	}

	es6Formatted := strconv.FormatFloat(ieeeF64, format, -1, 64)
	if exponent := strings.IndexByte(es6Formatted, 'e'); exponent > 0 {
		if es6Formatted[exponent+2] == '0' {
			es6Formatted = es6Formatted[:exponent+2] + es6Formatted[exponent+3:]
		}
	}
	return sign + es6Formatted, nil
}
