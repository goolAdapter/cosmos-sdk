package lottery

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 6

	CodeEmptyActor sdk.CodeType = 400
)

func ErrEmptyActor(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyActor, "actor can't be empty.")
}
