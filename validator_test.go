package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

type ruleTestData struct {
	ruleName                     string
	correctValue                 interface{}
	correctValueExpectedResult   interface{}
	incorrectValue               interface{}
	incorrectValueExpectedResult error
}

type rulesTestData []ruleTestData

var rulesTestDataList = []ruleTestData{
	{
		ruleName:                     "required",
		correctValue:                 gofakeit.Name(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "email",
		correctValue:                 gofakeit.Email(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "test@mailcom",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "url",
		correctValue:                 gofakeit.URL(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "http:/githubcom/gocondor/gocondor",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "alpha",
		correctValue:                 gofakeit.LoremIpsumWord(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "test232",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "digit",
		correctValue:                 gofakeit.Digit(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "d",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "alphaNumeric",
		correctValue:                 "abc3",
		correctValueExpectedResult:   nil,
		incorrectValue:               "!",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "lowerCase",
		correctValue:                 "abc",
		correctValueExpectedResult:   nil,
		incorrectValue:               "ABC",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "upperCase",
		correctValue:                 "ABC",
		correctValueExpectedResult:   nil,
		incorrectValue:               "abc",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "int",
		correctValue:                 "343",
		correctValueExpectedResult:   nil,
		incorrectValue:               "342.3",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "float",
		correctValue:                 "433.5",
		correctValueExpectedResult:   nil,
		incorrectValue:               gofakeit.LoremIpsumWord(),
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "uuid",
		correctValue:                 uuid.NewString(),
		correctValueExpectedResult:   nil,
		incorrectValue:               gofakeit.LoremIpsumWord(),
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "creditCard",
		correctValue:                 "4242 4242 4242 4242",
		correctValueExpectedResult:   nil,
		incorrectValue:               "dd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "json",
		correctValue:                 "{\"testKEy\": \"testVal\"}",
		correctValueExpectedResult:   nil,
		incorrectValue:               "dd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "base64",
		correctValue:                 "+rxVsR0pD0DU4XO4MZbXXg==",
		correctValueExpectedResult:   nil,
		incorrectValue:               "dd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "countryCode2",
		correctValue:                 "SD",
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "countryCode3",
		correctValue:                 "SDN",
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "isoCurrencyCode",
		correctValue:                 "USD",
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "mac",
		correctValue:                 gofakeit.MacAddress(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "ip",
		correctValue:                 gofakeit.IPv4Address(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "ipv4",
		correctValue:                 gofakeit.IPv4Address(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "ipv6",
		correctValue:                 gofakeit.IPv6Address(),
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "domain",
		correctValue:                 "site.com",
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "latitude",
		correctValue:                 fmt.Sprintf("%v", gofakeit.Latitude()),
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "longitude",
		correctValue:                 fmt.Sprintf("%v", gofakeit.Longitude()),
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
	{
		ruleName:                     "semver",
		correctValue:                 "3.2.4",
		correctValueExpectedResult:   nil,
		incorrectValue:               "ddd",
		incorrectValueExpectedResult: errors.New("mock error"),
	},
}

func TestValidatorValidate(t *testing.T) {
	validator := newValidator()
	v := validator.Validate(ValidatorData{
		"name": gofakeit.LoremIpsumWord(),
		"link": gofakeit.URL(),
	},
		ValidatorRules{
			"name": "required|alphaNumeric",
			"link": "required|url",
		},
	)
	if v.Failed() {
		t.Errorf("erro testing validator validate: '%v'", v.GetErrorMessagesJson())
	}
	v = validator.Validate(ValidatorData{
		"name": "",
		"link": gofakeit.URL(),
	},
		ValidatorRules{
			"name": "required|alphaNumeric",
			"link": "required|url",
		},
	)
	if !v.Failed() {
		t.Errorf("erro testing validator validate")
	}
	msgsMap := v.GetErrorMessagesMap()
	if !strings.Contains(msgsMap["name"], "cannot be blank") {
		t.Errorf("erro testing validator validate")
	}
	msgsJson := v.GetErrorMessagesJson()
	var masgsMapOfJ map[string]interface{}
	json.Unmarshal([]byte(msgsJson), &masgsMapOfJ)
	sval, _ := masgsMapOfJ["name"].(string)
	if !strings.Contains(sval, "cannot be blank") {
		t.Errorf("erro testing validator validate")
	}
}

func TestValidatorParseRules(t *testing.T) {
	_, err := parseRules("344")
	if err == nil {
		t.Errorf("failed testing validation parse rules")
	}
	rules, err := parseRules("required| min: 3|length: 3, 5")
	if err != nil {
		t.Errorf("failed testing validation parse rules: '%v'", err.Error())
	}
	if len(rules) != 3 {
		t.Errorf("failed testing validation parse rules")
	}
}

func TestValidatorGetRule(t *testing.T) {
	for _, td := range rulesTestDataList {
		r, err := getRule(td.ruleName)
		if err != nil {
			t.Errorf("failed testing validation rule '%v'", td.ruleName)
		}
		err = r.Validate(td.correctValue)
		if err != nil {
			t.Errorf("failed testing validation rule '%v'", td.ruleName)
		}
		err = r.Validate(td.incorrectValue)
		if err == nil {
			t.Errorf("failed testing validation rule %v", td.ruleName)
		}
	}
	_, err := getRule("unknownrule")
	if err == nil {
		t.Errorf("failed testing validation rule")
	}
}

func TestValidatorGetRuleMax(t *testing.T) {
	r, err := getRuleMax("max: 33")
	if err != nil {
		t.Errorf("failed testing validation rule 'max'")
	}
	err = r.Validate(30)
	if err != nil {
		t.Errorf("failed testing validation rule 'max': %v", err.Error())
	}
	err = r.Validate(40)
	if err == nil {
		t.Errorf("failed testing validation rule 'max'")
	}
}

func TestValidatorGetRuleMin(t *testing.T) {
	r, err := getRuleMin("min: 33")
	if err != nil {
		t.Errorf("failed testing validation rule 'min'")
	}
	err = r.Validate(34)
	if err != nil {
		t.Errorf("failed testing validation rule 'min': %v", err.Error())
	}
	err = r.Validate(3)
	if err == nil {
		t.Errorf("failed testing validation rule 'min'")
	}
}

func TestValidatorGetRuleIn(t *testing.T) {
	r, err := getRuleIn("in: a, b, c")
	if err != nil {
		t.Errorf("failed testing validation rule 'in'")
	}
	err = r.Validate("a")
	if err != nil {
		t.Errorf("failed testing validation rule 'in': %v", err)
	}
}

func TestGetValidationRuleDateLayout(t *testing.T) {
	r, err := getRuleDateLayout("dateLayout: 02 January 2006")
	if err != nil {
		t.Errorf("failed testing validation rule 'dateLayout'")
	}
	err = r.Validate("02 May 2023")
	if err != nil {
		t.Errorf("failed testing validation rule 'dateLayout': %v", err.Error())
	}
	err = r.Validate("02-04-2023")
	if err == nil {
		t.Errorf("failed testing validation rule 'dateLayout'")
	}
}

func TestValidatorGetRuleLenth(t *testing.T) {
	r, err := getRuleLength("length: 3, 5")
	if err != nil {
		t.Errorf("failed test validation rule 'length': %v", err.Error())
	}
	err = r.Validate("123")
	if err != nil {
		t.Errorf("failed test validation rule 'length': %v", err.Error())
	}
	err = r.Validate("12")
	if err == nil {
		t.Errorf("failed test validation rule 'length'")
	}
	err = r.Validate("123456")
	if err == nil {
		t.Errorf("failed test validation rule 'length'")
	}

	r, err = getRuleLength("length: 3dd, 5")
	if err == nil {
		t.Errorf("failed test validation rule 'length'")
	}
	r, err = getRuleLength("length: 3, 5dd")
	if err == nil {
		t.Errorf("failed test validation rule 'length'")
	}
}
