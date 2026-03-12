package snip

import "errors"

var ErrCategoryExists = errors.New("category already exists")
var ErrEntryNotFound = errors.New("entry not found")
var ErrEntryExists = errors.New("entry already exists")
