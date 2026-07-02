package main

import (
	"encoding/json"

	"github.com/flosch/pongo2/v6"
)

func init() {
	pongo2.RegisterFilter("json", filterJSON)
}

// filterJSON encodes a value as JSON, marked safe so Pongo2's HTML autoescape
// doesn't entity-mangle it — {{ business.name | json }} inside a JSON-LD
// script emits "Joe's \"Best\" Cameras" (quotes included), not
// Joe&#39;s &quot;…. Marshal \u-escapes <, >, and & by default, which also
// keeps a literal </script> from breaking out of the script block.
func filterJSON(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	b, err := json.Marshal(in.Interface())
	if err != nil {
		return nil, &pongo2.Error{Sender: "filter:json", OrigError: err}
	}
	return pongo2.AsSafeValue(string(b)), nil
}
