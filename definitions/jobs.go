package definitions

// ------------------------------------------------------------------------
// Util Jobs
// ------------------------------------------------------------------------

type Account struct {
	// (Required) address of the account which should be used as the default (if source) is
	// not given for future transactions. Will make sure the eris-keys has the public key
	// for the account. Generally account should be the first job called unless it is used
	// via a flag or environment variables to establish what default to use.
	Address string `mapstructure:"address" json:"address" yaml:"address" toml:"address"`
}

type Set struct {
	// (Required) value which should be saved along with the jobName (which will be the key)
	// this is useful to set variables which can be used throughout the epm definition file
	Value string `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}

// ------------------------------------------------------------------------
// Transaction Jobs
// ------------------------------------------------------------------------

type Send struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) address of the account to send the tokens
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	// (Required) amount of tokens to send from the `source` to the `destination`
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	// (Optional) wait for the transaction to be confirmed in the blockchain before proceeding
	Wait bool `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type RegisterName struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) name which will be registered
	Name string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	// (Optional, if data_file is used; otherwise required) data which will be stored at the `name` key
	Data string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	// (Optional) csv file in the form (name,data[,amount]) which can be used to bulk register names
	DataFile string `mapstructure:"data_file" json:"data_file" yaml:"data_file" toml:"data_file"`
	// (Optional) amount of blocks which the name entry will be reserved for the registering user
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional) validators' fee
	Fee string `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	// (Optional) wait for the transaction to be confirmed in the blockchain before proceeding
	Wait bool `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Permission struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) actions must be in the set ["set_base", "unset_base", "set_global", "add_role" "rm_role"]
	Action string `mapstructure:"action" json:"action" yaml:"action" toml:"action"`
	// (Required, unless add_role or rm_role action selected) the name of the permission flag which is to
	// be updated
	PermissionFlag string `mapstructure:"permission" json:"permission" yaml:"permission" toml:"permission"`
	// (Required) the value of the permission or role which is to be updated
	Value string `mapstructure:"value" json:"value" yaml:"value" toml:"value"`
	// (Required) the target account which is to be updated
	Target string `mapstructure:"target" json:"target" yaml:"target" toml:"target"`
	// (Required, if add_role or rm_role action selected) the role which should be given to the account
	Role string `mapstructure:"role" json:"role" yaml:"role" toml:"role"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	// (Optional) wait for the transaction to be confirmed in the blockchain before proceeding
	Wait bool `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Bond struct {
	// (Required) public key of the address which will be bonded
	PublicKey string `mapstructure:"pub_key" json:"pub_key" yaml:"pub_key" toml:"pub_key"`
	// (Required) address of the account which will be bonded
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) amount of tokens which will be bonded
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	// (Optional) wait for the transaction to be confirmed in the blockchain before proceeding
	Wait bool `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Unbond struct {
	// (Required) address of the account which to unbond
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) block on which the unbonding will take place (users may unbond at any
	// time >= currentBlock)
	Height string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
	// (Optional) wait for the transaction to be confirmed in the blockchain before proceeding
	Wait bool `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Rebond struct {
	// (Required) address of the account which to rebond
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) block on which the rebonding will take place (users may rebond at any
	// time >= (unbondBlock || currentBlock))
	Height string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
	// (Optional) wait for the transaction to be confirmed in the blockchain before proceeding
	Wait bool `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

// ------------------------------------------------------------------------
// Contracts Jobs
// ------------------------------------------------------------------------

type PackageDeploy struct {
	// TODO
}

type Deploy struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) the filepath to the contract file. this should be relative to the current path **or**
	// relative to the contracts path established via the --contracts-path flag or the $EPM_CONTRACTS_PATH
	// environment variable
	Contract string `mapstructure:"contract" json:"contract" yaml:"contract" toml:"contract"`
	// (Optional) the name of contract to instantiate (it has to be one of the contracts present)
	// in the file defined in Contract above.
	// When none is provided, the system will choose the contract with the same name as that file.
	Instance string `mapstructure:"instance" json:"instance" yaml:"instance" toml:"instance"`
	// (Optional) list of Name:Address separated by spaces of libraries (see solc --help)
	Libraries string `mapstructure:"libraries" json:"libraries" yaml:"libraries" toml:"libraries"`
	// (Optional) TODO: additional arguments to send along with the contract code
	Data string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	// (Optional) amount of tokens to send to the contract which will (after deployment) reside in the
	// contract's account
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional) validators' fee
	Fee string `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	// (Optional) amount of gas which should be sent along with the contract deployment transaction
	Gas string `mapstructure:"gas" json:"gas" yaml:"gas" toml:"gas"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	// (Optional) wait for the transaction to be confirmed in the blockchain before proceeding
	Wait bool `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Call struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) address of the contract which should be called
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	// (Required) data which should be called. will use the eris-abi tooling under the hood to formalize the
	// transaction
	Data string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	// (Optional) amount of tokens to send to the contract
	Amount string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	// (Optional) validators' fee
	Fee string `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	// (Optional) amount of gas which should be sent along with the call transaction
	Gas string `mapstructure:"gas" json:"gas" yaml:"gas" toml:"gas"`
	// (Optional, advanced only) nonce to use when eris-keys signs the transaction (do not use unless you
	// know what you're doing)
	Nonce string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	// (Optional) wait for the transaction to be confirmed in the blockchain before proceeding
	Wait bool `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

// ------------------------------------------------------------------------
// State Jobs
// ------------------------------------------------------------------------

type DumpState struct {
	WithValidators bool   `mapstructure:"include-validators" json:"include-validators" yaml:"include-validators" toml:"include-validators"`
	ToIPFS         bool   `mapstructure:"to-ipfs" json:"to-ipfs" yaml:"to-ipfs" toml:"to-ipfs"`
	ToFile         bool   `mapstructure:"to-file" json:"to-file" yaml:"to-file" toml:"to-file"`
	IPFSHost       string `mapstructure:"ipfs-host" json:"ipfs-host" yaml:"ipfs-host" toml:"ipfs-host"`
	FilePath       string `mapstructure:"file" json:"file" yaml:"file" toml:"file"`
}

type RestoreState struct {
	FromIPFS bool   `mapstructure:"from-ipfs" json:"from-ipfs" yaml:"from-ipfs" toml:"from-ipfs"`
	FromFile bool   `mapstructure:"from-file" json:"from-file" yaml:"from-file" toml:"from-file"`
	IPFSHost string `mapstructure:"ipfs-host" json:"ipfs-host" yaml:"ipfs-host" toml:"ipfs-host"`
	FilePath string `mapstructure:"file" json:"file" yaml:"file" toml:"file"`
}

// ------------------------------------------------------------------------
// Testing Jobs
// ------------------------------------------------------------------------

// aka. Simulated Call.
type QueryContract struct {
	// (Optional, if account job or global account set) address of the account from which to send (the
	// public key for the account must be available to eris-keys)
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	// (Required) address of the contract which should be called
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	// (Required) data which should be called. will use the eris-abi tooling under the hood to formalize the
	// transaction. QueryContract will usually be used with "accessor" functions in contracts
	Data string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
}

type QueryAccount struct {
	// (Required) address of the account which should be queried
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	// (Required) field which should be queried. If users are trying to query the permissions of the
	// account one can get either the `permissions.base` which will return the base permission of the
	// account, or one can get the `permissions.set` which will return the setBit of the account.
	Field string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

type QueryName struct {
	// (Required) name which should be queried
	Name string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	// (Required) field which should be quiried (generally will be "data" to get the registered "name")
	Field string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

type QueryVals struct {
	// (Required) should be of the set ["bonded_validators" or "unbonding_validators"] and it will
	// return a comma separated listing of the addresses which fall into one of those categories
	Field string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

type Assert struct {
	// (Required) key which should be used for the assertion. This is usually known as the "expected"
	// value in most testing suites
	Key string `mapstructure:"key" json:"key" yaml:"key" toml:"key"`
	// (Required) must be of the set ["eq", "ne", "ge", "gt", "le", "lt", "==", "!=", ">=", ">", "<=", "<"]
	// establishes the relation to be tested by the assertion. If a strings key:value pair is being used
	// only the equals or not-equals relations may be used as the key:value will try to be converted to
	// ints for the remainder of the relations. if strings are passed to them then eris:pm will return an
	// error
	Relation string `mapstructure:"relation" json:"relation" yaml:"relation" toml:"relation"`
	// (Required) value which should be used for the assertion. This is usually known as the "given"
	// value in most testing suites. Generally it will be a variable expansion from one of the query
	// jobs.
	Value string `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}
