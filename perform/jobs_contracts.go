package perform

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"

	log "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	compilers "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/eris-compilers"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
)

func PackageDeployJob(pkgDeploy *definitions.PackageDeploy, do *definitions.Do) (string, error) {
	// todo
	var result string

	return result, nil
}

func DeployJob(deploy *definitions.Deploy, do *definitions.Do) (result string, err error) {
	// Preprocess variables
	deploy.Source, _ = util.PreProcess(deploy.Source, do)
	deploy.Contract, _ = util.PreProcess(deploy.Contract, do)
	deploy.Instance, _ = util.PreProcess(deploy.Instance, do)
	deploy.Libraries, _ = util.PreProcess(deploy.Libraries, do)
	deploy.Amount, _ = util.PreProcess(deploy.Amount, do)
	deploy.Nonce, _ = util.PreProcess(deploy.Nonce, do)
	deploy.Fee, _ = util.PreProcess(deploy.Fee, do)
	deploy.Gas, _ = util.PreProcess(deploy.Gas, do)

	// trim the extension
	contractName := strings.TrimSuffix(deploy.Contract, filepath.Ext(deploy.Contract))
	// Use default
	deploy.Source = useDefault(deploy.Source, do.Package.Account)
	deploy.Instance = useDefault(deploy.Instance, contractName)
	deploy.Amount = useDefault(deploy.Amount, do.DefaultAmount)
	deploy.Fee = useDefault(deploy.Fee, do.DefaultFee)
	deploy.Gas = useDefault(deploy.Gas, do.DefaultGas)

	// assemble contract
	var p string
	if _, err := os.Stat(deploy.Contract); err == nil {
		p = deploy.Contract
	} else {
		p = filepath.Join(do.ContractsPath, deploy.Contract)
	}
	log.WithField("=>", p).Debug("Contract path")

	// use the proper compiler
	if do.Compiler != "" {
		log.WithField("=>", do.Compiler).Debug("Setting compiler path")
		if err := setCompiler(do, p); err != nil {
			return "", err
		}
	}

	// compile
	resp := compilers.Compile(p, deploy.Libraries)

	if resp.Error != "" {
		return "", fmt.Errorf(resp.Error)
	}

	// Don't use pubKey if account override
	var oldKey string
	if deploy.Source != do.Package.Account {
		oldKey = do.PublicKey
		do.PublicKey = ""
	}

	//fmt.Println(resp)
	for _, r := range resp.Objects {
		log.WithField("=>", string(r.ABI)).Debug("Abi spec")
		if r.Bytecode == nil {
			continue
		}
		contractCode := hex.EncodeToString(r.Bytecode)

		// additional data may be sent along with the contract
		// these are naively added to the end of the contract code using standard
		// mint packing
		if deploy.Data != "" {
			splitout := strings.Split(deploy.Data, " ")
			for _, s := range splitout {
				s, _ = util.PreProcess(s, do)
				addOns := common.LeftPadString(common.StripHex(common.Coerce2Hex(s)), 64)
				log.WithField("=>", contractCode).Debug("Contract Code")
				log.WithField("=>", addOns).Debug("Additional Data")
				contractCode = contractCode + addOns
			}
		}
		// Save ABI
		if _, err := os.Stat(do.ABIPath); os.IsNotExist(err) {
			if err := os.Mkdir(do.ABIPath, 0775); err != nil {
				return "", err
			}
		}

		// saving contract/library abi
		abiLocation := filepath.Join(do.ABIPath, r.Objectname)
		log.WithField("=>", abiLocation).Debug("Saving ABI")
		if err := ioutil.WriteFile(abiLocation, []byte(r.ABI), 0664); err != nil {
			return "", err
		}

		if strings.ToLower(r.Objectname) == strings.ToLower(deploy.Instance) {
			// Deploy contract
			log.WithFields(log.Fields{
				"source": deploy.Source,
				"code":   contractCode,
			}).Info("Deploying Contract")

			tx, err := core.Call(do.Chain, do.Signer, do.PublicKey, deploy.Source, "", deploy.Amount, deploy.Nonce, deploy.Gas, deploy.Fee, contractCode)
			if err != nil {
				return "", fmt.Errorf("Error deploying contract %s: %v", p, err)
			}

			// Sign, broadcast, display
			result, err = deployFinalize(do, tx, deploy.Wait)

			// saving contract/library abi at abi/address
			abiLocation := filepath.Join(do.ABIPath, result)
			log.WithField("=>", abiLocation).Debug("Saving ABI")
			if err := ioutil.WriteFile(abiLocation, []byte(r.ABI), 0664); err != nil {
				return "", err
			}
		}

	}

	// Don't use pubKey if account override
	if deploy.Source != do.Package.Account {
		do.PublicKey = oldKey
	}

	return result, nil
}

func CallJob(call *definitions.Call, do *definitions.Do) (string, error) {
	// Preprocess variables
	call.Source, _ = util.PreProcess(call.Source, do)
	call.Destination, _ = util.PreProcess(call.Destination, do)
	call.Amount, _ = util.PreProcess(call.Amount, do)
	call.Nonce, _ = util.PreProcess(call.Nonce, do)
	call.Fee, _ = util.PreProcess(call.Fee, do)
	call.Gas, _ = util.PreProcess(call.Gas, do)

	// Use default
	call.Source = useDefault(call.Source, do.Package.Account)
	call.Amount = useDefault(call.Amount, do.DefaultAmount)
	call.Fee = useDefault(call.Fee, do.DefaultFee)
	call.Gas = useDefault(call.Gas, do.DefaultGas)

	var err error
	call.Data, err = util.ReadAbiFormulateCall(call.Destination, call.Data, do)
	if err != nil {
		log.Error("Error Formulating Call from ABI.")
		return "", err
	}

	// Don't use pubKey if account override
	var oldKey string
	if call.Source != do.Package.Account {
		oldKey = do.PublicKey
		do.PublicKey = ""
	}

	log.WithFields(log.Fields{
		"destination": call.Destination,
		"data":        call.Data,
	}).Info("Calling")

	tx, err := core.Call(do.Chain, do.Signer, do.PublicKey, call.Source, call.Destination, call.Amount, call.Nonce, call.Gas, call.Fee, call.Data)
	if err != nil {
		return "", err
	}

	// Don't use pubKey if account override
	if call.Source != do.Package.Account {
		do.PublicKey = oldKey
	}

	// Sign, broadcast, display
	return txFinalize(do, tx, call.Wait)
}

func setCompiler(do *definitions.Do, tocompile string) error {
	lang, err := compilers.LangFromFile(tocompile)
	if err != nil {
		return err
	}

	url := do.Compiler + "/" + "compile"
	compilers.SetLanguageURL(lang, url)
	return nil
}

func deployFinalize(do *definitions.Do, tx interface{}, wait bool) (string, error) {
	var result string

	res, err := core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx.(types.Tx), true, true, wait)
	if err != nil {
		log.Error("ERROR =>")
		return "", err
	}

	if err := util.ReadTxSignAndBroadcast(res, err); err != nil {
		log.Error("ERROR =>")
		return "", err
	}

	result = fmt.Sprintf("%X", res.Address)
	return result, nil
}
