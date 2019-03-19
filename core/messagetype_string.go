// Code generated by "stringer -type=MessageType"; DO NOT EDIT.

package core

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeCallMethod-0]
	_ = x[TypeCallConstructor-1]
	_ = x[TypeReturnResults-2]
	_ = x[TypeExecutorResults-3]
	_ = x[TypeValidateCaseBind-4]
	_ = x[TypeValidationResults-5]
	_ = x[TypePendingFinished-6]
	_ = x[TypeStillExecuting-7]
	_ = x[TypeGetCode-8]
	_ = x[TypeGetObject-9]
	_ = x[TypeGetDelegate-10]
	_ = x[TypeGetChildren-11]
	_ = x[TypeUpdateObject-12]
	_ = x[TypeRegisterChild-13]
	_ = x[TypeJetDrop-14]
	_ = x[TypeSetRecord-15]
	_ = x[TypeValidateRecord-16]
	_ = x[TypeSetBlob-17]
	_ = x[TypeGetObjectIndex-18]
	_ = x[TypeGetPendingRequests-19]
	_ = x[TypeHotRecords-20]
	_ = x[TypeGetJet-21]
	_ = x[TypeAbandonedRequestsNotification-22]
	_ = x[TypeGetRequest-23]
	_ = x[TypeGetPendingRequestID-24]
	_ = x[TypeValidationCheck-25]
	_ = x[TypeHeavyStartStop-26]
	_ = x[TypeHeavyPayload-27]
	_ = x[TypeBootstrapRequest-28]
	_ = x[TypeNodeSignRequest-29]
}

const _MessageType_name = "TypeCallMethodTypeCallConstructorTypeReturnResultsTypeExecutorResultsTypeValidateCaseBindTypeValidationResultsTypePendingFinishedTypeStillExecutingTypeGetCodeTypeGetObjectTypeGetDelegateTypeGetChildrenTypeUpdateObjectTypeRegisterChildTypeJetDropTypeSetRecordTypeValidateRecordTypeSetBlobTypeGetObjectIndexTypeGetPendingRequestsTypeHotRecordsTypeGetJetTypeAbandonedRequestsNotificationTypeGetRequestTypeGetPendingRequestIDTypeValidationCheckTypeHeavyStartStopTypeHeavyPayloadTypeBootstrapRequestTypeNodeSignRequest"

var _MessageType_index = [...]uint16{0, 14, 33, 50, 69, 89, 110, 129, 147, 158, 171, 186, 201, 217, 234, 245, 258, 276, 287, 305, 327, 341, 351, 384, 398, 421, 440, 458, 474, 494, 513}

func (i MessageType) String() string {
	if i >= MessageType(len(_MessageType_index)-1) {
		return "MessageType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MessageType_name[_MessageType_index[i]:_MessageType_index[i+1]]
}
