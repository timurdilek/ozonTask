package service

import "errors"

var (
	ErrIncorrectPostLen    = errors.New("too long post")
	ErrIncorrectCommentLen = errors.New("too long comment")
	ErrIncorrectContentLen = errors.New("incorrect content")
)
