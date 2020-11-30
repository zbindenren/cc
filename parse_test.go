package cc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"gopkg.in/yaml.v3"
)

func TestParse(t *testing.T) {
	type tests []struct {
		Name     string `yaml:"name"`
		Message  string `yaml:"message"`
		Expected struct {
			Fail        bool    `yaml:"mustFail"`
			Breaking    bool    `yaml:"breaking"`
			Scope       string  `yaml:"scope"`
			Type        string  `yaml:"type"`
			Description string  `yaml:"description"`
			Body        string  `yaml:"body"`
			Footer      Footers `yaml:"footer"`
		} `yaml:"expected"`
	}

	allTests := tests{}
	d, err := ioutil.ReadFile("tests.yaml")
	require.NoError(t, err)

	err = yaml.Unmarshal(d, &allTests)
	require.NoError(t, err)

	for i := range allTests {
		tc := allTests[i]

		t.Run(tc.Name, func(t *testing.T) {
			// nolint: gocritic
			// fmt.Println(tc.Message)
			c, err := Parse(tc.Message)

			if tc.Expected.Fail {
				// nolint: gocritic
				// fmt.Println(err) // nolint: gocritic
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.Expected.Scope, c.Header.Scope, "invalid header scope")
			assert.Equal(t, tc.Expected.Type, c.Header.Type, "invalid header type")
			assert.Equal(t, tc.Expected.Breaking, c.isBreaking, "is breaking change but should not be")
			assert.Equal(t, tc.Expected.Body, c.Body, "invalid body")
			assert.EqualValues(t, tc.Expected.Footer, c.Footer, "invalid footer")
		})
	}
}

func TestBreakingMessage(t *testing.T) {
	tt := []struct {
		commit          Commit
		expectedMessage string
	}{
		{
			Commit{
				Header: Header{
					Description: "description",
				},
			},
			"",
		},
		{
			Commit{
				isBreaking: true,
				Header: Header{
					Description: "description",
				},
			},
			"description",
		},
		{
			Commit{
				isBreaking: true,
				Header: Header{
					Description: "description",
				},
				Footer: Footers{
					Footer{
						Token: "BREAKING CHANGE",
						Value: "breaking change",
					},
				},
			},
			"breaking change",
		},
		{
			Commit{
				isBreaking: true,
				Header: Header{
					Description: "description",
				},
				Footer: Footers{
					Footer{
						Token: "BREAKING-CHANGE",
						Value: "breaking change",
					},
				},
			},
			"breaking change",
		},
		{
			Commit{
				Header: Header{
					Description: "description",
				},
				Footer: Footers{
					Footer{
						Token: "BREAKING-CHANGE",
						Value: "breaking change",
					},
				},
			},
			"breaking change",
		},
	}

	for i := range tt {
		tc := tt[i]

		assert.Equal(t, tc.expectedMessage, tc.commit.BreakingMessage())
	}
}

func ExampleParse() {
	msg := `fix(compiler): correct minor typos in code

see the issue for details

on typos fixed.

Reviewed-by: Z
Refs #133`

	c, _ := Parse(msg)
	d, _ := json.MarshalIndent(c, "", "  ")

	fmt.Println(string(d))
	// Output:
	// {
	//   "Header": {
	//     "Type": "fix",
	//     "Scope": "compiler",
	//     "Description": "correct minor typos in code"
	//   },
	//   "Body": "see the issue for details\n\non typos fixed.",
	//   "Footer": [
	//     {
	//       "Token": "Reviewed-by",
	//       "Value": "Z"
	//     },
	//     {
	//       "Token": "Refs",
	//       "Value": "#133"
	//     }
	//   ]
	// }
}
