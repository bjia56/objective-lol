package stdlib

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleSocketCategories defines the order that categories should be rendered in documentation
var moduleSocketCategories = []string{
	"socket-creation",
	"socket-configuration",
	"tcp-server",
	"tcp-client",
	"udp-operations",
	"connection-management",
	"data-transfer",
	"connection-properties",
}

// SocketData stores the internal state of a SOKKIT
type SocketData struct {
	Host        string
	Port        int
	Protocol    string
	Timeout     time.Duration
	Listener    net.Listener   // For TCP server sockets
	PacketConn  net.PacketConn // For UDP sockets
	IsBound     bool
	IsListening bool
}

// WireData stores the internal state of a WIRE (connection)
type WireData struct {
	Conn        net.Conn
	IsConnected bool
}

// Global SOCKET class definitions - created once and reused
var socketClassesOnce = sync.Once{}
var socketClasses map[string]*environment.Class

func getSocketClasses() map[string]*environment.Class {
	socketClassesOnce.Do(func() {
		socketClasses = map[string]*environment.Class{
			"SOKKIT": {
				Name:          "SOKKIT",
				QualifiedName: "stdlib:SOCKET.SOKKIT",
				ModulePath:    "stdlib:SOCKET",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:SOCKET.SOKKIT"},
				Documentation: []string{
					"A network socket that provides TCP and UDP networking capabilities.",
					"Supports both client and server operations with configurable protocol, host, port, and timeout settings.",
					"",
					"@class SOKKIT",
					"@example TCP server setup",
					"I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT",
					"SERVER HOST ITZ \"0.0.0.0\"",
					"SERVER PORT ITZ 8080",
					"SERVER PROTOCOL ITZ \"TCP\"",
					"SERVER DO BIND",
					"SERVER DO LISTEN",
					"@example TCP client connection",
					"I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT",
					"CLIENT HOST ITZ \"127.0.0.1\"",
					"CLIENT PORT ITZ 8080",
					"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
					"@example UDP socket communication",
					"I HAS A VARIABLE UDP_SOCK TEH SOKKIT ITZ NEW SOKKIT",
					"UDP_SOCK PROTOCOL ITZ \"UDP\"",
					"UDP_SOCK PORT ITZ 9999",
					"UDP_SOCK DO BIND",
					"UDP_SOCK DO SEND_TO WIT \"Hello UDP!\" AN WIT \"localhost\" AN WIT 8888",
					"@note Default settings: TCP protocol, localhost host, port 8080, 30s timeout",
					"@note TCP sockets support BIND/LISTEN/ACCEPT for servers and CONNECT for clients",
					"@note UDP sockets support BIND/SEND_TO/RECEIVE_FROM for datagram communication",
					"@see WIRE",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"SOKKIT": {
						Name:       "SOKKIT",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Creates a new socket with default network settings.",
							"Initializes with TCP protocol, localhost host, port 8080, and 30-second timeout.",
							"",
							"@syntax NEW SOKKIT",
							"@returns {NOTHIN} No return value (constructor)",
							"@example Create default socket",
							"I HAS A VARIABLE SOCK TEH SOKKIT ITZ NEW SOKKIT",
							"BTW Socket created with TCP, localhost:8080, 30s timeout",
							"@example Create and configure socket",
							"I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT",
							"SERVER PROTOCOL ITZ \"TCP\"",
							"SERVER HOST ITZ \"0.0.0.0\" BTW Listen on all interfaces",
							"SERVER PORT ITZ 3000",
							"SERVER TIMEOUT ITZ 60 BTW 60 second timeout",
							"@note Use properties (PROTOCOL, HOST, PORT, TIMEOUT) to configure before operations",
							"@note Socket is not bound or connected after creation - use BIND or CONNECT",
							"@see BIND, CONNECT, PROTOCOL, HOST, PORT",
							"@category socket-creation",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							socketData := &SocketData{
								Host:        "localhost",
								Port:        8080,
								Protocol:    "TCP",
								IsBound:     false,
								IsListening: false,
							}
							this.NativeData = socketData

							return environment.NOTHIN, nil
						},
					},
					// BIND method
					"BIND": {
						Name:       "BIND",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Binds the socket to the configured host and port address.",
							"Prepares the socket for server operations (TCP) or datagram communication (UDP).",
							"",
							"@syntax <socket> DO BIND",
							"@returns {NOTHIN} No return value",
							"@example Bind TCP server socket",
							"I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT",
							"SERVER HOST ITZ \"0.0.0.0\" BTW Listen on all interfaces",
							"SERVER PORT ITZ 8080",
							"SERVER PROTOCOL ITZ \"TCP\"",
							"SERVER DO BIND",
							"SAYZ WIT \"Server bound to port 8080\"",
							"@example Bind UDP socket",
							"I HAS A VARIABLE UDP_SOCK TEH SOKKIT ITZ NEW SOKKIT",
							"UDP_SOCK PROTOCOL ITZ \"UDP\"",
							"UDP_SOCK PORT ITZ 9999",
							"UDP_SOCK DO BIND",
							"SAYZ WIT \"UDP socket bound to port 9999\"",
							"@example Bind with error handling",
							"MAYB",
							"    SERVER DO BIND",
							"    SAYZ WIT \"Successfully bound to port\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Failed to bind: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"@note Uses current HOST and PORT property values",
							"@note For TCP: enables LISTEN and ACCEPT operations",
							"@note For UDP: enables SEND_TO and RECEIVE_FROM operations",
							"@note Throws exception if port is already in use or insufficient permissions",
							"@see LISTEN, ACCEPT, SEND_TO, RECEIVE_FROM",
							"@category tcp-server",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "BIND: invalid socket context"}
							}

							host := socketData.Host
							port := socketData.Port
							protocol := socketData.Protocol

							socketData.Protocol = protocol
							address := fmt.Sprintf("%s:%d", host, port)

							switch protocol {
							case "TCP":
								listener, err := net.Listen("tcp", address)
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("BIND failed: %v", err)}
								}
								socketData.Listener = listener
							case "UDP":
								conn, err := net.ListenPacket("udp", address)
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("BIND failed: %v", err)}
								}
								socketData.PacketConn = conn
							default:
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("BIND: unsupported protocol: %s", protocol)}
							}

							socketData.IsBound = true
							return environment.NOTHIN, nil
						},
					},
					// LISTEN method (TCP only)
					"LISTEN": {
						Name:       "LISTEN",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Starts listening for incoming TCP connections on a bound socket.",
							"Enables the socket to accept client connections using ACCEPT method.",
							"",
							"@syntax <socket> DO LISTEN",
							"@returns {NOTHIN} No return value",
							"@example Complete TCP server setup",
							"I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT",
							"SERVER PROTOCOL ITZ \"TCP\"",
							"SERVER HOST ITZ \"0.0.0.0\"",
							"SERVER PORT ITZ 8080",
							"SERVER DO BIND",
							"SERVER DO LISTEN",
							"SAYZ WIT \"Server listening on port 8080\"",
							"@example Server with client acceptance loop",
							"SERVER DO LISTEN",
							"SAYZ WIT \"Waiting for connections...\"",
							"IM OUTTA UR LOOP",
							"    I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT",
							"    SAYZ WIT \"Client connected from \"",
							"    SAYZ WIT CLIENT REMOTE_HOST",
							"    CLIENT DO SEND WIT \"Welcome to server!\"",
							"    CLIENT DO CLOSE",
							"KTHX",
							"@note Only works with TCP protocol sockets",
							"@note Socket must be bound before calling LISTEN",
							"@note After LISTEN, use ACCEPT to handle incoming connections",
							"@note Does not block - ACCEPT is where blocking occurs",
							"@see BIND, ACCEPT",
							"@category tcp-server",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "LISTEN: invalid socket context"}
							}

							if !socketData.IsBound {
								return environment.NOTHIN, runtime.Exception{Message: "LISTEN: socket not bound"}
							}

							if socketData.Protocol != "TCP" {
								return environment.NOTHIN, runtime.Exception{Message: "LISTEN: only supported for TCP sockets"}
							}

							socketData.IsListening = true
							return environment.NOTHIN, nil
						},
					},
					// ACCEPT method (TCP only)
					"ACCEPT": {
						Name:       "ACCEPT",
						ReturnType: "WIRE",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Accepts an incoming TCP client connection on a listening socket.",
							"Blocks until a client connects, then returns a WIRE connection object.",
							"",
							"@syntax <socket> DO ACCEPT",
							"@returns {WIRE} Connection object for communicating with the client",
							"@example Simple echo server",
							"I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT",
							"SERVER PORT ITZ 8080",
							"SERVER DO BIND",
							"SERVER DO LISTEN",
							"SAYZ WIT \"Echo server listening...\"",
							"I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT",
							"SAYZ WIT \"Client connected!\"",
							"I HAS A VARIABLE MSG TEH STRIN ITZ CLIENT DO RECEIVE WIT 1024",
							"CLIENT DO SEND WIT \"Echo: \" MOAR MSG",
							"CLIENT DO CLOSE",
							"@example Multi-client server loop",
							"SERVER DO LISTEN",
							"IM OUTTA UR LOOP",
							"    I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT",
							"    SAYZ WIT \"New client from \"",
							"    SAYZ WIT CLIENT REMOTE_HOST",
							"    SAYZ WIT \":\"",
							"    SAYZ WIT CLIENT REMOTE_PORT",
							"    BTW Handle client in separate thread/process",
							"    CLIENT DO SEND WIT \"Hello from server!\"",
							"    CLIENT DO CLOSE",
							"KTHX",
							"@note Blocks execution until a client connects",
							"@note Socket must be in LISTEN state before calling ACCEPT",
							"@note Returns new WIRE object for each accepted connection",
							"@note Use returned WIRE object for data transfer with client",
							"@see LISTEN, WIRE",
							"@category tcp-server",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "ACCEPT: invalid socket context"}
							}

							if !socketData.IsListening {
								return environment.NOTHIN, runtime.Exception{Message: "ACCEPT: socket not listening"}
							}

							conn, err := socketData.Listener.Accept()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("ACCEPT failed: %v", err)}
							}

							return createWireInstance(conn)
						},
					},
					// CONNECT method
					"CONNECT": {
						Name:       "CONNECT",
						ReturnType: "WIRE",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Connects to a remote TCP server using configured HOST and PORT.",
							"Returns a WIRE connection object for data transfer, respects TIMEOUT setting.",
							"",
							"@syntax <socket> DO CONNECT",
							"@returns {WIRE} Connection object for communicating with the server",
							"@example Connect to TCP server",
							"I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT",
							"CLIENT HOST ITZ \"127.0.0.1\"",
							"CLIENT PORT ITZ 8080",
							"CLIENT TIMEOUT ITZ 10 BTW 10 second timeout",
							"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
							"SAYZ WIT \"Connected to server!\"",
							"@example HTTP client example",
							"CLIENT HOST ITZ \"httpbin.org\"",
							"CLIENT PORT ITZ 80",
							"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
							"CONN DO SEND WIT \"GET / HTTP/1.1\\r\\nHost: httpbin.org\\r\\n\\r\\n\"",
							"I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONN DO RECEIVE WIT 4096",
							"SAYZ WIT RESPONSE",
							"CONN DO CLOSE",
							"@example Connect with error handling",
							"MAYB",
							"    I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
							"    SAYZ WIT \"Successfully connected\"",
							"    BTW Use connection here",
							"    CONN DO CLOSE",
							"OOPSIE ERR",
							"    SAYZ WIT \"Connection failed: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"@note Only works with TCP protocol sockets",
							"@note Uses current HOST, PORT, and TIMEOUT property values",
							"@note Throws exception if connection fails or times out",
							"@note Returns WIRE object - same type as ACCEPT returns",
							"@see WIRE, HOST, PORT, TIMEOUT",
							"@category tcp-client",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							// Get connection parameters
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CONNECT: invalid socket context"}
							}

							host := socketData.Host
							port := socketData.Port
							protocol := strings.ToUpper(socketData.Protocol)
							timeout := time.Duration(socketData.Timeout) * time.Second

							if protocol != "TCP" {
								return environment.NOTHIN, runtime.Exception{Message: "CONNECT: only supported for TCP sockets"}
							}

							address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
							conn, err := net.DialTimeout("tcp", address, timeout)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("CONNECT failed: %v", err)}
							}

							return createWireInstance(conn)
						},
					},
					// SEND_TO method (UDP only)
					"SEND_TO": {
						Name: "SEND_TO",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
							{Name: "host", Type: "STRIN"},
							{Name: "port", Type: "INTEGR"},
						},
						Documentation: []string{
							"Sends data to a specific UDP address without establishing a connection.",
							"Used for UDP datagram communication - data is sent directly to the target.",
							"",
							"@syntax <socket> DO SEND_TO WIT <data> AN WIT <host> AN WIT <port>",
							"@param {STRIN} data - The data to send",
							"@param {STRIN} host - Target host address (IP or hostname)",
							"@param {INTEGR} port - Target port number",
							"@returns {NOTHIN} No return value",
							"@example UDP client sending data",
							"I HAS A VARIABLE UDP_CLIENT TEH SOKKIT ITZ NEW SOKKIT",
							"UDP_CLIENT PROTOCOL ITZ \"UDP\"",
							"UDP_CLIENT PORT ITZ 0 BTW Use any available port",
							"UDP_CLIENT DO BIND",
							"UDP_CLIENT DO SEND_TO WIT \"Hello UDP!\" AN WIT \"127.0.0.1\" AN WIT 9999",
							"SAYZ WIT \"UDP message sent\"",
							"@example Send to multiple destinations",
							"I HAS A VARIABLE SERVERS TEH BUKKIT ITZ NEW BUKKIT",
							"SERVERS DO PUSH WIT \"192.168.1.100\"",
							"SERVERS DO PUSH WIT \"192.168.1.101\"",
							"IM OUTTA UR SERVERS NERFIN SERVER_IP",
							"    UDP_CLIENT DO SEND_TO WIT \"Broadcast message\" AN WIT SERVER_IP AN WIT 8888",
							"IM IN UR SERVERS",
							"@example UDP logging client",
							"UDP_CLIENT DO SEND_TO WIT \"ERROR: Something went wrong\" AN WIT \"log-server.local\" AN WIT 514",
							"@note Only works with UDP protocol sockets",
							"@note Socket must be bound before sending",
							"@note No connection establishment - fire and forget",
							"@note No delivery guarantee - UDP is unreliable",
							"@see RECEIVE_FROM, BIND",
							"@category udp-operations",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "SEND_TO: invalid socket context"}
							}

							if socketData.Protocol != "UDP" {
								return environment.NOTHIN, runtime.Exception{Message: "SEND_TO: only supported for UDP sockets"}
							}

							if !socketData.IsBound {
								return environment.NOTHIN, runtime.Exception{Message: "SEND_TO: socket not bound"}
							}

							data := string(args[0].(environment.StringValue))
							host := string(args[1].(environment.StringValue))
							port := int(args[2].(environment.IntegerValue))

							address := fmt.Sprintf("%s:%d", host, port)
							addr, err := net.ResolveUDPAddr("udp", address)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("SEND_TO: invalid address: %v", err)}
							}

							_, err = socketData.PacketConn.WriteTo([]byte(data), addr)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("SEND_TO failed: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// RECEIVE_FROM method (UDP only)
					"RECEIVE_FROM": {
						Name:       "RECEIVE_FROM",
						ReturnType: "BASKIT",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Receives UDP data and returns both the data and sender information.",
							"Blocks until data is received, returns BASKIT with DATA, HOST, and PORT keys.",
							"",
							"@syntax <socket> DO RECEIVE_FROM",
							"@returns {BASKIT} Map containing DATA (received string), HOST (sender IP), PORT (sender port)",
							"@example UDP server receiving data",
							"I HAS A VARIABLE UDP_SERVER TEH SOKKIT ITZ NEW SOKKIT",
							"UDP_SERVER PROTOCOL ITZ \"UDP\"",
							"UDP_SERVER PORT ITZ 9999",
							"UDP_SERVER DO BIND",
							"SAYZ WIT \"UDP server listening on port 9999\"",
							"I HAS A VARIABLE PACKET TEH BASKIT ITZ UDP_SERVER DO RECEIVE_FROM",
							"SAYZ WIT \"Received: \"",
							"SAYZ WIT PACKET DO GET WIT \"DATA\"",
							"SAYZ WIT \"From: \"",
							"SAYZ WIT PACKET DO GET WIT \"HOST\"",
							"SAYZ WIT \":\"",
							"SAYZ WIT PACKET DO GET WIT \"PORT\"",
							"@example Echo UDP server",
							"IM OUTTA UR LOOP",
							"    I HAS A VARIABLE PACKET TEH BASKIT ITZ UDP_SERVER DO RECEIVE_FROM",
							"    I HAS A VARIABLE MSG TEH STRIN ITZ PACKET DO GET WIT \"DATA\"",
							"    I HAS A VARIABLE CLIENT_HOST TEH STRIN ITZ PACKET DO GET WIT \"HOST\"",
							"    I HAS A VARIABLE CLIENT_PORT TEH INTEGR ITZ PACKET DO GET WIT \"PORT\"",
							"    UDP_SERVER DO SEND_TO WIT \"Echo: \" MOAR MSG AN WIT CLIENT_HOST AN WIT CLIENT_PORT",
							"KTHX",
							"@example Process UDP packets",
							"I HAS A VARIABLE PACKET TEH BASKIT ITZ UDP_SERVER DO RECEIVE_FROM",
							"IZ (PACKET DO GET WIT \"DATA\") SAEM AS \"PING\"?",
							"    UDP_SERVER DO SEND_TO WIT \"PONG\" AN WIT (PACKET DO GET WIT \"HOST\") AN WIT (PACKET DO GET WIT \"PORT\")",
							"KTHX",
							"@note Only works with UDP protocol sockets",
							"@note Socket must be bound before receiving",
							"@note Blocks execution until data arrives",
							"@note Maximum packet size is 4096 bytes",
							"@see SEND_TO, BIND",
							"@category udp-operations",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "RECEIVE_FROM: invalid socket context"}
							}

							if socketData.Protocol != "UDP" {
								return environment.NOTHIN, runtime.Exception{Message: "RECEIVE_FROM: only supported for UDP sockets"}
							}

							if !socketData.IsBound {
								return environment.NOTHIN, runtime.Exception{Message: "RECEIVE_FROM: socket not bound"}
							}

							buffer := make([]byte, 4096)
							n, addr, err := socketData.PacketConn.ReadFrom(buffer)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("RECEIVE_FROM failed: %v", err)}
							}

							// Parse address
							udpAddr := addr.(*net.UDPAddr)
							host := udpAddr.IP.String()
							port := udpAddr.Port

							// Create result BASKIT
							result := NewBaskitInstance()
							baskitMap := result.NativeData.(BaskitMap)
							baskitMap["DATA"] = environment.StringValue(string(buffer[:n]))
							baskitMap["HOST"] = environment.StringValue(host)
							baskitMap["PORT"] = environment.IntegerValue(port)

							return result, nil
						},
					},
					// CLOSE method
					"CLOSE": {
						Name:       "CLOSE",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Closes the socket and releases all associated network resources.",
							"Stops listening, closes connections, and frees system resources.",
							"",
							"@syntax <socket> DO CLOSE",
							"@returns {NOTHIN} No return value",
							"@example Proper server shutdown",
							"I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT",
							"SERVER DO BIND",
							"SERVER DO LISTEN",
							"BTW Server operations here",
							"SERVER DO CLOSE",
							"SAYZ WIT \"Server shut down\"",
							"@example Cleanup in exception handling",
							"MAYB",
							"    SERVER DO BIND",
							"    SERVER DO LISTEN",
							"    BTW Server work here",
							"OOPSIE ERR",
							"    SAYZ WIT \"Server error: \"",
							"    SAYZ WIT ERR",
							"FINALLY",
							"    SERVER DO CLOSE BTW Always cleanup",
							"KTHX",
							"@note Safe to call multiple times - no error if already closed",
							"@note Automatically called when socket object is garbage collected",
							"@note For TCP servers: stops accepting new connections",
							"@note For UDP sockets: stops receiving packets",
							"@note Does not close existing WIRE connections from ACCEPT",
							"@see BIND, LISTEN, WIRE",
							"@category connection-management",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CLOSE: invalid socket context"}
							}

							if socketData.Listener != nil {
								socketData.Listener.Close()
								socketData.Listener = nil
							}
							if socketData.PacketConn != nil {
								socketData.PacketConn.Close()
								socketData.PacketConn = nil
							}

							socketData.IsBound = false
							socketData.IsListening = false

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"PROTOCOL": {
						Variable: environment.Variable{
							Name:     "PROTOCOL",
							Type:     "STRIN",
							IsLocked: false,
							IsPublic: true,
							Documentation: []string{
								"Socket protocol specification - either TCP or UDP.",
								"",
								"@property {STRIN} PROTOCOL - Network protocol (\"TCP\" or \"UDP\", default: \"TCP\")",
								"@example Set TCP protocol (default)",
								"I HAS A VARIABLE SOCK TEH SOKKIT ITZ NEW SOKKIT",
								"SOCK PROTOCOL ITZ \"TCP\"",
								"SAYZ WIT \"Using TCP protocol\"",
								"@example Set UDP protocol",
								"SOCK PROTOCOL ITZ \"UDP\"",
								"SAYZ WIT \"Using UDP protocol\"",
								"@example Check current protocol",
								"IZ (SOCK PROTOCOL) SAEM AS \"TCP\"?",
								"    SAYZ WIT \"Socket is configured for TCP\"",
								"NOPE",
								"    SAYZ WIT \"Socket is configured for UDP\"",
								"KTHX",
								"@note TCP provides reliable, ordered, connection-based communication",
								"@note UDP provides fast, connectionless, datagram-based communication",
								"@note Must be set before BIND operation",
								"@note Case-insensitive (converted to uppercase)",
								"@see BIND, CONNECT, SEND_TO",
								"@category socket-configuration",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return nil, runtime.Exception{Message: "PROTOCOL: invalid socket context"}
							}
							return environment.StringValue(socketData.Protocol), nil
						},
						NativeSet: func(this *environment.ObjectInstance, value environment.Value) error {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return runtime.Exception{Message: "PROTOCOL: invalid socket context"}
							}
							if strVal, err := value.Cast("STRIN"); err == nil {
								proto := strings.ToUpper(string(strVal.(environment.StringValue)))
								if proto != "TCP" && proto != "UDP" {
									return runtime.Exception{Message: "PROTOCOL must be 'TCP' or 'UDP'"}
								}
								socketData.Protocol = proto
								return nil
							}
							return runtime.Exception{Message: "PROTOCOL must be a STRIN"}
						},
					},
					"HOST": {
						Variable: environment.Variable{
							Name:     "HOST",
							Type:     "STRIN",
							IsLocked: false,
							IsPublic: true,
							Documentation: []string{
								"Target host address for network operations.",
								"",
								"@property {STRIN} HOST - Host address (IP or hostname, default: \"localhost\")",
								"@example Set specific IP address",
								"I HAS A VARIABLE SOCK TEH SOKKIT ITZ NEW SOKKIT",
								"SOCK HOST ITZ \"192.168.1.100\"",
								"SAYZ WIT \"Connecting to \"",
								"SAYZ WIT SOCK HOST",
								"@example Set hostname",
								"SOCK HOST ITZ \"example.com\"",
								"@example Server listening on all interfaces",
								"I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT",
								"SERVER HOST ITZ \"0.0.0.0\" BTW Listen on all network interfaces",
								"SERVER PORT ITZ 8080",
								"SERVER DO BIND",
								"@example Localhost connections only",
								"SERVER HOST ITZ \"127.0.0.1\" BTW Local connections only",
								"@note For servers: determines which network interface to bind to",
								"@note For clients: determines which host to connect to",
								"@note Use \"0.0.0.0\" to listen on all interfaces (servers only)",
								"@note Use \"127.0.0.1\" or \"localhost\" for local-only connections",
								"@see PORT, BIND, CONNECT",
								"@category socket-configuration",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return nil, runtime.Exception{Message: "HOST: invalid socket context"}
							}
							return environment.StringValue(socketData.Host), nil
						},
						NativeSet: func(this *environment.ObjectInstance, value environment.Value) error {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return runtime.Exception{Message: "HOST: invalid socket context"}
							}
							if strVal, err := value.Cast("STRIN"); err == nil {
								socketData.Host = string(strVal.(environment.StringValue))
								return nil
							}
							return runtime.Exception{Message: "HOST must be a STRIN"}
						},
					},
					"PORT": {
						Variable: environment.Variable{
							Name:     "PORT",
							Type:     "INTEGR",
							IsLocked: false,
							IsPublic: true,
							Documentation: []string{
								"Target port number for network operations.",
								"",
								"@property {INTEGR} PORT - Port number (default: 8080, valid range: 0-65535)",
								"@example Set web server port",
								"I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT",
								"SERVER PORT ITZ 80 BTW HTTP default port",
								"@example Set custom port",
								"SERVER PORT ITZ 3000",
								"@example Use ephemeral port (system assigns)",
								"I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT",
								"CLIENT PORT ITZ 0 BTW System will assign available port",
								"CLIENT DO BIND",
								"@example Common port numbers",
								"BTW SERVER PORT ITZ 21    BTW FTP",
								"BTW SERVER PORT ITZ 22    BTW SSH",
								"BTW SERVER PORT ITZ 80    BTW HTTP",
								"BTW SERVER PORT ITZ 443   BTW HTTPS",
								"BTW SERVER PORT ITZ 993   BTW IMAPS",
								"@note Ports 0-1023 are reserved and may require admin privileges",
								"@note Port 0 means \"assign any available port\" when binding",
								"@note Port must be in range 0-65535",
								"@note Common ports: 80 (HTTP), 443 (HTTPS), 22 (SSH), 21 (FTP)",
								"@see HOST, BIND, CONNECT",
								"@category socket-configuration",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return nil, runtime.Exception{Message: "PORT: invalid socket context"}
							}
							return environment.IntegerValue(socketData.Port), nil
						},
						NativeSet: func(this *environment.ObjectInstance, value environment.Value) error {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return runtime.Exception{Message: "PORT: invalid socket context"}
							}
							if intVal, err := value.Cast("INTEGR"); err == nil {
								intVal := int(intVal.(environment.IntegerValue))
								if intVal < 0 || intVal > 65535 {
									return runtime.Exception{Message: "PORT must be between 0 and 65535"}
								}
								socketData.Port = int(intVal)
								return nil
							}
							return runtime.Exception{Message: "PORT must be an INTEGER"}
						},
					},
					"TIMEOUT": {
						Variable: environment.Variable{
							Name:     "TIMEOUT",
							Type:     "INTEGR",
							IsLocked: false,
							IsPublic: true,
							Documentation: []string{
								"Connection timeout in seconds for TCP client connections.",
								"",
								"@property {INTEGR} TIMEOUT - Connection timeout in seconds (default: 30)",
								"@example Set short timeout",
								"I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT",
								"CLIENT HOST ITZ \"slow-server.com\"",
								"CLIENT TIMEOUT ITZ 5 BTW 5 second timeout",
								"MAYB",
								"    I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
								"    SAYZ WIT \"Connected successfully\"",
								"OOPSIE ERR",
								"    SAYZ WIT \"Connection timed out or failed\"",
								"KTHX",
								"@example Set long timeout for slow connections",
								"CLIENT TIMEOUT ITZ 120 BTW 2 minute timeout",
								"@example Disable timeout (wait indefinitely)",
								"CLIENT TIMEOUT ITZ 0 BTW No timeout",
								"@example Check current timeout",
								"SAYZ WIT \"Connection timeout: \"",
								"SAYZ WIT CLIENT TIMEOUT",
								"SAYZ WIT \" seconds\"",
								"@note Only affects TCP CONNECT operations, not UDP or server operations",
								"@note Timeout of 0 means wait indefinitely",
								"@note Must be non-negative integer",
								"@note Default is 30 seconds for new sockets",
								"@see CONNECT, HOST, PORT",
								"@category socket-configuration",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return nil, runtime.Exception{Message: "TIMEOUT: invalid socket context"}
							}
							return environment.IntegerValue(int64(socketData.Timeout.Seconds())), nil
						},
						NativeSet: func(this *environment.ObjectInstance, value environment.Value) error {
							socketData, ok := this.NativeData.(*SocketData)
							if !ok {
								return runtime.Exception{Message: "TIMEOUT: invalid socket context"}
							}
							if intVal, err := value.Cast("INTEGR"); err == nil {
								intVal := int(intVal.(environment.IntegerValue))
								if intVal < 0 {
									return runtime.Exception{Message: "TIMEOUT must be non-negative"}
								}
								socketData.Timeout = time.Duration(intVal) * time.Second
								return nil
							}
							return runtime.Exception{Message: "TIMEOUT must be an INTEGER"}
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"WIRE": {
				Name:          "WIRE",
				QualifiedName: "stdlib:SOCKET.WIRE",
				ModulePath:    "stdlib:SOCKET",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:SOCKET.WIRE"},
				Documentation: []string{
					"A TCP connection that provides bidirectional data transfer capabilities.",
					"Represents an active network connection between two endpoints for reliable data exchange.",
					"",
					"@class WIRE",
					"@example Client connection usage",
					"I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT",
					"CLIENT HOST ITZ \"127.0.0.1\"",
					"CLIENT PORT ITZ 8080",
					"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
					"CONN DO SEND WIT \"GET /api/data HTTP/1.1\\r\\n\\r\\n\"",
					"I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONN DO RECEIVE WIT 1024",
					"SAYZ WIT RESPONSE",
					"CONN DO CLOSE",
					"@example Server-side connection handling",
					"BTW From server ACCEPT",
					"I HAS A VARIABLE CLIENT_CONN TEH WIRE ITZ SERVER DO ACCEPT",
					"SAYZ WIT \"Client connected from \"",
					"SAYZ WIT CLIENT_CONN REMOTE_HOST",
					"CLIENT_CONN DO SEND WIT \"Welcome to server!\"",
					"I HAS A VARIABLE REQUEST TEH STRIN ITZ CLIENT_CONN DO RECEIVE WIT 512",
					"CLIENT_CONN DO CLOSE",
					"@example Bidirectional communication",
					"CONN DO SEND WIT \"HELLO\"",
					"I HAS A VARIABLE REPLY TEH STRIN ITZ CONN DO RECEIVE WIT 100",
					"IZ REPLY SAEM AS \"OK\"?",
					"    CONN DO SEND WIT \"DATA: important message\"",
					"KTHX",
					"@note Created by SOKKIT.CONNECT or SOKKIT.ACCEPT methods",
					"@note Provides reliable, ordered, bidirectional communication",
					"@note Connection must be closed when finished to free resources",
					"@see SOKKIT",
				},
				PublicFunctions: map[string]*environment.Function{
					// SEND method
					"SEND": {
						Name: "SEND",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
						},
						Documentation: []string{
							"Sends string data over the TCP connection to the remote endpoint.",
							"Data is transmitted immediately and may be buffered by the network stack.",
							"",
							"@syntax <wire> DO SEND WIT <data>",
							"@param {STRIN} data - The string data to send",
							"@returns {NOTHIN} No return value",
							"@example Send simple message",
							"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
							"CONN DO SEND WIT \"Hello, Server!\"",
							"@example Send HTTP request",
							"CONN DO SEND WIT \"GET /api/users HTTP/1.1\\r\\n\"",
							"CONN DO SEND WIT \"Host: api.example.com\\r\\n\"",
							"CONN DO SEND WIT \"Content-Length: 0\\r\\n\\r\\n\"",
							"@example Send JSON data",
							"I HAS A VARIABLE JSON_DATA TEH STRIN ITZ \"{\\\"name\\\":\\\"Alice\\\",\\\"age\\\":25}\"",
							"CONN DO SEND WIT \"POST /users HTTP/1.1\\r\\n\"",
							"CONN DO SEND WIT \"Content-Type: application/json\\r\\n\"",
							"CONN DO SEND WIT \"Content-Length: \"",
							"CONN DO SEND WIT JSON_DATA SIZ",
							"CONN DO SEND WIT \"\\r\\n\\r\\n\"",
							"CONN DO SEND WIT JSON_DATA",
							"@note Connection must be established (IS_CONNECTED = YEZ)",
							"@note Data is sent as-is - no automatic newlines or formatting added",
							"@note Large data may be sent in multiple network packets",
							"@note Throws exception if connection is broken or closed",
							"@see RECEIVE, RECEIVE_ALL, IS_CONNECTED",
							"@category data-transfer",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok || !wireData.IsConnected {
								return environment.NOTHIN, runtime.Exception{Message: "SEND: connection not established"}
							}

							data := string(args[0].(environment.StringValue))
							_, err := wireData.Conn.Write([]byte(data))
							if err != nil {
								wireData.IsConnected = false
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("SEND failed: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// RECEIVE method
					"RECEIVE": {
						Name:       "RECEIVE",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "length", Type: "INTEGR"},
						},
						Documentation: []string{
							"Receives up to the specified number of characters from the TCP connection.",
							"Blocks until data is available, returns received data (may be shorter than requested).",
							"",
							"@syntax <wire> DO RECEIVE WIT <length>",
							"@param {INTEGR} length - Maximum number of characters to receive",
							"@returns {STRIN} The received data (may be shorter than requested)",
							"@example Receive fixed amount of data",
							"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
							"I HAS A VARIABLE DATA TEH STRIN ITZ CONN DO RECEIVE WIT 1024",
							"SAYZ WIT \"Received: \"",
							"SAYZ WIT DATA",
							"@example Receive HTTP response",
							"CONN DO SEND WIT \"GET / HTTP/1.1\\r\\nHost: example.com\\r\\n\\r\\n\"",
							"I HAS A VARIABLE HEADERS TEH STRIN ITZ CONN DO RECEIVE WIT 4096",
							"SAYZ WIT \"Response headers: \"",
							"SAYZ WIT HEADERS",
							"@example Receive data in chunks",
							"I HAS A VARIABLE BUFFER TEH STRIN ITZ \"\"",
							"WHILE NO SAEM AS (BUFFER ENDS WIT \"END\")",
							"    I HAS A VARIABLE CHUNK TEH STRIN ITZ CONN DO RECEIVE WIT 256",
							"    IZ CHUNK SAEM AS \"\"?",
							"        OUTTA HERE BTW Connection closed",
							"    KTHX",
							"    BUFFER ITZ BUFFER MOAR CHUNK",
							"KTHX",
							"@note Blocks execution until data arrives or connection closes",
							"@note Returns empty string if connection is closed by remote end",
							"@note May return less data than requested if that's all that's available",
							"@note Connection must be established before receiving",
							"@see SEND, RECEIVE_ALL, IS_CONNECTED",
							"@category data-transfer",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok || !wireData.IsConnected {
								return environment.NOTHIN, runtime.Exception{Message: "RECEIVE: connection not established"}
							}

							length := int(args[0].(environment.IntegerValue))
							if length <= 0 {
								return environment.StringValue(""), nil
							}

							buffer := make([]byte, length)
							n, err := wireData.Conn.Read(buffer)
							if err != nil {
								if err.Error() == "EOF" {
									wireData.IsConnected = false
								}
								return environment.StringValue(string(buffer[:n])), nil
							}

							return environment.StringValue(string(buffer[:n])), nil
						},
					},
					// RECEIVE_ALL method
					"RECEIVE_ALL": {
						Name:       "RECEIVE_ALL",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Receives all available data from the TCP connection until the connection closes.",
							"Blocks until the remote end closes the connection, then returns all received data.",
							"",
							"@syntax <wire> DO RECEIVE_ALL",
							"@returns {STRIN} All data received from the connection",
							"@example Download entire web page",
							"I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT",
							"CLIENT HOST ITZ \"httpbin.org\"",
							"CLIENT PORT ITZ 80",
							"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
							"CONN DO SEND WIT \"GET /get HTTP/1.1\\r\\nHost: httpbin.org\\r\\nConnection: close\\r\\n\\r\\n\"",
							"I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONN DO RECEIVE_ALL",
							"SAYZ WIT \"Full response: \"",
							"SAYZ WIT RESPONSE",
							"@example Receive complete file transfer",
							"CONN DO SEND WIT \"GET_FILE document.txt\"",
							"I HAS A VARIABLE FILE_CONTENT TEH STRIN ITZ CONN DO RECEIVE_ALL",
							"SAYZ WIT \"File received, size: \"",
							"SAYZ WIT FILE_CONTENT SIZ",
							"@example Simple protocol with termination",
							"CONN DO SEND WIT \"DOWNLOAD_DATA\"",
							"I HAS A VARIABLE ALL_DATA TEH STRIN ITZ CONN DO RECEIVE_ALL",
							"BTW Server closes connection when done sending",
							"SAYZ WIT \"Received all data: \"",
							"SAYZ WIT ALL_DATA",
							"@note Blocks until remote end closes the connection",
							"@note Connection is automatically marked as closed after completion",
							"@note Useful for protocols where server closes connection to signal end",
							"@note May use significant memory for large data transfers",
							"@see RECEIVE, SEND, CLOSE",
							"@category data-transfer",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok || !wireData.IsConnected {
								return environment.NOTHIN, runtime.Exception{Message: "RECEIVE_ALL: connection not established"}
							}

							var result []byte
							buffer := make([]byte, 4096)

							for {
								n, err := wireData.Conn.Read(buffer)
								if n > 0 {
									result = append(result, buffer[:n]...)
								}
								if err != nil {
									break
								}
							}

							wireData.IsConnected = false
							return environment.StringValue(string(result)), nil
						},
					},
					// CLOSE method
					"CLOSE": {
						Name:       "CLOSE",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Closes the TCP connection and releases associated resources.",
							"Connection becomes unusable after closing and IS_CONNECTED becomes NO.",
							"",
							"@syntax <wire> DO CLOSE",
							"@returns {NOTHIN} No return value",
							"@example Proper connection cleanup",
							"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
							"CONN DO SEND WIT \"Hello, Server!\"",
							"I HAS A VARIABLE REPLY TEH STRIN ITZ CONN DO RECEIVE WIT 100",
							"CONN DO CLOSE",
							"SAYZ WIT \"Connection closed\"",
							"@example Connection cleanup in exception handling",
							"MAYB",
							"    I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
							"    CONN DO SEND WIT \"Important data\"",
							"    I HAS A VARIABLE RESULT TEH STRIN ITZ CONN DO RECEIVE WIT 500",
							"    BTW Process result here",
							"OOPSIE ERR",
							"    SAYZ WIT \"Connection error: \"",
							"    SAYZ WIT ERR",
							"FINALLY",
							"    IZ CONN IS_CONNECTED?",
							"        CONN DO CLOSE BTW Always cleanup",
							"    KTHX",
							"KTHX",
							"@example Server client handling",
							"I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT",
							"CLIENT DO SEND WIT \"Welcome!\"",
							"I HAS A VARIABLE REQUEST TEH STRIN ITZ CLIENT DO RECEIVE WIT 1024",
							"BTW Process request and send response",
							"CLIENT DO CLOSE BTW Close client connection",
							"@note Safe to call multiple times - no error if already closed",
							"@note Connection cannot be reused after closing",
							"@note Always close connections to prevent resource leaks",
							"@note IS_CONNECTED property becomes NO after closing",
							"@see IS_CONNECTED, SEND, RECEIVE",
							"@category connection-management",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CLOSE: invalid connection context"}
							}

							if wireData.Conn != nil {
								wireData.Conn.Close()
								wireData.IsConnected = false
							}

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"REMOTE_HOST": {
						Variable: environment.Variable{
							Name:     "REMOTE_HOST",
							Type:     "STRIN",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property containing the remote endpoint's IP address.",
								"",
								"@property {STRIN} REMOTE_HOST - IP address of the connected remote host",
								"@example Check client connection details",
								"I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT",
								"SAYZ WIT \"Client connected from: \"",
								"SAYZ WIT CLIENT REMOTE_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT REMOTE_PORT",
								"@example Log connection information",
								"SAYZ WIT \"New connection: \"",
								"SAYZ WIT CLIENT REMOTE_HOST",
								"SAYZ WIT \" -> \"",
								"SAYZ WIT CLIENT LOCAL_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT LOCAL_PORT",
								"@example Access control by IP",
								"IZ (CLIENT REMOTE_HOST) SAEM AS \"127.0.0.1\"?",
								"    SAYZ WIT \"Local connection allowed\"",
								"NOPE",
								"    SAYZ WIT \"Remote connection from \"",
								"    SAYZ WIT CLIENT REMOTE_HOST",
								"    CLIENT DO CLOSE BTW Block external connections",
								"KTHX",
								"@note Returns empty string if connection is not established",
								"@note Shows actual IP address, not hostname",
								"@note IPv4 addresses shown in dotted decimal notation (e.g., \"192.168.1.100\")",
								"@see REMOTE_PORT, LOCAL_HOST, IS_CONNECTED",
								"@category connection-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok {
								return nil, runtime.Exception{Message: "REMOTE_HOST: invalid connection context"}
							}
							if !wireData.IsConnected {
								return environment.StringValue(""), nil
							}
							remoteAddr := wireData.Conn.RemoteAddr().(*net.TCPAddr)
							return environment.StringValue(remoteAddr.IP.String()), nil
						},
						NativeSet: nil, // Read-only
					},
					"REMOTE_PORT": {
						Variable: environment.Variable{
							Name:     "REMOTE_PORT",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property containing the remote endpoint's port number.",
								"",
								"@property {INTEGR} REMOTE_PORT - Port number of the connected remote host",
								"@example Display connection details",
								"I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT",
								"SAYZ WIT \"Connection from \"",
								"SAYZ WIT CLIENT REMOTE_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT REMOTE_PORT",
								"@example Connection logging",
								"SAYZ WIT \"Accepted connection from \"",
								"SAYZ WIT CLIENT REMOTE_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT REMOTE_PORT",
								"SAYZ WIT \" on local port \"",
								"SAYZ WIT CLIENT LOCAL_PORT",
								"@example Port-based filtering",
								"IZ (CLIENT REMOTE_PORT) BIGGR DAN 1024?",
								"    SAYZ WIT \"Client using non-privileged port\"",
								"NOPE",
								"    SAYZ WIT \"Client using privileged port: \"",
								"    SAYZ WIT CLIENT REMOTE_PORT",
								"KTHX",
								"@note Returns 0 if connection is not established",
								"@note Shows the actual port number the remote client is using",
								"@note Remote port is typically assigned randomly by the client's OS",
								"@see REMOTE_HOST, LOCAL_PORT, IS_CONNECTED",
								"@category connection-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok {
								return nil, runtime.Exception{Message: "REMOTE_PORT: invalid connection context"}
							}
							if !wireData.IsConnected {
								return environment.IntegerValue(0), nil
							}
							remoteAddr := wireData.Conn.RemoteAddr().(*net.TCPAddr)
							return environment.IntegerValue(remoteAddr.Port), nil
						},
						NativeSet: nil, // Read-only
					},
					"LOCAL_HOST": {
						Variable: environment.Variable{
							Name:     "LOCAL_HOST",
							Type:     "STRIN",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property containing the local endpoint's IP address.",
								"",
								"@property {STRIN} LOCAL_HOST - IP address of the local endpoint",
								"@example Show server's local address",
								"I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT",
								"SAYZ WIT \"Server address: \"",
								"SAYZ WIT CLIENT LOCAL_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT LOCAL_PORT",
								"@example Full connection details",
								"SAYZ WIT \"Connection: \"",
								"SAYZ WIT CLIENT LOCAL_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT LOCAL_PORT",
								"SAYZ WIT \" <-> \"",
								"SAYZ WIT CLIENT REMOTE_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT REMOTE_PORT",
								"@example Check local interface",
								"IZ (CLIENT LOCAL_HOST) SAEM AS \"127.0.0.1\"?",
								"    SAYZ WIT \"Local loopback connection\"",
								"NOPE IZ (CLIENT LOCAL_HOST) SAEM AS \"0.0.0.0\"?",
								"    SAYZ WIT \"Listening on all interfaces\"",
								"NOPE",
								"    SAYZ WIT \"Specific interface: \"",
								"    SAYZ WIT CLIENT LOCAL_HOST",
								"KTHX",
								"@note Returns empty string if connection is not established",
								"@note Shows the actual IP address of the local endpoint",
								"@note For servers: shows which interface accepted the connection",
								"@see LOCAL_PORT, REMOTE_HOST, IS_CONNECTED",
								"@category connection-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok {
								return nil, runtime.Exception{Message: "LOCAL_HOST: invalid connection context"}
							}
							if !wireData.IsConnected {
								return environment.StringValue(""), nil
							}
							localAddr := wireData.Conn.LocalAddr().(*net.TCPAddr)
							return environment.StringValue(localAddr.IP.String()), nil
						},
						NativeSet: nil, // Read-only
					},
					"LOCAL_PORT": {
						Variable: environment.Variable{
							Name:     "LOCAL_PORT",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property containing the local endpoint's port number.",
								"",
								"@property {INTEGR} LOCAL_PORT - Port number of the local endpoint",
								"@example Display server port information",
								"I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT",
								"SAYZ WIT \"Server running on port: \"",
								"SAYZ WIT CLIENT LOCAL_PORT",
								"@example Connection summary",
								"SAYZ WIT \"Local: \"",
								"SAYZ WIT CLIENT LOCAL_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT LOCAL_PORT",
								"SAYZ WIT \" | Remote: \"",
								"SAYZ WIT CLIENT REMOTE_HOST",
								"SAYZ WIT \":\"",
								"SAYZ WIT CLIENT REMOTE_PORT",
								"@example Service identification",
								"I HAS A VARIABLE PORT TEH INTEGR ITZ CLIENT LOCAL_PORT",
								"IZ PORT SAEM AS 80?",
								"    SAYZ WIT \"HTTP service\"",
								"NOPE IZ PORT SAEM AS 443?",
								"    SAYZ WIT \"HTTPS service\"",
								"NOPE IZ PORT SAEM AS 22?",
								"    SAYZ WIT \"SSH service\"",
								"NOPE",
								"    SAYZ WIT \"Custom service on port \"",
								"    SAYZ WIT PORT",
								"KTHX",
								"@note Returns 0 if connection is not established",
								"@note Shows the actual port the server is listening on",
								"@note Useful for logging and connection management",
								"@see LOCAL_HOST, REMOTE_PORT, IS_CONNECTED",
								"@category connection-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok {
								return nil, runtime.Exception{Message: "LOCAL_PORT: invalid connection context"}
							}
							if !wireData.IsConnected {
								return environment.IntegerValue(0), nil
							}
							localAddr := wireData.Conn.LocalAddr().(*net.TCPAddr)
							return environment.IntegerValue(localAddr.Port), nil
						},
						NativeSet: nil, // Read-only
					},
					"IS_CONNECTED": {
						Variable: environment.Variable{
							Name:     "IS_CONNECTED",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property indicating whether the connection is currently active.",
								"",
								"@property {BOOL} IS_CONNECTED - YEZ if connection is active, NO if closed",
								"@example Check connection before operations",
								"I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT",
								"IZ CONN IS_CONNECTED?",
								"    CONN DO SEND WIT \"Hello!\"",
								"    I HAS A VARIABLE REPLY TEH STRIN ITZ CONN DO RECEIVE WIT 100",
								"NOPE",
								"    SAYZ WIT \"Connection not available\"",
								"KTHX",
								"@example Connection monitoring loop",
								"WHILE (CONN IS_CONNECTED)",
								"    I HAS A VARIABLE DATA TEH STRIN ITZ CONN DO RECEIVE WIT 256",
								"    IZ DATA SAEM AS \"\"?",
								"        OUTTA HERE BTW Connection closed by remote",
								"    KTHX",
								"    BTW Process received data",
								"    SAYZ WIT \"Received: \"",
								"    SAYZ WIT DATA",
								"KTHX",
								"@example Safe connection cleanup",
								"MAYB",
								"    BTW Some network operations",
								"    CONN DO SEND WIT \"Important message\"",
								"OOPSIE ERR",
								"    SAYZ WIT \"Network error: \"",
								"    SAYZ WIT ERR",
								"FINALLY",
								"    IZ CONN IS_CONNECTED?",
								"        CONN DO CLOSE",
								"    KTHX",
								"KTHX",
								"@note Automatically set to NO when connection fails or is closed",
								"@note Use this to avoid exceptions from operations on closed connections",
								"@note Connection may become disconnected due to network errors or remote closure",
								"@see CLOSE, SEND, RECEIVE",
								"@category connection-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							wireData, ok := this.NativeData.(*WireData)
							if !ok {
								return nil, runtime.Exception{Message: "IS_CONNECTED: invalid connection context"}
							}
							return environment.BoolValue(wireData.IsConnected), nil
						},
						NativeSet: nil, // Read-only
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
		}
	})
	return socketClasses
}

// createWireInstance creates a new WIRE object instance from a net.Conn
func createWireInstance(conn net.Conn) (environment.Value, error) {
	wireClass := getSocketClasses()["WIRE"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(wireClass)

	wireInstance := &environment.ObjectInstance{
		Environment: env,
		Class:       wireClass,
		Variables:   make(map[string]*environment.MemberVariable),
		NativeData: &WireData{
			Conn:        conn,
			IsConnected: true,
		},
	}
	env.InitializeInstanceVariablesWithMRO(wireInstance)

	return wireInstance, nil
}

// RegisterSOCKETInEnv registers SOCKET classes in the given environment
func RegisterSOCKETInEnv(env *environment.Environment, declarations ...string) error {
	socketClasses := getSocketClasses()

	// If declarations is empty, import all classes
	if len(declarations) == 0 {
		for _, class := range socketClasses {
			env.DefineClass(class)
		}
		return nil
	}

	// Otherwise, import only specified classes
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := socketClasses[declUpper]; exists {
			env.DefineClass(class)
			// If importing SOKKIT, also import WIRE (required dependency)
			if declUpper == "SOKKIT" {
				if wireClass, exists := socketClasses["WIRE"]; exists {
					env.DefineClass(wireClass)
				}
			}
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown SOCKET declaration: %s", decl)}
		}
	}

	return nil
}
