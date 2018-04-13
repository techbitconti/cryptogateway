package network

import (
	"api"
	"logic"
	"net/http"
)

type fn func(string, http.ResponseWriter, []byte)

var PROCCESSING_MAP = map[string]fn{
	api.BALANCE:  logic.Do_GetBalance,
	api.DEPOSIT:  logic.Do_Deposit,
	api.WITHDRAW: logic.Do_Withdraw,
}
