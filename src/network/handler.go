package network

import (
	"api"
	"logic"
	"net/http"
)

type fn func(string, http.ResponseWriter, []byte)

var PROCCESSING_MAP = map[string]fn{
	api.ADDRESS: logic.Do_GenAddress,
	api.BALANCE: logic.Do_GetBalance,
}
