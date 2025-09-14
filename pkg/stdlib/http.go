package stdlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleHTTPCategories defines the order that categories should be rendered in documentation
var moduleHTTPCategories = []string{
	"http-requests",
	"http-configuration",
	"response-properties",
	"response-parsing",
}

// InterwebData stores the internal state of an INTERWEB client
type InterwebData struct {
	Client *http.Client
}

// ResponseData stores the internal state of a RESPONSE
type ResponseData struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// Global HTTP class definitions - created once and reused
var httpClassesOnce = sync.Once{}
var httpClasses map[string]*environment.Class

func getHTTPClasses() map[string]*environment.Class {
	httpClassesOnce.Do(func() {
		httpClasses = map[string]*environment.Class{
			"INTERWEB": {
				Name: "INTERWEB",
				Documentation: []string{
					"HTTP client that can make requests to web servers.",
					"Supports HTTP and HTTPS protocols with configurable timeout and headers.",
					"",
					"@class INTERWEB",
					"@example Basic GET request",
					"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
					"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/get\"",
					"IZ RESPONSE IS_SUCCESS?",
					"    SAYZ WIT \"Request successful!\"",
					"    SAYZ WIT RESPONSE BODY",
					"NOPE",
					"    SAYZ WIT \"Request failed with status: \"",
					"    SAYZ WIT RESPONSE STATUS",
					"KTHX",
					"@example POST request with JSON data",
					"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
					"I HAS A VARIABLE DATA TEH STRIN ITZ \"{\\\"name\\\":\\\"Alice\\\",\\\"age\\\":30}\"",
					"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://httpbin.org/post\" AN WIT DATA",
					"IZ RESPONSE IS_SUCCESS?",
					"    I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON",
					"    SAYZ WIT \"Response received\"",
					"NOPE",
					"    SAYZ WIT \"POST request failed\"",
					"KTHX",
					"@example Configure custom headers and timeout",
					"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
					"CLIENT HEADERS DO PUT WIT \"Authorization\" AN WIT \"Bearer token123\"",
					"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/json\"",
					"CLIENT TIMEOUT ITZ 10",
					"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/data\"",
					"SAYZ WIT \"Response status: \"",
					"SAYZ WIT RESPONSE STATUS",
					"@note Default timeout is 30 seconds",
					"@note Supports GET, POST, PUT, and DELETE methods",
					"@note Headers are applied to all requests made by this client",
					"@see RESPONSE",
				},
				QualifiedName: "stdlib:HTTP.INTERWEB",
				ModulePath:    "stdlib:HTTP",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:HTTP.INTERWEB"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"INTERWEB": {
						Name: "INTERWEB",
						Documentation: []string{
							"Initializes an INTERWEB HTTP client with default settings.",
							"Default timeout is 30 seconds with empty headers.",
							"",
							"@syntax NEW INTERWEB",
							"@returns {NOTHIN} No return value (constructor)",
							"@example Create basic HTTP client",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"SAYZ WIT \"HTTP client created with default settings\"",
							"@example Create client and immediately configure",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT TIMEOUT ITZ 60",
							"CLIENT HEADERS DO PUT WIT \"User-Agent\" AN WIT \"MyApp/1.0\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/user-agent\"",
							"SAYZ WIT RESPONSE BODY",
							"@example Create multiple clients with different configurations",
							"I HAS A VARIABLE API_CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"API_CLIENT HEADERS DO PUT WIT \"Authorization\" AN WIT \"Bearer token123\"",
							"API_CLIENT TIMEOUT ITZ 10",
							"I HAS A VARIABLE WEB_CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"WEB_CLIENT HEADERS DO PUT WIT \"User-Agent\" AN WIT \"WebScraper/1.0\"",
							"WEB_CLIENT TIMEOUT ITZ 30",
							"SAYZ WIT \"Created two clients with different configurations\"",
							"@note Initializes with 30-second timeout",
							"@note Creates empty headers BASKIT that can be populated",
							"@note Client is ready to make requests immediately",
							"@see TIMEOUT, HEADERS, GET, POST",
							"@category http-client",
						},
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							// Initialize HTTP client
							client := &http.Client{
								Timeout: 30 * time.Second,
							}

							interwebData := &InterwebData{
								Client: client,
							}
							this.NativeData = interwebData

							// Initialize empty headers BASKIT
							headers := NewBaskitInstance()
							this.Variables["HEADERS"].Value = headers

							return environment.NOTHIN, nil
						},
					},
					// GET method
					"GET": {
						Name: "GET",
						Documentation: []string{
							"Makes an HTTP GET request to the specified URL.",
							"Returns a RESPONSE object containing status, body, and headers.",
							"",
							"@syntax <client> DO GET WIT <url>",
							"@param {STRIN} url - The URL to request (must include http:// or https://)",
							"@returns {RESPONSE} Response object with status, body, and headers",
							"@example Simple GET request",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/get\"",
							"IZ RESPONSE IS_SUCCESS?",
							"    SAYZ WIT \"Success! Status: \"",
							"    SAYZ WIT RESPONSE STATUS",
							"    SAYZ WIT \"Body length: \"",
							"    SAYZ WIT RESPONSE BODY LENGTH",
							"NOPE",
							"    SAYZ WIT \"Request failed with status: \"",
							"    SAYZ WIT RESPONSE STATUS",
							"KTHX",
							"@example GET with custom headers",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"Accept\" AN WIT \"application/json\"",
							"CLIENT HEADERS DO PUT WIT \"User-Agent\" AN WIT \"MyApp/1.0\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.github.com/user\"",
							"SAYZ WIT \"Response status: \"",
							"SAYZ WIT RESPONSE STATUS",
							"@example Handle different response types",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/json\"",
							"IZ RESPONSE IS_SUCCESS?",
							"    I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON",
							"    SAYZ WIT \"Parsed JSON response\"",
							"NOPE",
							"    SAYZ WIT \"Request failed or returned non-JSON content\"",
							"    SAYZ WIT RESPONSE BODY",
							"KTHX",
							"@example Check response headers",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/headers\"",
							"I HAS A VARIABLE CONTENT_TYPE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"content-type\"",
							"SAYZ WIT \"Content-Type: \"",
							"SAYZ WIT CONTENT_TYPE",
							"@throws Exception if URL is invalid or unreachable",
							"@throws Exception if network error occurs",
							"@note URL must include protocol (http:// or https://)",
							"@note Applies any headers set on the client",
							"@note Respects the client's timeout setting",
							"@see POST, PUT, DELETE, RESPONSE",
							"@category http-requests",
						},
						ReturnType: "RESPONSE",
						Parameters: []environment.Parameter{
							{Name: "url", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return executeHTTPRequest(this, "GET", args[0], environment.StringValue(""))
						},
					},
					// POST method
					"POST": {
						Name: "POST",
						Documentation: []string{
							"Makes an HTTP POST request with data in the request body to the specified URL.",
							"Returns a RESPONSE object containing status, body, and headers.",
							"",
							"@syntax <client> DO POST WIT <url> AN WIT <data>",
							"@param {STRIN} url - The URL to request (must include http:// or https://)",
							"@param {STRIN} data - The data to send in the request body",
							"@returns {RESPONSE} Response object with status, body, and headers",
							"@example POST JSON data",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/json\"",
							"I HAS A VARIABLE JSON_DATA TEH STRIN ITZ \"{\\\"name\\\":\\\"Alice\\\",\\\"email\\\":\\\"alice@example.com\\\"}\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://httpbin.org/post\" AN WIT JSON_DATA",
							"IZ RESPONSE IS_SUCCESS?",
							"    SAYZ WIT \"Data posted successfully\"",
							"    I HAS A VARIABLE RESPONSE_JSON TEH BASKIT ITZ RESPONSE DO TO_JSON",
							"    SAYZ WIT \"Server received our data\"",
							"NOPE",
							"    SAYZ WIT \"POST failed with status: \"",
							"    SAYZ WIT RESPONSE STATUS",
							"KTHX",
							"@example POST form data",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/x-www-form-urlencoded\"",
							"I HAS A VARIABLE FORM_DATA TEH STRIN ITZ \"username=alice&password=secret123\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://httpbin.org/post\" AN WIT FORM_DATA",
							"SAYZ WIT \"Form posted, status: \"",
							"SAYZ WIT RESPONSE STATUS",
							"@example POST with authentication",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"Authorization\" AN WIT \"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...\"",
							"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/json\"",
							"I HAS A VARIABLE PAYLOAD TEH STRIN ITZ \"{\\\"action\\\":\\\"create\\\",\\\"resource\\\":\\\"user\\\"}\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://api.example.com/users\" AN WIT PAYLOAD",
							"IZ RESPONSE IS_SUCCESS?",
							"    SAYZ WIT \"User created successfully\"",
							"NOPE",
							"    SAYZ WIT \"Failed to create user: \"",
							"    SAYZ WIT RESPONSE BODY",
							"KTHX",
							"@example POST large data",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT TIMEOUT ITZ 60", // Increase timeout for large data
							"I HAS A VARIABLE LARGE_DATA TEH STRIN ITZ \"<very large XML or JSON payload>\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://api.example.com/upload\" AN WIT LARGE_DATA",
							"SAYZ WIT \"Upload completed with status: \"",
							"SAYZ WIT RESPONSE STATUS",
							"@throws Exception if URL is invalid or unreachable",
							"@throws Exception if network error occurs",
							"@note URL must include protocol (http:// or https://)",
							"@note Data is sent as the request body",
							"@note Applies any headers set on the client",
							"@note Content-Type header should be set appropriately",
							"@see GET, PUT, DELETE, RESPONSE",
							"@category http-requests",
						},
						ReturnType: "RESPONSE",
						Parameters: []environment.Parameter{
							{Name: "url", Type: "STRIN"},
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return executeHTTPRequest(this, "POST", args[0], args[1])
						},
					},
					// PUT method
					"PUT": {
						Name: "PUT",
						Documentation: []string{
							"Makes an HTTP PUT request with data in the request body to the specified URL.",
							"Returns a RESPONSE object containing status, body, and headers.",
							"",
							"@syntax <client> DO PUT WIT <url> AN WIT <data>",
							"@param {STRIN} url - The URL to request (must include http:// or https://)",
							"@param {STRIN} data - The data to send in the request body",
							"@returns {RESPONSE} Response object with status, body, and headers",
							"@example Update resource with PUT",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/json\"",
							"I HAS A VARIABLE UPDATED_DATA TEH STRIN ITZ \"{\\\"name\\\":\\\"Alice\\\",\\\"email\\\":\\\"alice.smith@example.com\\\"}\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT \"https://api.example.com/users/123\" AN WIT UPDATED_DATA",
							"IZ RESPONSE IS_SUCCESS?",
							"    SAYZ WIT \"User updated successfully\"",
							"NOPE",
							"    SAYZ WIT \"Update failed with status: \"",
							"    SAYZ WIT RESPONSE STATUS",
							"KTHX",
							"@example PUT with version control",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"If-Match\" AN WIT \"\\\"etag123\\\"\"",
							"I HAS A VARIABLE RESOURCE_DATA TEH STRIN ITZ \"{\\\"content\\\":\\\"Updated content\\\",\\\"version\\\":2}\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT \"https://api.example.com/documents/456\" AN WIT RESOURCE_DATA",
							"IZ RESPONSE STATUS SAEM AS 412?",
							"    SAYZ WIT \"Resource was modified by another client\"",
							"NOPE",
							"    IZ RESPONSE IS_SUCCESS?",
							"        SAYZ WIT \"Document updated successfully\"",
							"    NOPE",
							"        SAYZ WIT \"Update failed\"",
							"    KTHX",
							"KTHX",
							"@example Replace entire resource",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/json\"",
							"I HAS A VARIABLE COMPLETE_RESOURCE TEH STRIN ITZ \"{\\\"id\\\":789,\\\"name\\\":\\\"New Name\\\",\\\"active\\\":YEZ,\\\"tags\\\":[\\\"tag1\\\",\\\"tag2\\\"]}\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT \"https://api.example.com/resources/789\" AN WIT COMPLETE_RESOURCE",
							"SAYZ WIT \"Resource replacement status: \"",
							"SAYZ WIT RESPONSE STATUS",
							"@example PUT binary data",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/octet-stream\"",
							"I HAS A VARIABLE BINARY_DATA TEH STRIN ITZ \"<binary data as string>\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO PUT WIT \"https://api.example.com/files/document.pdf\" AN WIT BINARY_DATA",
							"IZ RESPONSE IS_SUCCESS?",
							"    SAYZ WIT \"File uploaded successfully\"",
							"NOPE",
							"    SAYZ WIT \"Upload failed\"",
							"KTHX",
							"@throws Exception if URL is invalid or unreachable",
							"@throws Exception if network error occurs",
							"@note URL must include protocol (http:// or https://)",
							"@note PUT typically replaces the entire resource",
							"@note Use for idempotent operations (multiple calls have same effect)",
							"@note Applies any headers set on the client",
							"@see GET, POST, DELETE, RESPONSE",
							"@category http-requests",
						},
						ReturnType: "RESPONSE",
						Parameters: []environment.Parameter{
							{Name: "url", Type: "STRIN"},
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return executeHTTPRequest(this, "PUT", args[0], args[1])
						},
					},
					// DELETE method
					"DELETE": {
						Name: "DELETE",
						Documentation: []string{
							"Makes an HTTP DELETE request to the specified URL.",
							"Returns a RESPONSE object containing status, body, and headers.",
							"",
							"@syntax <client> DO DELETE WIT <url>",
							"@param {STRIN} url - The URL to request (must include http:// or https://)",
							"@returns {RESPONSE} Response object with status, body, and headers",
							"@example Delete a resource",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT \"https://api.example.com/users/123\"",
							"IZ RESPONSE IS_SUCCESS?",
							"    SAYZ WIT \"User deleted successfully\"",
							"NOPE",
							"    SAYZ WIT \"Delete failed with status: \"",
							"    SAYZ WIT RESPONSE STATUS",
							"KTHX",
							"@example Delete with authentication",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"CLIENT HEADERS DO PUT WIT \"Authorization\" AN WIT \"Bearer token123\"",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT \"https://api.example.com/posts/456\"",
							"IZ RESPONSE STATUS SAEM AS 204?",
							"    SAYZ WIT \"Post deleted (no content returned)\"",
							"NOPE",
							"    IZ RESPONSE IS_SUCCESS?",
							"        SAYZ WIT \"Post deleted successfully\"",
							"    NOPE",
							"        SAYZ WIT \"Delete failed\"",
							"    KTHX",
							"KTHX",
							"@example Check if resource exists before deleting",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE CHECK_RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/items/789\"",
							"IZ CHECK_RESPONSE STATUS SAEM AS 404?",
							"    SAYZ WIT \"Resource does not exist\"",
							"NOPE",
							"    I HAS A VARIABLE DELETE_RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT \"https://api.example.com/items/789\"",
							"    IZ DELETE_RESPONSE IS_SUCCESS?",
							"        SAYZ WIT \"Resource deleted successfully\"",
							"    NOPE",
							"        SAYZ WIT \"Delete failed\"",
							"    KTHX",
							"KTHX",
							"@example Batch delete with error handling",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE IDS TEH BUKKIT ITZ NEW BUKKIT",
							"IDS DO PUSH WIT \"101\"",
							"IDS DO PUSH WIT \"102\"",
							"IDS DO PUSH WIT \"103\"",
							"I HAS A VARIABLE SUCCESS_COUNT TEH INTEGR ITZ 0",
							"I HAS A VARIABLE ERROR_COUNT TEH INTEGR ITZ 0",
							"WHILE NO SAEM AS (IDS LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE ID TEH STRIN ITZ IDS DO POP",
							"    I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO DELETE WIT \"https://api.example.com/resources/\" MOAR ID",
							"    IZ RESPONSE IS_SUCCESS?",
							"        SUCCESS_COUNT ITZ SUCCESS_COUNT UP 1",
							"    NOPE",
							"        ERROR_COUNT ITZ ERROR_COUNT UP 1",
							"        SAYZ WIT \"Failed to delete resource \"",
							"        SAYZ WIT ID",
							"        SAYZ WIT \" (status: \"",
							"        SAYZ WIT RESPONSE STATUS",
							"        SAYZ WIT \")\"",
							"    KTHX",
							"KTHX",
							"SAYZ WIT \"Batch delete completed: \"",
							"SAYZ WIT SUCCESS_COUNT",
							"SAYZ WIT \" successful, \"",
							"SAYZ WIT ERROR_COUNT",
							"SAYZ WIT \" failed\"",
							"@throws Exception if URL is invalid or unreachable",
							"@throws Exception if network error occurs",
							"@note URL must include protocol (http:// or https://)",
							"@note DELETE operations should be idempotent",
							"@note Response body may be empty (204 No Content) on success",
							"@note Applies any headers set on the client",
							"@see GET, POST, PUT, RESPONSE",
							"@category http-requests",
						},
						ReturnType: "RESPONSE",
						Parameters: []environment.Parameter{
							{Name: "url", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return executeHTTPRequest(this, "DELETE", args[0], environment.StringValue(""))
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"TIMEOUT": {
						Variable: environment.Variable{
							Name: "TIMEOUT",
							Documentation: []string{
								"Request timeout in seconds.",
								"Default is 30 seconds. Set to 0 for no timeout.",
								"",
								"@property {INTEGR} TIMEOUT - Timeout duration in seconds",
								"@example Set custom timeout",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT TIMEOUT ITZ 60",
								"SAYZ WIT \"Timeout set to 60 seconds\"",
								"@example Short timeout for quick requests",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT TIMEOUT ITZ 5",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/delay/1\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    SAYZ WIT \"Request completed within 5 seconds\"",
								"NOPE",
								"    SAYZ WIT \"Request timed out or failed\"",
								"KTHX",
								"@example Long timeout for slow operations",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT TIMEOUT ITZ 300", // 5 minutes
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://api.example.com/slow-operation\" AN WIT \"<large data>\"",
								"SAYZ WIT \"Long operation completed with status: \"",
								"SAYZ WIT RESPONSE STATUS",
								"@example Disable timeout",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT TIMEOUT ITZ 0",
								"SAYZ WIT \"Timeout disabled - request will wait indefinitely\"",
								"@example Check current timeout",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"SAYZ WIT \"Current timeout: \"",
								"SAYZ WIT CLIENT TIMEOUT",
								"SAYZ WIT \" seconds\"",
								"@note Default timeout is 30 seconds",
								"@note Setting to 0 disables timeout (not recommended)",
								"@note Timeout applies to entire request (connect + read + write)",
								"@note Network errors may occur before timeout",
								"@see HEADERS, GET, POST",
								"@category http-configuration",
							},
							Type:     "INTEGR",
							IsLocked: false,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if interwebData, ok := this.NativeData.(*InterwebData); ok {
								return environment.IntegerValue(int(interwebData.Client.Timeout.Seconds())), nil
							}
							return environment.IntegerValue(0), runtime.Exception{Message: "TIMEOUT: invalid context"}
						},
						NativeSet: func(this *environment.ObjectInstance, value environment.Value) error {
							timeoutValue, err := value.Cast("INTEGR")
							if err != nil {
								return fmt.Errorf("TIMEOUT expects INTEGR value, got %s", value.Type())
							}
							if interwebData, ok := this.NativeData.(*InterwebData); ok {
								intVal := timeoutValue.(environment.IntegerValue)
								interwebData.Client.Timeout = time.Duration(int(intVal)) * time.Second
								return nil
							}
							return fmt.Errorf("TIMEOUT: invalid context")
						},
					},
					"HEADERS": {
						Variable: environment.Variable{
							Name: "HEADERS",
							Documentation: []string{
								"Request headers as key-value pairs stored in a BASKIT.",
								"Headers are applied to all HTTP requests made by this client.",
								"",
								"@property {BASKIT} HEADERS - HTTP headers as key-value map",
								"@example Set authorization header",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT HEADERS DO PUT WIT \"Authorization\" AN WIT \"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...\"",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/protected\"",
								"SAYZ WIT \"Authenticated request status: \"",
								"SAYZ WIT RESPONSE STATUS",
								"@example Set content type and user agent",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/json\"",
								"CLIENT HEADERS DO PUT WIT \"User-Agent\" AN WIT \"MyApp/1.0\"",
								"CLIENT HEADERS DO PUT WIT \"Accept\" AN WIT \"application/json\"",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/data\"",
								"SAYZ WIT \"Request completed with proper headers\"",
								"@example Add custom headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT HEADERS DO PUT WIT \"X-API-Key\" AN WIT \"secret-api-key-123\"",
								"CLIENT HEADERS DO PUT WIT \"X-Request-ID\" AN WIT \"req-456\"",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://api.example.com/webhook\" AN WIT \"{\\\"event\\\":\\\"test\\\"}\"",
								"SAYZ WIT \"Custom headers sent\"",
								"@example Modify existing headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT HEADERS DO PUT WIT \"Authorization\" AN WIT \"Bearer old-token\"",
								"I HAS A VARIABLE RESPONSE1 TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/data\"",
								"CLIENT HEADERS DO PUT WIT \"Authorization\" AN WIT \"Bearer new-token\"",
								"I HAS A VARIABLE RESPONSE2 TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/data\"",
								"SAYZ WIT \"Used different auth tokens for requests\"",
								"@example Check current headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT HEADERS DO PUT WIT \"Accept\" AN WIT \"application/json\"",
								"I HAS A VARIABLE ACCEPT_HEADER TEH STRIN ITZ CLIENT HEADERS DO GET WIT \"Accept\"",
								"SAYZ WIT \"Accept header: \"",
								"SAYZ WIT ACCEPT_HEADER",
								"@example Clear all headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"CLIENT HEADERS DO PUT WIT \"Authorization\" AN WIT \"Bearer token123\"",
								"CLIENT HEADERS DO PUT WIT \"Content-Type\" AN WIT \"application/json\"",
								"I HAS A VARIABLE NEW_HEADERS TEH BASKIT ITZ NEW BASKIT",
								"CLIENT HEADERS ITZ NEW_HEADERS",
								"SAYZ WIT \"All headers cleared\"",
								"@note Headers are applied to all requests made by this client",
								"@note Common headers: Authorization, Content-Type, Accept, User-Agent",
								"@note Header names are case-insensitive in HTTP",
								"@note Use BASKIT operations (PUT, GET, HAS) to manage headers",
								"@see TIMEOUT, GET, POST",
								"@category http-configuration",
							},
							Type:     "BASKIT",
							Value:    NewBaskitInstance(),
							IsLocked: false,
							IsPublic: true,
						},
						// Since we are returning an object (BASKIT), we cannot
						// use NativeGet/NativeSet.
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"RESPONSE": {
				Name: "RESPONSE",
				Documentation: []string{
					"HTTP response object containing status information, headers, and response body.",
					"Returned by all HTTP request methods (GET, POST, PUT, DELETE).",
					"",
					"@class RESPONSE",
					"@example Basic response handling",
					"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
					"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/get\"",
					"IZ RESPONSE IS_SUCCESS?",
					"    SAYZ WIT \"Request successful!\"",
					"    SAYZ WIT \"Status: \"",
					"    SAYZ WIT RESPONSE STATUS",
					"    SAYZ WIT \"Body: \"",
					"    SAYZ WIT RESPONSE BODY",
					"NOPE",
					"    SAYZ WIT \"Request failed with status: \"",
					"    SAYZ WIT RESPONSE STATUS",
					"KTHX",
					"@example Parse JSON response",
					"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
					"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/json\"",
					"IZ RESPONSE IS_SUCCESS?",
					"    I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON",
					"    SAYZ WIT \"JSON parsed successfully\"",
					"NOPE",
					"    SAYZ WIT \"Failed to get JSON data\"",
					"KTHX",
					"@example Check response headers",
					"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
					"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/headers\"",
					"I HAS A VARIABLE CONTENT_TYPE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"content-type\"",
					"I HAS A VARIABLE SERVER TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"server\"",
					"SAYZ WIT \"Content-Type: \"",
					"SAYZ WIT CONTENT_TYPE",
					"SAYZ WIT \"Server: \"",
					"SAYZ WIT SERVER",
					"@example Handle different HTTP status codes",
					"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
					"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/status/404\"",
					"IZ RESPONSE STATUS SAEM AS 200?",
					"    SAYZ WIT \"OK\"",
					"NOPE",
					"    IZ RESPONSE STATUS SAEM AS 404?",
					"        SAYZ WIT \"Not Found\"",
					"    NOPE",
					"        IZ RESPONSE STATUS SAEM AS 500?",
					"            SAYZ WIT \"Server Error\"",
					"        NOPE",
					"            SAYZ WIT \"Other status: \"",
					"            SAYZ WIT RESPONSE STATUS",
					"        KTHX",
					"    KTHX",
					"KTHX",
					"@note Contains status code, response body, and headers",
					"@note Use IS_SUCCESS and IS_ERROR for easy status checking",
					"@note Body contains raw response data as string",
					"@note Headers are stored as BASKIT for easy access",
					"@see INTERWEB, TO_JSON",
				},
				QualifiedName: "stdlib:HTTP.RESPONSE",
				ModulePath:    "stdlib:HTTP",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:HTTP.RESPONSE"},
				PublicFunctions: map[string]*environment.Function{
					// TO_JSON method
					"TO_JSON": {
						Name: "TO_JSON",
						Documentation: []string{
							"Parses the response body as JSON and returns a BASKIT.",
							"Throws an exception if the response body is not valid JSON.",
							"",
							"@syntax <response> DO TO_JSON",
							"@returns {BASKIT} Parsed JSON data as nested BASKIT structure",
							"@example Parse simple JSON object",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/json\"",
							"IZ RESPONSE IS_SUCCESS?",
							"    I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON",
							"    SAYZ WIT \"JSON parsed successfully\"",
							"NOPE",
							"    SAYZ WIT \"Failed to get JSON response\"",
							"KTHX",
							"@example Access JSON object properties",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/user/123\"",
							"IZ RESPONSE IS_SUCCESS?",
							"    I HAS A VARIABLE USER_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON",
							"    I HAS A VARIABLE USERNAME TEH STRIN ITZ USER_DATA DO GET WIT \"username\"",
							"    I HAS A VARIABLE EMAIL TEH STRIN ITZ USER_DATA DO GET WIT \"email\"",
							"    SAYZ WIT \"User: \"",
							"    SAYZ WIT USERNAME",
							"    SAYZ WIT \" (\"",
							"    SAYZ WIT EMAIL",
							"    SAYZ WIT \")\"",
							"NOPE",
							"    SAYZ WIT \"Failed to get user data\"",
							"KTHX",
							"@example Handle JSON arrays",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/users\"",
							"IZ RESPONSE IS_SUCCESS?",
							"    I HAS A VARIABLE USERS TEH BUKKIT ITZ RESPONSE DO TO_JSON",
							"    SAYZ WIT \"Found \"",
							"    SAYZ WIT USERS LENGTH",
							"    SAYZ WIT \" users\"",
							"    WHILE NO SAEM AS (USERS LENGTH SAEM AS 0)",
							"        I HAS A VARIABLE USER TEH BASKIT ITZ USERS DO POP",
							"        I HAS A VARIABLE NAME TEH STRIN ITZ USER DO GET WIT \"name\"",
							"        SAYZ WIT \"- \"",
							"        SAYZ WIT NAME",
							"    KTHX",
							"NOPE",
							"    SAYZ WIT \"Failed to get users list\"",
							"KTHX",
							"@example Parse nested JSON structures",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/complex-data\"",
							"IZ RESPONSE IS_SUCCESS?",
							"    I HAS A VARIABLE DATA TEH BASKIT ITZ RESPONSE DO TO_JSON",
							"    I HAS A VARIABLE METADATA TEH BASKIT ITZ DATA DO GET WIT \"metadata\"",
							"    I HAS A VARIABLE VERSION TEH STRIN ITZ METADATA DO GET WIT \"version\"",
							"    I HAS A VARIABLE ITEMS TEH BUKKIT ITZ DATA DO GET WIT \"items\"",
							"    SAYZ WIT \"API version: \"",
							"    SAYZ WIT VERSION",
							"    SAYZ WIT \"Number of items: \"",
							"    SAYZ WIT ITEMS LENGTH",
							"NOPE",
							"    SAYZ WIT \"Failed to parse complex JSON\"",
							"KTHX",
							"@example Error handling for invalid JSON",
							"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
							"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/html\"",
							"IZ RESPONSE IS_SUCCESS?",
							"    MAYB",
							"        I HAS A VARIABLE JSON_DATA TEH BASKIT ITZ RESPONSE DO TO_JSON",
							"        SAYZ WIT \"Unexpected JSON parsing success\"",
							"    OOPSIE ERR",
							"        SAYZ WIT \"Expected error parsing HTML as JSON: \"",
							"        SAYZ WIT ERR",
							"    KTHX",
							"NOPE",
							"    SAYZ WIT \"Request failed\"",
							"KTHX",
							"@throws Exception if response body is not valid JSON",
							"@note Converts JSON objects to BASKIT, arrays to BUKKIT",
							"@note Primitive values (strings, numbers, booleans) remain as-is",
							"@note null values become NOTHIN",
							"@note Nested structures are preserved",
							"@see RESPONSE, BODY",
							"@category response-parsing",
						},
						ReturnType: "BASKIT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							bodyVar, exists := this.Variables["BODY"]
							if !exists {
								return environment.NOTHIN, runtime.Exception{Message: "TO_JSON: BODY variable not found"}
							}

							bodyVal, ok := bodyVar.Value.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "TO_JSON: BODY is not a string"}
							}

							var jsonData any
							err := json.Unmarshal([]byte(bodyVal), &jsonData)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("TO_JSON: failed to parse JSON: %v", err)}
							}

							// Convert JSON to BASKIT
							result := jsonToBaskit(jsonData)
							return result, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"STATUS": {
						Variable: environment.Variable{
							Name: "STATUS",
							Documentation: []string{
								"HTTP status code returned by the server.",
								"Examples: 200 (OK), 404 (Not Found), 500 (Server Error).",
								"",
								"@property {INTEGR} STATUS - HTTP status code (100-599)",
								"@example Check specific status codes",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/status/201\"",
								"IZ RESPONSE STATUS SAEM AS 201?",
								"    SAYZ WIT \"Resource created successfully\"",
								"NOPE",
								"    IZ RESPONSE STATUS SAEM AS 404?",
								"        SAYZ WIT \"Resource not found\"",
								"    NOPE",
								"        SAYZ WIT \"Other status: \"",
								"        SAYZ WIT RESPONSE STATUS",
								"    KTHX",
								"KTHX",
								"@example Categorize status codes",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/status/500\"",
								"I HAS A VARIABLE CODE TEH INTEGR ITZ RESPONSE STATUS",
								"IZ CODE BIGGR THAN 499?",
								"    SAYZ WIT \"Server error (5xx)\"",
								"NOPE",
								"    IZ CODE BIGGR THAN 399?",
								"        SAYZ WIT \"Client error (4xx)\"",
								"    NOPE",
								"        IZ CODE BIGGR THAN 299?",
								"            SAYZ WIT \"Redirection (3xx)\"",
								"        NOPE",
								"            IZ CODE BIGGR THAN 199?",
								"                SAYZ WIT \"Success (2xx)\"",
								"            NOPE",
								"                SAYZ WIT \"Informational (1xx)\"",
								"            KTHX",
								"        KTHX",
								"    KTHX",
								"KTHX",
								"@example Handle common HTTP statuses",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://api.example.com/login\" AN WIT \"{\\\"username\\\":\\\"bad\\\",\\\"password\\\":\\\"wrong\\\"}\"",
								"IZ RESPONSE STATUS SAEM AS 401?",
								"    SAYZ WIT \"Authentication failed - check credentials\"",
								"NOPE",
								"    IZ RESPONSE STATUS SAEM AS 403?",
								"        SAYZ WIT \"Access forbidden - insufficient permissions\"",
								"    NOPE",
								"        IZ RESPONSE STATUS SAEM AS 429?",
								"            SAYZ WIT \"Too many requests - rate limited\"",
								"        NOPE",
								"            IZ RESPONSE IS_SUCCESS?",
								"                SAYZ WIT \"Login successful\"",
								"            NOPE",
								"                SAYZ WIT \"Login failed with status: \"",
								"                SAYZ WIT RESPONSE STATUS",
								"            KTHX",
								"        KTHX",
								"    KTHX",
								"KTHX",
								"@example Log status codes for debugging",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE URLS TEH BUKKIT ITZ NEW BUKKIT",
								"URLS DO PUSH WIT \"https://httpbin.org/status/200\"",
								"URLS DO PUSH WIT \"https://httpbin.org/status/404\"",
								"URLS DO PUSH WIT \"https://httpbin.org/status/500\"",
								"WHILE NO SAEM AS (URLS LENGTH SAEM AS 0)",
								"    I HAS A VARIABLE URL TEH STRIN ITZ URLS DO POP",
								"    I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL",
								"    SAYZ WIT \"URL: \"",
								"    SAYZ WIT URL",
								"    SAYZ WIT \" -> Status: \"",
								"    SAYZ WIT RESPONSE STATUS",
								"KTHX",
								"@note 200-299 = Success, 300-399 = Redirection, 400-499 = Client Error, 500-599 = Server Error",
								"@note Common codes: 200 OK, 201 Created, 301 Moved, 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 500 Internal Server Error",
								"@note Use IS_SUCCESS and IS_ERROR for easy status checking",
								"@note Status codes are standardized by HTTP specification",
								"@see IS_SUCCESS, IS_ERROR, BODY",
								"@category response-properties",
							},
							Type:     "INTEGR",
							Value:    environment.IntegerValue(0),
							IsLocked: true,
							IsPublic: true,
						},
					},
					"BODY": {
						Variable: environment.Variable{
							Name: "BODY",
							Documentation: []string{
								"Raw response body as a string.",
								"Contains the full response content from the server.",
								"",
								"@property {STRIN} BODY - Raw response body content",
								"@example Get and display response body",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/json\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    SAYZ WIT \"Response body:\"",
								"    SAYZ WIT RESPONSE BODY",
								"NOPE",
								"    SAYZ WIT \"Request failed with status: \"",
								"    SAYZ WIT RESPONSE STATUS",
								"KTHX",
								"@example Parse JSON response body",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/json\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    I HAS A VARIABLE JSON_DATA TEH OBJECT ITZ RESPONSE DO TO_JSON",
								"    SAYZ WIT \"JSON parsed successfully\"",
								"NOPE",
								"    SAYZ WIT \"Failed to get JSON data\"",
								"KTHX",
								"@example Handle different content types",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/html\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    I HAS A VARIABLE CONTENT_TYPE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"content-type\"",
								"    SAYZ WIT \"Content-Type: \"",
								"    SAYZ WIT CONTENT_TYPE",
								"    SAYZ WIT \"Body length: \"",
								"    SAYZ WIT RESPONSE BODY LENGTH",
								"    SAYZ WIT \"Body preview:\"",
								"    SAYZ WIT RESPONSE BODY SUBSTRIN WIT 0 AN WIT 100",
								"NOPE",
								"    SAYZ WIT \"Failed to fetch content\"",
								"KTHX",
								"@example Save response body to file",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/bytes/1024\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    I HAS A VARIABLE FILE TEH DOCUMENT ITZ NEW DOCUMENT WIT \"response.bin\"",
								"    FILE DO WRITE WIT RESPONSE BODY",
								"    FILE DO CLOSE",
								"    SAYZ WIT \"Response saved to file\"",
								"NOPE",
								"    SAYZ WIT \"Failed to download file\"",
								"KTHX",
								"@example Check body content for specific text",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/get\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    I HAS A VARIABLE HAS_JSON TEH BOOL ITZ RESPONSE BODY HAS \"json\"",
								"    I HAS A VARIABLE HAS_URL TEH BOOL ITZ RESPONSE BODY HAS \"url\"",
								"    SAYZ WIT \"Contains 'json': \"",
								"    SAYZ WIT HAS_JSON",
								"    SAYZ WIT \"Contains 'url': \"",
								"    SAYZ WIT HAS_URL",
								"NOPE",
								"    SAYZ WIT \"Request failed\"",
								"KTHX",
								"@note Body content is always returned as a string, even for binary data",
								"@note For JSON responses, use TO_JSON method to parse into objects",
								"@note Large response bodies may impact memory usage",
								"@note Check Content-Type header to understand the response format",
								"@note Body may be empty for some responses (e.g., 204 No Content)",
								"@see TO_JSON, HEADERS, STATUS",
								"@category response-properties",
							},
							Type:     "STRIN",
							Value:    environment.StringValue(""),
							IsLocked: true,
							IsPublic: true,
						},
					},
					"HEADERS": {
						Variable: environment.Variable{
							Name: "HEADERS",
							Documentation: []string{
								"Response headers as key-value pairs stored in a BASKIT.",
								"Contains all HTTP headers returned by the server.",
								"",
								"@property {BASKIT} HEADERS - HTTP response headers as key-value map",
								"@example Access common response headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/headers\"",
								"I HAS A VARIABLE CONTENT_TYPE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"content-type\"",
								"I HAS A VARIABLE SERVER TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"server\"",
								"I HAS A VARIABLE CONTENT_LENGTH TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"content-length\"",
								"SAYZ WIT \"Content-Type: \"",
								"SAYZ WIT CONTENT_TYPE",
								"SAYZ WIT \"Server: \"",
								"SAYZ WIT SERVER",
								"SAYZ WIT \"Content-Length: \"",
								"SAYZ WIT CONTENT_LENGTH",
								"@example Check for specific headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/json\"",
								"IZ RESPONSE HEADERS DO HAS WIT \"content-type\"?",
								"    I HAS A VARIABLE CT TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"content-type\"",
								"    SAYZ WIT \"Content-Type header found: \"",
								"    SAYZ WIT CT",
								"NOPE",
								"    SAYZ WIT \"Content-Type header not found\"",
								"KTHX",
								"@example Iterate through all response headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/headers\"",
								"I HAS A VARIABLE HEADER_KEYS TEH BUKKIT ITZ RESPONSE HEADERS DO KEYS",
								"SAYZ WIT \"Response headers:\"",
								"WHILE NO SAEM AS (HEADER_KEYS LENGTH SAEM AS 0)",
								"    I HAS A VARIABLE KEY TEH STRIN ITZ HEADER_KEYS DO POP",
								"    I HAS A VARIABLE VALUE TEH STRIN ITZ RESPONSE HEADERS DO GET WIT KEY",
								"    SAYZ WIT KEY",
								"    SAYZ WIT \": \"",
								"    SAYZ WIT VALUE",
								"KTHX",
								"@example Handle CORS headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/data\"",
								"I HAS A VARIABLE ALLOW_ORIGIN TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"access-control-allow-origin\"",
								"I HAS A VARIABLE ALLOW_METHODS TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"access-control-allow-methods\"",
								"IZ ALLOW_ORIGIN SAEM AS \"*\"?",
								"    SAYZ WIT \"CORS allows all origins\"",
								"NOPE",
								"    SAYZ WIT \"CORS origin: \"",
								"    SAYZ WIT ALLOW_ORIGIN",
								"KTHX",
								"SAYZ WIT \"Allowed methods: \"",
								"SAYZ WIT ALLOW_METHODS",
								"@example Check cache headers",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/cache/300\"",
								"I HAS A VARIABLE CACHE_CONTROL TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"cache-control\"",
								"I HAS A VARIABLE ETAG TEH STRIN ITZ RESPONSE HEADERS DO GET WIT \"etag\"",
								"SAYZ WIT \"Cache-Control: \"",
								"SAYZ WIT CACHE_CONTROL",
								"SAYZ WIT \"ETag: \"",
								"SAYZ WIT ETAG",
								"@note Header names are case-insensitive in HTTP",
								"@note Common headers: content-type, content-length, server, date, cache-control",
								"@note Use BASKIT operations (GET, HAS, KEYS) to access headers",
								"@note Headers may not be present for all responses",
								"@note Some headers may have multiple values (comma-separated)",
								"@see BODY, STATUS, TO_JSON",
								"@category response-properties",
							},
							Type:     "BASKIT",
							Value:    NewBaskitInstance(),
							IsLocked: true,
							IsPublic: true,
						},
					},
					"IS_SUCCESS": {
						Variable: environment.Variable{
							Name: "IS_SUCCESS",
							Documentation: []string{
								"YEZ if the HTTP status code indicates success (200-299), NO otherwise.",
								"Provides a convenient way to check if the request was successful.",
								"",
								"@property {BOOL} IS_SUCCESS - True if status code is 200-299",
								"@example Basic success check",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/get\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    SAYZ WIT \"Request successful!\"",
								"NOPE",
								"    SAYZ WIT \"Request failed with status: \"",
								"    SAYZ WIT RESPONSE STATUS",
								"KTHX",
								"@example Handle different success scenarios",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE URLS TEH BUKKIT ITZ NEW BUKKIT",
								"URLS DO PUSH WIT \"https://httpbin.org/status/200\"",
								"URLS DO PUSH WIT \"https://httpbin.org/status/201\"",
								"URLS DO PUSH WIT \"https://httpbin.org/status/204\"",
								"I HAS A VARIABLE SUCCESS_COUNT TEH INTEGR ITZ 0",
								"WHILE NO SAEM AS (URLS LENGTH SAEM AS 0)",
								"    I HAS A VARIABLE URL TEH STRIN ITZ URLS DO POP",
								"    I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL",
								"    IZ RESPONSE IS_SUCCESS?",
								"        SUCCESS_COUNT ITZ SUCCESS_COUNT UP 1",
								"        SAYZ WIT \"Success: \"",
								"        SAYZ WIT RESPONSE STATUS",
								"    NOPE",
								"        SAYZ WIT \"Failed: \"",
								"        SAYZ WIT RESPONSE STATUS",
								"    KTHX",
								"KTHX",
								"SAYZ WIT \"Total successful requests: \"",
								"SAYZ WIT SUCCESS_COUNT",
								"@example Combine with error checking",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://api.example.com/login\" AN WIT \"{\\\"username\\\":\\\"user\\\",\\\"password\\\":\\\"pass\\\"}\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    SAYZ WIT \"Login successful\"",
								"    I HAS A VARIABLE USER_DATA TEH OBJECT ITZ RESPONSE DO TO_JSON",
								"NOPE",
								"    IZ RESPONSE IS_ERROR?",
								"        SAYZ WIT \"Login failed - authentication error\"",
								"    NOPE",
								"        SAYZ WIT \"Login failed - unknown error (status: \"",
								"        SAYZ WIT RESPONSE STATUS",
								"        SAYZ WIT \")\"",
								"    KTHX",
								"KTHX",
								"@example Use in conditional logic",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/user/profile\"",
								"IZ RESPONSE IS_SUCCESS AN RESPONSE STATUS SAEM AS 200?",
								"    SAYZ WIT \"Profile retrieved successfully\"",
								"NOPE",
								"    IZ RESPONSE IS_SUCCESS AN RESPONSE STATUS SAEM AS 201?",
								"        SAYZ WIT \"Profile created successfully\"",
								"    NOPE",
								"        IZ RESPONSE IS_SUCCESS?",
								"            SAYZ WIT \"Profile operation successful (status: \"",
								"            SAYZ WIT RESPONSE STATUS",
								"            SAYZ WIT \")\"",
								"        NOPE",
								"            SAYZ WIT \"Profile operation failed\"",
								"        KTHX",
								"    KTHX",
								"KTHX",
								"@example Batch processing with success tracking",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE IDS TEH BUKKIT ITZ NEW BUKKIT",
								"IDS DO PUSH WIT \"1\"",
								"IDS DO PUSH WIT \"2\"",
								"IDS DO PUSH WIT \"3\"",
								"I HAS A VARIABLE RESULTS TEH BUKKIT ITZ NEW BUKKIT",
								"WHILE NO SAEM AS (IDS LENGTH SAEM AS 0)",
								"    I HAS A VARIABLE ID TEH STRIN ITZ IDS DO POP",
								"    I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/items/\" MOAR ID",
								"    IZ RESPONSE IS_SUCCESS?",
								"        RESULTS DO PUSH WIT \"SUCCESS: \" MOAR ID",
								"    NOPE",
								"        RESULTS DO PUSH WIT \"FAILED: \" MOAR ID MOAR \" (\" MOAR RESPONSE STATUS MOAR \")\"",
								"    KTHX",
								"KTHX",
								"SAYZ WIT \"Batch results:\"",
								"WHILE NO SAEM AS (RESULTS LENGTH SAEM AS 0)",
								"    SAYZ WIT RESULTS DO POP",
								"KTHX",
								"@note Returns YEZ for status codes 200-299",
								"@note Equivalent to checking if STATUS >= 200 AND STATUS < 300",
								"@note Use IS_ERROR for checking error conditions (400+)",
								"@note Success doesn't guarantee the response body contains expected data",
								"@note Some APIs return 200 with error information in the body",
								"@see IS_ERROR, STATUS, BODY",
								"@category response-properties",
							},
							Type:     "BOOL",
							Value:    environment.NO,
							IsLocked: true,
							IsPublic: true,
						},
					},
					"IS_ERROR": {
						Variable: environment.Variable{
							Name: "IS_ERROR",
							Documentation: []string{
								"YEZ if the HTTP status code indicates an error (400+), NO otherwise.",
								"Includes both client errors (4xx) and server errors (5xx).",
								"",
								"@property {BOOL} IS_ERROR - True if status code is 400 or higher",
								"@example Basic error checking",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://httpbin.org/status/404\"",
								"IZ RESPONSE IS_ERROR?",
								"    SAYZ WIT \"Request failed with status: \"",
								"    SAYZ WIT RESPONSE STATUS",
								"NOPE",
								"    SAYZ WIT \"Request successful\"",
								"KTHX",
								"@example Handle different error types",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE URLS TEH BUKKIT ITZ NEW BUKKIT",
								"URLS DO PUSH WIT \"https://httpbin.org/status/400\"",
								"URLS DO PUSH WIT \"https://httpbin.org/status/404\"",
								"URLS DO PUSH WIT \"https://httpbin.org/status/500\"",
								"WHILE NO SAEM AS (URLS LENGTH SAEM AS 0)",
								"    I HAS A VARIABLE URL TEH STRIN ITZ URLS DO POP",
								"    I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT URL",
								"    IZ RESPONSE IS_ERROR?",
								"        I HAS A VARIABLE CODE TEH INTEGR ITZ RESPONSE STATUS",
								"        IZ CODE BIGGR THAN 499?",
								"            SAYZ WIT \"Server error (5xx): \"",
								"            SAYZ WIT CODE",
								"        NOPE",
								"            SAYZ WIT \"Client error (4xx): \"",
								"            SAYZ WIT CODE",
								"        KTHX",
								"    NOPE",
								"        SAYZ WIT \"Success: \"",
								"        SAYZ WIT RESPONSE STATUS",
								"    KTHX",
								"KTHX",
								"@example Combine with success checking",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO POST WIT \"https://api.example.com/login\" AN WIT \"{\\\"username\\\":\\\"bad\\\",\\\"password\\\":\\\"wrong\\\"}\"",
								"IZ RESPONSE IS_SUCCESS?",
								"    SAYZ WIT \"Login successful\"",
								"NOPE",
								"    IZ RESPONSE IS_ERROR?",
								"        I HAS A VARIABLE CODE TEH INTEGR ITZ RESPONSE STATUS",
								"        IZ CODE SAEM AS 401?",
								"            SAYZ WIT \"Authentication failed - invalid credentials\"",
								"        NOPE",
								"            IZ CODE SAEM AS 403?",
								"                SAYZ WIT \"Access forbidden - insufficient permissions\"",
								"            NOPE",
								"                IZ CODE SAEM AS 429?",
								"                    SAYZ WIT \"Too many requests - rate limited\"",
								"                NOPE",
								"                    SAYZ WIT \"Other error: \"",
								"                    SAYZ WIT CODE",
								"                KTHX",
								"            KTHX",
								"        KTHX",
								"    NOPE",
								"        SAYZ WIT \"Unexpected status: \"",
								"        SAYZ WIT RESPONSE STATUS",
								"    KTHX",
								"KTHX",
								"@example Error handling in API calls",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/user/999999\"",
								"IZ RESPONSE IS_ERROR?",
								"    I HAS A VARIABLE ERROR_MSG TEH STRIN",
								"    I HAS A VARIABLE CODE TEH INTEGR ITZ RESPONSE STATUS",
								"    IZ CODE SAEM AS 404?",
								"        ERROR_MSG ITZ \"User not found\"",
								"    NOPE",
								"        IZ CODE SAEM AS 403?",
								"            ERROR_MSG ITZ \"Access denied\"",
								"        NOPE",
								"            IZ CODE BIGGR THAN 499?",
								"                ERROR_MSG ITZ \"Server error - please try again later\"",
								"            NOPE",
								"                ERROR_MSG ITZ \"Request failed\"",
								"            KTHX",
								"        KTHX",
								"    KTHX",
								"    SAYZ WIT \"Error: \"",
								"    SAYZ WIT ERROR_MSG",
								"    SAYZ WIT \" (status: \"",
								"    SAYZ WIT CODE",
								"    SAYZ WIT \")\"",
								"NOPE",
								"    SAYZ WIT \"User data retrieved successfully\"",
								"KTHX",
								"@example Retry logic with error checking",
								"I HAS A VARIABLE CLIENT TEH INTERWEB ITZ NEW INTERWEB",
								"I HAS A VARIABLE MAX_RETRIES TEH INTEGR ITZ 3",
								"I HAS A VARIABLE RETRY_COUNT TEH INTEGR ITZ 0",
								"I HAS A VARIABLE SUCCESS TEH BOOL ITZ NO",
								"WHILE NO SAEM AS (SUCCESS SAEM AS YEZ) AN NO SAEM AS (RETRY_COUNT SAEM AS MAX_RETRIES)",
								"    I HAS A VARIABLE RESPONSE TEH RESPONSE ITZ CLIENT DO GET WIT \"https://api.example.com/unreliable-endpoint\"",
								"    IZ RESPONSE IS_ERROR?",
								"        RETRY_COUNT ITZ RETRY_COUNT UP 1",
								"        IZ RETRY_COUNT SAEM AS MAX_RETRIES?",
								"            SAYZ WIT \"Max retries reached - giving up\"",
								"        NOPE",
								"            SAYZ WIT \"Request failed (attempt \"",
								"            SAYZ WIT RETRY_COUNT",
								"            SAYZ WIT \"/\"",
								"            SAYZ WIT MAX_RETRIES",
								"            SAYZ WIT \") - retrying...\"",
								"        KTHX",
								"    NOPE",
								"        SUCCESS ITZ YEZ",
								"        SAYZ WIT \"Request successful after \"",
								"        SAYZ WIT RETRY_COUNT",
								"        SAYZ WIT \" retries\"",
								"    KTHX",
								"KTHX",
								"@note Returns YEZ for status codes 400 and above",
								"@note Includes both client errors (4xx) and server errors (5xx)",
								"@note Use IS_SUCCESS for checking success conditions (200-299)",
								"@note Error responses may still contain useful information in the body",
								"@note Some APIs return error details in JSON format",
								"@see IS_SUCCESS, STATUS, BODY",
								"@category response-properties",
							},
							Type:     "BOOL",
							Value:    environment.NO,
							IsLocked: true,
							IsPublic: true,
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
		}
	})
	return httpClasses
}

// executeHTTPRequest performs the actual HTTP request
func executeHTTPRequest(this *environment.ObjectInstance, method string, urlArg environment.Value, dataArg environment.Value) (environment.Value, error) {
	// Get URL
	urlVal, ok := urlArg.(environment.StringValue)
	if !ok {
		return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("%s expects STRIN url, got %s", method, urlArg.Type())}
	}

	// Get data (for POST/PUT)
	var dataVal environment.StringValue
	if dataArg != nil {
		var ok bool
		dataVal, ok = dataArg.(environment.StringValue)
		if !ok {
			return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("%s expects STRIN data, got %s", method, dataArg.Type())}
		}
	}

	// Get INTERWEB data
	interwebData, ok := this.NativeData.(*InterwebData)
	if !ok {
		return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("%s: invalid context", method)}
	}

	// Create request
	var req *http.Request
	var err error

	if method == "POST" || method == "PUT" {
		req, err = http.NewRequest(method, string(urlVal), bytes.NewReader([]byte(dataVal)))
	} else {
		req, err = http.NewRequest(method, string(urlVal), nil)
	}

	if err != nil {
		return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("%s: HTTP request creation failed: %v", method, err)}
	}

	// Apply headers from HEADERS baskit
	headersVar, exists := this.Variables["HEADERS"]
	if exists {
		if headersBaskit, ok := headersVar.Value.(*environment.ObjectInstance); ok {
			if baskitMap, ok := headersBaskit.NativeData.(BaskitMap); ok {
				for key, value := range baskitMap {
					if valueStr, ok := value.(environment.StringValue); ok {
						req.Header.Set(key, string(valueStr))
					}
				}
			}
		}
	}

	// Execute request
	resp, err := interwebData.Client.Do(req)
	if err != nil {
		return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("%s: HTTP request failed: %v", method, err)}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("%s: failed to read response body: %v", method, err)}
	}

	// Create RESPONSE object
	responseClass := getHTTPClasses()["RESPONSE"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(responseClass)
	responseInstance := &environment.ObjectInstance{
		Environment: env,
		Class:       responseClass,
		Variables:   make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(responseInstance)

	// Set response data
	responseInstance.Variables["STATUS"].Value = environment.IntegerValue(resp.StatusCode)
	responseInstance.Variables["BODY"].Value = environment.StringValue(string(body))

	// Convert response headers to BASKIT
	headerBaskit := NewBaskitInstance()
	headerMap := headerBaskit.NativeData.(BaskitMap)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headerMap[key] = environment.StringValue(values[0])
		}
	}
	responseInstance.Variables["HEADERS"].Value = headerBaskit

	// Set success/error flags
	isSuccess := resp.StatusCode >= 200 && resp.StatusCode < 300
	isError := resp.StatusCode >= 400

	responseInstance.Variables["IS_SUCCESS"].Value = environment.BoolValue(isSuccess)
	responseInstance.Variables["IS_ERROR"].Value = environment.BoolValue(isError)

	return responseInstance, nil
}

// jsonToBaskit converts a JSON interface{} to a BASKIT
func jsonToBaskit(data any) environment.Value {
	switch v := data.(type) {
	case map[string]any:
		baskit := NewBaskitInstance()
		baskitMap := baskit.NativeData.(BaskitMap)
		for key, value := range v {
			baskitMap[key] = jsonToBaskit(value)
		}
		return baskit
	case []any:
		bukkit := NewBukkitInstance()
		bukkitSlice := bukkit.NativeData.(BukkitSlice)
		for _, value := range v {
			bukkitSlice = append(bukkitSlice, jsonToBaskit(value))
		}
		bukkit.NativeData = bukkitSlice
		return bukkit
	case string:
		return environment.StringValue(v)
	case float64:
		// JSON numbers are float64, convert based on whether it's an integer
		if v == float64(int(v)) {
			return environment.IntegerValue(int(v))
		}
		return environment.DoubleValue(v)
	case bool:
		return environment.BoolValue(v)
	case nil:
		return environment.NOTHIN
	default:
		return environment.StringValue(fmt.Sprintf("%v", v))
	}
}

// RegisterHTTPInEnv registers HTTP classes in the given environment
func RegisterHTTPInEnv(env *environment.Environment, declarations ...string) error {
	httpClasses := getHTTPClasses()

	// If declarations is empty, import all classes
	if len(declarations) == 0 {
		for _, class := range httpClasses {
			env.DefineClass(class)
		}
		return nil
	}

	// Otherwise, import only specified classes
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := httpClasses[declUpper]; exists {
			env.DefineClass(class)
			// If importing INTERWEB, also import RESPONSE (required dependency)
			if declUpper == "INTERWEB" {
				if responseClass, exists := httpClasses["RESPONSE"]; exists {
					env.DefineClass(responseClass)
				}
			}
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown HTTP declaration: %s", decl)}
		}
	}

	return nil
}
