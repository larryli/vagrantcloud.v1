package vagrantcloud

import (
	"encoding/json"
	"errors"
)

// Standard HTTP response codes are returned.
// 404 Not Found codes are returned for all resources that a user does not have access to,
// as well as for resources that don't exist.
// This is done to avoid a potential attacker discovery the existence of a resource.
//
// Errors format:
//
// 		{
// 			"name": [
// 				"has already been taken"
// 			]
// 		}
type Error struct {
	Msg    string
	Errors interface{}
}

func NewError(msg string, errs string) error {
	var info map[string]interface{}
	err := json.Unmarshal([]byte(errs), &info)
	if err == nil {
		if errs, ok1 := info["errors"]; ok1 {
			return &Error{
				Msg:    msg,
				Errors: errs,
			}
		}
	}
	return errors.New(msg)
}

func (e *Error) Error() string {
	return e.Msg
}
