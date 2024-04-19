package hd_wallet

import (
	"abao/pkg/common"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/btcutil/bech32"
	"github.com/ethereum/go-ethereum/crypto"
	tronCommon "github.com/fbsobreira/gotron-sdk/pkg/common"
)

type Account struct {
	privateKey *btcec.PrivateKey
}

// @title	创建账户
// @param	privateKeyHex	string		节点列表
// @return	_				*Account	账户
// @return	_				error		异常信息
func NewAccountFromPrivateKey(privateKey *btcec.PrivateKey) (*Account, error) {
	return &Account{
		privateKey: privateKey,
	}, nil
}

// @title	创建账户
// @param	privateKeyHex	string		节点列表
// @return	_				*Account	账户
// @return	_				error		异常信息
func NewAccountFromPrivateKeyHex(privateKeyHex string) (*Account, error) {
	bytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), bytes)
	return NewAccountFromPrivateKey(privateKey)
}

/*
@title 	获取钱包地址
@param 	Self   		*Account
@param 	addressType common.AddressType	地址类型
@return _ 			string 				地址
@return _ 			error 				异常信息
*/
func (Self *Account) GetAddress(addressType common.AddressType) (string, error) {
	address := ""
	switch addressType {
	case common.AddressType_BTC_LEGACY:
		bytes := btcutil.Hash160(Self.privateKey.PubKey().SerializeCompressed())
		address = base58.CheckEncode(bytes, 0x00)
	case common.AddressType_BTC_SEGWIT:
		bytes := btcutil.Hash160(Self.privateKey.PubKey().SerializeCompressed())
		bytes = append([]byte{0x00, 0x14}, bytes...)
		bytes = btcutil.Hash160(bytes)
		address = base58.CheckEncode(bytes, 0x05)
	case common.AddressType_BTC_NATIVE_SEGWIT:
		bytes := btcutil.Hash160(Self.privateKey.PubKey().SerializeCompressed())
		converted, err := bech32.ConvertBits(bytes, 8, 5, true)
		if err != nil {
			break
		}
		combined := make([]byte, len(converted)+1)
		combined[0] = 0x00
		copy(combined[1:], converted)
		address, err = bech32.Encode("bc", combined)
		if err != nil {
			break
		}
	case common.AddressType_ETH:
		address = crypto.PubkeyToAddress(Self.privateKey.ToECDSA().PublicKey).Hex()
	case common.AddressType_TRON:
		ethAddress := crypto.PubkeyToAddress(Self.privateKey.ToECDSA().PublicKey)
		tronAddress := make([]byte, 0)
		tronAddress = append(tronAddress, byte(0x41))
		tronAddress = append(tronAddress, ethAddress.Bytes()...)
		address = tronCommon.EncodeCheck(tronAddress)
	case common.AddressType_BSC:
		address = crypto.PubkeyToAddress(Self.privateKey.ToECDSA().PublicKey).Hex()
	default:
		return "", errors.New("unsupported address type")
	}
	return address, nil
}
