package main

import (
	"testing"
	"time"

	"github.com/flosch/pongo2/v6"
)

func TestDateDefaultRender(t *testing.T) {
	tpl, err := pongo2.FromString("{{ date }}")
	if err != nil {
		t.Fatal(err)
	}
	out, err := tpl.Execute(pongo2.Context{"date": Date{time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC)}})
	if err != nil {
		t.Fatal(err)
	}
	if out != "2026-05-22" {
		t.Errorf("bare {{ date }} = %q, want ISO 2026-05-22", out)
	}
}

func TestDateFilter(t *testing.T) {
	cases := map[string]string{
		`{{ date | date:"January 2, 2006" }}`: "May 22, 2026",
		`{{ date | date:"2006/01/02" }}`:      "2026/05/22",
		`{{ date | date }}`:                   "2026-05-22", // no arg → ISO default
	}
	ctx := pongo2.Context{"date": Date{time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC)}}
	for tmpl, want := range cases {
		tpl, err := pongo2.FromString(tmpl)
		if err != nil {
			t.Fatal(err)
		}
		out, err := tpl.Execute(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if out != want {
			t.Errorf("%s = %q, want %q", tmpl, out, want)
		}
	}
}

// A collection mixes dated and undated pages — {{ p.date | date }} over it
// must render "" for the undated ones, not fail the build.
func TestDateFilterNilRendersEmpty(t *testing.T) {
	for _, tmpl := range []string{`{{ date | date }}`, `{{ date | date:"January 2, 2006" }}`} {
		tpl, err := pongo2.FromString(tmpl)
		if err != nil {
			t.Fatal(err)
		}
		out, err := tpl.Execute(pongo2.Context{"date": nil})
		if err != nil {
			t.Fatalf("%s: %v", tmpl, err)
		}
		if out != "" {
			t.Errorf("%s = %q, want empty", tmpl, out)
		}
	}
}

func TestDateFilterAcceptsRawTime(t *testing.T) {
	tpl, _ := pongo2.FromString(`{{ date | date:"2006" }}`)
	out, err := tpl.Execute(pongo2.Context{"date": time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if out != "2026" {
		t.Errorf("out = %q, want 2026", out)
	}
}

func TestNormalizeDates(t *testing.T) {
	tm := time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC)
	in := map[string]any{
		"date":  tm,
		"title": "Post",
		"meta":  map[string]any{"updated": tm},
		"events": []any{
			map[string]any{"on": tm},
		},
	}
	normalizeDates(in)

	if _, ok := in["date"].(Date); !ok {
		t.Errorf("top-level date = %T, want Date", in["date"])
	}
	if in["title"] != "Post" {
		t.Errorf("non-date value mutated: %v", in["title"])
	}
	if _, ok := in["meta"].(map[string]any)["updated"].(Date); !ok {
		t.Error("nested map date not normalized")
	}
	if _, ok := in["events"].([]any)[0].(map[string]any)["on"].(Date); !ok {
		t.Error("date inside slice not normalized")
	}
}
