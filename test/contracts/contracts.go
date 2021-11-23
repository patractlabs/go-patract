package contracts

import "github.com/patractlabs/go-patract/types"

var (
	CodeHashContractTerminate = types.NewHash(
		types.MustHexDecodeString("0x95646611d2b62cbe817692db6624fdac839404d21d5ec5430de72ff7f4ac775c"))

	// CodeHashDelegator (CodeHashAccumulator and CodeHashAdder and CodeHashSubber)
	CodeHashDelegator = types.NewHash(
		types.MustHexDecodeString("0x2f424bbdbeee2f358af97ae0f105c97285f3127676e2a2528c7922922aa4721d"))
	CodeHashAccumulator = types.NewHash(
		types.MustHexDecodeString("0xf5d00b6cd1bf93773c0d111dcef790b47c741e46813e0489811251c1adfe2dba"))
	CodeHashAdder = types.NewHash(
		types.MustHexDecodeString("0xc4a25652cf653c5b2c79d69180f87684a289b15242986faf0f6846c3f30baf68"))
	CodeHashSubber = types.NewHash(
		types.MustHexDecodeString("0x162fce7372d8555653a52725d0d8b223abbad3e46719ed85ada5c892391afe36"))

	CodeHashDNS = types.NewHash(
		types.MustHexDecodeString("0x95fbcfd248193474b11b56a8d61db8d2b1d67f5db9ababd5dfbf4836d56b2f03"))
	// CodeHashERC20 hash for erc20.wasm
	CodeHashERC20 = types.NewHash(
		types.MustHexDecodeString("0x2d8214d5a8b920fc351d2a186cd15e842f4bf4d35da8b210a46beaeabbea62bf"))
	CodeHashERC721 = types.NewHash(
		types.MustHexDecodeString("0xe403d76a74a8a434e74087f6c773d3e9e54702b1585e2a281cab41b5e02a66fb"))
	CodeHashERC1155 = types.NewHash(
		types.MustHexDecodeString("0x905d2fc45938227edbd2c020697f006ba9749167f88b80e08141d04597e878ad"))
	CodeHashFlipper = types.NewHash(
		types.MustHexDecodeString("0x88af2390f1eb7c7adbf9bc726e454e94507995627e1e38d6df4d5071cc274849"))
	CodeHashIncrementer = types.NewHash(
		types.MustHexDecodeString("0xc9ecf4f718b10d3b6e4a1ba244fe3fbca3fc084baff276e925c111c78d57bf47"))
	CodeHashTraitFlipper = types.NewHash(
		types.MustHexDecodeString("0x4a13efbd917a28bceb3c47498b20072abed35bcf9131d5a2910847f4940c1e86"))
	CodeHashTraitIncrementer = types.NewHash(
		types.MustHexDecodeString("0x5ce5992896f5ba4cb10babde6e557e7842fbe2aaeb3c2d7a3a7ac83786a55d30"))
	CodeHashTraitERC20 = types.NewHash(
		types.MustHexDecodeString("0x66f7b38626227ec7c880740904ae741d3a2dc09a5103f9913a75cbe2e2e02b9f"))
)
