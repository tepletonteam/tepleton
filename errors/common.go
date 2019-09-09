package errors

/**
*    Copyright (C) 2017 Ethan Frey
**/

import wrsp "github.com/tepleton/wrsp/types"

const (
	msgDecoding     = "Error decoding input"
	msgUnauthorized = "Unauthorized"
)

func DecodingError() TMError {
	return New(msgDecoding, wrsp.CodeType_EncodingError)
}

func Unauthorized() TMError {
	return New(msgUnauthorized, wrsp.CodeType_Unauthorized)
}
