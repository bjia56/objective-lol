# HTTP Module - HTTP Client Operations

The HTTP module provides HTTP client functionality through the INTERWEB class for making web requests and the RESPONSE class for handling responses.

## Importing HTTP Module

```lol
BTW Import entire module
I CAN HAS HTTP?

BTW Selective import
I CAN HAS INTERWEB FROM HTTP?
I CAN HAS RESPONSE FROM HTTP?
```

**Note:** The HTTP module automatically imports the RESPONSE class when INTERWEB is imported.

## INTERWEB Class

The INTERWEB class represents an HTTP client that can make requests to web servers. It supports HTTP and HTTPS protocols transparently and provides configurable timeout and header management.

### Constructor

```lol
I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
```

The constructor creates an HTTP client with default settings:
- **Timeout**: 30 seconds
- **Headers**: Empty BASKIT

### Properties

- **TIMEOUT**: INTEGR - Request timeout in seconds (default: 30)
- **HEADERS**: BASKIT - Request headers as key-value pairs

### Methods

#### HTTP Request Methods

##### GET - Make GET Request

Performs an HTTP GET request to the specified URL.

```lol
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ client DO GET WIT <url>
```

**Parameters:**
- **url**: STRIN - The target URL (HTTP or HTTPS)

**Returns:** RESPONSE - The HTTP response object

##### POST - Make POST Request

Performs an HTTP POST request with data to the specified URL.

```lol
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ client DO POST WIT <url> AN WIT <data>
```

**Parameters:**
- **url**: STRIN - The target URL (HTTP or HTTPS)
- **data**: STRIN - The request body data

**Returns:** RESPONSE - The HTTP response object

##### PUT - Make PUT Request

Performs an HTTP PUT request with data to the specified URL.

```lol
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ client DO PUT WIT <url> AN WIT <data>
```

**Parameters:**
- **url**: STRIN - The target URL (HTTP or HTTPS)
- **data**: STRIN - The request body data

**Returns:** RESPONSE - The HTTP response object

##### DELETE - Make DELETE Request

Performs an HTTP DELETE request to the specified URL.

```lol
I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ client DO DELETE WIT <url>
```

**Parameters:**
- **url**: STRIN - The target URL (HTTP or HTTPS)

**Returns:** RESPONSE - The HTTP response object

## RESPONSE Class

The RESPONSE class represents an HTTP response containing status information, headers, and response body.

### Properties

- **STATUS**: INTEGR (read-only) - HTTP status code (e.g., 200, 404, 500)
- **BODY**: STRIN (read-only) - Response body content
- **HEADERS**: BASKIT (read-only) - Response headers as key-value pairs
- **IS_SUCCESS**: BOOL (read-only) - True for status codes 200-299
- **IS_ERROR**: BOOL (read-only) - True for status codes 400+

### Methods

##### TO_JSON - Parse JSON Response

Parses the response body as JSON and returns a BASKIT.

```lol
I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ response DO TO_JSON
```

**Returns:** BASKIT - The parsed JSON data as a BASKIT

**Throws:** Exception if response body is not valid JSON

## Basic HTTP Operations

### Simple GET Request

```lol
I CAN HAS HTTP?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN SIMPLE_GET_REQUEST WIT URL TEH STRIN
    SAYZ WIT "=== Simple GET Request ==="

    BTW Create HTTP client
    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB

    MAYB
        BTW Make GET request
        I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL

        SAYZ WIT "Request successful!"
        SAY WIT "Status Code: "
        SAYZ WIT RESPONSE STATUS

        SAY WIT "Response Body: "
        SAYZ WIT RESPONSE BODY

    OOPSIE HTTP_ERROR
        SAY WIT "Request failed: "
        SAYZ WIT HTTP_ERROR
    KTHX
KTHXBAI
```

### POST Request with Data

```lol
I CAN HAS HTTP?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN POST_JSON_DATA WIT URL TEH STRIN AN WIT JSON_DATA TEH STRIN
    SAYZ WIT "=== POST JSON Data ==="

    BTW Create HTTP client
    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB

    BTW Configure headers
    I HAS A VARIABLE HEADERS TEH BASKIT ITZ CLIENT HEADERS
    HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
    HEADERS DO PUT WIT "Accept" AN WIT "application/json"

    MAYB
        BTW Make POST request
        I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT URL AN WIT JSON_DATA

        IZ RESPONSE IS_SUCCESS?
            SAYZ WIT "POST successful!"
            SAY WIT "Status Code: "
            SAYZ WIT RESPONSE STATUS

            BTW Parse JSON response if applicable
            MAYB
                I HAS A VARIABLE RESULT TEH BASKIT ITZ RESPONSE DO TO_JSON
                SAYZ WIT "JSON response received"
            OOPSIE JSON_ERROR
                SAYZ WIT "Response is not JSON"
            KTHX

        NOPE
            SAYZ WIT "POST failed"
            SAY WIT "Status Code: "
            SAYZ WIT RESPONSE STATUS
        KTHX

    OOPSIE HTTP_ERROR
        SAY WIT "Request failed: "
        SAYZ WIT HTTP_ERROR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_POST_REQUEST
    I HAS A VARIABLE JSON_DATA TEH STRIN ITZ "{\"name\": \"Alice\", \"age\": 30}"
    POST_JSON_DATA WIT "https://httpbin.org/post" AN WIT JSON_DATA
KTHXBAI
```

### Request with Custom Headers

```lol
I CAN HAS HTTP?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN AUTHENTICATED_REQUEST WIT URL TEH STRIN AN WIT TOKEN TEH STRIN
    SAYZ WIT "=== Authenticated Request ==="

    BTW Create HTTP client with custom configuration
    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB

    BTW Set timeout to 60 seconds
    CLIENT TIMEOUT ITZ 60

    BTW Configure authentication headers
    I HAS A VARIABLE HEADERS TEH BASKIT ITZ CLIENT HEADERS
    HEADERS DO PUT WIT "Authorization" AN WIT "Bearer " + TOKEN
    HEADERS DO PUT WIT "User-Agent" AN WIT "ObjectiveLOL-Client/1.0"
    HEADERS DO PUT WIT "Accept" AN WIT "application/json"

    MAYB
        I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL

        IZ RESPONSE IS_SUCCESS?
            SAYZ WIT "Authenticated request successful!"

            BTW Show response headers
            I HAS A VARIABLE RESP_HEADERS TEH BASKIT ITZ RESPONSE HEADERS
            MAYB
                SAY WIT "Server: "
                SAYZ WIT RESP_HEADERS DO GET WIT "Server"
            OOPSIE HEADER_ERROR
                SAYZ WIT "Server header not found"
            KTHX

        NOPE
            SAYZ WIT "Authentication failed"
            SAY WIT "Status: "
            SAYZ WIT RESPONSE STATUS
        KTHX

    OOPSIE HTTP_ERROR
        SAY WIT "Request failed: "
        SAYZ WIT HTTP_ERROR
    KTHX
KTHXBAI
```

## Advanced HTTP Operations

### JSON API Client

```lol
I CAN HAS HTTP?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN API_CLIENT_DEMO WIT BASE_URL TEH STRIN
    SAYZ WIT "=== JSON API Client Demo ==="

    BTW Create configured HTTP client
    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
    CLIENT TIMEOUT ITZ 30

    I HAS A VARIABLE HEADERS TEH BASKIT ITZ CLIENT HEADERS
    HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
    HEADERS DO PUT WIT "Accept" AN WIT "application/json"
    HEADERS DO PUT WIT "User-Agent" AN WIT "ObjectiveLOL-API-Client/1.0"

    MAYB
        BTW Get API information
        I HAS A VARIABLE INFO_URL TEH STRIN ITZ BASE_URL
        I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT INFO_URL

        IZ RESPONSE IS_SUCCESS?
            SAYZ WIT "API connection successful"

            BTW Parse JSON response
            I HAS A VARIABLE API_INFO TEH BASKIT ITZ RESPONSE DO TO_JSON

            BTW Display API information
            MAYB
                SAY WIT "Current User URL: "
                SAYZ WIT API_INFO DO GET WIT "current_user_url"
            OOPSIE KEY_ERROR
                SAYZ WIT "API structure may have changed"
            KTHX

            BTW Show response metadata
            SAY WIT "Response size: "
            SAYZ WIT LEN WIT RESPONSE BODY

            I HAS A VARIABLE RESP_HEADERS TEH BASKIT ITZ RESPONSE HEADERS
            MAYB
                SAY WIT "Content-Type: "
                SAYZ WIT RESP_HEADERS DO GET WIT "Content-Type"
            OOPSIE HEADER_ERROR
                SAYZ WIT "Content-Type header not found"
            KTHX

        NOPE
            SAYZ WIT "API connection failed"
            IZ RESPONSE IS_ERROR?
                SAY WIT "Error status: "
                SAYZ WIT RESPONSE STATUS
            KTHX
        KTHX

    OOPSIE HTTP_ERROR
        SAY WIT "Network error: "
        SAYZ WIT HTTP_ERROR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_GITHUB_API
    API_CLIENT_DEMO WIT "https://api.github.com"
KTHXBAI
```

### File Upload Simulation

```lol
I CAN HAS HTTP?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN UPLOAD_DATA WIT UPLOAD_URL TEH STRIN AN WIT DATA TEH STRIN
    SAYZ WIT "=== File Upload Simulation ==="

    BTW Create HTTP client for upload
    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
    CLIENT TIMEOUT ITZ 120  BTW Longer timeout for uploads

    BTW Configure upload headers
    I HAS A VARIABLE HEADERS TEH BASKIT ITZ CLIENT HEADERS
    HEADERS DO PUT WIT "Content-Type" AN WIT "text/plain"
    HEADERS DO PUT WIT "Content-Length" AN WIT LEN WIT DATA

    MAYB
        BTW Simulate file upload with PUT
        I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT UPLOAD_URL AN WIT DATA

        IZ RESPONSE IS_SUCCESS?
            SAYZ WIT "Upload successful!"
            SAY WIT "Status: "
            SAYZ WIT RESPONSE STATUS

            BTW Show upload confirmation
            SAY WIT "Response: "
            SAYZ WIT RESPONSE BODY

        NOPE
            SAYZ WIT "Upload failed"
            SAY WIT "Status: "
            SAYZ WIT RESPONSE STATUS
            SAY WIT "Error: "
            SAYZ WIT RESPONSE BODY
        KTHX

    OOPSIE UPLOAD_ERROR
        SAY WIT "Upload error: "
        SAYZ WIT UPLOAD_ERROR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_UPLOAD
    I HAS A VARIABLE SAMPLE_DATA TEH STRIN ITZ "Hello, World!\nThis is sample upload data.\nGenerated by Objective-LOL!"
    UPLOAD_DATA WIT "https://httpbin.org/put" AN WIT SAMPLE_DATA
KTHXBAI
```

### RESTful API Operations

```lol
I CAN HAS HTTP?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN REST_API_DEMO WIT API_BASE TEH STRIN
    SAYZ WIT "=== RESTful API Operations ==="

    BTW Create API client
    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
    CLIENT TIMEOUT ITZ 30

    I HAS A VARIABLE HEADERS TEH BASKIT ITZ CLIENT HEADERS
    HEADERS DO PUT WIT "Content-Type" AN WIT "application/json"
    HEADERS DO PUT WIT "Accept" AN WIT "application/json"

    BTW GET - List resources
    SAYZ WIT "1. GET - Listing resources"
    MAYB
        I HAS A VARIABLE LIST_URL TEH STRIN ITZ API_BASE + "/users"
        I HAS A VARIABLE GET_RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT LIST_URL

        IZ GET_RESPONSE IS_SUCCESS?
            SAYZ WIT "GET successful"
            SAY WIT "Status: "
            SAYZ WIT GET_RESPONSE STATUS
        NOPE
            SAYZ WIT "GET failed"
        KTHX
    OOPSIE GET_ERROR
        SAYZ WIT "GET operation failed"
    KTHX

    BTW POST - Create resource
    SAYZ WIT "2. POST - Creating resource"
    MAYB
        I HAS A VARIABLE CREATE_URL TEH STRIN ITZ API_BASE + "/users"
        I HAS A VARIABLE USER_DATA TEH STRIN ITZ "{\"name\": \"Test User\", \"email\": \"test@example.com\"}"
        I HAS A VARIABLE POST_RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT CREATE_URL AN WIT USER_DATA

        IZ POST_RESPONSE IS_SUCCESS?
            SAYZ WIT "POST successful"
            SAY WIT "Status: "
            SAYZ WIT POST_RESPONSE STATUS
        NOPE
            SAYZ WIT "POST failed"
        KTHX
    OOPSIE POST_ERROR
        SAYZ WIT "POST operation failed"
    KTHX

    BTW PUT - Update resource
    SAYZ WIT "3. PUT - Updating resource"
    MAYB
        I HAS A VARIABLE UPDATE_URL TEH STRIN ITZ API_BASE + "/users/1"
        I HAS A VARIABLE UPDATE_DATA TEH STRIN ITZ "{\"name\": \"Updated User\", \"email\": \"updated@example.com\"}"
        I HAS A VARIABLE PUT_RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT UPDATE_URL AN WIT UPDATE_DATA

        IZ PUT_RESPONSE IS_SUCCESS?
            SAYZ WIT "PUT successful"
            SAY WIT "Status: "
            SAYZ WIT PUT_RESPONSE STATUS
        NOPE
            SAYZ WIT "PUT failed"
        KTHX
    OOPSIE PUT_ERROR
        SAYZ WIT "PUT operation failed"
    KTHX

    BTW DELETE - Remove resource
    SAYZ WIT "4. DELETE - Removing resource"
    MAYB
        I HAS A VARIABLE DELETE_URL TEH STRIN ITZ API_BASE + "/users/1"
        I HAS A VARIABLE DELETE_RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT DELETE_URL

        IZ DELETE_RESPONSE IS_SUCCESS?
            SAYZ WIT "DELETE successful"
            SAY WIT "Status: "
            SAYZ WIT DELETE_RESPONSE STATUS
        NOPE
            SAYZ WIT "DELETE failed"
        KTHX
    OOPSIE DELETE_ERROR
        SAYZ WIT "DELETE operation failed"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_REST_API
    REST_API_DEMO WIT "https://jsonplaceholder.typicode.com"
KTHXBAI
```

## Error Handling

### Network Error Handling

```lol
I CAN HAS HTTP?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN ROBUST_HTTP_REQUEST WIT URL TEH STRIN
    SAYZ WIT "=== Robust HTTP Request ==="

    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB
    CLIENT TIMEOUT ITZ 10  BTW Short timeout for demo

    MAYB
        I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL

        BTW Check response status
        IZ RESPONSE IS_SUCCESS?
            SAYZ WIT "Request successful!"
            SAY WIT "Status: "
            SAYZ WIT RESPONSE STATUS

        NOPE IZ RESPONSE IS_ERROR?
            SAYZ WIT "Client or server error"
            SAY WIT "Status: "
            SAYZ WIT RESPONSE STATUS

            BTW Handle specific error codes
            IZ RESPONSE STATUS SAEM AS 404?
                SAYZ WIT "Resource not found"
            NOPE IZ RESPONSE STATUS SAEM AS 401?
                SAYZ WIT "Authentication required"
            NOPE IZ RESPONSE STATUS SAEM AS 500?
                SAYZ WIT "Server internal error"
            NOPE
                SAYZ WIT "Other error occurred"
            KTHX

        NOPE
            SAYZ WIT "Unexpected response status"
            SAY WIT "Status: "
            SAYZ WIT RESPONSE STATUS
        KTHX

    OOPSIE HTTP_ERROR
        SAYZ WIT "HTTP operation failed: "
        SAYZ WIT HTTP_ERROR
        SAYZ WIT "This could be due to:"
        SAYZ WIT "- Network connectivity issues"
        SAYZ WIT "- Invalid URL"
        SAYZ WIT "- Timeout"
        SAYZ WIT "- DNS resolution failure"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_ERROR_HANDLING
    BTW Test with valid URL
    ROBUST_HTTP_REQUEST WIT "https://httpbin.org/status/200"

    BTW Test with error status
    ROBUST_HTTP_REQUEST WIT "https://httpbin.org/status/404"

    BTW Test with invalid URL
    ROBUST_HTTP_REQUEST WIT "invalid-url"

    BTW Test with timeout (very slow response)
    ROBUST_HTTP_REQUEST WIT "https://httpbin.org/delay/15"
KTHXBAI
```

### JSON Parsing Error Handling

```lol
I CAN HAS HTTP?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN SAFE_JSON_PARSING WIT URL TEH STRIN
    SAYZ WIT "=== Safe JSON Parsing ==="

    I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB

    MAYB
        I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL

        IZ RESPONSE IS_SUCCESS?
            SAYZ WIT "Response received, attempting JSON parse..."

            MAYB
                I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON
                SAYZ WIT "JSON parsing successful!"

                BTW Safe access to JSON fields
                MAYB
                    I HAS A VARIABLE NAME TEH STRIN ITZ JSON_DATA DO GET WIT "name"
                    SAY WIT "Name: "
                    SAYZ WIT NAME
                OOPSIE FIELD_ERROR
                    SAYZ WIT "Name field not found in JSON"
                KTHX

            OOPSIE JSON_ERROR
                SAYZ WIT "JSON parsing failed: "
                SAYZ WIT JSON_ERROR
                SAYZ WIT "Response may not be valid JSON"
                SAYZ WIT "Response body preview:"
                SAYZ WIT RESPONSE BODY
            KTHX

        NOPE
            SAYZ WIT "HTTP request failed, no JSON to parse"
        KTHX

    OOPSIE REQUEST_ERROR
        SAYZ WIT "Request failed: "
        SAYZ WIT REQUEST_ERROR
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_JSON_HANDLING
    BTW Test with valid JSON
    SAFE_JSON_PARSING WIT "https://httpbin.org/json"

    BTW Test with non-JSON response
    SAFE_JSON_PARSING WIT "https://httpbin.org/html"
KTHXBAI
```

## Quick Reference

### Constructor

| Usage | Description |
|-------|-------------|
| `NEW INTERWEB` | Create HTTP client with default settings |

### Configuration

| Property | Type | Description |
|----------|------|-------------|
| `TIMEOUT` | INTEGR | Request timeout in seconds |
| `HEADERS` | BASKIT | Request headers |

### HTTP Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `GET WIT url` | url: STRIN | RESPONSE | Make GET request |
| `POST WIT url AN WIT data` | url: STRIN, data: STRIN | RESPONSE | Make POST request |
| `PUT WIT url AN WIT data` | url: STRIN, data: STRIN | RESPONSE | Make PUT request |
| `DELETE WIT url` | url: STRIN | RESPONSE | Make DELETE request |

### Response Properties

| Property | Type | Description |
|----------|------|-------------|
| `STATUS` | INTEGR | HTTP status code |
| `BODY` | STRIN | Response body content |
| `HEADERS` | BASKIT | Response headers |
| `IS_SUCCESS` | BOOL | True for 200-299 status |
| `IS_ERROR` | BOOL | True for 400+ status |

### Response Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `TO_JSON` | BASKIT | Parse response body as JSON |

## Related

- [STDIO Module](stdio.md) - Console input/output for debugging
- [String Module](string.md) - String manipulation for URLs and data
- [Collections](collections.md) - BASKIT operations for headers and JSON
- [Control Flow](../language-guide/control-flow.md) - Exception handling patterns