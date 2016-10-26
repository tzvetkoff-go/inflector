package inflector

import (
	"fmt"
	"path"
	"runtime"
	"testing"
)

type test struct {
	where		string
	f			func(string) string
	arg			string
	expected	string
}

var tests = []test{
	{
		where:		here(),
		f:			Camelize,
		arg:		"add_user_id_to_users",
		expected:	"AddUserIDToUsers",
	},
	{
		where:		here(),
		f:			Underscore,
		arg:		"AddUserIDToUsers",
		expected:	"add_user_id_to_users",
	},

	{
		where:		here(),
		f:			Camelize,
		arg:		"restful_api",
		expected:	"RESTfulAPI",
	},
	{
		where:		here(),
		f:			Underscore,
		arg:		"RESTfulAPI",
		expected:	"restful_api",
	},

	{
		where:		here(),
		f:			Parameterize,
		arg:		"RESTfulAPI",
		expected:	"restful-api",
	},
	{
		where:		here(),
		f:			Parameterize,
		arg:		"Крали Марко (ID: 31337) пише RESTful API",
		expected:	"krali-marko-id-31337-pishe-restful-api",
	},

	{
		where:		here(),
		f:			Pluralize,
		arg:		"person",
		expected:	"people",
	},
	{
		where:		here(),
		f:			Pluralize,
		arg:		"octopus",
		expected:	"octopi",
	},

	{
		where:		here(),
		f:			Singularize,
		arg:		"people",
		expected:	"person",
	},
	{
		where:		here(),
		f:			Singularize,
		arg:		"money",
		expected:	"money",
	},
}

func init() {
	DefaultInflector.AddAcronym("id", "ID")
	DefaultInflector.AddAcronym("restful", "RESTful")
	DefaultInflector.AddAcronym("api", "API")
}

func TestInflector(t *testing.T) {
	for _, tt := range tests {
		r := tt.f(tt.arg)

		if r != tt.expected {
			t.Errorf("%s: got %v, expected %v", tt.where, r, tt.expected)
		}
	}
}

func here() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", path.Base(file), line)
}
