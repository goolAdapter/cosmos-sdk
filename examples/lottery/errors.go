package lottery

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 7

	CodeEmptyActor    sdk.CodeType = 400
	CodeErrorStatus   sdk.CodeType = 401
	CodeErrorSequence sdk.CodeType = 402
	CodeErrorValue    sdk.CodeType = 403
)

func ErrEmptyActor(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyActor, "actor can't be empty")
}

func ErrStatusNotMatch(codespace sdk.CodespaceType, curStatus int64) sdk.Error {
	return sdk.NewError(codespace, CodeErrorStatus, fmt.Sprintf("status can't be %d", curStatus))
}

func ErrSequenceNotMatch(codespace sdk.CodespaceType, preSequence, curSequence int64) sdk.Error {
	return sdk.NewError(codespace, CodeErrorSequence, fmt.Sprintf("Sequence don't match, preSequence is %d, curSequence is %d", preSequence, curSequence))
}

func ErrInvalidValue(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeErrorValue, fmt.Sprintf("Value is invalid"))
}
