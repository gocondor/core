package core

import (
	"encoding/json"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
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
	// TODO map rules
	vr = validationResult{}
	vr.hasFailed = false
	res := map[string]string{}
	for key, val := range data {
		err := validation.Validate(val, validation.Required)
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
