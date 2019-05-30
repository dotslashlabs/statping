// Statping
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/hunterlong/statping
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"bytes"
	"text/template"
	"time"
)

// GetTemplateFuncMap returns a function map for a template to apply.
func GetTemplateFuncMap() *template.FuncMap {
	defaultDateFormat := "2006-01-02"

	resolveTimeArgs := func(layout string, timezone string) (string, *time.Location) {
		if layout == "" {
			layout = defaultDateFormat
		}

		if timezone == "" {
			timezone = "UTC"
		}

		loc, err := time.LoadLocation(timezone)
		if err != nil {
			// user-specified location not available, fallback silently to UTC
			loc, _ = time.LoadLocation("UTC")
		}

		return layout, loc
	}

	return &template.FuncMap{
		"today_date": func() string {
			return time.Now().Format(defaultDateFormat)
		},
		"today": func(layout string, timezone string) string {
			dateFmt, loc := resolveTimeArgs(layout, timezone)
			return time.Now().In(loc).Format(dateFmt)
		},
		"today_mins_ago": func(minsAgo int, layout string, timezone string) string {
			dateFmt, loc := resolveTimeArgs(layout, timezone)
			t := time.Now().In(loc)
			if minsAgo <= 0 {
				return t.Format(dateFmt)
			}
			return t.Add(time.Duration(-minsAgo) * time.Minute).Format(dateFmt)
		},
		"yesterday_date": func() string {
			return time.Now().AddDate(0, 0, -1).Format(defaultDateFormat)
		},
		"yesterday": func(layout string, timezone string) string {
			dateFmt, loc := resolveTimeArgs(layout, timezone)
			return time.Now().In(loc).AddDate(0, 0, -1).Format(dateFmt)
		},
	}
}

// ReplaceTemplateVars parses the supplied string for template functions
// and replaces them with the function's return value.
func ReplaceTemplateVars(msg string, funcMap *template.FuncMap) string {
	tmpl, err := template.New("tmpl").Delims("${", "}").Funcs(*funcMap).Parse(msg)
	if err != nil {
		return msg
	}

	var substData bytes.Buffer
	err = tmpl.Execute(&substData, "")
	if err != nil {
		return msg
	}

	return substData.String()
}
