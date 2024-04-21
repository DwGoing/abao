package eth

import (
	"abao/pkg/common"
	"abao/pkg/hd_wallet"
	"context"
	"errors"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	goEthereumCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	common.ChainClient
	currencies map[string]EthCurrency
}

// @title	创建Eth客户端
// @param	nodes		map[string]int			节点列表
// @param	currencies	map[string]TronCurrency	币种列表
// @return	_			*EthClient				Eth客户端
func NewEthClient(nodes map[string]int, currencies map[string]EthCurrency) *EthClient {
	return &EthClient{
		ChainClient: *common.NewChainClient(common.Chain_TRON, nodes),
		currencies:  currencies,
	}
}

// @title	获取当前高度
// @param	Self		*EthClient
// @return	_			int64			当前高度
// @return	_			error			异常信息
func (Self *EthClient) GetCurrentHeight() (int64, error) {
	client, err := Self.GetEthClient()
	if err != nil {
		return 0, err
	}
	height, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}
	return int64(height), nil
}

// @title	查询余额
// @param	Self		*EthClient
// @param	address		string		地址
// @param	currency	string		币种
// @param	args		any			参数
// @return	_			float64		余额
// @return	_			error		异常信息
func (Self *EthClient) GetBalance(address string, currency string, args any) (float64, error) {
	currencyInfo, ok := Self.currencies[currency]
	if !ok {
		return 0, errors.New("unsupported currency")
	}
	_, ok = args.(EthGetBalanceParameter)
	if !ok {
		return 0, nil
	}
	client, err := Self.GetEthClient()
	if err != nil {
		return 0, err
	}
	var balanceBigInt *big.Int
	if currencyInfo.Contract == "" {
		balanceBigInt, err = client.BalanceAt(context.Background(), goEthereumCommon.HexToAddress(address), nil)
		if err != nil {
			return 0, err
		}
	} else {
		erc20, err := NewErc20(goEthereumCommon.HexToAddress(currencyInfo.Contract), client)
		if err != nil {
			return 0, err
		}
		balanceBigInt, err = erc20.BalanceOf(nil, goEthereumCommon.HexToAddress(address))
		if err != nil {
			return 0, err
		}
	}
	balance, _ := new(big.Float).Quo(new(big.Float).SetInt(balanceBigInt), big.NewFloat(math.Pow10(currencyInfo.Decimals))).Float64()
	return balance, nil
}

// @title	转账
// @param	Self		*EthClient
// @param	privateKey	*ecdsa.PrivateKey	私钥
// @param	to			string				交易
// @param	currency	string				币种
// @param	value		float64				金额
// @param	args		any					参数
// @return	_			string				交易哈希
// @return	_			error				异常信息
func (Self *EthClient) Transfer(privateKey string, to string, currency string, value float64, args any) (string, error) {
	account, err := hd_wallet.NewAccountFromPrivateKeyHex(common.AddressType_ETH, privateKey)
	if err != nil {
		return "", err
	}
	currencyInfo, ok := Self.currencies[currency]
	if !ok {
		return "", errors.New("unsupported currency")
	}
	_, ok = args.(EthTransferParameter)
	if !ok {
		return "", nil
	}
	account.GetPrivateKey().ToECDSA()
	client, err := Self.GetEthClient()
	if err != nil {
		return "", err
	}
	var signedTx *types.Transaction
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return "", err
	}
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(account.GetPrivateKey().ToECDSA().PublicKey))
	if err != nil {
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	valueBigInt, _ := new(big.Float).Mul(big.NewFloat(value), big.NewFloat(math.Pow10(currencyInfo.Decimals))).Int(new(big.Int))
	if currencyInfo.Contract == "" {
		tx := types.NewTransaction(nonce, goEthereumCommon.HexToAddress(to), valueBigInt, 21000, gasPrice, nil)
		signedTx, err = types.SignTx(tx, types.LatestSignerForChainID(chainId), account.GetPrivateKey().ToECDSA())
		if err != nil {
			return "", err
		}
	} else {
		ierc20, err := NewErc20(goEthereumCommon.HexToAddress(currencyInfo.Contract), client)
		if err != nil {
			return "", err
		}
		transactOpts, err := bind.NewKeyedTransactorWithChainID(account.GetPrivateKey().ToECDSA(), chainId)
		if err != nil {
			return "", err
		}
		transactOpts.NoSend = true
		transactOpts.Nonce = big.NewInt(int64(nonce))
		transactOpts.GasLimit = uint64(300000)
		transactOpts.GasPrice = gasPrice
		signedTx, err = ierc20.Transfer(transactOpts, goEthereumCommon.HexToAddress(to), valueBigInt)
		if err != nil {
			return "", err
		}
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}

// @title	查询交易
// @param	Self		*EthClient
// @param	txHash		string			交易Hash
// @return	_			*Transaction	交易信息
// @return	_			error			异常信息
func (Self *EthClient) GetTransaction(txHash string) (*common.Transaction, error) {
	transaction := common.Transaction{
		Chain: common.Chain_ETH,
		Hash:  txHash,
	}
	client, err := Self.GetEthClient()
	if err != nil {
		return nil, err
	}
	receipt, err := client.TransactionReceipt(context.Background(), goEthereumCommon.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	transaction.Result = receipt.Status == 1
	transaction.Height = receipt.BlockNumber.Int64()
	tx, isPending, err := client.TransactionByHash(context.Background(), goEthereumCommon.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		return nil, err
	}
	transaction.TimeStamp = int64(block.Time())
	var currency string
	var currencyInfo EthCurrency
	var from string
	var to string
	var valueBigInt *big.Int
	if receipt.Logs == nil || len(receipt.Logs) == 0 {
		fromAddress, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
		if err != nil {
			return nil, err
		}
		from = fromAddress.Hex()
		to = tx.To().Hex()
		currency = "ETH"
		valueBigInt = tx.Value()
	} else {
		matchCurrency := linq.From(Self.currencies).FirstWithT(func(item linq.KeyValue) bool {
			return goEthereumCommon.HexToAddress(item.Value.(EthCurrency).Contract) == *tx.To()
		})
		currency = matchCurrency.(linq.KeyValue).Key.(string)
		currencyInfo = matchCurrency.(linq.KeyValue).Value.(EthCurrency)
		erc20, err := NewErc20(goEthereumCommon.HexToAddress(currencyInfo.Contract), client)
		if err != nil {
			return nil, err
		}
		for _, log := range receipt.Logs {
			filterer, err := erc20.ParseTransfer(*log)
			if err != nil {
				continue
			}
			from = filterer.From.Hex()
			to = filterer.To.Hex()
			valueBigInt = filterer.Value
			break
		}
	}
	transaction.Currency = currency
	transaction.From = from
	transaction.To = to
	transaction.Value, _ = new(big.Float).Quo(new(big.Float).SetInt(valueBigInt), big.NewFloat(math.Pow10(Self.currencies[currency].Decimals))).Float64()
	if !isPending && transaction.Result {
		height, err := Self.GetCurrentHeight()
		if err != nil {
			return nil, err
		}
		transaction.Confirms = height - transaction.Height
	}
	return &transaction, nil
}

// @title	获取Eth客户端
// @param 	Self 	*EthClient
// @return 	_ 		*ethclient.Client 	Eth客户端
// @return 	_ 		error 				异常信息
func (Self *EthClient) GetEthClient() (*ethclient.Client, error) {
	sum := 0
	for _, v := range Self.GetNodes() {
		sum += v
	}
	i := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(sum)
	var node string
	for k, v := range Self.GetNodes() {
		if v >= i {
			node = k
			break
		}
		i = i - v
	}
	client, err := ethclient.Dial(node)
	if err != nil {
		return nil, err
	}
	return client, nil
}
