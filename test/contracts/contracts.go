package contracts

import "github.com/patractlabs/go-patract/types"

var (
	// CodeHashERC20 hash for erc20.wasm
	CodeHashERC20 = types.NewHash(
		types.MustHexDecodeString("0x57f26a48169e57f118f6a0bca610a09b85cac7892f73c48f0c409cc8971817e7"))
	CodeHashFlipper = types.NewHash(
		types.MustHexDecodeString("0x0f91ab17cd254fb8a00b22396e8c1be5067686469064216e242a13140d83248e"))
)
