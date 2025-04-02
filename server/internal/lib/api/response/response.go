package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be of the email type", err.Field()))
		// case "min":
		//  	errMsgs = append(errMsgs, fmt.Sprintf("field %s must have more than %s characters", err.Field(), err.Param()))
		// case "gte":
		// 	errMsgs = append(errMsgs, fmt.Sprintf("field %s must have more than %s characters", err.Field(), err.Param()))
		case "gte":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must have at least %s characters", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: http.StatusConflict,
		Error:  strings.Join(errMsgs, ", "),
	}
}
