package config

var PATH_ETH = "/Users/A/ethereum/private/keystore"
var PATH_BTC = "/Users/A/bitcoin/btcwallet"

func SetPATH(net string) {
	switch net {
	case "local":
		PATH_ETH = "/Users/A/ethereum/private/keystore"
		PATH_BTC = "/Users/A/bitcoin/btcwallet"
	case "server":
		PATH_ETH = "/home/ramost/ethereum/private/keystore"
		PATH_BTC = "/home//bitcoin/btcwallet"
	}
}

var BTC_SIM = struct {
	PrivKey string
	Address string
}{
	`FwMCdKjGEMYe1VPL2tXEqH7ecXhZXshqzpBoxQieCXDG5yGQvGuZ`,
	`ShGn19iCHHhGjy7huZeQ2T8vtBo7DFpA9U`,
}

var BTC_TEST = struct {
	PrivKey string
	Address string
}{
	`cQJYynSnzuUbNisDb7FsM2tpKi7Hu3HKtxegWohemwf8YU1EDduD`,
	`2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh`,
}

var ETH_SIM = struct {
	PrivKey string
	Address string
}{
	`47de15108b35169c4aff4826d5c413fe117e361a900325f6d3df1f0e04cbd706`,
	`0x8dD75F7c03A048C0a66a53dbf9ED76d04E9a9eA3`,
}

var ETH_TEST = struct {
	PrivKey string
	Address string
}{
	`47de15108b35169c4aff4826d5c413fe117e361a900325f6d3df1f0e04cbd706`,
	`0x8dD75F7c03A048C0a66a53dbf9ED76d04E9a9eA3`,
}

var ERC20_SIM = struct {
	PrivKey string
	Address string
	Abi     string
}{
	ETH_SIM.PrivKey,
	`0x1da98ecccd7fca0e38d8b0732b53a1ce6a382ce7`,
	`[
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "spender",
				"type": "address"
			},
			{
				"name": "tokens",
				"type": "uint256"
			}
		],
		"name": "approve",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "from",
				"type": "address"
			},
			{
				"name": "to",
				"type": "address"
			},
			{
				"name": "tokens",
				"type": "uint256"
			}
		],
		"name": "transferFrom",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [
			{
				"name": "",
				"type": "uint8"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_totalSupply",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenOwner",
				"type": "address"
			}
		],
		"name": "balanceOf",
		"outputs": [
			{
				"name": "balance",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "acceptOwnership",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "to",
				"type": "address"
			},
			{
				"name": "tokens",
				"type": "uint256"
			}
		],
		"name": "transfer",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "spender",
				"type": "address"
			},
			{
				"name": "tokens",
				"type": "uint256"
			},
			{
				"name": "data",
				"type": "bytes"
			}
		],
		"name": "approveAndCall",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "newOwner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "tokenAddress",
				"type": "address"
			},
			{
				"name": "tokens",
				"type": "uint256"
			}
		],
		"name": "transferAnyERC20Token",
		"outputs": [
			{
				"name": "success",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "tokenOwner",
				"type": "address"
			},
			{
				"name": "spender",
				"type": "address"
			}
		],
		"name": "allowance",
		"outputs": [
			{
				"name": "remaining",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_newOwner",
				"type": "address"
			}
		],
		"name": "transferOwnership",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"payable": true,
		"stateMutability": "payable",
		"type": "fallback"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_from",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "_to",
				"type": "address"
			}
		],
		"name": "OwnershipTransferred",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "from",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "to",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "tokens",
				"type": "uint256"
			}
		],
		"name": "Transfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "tokenOwner",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "spender",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "tokens",
				"type": "uint256"
			}
		],
		"name": "Approval",
		"type": "event"
	}
]`,
}

var ERC20_TEST = struct {
	PrivKey string
	Address string
	Abi     string
}{
	ETH_TEST.PrivKey,
	``,
	``,
}

var MASS = struct {
	PrivKey string
	Address string
	Abi     string
}{
	ETH_TEST.PrivKey,
	`0x25a12cda4cdd61bff3be3b1e585dc3c5651be217`,
	`[
	{
		"constant": false,
		"inputs": [
			{
				"name": "_a",
				"type": "string"
			}
		],
		"name": "stringToAddress",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "tmp",
				"type": "bytes"
			}
		],
		"name": "bytesToAddress",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_arr",
				"type": "string"
			},
			{
				"name": "num",
				"type": "uint256"
			}
		],
		"name": "stringToListAddress",
		"outputs": [
			{
				"name": "",
				"type": "address[]"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_addrDeposit",
				"type": "address"
			},
			{
				"name": "_receipts",
				"type": "address[]"
			},
			{
				"name": "_tokens",
				"type": "uint256[]"
			}
		],
		"name": "uploadETH",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_addrDeposit",
				"type": "address"
			},
			{
				"name": "_erc20",
				"type": "address[]"
			},
			{
				"name": "_receipts",
				"type": "address[]"
			},
			{
				"name": "_tokens",
				"type": "uint256[]"
			}
		],
		"name": "uploadERC20",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_addrDeposit",
				"type": "address"
			}
		],
		"name": "transferERC20",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "newOwner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_addrDeposit",
				"type": "address"
			}
		],
		"name": "transferETH",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_a",
				"type": "string"
			}
		],
		"name": "stringToInt",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": true,
		"stateMutability": "payable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "_addrDeposit",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_receipts",
				"type": "address[]"
			},
			{
				"indexed": false,
				"name": "_tokens",
				"type": "uint256[]"
			}
		],
		"name": "LogUpLoadETH",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "_addrDeposit",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_erc20",
				"type": "address[]"
			},
			{
				"indexed": false,
				"name": "_receipts",
				"type": "address[]"
			},
			{
				"indexed": false,
				"name": "_tokens",
				"type": "uint256[]"
			}
		],
		"name": "LogUpLoadERC20",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "_ok",
				"type": "bool"
			},
			{
				"indexed": false,
				"name": "_addrDeposit",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_receipts",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_coin",
				"type": "uint256"
			}
		],
		"name": "LogTranferETH",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "_addrDeposit",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_erc20",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_receipts",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_tokens",
				"type": "uint256"
			}
		],
		"name": "LogTranferERC20",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "_from",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "_to",
				"type": "address"
			}
		],
		"name": "OwnershipTransferred",
		"type": "event"
	}
]`,
}
