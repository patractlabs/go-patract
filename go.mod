module github.com/patractlabs/go-patract

go 1.15

require (
	github.com/boltdb/bolt v1.3.1
	github.com/btcsuite/btcutil v1.0.2
	github.com/centrifuge/go-substrate-rpc-client/v2 v2.1.1-0.20210302023953-904cb0b931a9
	github.com/ethereum/go-ethereum v1.10.1 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/jesselucas/executil v0.0.0-20151120044647-dde271ce6a5c
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210317152858-513c2a44f670
	golang.org/x/sys v0.0.0-20210317225723-c4fcb01b228e // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.1.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/centrifuge/go-substrate-rpc-client/v2 v2.1.1-0.20210302023953-904cb0b931a9 => github.com/Snowfork/go-substrate-rpc-client/v2 v2.0.2-0.20210115165558-f6ad0aceb9bc
