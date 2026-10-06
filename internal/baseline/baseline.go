package baseline

import _ "embed"

// Raw contains the baseline JSON embedded into the scanner binary.
//
//go:embed baseline.json
var Raw []byte
