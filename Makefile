.PHONY:

gen-proto:
	protoc --proto_path=proto --go_out=internal/models --go_opt=paths=source_relative poolpair.proto

gen-abi:
	abigen --abi ./abi/UniswapV3/UniswapV3Factory.json --pkg v3 --type UniswapV3Factory --out internal/dex/uniswap/v3/UniswapV3Factory_abi.go
	abigen --abi ./abi/UniswapV3/UniswapV3Pool.json --pkg v3 --type UniswapV3Pool --out internal/dex/uniswap/v3/UniswapV3Pool_abi.go
	abigen --abi ./abi/UniswapV3/QuoterV2.json --pkg v3 --type QuoterV2 --out internal/dex/uniswap/v3/QuoterV2_abi.go
	abigen --abi ./abi/UniswapV2/UniswapV2Pair.json --pkg v2 --type UniswapV2Pair --out internal/dex/uniswap/v2/UniswapV2Pair_abi.go
	abigen --abi ./abi/ERC20.json --pkg erc20 --type ERC20 --out pkg/erc20/ERC20_abi.go