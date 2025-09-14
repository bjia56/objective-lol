# HTTP Module

## Import

```lol
BTW Full import
I CAN HAS HTTP?

BTW Selective import examples
```

## Miscellaneous

### INTERWEB Class

HTTP client that can make requests to web servers.
Supports HTTP and HTTPS protocols with configurable timeout and headers.

**Methods:**

#### DELETE

Makes an HTTP DELETE request to the specified URL.
Returns a RESPONSE object containing status, body, and headers.

**Syntax:** `<client> DO DELETE WIT <url>`
**Parameters:**
- `url` (STRIN): The URL to request (must include http:// or https://)

**Example: Delete a resource**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT "https://api.example.com/users/123"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "User deleted successfully"
NOPE
SAYZ WIT "Delete failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Example: Delete with authentication**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Authorization" AN WIT "Bearer token123"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT "https://api.example.com/posts/456"
IZ RESPONSE STATUS SAEM AS 204?
SAYZ WIT "Post deleted (no content returned)"
NOPE
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Post deleted successfully"
NOPE
SAYZ WIT "Delete failed"
KTHX
KTHX
```

**Example: Check if resource exists before deleting**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE CHECK_RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/items/789"
IZ CHECK_RESPONSE STATUS SAEM AS 404?
SAYZ WIT "Resource does not exist"
NOPE
I HAS A VARIABLE DELETE_RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT "https://api.example.com/items/789"
IZ DELETE_RESPONSE IS_SUCCESS?
SAYZ WIT "Resource deleted successfully"
NOPE
SAYZ WIT "Delete failed"
KTHX
KTHX
```

**Example: Batch delete with error handling**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE IDS TEH BUKKIT ITZ NEW BUKKIT
IDS DO PUSH WIT "101"
IDS DO PUSH WIT "102"
IDS DO PUSH WIT "103"
I HAS A VARIABLE SUCCESS_COUNT TEH INTEGR ITZ 0
I HAS A VARIABLE ERROR_COUNT TEH INTEGR ITZ 0
WHILE NO SAEM AS (IDS LENGTH SAEM AS 0)
I HAS A VARIABLE ID TEH STRIN ITZ IDS DO POP
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT "https://api.example.com/resources/" MOAR ID
IZ RESPONSE IS_SUCCESS?
SUCCESS_COUNT ITZ SUCCESS_COUNT UP 1
NOPE
ERROR_COUNT ITZ ERROR_COUNT UP 1
SAYZ WIT "Failed to delete resource "
SAYZ WIT ID
SAYZ WIT " (status: "
SAYZ WIT RESPONSE STATUS
SAYZ WIT ")"
KTHX
KTHX
SAYZ WIT "Batch delete completed: "
SAYZ WIT SUCCESS_COUNT
SAYZ WIT " successful, "
SAYZ WIT ERROR_COUNT
SAYZ WIT " failed"
```

**Note:** URL must include protocol (http:// or https://)

**Note:** DELETE operations should be idempotent

**Note:** Response body may be empty (204 No Content) on success

**Note:** Applies any headers set on the client

#### GET

Makes an HTTP GET request to the specified URL.
Returns a RESPONSE object containing status, body, and headers.

**Syntax:** `<client> DO GET WIT <url>`
**Parameters:**
- `url` (STRIN): The URL to request (must include http:// or https://)

**Example: Simple GET request**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/get"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Success! Status: "
SAYZ WIT RESPONSE STATUS
SAYZ WIT "Body length: "
SAYZ WIT RESPONSE BODY LENGTH
NOPE
SAYZ WIT "Request failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Example: GET with custom headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Accept" AN WIT "application/json"
CLIENT HEADERS DO PUT WIT "User-Agent" AN WIT "MyApp/1.0"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.github.com/user"
SAYZ WIT "Response status: "
SAYZ WIT RESPONSE STATUS
```

**Example: Handle different response types**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/json"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON
SAYZ WIT "Parsed JSON response"
NOPE
SAYZ WIT "Request failed or returned non-JSON content"
SAYZ WIT RESPONSE BODY
KTHX
```

**Example: Check response headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/headers"
I HAS A VARIABLE CONTENT_TYPE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "content-type"
SAYZ WIT "Content-Type: "
SAYZ WIT CONTENT_TYPE
```

**Note:** URL must include protocol (http:// or https://)

**Note:** Applies any headers set on the client

**Note:** Respects the client's timeout setting

#### INTERWEB

Initializes an INTERWEB HTTP client with default settings.
Default timeout is 30 seconds with empty headers.

**Syntax:** `NEW INTERWEB`
**Example: Create basic HTTP client**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
SAYZ WIT "HTTP client created with default settings"
```

**Example: Create client and immediately configure**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT TIMEOUT ITZ 60
CLIENT HEADERS DO PUT WIT "User-Agent" AN WIT "MyApp/1.0"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/user-agent"
SAYZ WIT RESPONSE BODY
```

**Example: Create multiple clients with different configurations**

```lol
I HAS A VARIABLE API_CLIENT TEH INTERWEB ITZ NEW INTERWEB
API_CLIENT HEADERS DO PUT WIT "Authorization" AN WIT "Bearer token123"
API_CLIENT TIMEOUT ITZ 10
I HAS A VARIABLE WEB_CLIENT TEH INTERWEB ITZ NEW INTERWEB
WEB_CLIENT HEADERS DO PUT WIT "User-Agent" AN WIT "WebScraper/1.0"
WEB_CLIENT TIMEOUT ITZ 30
SAYZ WIT "Created two clients with different configurations"
```

**Note:** Initializes with 30-second timeout

**Note:** Creates empty headers BASKIT that can be populated

**Note:** Client is ready to make requests immediately

#### POST

Makes an HTTP POST request with data in the request body to the specified URL.
Returns a RESPONSE object containing status, body, and headers.

**Syntax:** `<client> DO POST WIT <url> AN WIT <data>`
**Parameters:**
- `url` (STRIN): The URL to request (must include http:// or https://)
- `data` (STRIN): The data to send in the request body

**Example: POST JSON data**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
I HAS A VARIABLE JSON_DATA TEH STRIN ITZ "{\"name\":\"Alice\",\"email\":\"alice@example.com\"}"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://httpbin.org/post" AN WIT JSON_DATA
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Data posted successfully"
I HAS A VARIABLE RESPONSE_JSON TEH BASKIT ITZ RESPONSE DO TO_JSON
SAYZ WIT "Server received our data"
NOPE
SAYZ WIT "POST failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Example: POST form data**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/x-www-form-urlencoded"
I HAS A VARIABLE FORM_DATA TEH STRIN ITZ "username=alice&password=secret123"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://httpbin.org/post" AN WIT FORM_DATA
SAYZ WIT "Form posted, status: "
SAYZ WIT RESPONSE STATUS
```

**Example: POST with authentication**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Authorization" AN WIT "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
I HAS A VARIABLE PAYLOAD TEH STRIN ITZ "{\"action\":\"create\",\"resource\":\"user\"}"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://api.example.com/users" AN WIT PAYLOAD
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "User created successfully"
NOPE
SAYZ WIT "Failed to create user: "
SAYZ WIT RESPONSE BODY
KTHX
```

**Example: POST large data**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT TIMEOUT ITZ 60
I HAS A VARIABLE LARGE_DATA TEH STRIN ITZ "<very large XML or JSON payload>"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://api.example.com/upload" AN WIT LARGE_DATA
SAYZ WIT "Upload completed with status: "
SAYZ WIT RESPONSE STATUS
```

**Note:** URL must include protocol (http:// or https://)

**Note:** Data is sent as the request body

**Note:** Applies any headers set on the client

**Note:** Content-Type header should be set appropriately

#### PUT

Makes an HTTP PUT request with data in the request body to the specified URL.
Returns a RESPONSE object containing status, body, and headers.

**Syntax:** `<client> DO PUT WIT <url> AN WIT <data>`
**Parameters:**
- `url` (STRIN): The URL to request (must include http:// or https://)
- `data` (STRIN): The data to send in the request body

**Example: Update resource with PUT**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
I HAS A VARIABLE UPDATED_DATA TEH STRIN ITZ "{\"name\":\"Alice\",\"email\":\"alice.smith@example.com\"}"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT "https://api.example.com/users/123" AN WIT UPDATED_DATA
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "User updated successfully"
NOPE
SAYZ WIT "Update failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Example: PUT with version control**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "If-Match" AN WIT "\"etag123\""
I HAS A VARIABLE RESOURCE_DATA TEH STRIN ITZ "{\"content\":\"Updated content\",\"version\":2}"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT "https://api.example.com/documents/456" AN WIT RESOURCE_DATA
IZ RESPONSE STATUS SAEM AS 412?
SAYZ WIT "Resource was modified by another client"
NOPE
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Document updated successfully"
NOPE
SAYZ WIT "Update failed"
KTHX
KTHX
```

**Example: Replace entire resource**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
I HAS A VARIABLE COMPLETE_RESOURCE TEH STRIN ITZ "{\"id\":789,\"name\":\"New Name\",\"active\":YEZ,\"tags\":[\"tag1\",\"tag2\"]}"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT "https://api.example.com/resources/789" AN WIT COMPLETE_RESOURCE
SAYZ WIT "Resource replacement status: "
SAYZ WIT RESPONSE STATUS
```

**Example: PUT binary data**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/octet-stream"
I HAS A VARIABLE BINARY_DATA TEH STRIN ITZ "<binary data as string>"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT "https://api.example.com/files/document.pdf" AN WIT BINARY_DATA
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "File uploaded successfully"
NOPE
SAYZ WIT "Upload failed"
KTHX
```

**Note:** URL must include protocol (http:// or https://)

**Note:** PUT typically replaces the entire resource

**Note:** Use for idempotent operations (multiple calls have same effect)

**Note:** Applies any headers set on the client

**Member Variables:**

#### HEADERS

Request headers as key-value pairs stored in a BASKIT.
Headers are applied to all HTTP requests made by this client.


**Example: Set authorization header**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Authorization" AN WIT "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/protected"
SAYZ WIT "Authenticated request status: "
SAYZ WIT RESPONSE STATUS
```

**Example: Set content type and user agent**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
CLIENT HEADERS DO PUT WIT "User-Agent" AN WIT "MyApp/1.0"
CLIENT HEADERS DO PUT WIT "Accept" AN WIT "application/json"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/data"
SAYZ WIT "Request completed with proper headers"
```

**Example: Add custom headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "X-API-Key" AN WIT "secret-api-key-123"
CLIENT HEADERS DO PUT WIT "X-Request-ID" AN WIT "req-456"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://api.example.com/webhook" AN WIT "{\"event\":\"test\"}"
SAYZ WIT "Custom headers sent"
```

**Example: Modify existing headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Authorization" AN WIT "Bearer old-token"
I HAS A VARIABLE RESPONSE1 TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/data"
CLIENT HEADERS DO PUT WIT "Authorization" AN WIT "Bearer new-token"
I HAS A VARIABLE RESPONSE2 TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/data"
SAYZ WIT "Used different auth tokens for requests"
```

**Example: Check current headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Accept" AN WIT "application/json"
I HAS A VARIABLE ACCEPT_HEADER TEH STRIN ITZ CLIENT HEADERS DO GET WIT "Accept"
SAYZ WIT "Accept header: "
SAYZ WIT ACCEPT_HEADER
```

**Example: Clear all headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Authorization" AN WIT "Bearer token123"
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
I HAS A VARIABLE NEW_HEADERS TEH BASKIT ITZ NEW BASKIT
CLIENT HEADERS ITZ NEW_HEADERS
SAYZ WIT "All headers cleared"
```

**Note:** Headers are applied to all requests made by this client

**Note:** Common headers: Authorization, Content-Type, Accept, User-Agent

**Note:** Header names are case-insensitive in HTTP

**Note:** Use BASKIT operations (PUT, GET, HAS) to manage headers

#### TIMEOUT

Request timeout in seconds.
Default is 30 seconds. Set to 0 for no timeout.


**Example: Set custom timeout**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT TIMEOUT ITZ 60
SAYZ WIT "Timeout set to 60 seconds"
```

**Example: Short timeout for quick requests**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT TIMEOUT ITZ 5
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/delay/1"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Request completed within 5 seconds"
NOPE
SAYZ WIT "Request timed out or failed"
KTHX
```

**Example: Long timeout for slow operations**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT TIMEOUT ITZ 300
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://api.example.com/slow-operation" AN WIT "<large data>"
SAYZ WIT "Long operation completed with status: "
SAYZ WIT RESPONSE STATUS
```

**Example: Disable timeout**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT TIMEOUT ITZ 0
SAYZ WIT "Timeout disabled - request will wait indefinitely"
```

**Example: Check current timeout**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
SAYZ WIT "Current timeout: "
SAYZ WIT CLIENT TIMEOUT
SAYZ WIT " seconds"
```

**Note:** Default timeout is 30 seconds

**Note:** Setting to 0 disables timeout (not recommended)

**Note:** Timeout applies to entire request (connect + read + write)

**Note:** Network errors may occur before timeout

**Example: Basic GET request**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/get"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Request successful!"
SAYZ WIT RESPONSE BODY
NOPE
SAYZ WIT "Request failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Example: POST request with JSON data**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE DATA TEH STRIN ITZ "{\"name\":\"Alice\",\"age\":30}"
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://httpbin.org/post" AN WIT DATA
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON
SAYZ WIT "Response received"
NOPE
SAYZ WIT "POST request failed"
KTHX
```

**Example: Configure custom headers and timeout**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
CLIENT HEADERS DO PUT WIT "Authorization" AN WIT "Bearer token123"
CLIENT HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
CLIENT TIMEOUT ITZ 10
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/data"
SAYZ WIT "Response status: "
SAYZ WIT RESPONSE STATUS
```

### RESPONSE Class

HTTP response object containing status information, headers, and response body.
Returned by all HTTP request methods (GET, POST, PUT, DELETE).

**Methods:**

#### TO_JSON

Parses the response body as JSON and returns a BASKIT.
Throws an exception if the response body is not valid JSON.

**Syntax:** `<response> DO TO_JSON`
**Example: Parse simple JSON object**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/json"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON
SAYZ WIT "JSON parsed successfully"
NOPE
SAYZ WIT "Failed to get JSON response"
KTHX
```

**Example: Access JSON object properties**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/user/123"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE USER_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON
I HAS A VARIABLE USERNAME TEH STRIN ITZ USER_DATA DO GET WIT "username"
I HAS A VARIABLE EMAIL TEH STRIN ITZ USER_DATA DO GET WIT "email"
SAYZ WIT "User: "
SAYZ WIT USERNAME
SAYZ WIT " ("
SAYZ WIT EMAIL
SAYZ WIT ")"
NOPE
SAYZ WIT "Failed to get user data"
KTHX
```

**Example: Handle JSON arrays**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/users"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE USERS TEH BUKKIT ITZ RESPONSE DO TO_JSON
SAYZ WIT "Found "
SAYZ WIT USERS LENGTH
SAYZ WIT " users"
WHILE NO SAEM AS (USERS LENGTH SAEM AS 0)
I HAS A VARIABLE USER TEH BASKIT ITZ USERS DO POP
I HAS A VARIABLE NAME TEH STRIN ITZ USER DO GET WIT "name"
SAYZ WIT "- "
SAYZ WIT NAME
KTHX
NOPE
SAYZ WIT "Failed to get users list"
KTHX
```

**Example: Parse nested JSON structures**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/complex-data"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE DATA TEH BASKIT ITZ RESPONSE DO TO_JSON
I HAS A VARIABLE METADATA TEH BASKIT ITZ DATA DO GET WIT "metadata"
I HAS A VARIABLE VERSION TEH STRIN ITZ METADATA DO GET WIT "version"
I HAS A VARIABLE ITEMS TEH BUKKIT ITZ DATA DO GET WIT "items"
SAYZ WIT "API version: "
SAYZ WIT VERSION
SAYZ WIT "Number of items: "
SAYZ WIT ITEMS LENGTH
NOPE
SAYZ WIT "Failed to parse complex JSON"
KTHX
```

**Example: Error handling for invalid JSON**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/html"
IZ RESPONSE IS_SUCCESS?
MAYB
I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON
SAYZ WIT "Unexpected JSON parsing success"
OOPSIE ERR
SAYZ WIT "Expected error parsing HTML as JSON: "
SAYZ WIT ERR
KTHX
NOPE
SAYZ WIT "Request failed"
KTHX
```

**Note:** Converts JSON objects to BASKIT, arrays to BUKKIT

**Note:** Primitive values (strings, numbers, booleans) remain as-is

**Note:** null values become NOTHIN

**Note:** Nested structures are preserved

**Member Variables:**

#### BODY

Raw response body as a string.
Contains the full response content from the server.


**Example: Get and display response body**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/json"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Response body:"
SAYZ WIT RESPONSE BODY
NOPE
SAYZ WIT "Request failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Example: Parse JSON response body**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/json"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE JSON_DATA TEH OBJECT ITZ RESPONSE DO TO_JSON
SAYZ WIT "JSON parsed successfully"
NOPE
SAYZ WIT "Failed to get JSON data"
KTHX
```

**Example: Handle different content types**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/html"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE CONTENT_TYPE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "content-type"
SAYZ WIT "Content-Type: "
SAYZ WIT CONTENT_TYPE
SAYZ WIT "Body length: "
SAYZ WIT RESPONSE BODY LENGTH
SAYZ WIT "Body preview:"
SAYZ WIT RESPONSE BODY SUBSTRIN WIT 0 AN WIT 100
NOPE
SAYZ WIT "Failed to fetch content"
KTHX
```

**Example: Save response body to file**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/bytes/1024"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE FILE TEH DOCUMENT ITZ NEW DOCUMENT WIT "response.bin"
FILE DO WRITE WIT RESPONSE BODY
FILE DO CLOSE
SAYZ WIT "Response saved to file"
NOPE
SAYZ WIT "Failed to download file"
KTHX
```

**Example: Check body content for specific text**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/get"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE HAS_JSON TEH BOOL ITZ RESPONSE BODY HAS "json"
I HAS A VARIABLE HAS_URL TEH BOOL ITZ RESPONSE BODY HAS "url"
SAYZ WIT "Contains 'json': "
SAYZ WIT HAS_JSON
SAYZ WIT "Contains 'url': "
SAYZ WIT HAS_URL
NOPE
SAYZ WIT "Request failed"
KTHX
```

**Note:** Body content is always returned as a string, even for binary data

**Note:** For JSON responses, use TO_JSON method to parse into objects

**Note:** Large response bodies may impact memory usage

**Note:** Check Content-Type header to understand the response format

**Note:** Body may be empty for some responses (e.g., 204 No Content)

#### HEADERS

Response headers as key-value pairs stored in a BASKIT.
Contains all HTTP headers returned by the server.


**Example: Access common response headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/headers"
I HAS A VARIABLE CONTENT_TYPE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "content-type"
I HAS A VARIABLE SERVER TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "server"
I HAS A VARIABLE CONTENT_LENGTH TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "content-length"
SAYZ WIT "Content-Type: "
SAYZ WIT CONTENT_TYPE
SAYZ WIT "Server: "
SAYZ WIT SERVER
SAYZ WIT "Content-Length: "
SAYZ WIT CONTENT_LENGTH
```

**Example: Check for specific headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/json"
IZ RESPONSE HEADERS DO HAS WIT "content-type"?
I HAS A VARIABLE CT TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "content-type"
SAYZ WIT "Content-Type header found: "
SAYZ WIT CT
NOPE
SAYZ WIT "Content-Type header not found"
KTHX
```

**Example: Iterate through all response headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/headers"
I HAS A VARIABLE HEADER_KEYS TEH BUKKIT ITZ RESPONSE HEADERS DO KEYS
SAYZ WIT "Response headers:"
WHILE NO SAEM AS (HEADER_KEYS LENGTH SAEM AS 0)
I HAS A VARIABLE KEY TEH STRIN ITZ HEADER_KEYS DO POP
I HAS A VARIABLE VALUE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT KEY
SAYZ WIT KEY
SAYZ WIT ": "
SAYZ WIT VALUE
KTHX
```

**Example: Handle CORS headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/data"
I HAS A VARIABLE ALLOW_ORIGIN TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "access-control-allow-origin"
I HAS A VARIABLE ALLOW_METHODS TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "access-control-allow-methods"
IZ ALLOW_ORIGIN SAEM AS "*"?
SAYZ WIT "CORS allows all origins"
NOPE
SAYZ WIT "CORS origin: "
SAYZ WIT ALLOW_ORIGIN
KTHX
SAYZ WIT "Allowed methods: "
SAYZ WIT ALLOW_METHODS
```

**Example: Check cache headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/cache/300"
I HAS A VARIABLE CACHE_CONTROL TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "cache-control"
I HAS A VARIABLE ETAG TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "etag"
SAYZ WIT "Cache-Control: "
SAYZ WIT CACHE_CONTROL
SAYZ WIT "ETag: "
SAYZ WIT ETAG
```

**Note:** Header names are case-insensitive in HTTP

**Note:** Common headers: content-type, content-length, server, date, cache-control

**Note:** Use BASKIT operations (GET, HAS, KEYS) to access headers

**Note:** Headers may not be present for all responses

**Note:** Some headers may have multiple values (comma-separated)

#### IS_ERROR

YEZ if the HTTP status code indicates an error (400+), NO otherwise.
Includes both client errors (4xx) and server errors (5xx).


**Example: Basic error checking**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/status/404"
IZ RESPONSE IS_ERROR?
SAYZ WIT "Request failed with status: "
SAYZ WIT RESPONSE STATUS
NOPE
SAYZ WIT "Request successful"
KTHX
```

**Example: Handle different error types**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE URLS TEH BUKKIT ITZ NEW BUKKIT
URLS DO PUSH WIT "https://httpbin.org/status/400"
URLS DO PUSH WIT "https://httpbin.org/status/404"
URLS DO PUSH WIT "https://httpbin.org/status/500"
WHILE NO SAEM AS (URLS LENGTH SAEM AS 0)
I HAS A VARIABLE URL TEH STRIN ITZ URLS DO POP
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL
IZ RESPONSE IS_ERROR?
I HAS A VARIABLE CODE TEH INTEGR ITZ RESPONSE STATUS
IZ CODE BIGGR THAN 499?
SAYZ WIT "Server error (5xx): "
SAYZ WIT CODE
NOPE
SAYZ WIT "Client error (4xx): "
SAYZ WIT CODE
KTHX
NOPE
SAYZ WIT "Success: "
SAYZ WIT RESPONSE STATUS
KTHX
KTHX
```

**Example: Combine with success checking**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://api.example.com/login" AN WIT "{\"username\":\"bad\",\"password\":\"wrong\"}"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Login successful"
NOPE
IZ RESPONSE IS_ERROR?
I HAS A VARIABLE CODE TEH INTEGR ITZ RESPONSE STATUS
IZ CODE SAEM AS 401?
SAYZ WIT "Authentication failed - invalid credentials"
NOPE
IZ CODE SAEM AS 403?
SAYZ WIT "Access forbidden - insufficient permissions"
NOPE
IZ CODE SAEM AS 429?
SAYZ WIT "Too many requests - rate limited"
NOPE
SAYZ WIT "Other error: "
SAYZ WIT CODE
KTHX
KTHX
KTHX
NOPE
SAYZ WIT "Unexpected status: "
SAYZ WIT RESPONSE STATUS
KTHX
KTHX
```

**Example: Error handling in API calls**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/user/999999"
IZ RESPONSE IS_ERROR?
I HAS A VARIABLE ERROR_MSG TEH STRIN
I HAS A VARIABLE CODE TEH INTEGR ITZ RESPONSE STATUS
IZ CODE SAEM AS 404?
ERROR_MSG ITZ "User not found"
NOPE
IZ CODE SAEM AS 403?
ERROR_MSG ITZ "Access denied"
NOPE
IZ CODE BIGGR THAN 499?
ERROR_MSG ITZ "Server error - please try again later"
NOPE
ERROR_MSG ITZ "Request failed"
KTHX
KTHX
KTHX
SAYZ WIT "Error: "
SAYZ WIT ERROR_MSG
SAYZ WIT " (status: "
SAYZ WIT CODE
SAYZ WIT ")"
NOPE
SAYZ WIT "User data retrieved successfully"
KTHX
```

**Example: Retry logic with error checking**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE MAX_RETRIES TEH INTEGR ITZ 3
I HAS A VARIABLE RETRY_COUNT TEH INTEGR ITZ 0
I HAS A VARIABLE SUCCESS TEH BOOL ITZ NO
WHILE NO SAEM AS (SUCCESS SAEM AS YEZ) AN NO SAEM AS (RETRY_COUNT SAEM AS MAX_RETRIES)
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/unreliable-endpoint"
IZ RESPONSE IS_ERROR?
RETRY_COUNT ITZ RETRY_COUNT UP 1
IZ RETRY_COUNT SAEM AS MAX_RETRIES?
SAYZ WIT "Max retries reached - giving up"
NOPE
SAYZ WIT "Request failed (attempt "
SAYZ WIT RETRY_COUNT
SAYZ WIT "/"
SAYZ WIT MAX_RETRIES
SAYZ WIT ") - retrying..."
KTHX
NOPE
SUCCESS ITZ YEZ
SAYZ WIT "Request successful after "
SAYZ WIT RETRY_COUNT
SAYZ WIT " retries"
KTHX
KTHX
```

**Note:** Returns YEZ for status codes 400 and above

**Note:** Includes both client errors (4xx) and server errors (5xx)

**Note:** Use IS_SUCCESS for checking success conditions (200-299)

**Note:** Error responses may still contain useful information in the body

**Note:** Some APIs return error details in JSON format

#### IS_SUCCESS

YEZ if the HTTP status code indicates success (200-299), NO otherwise.
Provides a convenient way to check if the request was successful.


**Example: Basic success check**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/get"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Request successful!"
NOPE
SAYZ WIT "Request failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Example: Handle different success scenarios**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE URLS TEH BUKKIT ITZ NEW BUKKIT
URLS DO PUSH WIT "https://httpbin.org/status/200"
URLS DO PUSH WIT "https://httpbin.org/status/201"
URLS DO PUSH WIT "https://httpbin.org/status/204"
I HAS A VARIABLE SUCCESS_COUNT TEH INTEGR ITZ 0
WHILE NO SAEM AS (URLS LENGTH SAEM AS 0)
I HAS A VARIABLE URL TEH STRIN ITZ URLS DO POP
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL
IZ RESPONSE IS_SUCCESS?
SUCCESS_COUNT ITZ SUCCESS_COUNT UP 1
SAYZ WIT "Success: "
SAYZ WIT RESPONSE STATUS
NOPE
SAYZ WIT "Failed: "
SAYZ WIT RESPONSE STATUS
KTHX
KTHX
SAYZ WIT "Total successful requests: "
SAYZ WIT SUCCESS_COUNT
```

**Example: Combine with error checking**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://api.example.com/login" AN WIT "{\"username\":\"user\",\"password\":\"pass\"}"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Login successful"
I HAS A VARIABLE USER_DATA TEH OBJECT ITZ RESPONSE DO TO_JSON
NOPE
IZ RESPONSE IS_ERROR?
SAYZ WIT "Login failed - authentication error"
NOPE
SAYZ WIT "Login failed - unknown error (status: "
SAYZ WIT RESPONSE STATUS
SAYZ WIT ")"
KTHX
KTHX
```

**Example: Use in conditional logic**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/user/profile"
IZ RESPONSE IS_SUCCESS AN RESPONSE STATUS SAEM AS 200?
SAYZ WIT "Profile retrieved successfully"
NOPE
IZ RESPONSE IS_SUCCESS AN RESPONSE STATUS SAEM AS 201?
SAYZ WIT "Profile created successfully"
NOPE
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Profile operation successful (status: "
SAYZ WIT RESPONSE STATUS
SAYZ WIT ")"
NOPE
SAYZ WIT "Profile operation failed"
KTHX
KTHX
KTHX
```

**Example: Batch processing with success tracking**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE IDS TEH BUKKIT ITZ NEW BUKKIT
IDS DO PUSH WIT "1"
IDS DO PUSH WIT "2"
IDS DO PUSH WIT "3"
I HAS A VARIABLE RESULTS TEH BUKKIT ITZ NEW BUKKIT
WHILE NO SAEM AS (IDS LENGTH SAEM AS 0)
I HAS A VARIABLE ID TEH STRIN ITZ IDS DO POP
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://api.example.com/items/" MOAR ID
IZ RESPONSE IS_SUCCESS?
RESULTS DO PUSH WIT "SUCCESS: " MOAR ID
NOPE
RESULTS DO PUSH WIT "FAILED: " MOAR ID MOAR " (" MOAR RESPONSE STATUS MOAR ")"
KTHX
KTHX
SAYZ WIT "Batch results:"
WHILE NO SAEM AS (RESULTS LENGTH SAEM AS 0)
SAYZ WIT RESULTS DO POP
KTHX
```

**Note:** Returns YEZ for status codes 200-299

**Note:** Equivalent to checking if STATUS >= 200 AND STATUS < 300

**Note:** Use IS_ERROR for checking error conditions (400+)

**Note:** Success doesn't guarantee the response body contains expected data

**Note:** Some APIs return 200 with error information in the body

#### STATUS

HTTP status code returned by the server.
Examples: 200 (OK), 404 (Not Found), 500 (Server Error).


**Example: Check specific status codes**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/status/201"
IZ RESPONSE STATUS SAEM AS 201?
SAYZ WIT "Resource created successfully"
NOPE
IZ RESPONSE STATUS SAEM AS 404?
SAYZ WIT "Resource not found"
NOPE
SAYZ WIT "Other status: "
SAYZ WIT RESPONSE STATUS
KTHX
KTHX
```

**Example: Categorize status codes**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/status/500"
I HAS A VARIABLE CODE TEH INTEGR ITZ RESPONSE STATUS
IZ CODE BIGGR THAN 499?
SAYZ WIT "Server error (5xx)"
NOPE
IZ CODE BIGGR THAN 399?
SAYZ WIT "Client error (4xx)"
NOPE
IZ CODE BIGGR THAN 299?
SAYZ WIT "Redirection (3xx)"
NOPE
IZ CODE BIGGR THAN 199?
SAYZ WIT "Success (2xx)"
NOPE
SAYZ WIT "Informational (1xx)"
KTHX
KTHX
KTHX
KTHX
```

**Example: Handle common HTTP statuses**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT "https://api.example.com/login" AN WIT "{\"username\":\"bad\",\"password\":\"wrong\"}"
IZ RESPONSE STATUS SAEM AS 401?
SAYZ WIT "Authentication failed - check credentials"
NOPE
IZ RESPONSE STATUS SAEM AS 403?
SAYZ WIT "Access forbidden - insufficient permissions"
NOPE
IZ RESPONSE STATUS SAEM AS 429?
SAYZ WIT "Too many requests - rate limited"
NOPE
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Login successful"
NOPE
SAYZ WIT "Login failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
KTHX
KTHX
KTHX
```

**Example: Log status codes for debugging**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE URLS TEH BUKKIT ITZ NEW BUKKIT
URLS DO PUSH WIT "https://httpbin.org/status/200"
URLS DO PUSH WIT "https://httpbin.org/status/404"
URLS DO PUSH WIT "https://httpbin.org/status/500"
WHILE NO SAEM AS (URLS LENGTH SAEM AS 0)
I HAS A VARIABLE URL TEH STRIN ITZ URLS DO POP
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL
SAYZ WIT "URL: "
SAYZ WIT URL
SAYZ WIT " -> Status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Note:** 200-299 = Success, 300-399 = Redirection, 400-499 = Client Error, 500-599 = Server Error

**Note:** Common codes: 200 OK, 201 Created, 301 Moved, 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 500 Internal Server Error

**Note:** Use IS_SUCCESS and IS_ERROR for easy status checking

**Note:** Status codes are standardized by HTTP specification

**Example: Basic response handling**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/get"
IZ RESPONSE IS_SUCCESS?
SAYZ WIT "Request successful!"
SAYZ WIT "Status: "
SAYZ WIT RESPONSE STATUS
SAYZ WIT "Body: "
SAYZ WIT RESPONSE BODY
NOPE
SAYZ WIT "Request failed with status: "
SAYZ WIT RESPONSE STATUS
KTHX
```

**Example: Parse JSON response**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/json"
IZ RESPONSE IS_SUCCESS?
I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON
SAYZ WIT "JSON parsed successfully"
NOPE
SAYZ WIT "Failed to get JSON data"
KTHX
```

**Example: Check response headers**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/headers"
I HAS A VARIABLE CONTENT_TYPE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "content-type"
I HAS A VARIABLE SERVER TEH STRIN ITZ RESPONSE HEADERS DO GET WIT "server"
SAYZ WIT "Content-Type: "
SAYZ WIT CONTENT_TYPE
SAYZ WIT "Server: "
SAYZ WIT SERVER
```

**Example: Handle different HTTP status codes**

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT "https://httpbin.org/status/404"
IZ RESPONSE STATUS SAEM AS 200?
SAYZ WIT "OK"
NOPE
IZ RESPONSE STATUS SAEM AS 404?
SAYZ WIT "Not Found"
NOPE
IZ RESPONSE STATUS SAEM AS 500?
SAYZ WIT "Server Error"
NOPE
SAYZ WIT "Other status: "
SAYZ WIT RESPONSE STATUS
KTHX
KTHX
KTHX
```

