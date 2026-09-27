//go:build js

package stdlib

import (
	"fmt"
	"net"
)

// listenPacketConn is the GOOS=js counterpart of listenpacket_other.go's
// version: there is no net.ListenPacket under wasm (browser or Node), so
// binding a UDP SOKKIT fails with an ordinary Objective-LOL exception
// instead of failing to compile.
func listenPacketConn(network, address string) (net.PacketConn, error) {
	return nil, fmt.Errorf("UDP sockets are not supported in this environment")
}
