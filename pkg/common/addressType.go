package common

import "errors"

type AddressType int8

const (
	AddressType_BTC_LEGACY        AddressType = 1
	AddressType_BTC_SEGWIT        AddressType = 2
	AddressType_BTC_NATIVE_SEGWIT AddressType = 3
	AddressType_ETH               AddressType = 4
	AddressType_TRON              AddressType = 5
	AddressType_BSC               AddressType = 6
)

func (AddressType AddressType) ToString() (string, error) {
	switch AddressType {
	case AddressType_BTC_LEGACY:
		return "BTC Legacy", nil
	case AddressType_BTC_SEGWIT:
		return "BTC SegWit", nil
	case AddressType_BTC_NATIVE_SEGWIT:
		return "BTC Native SegWit", nil
	case AddressType_ETH:
		return "ETH", nil
	case AddressType_TRON:
		return "TRON", nil
	case AddressType_BSC:
		return "BSC", nil
	default:
		return "", errors.New("unsupported address type")
	}
}

func ParseAddressType(AddressType string) (AddressType, error) {
	switch AddressType {
	case "BTC Legacy":
		return AddressType_BTC_LEGACY, nil
	case "BTC SegWit":
		return AddressType_BTC_SEGWIT, nil
	case "BTC Native SegWit":
		return AddressType_BTC_NATIVE_SEGWIT, nil
	case "ETH":
		return AddressType_ETH, nil
	case "TRON":
		return AddressType_TRON, nil
	case "BSC":
		return AddressType_BSC, nil
	default:
		return 0, errors.New("unsupported address type")
	}
}
