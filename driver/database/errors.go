package database

import "errors"

var ErrMissingOptionConnection = errors.New("missing option \"connection\"")
var ErrMissingOptionName = errors.New("missing option \"name\"")
