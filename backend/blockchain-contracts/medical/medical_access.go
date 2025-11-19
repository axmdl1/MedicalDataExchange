// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package medical

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MedicalAccessMetaData contains all meta data concerning the MedicalAccess contract.
var MedicalAccessMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"tokenHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiresAt\",\"type\":\"uint256\"}],\"name\":\"RequestApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"RequestConsumed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recordId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"patient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clinic\",\"type\":\"string\"}],\"name\":\"RequestCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"RequestRevoked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"tokenHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiresAt\",\"type\":\"uint256\"}],\"name\":\"approveAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"}],\"name\":\"consumeAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"}],\"name\":\"getRequest\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"recordId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"patient\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"clinic\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"tokenHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiresAt\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"status\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recordId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"patient\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"clinic\",\"type\":\"string\"}],\"name\":\"requestAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"tokenHash\",\"type\":\"bytes32\"}],\"name\":\"verifyToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506115918061001f6000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80632998bb171461005c5780633cad02591461008c57806366b2db3a146100a857806385b92413146100c4578063a5b2c15b146100f9575b600080fd5b610076600480360381019061007191906109c7565b610115565b6040516100839190610a42565b60405180910390f35b6100a660048036038101906100a19190610a5d565b6101f9565b005b6100c260048036038101906100bd9190610ae0565b61038c565b005b6100de60048036038101906100d99190610a5d565b610535565b6040516100f096959493929190610c02565b60405180910390f35b610113600480360381019061010e9190610c7f565b6107c3565b005b6000806000858560405161012a929190610da7565b90815260200160405180910390209050600081600501805461014b90610def565b90500361015c5760009150506101f2565b828160030154146101715760009150506101f2565b80600401544211156101875760009150506101f2565b6040518060400160405280600881526020017f617070726f76656400000000000000000000000000000000000000000000000081525080519060200120816005016040516101d59190610ec3565b6040518091039020146101ec5760009150506101f2565b60019150505b9392505050565b600080838360405161020c929190610da7565b90815260200160405180910390209050600081600501805461022d90610def565b90500361026f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161026690610f26565b60405180910390fd5b6040518060400160405280600881526020017f617070726f76656400000000000000000000000000000000000000000000000081525080519060200120816005016040516102bd9190610ec3565b604051809103902014610305576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102fc90610f92565b60405180910390fd5b6040518060400160405280600881526020017f636f6e73756d656400000000000000000000000000000000000000000000000081525081600501908161034b919061118d565b507ff0b51072dcb1f05229e3efb58f6fec5511df95a8de7ff025b8735e723fb1ba4c83833360405161037f939291906112cd565b60405180910390a1505050565b600080858560405161039f929190610da7565b9081526020016040518091039020905060008160050180546103c090610def565b905003610402576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103f990610f26565b60405180910390fd5b6040518060400160405280600781526020017f70656e64696e670000000000000000000000000000000000000000000000000081525080519060200120816005016040516104509190610ec3565b604051809103902014610498576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161048f9061134b565b60405180910390fd5b8281600301819055508181600401819055506040518060400160405280600881526020017f617070726f7665640000000000000000000000000000000000000000000000008152508160050190816104f0919061118d565b507f378120d9127d115d83c0f0c1e108ff1610c9cc3899e4bb4b88ea6b40a90f19c585858585604051610526949392919061136b565b60405180910390a15050505050565b606080606060008060606000808989604051610552929190610da7565b90815260200160405180910390209050806000018160010182600201836003015484600401548560050185805461058890610def565b80601f01602080910402602001604051908101604052809291908181526020018280546105b490610def565b80156106015780601f106105d657610100808354040283529160200191610601565b820191906000526020600020905b8154815290600101906020018083116105e457829003601f168201915b5050505050955084805461061490610def565b80601f016020809104026020016040519081016040528092919081815260200182805461064090610def565b801561068d5780601f106106625761010080835404028352916020019161068d565b820191906000526020600020905b81548152906001019060200180831161067057829003601f168201915b505050505094508380546106a090610def565b80601f01602080910402602001604051908101604052809291908181526020018280546106cc90610def565b80156107195780601f106106ee57610100808354040283529160200191610719565b820191906000526020600020905b8154815290600101906020018083116106fc57829003601f168201915b5050505050935080805461072c90610def565b80601f016020809104026020016040519081016040528092919081815260200182805461075890610def565b80156107a55780601f1061077a576101008083540402835291602001916107a5565b820191906000526020600020905b81548152906001019060200180831161078857829003601f168201915b50505050509050965096509650965096509650509295509295509295565b60008089896040516107d6929190610da7565b9081526020016040518091039020905060008160050180546107f790610def565b905014610839576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610830906113f7565b60405180910390fd5b868682600001918261084c929190611422565b508484826001019182610860929190611422565b508282826002019182610874929190611422565b506000801b8160030181905550600081600401819055506040518060400160405280600781526020017f70656e64696e67000000000000000000000000000000000000000000000000008152508160050190816108d1919061118d565b507fd3857f86ae3eb5b776d6168ef81042dc9e481ad045554ff37ce44ae19b8e84c7898989898989898960405161090f9897969594939291906114f2565b60405180910390a1505050505050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f8401126109515761095061092c565b5b8235905067ffffffffffffffff81111561096e5761096d610931565b5b60208301915083600182028301111561098a57610989610936565b5b9250929050565b6000819050919050565b6109a481610991565b81146109af57600080fd5b50565b6000813590506109c18161099b565b92915050565b6000806000604084860312156109e0576109df610922565b5b600084013567ffffffffffffffff8111156109fe576109fd610927565b5b610a0a8682870161093b565b93509350506020610a1d868287016109b2565b9150509250925092565b60008115159050919050565b610a3c81610a27565b82525050565b6000602082019050610a576000830184610a33565b92915050565b60008060208385031215610a7457610a73610922565b5b600083013567ffffffffffffffff811115610a9257610a91610927565b5b610a9e8582860161093b565b92509250509250929050565b6000819050919050565b610abd81610aaa565b8114610ac857600080fd5b50565b600081359050610ada81610ab4565b92915050565b60008060008060608587031215610afa57610af9610922565b5b600085013567ffffffffffffffff811115610b1857610b17610927565b5b610b248782880161093b565b94509450506020610b37878288016109b2565b9250506040610b4887828801610acb565b91505092959194509250565b600081519050919050565b600082825260208201905092915050565b60005b83811015610b8e578082015181840152602081019050610b73565b60008484015250505050565b6000601f19601f8301169050919050565b6000610bb682610b54565b610bc08185610b5f565b9350610bd0818560208601610b70565b610bd981610b9a565b840191505092915050565b610bed81610991565b82525050565b610bfc81610aaa565b82525050565b600060c0820190508181036000830152610c1c8189610bab565b90508181036020830152610c308188610bab565b90508181036040830152610c448187610bab565b9050610c536060830186610be4565b610c606080830185610bf3565b81810360a0830152610c728184610bab565b9050979650505050505050565b6000806000806000806000806080898b031215610c9f57610c9e610922565b5b600089013567ffffffffffffffff811115610cbd57610cbc610927565b5b610cc98b828c0161093b565b9850985050602089013567ffffffffffffffff811115610cec57610ceb610927565b5b610cf88b828c0161093b565b9650965050604089013567ffffffffffffffff811115610d1b57610d1a610927565b5b610d278b828c0161093b565b9450945050606089013567ffffffffffffffff811115610d4a57610d49610927565b5b610d568b828c0161093b565b92509250509295985092959890939650565b600081905092915050565b82818337600083830152505050565b6000610d8e8385610d68565b9350610d9b838584610d73565b82840190509392505050565b6000610db4828486610d82565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610e0757607f821691505b602082108103610e1a57610e19610dc0565b5b50919050565b600081905092915050565b60008190508160005260206000209050919050565b60008154610e4d81610def565b610e578186610e20565b94506001821660008114610e725760018114610e8757610eba565b60ff1983168652811515820286019350610eba565b610e9085610e2b565b60005b83811015610eb257815481890152600182019150602081019050610e93565b838801955050505b50505092915050565b6000610ecf8284610e40565b915081905092915050565b7f72657175657374206e6f7420666f756e64000000000000000000000000000000600082015250565b6000610f10601183610b5f565b9150610f1b82610eda565b602082019050919050565b60006020820190508181036000830152610f3f81610f03565b9050919050565b7f6e6f7420617070726f7665640000000000000000000000000000000000000000600082015250565b6000610f7c600c83610b5f565b9150610f8782610f46565b602082019050919050565b60006020820190508181036000830152610fab81610f6f565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026110437fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82611006565b61104d8683611006565b95508019841693508086168417925050509392505050565b6000819050919050565b600061108a61108561108084610aaa565b611065565b610aaa565b9050919050565b6000819050919050565b6110a48361106f565b6110b86110b082611091565b848454611013565b825550505050565b600090565b6110cd6110c0565b6110d881848461109b565b505050565b5b818110156110fc576110f16000826110c5565b6001810190506110de565b5050565b601f8211156111415761111281610fe1565b61111b84610ff6565b8101602085101561112a578190505b61113e61113685610ff6565b8301826110dd565b50505b505050565b600082821c905092915050565b600061116460001984600802611146565b1980831691505092915050565b600061117d8383611153565b9150826002028217905092915050565b61119682610b54565b67ffffffffffffffff8111156111af576111ae610fb2565b5b6111b98254610def565b6111c4828285611100565b600060209050601f8311600181146111f757600084156111e5578287015190505b6111ef8582611171565b865550611257565b601f19841661120586610fe1565b60005b8281101561122d57848901518255600182019150602085019450602081019050611208565b8683101561124a5784890151611246601f891682611153565b8355505b6001600288020188555050505b505050505050565b600061126b8385610b5f565b9350611278838584610d73565b61128183610b9a565b840190509392505050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006112b78261128c565b9050919050565b6112c7816112ac565b82525050565b600060408201905081810360008301526112e881858761125f565b90506112f760208301846112be565b949350505050565b7f6e6f742070656e64696e67000000000000000000000000000000000000000000600082015250565b6000611335600b83610b5f565b9150611340826112ff565b602082019050919050565b6000602082019050818103600083015261136481611328565b9050919050565b6000606082019050818103600083015261138681868861125f565b90506113956020830185610be4565b6113a26040830184610bf3565b95945050505050565b7f7265717565737420657869737473000000000000000000000000000000000000600082015250565b60006113e1600e83610b5f565b91506113ec826113ab565b602082019050919050565b60006020820190508181036000830152611410816113d4565b9050919050565b600082905092915050565b61142c8383611417565b67ffffffffffffffff81111561144557611444610fb2565b5b61144f8254610def565b61145a828285611100565b6000601f8311600181146114895760008415611477578287013590505b6114818582611171565b8655506114e9565b601f19841661149786610fe1565b60005b828110156114bf5784890135825560018201915060208501945060208101905061149a565b868310156114dc57848901356114d8601f891682611153565b8355505b6001600288020188555050505b50505050505050565b6000608082019050818103600083015261150d818a8c61125f565b9050818103602083015261152281888a61125f565b9050818103604083015261153781868861125f565b9050818103606083015261154c81848661125f565b9050999850505050505050505056fea264697066735822122051ce1cb8026cdbee09b92b0914e41f9c9ce94643a523fb7a40d87cfbe921b5ad64736f6c634300081e0033",
}

// MedicalAccessABI is the input ABI used to generate the binding from.
// Deprecated: Use MedicalAccessMetaData.ABI instead.
var MedicalAccessABI = MedicalAccessMetaData.ABI

// MedicalAccessBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MedicalAccessMetaData.Bin instead.
var MedicalAccessBin = MedicalAccessMetaData.Bin

// DeployMedicalAccess deploys a new Ethereum contract, binding an instance of MedicalAccess to it.
func DeployMedicalAccess(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MedicalAccess, error) {
	parsed, err := MedicalAccessMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MedicalAccessBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MedicalAccess{MedicalAccessCaller: MedicalAccessCaller{contract: contract}, MedicalAccessTransactor: MedicalAccessTransactor{contract: contract}, MedicalAccessFilterer: MedicalAccessFilterer{contract: contract}}, nil
}

// MedicalAccess is an auto generated Go binding around an Ethereum contract.
type MedicalAccess struct {
	MedicalAccessCaller     // Read-only binding to the contract
	MedicalAccessTransactor // Write-only binding to the contract
	MedicalAccessFilterer   // Log filterer for contract events
}

// MedicalAccessCaller is an auto generated read-only Go binding around an Ethereum contract.
type MedicalAccessCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedicalAccessTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MedicalAccessTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedicalAccessFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MedicalAccessFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MedicalAccessSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MedicalAccessSession struct {
	Contract     *MedicalAccess    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MedicalAccessCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MedicalAccessCallerSession struct {
	Contract *MedicalAccessCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// MedicalAccessTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MedicalAccessTransactorSession struct {
	Contract     *MedicalAccessTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// MedicalAccessRaw is an auto generated low-level Go binding around an Ethereum contract.
type MedicalAccessRaw struct {
	Contract *MedicalAccess // Generic contract binding to access the raw methods on
}

// MedicalAccessCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MedicalAccessCallerRaw struct {
	Contract *MedicalAccessCaller // Generic read-only contract binding to access the raw methods on
}

// MedicalAccessTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MedicalAccessTransactorRaw struct {
	Contract *MedicalAccessTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMedicalAccess creates a new instance of MedicalAccess, bound to a specific deployed contract.
func NewMedicalAccess(address common.Address, backend bind.ContractBackend) (*MedicalAccess, error) {
	contract, err := bindMedicalAccess(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MedicalAccess{MedicalAccessCaller: MedicalAccessCaller{contract: contract}, MedicalAccessTransactor: MedicalAccessTransactor{contract: contract}, MedicalAccessFilterer: MedicalAccessFilterer{contract: contract}}, nil
}

// NewMedicalAccessCaller creates a new read-only instance of MedicalAccess, bound to a specific deployed contract.
func NewMedicalAccessCaller(address common.Address, caller bind.ContractCaller) (*MedicalAccessCaller, error) {
	contract, err := bindMedicalAccess(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MedicalAccessCaller{contract: contract}, nil
}

// NewMedicalAccessTransactor creates a new write-only instance of MedicalAccess, bound to a specific deployed contract.
func NewMedicalAccessTransactor(address common.Address, transactor bind.ContractTransactor) (*MedicalAccessTransactor, error) {
	contract, err := bindMedicalAccess(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MedicalAccessTransactor{contract: contract}, nil
}

// NewMedicalAccessFilterer creates a new log filterer instance of MedicalAccess, bound to a specific deployed contract.
func NewMedicalAccessFilterer(address common.Address, filterer bind.ContractFilterer) (*MedicalAccessFilterer, error) {
	contract, err := bindMedicalAccess(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MedicalAccessFilterer{contract: contract}, nil
}

// bindMedicalAccess binds a generic wrapper to an already deployed contract.
func bindMedicalAccess(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MedicalAccessMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MedicalAccess *MedicalAccessRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MedicalAccess.Contract.MedicalAccessCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MedicalAccess *MedicalAccessRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MedicalAccess.Contract.MedicalAccessTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MedicalAccess *MedicalAccessRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MedicalAccess.Contract.MedicalAccessTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MedicalAccess *MedicalAccessCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MedicalAccess.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MedicalAccess *MedicalAccessTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MedicalAccess.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MedicalAccess *MedicalAccessTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MedicalAccess.Contract.contract.Transact(opts, method, params...)
}

// GetRequest is a free data retrieval call binding the contract method 0x85b92413.
//
// Solidity: function getRequest(string requestId) view returns(string recordId, string patient, string clinic, bytes32 tokenHash, uint256 expiresAt, string status)
func (_MedicalAccess *MedicalAccessCaller) GetRequest(opts *bind.CallOpts, requestId string) (struct {
	RecordId  string
	Patient   string
	Clinic    string
	TokenHash [32]byte
	ExpiresAt *big.Int
	Status    string
}, error) {
	var out []interface{}
	err := _MedicalAccess.contract.Call(opts, &out, "getRequest", requestId)

	outstruct := new(struct {
		RecordId  string
		Patient   string
		Clinic    string
		TokenHash [32]byte
		ExpiresAt *big.Int
		Status    string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RecordId = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Patient = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Clinic = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.TokenHash = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.ExpiresAt = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[5], new(string)).(*string)

	return *outstruct, err

}

// GetRequest is a free data retrieval call binding the contract method 0x85b92413.
//
// Solidity: function getRequest(string requestId) view returns(string recordId, string patient, string clinic, bytes32 tokenHash, uint256 expiresAt, string status)
func (_MedicalAccess *MedicalAccessSession) GetRequest(requestId string) (struct {
	RecordId  string
	Patient   string
	Clinic    string
	TokenHash [32]byte
	ExpiresAt *big.Int
	Status    string
}, error) {
	return _MedicalAccess.Contract.GetRequest(&_MedicalAccess.CallOpts, requestId)
}

// GetRequest is a free data retrieval call binding the contract method 0x85b92413.
//
// Solidity: function getRequest(string requestId) view returns(string recordId, string patient, string clinic, bytes32 tokenHash, uint256 expiresAt, string status)
func (_MedicalAccess *MedicalAccessCallerSession) GetRequest(requestId string) (struct {
	RecordId  string
	Patient   string
	Clinic    string
	TokenHash [32]byte
	ExpiresAt *big.Int
	Status    string
}, error) {
	return _MedicalAccess.Contract.GetRequest(&_MedicalAccess.CallOpts, requestId)
}

// VerifyToken is a free data retrieval call binding the contract method 0x2998bb17.
//
// Solidity: function verifyToken(string requestId, bytes32 tokenHash) view returns(bool)
func (_MedicalAccess *MedicalAccessCaller) VerifyToken(opts *bind.CallOpts, requestId string, tokenHash [32]byte) (bool, error) {
	var out []interface{}
	err := _MedicalAccess.contract.Call(opts, &out, "verifyToken", requestId, tokenHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyToken is a free data retrieval call binding the contract method 0x2998bb17.
//
// Solidity: function verifyToken(string requestId, bytes32 tokenHash) view returns(bool)
func (_MedicalAccess *MedicalAccessSession) VerifyToken(requestId string, tokenHash [32]byte) (bool, error) {
	return _MedicalAccess.Contract.VerifyToken(&_MedicalAccess.CallOpts, requestId, tokenHash)
}

// VerifyToken is a free data retrieval call binding the contract method 0x2998bb17.
//
// Solidity: function verifyToken(string requestId, bytes32 tokenHash) view returns(bool)
func (_MedicalAccess *MedicalAccessCallerSession) VerifyToken(requestId string, tokenHash [32]byte) (bool, error) {
	return _MedicalAccess.Contract.VerifyToken(&_MedicalAccess.CallOpts, requestId, tokenHash)
}

// ApproveAccess is a paid mutator transaction binding the contract method 0x66b2db3a.
//
// Solidity: function approveAccess(string requestId, bytes32 tokenHash, uint256 expiresAt) returns()
func (_MedicalAccess *MedicalAccessTransactor) ApproveAccess(opts *bind.TransactOpts, requestId string, tokenHash [32]byte, expiresAt *big.Int) (*types.Transaction, error) {
	return _MedicalAccess.contract.Transact(opts, "approveAccess", requestId, tokenHash, expiresAt)
}

// ApproveAccess is a paid mutator transaction binding the contract method 0x66b2db3a.
//
// Solidity: function approveAccess(string requestId, bytes32 tokenHash, uint256 expiresAt) returns()
func (_MedicalAccess *MedicalAccessSession) ApproveAccess(requestId string, tokenHash [32]byte, expiresAt *big.Int) (*types.Transaction, error) {
	return _MedicalAccess.Contract.ApproveAccess(&_MedicalAccess.TransactOpts, requestId, tokenHash, expiresAt)
}

// ApproveAccess is a paid mutator transaction binding the contract method 0x66b2db3a.
//
// Solidity: function approveAccess(string requestId, bytes32 tokenHash, uint256 expiresAt) returns()
func (_MedicalAccess *MedicalAccessTransactorSession) ApproveAccess(requestId string, tokenHash [32]byte, expiresAt *big.Int) (*types.Transaction, error) {
	return _MedicalAccess.Contract.ApproveAccess(&_MedicalAccess.TransactOpts, requestId, tokenHash, expiresAt)
}

// ConsumeAccess is a paid mutator transaction binding the contract method 0x3cad0259.
//
// Solidity: function consumeAccess(string requestId) returns()
func (_MedicalAccess *MedicalAccessTransactor) ConsumeAccess(opts *bind.TransactOpts, requestId string) (*types.Transaction, error) {
	return _MedicalAccess.contract.Transact(opts, "consumeAccess", requestId)
}

// ConsumeAccess is a paid mutator transaction binding the contract method 0x3cad0259.
//
// Solidity: function consumeAccess(string requestId) returns()
func (_MedicalAccess *MedicalAccessSession) ConsumeAccess(requestId string) (*types.Transaction, error) {
	return _MedicalAccess.Contract.ConsumeAccess(&_MedicalAccess.TransactOpts, requestId)
}

// ConsumeAccess is a paid mutator transaction binding the contract method 0x3cad0259.
//
// Solidity: function consumeAccess(string requestId) returns()
func (_MedicalAccess *MedicalAccessTransactorSession) ConsumeAccess(requestId string) (*types.Transaction, error) {
	return _MedicalAccess.Contract.ConsumeAccess(&_MedicalAccess.TransactOpts, requestId)
}

// RequestAccess is a paid mutator transaction binding the contract method 0xa5b2c15b.
//
// Solidity: function requestAccess(string requestId, string recordId, string patient, string clinic) returns()
func (_MedicalAccess *MedicalAccessTransactor) RequestAccess(opts *bind.TransactOpts, requestId string, recordId string, patient string, clinic string) (*types.Transaction, error) {
	return _MedicalAccess.contract.Transact(opts, "requestAccess", requestId, recordId, patient, clinic)
}

// RequestAccess is a paid mutator transaction binding the contract method 0xa5b2c15b.
//
// Solidity: function requestAccess(string requestId, string recordId, string patient, string clinic) returns()
func (_MedicalAccess *MedicalAccessSession) RequestAccess(requestId string, recordId string, patient string, clinic string) (*types.Transaction, error) {
	return _MedicalAccess.Contract.RequestAccess(&_MedicalAccess.TransactOpts, requestId, recordId, patient, clinic)
}

// RequestAccess is a paid mutator transaction binding the contract method 0xa5b2c15b.
//
// Solidity: function requestAccess(string requestId, string recordId, string patient, string clinic) returns()
func (_MedicalAccess *MedicalAccessTransactorSession) RequestAccess(requestId string, recordId string, patient string, clinic string) (*types.Transaction, error) {
	return _MedicalAccess.Contract.RequestAccess(&_MedicalAccess.TransactOpts, requestId, recordId, patient, clinic)
}

// MedicalAccessRequestApprovedIterator is returned from FilterRequestApproved and is used to iterate over the raw logs and unpacked data for RequestApproved events raised by the MedicalAccess contract.
type MedicalAccessRequestApprovedIterator struct {
	Event *MedicalAccessRequestApproved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MedicalAccessRequestApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MedicalAccessRequestApproved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MedicalAccessRequestApproved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MedicalAccessRequestApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MedicalAccessRequestApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MedicalAccessRequestApproved represents a RequestApproved event raised by the MedicalAccess contract.
type MedicalAccessRequestApproved struct {
	RequestId string
	TokenHash [32]byte
	ExpiresAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestApproved is a free log retrieval operation binding the contract event 0x378120d9127d115d83c0f0c1e108ff1610c9cc3899e4bb4b88ea6b40a90f19c5.
//
// Solidity: event RequestApproved(string requestId, bytes32 tokenHash, uint256 expiresAt)
func (_MedicalAccess *MedicalAccessFilterer) FilterRequestApproved(opts *bind.FilterOpts) (*MedicalAccessRequestApprovedIterator, error) {

	logs, sub, err := _MedicalAccess.contract.FilterLogs(opts, "RequestApproved")
	if err != nil {
		return nil, err
	}
	return &MedicalAccessRequestApprovedIterator{contract: _MedicalAccess.contract, event: "RequestApproved", logs: logs, sub: sub}, nil
}

// WatchRequestApproved is a free log subscription operation binding the contract event 0x378120d9127d115d83c0f0c1e108ff1610c9cc3899e4bb4b88ea6b40a90f19c5.
//
// Solidity: event RequestApproved(string requestId, bytes32 tokenHash, uint256 expiresAt)
func (_MedicalAccess *MedicalAccessFilterer) WatchRequestApproved(opts *bind.WatchOpts, sink chan<- *MedicalAccessRequestApproved) (event.Subscription, error) {

	logs, sub, err := _MedicalAccess.contract.WatchLogs(opts, "RequestApproved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MedicalAccessRequestApproved)
				if err := _MedicalAccess.contract.UnpackLog(event, "RequestApproved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestApproved is a log parse operation binding the contract event 0x378120d9127d115d83c0f0c1e108ff1610c9cc3899e4bb4b88ea6b40a90f19c5.
//
// Solidity: event RequestApproved(string requestId, bytes32 tokenHash, uint256 expiresAt)
func (_MedicalAccess *MedicalAccessFilterer) ParseRequestApproved(log types.Log) (*MedicalAccessRequestApproved, error) {
	event := new(MedicalAccessRequestApproved)
	if err := _MedicalAccess.contract.UnpackLog(event, "RequestApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MedicalAccessRequestConsumedIterator is returned from FilterRequestConsumed and is used to iterate over the raw logs and unpacked data for RequestConsumed events raised by the MedicalAccess contract.
type MedicalAccessRequestConsumedIterator struct {
	Event *MedicalAccessRequestConsumed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MedicalAccessRequestConsumedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MedicalAccessRequestConsumed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MedicalAccessRequestConsumed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MedicalAccessRequestConsumedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MedicalAccessRequestConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MedicalAccessRequestConsumed represents a RequestConsumed event raised by the MedicalAccess contract.
type MedicalAccessRequestConsumed struct {
	RequestId string
	Consumer  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestConsumed is a free log retrieval operation binding the contract event 0xf0b51072dcb1f05229e3efb58f6fec5511df95a8de7ff025b8735e723fb1ba4c.
//
// Solidity: event RequestConsumed(string requestId, address consumer)
func (_MedicalAccess *MedicalAccessFilterer) FilterRequestConsumed(opts *bind.FilterOpts) (*MedicalAccessRequestConsumedIterator, error) {

	logs, sub, err := _MedicalAccess.contract.FilterLogs(opts, "RequestConsumed")
	if err != nil {
		return nil, err
	}
	return &MedicalAccessRequestConsumedIterator{contract: _MedicalAccess.contract, event: "RequestConsumed", logs: logs, sub: sub}, nil
}

// WatchRequestConsumed is a free log subscription operation binding the contract event 0xf0b51072dcb1f05229e3efb58f6fec5511df95a8de7ff025b8735e723fb1ba4c.
//
// Solidity: event RequestConsumed(string requestId, address consumer)
func (_MedicalAccess *MedicalAccessFilterer) WatchRequestConsumed(opts *bind.WatchOpts, sink chan<- *MedicalAccessRequestConsumed) (event.Subscription, error) {

	logs, sub, err := _MedicalAccess.contract.WatchLogs(opts, "RequestConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MedicalAccessRequestConsumed)
				if err := _MedicalAccess.contract.UnpackLog(event, "RequestConsumed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestConsumed is a log parse operation binding the contract event 0xf0b51072dcb1f05229e3efb58f6fec5511df95a8de7ff025b8735e723fb1ba4c.
//
// Solidity: event RequestConsumed(string requestId, address consumer)
func (_MedicalAccess *MedicalAccessFilterer) ParseRequestConsumed(log types.Log) (*MedicalAccessRequestConsumed, error) {
	event := new(MedicalAccessRequestConsumed)
	if err := _MedicalAccess.contract.UnpackLog(event, "RequestConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MedicalAccessRequestCreatedIterator is returned from FilterRequestCreated and is used to iterate over the raw logs and unpacked data for RequestCreated events raised by the MedicalAccess contract.
type MedicalAccessRequestCreatedIterator struct {
	Event *MedicalAccessRequestCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MedicalAccessRequestCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MedicalAccessRequestCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MedicalAccessRequestCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MedicalAccessRequestCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MedicalAccessRequestCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MedicalAccessRequestCreated represents a RequestCreated event raised by the MedicalAccess contract.
type MedicalAccessRequestCreated struct {
	RequestId string
	RecordId  string
	Patient   string
	Clinic    string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestCreated is a free log retrieval operation binding the contract event 0xd3857f86ae3eb5b776d6168ef81042dc9e481ad045554ff37ce44ae19b8e84c7.
//
// Solidity: event RequestCreated(string requestId, string recordId, string patient, string clinic)
func (_MedicalAccess *MedicalAccessFilterer) FilterRequestCreated(opts *bind.FilterOpts) (*MedicalAccessRequestCreatedIterator, error) {

	logs, sub, err := _MedicalAccess.contract.FilterLogs(opts, "RequestCreated")
	if err != nil {
		return nil, err
	}
	return &MedicalAccessRequestCreatedIterator{contract: _MedicalAccess.contract, event: "RequestCreated", logs: logs, sub: sub}, nil
}

// WatchRequestCreated is a free log subscription operation binding the contract event 0xd3857f86ae3eb5b776d6168ef81042dc9e481ad045554ff37ce44ae19b8e84c7.
//
// Solidity: event RequestCreated(string requestId, string recordId, string patient, string clinic)
func (_MedicalAccess *MedicalAccessFilterer) WatchRequestCreated(opts *bind.WatchOpts, sink chan<- *MedicalAccessRequestCreated) (event.Subscription, error) {

	logs, sub, err := _MedicalAccess.contract.WatchLogs(opts, "RequestCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MedicalAccessRequestCreated)
				if err := _MedicalAccess.contract.UnpackLog(event, "RequestCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestCreated is a log parse operation binding the contract event 0xd3857f86ae3eb5b776d6168ef81042dc9e481ad045554ff37ce44ae19b8e84c7.
//
// Solidity: event RequestCreated(string requestId, string recordId, string patient, string clinic)
func (_MedicalAccess *MedicalAccessFilterer) ParseRequestCreated(log types.Log) (*MedicalAccessRequestCreated, error) {
	event := new(MedicalAccessRequestCreated)
	if err := _MedicalAccess.contract.UnpackLog(event, "RequestCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MedicalAccessRequestRevokedIterator is returned from FilterRequestRevoked and is used to iterate over the raw logs and unpacked data for RequestRevoked events raised by the MedicalAccess contract.
type MedicalAccessRequestRevokedIterator struct {
	Event *MedicalAccessRequestRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MedicalAccessRequestRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MedicalAccessRequestRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MedicalAccessRequestRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MedicalAccessRequestRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MedicalAccessRequestRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MedicalAccessRequestRevoked represents a RequestRevoked event raised by the MedicalAccess contract.
type MedicalAccessRequestRevoked struct {
	RequestId string
	Reason    string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestRevoked is a free log retrieval operation binding the contract event 0x35c0e0f20d08109ba2a261c71d78cdce8544e844916a4bc233048194f946eb6b.
//
// Solidity: event RequestRevoked(string requestId, string reason)
func (_MedicalAccess *MedicalAccessFilterer) FilterRequestRevoked(opts *bind.FilterOpts) (*MedicalAccessRequestRevokedIterator, error) {

	logs, sub, err := _MedicalAccess.contract.FilterLogs(opts, "RequestRevoked")
	if err != nil {
		return nil, err
	}
	return &MedicalAccessRequestRevokedIterator{contract: _MedicalAccess.contract, event: "RequestRevoked", logs: logs, sub: sub}, nil
}

// WatchRequestRevoked is a free log subscription operation binding the contract event 0x35c0e0f20d08109ba2a261c71d78cdce8544e844916a4bc233048194f946eb6b.
//
// Solidity: event RequestRevoked(string requestId, string reason)
func (_MedicalAccess *MedicalAccessFilterer) WatchRequestRevoked(opts *bind.WatchOpts, sink chan<- *MedicalAccessRequestRevoked) (event.Subscription, error) {

	logs, sub, err := _MedicalAccess.contract.WatchLogs(opts, "RequestRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MedicalAccessRequestRevoked)
				if err := _MedicalAccess.contract.UnpackLog(event, "RequestRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestRevoked is a log parse operation binding the contract event 0x35c0e0f20d08109ba2a261c71d78cdce8544e844916a4bc233048194f946eb6b.
//
// Solidity: event RequestRevoked(string requestId, string reason)
func (_MedicalAccess *MedicalAccessFilterer) ParseRequestRevoked(log types.Log) (*MedicalAccessRequestRevoked, error) {
	event := new(MedicalAccessRequestRevoked)
	if err := _MedicalAccess.contract.UnpackLog(event, "RequestRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
