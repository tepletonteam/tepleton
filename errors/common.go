package errors

/**
*    Copyright (C) 2017 Ethan Frey
**/

import wrsp "github.com/tepleton/wrsp/types"

const (
	msgDecoding          = "Error decoding input"
	msgUnauthorized      = "Unauthorized"
	msgInvalidAddress    = "Invalid Address"
	msgInvalidCoins      = "Invalid Coins"
	msgInvalidSequence   = "Invalid Sequence"
	msgInvalidSignature  = "Invalid Signature"
	msgNoInputs          = "No Input Coins"
	msgNoOutputs         = "No Output Coins"
	msgTooLarge          = "Input size too large"
	msgMissingSignature  = "Signature missing"
	msgTooManySignatures = "Too many signatures"
)

func DecodingError() TMError {
	return New(msgDecoding, wrsp.CodeType_EncodingError)
}

func Unauthorized() TMError {
	return New(msgUnauthorized, wrsp.CodeType_Unauthorized)
}

func MissingSignature() TMError {
	return New(msgMissingSignature, wrsp.CodeType_Unauthorized)
}

func TooManySignatures() TMError {
	return New(msgTooManySignatures, wrsp.CodeType_Unauthorized)
}

func InvalidSignature() TMError {
	return New(msgInvalidSignature, wrsp.CodeType_Unauthorized)
}

func InvalidAddress() TMError {
	return New(msgInvalidAddress, wrsp.CodeType_BaseInvalidInput)
}

func InvalidCoins() TMError {
	return New(msgInvalidCoins, wrsp.CodeType_BaseInvalidInput)
}

func InvalidSequence() TMError {
	return New(msgInvalidSequence, wrsp.CodeType_BaseInvalidInput)
}

func NoInputs() TMError {
	return New(msgNoInputs, wrsp.CodeType_BaseInvalidInput)
}

func NoOutputs() TMError {
	return New(msgNoOutputs, wrsp.CodeType_BaseInvalidOutput)
}

func TooLarge() TMError {
	return New(msgTooLarge, wrsp.CodeType_EncodingError)
}
