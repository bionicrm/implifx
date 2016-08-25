# Implifx
Server side API for LIFX device implementations.

This API extends [Controlifx](https://github.com/golifx/controlifx), the client side API for LIFX device control. Controlifx is responsible for sending requests to LIFX devices (or machines capable of understanding the LAN protocol) and receiving responses. Implementations of Implifx are responsible for receiving requests, updating a virtual light bulb or some other hardware, and sending back responses.

Implifx simply defines methods for marshalling and unmarshalling binary into and from responses and requests, all respectively. This means that handling each message and sending back appropriate responses is up to the implementation.

**Built with Implifx:**
- [Emulifx](https://github.com/golifx/emulifx) &ndash; LIFX device emulator

**Contents:**
- [Installation](#installation)
- [Getting Started](#getting-started)

## Installation
Just run `go get -u gopkg.in/golifx/implifx.v1` to get the latest version.

## Getting Started
TODO
