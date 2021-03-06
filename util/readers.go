package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"

	log "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	ebi "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/eris-abi/core"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
)

// This is a closer function which is called by most of the tx_run functions
func ReadTxSignAndBroadcast(result *core.TxResult, err error) error {
	// if there's an error just return.
	if err != nil {
		return err
	}

	// if there is nothing to unpack then just return.
	if result == nil {
		return nil
	}

	// Unpack and display for the user.
	addr := fmt.Sprintf("%X", result.Address)
	hash := fmt.Sprintf("%X", result.Hash)
	blkHash := fmt.Sprintf("%X", result.BlockHash)
	ret := fmt.Sprintf("%X", result.Return)

	if result.Address != nil {
		log.WithField("=>", addr).Warn("Contract Address")
		log.WithField("=>", hash).Warn("Transaction Hash")
	} else {
		log.WithField("=>", hash).Warn("Transaction Hash")
		log.WithField("=>", blkHash).Debug("Block Hash")
		if len(result.Return) != 0 {
			log.WithField("=>", ret).Warn("Return Value")
			log.WithField("=>", result.Exception).Debug("Exception")
		}
	}

	return nil
}

func ReadAbiFormulateCall(contract, dataRaw string, do *definitions.Do) (string, error) {
	abiSpecBytes, err := readAbi(do.ABIPath, contract)
	if err != nil {
		return "", err
	}
	log.WithField("=>", abiSpecBytes).Debug("ABI Spec")

	// Process and Pack the Call
	funcName, args := abiPreProcess(dataRaw, do)

	var totalArgs []string
	totalArgs = append(totalArgs, funcName)
	totalArgs = append(totalArgs, args...)
	arg := strings.Join(totalArgs, "\n")
	log.WithFields(log.Fields{
		"function name": funcName,
		"arguments":     arg,
	}).Debug("Packing Call via ABI =>")

	return ebi.Packer(abiSpecBytes, totalArgs...)
}

func ReadAndDecodeContractReturn(contract, dataRaw, resultRaw string, do *definitions.Do) (string, error) {
	abiSpecBytes, err := readAbi(do.ABIPath, contract)
	if err != nil {
		return "", err
	}
	log.WithField("=>", abiSpecBytes).Debug("ABI Spec")

	// Process and Pack the Call
	funcName, _ := abiPreProcess(dataRaw, do)

	// Unpack the result
	res, err := ebi.UnPacker(abiSpecBytes, funcName, resultRaw, false)
	if err != nil {
		return "", err
	}

	// Wrangle these returns
	type ContractReturn struct {
		Name  string `mapstructure:"," json:","`
		Type  string `mapstructure:"," json:","`
		Value string `mapstructure:"," json:","`
	}
	var resTotal []ContractReturn
	err = json.Unmarshal([]byte(res), &resTotal)
	if err != nil {
		return "", err
	}

	// Get the value
	result := resTotal[0].Value
	return result, nil
}

func abiPreProcess(dataRaw string, do *definitions.Do) (string, []string) {
	var dataNew []string

	data := strings.Split(dataRaw, " ")
	for _, d := range data {
		d, _ = PreProcess(d, do)
		dataNew = append(dataNew, d)
	}

	funcName := dataNew[0]
	args := dataNew[1:]

	return funcName, args
}

func readAbi(root, contract string) ([]byte, error) {
	p := path.Join(root, common.StripHex(contract))
	if _, err := os.Stat(p); err != nil {
		return []byte{}, fmt.Errorf("Abi doesn't exist for =>\t%s", p)
	}

	b, err := ioutil.ReadFile(p)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
