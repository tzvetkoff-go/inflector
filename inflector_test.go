package inflector_test

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/tzvetkoff-go/inflector"
)

type test struct {
	where    string
	f        func(string) string
	arg      string
	expected string
}

var tests = []test{
	{
		where:    here(),
		f:        inflector.Camelize,
		arg:      "add_user_id_to_users",
		expected: "AddUserIDToUsers",
	},
	{
		where:    here(),
		f:        inflector.Underscore,
		arg:      "AddUserIDToUsers",
		expected: "add_user_id_to_users",
	},

	{
		where:    here(),
		f:        inflector.Camelize,
		arg:      "restful_api",
		expected: "RESTfulAPI",
	},
	{
		where:    here(),
		f:        inflector.Underscore,
		arg:      "RESTfulAPI",
		expected: "restful_api",
	},

	{
		where:    here(),
		f:        inflector.Parameterize,
		arg:      "RESTfulAPI",
		expected: "restful-api",
	},
	{
		where:    here(),
		f:        inflector.Parameterize,
		arg:      "Крали Марко (ID: 31337) пише RESTful API",
		expected: "krali-marko-id-31337-pishe-restful-api",
	},

	{
		where:    here(),
		f:        inflector.Pluralize,
		arg:      "person",
		expected: "people",
	},
	{
		where:    here(),
		f:        inflector.Pluralize,
		arg:      "octopus",
		expected: "octopi",
	},

	{
		where:    here(),
		f:        inflector.Singularize,
		arg:      "people",
		expected: "person",
	},
	{
		where:    here(),
		f:        inflector.Singularize,
		arg:      "money",
		expected: "money",
	},
}

func init() {
	inflector.DefaultInflector.AddAcronym("id", "ID")
	inflector.DefaultInflector.AddAcronym("restful", "RESTful")
	inflector.DefaultInflector.AddAcronym("api", "API")
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
