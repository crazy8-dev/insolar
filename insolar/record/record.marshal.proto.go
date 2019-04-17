// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: insolar/record/record.proto

package record

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_pkg_errors "github.com/pkg/errors"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Genesis X_GenesisRecord
type Child X_ChildRecord
type Jet X_JetRecord
type Request X_RequestRecord
type Result X_ResultRecord
type Type X_TypeRecord
type Code X_CodeRecord
type Activate X_ActivateRecord
type Amend X_AmendRecord
type Deactivate X_DeactivateRecord

// Returns pointer to any sub-record type or error
func UnmarshalRecord(data []byte) (Record, error) {
	base := X_Record{}

	if error := base.Unmarshal(data); error != nil {
		return nil, github_com_pkg_errors.Wrap(error, "Failed to unmarshal request")
	}

	union := base.GetUnion()

	if union == nil {
		return nil, github_com_pkg_errors.New("We go empty request")
	}

	var subRecord Record

	switch union.(type) {
	case *X_Record_Genesis:
		subRecord = (*Genesis)(((union).(*X_Record_Genesis).Genesis))
	case *X_Record_Child:
		subRecord = (*Child)(((union).(*X_Record_Child).Child))
	case *X_Record_Jet:
		subRecord = (*Jet)(((union).(*X_Record_Jet).Jet))
	case *X_Record_Request:
		subRecord = (*Request)(((union).(*X_Record_Request).Request))
	case *X_Record_Result:
		subRecord = (*Result)(((union).(*X_Record_Result).Result))
	case *X_Record_Type:
		subRecord = (*Type)(((union).(*X_Record_Type).Type))
	case *X_Record_Code:
		subRecord = (*Code)(((union).(*X_Record_Code).Code))
	case *X_Record_Activate:
		subRecord = (*Activate)(((union).(*X_Record_Activate).Activate))
	case *X_Record_Amend:
		subRecord = (*Amend)(((union).(*X_Record_Amend).Amend))
	case *X_Record_Deactivate:
		subRecord = (*Deactivate)(((union).(*X_Record_Deactivate).Deactivate))
	default:
		return nil, fmt.Errorf("__Record.union has unexpected type %T", subRecord)
	}
	return subRecord, nil
}

// Puts sub-record into record and convert it to binary (if any)
func MarshalRecord(subRecord Record) ([]byte, error) {
	base := X_Record{}

	switch subRecord.(type) {
	case *Genesis:
		base.Union = &X_Record_Genesis{(*X_GenesisRecord)(subRecord.(*Genesis))}
	case *Child:
		base.Union = &X_Record_Child{(*X_ChildRecord)(subRecord.(*Child))}
	case *Jet:
		base.Union = &X_Record_Jet{(*X_JetRecord)(subRecord.(*Jet))}
	case *Request:
		base.Union = &X_Record_Request{(*X_RequestRecord)(subRecord.(*Request))}
	case *Result:
		base.Union = &X_Record_Result{(*X_ResultRecord)(subRecord.(*Result))}
	case *Type:
		base.Union = &X_Record_Type{(*X_TypeRecord)(subRecord.(*Type))}
	case *Code:
		base.Union = &X_Record_Code{(*X_CodeRecord)(subRecord.(*Code))}
	case *Activate:
		base.Union = &X_Record_Activate{(*X_ActivateRecord)(subRecord.(*Activate))}
	case *Amend:
		base.Union = &X_Record_Amend{(*X_AmendRecord)(subRecord.(*Amend))}
	case *Deactivate:
		base.Union = &X_Record_Deactivate{(*X_DeactivateRecord)(subRecord.(*Deactivate))}
	default:
		return nil, fmt.Errorf("__Record.union has unexpected type %T", subRecord)
	}
	return base.Marshal()
}
