# Himalia PatractGo

Substrate Contract SDK for Golang As a part of Himalia

----------

PatractGo is a Golang contract SDK. In addition to supporting the basic interactive API with the contract, it mainly supports the development of micro-services based on the contract status. For example, for common ERC20 contracts, a service can be developed based on PatractGo to synchronize all transfer information to the database, and based on the database to implemente the statistics and processing of transfer data. In addition, developers can also develop some command-line tools based on Golang to assist in testing and debugging.

PatractGo will be based on [Centrifuge's GSRPC](https://github.com/centrifuge/go-substrate-rpc-client), which is a Go sdk for Substrate.

### Intruduction

Most contract behaviors are highly related to context. In addition to interacting with the chain, user-oriented contract applications also need to provide users with current relevant context status information:

```
+--DAPP-Front-End--------------+        +---Chain-------------------------+
|                              |        |                                 |
| +----+  +------------------+ |        | +-------+     +-------+         |
| |    |  |                  | | Commit | |       |     |       |         |
| |    |  |   Polkadot-JS    +------------> Node  +---->+ Node  |         |
| |    +->+                  | |   Tx   | |       |     |       |         |
| |    |  |                  | |        | +-------+     +----+-++         |
| |    |  +------------------+ |        |                    ^ |          |
| | UI |                       |        +---------------------------------+
| |    |  +------------------+ |                             | |
| |    |  |                  | |        +--DAPP-Server--------------------+
| |    |  |                  | |  Push  | +--------+     +-----v-------+  |
| |    +<-+   Model          +<-----------+        +-----+             |  |
| |    |  |                  | |        | | Server |     |  PatractGo  |  |
| |    |  |                  +------------>        +-----+             |  |
| +----+  +------------------+ | Query  | +----+---+     +-----+-------+  |
+------------------------------+        |      |               |          |
                                        |      |         +-----v-------+  |
                                        |      |         |             |  |
                                        |      +-------->+   DataBase  |  |
                                        |                |             |  |
                                        |                +-------------+  |
                                        |                                 |
                                        +---------------------------------+
```

PatractGo is mainly responsible for implementing micro-services in a DApp. Unlike querying the state of the chain API, PatractGo can monitor the calls and events generated by the specified contract. Developers can obtain the state storage based on this information to maintain consistent state with the chain. Through data services based on a typical API-DB architecture, the front-end DApp can efficiently and concisely obtain the state on the chain as context information.

Based on the API of chain nodes, PatractGo obtains block information and summarizes and filters it, and sends contract-related messages and events based on metadata analysis to the handler protocol specified by the developer. For example, for a typical ERC20 contract, the developer can use the channel to subscribe to all transfer events that occur, and then synchronize them into the database, so that other microservices can provide services corresponding to the token data of the account, such as querying the current token holding distribution and other logics.

Therefor, PatractGo will achieve the following support:

* Complete the secondary packaging of the contract module interface, complete operations such as `put_code`, `call`, `instantiate`, etc.
* Parse the metadata.json information of the contract, and support the automatic generation of http service interface for the metadata corresponding contract
* Scanning and monitoring support of the contract status on the chain for statistics and analysis
* Basic command line tool support for native interaction with the contract, mainly used to test the security of the contract
* SDK development examples for ERC20 contract support

## Getting Start

PatractGo based on [GSRPC](https://github.com/centrifuge/go-substrate-rpc-client), So we need install some depends:

First is subkey:

```bash
cargo install --force subkey --git https://github.com/paritytech/substrate --version 2.0.0
subkey --version
```

For Now, the sdk examples will connect to the [canvas](https://github.com/paritytech/canvas-node), and also need cli tools:

```bash
cargo install canvas-node --git https://github.com/paritytech/canvas-node.git --tag v0.1.3 --force --locked
cargo install cargo-contract --vers 0.7.1 --force --locked
canvas -V
```

For some examples, we can simply run a canvas node:

```bash
canvas --dev --tmp
```

## Thanks

- [Centrifuge's GSRPC](https://github.com/centrifuge/go-substrate-rpc-client)
