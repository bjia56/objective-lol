//go:build !js

package stdlib

import "net"

// listenPacketConn opens a UDP (packet-oriented) listener. Split out
// behind a build tag because net.ListenPacket is unimplemented for
// GOOS=js (both the standard wasm port and TinyGo), which otherwise
// breaks compiling this package for the browser.
func listenPacketConn(network, address string) (net.PacketConn, error) {
	return net.ListenPacket(network, address)
}
