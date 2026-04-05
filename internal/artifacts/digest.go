package artifacts

import "regexp"

var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
