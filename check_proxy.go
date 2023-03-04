package main

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 특정 문자를 문자열에서 모두 찾기
func findAllOccurrences(s, substr string) []string {
	var occurrences []string
	for i := 0; i < len(s); {
		index := strings.Index(s[i:], substr)
		if index == -1 {
			break
		}
		result := i + index + 2
		if result+64 > len(s) {
			break
		}
		mes := s[result : result+64]
		occurrences = append(occurrences, mes)
		i += index + len(substr)
	}
	return occurrences
}

// 슬라이스에서 중복된 요소 제거
func makeSliceUnique(s []string) []string {
	keys := make(map[string]struct{})
	res := make([]string, 0)
	for _, val := range s {
		if _, ok := keys[val]; ok {
			continue
		} else {
			keys[val] = struct{}{}
			res = append(res, val)
		}
	}
	return res
}
func is_TRN(ctx context.Context, ethclient *ethclient.Client, account common.Address, code []byte, encodeCode string) string {
	// Logic contract lookup portion of Transparent proxy contract
	data, _ := hexutil.Decode("0x5c60da1b") // First 4 bytes of keccak256("implementation()")
	callMsg := ethereum.CallMsg{            // callMsg settings
		To:   &account,
		Data: data,
	}@
	response, _ := ethclient.CallContract(ctx, callMsg, nil) // Contract Call

	if len(response) != 0 {
		logic_contract := common.HexToAddress(hex.EncodeToString(response))
		// code, _ := ethclient.CodeAt(ctx, logic_contract, nil)
		// encode_Code := hex.EncodeToString(code)
		return logic_contract.Hex()
	} else {
		charstr := "7f"
		storage := makeSliceUnique(findAllOccurrences(encodeCode, charstr))
		for _, check := range storage {
			address := common.HexToAddress(account.Hex())
			position := common.HexToHash(check)
			result, _ := ethclient.StorageAt(ctx, address, position, nil)
			encodedString := hex.EncodeToString(result)
			logic_contract := common.HexToAddress(encodedString)
			var logic_contract_hex string
			if logic_contract.Hex() == "0x0000000000000000000000000000000000000000" {
				continue
			} else {
				logic_contract_hex = logic_contract.Hex()
			}
			code, _ := ethclient.CodeAt(ctx, logic_contract, nil)
			// encode_Code := hex.EncodeToString(code)
			if len(code) > 2 {
				return logic_contract_hex
			}
		}
	}
	return ""
}

func is_UUP(ctx context.Context, ethclient *ethclient.Client, account common.Address, code []byte, encodeCode string) string {
	address := common.HexToAddress(account.Hex())
	position := common.HexToHash("0xc5f16f0fcc639fa48a6947836d9850f504798523bf8c9a3a87d5876cf622bcf7")
	result, err := ethclient.StorageAt(ctx, address, position, nil)
	if err == nil {
		encodedString := hex.EncodeToString(result)
		logic_contract := common.HexToAddress(encodedString)

		return logic_contract.Hex()
	} else {
		return ""
	}
}

// 로직 컨트랙트 주소 호출
