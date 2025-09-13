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
					"A network socket that can operate in TCP or UDP mode.",
					"Provides server and client functionality with configurable network parameters.",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"SOKKIT": {
						Name:       "SOKKIT",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Creates a socket with default settings: TCP protocol, localhost host, port 8080.",
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
							"Binds the socket to the specified host and port for server operations.",
							"Uses HOST and PORT properties, throws exception if binding fails.",
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
							"Starts listening for incoming connections on a bound TCP socket.",
							"Throws exception if socket is not bound or not TCP.",
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
							"Accepts an incoming connection on a listening TCP socket.",
							"Returns WIRE connection object, throws exception if socket is not listening.",
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
							"Connects to a remote TCP server and returns a WIRE connection object.",
							"Uses HOST, PORT, and TIMEOUT properties, throws exception if connection fails.",
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
							"Sends data to a specific UDP address.",
							"Throws exception if socket is not bound or not UDP.",
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
							"Receives data from a UDP socket and returns sender information.",
							"Returns BASKIT with DATA, HOST, and PORT keys, throws exception if socket is not bound or not UDP.",
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
						Name:          "CLOSE",
						Parameters:    []environment.Parameter{},
						Documentation: []string{"Closes the socket and releases resources."},
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
							Name:          "PROTOCOL",
							Type:          "STRIN",
							IsLocked:      false,
							IsPublic:      true,
							Documentation: []string{"Socket protocol, either TCP or UDP (default: TCP)."},
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
							Name:          "HOST",
							Type:          "STRIN",
							IsLocked:      false,
							IsPublic:      true,
							Documentation: []string{"Target host address (default: localhost)."},
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
							Name:          "PORT",
							Type:          "INTEGR",
							IsLocked:      false,
							IsPublic:      true,
							Documentation: []string{"Target port number (default: 8080, valid range: 0-65535)."},
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
							Name:          "TIMEOUT",
							Type:          "INTEGR",
							IsLocked:      false,
							IsPublic:      true,
							Documentation: []string{"Connection timeout in seconds (default: 30)."},
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
				Documentation: []string{"A TCP connection providing bidirectional data transfer capabilities."},
				PublicFunctions: map[string]*environment.Function{
					// SEND method
					"SEND": {
						Name: "SEND",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
						},
						Documentation: []string{
							"Sends data over the TCP connection.",
							"Throws exception if connection is not established.",
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
							"Receives a specific amount of data from the TCP connection.",
							"Returns received data (may be shorter than requested).",
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
							"Receives all available data from the TCP connection until it closes.",
							"Returns all received data as a single string.",
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
						Name:          "CLOSE",
						Parameters:    []environment.Parameter{},
						Documentation: []string{"Closes the TCP connection."},
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
							Name:          "REMOTE_HOST",
							Type:          "STRIN",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"Remote endpoint IP address."},
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
							Name:          "REMOTE_PORT",
							Type:          "INTEGR",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"Remote endpoint port number."},
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
							Name:          "LOCAL_HOST",
							Type:          "STRIN",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"Local endpoint IP address"},
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
							Name:          "LOCAL_PORT",
							Type:          "INTEGR",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"Local endpoint port number."},
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
							Name:          "IS_CONNECTED",
							Type:          "BOOL",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"Connection status."},
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
