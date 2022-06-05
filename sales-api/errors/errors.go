package errors

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/go-playground/validator/v10"
)

const (
	validationErrTypePOST = "any.required"
)

type BodyValidationError struct {
	Ctx     ValidationErrorCtx `json:"context"`
	Message string             `json:"message"`
	Type    string             `json:"type"`
	Path    []string           `json:"path"`
}

type ValidationErrorCtx struct {
	Value           *struct{} `json:"value,omitempty"`
	Label           string    `json:"label"`
	Key             string    `json:"key,omitempty"`
	Peers           []string  `json:"peers,omitempty"`
	PeersWithLabels []string  `json:"peersWithLabels,omitempty"`
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
			Value:           &struct{}{},
		},
	})

	return errs, msg
}
