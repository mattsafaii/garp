package main

import (
	"testing"

	"github.com/flosch/pongo2/v6"
)

// The json filter must survive Pongo2's HTML autoescape untouched — that's
// the whole point (quotes as \", not &quot;).
func TestJSONFilter(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{`Joe's "Best" Cameras`, `"Joe's \"Best\" Cameras"`},
		// Marshal escapes HTML — a value can't break out of the script block
		{`</script>`, `"\u003c/script\u003e"`},
		{[]any{"a", "b"}, `["a","b"]`},
		{nil, `null`},
	}
	for _, c := range cases {
		tpl, err := pongo2.FromString(`{{ v | json }}`)
		if err != nil {
			t.Fatal(err)
		}
		out, err := tpl.Execute(pongo2.Context{"v": c.in})
		if err != nil {
			t.Fatal(err)
		}
		if out != c.want {
			t.Errorf("json(%v) = %s, want %s", c.in, out, c.want)
		}
	}
}
