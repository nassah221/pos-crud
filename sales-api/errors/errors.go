package errors

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/go-playground/validator/v10"
)

const (
	validationErrTypePOST = "any.required"
	validationErrTypePUT  = "object.missing"
)

type BodyValidationError struct {
	Message string             `json:"message"`
	Path    []string           `json:"path"`
	Type    string             `json:"type"`
	Ctx     ValidationErrorCtx `json:"context"`
}

type ValidationErrorCtx struct {
	Label           string   `json:"label"`
	Key             string   `json:"key,omitempty"`
	Peers           []string `json:"peers,omitempty"`
	PeersWithLabels []string `json:"peersWithLabels,omitempty"`
	Value           struct{} `json:"value,omitempty"`
}

func FromFieldValidationErrorPOST(err error) ([]BodyValidationError, string) {
	errs := []BodyValidationError{}
	msg := "body ValidationError:"

	for _, fieldErr := range err.(validator.ValidationErrors) {
		fieldName := strcase.ToCamel(fieldErr.Field())
		m := fmt.Sprintf("%q is required", fieldName)
		msg += " " + m + "."
		vErr := BodyValidationError{
			Message: fmt.Sprintf("%q is required", fieldName),
			Path:    []string{fieldName},
			Type:    validationErrTypePOST,
			Ctx: ValidationErrorCtx{
				Label: fieldName,
				Key:   fieldName,
			},
		}
		errs = append(errs, vErr)
	}

	return errs, msg
}

func FromFieldValidationErrorPUT(err error) ([]BodyValidationError, string) {
	errs := []BodyValidationError{}
	msg := "body ValidationError: \"value\" must contain at least one of "

	missingVals := []string{}
	for _, fieldErr := range err.(validator.ValidationErrors) {
		fieldName := strcase.ToCamel(fieldErr.Field())
		missingVals = append(missingVals, fieldName)
	}

	msg += "[" + strings.Join(missingVals, ", ") + "]"
	errs = append(errs, BodyValidationError{
		Message: msg,
		Path:    []string{},
		Ctx: ValidationErrorCtx{
			Label:           "value",
			Peers:           missingVals,
			PeersWithLabels: missingVals,
		},
	})

	return errs, msg
}
