package autoEncode

import "errors"

var (
	ErrFileNotFound       = errors.New("file is not found")
	ErrTargetNotFound     = errors.New("target title not found")
	ErrAlreadyExists      = errors.New("target is already exists")
	ErrStatusUnchanged    = errors.New("target status is already updated")
	ErrZeroRecord         = errors.New("status record is zero")
	ErrException          = errors.New("exception err")
	ErrTargetIsNotDir     = errors.New("target is not directory")
	ErrFailedCMD          = errors.New("cmd is failed")
	ErrTargetPathNotFound = errors.New("target path is not found")
	ErrAddAmatsukaze      = errors.New("Add amatsukaze failed")
)
