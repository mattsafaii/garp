package main

import (
	"errors"
	"time"

	"github.com/flosch/pongo2/v6"
)

// Date wraps a date parsed from YAML (frontmatter or data) so a bare
// {{ date }} renders as a clean ISO date instead of Go's default time.Time
// string ("2026-05-22 00:00:00 +0000 UTC"). Pongo2 honors fmt.Stringer when
// rendering a value, so String() controls the default; the date/time filters
// below still format it with any Go layout.
type Date struct{ time.Time }

func (d Date) String() string { return d.Format("2006-01-02") }

func init() {
	pongo2.ReplaceFilter("date", filterDate)
	pongo2.ReplaceFilter("time", filterDate)
}

// filterDate formats a Date or time.Time with a Go reference layout, e.g.
// {{ post.date | date:"January 2, 2006" }}. With no argument it falls back to
// the ISO default, so {{ post.date | date }} is also safe. A missing value
// renders as "" — collections mix dated and undated pages, and one undated
// entry must not fail the build.
func filterDate(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var t time.Time
	switch v := in.Interface().(type) {
	case Date:
		t = v.Time
	case time.Time:
		t = v
	case nil:
		return pongo2.AsValue(""), nil
	default:
		return nil, &pongo2.Error{
			Sender:    "filter:date",
			OrigError: errors.New("filter input must be a date"),
		}
	}
	layout := param.String()
	if layout == "" {
		layout = "2006-01-02"
	}
	return pongo2.AsValue(t.Format(layout)), nil
}

// normalizeDates walks parsed YAML (maps, slices, scalars) and wraps every
// time.Time as a Date, in place, so default rendering is clean wherever a
// date appears. Returns the same value with dates replaced.
func normalizeDates(v any) any {
	switch t := v.(type) {
	case time.Time:
		return Date{t}
	case map[string]any:
		for k, val := range t {
			t[k] = normalizeDates(val)
		}
	case []any:
		for i, val := range t {
			t[i] = normalizeDates(val)
		}
	}
	return v
}
