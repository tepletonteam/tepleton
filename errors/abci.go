package errors

import (
	wrsp "github.com/tepleton/wrsp/types"
)

type causer interface {
	Cause() error
}

func getWRSPError(err error) (WRSPError, bool) {
	if err, ok := err.(WRSPError); ok {
		return err, true
	}
	if causer, ok := err.(causer); ok {
		err := causer.Cause()
		if err, ok := err.(WRSPError); ok {
			return err, true
		}
	}
	return nil, false
}

func ResponseDeliverTxFromErr(err error) *wrsp.ResponseDeliverTx {
	var code = CodeInternalError
	var log = codeToDefaultLog(code)

	wrspErr, ok := getWRSPError(err)
	if ok {
		code = wrspErr.WRSPCode()
		log = wrspErr.WRSPLog()
	}

	return &wrsp.ResponseDeliverTx{
		Code: code,
		Data: nil,
		Log:  log,
		Tags: nil,
	}
}

func ResponseCheckTxFromErr(err error) *wrsp.ResponseCheckTx {
	var code = CodeInternalError
	var log = codeToDefaultLog(code)

	wrspErr, ok := getWRSPError(err)
	if ok {
		code = wrspErr.WRSPCode()
		log = wrspErr.WRSPLog()
	}

	return &wrsp.ResponseCheckTx{
		Code: code,
		Data: nil,
		Log:  log,
		Gas:  0, // TODO
		Fee:  0, // TODO
	}
}
