package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Validator struct{}
type validationResult struct {
	hasFailed     bool
	errorMessages map[string]string
}
type ValidatorData map[string]interface{}
type ValidatorRules map[string]interface{}

var vr validationResult
var v *Validator

func newValidator() *Validator {
	v := &Validator{}
	return v
}

func resolve() *Validator {
	return v
}

func (v *Validator) Validate(data ValidatorData, rules ValidatorRules) validationResult {
	vr = validationResult{}
	vr.hasFailed = false
	res := map[string]string{}
	for key, val := range data {
		_, ok := rules[key]
		if !ok {
			continue
		}
		err := validation.Validate(val, parseRules(rules[key])...)
		if err != nil {
			res[key] = fmt.Sprintf("%v: %v", key, err.Error())

		}
	}
	if len(res) != 0 {
		vr.hasFailed = true
		vr.errorMessages = res
	}
	return vr
}

func (vr *validationResult) Failed() bool {
	return vr.hasFailed
}

func (vr *validationResult) GetErrorMessagesMap() map[string]string {
	return vr.errorMessages
}

func (vr *validationResult) GetErrorMessagesJson() string {
	j, err := json.Marshal(vr.GetErrorMessagesMap())
	if err != nil {
		panic("error converting validation error messages to json")
	}
	return string(j)
}

func parseRules(rawRules interface{}) []validation.Rule {
	var res []validation.Rule
	rulesStr, ok := rawRules.(string)
	if !ok {
		panic("invalid validation rule")
	}
	rules := strings.Split(rulesStr, "|")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		res = append(res, getRule(rule))
	}

	return res
}

// TODO handle all rules
func getRule(rule string) validation.Rule {
	if strings.Contains(rule, "max:") {
		// max: 44
		rr := strings.ReplaceAll(rule, "max:", "")
		m := strings.TrimSpace(rr)
		n, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			panic("invalid value for validation rule 'max'")
		}
		return validation.Max(n)
	}

	if strings.Contains(rule, "min:") {
		// min: 33
		rr := strings.ReplaceAll(rule, "min:", "")
		m := strings.TrimSpace(rr)
		n, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			panic("invalid value for validation rule 'min'")
		}
		return validation.Min(n)
	}

	if strings.Contains(rule, "in:") {
		// in: first, second, third
		var readyElms []interface{}
		rr := strings.ReplaceAll(rule, "in:", "")
		elms := strings.Split(rr, ",")
		for _, elm := range elms {
			readyElms = append(readyElms, strings.TrimSpace(elm))
		}
		return validation.In(readyElms...)
	}

	//https://programming.guide/go/format-parse-string-time-date-example.html
	if strings.Contains(rule, "dateLayout:") {
		// dateLayout: 02 January 2006
		rr := rule
		rr = strings.TrimSpace(strings.Replace(rr, "dateLayout:", "", -1))
		return validation.Date(rr)
	}

	if strings.Contains(rule, "length:") {
		// length: 3, 7
		rr := rule
		rr = strings.Replace(rr, "length:", "", -1)
		lengthRange := strings.Split(rr, ",")
		if len(lengthRange) < 0 {
			panic("min value is not set for validation rule 'length'")
		}
		min, err := strconv.Atoi(strings.TrimSpace(lengthRange[0]))
		if err != nil {
			panic("min value is not set for validation rule 'length'")
		}
		if len(lengthRange) < 1 {
			panic("max value is not set for validation rule 'length'")
		}
		max, err := strconv.Atoi(strings.TrimSpace(lengthRange[1]))
		if err != nil {
			panic("max value is not set for validation rule 'length'")
		}
		return validation.Length(min, max)
	}

	switch rule {
	case "required":
		return validation.Required
	case "email":
		return is.Email
	case "url":
		return is.URL
	case "alpha":
		return is.Alpha
	case "digit":
		return is.Digit
	case "alphaNumeric":
		return is.Alphanumeric
	case "lowerCase":
		return is.LowerCase
	case "upperCase":
		return is.UpperCase
	case "int":
		return is.Int
	case "float":
		return is.Float
	case "uuid":
		return is.UUID
	case "creditCard":
		return is.CreditCard
	case "json":
		return is.JSON
	case "base64":
		return is.Base64
	case "countryCode2":
		return is.CountryCode2
	case "countryCode3":
		return is.CountryCode3
	case "isoCurrencyCode":
		return is.CurrencyCode
	case "mac":
		return is.MAC
	case "ip":
		return is.IP
	case "ipv4":
		return is.IPv4
	case "ipv6":
		return is.IPv6
	case "subdomain":
		return is.Subdomain
	case "domain":
		return is.Domain
	case "dnsName":
		return is.DNSName
	case "host":
		return is.Host
	case "port":
		return is.Port
	case "mongoDbId":
		return is.MongoID
	case "latitude":
		return is.Latitude
	case "longitude":
		return is.Longitude
	case "ssn":
		return is.SSN
	case "semver":
		return is.Semver
	default:
		panic(fmt.Sprintf("invalid validation rule: %v", rule))
	}
}
