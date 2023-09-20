package core

import (
	"encoding/json"
	"errors"
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

var vr validationResult
var v *Validator

func newValidator() *Validator {
	v := &Validator{}
	return v
}

func resolveValidator() *Validator {
	return v
}

func (v *Validator) Validate(data map[string]interface{}, rules map[string]interface{}) validationResult {
	vr = validationResult{}
	vr.hasFailed = false
	res := map[string]string{}
	for key, val := range data {
		_, ok := rules[key]
		if !ok {
			continue
		}
		rls, err := parseRules(rules[key])
		if err != nil {
			panic(err.Error())
		}
		err = validation.Validate(val, rls...)
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

func parseRules(rawRules interface{}) ([]validation.Rule, error) {
	var res []validation.Rule
	rulesStr, ok := rawRules.(string)
	if !ok {
		return nil, errors.New("invalid validation rule")
	}
	rules := strings.Split(rulesStr, "|")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		r, err := getRule(rule)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func getRule(rule string) (validation.Rule, error) {
	switch {
	case strings.Contains(rule, "max:"):
		return getRuleMax(rule)
	case strings.Contains(rule, "min:"):
		return getRuleMin(rule)
	case strings.Contains(rule, "in:"):
		return getRuleIn(rule)
	case strings.Contains(rule, "dateLayout:"):
		return getRuleDateLayout(rule)
	case strings.Contains(rule, "length:"):
		return getRuleLength(rule)
	}

	switch rule {
	case "required":
		return validation.Required, nil
	case "email":
		return is.Email, nil
	case "url":
		return is.URL, nil
	case "alpha":
		return is.Alpha, nil
	case "digit":
		return is.Digit, nil
	case "alphaNumeric":
		return is.Alphanumeric, nil
	case "lowerCase":
		return is.LowerCase, nil
	case "upperCase":
		return is.UpperCase, nil
	case "int":
		return is.Int, nil
	case "float":
		return is.Float, nil
	case "uuid":
		return is.UUID, nil
	case "creditCard":
		return is.CreditCard, nil
	case "json":
		return is.JSON, nil
	case "base64":
		return is.Base64, nil
	case "countryCode2":
		return is.CountryCode2, nil
	case "countryCode3":
		return is.CountryCode3, nil
	case "isoCurrencyCode":
		return is.CurrencyCode, nil
	case "mac":
		return is.MAC, nil
	case "ip":
		return is.IP, nil
	case "ipv4":
		return is.IPv4, nil
	case "ipv6":
		return is.IPv6, nil
	case "subdomain":
		return is.Subdomain, nil
	case "domain":
		return is.Domain, nil
	case "dnsName":
		return is.DNSName, nil
	case "host":
		return is.Host, nil
	case "port":
		return is.Port, nil
	case "mongoID":
		return is.MongoID, nil
	case "latitude":
		return is.Latitude, nil
	case "longitude":
		return is.Longitude, nil
	case "ssn":
		return is.SSN, nil
	case "semver":
		return is.Semver, nil
	default:
		err := errors.New(fmt.Sprintf("invalid validation rule: %v", rule))
		return nil, err
	}
}

func getRuleMax(rule string) (validation.Rule, error) {
	// max: 44
	rr := strings.ReplaceAll(rule, "max:", "")
	m := strings.TrimSpace(rr)
	n, err := strconv.ParseInt(m, 10, 64)
	if err != nil {
		err := errors.New("invalid value for validation rule 'max'")
		return nil, err
	}
	return validation.Max(n), err
}

func getRuleMin(rule string) (validation.Rule, error) {
	// min: 33
	rr := strings.ReplaceAll(rule, "min:", "")
	m := strings.TrimSpace(rr)
	n, err := strconv.ParseInt(m, 10, 64)
	if err != nil {
		err := errors.New("invalid value for validation rule 'min'")
		return nil, err
	}
	return validation.Min(n), nil
}

func getRuleIn(rule string) (validation.Rule, error) {
	// in: first, second, third
	var readyElms []interface{}
	rr := strings.ReplaceAll(rule, "in:", "")
	elms := strings.Split(rr, ",")
	for _, elm := range elms {
		readyElms = append(readyElms, strings.TrimSpace(elm))
	}
	return validation.In(readyElms...), nil
}

// example date layouts: https://programming.guide/go/format-parse-string-time-date-example.html
func getRuleDateLayout(rule string) (validation.Rule, error) {
	// dateLayout: 02 January 2006
	rr := rule
	rr = strings.TrimSpace(strings.Replace(rr, "dateLayout:", "", -1))
	return validation.Date(rr), nil
}

func getRuleLength(rule string) (validation.Rule, error) {
	// length: 3, 7
	rr := rule
	rr = strings.Replace(rr, "length:", "", -1)
	lengthRange := strings.Split(rr, ",")
	if len(lengthRange) < 0 {
		err := errors.New("min value is not set for validation rule 'length'")
		return nil, err
	}
	min, err := strconv.Atoi(strings.TrimSpace(lengthRange[0]))
	if err != nil {
		err := errors.New("min value is not set for validation rule 'length'")
		return nil, err
	}
	if len(lengthRange) < 1 {
		err := errors.New("max value is not set for validation rule 'length'")
		return nil, err
	}
	max, err := strconv.Atoi(strings.TrimSpace(lengthRange[1]))
	if err != nil {
		err := errors.New("max value is not set for validation rule 'length'")
		return nil, err
	}
	return validation.Length(min, max), nil
}
