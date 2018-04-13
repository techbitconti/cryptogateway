package btc

import (
	"config"
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcutil"
	pb "github.com/btcsuite/btcwallet/rpc/walletrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var certificateFile = filepath.Join(config.PATH_BTC, "rpc.cert")

var c pb.WalletServiceClient

func Connect_btcwallet_gRPC(net string) {

	ip := "localhost"
	port := "18554"

	switch net {
	case "mainnet":
		port = "8332"
	case "testnet":
		port = "18332"
	case "simnet":
		port = "18554"
	}
	host := ip + ":" + port

	creds, err := credentials.NewClientTLSFromFile(certificateFile, ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer conn.Close()
	c = pb.NewWalletServiceClient(conn)
}

func Ping_gRPC() {
	pingRequest := &pb.PingRequest{}
	pingResponse, err := c.Ping(context.Background(), pingRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Ping : ", pingResponse)
}

func Network_gRPC() {
	netWorkRequest := &pb.NetworkRequest{}
	netWorkRespone, err := c.Network(context.Background(), netWorkRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Network : ", netWorkRespone)
}

func AccountNumber_gRPC(accountName string) {
	accountNumberRequest := &pb.AccountNumberRequest{AccountName: accountName}
	accountNumberRespone, err := c.AccountNumber(context.Background(), accountNumberRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("AccountNumber : ", accountNumberRespone)
}

func Accounts_gRPC() {
	accountsRequest := &pb.AccountsRequest{}
	accountsRespone, err := c.Accounts(context.Background(), accountsRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Accounts : ", accountsRespone.GetAccounts())
}

func Balance_gRPC(indexAccount uint32, confirm int32) {
	balanceRequest := &pb.BalanceRequest{
		AccountNumber:         indexAccount,
		RequiredConfirmations: confirm,
	}
	balanceResponse, err := c.Balance(context.Background(), balanceRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Balance ", indexAccount, " : ", btcutil.Amount(balanceResponse.Spendable))
}

func GetTransactions_gRPC(startBlock, endBlock int32) {
	getTransRequest := &pb.GetTransactionsRequest{
		StartingBlockHeight: startBlock,
		EndingBlockHeight:   endBlock,
	}
	getTransRespone, err := c.GetTransactions(context.Background(), getTransRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("GetTransactions : ", getTransRespone.GetMinedTransactions())
}

func NextAccount_gRPC(accountName, passPhrase string) {
	nextAccountRequest := &pb.NextAccountRequest{
		Passphrase:  []byte(passPhrase),
		AccountName: accountName,
	}
	nextAccountRespone, err := c.NextAccount(context.Background(), nextAccountRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("NextAccount : ", nextAccountRespone)
}

func NextAddress_gRPC(indexAccount uint32) {

	nextAddressRequest := &pb.NextAddressRequest{Account: indexAccount}
	nextAddressRespone, err := c.NextAddress(context.Background(), nextAddressRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("NextAddress : ", nextAddressRespone)
}

//c.ImportPrivateKey()

//c.FundTransaction()

//c.SignTransaction()

//c.PublishTransaction()
