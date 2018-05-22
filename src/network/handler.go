package network

import (
	"api"
	"logic"
	"net/http"
)

type fn func(string, http.ResponseWriter, []byte)

var PROCCESSING_MAP = map[string]fn{
	api.BALANCE:       logic.Do_GetBalance,
	api.DEPOSIT:       logic.Do_Deposit,
	api.WITHDRAW:      logic.Do_Withdraw,
	api.WITHDRAW_MAX:  logic.Do_WithdrawMax,
	api.REGIS_APP:     logic.Do_RegisApp,
	api.CHANGE_IP:     logic.Do_ChangeIP,
	api.CHANGE_NOTIFY: logic.Do_ChangeNofify,
	api.TRANSFER:      logic.Do_Transfer,
	api.LIST_ADDRESS:  logic.Do_ListAddress,
}
