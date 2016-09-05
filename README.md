# Implifx
Server side API for LIFX device implementations.

This API extends [Controlifx](https://github.com/lifx-tools/controlifx), the client side API for LIFX device control. Controlifx is responsible for sending requests to LIFX devices (or machines capable of understanding the LAN protocol) and receiving responses. Implementations of Implifx are responsible for receiving requests, updating a virtual light bulb or some other hardware, and sending back responses.

Implifx simply defines methods for marshalling and unmarshalling binary into and from responses and requests, all respectively. This means that handling each message and sending back appropriate responses is up to the implementation.

**Built with Implifx:**
- [Emulifx](https://github.com/lifx-tools/emulifx) &ndash; LIFX device emulator

**Contents:**
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Testing](#testing)
- [Additional Help](#additional-help)

## Installation
Just run `go get -u gopkg.in/lifx-tools/implifx.v1` to get the latest version.

## Getting Started
If you prefer a fully functioning implementation, [click here](https://github.com/lifx-tools/emulifx/blob/master/server/server.go) to see how Emulifx uses this API for emulating LIFX bulbs. Emulifx acts as if it were an actual LIFX light bulb, which is what you're presumably trying to do if you use this API.

If you had a look at Emulifx or would like to get straight to making your own implementation, here's some instructions and code to go along with it... first, just like Controlifx, you'll need to open a UDP socket for sending and receiving messages.

Note that in this case, we use `ListenOnOtherPort(...)` instead of `Listen(...)`, which uses port 56700, because we'll eventually want to test the server on our local machine. When we use [Clifx](https://github.com/lifx-tools/clifx) later on, it will bind to 0.0.0.0:56700, and so we don't want to bind to that same address here or else we'll get an error:

```go
conn, err := implifx.ListenOnOtherPort("127.0.0.1", "0")
if err != nil {
	log.Fatalln(err)
}
defer conn.Close()

log.Println("Listening @", conn.LocalAddr().String())
```

As an aside, you can set the MAC address of your server to be sent with responses with:

```go
conn.Mac = ...
```

Now that we have a connection open, we can make our main loop for receiving and processing messages. `Connection.Receive()` will block until it receives a valid message, and respectively returns the number of bytes read, the remote address of the client that sent the message, the parsed message, and an error if there was one:

```go
for {
	_, raddr, recMsg, err := conn.Receive()
	if err != nil {
		if err.(net.Error).Temporary() {
			continue
		}
		log.Fatalln(err)
	}

	// ...
}
```

At this point, you'll want to process the received message. Handling messages and acting on them is non-trivial and completely different for every implementation. If you haven't already, take a look at Emulifx to get some ideas on how you might want to go about it. With that said, we'll assume you've done what you needed with the messages (or just printed them out for testing purposes).

However, sending responses is handled by the API. How and when responses are sent is dependant upon the protocol, not really the implementation, since certain messages constitute a response or acknowledgement regardless of how the server processes them. The following code sample will show you how to do that.

The example is largely incomplete, but does work properly for messages with `AckRequired` set to `true` and `ResRequired` set to `false`. Here's the method broken down in the order of parameters:

1. This `bool` value describes whether or not a response should be sent even if it was not explicitly required in the `ResRequired` header field. For example, when a device sends a `GetLabel` request, this should be `true` since it wouldn't make sense for the server to simply not respond to a `GetLabel` request. However, for a `SetPower` request, this should be `false` since a `StatePower` response should only be sent if it was explicitly required.
2. The remote address of the sender of the message you're responding to. This is technically just the address the response will be sent to.
3. The received message that you're responding to. This is necessary in order to properly set the `Source` and `Sequence` header fields.
4. The type of message you're responding with. In this example, we're only going to respond to requests with acknowledgements.
5. The payload to be sent with the response. All responses should have payloads, and this field should never be `nil`, or else messages sent with `ResRequired` set to `true` will cause a panic. For this simple example, we'll risk it and hope we only have to respond to requests with `ResRequired` set to `false` and `AckRequired` set to `true`.

```go
if _, err := conn.Respond(false, raddr, recMsg,
	controlifx.AcknowledgementType, nil); err != nil {
	log.Fatalln(err)
}
```

**Completed example:**
```go
package main

import (
	"gopkg.in/lifx-tools/controlifx.v1"
	"gopkg.in/lifx-tools/implifx.v1"
	"log"
	"net"
)

func main() {
	conn, err := implifx.ListenOnOtherPort("127.0.0.1", "0")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	log.Println("Listening @", conn.LocalAddr().String())

	for {
		_, raddr, recMsg, err := conn.Receive()
		if err != nil {
			if err.(net.Error).Temporary() {
				continue
			}
			log.Fatalln(err)
		}

		// Do some processing here...
		log.Println("Received message of type",
			recMsg.Header.ProtocolHeader.Type, "from",
			raddr.String())

		if _, err := conn.Respond(false, raddr, recMsg,
			controlifx.AcknowledgementType, nil); err != nil {
			log.Fatalln(err)
		}
	}
}
```

## Testing
To test your implementation, use [Clifx](https://github.com/lifx-tools/clifx) on the same machine that you're developing on. The CLI will, by default, emit and receive messages on the LAN. However, you want to send and receive messages on your local machine instead.

In the following command, replace `59311` with the port that your server is listening on. If you have the same code as the completed example from [Getting Started](#getting-started), this can be found when it prints out `2016/08/25 21:38:09 Listening @ 127.0.0.1:59311` or similar. The other parts of the command can be explained by checking out the Clifx read-me.

```bash
$ clifx --broadcast-addr 127.0.0.1:59311 -a info
```

**Example output from Clifx:**
```
[{"Device":{"Addr":{"IP":"127.0.0.1","Port":59311,"Zone":""},"Mac":0}}]
```

**Example output from your server implementation:**
```
2016/08/25 21:43:38 Received message of type 34 from 127.0.0.1:56700
```

## Additional Help
Visit [#lifx-tools](http://webchat.freenode.net?randomnick=1&channels=%23lifx-tools&prompt=1) on chat.freenode.net to get help, ask questions, or discuss ideas.
