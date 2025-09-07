# SOCKET Module - Network Socket Operations

The SOCKET module provides TCP and UDP socket functionality through the SOKKIT class for socket management and the WIRE class for TCP connections.

## Importing SOCKET Module

```lol
BTW Import entire module
I CAN HAS SOCKET?

BTW Selective import
I CAN HAS SOKKIT FROM SOCKET?
I CAN HAS WIRE FROM SOCKET?
```

**Note:** The SOCKET module automatically imports the WIRE class when SOKKIT is imported.

## SOKKIT Class

The SOKKIT class represents a network socket that can operate in TCP or UDP mode. It provides server and client functionality with configurable network parameters.

### Constructor

```lol
I HAS A VARIABLE SOCKET TEH SOKKIT ITZ NEW SOKKIT
```

The constructor creates a socket with default settings:
- **Protocol**: TCP
- **Host**: localhost
- **Port**: 8080
- **Timeout**: 30 seconds

### Properties

- **PROTOCOL**: STRIN - Socket protocol, either "TCP" or "UDP" (default: "TCP")
- **HOST**: STRIN - Target host address (default: "localhost")
- **PORT**: INTEGR - Target port number (default: 8080)
- **TIMEOUT**: INTEGR - Connection timeout in seconds (default: 30)

### Methods

#### Socket Management

##### BIND - Bind Socket to Address

Binds the socket to the specified host and port for server operations.

```lol
socket DO BIND
```

**Parameters:** None (uses HOST and PORT properties)

**Throws:** Exception if binding fails

##### LISTEN - Start Listening (TCP only)

Starts listening for incoming connections on a bound TCP socket.

```lol
socket DO LISTEN
```

**Parameters:** None

**Throws:** Exception if socket is not bound or not TCP

##### CLOSE - Close Socket

Closes the socket and releases resources.

```lol
socket DO CLOSE
```

**Parameters:** None

#### TCP Operations

##### CONNECT - Connect to Server

Connects to a remote TCP server and returns a WIRE connection object.

```lol
I HAS A VARIABLE CONNECTION TEH WIRE ITZ socket DO CONNECT
```

**Parameters:** None (uses HOST, PORT, and TIMEOUT properties)

**Returns:** WIRE - The connection object

**Throws:** Exception if connection fails

##### ACCEPT - Accept Incoming Connection

Accepts an incoming connection on a listening TCP socket.

```lol
I HAS A VARIABLE CLIENT TEH WIRE ITZ socket DO ACCEPT
```

**Parameters:** None

**Returns:** WIRE - The client connection object

**Throws:** Exception if socket is not listening

#### UDP Operations

##### SEND_TO - Send UDP Data

Sends data to a specific UDP address.

```lol
socket DO SEND_TO WIT <data> AN WIT <host> AN WIT <port>
```

**Parameters:**
- **data**: STRIN - The data to send
- **host**: STRIN - Target host address
- **port**: INTEGR - Target port number

**Throws:** Exception if socket is not bound or not UDP

##### RECEIVE_FROM - Receive UDP Data

Receives data from a UDP socket and returns sender information.

```lol
I HAS A VARIABLE PACKET TEH BASKIT ITZ socket DO RECEIVE_FROM
```

**Parameters:** None

**Returns:** BASKIT - Contains "DATA", "HOST", and "PORT" keys

**Throws:** Exception if socket is not bound or not UDP

## WIRE Class

The WIRE class represents a TCP connection providing bidirectional data transfer capabilities.

### Properties

- **REMOTE_HOST**: STRIN (read-only) - Remote endpoint IP address
- **REMOTE_PORT**: INTEGR (read-only) - Remote endpoint port number
- **LOCAL_HOST**: STRIN (read-only) - Local endpoint IP address
- **LOCAL_PORT**: INTEGR (read-only) - Local endpoint port number
- **IS_CONNECTED**: BOOL (read-only) - Connection status

### Methods

##### SEND - Send Data

Sends data over the TCP connection.

```lol
connection DO SEND WIT <data>
```

**Parameters:**
- **data**: STRIN - The data to send

**Throws:** Exception if connection is not established

##### RECEIVE - Receive Data

Receives a specific amount of data from the TCP connection.

```lol
I HAS A VARIABLE DATA TEH STRIN ITZ connection DO RECEIVE WIT <length>
```

**Parameters:**
- **length**: INTEGR - Maximum number of bytes to receive

**Returns:** STRIN - The received data (may be shorter than requested)

##### RECEIVE_ALL - Receive All Data

Receives all available data from the TCP connection until it closes.

```lol
I HAS A VARIABLE DATA TEH STRIN ITZ connection DO RECEIVE_ALL
```

**Parameters:** None

**Returns:** STRIN - All received data

##### CLOSE - Close Connection

Closes the TCP connection.

```lol
connection DO CLOSE
```

**Parameters:** None

## TCP Socket Operations

### TCP Server

```lol
I CAN HAS SOCKET?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN TCP_SERVER WIT PORT TEH INTEGR
    SAYZ WIT "=== TCP Server Example ==="

    BTW Create server socket
    I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
    SERVER PROTOCOL ITZ "TCP"
    SERVER HOST ITZ "localhost"
    SERVER PORT ITZ PORT

    MAYB
        BTW Bind and listen
        SERVER DO BIND
        SAYZ WIT "Server bound to port " + PORT
        SERVER DO LISTEN
        SAYZ WIT "Server listening for connections..."

        BTW Accept client connection
        I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
        SAYZ WIT "Client connected!"

        BTW Show client information
        SAY WIT "Client address: "
        SAY WIT CLIENT REMOTE_HOST
        SAY WIT ":"
        SAYZ WIT CLIENT REMOTE_PORT

        BTW Receive message from client
        I HAS A VARIABLE MESSAGE TEH STRIN ITZ CLIENT DO RECEIVE WIT 1024
        SAY WIT "Received: "
        SAYZ WIT MESSAGE

        BTW Send response
        CLIENT DO SEND WIT "Hello from server!"

        BTW Clean up
        CLIENT DO CLOSE
        SERVER DO CLOSE
        SAYZ WIT "Server closed"

    OOPSIE ERR
        SAYZ WIT "Server error: " + ERR
        SERVER DO CLOSE
    KTHX
KTHXBAI
```

### TCP Client

```lol
I CAN HAS SOCKET?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN TCP_CLIENT WIT HOST TEH STRIN AN WIT PORT TEH INTEGR
    SAYZ WIT "=== TCP Client Example ==="

    BTW Create client socket
    I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
    CLIENT PROTOCOL ITZ "TCP"
    CLIENT HOST ITZ HOST
    CLIENT PORT ITZ PORT
    CLIENT TIMEOUT ITZ 10

    MAYB
        BTW Connect to server
        I HAS A VARIABLE CONNECTION TEH WIRE ITZ CLIENT DO CONNECT
        SAYZ WIT "Connected to server!"

        BTW Show connection information
        SAY WIT "Connected to: "
        SAY WIT CONNECTION REMOTE_HOST
        SAY WIT ":"
        SAYZ WIT CONNECTION REMOTE_PORT

        BTW Send message
        CONNECTION DO SEND WIT "Hello from client!"
        SAYZ WIT "Message sent"

        BTW Receive response
        I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONNECTION DO RECEIVE WIT 1024
        SAY WIT "Server responded: "
        SAYZ WIT RESPONSE

        BTW Close connection
        CONNECTION DO CLOSE
        SAYZ WIT "Connection closed"

    OOPSIE ERR
        SAYZ WIT "Client error: " + ERR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_TCP_CLIENT
    TCP_CLIENT WIT "localhost" AN WIT 8080
KTHXBAI
```

## UDP Socket Operations

### UDP Server

```lol
I CAN HAS SOCKET?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN UDP_SERVER WIT PORT TEH INTEGR
    SAYZ WIT "=== UDP Server Example ==="

    BTW Create UDP server socket
    I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
    SERVER PROTOCOL ITZ "UDP"
    SERVER HOST ITZ "localhost"
    SERVER PORT ITZ PORT

    MAYB
        BTW Bind UDP socket
        SERVER DO BIND
        SAYZ WIT "UDP server listening on port " + PORT

        BTW Receive data from client
        I HAS A VARIABLE PACKET TEH BASKIT ITZ SERVER DO RECEIVE_FROM
        
        I HAS A VARIABLE DATA TEH STRIN ITZ PACKET DO GET WIT "DATA"
        I HAS A VARIABLE CLIENT_HOST TEH STRIN ITZ PACKET DO GET WIT "HOST"
        I HAS A VARIABLE CLIENT_PORT TEH INTEGR ITZ PACKET DO GET WIT "PORT"

        SAY WIT "Received from "
        SAY WIT CLIENT_HOST
        SAY WIT ":"
        SAY WIT CLIENT_PORT
        SAY WIT " - "
        SAYZ WIT DATA

        BTW Send response back to client
        SERVER DO SEND_TO WIT "Server response: " + DATA AN WIT CLIENT_HOST AN WIT CLIENT_PORT
        SAYZ WIT "Response sent back to client"

        BTW Clean up
        SERVER DO CLOSE
        SAYZ WIT "UDP server closed"

    OOPSIE ERR
        SAYZ WIT "UDP server error: " + ERR
        SERVER DO CLOSE
    KTHX
KTHXBAI
```

### UDP Client

```lol
I CAN HAS SOCKET?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN UDP_CLIENT WIT HOST TEH STRIN AN WIT PORT TEH INTEGR
    SAYZ WIT "=== UDP Client Example ==="

    BTW Create UDP client socket
    I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
    CLIENT PROTOCOL ITZ "UDP"
    CLIENT HOST ITZ "localhost"
    CLIENT PORT ITZ 0  BTW Use any available port

    MAYB
        BTW Bind to any available port
        CLIENT DO BIND
        SAYZ WIT "UDP client ready"

        BTW Send data to server
        I HAS A VARIABLE MESSAGE TEH STRIN ITZ "Hello UDP server!"
        CLIENT DO SEND_TO WIT MESSAGE AN WIT HOST AN WIT PORT
        SAYZ WIT "Message sent to server"

        BTW Receive response
        I HAS A VARIABLE RESPONSE TEH BASKIT ITZ CLIENT DO RECEIVE_FROM
        I HAS A VARIABLE REPLY TEH STRIN ITZ RESPONSE DO GET WIT "DATA"
        
        SAY WIT "Server replied: "
        SAYZ WIT REPLY

        BTW Clean up
        CLIENT DO CLOSE
        SAYZ WIT "UDP client closed"

    OOPSIE ERR
        SAYZ WIT "UDP client error: " + ERR
        CLIENT DO CLOSE
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_UDP_CLIENT
    UDP_CLIENT WIT "localhost" AN WIT 8081
KTHXBAI
```

## Advanced Socket Operations

### Multi-Client TCP Server

```lol
I CAN HAS SOCKET?
I CAN HAS STDIO?
I CAN HAS THREAD?

HAI ME TEH FUNCSHUN HANDLE_CLIENT WIT CLIENT_CONN TEH WIRE
    MAYB
        BTW Get client information
        I HAS A VARIABLE CLIENT_ADDR TEH STRIN ITZ CLIENT_CONN REMOTE_HOST + ":" + CLIENT_CONN REMOTE_PORT
        SAYZ WIT "Handling client: " + CLIENT_ADDR

        BTW Echo server - receive and send back
        I HAS A VARIABLE DATA TEH STRIN ITZ CLIENT_CONN DO RECEIVE WIT 1024
        IZ DATA SAEM AS ""?
            SAYZ WIT "Client " + CLIENT_ADDR + " disconnected"
        NOPE
            SAYZ WIT "Echoing to " + CLIENT_ADDR + ": " + DATA
            CLIENT_CONN DO SEND WIT "Echo: " + DATA
        KTHX

        CLIENT_CONN DO CLOSE
        SAYZ WIT "Client " + CLIENT_ADDR + " handler finished"

    OOPSIE ERR
        SAYZ WIT "Client handling error: " + ERR
        CLIENT_CONN DO CLOSE
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN MULTI_CLIENT_SERVER WIT PORT TEH INTEGR
    SAYZ WIT "=== Multi-Client TCP Server ==="

    I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
    SERVER PROTOCOL ITZ "TCP"
    SERVER HOST ITZ "localhost"
    SERVER PORT ITZ PORT

    MAYB
        SERVER DO BIND
        SERVER DO LISTEN
        SAYZ WIT "Multi-client server listening on port " + PORT

        BTW Accept multiple clients (simplified example)
        BTW In practice, you'd use threading for each client
        I HAS A VARIABLE CLIENT_COUNT TEH INTEGR ITZ 0
        
        BTW Accept and handle first client
        I HAS A VARIABLE CLIENT1 TEH WIRE ITZ SERVER DO ACCEPT
        SAYZ WIT "Client 1 connected"
        HANDLE_CLIENT WIT CLIENT1

        SERVER DO CLOSE
        SAYZ WIT "Multi-client server closed"

    OOPSIE ERR
        SAYZ WIT "Multi-client server error: " + ERR
        SERVER DO CLOSE
    KTHX
KTHXBAI
```

### Socket with Custom Configuration

```lol
I CAN HAS SOCKET?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN CONFIGURED_SOCKET_CLIENT WIT HOST TEH STRIN AN WIT PORT TEH INTEGR
    SAYZ WIT "=== Configured Socket Client ==="

    BTW Create socket with custom settings
    I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
    CLIENT PROTOCOL ITZ "TCP"
    CLIENT HOST ITZ HOST
    CLIENT PORT ITZ PORT
    CLIENT TIMEOUT ITZ 5  BTW Short timeout for demo

    MAYB
        SAYZ WIT "Attempting connection with 5 second timeout..."
        I HAS A VARIABLE CONNECTION TEH WIRE ITZ CLIENT DO CONNECT

        SAYZ WIT "Connection established!"
        SAY WIT "Local address: "
        SAY WIT CONNECTION LOCAL_HOST
        SAY WIT ":"
        SAYZ WIT CONNECTION LOCAL_PORT

        SAY WIT "Remote address: "
        SAY WIT CONNECTION REMOTE_HOST
        SAY WIT ":"
        SAYZ WIT CONNECTION REMOTE_PORT

        SAY WIT "Connection status: "
        SAYZ WIT CONNECTION IS_CONNECTED

        CONNECTION DO CLOSE
        SAYZ WIT "Connection closed"

    OOPSIE ERR
        SAYZ WIT "Connection failed: " + ERR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_CONFIGURED_CLIENT
    CONFIGURED_SOCKET_CLIENT WIT "httpbin.org" AN WIT 80
KTHXBAI
```

## Error Handling

### Network Error Handling

```lol
I CAN HAS SOCKET?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN ROBUST_SOCKET_CLIENT WIT HOST TEH STRIN AN WIT PORT TEH INTEGR
    SAYZ WIT "=== Robust Socket Client ==="

    I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
    CLIENT PROTOCOL ITZ "TCP"
    CLIENT HOST ITZ HOST
    CLIENT PORT ITZ PORT
    CLIENT TIMEOUT ITZ 10

    MAYB
        I HAS A VARIABLE CONNECTION TEH WIRE ITZ CLIENT DO CONNECT
        SAYZ WIT "Connection successful!"

        BTW Test connection status before operations
        IZ CONNECTION IS_CONNECTED?
            SAYZ WIT "Connection is active, sending data..."
            CONNECTION DO SEND WIT "Test message"
            
            I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONNECTION DO RECEIVE WIT 1024
            IZ RESPONSE SAEM AS ""?
                SAYZ WIT "No data received (connection may be closed)"
            NOPE
                SAYZ WIT "Received: " + RESPONSE
            KTHX
        NOPE
            SAYZ WIT "Connection is not active"
        KTHX

        CONNECTION DO CLOSE

    OOPSIE ERR
        SAYZ WIT "Socket error: " + ERR
        SAYZ WIT "This could be due to:"
        SAYZ WIT "- Network connectivity issues"
        SAYZ WIT "- Invalid host or port"
        SAYZ WIT "- Connection timeout"
        SAYZ WIT "- Server not responding"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_ERROR_HANDLING
    BTW Test with unreachable host
    ROBUST_SOCKET_CLIENT WIT "192.0.2.1" AN WIT 12345

    BTW Test with invalid port
    ROBUST_SOCKET_CLIENT WIT "localhost" AN WIT 99999
KTHXBAI
```

### UDP Error Scenarios

```lol
I CAN HAS SOCKET?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN UDP_ERROR_SCENARIOS
    SAYZ WIT "=== UDP Error Scenarios ==="

    I HAS A VARIABLE UDP_SOCKET TEH SOKKIT ITZ NEW SOKKIT
    UDP_SOCKET PROTOCOL ITZ "UDP"

    BTW Test operations without binding
    SAYZ WIT "Testing operations without binding..."
    MAYB
        UDP_SOCKET DO SEND_TO WIT "test" AN WIT "localhost" AN WIT 9999
        SAYZ WIT "ERROR: SEND_TO should have failed without bind"
    OOPSIE ERR
        SAYZ WIT "✓ SEND_TO correctly failed without bind: " + ERR
    KTHX

    MAYB
        I HAS A VARIABLE DATA TEH BASKIT ITZ UDP_SOCKET DO RECEIVE_FROM
        SAYZ WIT "ERROR: RECEIVE_FROM should have failed without bind"
    OOPSIE ERR
        SAYZ WIT "✓ RECEIVE_FROM correctly failed without bind: " + ERR
    KTHX

    BTW Test with proper binding
    SAYZ WIT "Testing with proper binding..."
    MAYB
        UDP_SOCKET HOST ITZ "localhost"
        UDP_SOCKET PORT ITZ 0  BTW Any available port
        UDP_SOCKET DO BIND
        SAYZ WIT "✓ UDP socket bound successfully"

        BTW Test sending to invalid address
        MAYB
            UDP_SOCKET DO SEND_TO WIT "test" AN WIT "invalid.host.name" AN WIT 9999
        OOPSIE ERR
            SAYZ WIT "✓ Invalid address correctly handled: " + ERR
        KTHX

        UDP_SOCKET DO CLOSE
        SAYZ WIT "✓ UDP socket closed"

    OOPSIE ERR
        SAYZ WIT "Bind error: " + ERR
    KTHX
KTHXBAI
```

## Quick Reference

### Constructor

| Usage | Description |
|-------|-------------|
| `NEW SOKKIT` | Create socket with default TCP settings |

### SOKKIT Configuration

| Property | Type | Description |
|----------|------|-------------|
| `PROTOCOL` | STRIN | "TCP" or "UDP" |
| `HOST` | STRIN | Target host address |
| `PORT` | INTEGR | Target port number |
| `TIMEOUT` | INTEGR | Connection timeout in seconds |

### SOKKIT Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `BIND` | None | None | Bind socket to address |
| `LISTEN` | None | None | Start listening (TCP only) |
| `CONNECT` | None | WIRE | Connect to server (TCP only) |
| `ACCEPT` | None | WIRE | Accept connection (TCP only) |
| `SEND_TO WIT data AN WIT host AN WIT port` | data: STRIN, host: STRIN, port: INTEGR | None | Send UDP data |
| `RECEIVE_FROM` | None | BASKIT | Receive UDP data |
| `CLOSE` | None | None | Close socket |

### WIRE Properties

| Property | Type | Description |
|----------|------|-------------|
| `REMOTE_HOST` | STRIN | Remote IP address |
| `REMOTE_PORT` | INTEGR | Remote port number |
| `LOCAL_HOST` | STRIN | Local IP address |
| `LOCAL_PORT` | INTEGR | Local port number |
| `IS_CONNECTED` | BOOL | Connection status |

### WIRE Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `SEND WIT data` | data: STRIN | None | Send TCP data |
| `RECEIVE WIT length` | length: INTEGR | STRIN | Receive TCP data |
| `RECEIVE_ALL` | None | STRIN | Receive all available data |
| `CLOSE` | None | None | Close connection |

## Related

- [HTTP Module](http.md) - High-level HTTP client operations
- [THREAD Module](threading.md) - Concurrency for multi-client servers
- [STDIO Module](stdio.md) - Console output for debugging
- [Collections](collections.md) - BASKIT operations for UDP packets
- [Control Flow](../language-guide/control-flow.md) - Exception handling patterns