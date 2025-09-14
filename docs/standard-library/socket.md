# SOCKET Module

## Import

```lol
BTW Full import
I CAN HAS SOCKET?

BTW Selective import examples
```

## Miscellaneous

### SOKKIT Class

A network socket that provides TCP and UDP networking capabilities.
Supports both client and server operations with configurable protocol, host, port, and timeout settings.

**Methods:**

#### ACCEPT

Accepts an incoming TCP client connection on a listening socket.
Blocks until a client connects, then returns a WIRE connection object.

**Syntax:** `<socket> DO ACCEPT`
**Example: Simple echo server**

```lol
I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
SERVER PORT ITZ 8080
SERVER DO BIND
SERVER DO LISTEN
SAYZ WIT "Echo server listening..."
I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
SAYZ WIT "Client connected!"
I HAS A VARIABLE MSG TEH STRIN ITZ CLIENT DO RECEIVE WIT 1024
CLIENT DO SEND WIT "Echo: " MOAR MSG
CLIENT DO CLOSE
```

**Example: Multi-client server loop**

```lol
SERVER DO LISTEN
IM OUTTA UR LOOP
I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
SAYZ WIT "New client from "
SAYZ WIT CLIENT REMOTE_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT REMOTE_PORT
BTW Handle client in separate thread/process
CLIENT DO SEND WIT "Hello from server!"
CLIENT DO CLOSE
KTHX
```

**Note:** Blocks execution until a client connects

**Note:** Socket must be in LISTEN state before calling ACCEPT

**Note:** Returns new WIRE object for each accepted connection

**Note:** Use returned WIRE object for data transfer with client

#### BIND

Binds the socket to the configured host and port address.
Prepares the socket for server operations (TCP) or datagram communication (UDP).

**Syntax:** `<socket> DO BIND`
**Example: Bind TCP server socket**

```lol
I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
SERVER HOST ITZ "0.0.0.0" BTW Listen on all interfaces
SERVER PORT ITZ 8080
SERVER PROTOCOL ITZ "TCP"
SERVER DO BIND
SAYZ WIT "Server bound to port 8080"
```

**Example: Bind UDP socket**

```lol
I HAS A VARIABLE UDP_SOCK TEH SOKKIT ITZ NEW SOKKIT
UDP_SOCK PROTOCOL ITZ "UDP"
UDP_SOCK PORT ITZ 9999
UDP_SOCK DO BIND
SAYZ WIT "UDP socket bound to port 9999"
```

**Example: Bind with error handling**

```lol
MAYB
SERVER DO BIND
SAYZ WIT "Successfully bound to port"
OOPSIE ERR
SAYZ WIT "Failed to bind: "
SAYZ WIT ERR
KTHX
```

**Note:** Uses current HOST and PORT property values

**Note:** For TCP: enables LISTEN and ACCEPT operations

**Note:** For UDP: enables SEND_TO and RECEIVE_FROM operations

**Note:** Throws exception if port is already in use or insufficient permissions

#### CLOSE

Closes the socket and releases all associated network resources.
Stops listening, closes connections, and frees system resources.

**Syntax:** `<socket> DO CLOSE`
**Example: Proper server shutdown**

```lol
I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
SERVER DO BIND
SERVER DO LISTEN
BTW Server operations here
SERVER DO CLOSE
SAYZ WIT "Server shut down"
```

**Example: Cleanup in exception handling**

```lol
MAYB
SERVER DO BIND
SERVER DO LISTEN
BTW Server work here
OOPSIE ERR
SAYZ WIT "Server error: "
SAYZ WIT ERR
FINALLY
SERVER DO CLOSE BTW Always cleanup
KTHX
```

**Note:** Safe to call multiple times - no error if already closed

**Note:** Automatically called when socket object is garbage collected

**Note:** For TCP servers: stops accepting new connections

**Note:** For UDP sockets: stops receiving packets

**Note:** Does not close existing WIRE connections from ACCEPT

#### CONNECT

Connects to a remote TCP server using configured HOST and PORT.
Returns a WIRE connection object for data transfer, respects TIMEOUT setting.

**Syntax:** `<socket> DO CONNECT`
**Example: Connect to TCP server**

```lol
I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
CLIENT HOST ITZ "127.0.0.1"
CLIENT PORT ITZ 8080
CLIENT TIMEOUT ITZ 10 BTW 10 second timeout
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
SAYZ WIT "Connected to server!"
```

**Example: HTTP client example**

```lol
CLIENT HOST ITZ "httpbin.org"
CLIENT PORT ITZ 80
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
CONN DO SEND WIT "GET / HTTP/1.1\r\nHost: httpbin.org\r\n\r\n"
I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONN DO RECEIVE WIT 4096
SAYZ WIT RESPONSE
CONN DO CLOSE
```

**Example: Connect with error handling**

```lol
MAYB
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
SAYZ WIT "Successfully connected"
BTW Use connection here
CONN DO CLOSE
OOPSIE ERR
SAYZ WIT "Connection failed: "
SAYZ WIT ERR
KTHX
```

**Note:** Only works with TCP protocol sockets

**Note:** Uses current HOST, PORT, and TIMEOUT property values

**Note:** Throws exception if connection fails or times out

**Note:** Returns WIRE object - same type as ACCEPT returns

#### LISTEN

Starts listening for incoming TCP connections on a bound socket.
Enables the socket to accept client connections using ACCEPT method.

**Syntax:** `<socket> DO LISTEN`
**Example: Complete TCP server setup**

```lol
I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
SERVER PROTOCOL ITZ "TCP"
SERVER HOST ITZ "0.0.0.0"
SERVER PORT ITZ 8080
SERVER DO BIND
SERVER DO LISTEN
SAYZ WIT "Server listening on port 8080"
```

**Example: Server with client acceptance loop**

```lol
SERVER DO LISTEN
SAYZ WIT "Waiting for connections..."
IM OUTTA UR LOOP
I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
SAYZ WIT "Client connected from "
SAYZ WIT CLIENT REMOTE_HOST
CLIENT DO SEND WIT "Welcome to server!"
CLIENT DO CLOSE
KTHX
```

**Note:** Only works with TCP protocol sockets

**Note:** Socket must be bound before calling LISTEN

**Note:** After LISTEN, use ACCEPT to handle incoming connections

**Note:** Does not block - ACCEPT is where blocking occurs

#### RECEIVE_FROM

Receives UDP data and returns both the data and sender information.
Blocks until data is received, returns BASKIT with DATA, HOST, and PORT keys.

**Syntax:** `<socket> DO RECEIVE_FROM`
**Example: UDP server receiving data**

```lol
I HAS A VARIABLE UDP_SERVER TEH SOKKIT ITZ NEW SOKKIT
UDP_SERVER PROTOCOL ITZ "UDP"
UDP_SERVER PORT ITZ 9999
UDP_SERVER DO BIND
SAYZ WIT "UDP server listening on port 9999"
I HAS A VARIABLE PACKET TEH BASKIT ITZ UDP_SERVER DO RECEIVE_FROM
SAYZ WIT "Received: "
SAYZ WIT PACKET DO GET WIT "DATA"
SAYZ WIT "From: "
SAYZ WIT PACKET DO GET WIT "HOST"
SAYZ WIT ":"
SAYZ WIT PACKET DO GET WIT "PORT"
```

**Example: Echo UDP server**

```lol
IM OUTTA UR LOOP
I HAS A VARIABLE PACKET TEH BASKIT ITZ UDP_SERVER DO RECEIVE_FROM
I HAS A VARIABLE MSG TEH STRIN ITZ PACKET DO GET WIT "DATA"
I HAS A VARIABLE CLIENT_HOST TEH STRIN ITZ PACKET DO GET WIT "HOST"
I HAS A VARIABLE CLIENT_PORT TEH INTEGR ITZ PACKET DO GET WIT "PORT"
UDP_SERVER DO SEND_TO WIT "Echo: " MOAR MSG AN WIT CLIENT_HOST AN WIT CLIENT_PORT
KTHX
```

**Example: Process UDP packets**

```lol
I HAS A VARIABLE PACKET TEH BASKIT ITZ UDP_SERVER DO RECEIVE_FROM
IZ (PACKET DO GET WIT "DATA") SAEM AS "PING"?
UDP_SERVER DO SEND_TO WIT "PONG" AN WIT (PACKET DO GET WIT "HOST") AN WIT (PACKET DO GET WIT "PORT")
KTHX
```

**Note:** Only works with UDP protocol sockets

**Note:** Socket must be bound before receiving

**Note:** Blocks execution until data arrives

**Note:** Maximum packet size is 4096 bytes

#### SEND_TO

Sends data to a specific UDP address without establishing a connection.
Used for UDP datagram communication - data is sent directly to the target.

**Syntax:** `<socket> DO SEND_TO WIT <data> AN WIT <host> AN WIT <port>`
**Parameters:**
- `data` (STRIN): The data to send
- `host` (STRIN): Target host address (IP or hostname)
- `port` (INTEGR): Target port number

**Example: UDP client sending data**

```lol
I HAS A VARIABLE UDP_CLIENT TEH SOKKIT ITZ NEW SOKKIT
UDP_CLIENT PROTOCOL ITZ "UDP"
UDP_CLIENT PORT ITZ 0 BTW Use any available port
UDP_CLIENT DO BIND
UDP_CLIENT DO SEND_TO WIT "Hello UDP!" AN WIT "127.0.0.1" AN WIT 9999
SAYZ WIT "UDP message sent"
```

**Example: Send to multiple destinations**

```lol
I HAS A VARIABLE SERVERS TEH BUKKIT ITZ NEW BUKKIT
SERVERS DO PUSH WIT "192.168.1.100"
SERVERS DO PUSH WIT "192.168.1.101"
IM OUTTA UR SERVERS NERFIN SERVER_IP
UDP_CLIENT DO SEND_TO WIT "Broadcast message" AN WIT SERVER_IP AN WIT 8888
IM IN UR SERVERS
```

**Example: UDP logging client**

```lol
UDP_CLIENT DO SEND_TO WIT "ERROR: Something went wrong" AN WIT "log-server.local" AN WIT 514
```

**Note:** Only works with UDP protocol sockets

**Note:** Socket must be bound before sending

**Note:** No connection establishment - fire and forget

**Note:** No delivery guarantee - UDP is unreliable

#### SOKKIT

Creates a new socket with default network settings.
Initializes with TCP protocol, localhost host, port 8080, and 30-second timeout.

**Syntax:** `NEW SOKKIT`
**Example: Create default socket**

```lol
I HAS A VARIABLE SOCK TEH SOKKIT ITZ NEW SOKKIT
BTW Socket created with TCP, localhost:8080, 30s timeout
```

**Example: Create and configure socket**

```lol
I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
SERVER PROTOCOL ITZ "TCP"
SERVER HOST ITZ "0.0.0.0" BTW Listen on all interfaces
SERVER PORT ITZ 3000
SERVER TIMEOUT ITZ 60 BTW 60 second timeout
```

**Note:** Use properties (PROTOCOL, HOST, PORT, TIMEOUT) to configure before operations

**Note:** Socket is not bound or connected after creation - use BIND or CONNECT

**Member Variables:**

#### HOST

Target host address for network operations.


**Example: Set specific IP address**

```lol
I HAS A VARIABLE SOCK TEH SOKKIT ITZ NEW SOKKIT
SOCK HOST ITZ "192.168.1.100"
SAYZ WIT "Connecting to "
SAYZ WIT SOCK HOST
```

**Example: Set hostname**

```lol
SOCK HOST ITZ "example.com"
```

**Example: Server listening on all interfaces**

```lol
I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
SERVER HOST ITZ "0.0.0.0" BTW Listen on all network interfaces
SERVER PORT ITZ 8080
SERVER DO BIND
```

**Example: Localhost connections only**

```lol
SERVER HOST ITZ "127.0.0.1" BTW Local connections only
```

**Note:** For servers: determines which network interface to bind to

**Note:** For clients: determines which host to connect to

**Note:** Use "0.0.0.0" to listen on all interfaces (servers only)

**Note:** Use "127.0.0.1" or "localhost" for local-only connections

#### PORT

Target port number for network operations.


**Example: Set web server port**

```lol
I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
SERVER PORT ITZ 80 BTW HTTP default port
```

**Example: Set custom port**

```lol
SERVER PORT ITZ 3000
```

**Example: Use ephemeral port (system assigns)**

```lol
I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
CLIENT PORT ITZ 0 BTW System will assign available port
CLIENT DO BIND
```

**Example: Common port numbers**

```lol
BTW SERVER PORT ITZ 21    BTW FTP
BTW SERVER PORT ITZ 22    BTW SSH
BTW SERVER PORT ITZ 80    BTW HTTP
BTW SERVER PORT ITZ 443   BTW HTTPS
BTW SERVER PORT ITZ 993   BTW IMAPS
```

**Note:** Ports 0-1023 are reserved and may require admin privileges

**Note:** Port 0 means "assign any available port" when binding

**Note:** Port must be in range 0-65535

**Note:** Common ports: 80 (HTTP), 443 (HTTPS), 22 (SSH), 21 (FTP)

#### PROTOCOL

Socket protocol specification - either TCP or UDP.


**Example: Set TCP protocol (default)**

```lol
I HAS A VARIABLE SOCK TEH SOKKIT ITZ NEW SOKKIT
SOCK PROTOCOL ITZ "TCP"
SAYZ WIT "Using TCP protocol"
```

**Example: Set UDP protocol**

```lol
SOCK PROTOCOL ITZ "UDP"
SAYZ WIT "Using UDP protocol"
```

**Example: Check current protocol**

```lol
IZ (SOCK PROTOCOL) SAEM AS "TCP"?
SAYZ WIT "Socket is configured for TCP"
NOPE
SAYZ WIT "Socket is configured for UDP"
KTHX
```

**Note:** TCP provides reliable, ordered, connection-based communication

**Note:** UDP provides fast, connectionless, datagram-based communication

**Note:** Must be set before BIND operation

**Note:** Case-insensitive (converted to uppercase)

#### TIMEOUT

Connection timeout in seconds for TCP client connections.


**Example: Set short timeout**

```lol
I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
CLIENT HOST ITZ "slow-server.com"
CLIENT TIMEOUT ITZ 5 BTW 5 second timeout
MAYB
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
SAYZ WIT "Connected successfully"
OOPSIE ERR
SAYZ WIT "Connection timed out or failed"
KTHX
```

**Example: Set long timeout for slow connections**

```lol
CLIENT TIMEOUT ITZ 120 BTW 2 minute timeout
```

**Example: Disable timeout (wait indefinitely)**

```lol
CLIENT TIMEOUT ITZ 0 BTW No timeout
```

**Example: Check current timeout**

```lol
SAYZ WIT "Connection timeout: "
SAYZ WIT CLIENT TIMEOUT
SAYZ WIT " seconds"
```

**Note:** Only affects TCP CONNECT operations, not UDP or server operations

**Note:** Timeout of 0 means wait indefinitely

**Note:** Must be non-negative integer

**Note:** Default is 30 seconds for new sockets

**Example: TCP server setup**

```lol
I HAS A VARIABLE SERVER TEH SOKKIT ITZ NEW SOKKIT
SERVER HOST ITZ "0.0.0.0"
SERVER PORT ITZ 8080
SERVER PROTOCOL ITZ "TCP"
SERVER DO BIND
SERVER DO LISTEN
```

**Example: TCP client connection**

```lol
I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
CLIENT HOST ITZ "127.0.0.1"
CLIENT PORT ITZ 8080
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
```

**Example: UDP socket communication**

```lol
I HAS A VARIABLE UDP_SOCK TEH SOKKIT ITZ NEW SOKKIT
UDP_SOCK PROTOCOL ITZ "UDP"
UDP_SOCK PORT ITZ 9999
UDP_SOCK DO BIND
UDP_SOCK DO SEND_TO WIT "Hello UDP!" AN WIT "localhost" AN WIT 8888
```

### WIRE Class

A TCP connection that provides bidirectional data transfer capabilities.
Represents an active network connection between two endpoints for reliable data exchange.

**Methods:**

#### CLOSE

Closes the TCP connection and releases associated resources.
Connection becomes unusable after closing and IS_CONNECTED becomes NO.

**Syntax:** `<wire> DO CLOSE`
**Example: Proper connection cleanup**

```lol
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
CONN DO SEND WIT "Hello, Server!"
I HAS A VARIABLE REPLY TEH STRIN ITZ CONN DO RECEIVE WIT 100
CONN DO CLOSE
SAYZ WIT "Connection closed"
```

**Example: Connection cleanup in exception handling**

```lol
MAYB
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
CONN DO SEND WIT "Important data"
I HAS A VARIABLE RESULT TEH STRIN ITZ CONN DO RECEIVE WIT 500
BTW Process result here
OOPSIE ERR
SAYZ WIT "Connection error: "
SAYZ WIT ERR
FINALLY
IZ CONN IS_CONNECTED?
CONN DO CLOSE BTW Always cleanup
KTHX
KTHX
```

**Example: Server client handling**

```lol
I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
CLIENT DO SEND WIT "Welcome!"
I HAS A VARIABLE REQUEST TEH STRIN ITZ CLIENT DO RECEIVE WIT 1024
BTW Process request and send response
CLIENT DO CLOSE BTW Close client connection
```

**Note:** Safe to call multiple times - no error if already closed

**Note:** Connection cannot be reused after closing

**Note:** Always close connections to prevent resource leaks

**Note:** IS_CONNECTED property becomes NO after closing

#### RECEIVE

Receives up to the specified number of characters from the TCP connection.
Blocks until data is available, returns received data (may be shorter than requested).

**Syntax:** `<wire> DO RECEIVE WIT <length>`
**Parameters:**
- `length` (INTEGR): Maximum number of characters to receive

**Example: Receive fixed amount of data**

```lol
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
I HAS A VARIABLE DATA TEH STRIN ITZ CONN DO RECEIVE WIT 1024
SAYZ WIT "Received: "
SAYZ WIT DATA
```

**Example: Receive HTTP response**

```lol
CONN DO SEND WIT "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"
I HAS A VARIABLE HEADERS TEH STRIN ITZ CONN DO RECEIVE WIT 4096
SAYZ WIT "Response headers: "
SAYZ WIT HEADERS
```

**Example: Receive data in chunks**

```lol
I HAS A VARIABLE BUFFER TEH STRIN ITZ ""
WHILE NO SAEM AS (BUFFER ENDS WIT "END")
I HAS A VARIABLE CHUNK TEH STRIN ITZ CONN DO RECEIVE WIT 256
IZ CHUNK SAEM AS ""?
OUTTA HERE BTW Connection closed
KTHX
BUFFER ITZ BUFFER MOAR CHUNK
KTHX
```

**Note:** Blocks execution until data arrives or connection closes

**Note:** Returns empty string if connection is closed by remote end

**Note:** May return less data than requested if that's all that's available

**Note:** Connection must be established before receiving

#### RECEIVE_ALL

Receives all available data from the TCP connection until the connection closes.
Blocks until the remote end closes the connection, then returns all received data.

**Syntax:** `<wire> DO RECEIVE_ALL`
**Example: Download entire web page**

```lol
I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
CLIENT HOST ITZ "httpbin.org"
CLIENT PORT ITZ 80
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
CONN DO SEND WIT "GET /get HTTP/1.1\r\nHost: httpbin.org\r\nConnection: close\r\n\r\n"
I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONN DO RECEIVE_ALL
SAYZ WIT "Full response: "
SAYZ WIT RESPONSE
```

**Example: Receive complete file transfer**

```lol
CONN DO SEND WIT "GET_FILE document.txt"
I HAS A VARIABLE FILE_CONTENT TEH STRIN ITZ CONN DO RECEIVE_ALL
SAYZ WIT "File received, size: "
SAYZ WIT FILE_CONTENT SIZ
```

**Example: Simple protocol with termination**

```lol
CONN DO SEND WIT "DOWNLOAD_DATA"
I HAS A VARIABLE ALL_DATA TEH STRIN ITZ CONN DO RECEIVE_ALL
BTW Server closes connection when done sending
SAYZ WIT "Received all data: "
SAYZ WIT ALL_DATA
```

**Note:** Blocks until remote end closes the connection

**Note:** Connection is automatically marked as closed after completion

**Note:** Useful for protocols where server closes connection to signal end

**Note:** May use significant memory for large data transfers

#### SEND

Sends string data over the TCP connection to the remote endpoint.
Data is transmitted immediately and may be buffered by the network stack.

**Syntax:** `<wire> DO SEND WIT <data>`
**Parameters:**
- `data` (STRIN): The string data to send

**Example: Send simple message**

```lol
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
CONN DO SEND WIT "Hello, Server!"
```

**Example: Send HTTP request**

```lol
CONN DO SEND WIT "GET /api/users HTTP/1.1\r\n"
CONN DO SEND WIT "Host: api.example.com\r\n"
CONN DO SEND WIT "Content-Length: 0\r\n\r\n"
```

**Example: Send JSON data**

```lol
I HAS A VARIABLE JSON_DATA TEH STRIN ITZ "{\"name\":\"Alice\",\"age\":25}"
CONN DO SEND WIT "POST /users HTTP/1.1\r\n"
CONN DO SEND WIT "Content-Type: application/json\r\n"
CONN DO SEND WIT "Content-Length: "
CONN DO SEND WIT JSON_DATA SIZ
CONN DO SEND WIT "\r\n\r\n"
CONN DO SEND WIT JSON_DATA
```

**Note:** Connection must be established (IS_CONNECTED = YEZ)

**Note:** Data is sent as-is - no automatic newlines or formatting added

**Note:** Large data may be sent in multiple network packets

**Note:** Throws exception if connection is broken or closed

**Member Variables:**

#### IS_CONNECTED

Read-only property indicating whether the connection is currently active.


**Example: Check connection before operations**

```lol
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
IZ CONN IS_CONNECTED?
CONN DO SEND WIT "Hello!"
I HAS A VARIABLE REPLY TEH STRIN ITZ CONN DO RECEIVE WIT 100
NOPE
SAYZ WIT "Connection not available"
KTHX
```

**Example: Connection monitoring loop**

```lol
WHILE (CONN IS_CONNECTED)
I HAS A VARIABLE DATA TEH STRIN ITZ CONN DO RECEIVE WIT 256
IZ DATA SAEM AS ""?
OUTTA HERE BTW Connection closed by remote
KTHX
BTW Process received data
SAYZ WIT "Received: "
SAYZ WIT DATA
KTHX
```

**Example: Safe connection cleanup**

```lol
MAYB
BTW Some network operations
CONN DO SEND WIT "Important message"
OOPSIE ERR
SAYZ WIT "Network error: "
SAYZ WIT ERR
FINALLY
IZ CONN IS_CONNECTED?
CONN DO CLOSE
KTHX
KTHX
```

**Note:** Automatically set to NO when connection fails or is closed

**Note:** Use this to avoid exceptions from operations on closed connections

**Note:** Connection may become disconnected due to network errors or remote closure

#### LOCAL_HOST

Read-only property containing the local endpoint's IP address.


**Example: Show server's local address**

```lol
I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
SAYZ WIT "Server address: "
SAYZ WIT CLIENT LOCAL_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT LOCAL_PORT
```

**Example: Full connection details**

```lol
SAYZ WIT "Connection: "
SAYZ WIT CLIENT LOCAL_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT LOCAL_PORT
SAYZ WIT " <-> "
SAYZ WIT CLIENT REMOTE_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT REMOTE_PORT
```

**Example: Check local interface**

```lol
IZ (CLIENT LOCAL_HOST) SAEM AS "127.0.0.1"?
SAYZ WIT "Local loopback connection"
NOPE IZ (CLIENT LOCAL_HOST) SAEM AS "0.0.0.0"?
SAYZ WIT "Listening on all interfaces"
NOPE
SAYZ WIT "Specific interface: "
SAYZ WIT CLIENT LOCAL_HOST
KTHX
```

**Note:** Returns empty string if connection is not established

**Note:** Shows the actual IP address of the local endpoint

**Note:** For servers: shows which interface accepted the connection

#### LOCAL_PORT

Read-only property containing the local endpoint's port number.


**Example: Display server port information**

```lol
I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
SAYZ WIT "Server running on port: "
SAYZ WIT CLIENT LOCAL_PORT
```

**Example: Connection summary**

```lol
SAYZ WIT "Local: "
SAYZ WIT CLIENT LOCAL_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT LOCAL_PORT
SAYZ WIT " | Remote: "
SAYZ WIT CLIENT REMOTE_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT REMOTE_PORT
```

**Example: Service identification**

```lol
I HAS A VARIABLE PORT TEH INTEGR ITZ CLIENT LOCAL_PORT
IZ PORT SAEM AS 80?
SAYZ WIT "HTTP service"
NOPE IZ PORT SAEM AS 443?
SAYZ WIT "HTTPS service"
NOPE IZ PORT SAEM AS 22?
SAYZ WIT "SSH service"
NOPE
SAYZ WIT "Custom service on port "
SAYZ WIT PORT
KTHX
```

**Note:** Returns 0 if connection is not established

**Note:** Shows the actual port the server is listening on

**Note:** Useful for logging and connection management

#### REMOTE_HOST

Read-only property containing the remote endpoint's IP address.


**Example: Check client connection details**

```lol
I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
SAYZ WIT "Client connected from: "
SAYZ WIT CLIENT REMOTE_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT REMOTE_PORT
```

**Example: Log connection information**

```lol
SAYZ WIT "New connection: "
SAYZ WIT CLIENT REMOTE_HOST
SAYZ WIT " -> "
SAYZ WIT CLIENT LOCAL_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT LOCAL_PORT
```

**Example: Access control by IP**

```lol
IZ (CLIENT REMOTE_HOST) SAEM AS "127.0.0.1"?
SAYZ WIT "Local connection allowed"
NOPE
SAYZ WIT "Remote connection from "
SAYZ WIT CLIENT REMOTE_HOST
CLIENT DO CLOSE BTW Block external connections
KTHX
```

**Note:** Returns empty string if connection is not established

**Note:** Shows actual IP address, not hostname

**Note:** IPv4 addresses shown in dotted decimal notation (e.g., "192.168.1.100")

#### REMOTE_PORT

Read-only property containing the remote endpoint's port number.


**Example: Display connection details**

```lol
I HAS A VARIABLE CLIENT TEH WIRE ITZ SERVER DO ACCEPT
SAYZ WIT "Connection from "
SAYZ WIT CLIENT REMOTE_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT REMOTE_PORT
```

**Example: Connection logging**

```lol
SAYZ WIT "Accepted connection from "
SAYZ WIT CLIENT REMOTE_HOST
SAYZ WIT ":"
SAYZ WIT CLIENT REMOTE_PORT
SAYZ WIT " on local port "
SAYZ WIT CLIENT LOCAL_PORT
```

**Example: Port-based filtering**

```lol
IZ (CLIENT REMOTE_PORT) BIGGR DAN 1024?
SAYZ WIT "Client using non-privileged port"
NOPE
SAYZ WIT "Client using privileged port: "
SAYZ WIT CLIENT REMOTE_PORT
KTHX
```

**Note:** Returns 0 if connection is not established

**Note:** Shows the actual port number the remote client is using

**Note:** Remote port is typically assigned randomly by the client's OS

**Example: Client connection usage**

```lol
I HAS A VARIABLE CLIENT TEH SOKKIT ITZ NEW SOKKIT
CLIENT HOST ITZ "127.0.0.1"
CLIENT PORT ITZ 8080
I HAS A VARIABLE CONN TEH WIRE ITZ CLIENT DO CONNECT
CONN DO SEND WIT "GET /api/data HTTP/1.1\r\n\r\n"
I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONN DO RECEIVE WIT 1024
SAYZ WIT RESPONSE
CONN DO CLOSE
```

**Example: Server-side connection handling**

```lol
BTW From server ACCEPT
I HAS A VARIABLE CLIENT_CONN TEH WIRE ITZ SERVER DO ACCEPT
SAYZ WIT "Client connected from "
SAYZ WIT CLIENT_CONN REMOTE_HOST
CLIENT_CONN DO SEND WIT "Welcome to server!"
I HAS A VARIABLE REQUEST TEH STRIN ITZ CLIENT_CONN DO RECEIVE WIT 512
CLIENT_CONN DO CLOSE
```

**Example: Bidirectional communication**

```lol
CONN DO SEND WIT "HELLO"
I HAS A VARIABLE REPLY TEH STRIN ITZ CONN DO RECEIVE WIT 100
IZ REPLY SAEM AS "OK"?
CONN DO SEND WIT "DATA: important message"
KTHX
```

