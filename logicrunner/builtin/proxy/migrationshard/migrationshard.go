//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package migrationshard

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
)

// PrototypeReference to prototype of this contract
// error checking hides in generator
var PrototypeReference, _ = insolar.NewReferenceFromBase58("1111A7FNYLZLYXYWZPbkMhCAPwV9nYrWWE7L57CtdJCj")

// MigrationShard holds proxy type
type MigrationShard struct {
	Reference insolar.Reference
	Prototype insolar.Reference
	Code      insolar.Reference
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef insolar.Reference) (*MigrationShard, error) {
	ref, ret, err := common.CurrentProxyCtx.SaveAsChild(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}

	var constructorError *foundation.Error
	resultContainer := foundation.Result{
		Returns: []interface{}{&constructorError},
	}
	err = common.CurrentProxyCtx.Deserialize(ret, &resultContainer)
	if err != nil {
		return nil, err
	}

	if resultContainer.Error != nil {
		return nil, resultContainer.Error
	}

	if constructorError != nil {
		return nil, constructorError
	}

	return &MigrationShard{Reference: *ref}, nil
}

// GetObject returns proxy object
func GetObject(ref insolar.Reference) (r *MigrationShard) {
	return &MigrationShard{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() insolar.Reference {
	return *PrototypeReference
}

// New is constructor
func New(migrationAddresses []string) *ContractConstructorHolder {
	var args [1]interface{}
	args[0] = migrationAddresses

	var argsSerialized []byte
	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *MigrationShard) GetReference() insolar.Reference {
	return r.Reference
}

// GetPrototype returns reference to the code
func (r *MigrationShard) GetPrototype() (insolar.Reference, error) {
	if r.Prototype.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 insolar.Reference
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetPrototype", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = common.CurrentProxyCtx.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Prototype = ret0
	}

	return r.Prototype, nil

}

// GetCode returns reference to the code
func (r *MigrationShard) GetCode() (insolar.Reference, error) {
	if r.Code.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 insolar.Reference
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetCode", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = common.CurrentProxyCtx.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Code = ret0
	}

	return r.Code, nil
}

// GetMigrationAddressesAmount is proxy generated method
func (r *MigrationShard) GetMigrationAddressesAmountAsMutable(migrationAddresses []string) (int, error) {
	var args [1]interface{}
	args[0] = migrationAddresses

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 int
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetMigrationAddressesAmount", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetMigrationAddressesAmountNoWait is proxy generated method
func (r *MigrationShard) GetMigrationAddressesAmountNoWait(migrationAddresses []string) error {
	var args [1]interface{}
	args[0] = migrationAddresses

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetMigrationAddressesAmount", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetMigrationAddressesAmountAsImmutable is proxy generated method
func (r *MigrationShard) GetMigrationAddressesAmount(migrationAddresses []string) (int, error) {
	var args [1]interface{}
	args[0] = migrationAddresses

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 int
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetMigrationAddressesAmount", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// AddFreeMigrationAddresses is proxy generated method
func (r *MigrationShard) AddFreeMigrationAddresses(migrationAddresses []string) error {
	var args [1]interface{}
	args[0] = migrationAddresses

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "AddFreeMigrationAddresses", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// AddFreeMigrationAddressesNoWait is proxy generated method
func (r *MigrationShard) AddFreeMigrationAddressesNoWait(migrationAddresses []string) error {
	var args [1]interface{}
	args[0] = migrationAddresses

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "AddFreeMigrationAddresses", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// AddFreeMigrationAddressesAsImmutable is proxy generated method
func (r *MigrationShard) AddFreeMigrationAddressesAsImmutable(migrationAddresses []string) error {
	var args [1]interface{}
	args[0] = migrationAddresses

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "AddFreeMigrationAddresses", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// GetFreeMigrationAddress is proxy generated method
func (r *MigrationShard) GetFreeMigrationAddress() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetFreeMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetFreeMigrationAddressNoWait is proxy generated method
func (r *MigrationShard) GetFreeMigrationAddressNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetFreeMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetFreeMigrationAddressAsImmutable is proxy generated method
func (r *MigrationShard) GetFreeMigrationAddressAsImmutable() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetFreeMigrationAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetRef is proxy generated method
func (r *MigrationShard) GetRefAsMutable(key string) (string, error) {
	var args [1]interface{}
	args[0] = key

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetRefNoWait is proxy generated method
func (r *MigrationShard) GetRefNoWait(key string) error {
	var args [1]interface{}
	args[0] = key

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetRefAsImmutable is proxy generated method
func (r *MigrationShard) GetRef(key string) (string, error) {
	var args [1]interface{}
	args[0] = key

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// SetRef is proxy generated method
func (r *MigrationShard) SetRef(ma string, ref string) error {
	var args [2]interface{}
	args[0] = ma
	args[1] = ref

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "SetRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// SetRefNoWait is proxy generated method
func (r *MigrationShard) SetRefNoWait(ma string, ref string) error {
	var args [2]interface{}
	args[0] = ma
	args[1] = ref

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "SetRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// SetRefAsImmutable is proxy generated method
func (r *MigrationShard) SetRefAsImmutable(ma string, ref string) error {
	var args [2]interface{}
	args[0] = ma
	args[1] = ref

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "SetRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}
