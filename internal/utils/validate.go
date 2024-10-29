package utils

import "github.com/go-playground/validator/v10"

// Validate is a globally accessible instance of the validator package, used for struct validation.
var Validate = validator.New()
