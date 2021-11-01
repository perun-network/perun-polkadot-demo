<h1 align="center"><br>
    <a href="https://perun.network/"><img src=".assets/logo.png" alt="Perun" width="196"></a>
<br></h1>

<h4 align="center">Polkadot Demo CLI</h4>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License: Apache 2.0"></a>
  <a href="https://github.com/perun-network/perun-polkadot-node/actions/workflows/rust.yml"><img src="https://github.com/perun-network/perun-polkadot-node/actions/workflows/rust.yml/badge.svg" alt="CI status"></a>
</p>

_perun-polkadot-demo_ allows you to interact with [perun](https://perun.network/) Payment-Channels over a CLI powered by [go-perun](https://github.com/perun-network/go-perun).  
It connects to our [Pallet] that runs on a [Polkadot Node] by using our [Polkadot Backend].

## Security Disclaimer
The authors take no responsibility for any loss of digital assets or other damage caused by the use of this software.  
**Do not use this software with real funds**.

## Getting Started

Running _perun-polkadot-demo_ requires a recent [Go installation](https://golang.org), see `go.mod` for the required version. 
```sh
# Clone the repository into a directory of your choice
git clone https://github.com/perun-network/perun-polkadot-demo
cd perun-polkadot-demo
# Compile with
go build
# Check that the binary works
./perun-polkadot-demo --help
```

## Demo

Currently, the only available sub-command of _perun-polkadot-demo_ is `demo`, which starts the CLI node. The node's
configuration file can be chosen with the `--config` flag. Two sample
configurations `alice.yaml` and `bob.yaml` are provided. A default network
configuration for Alice and Bob is provided in file `network.yaml`.

## Example Walkthrough
In a first terminal, start a development [Polkadot Node]:
```sh
docker run --rm -it -p 127.0.0.1:9944:9944/tcp ghcr.io/perun-network/polkadot-test-node:docker-push
```

In a second terminal, start the node of Alice with
```sh
./perun-polkadot-demo demo --config alice.yaml
```
and in a third terminal, start the node of Bob with
```sh
./perun-polkadot-demo demo --config bob.yaml
```

Once both CLIs are running, e.g. in Alice's terminal, propose a payment channel
to Bob with 100 *Dot* deposit from both sides via the following command.
```
> open bob 100 100
```
In Bobs terminal, accept the appearing channel proposal.
```
🔁 Incoming channel proposal from alice with funding [My: 100 Dot, Peer: 100 Dot].
Accept (y/n)? > y
```

You will see a message like `"Channel established with…"` in both terminals
after the on-chain funding completed.

Now you can execute off-chain payments, e.g. in Bob's terminal with
```
> send alice 10
```
The updated balance will immediately be printed in both terminals, but no
transaction will be visible in the ganache's terminal.

You may always check the current status with command `info`.

You can also run a performance benchmark with command
```
> benchmark alice 10 100
```
which will send 10 *Dot* in 100 micro-transactions from Bob to Alice. Transaction performance will be printed in a table.

Finally, you can settle the channel on either side with
```
> close alice
```

Now you can exit the CLI with command `exit` or `Ctrl+D`.

## Copyright

Copyright 2021 PolyCrypt GmbH.  
Copyright 2021 Chair of Applied Cryptography, Technische Universität Darmstadt, Germany.  
All rights reserved.
Use of the source code is governed by the Apache 2.0 license that can be found in the [LICENSE file](LICENSE).

Contact us at [info@perun.network](mailto:info@perun.network).

<!-- Links -->
[Pallet]: https://github.com/perun-network/perun-polkadot-pallet/
[Polkadot Backend]: https://github.com/perun-network/perun-polkadot-backend
[Polkadot Node]: https://github.com/perun-network/perun-polkadot-node
