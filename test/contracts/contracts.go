package contracts

import "github.com/patractlabs/go-patract/types"

var (
	CodeHashContractTerminate = types.NewHash(
		types.MustHexDecodeString("0x95646611d2b62cbe817692db6624fdac839404d21d5ec5430de72ff7f4ac775c"))
	CodeHashDNS = types.NewHash(
		types.MustHexDecodeString("0x95fbcfd248193474b11b56a8d61db8d2b1d67f5db9ababd5dfbf4836d56b2f03"))
	// CodeHashERC20 hash for erc20.wasm
	CodeHashERC20 = types.NewHash(
		types.MustHexDecodeString("0x2d8214d5a8b920fc351d2a186cd15e842f4bf4d35da8b210a46beaeabbea62bf"))
	CodeHashERC721 = types.NewHash(
		types.MustHexDecodeString("0xe403d76a74a8a434e74087f6c773d3e9e54702b1585e2a281cab41b5e02a66fb"))
	CodeHashFlipper = types.NewHash(
		types.MustHexDecodeString("0x88af2390f1eb7c7adbf9bc726e454e94507995627e1e38d6df4d5071cc274849"))
	CodeHashTraitFlipper = types.NewHash(
		types.MustHexDecodeString("0x4a13efbd917a28bceb3c47498b20072abed35bcf9131d5a2910847f4940c1e86"))

	CodeHashTraitERC20 = types.NewHash(
		types.MustHexDecodeString("0x66f7b38626227ec7c880740904ae741d3a2dc09a5103f9913a75cbe2e2e02b9f"))
)
