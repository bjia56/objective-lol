package stdlib

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleCacheCategories defines the order that categories should be rendered in documentation
var moduleCacheCategories = []string{
	"cache-interface",
	"lru-cache",
	"ttl-cache",
	"cache-operations",
	"cache-properties",
}

// MemStashData stores the internal state of a MEMSTASH (LRU cache)
type MemStashData struct {
	cache    map[string]*LRUNode
	head     *LRUNode
	tail     *LRUNode
	capacity int
	size     int
	mutex    sync.RWMutex
}

// LRUNode represents a node in the doubly-linked list for LRU cache
type LRUNode struct {
	key   string
	value string
	prev  *LRUNode
	next  *LRUNode
}

// TimeStashData stores the internal state of a TIMESTASH (TTL cache)
type TimeStashData struct {
	cache map[string]*TTLNode
	ttl   time.Duration
	mutex sync.RWMutex
}

// TTLNode represents a cache entry with expiration time
type TTLNode struct {
	value     string
	expiresAt time.Time
}

// Global CACHE class definitions - created once and reused
var cacheClassesOnce = sync.Once{}
var cacheClasses map[string]*environment.Class

func getCacheClasses() map[string]*environment.Class {
	cacheClassesOnce.Do(func() {
		cacheClasses = map[string]*environment.Class{
			"STASH": {
				Name: "STASH",
				Documentation: []string{
					"Abstract base class for all cache implementations in the CACHE module.",
					"Provides a common interface for cache operations like storing, retrieving, and managing cached data.",
					"",
					"@class STASH",
					"@abstract",
					"@example Basic cache usage pattern",
					"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
					"CACHE DO PUT WIT \"key1\" AN WIT \"value1\"",
					"I HAS A VARIABLE VALUE TEH STRIN ITZ CACHE DO GET WIT \"key1\"",
					"SAYZ WIT \"Retrieved: \" MOAR VALUE",
					"@example Cache operations",
					"IZ CACHE DO CONTAINS WIT \"key1\"?",
					"    SAYZ WIT \"Key exists\"",
					"KTHX",
					"CACHE DO DELETE WIT \"key1\"",
					"CACHE DO CLEAR",
					"@note This is an abstract class - use MEMSTASH or TIMESTASH instead",
					"@note All cache implementations inherit from STASH",
					"@note Thread-safe implementations available in subclasses",
					"@see MEMSTASH, TIMESTASH",
					"@category cache-interface",
				},
				QualifiedName: "stdlib:CACHE.STASH",
				ModulePath:    "stdlib:CACHE",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:CACHE.STASH"},
				PublicFunctions: map[string]*environment.Function{
					"PUT": {
						Name: "PUT",
						Documentation: []string{
							"Stores a key-value pair in the cache.",
							"Overwrites existing keys with new values.",
							"",
							"@syntax <cache> DO PUT WIT <key> AN WIT <value>",
							"@param {STRIN} key - The cache key to store",
							"@param {STRIN} value - The value to associate with the key",
							"@returns {NOTHIN}",
							"@example Store a value",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"CACHE DO PUT WIT \"username\" AN WIT \"john_doe\"",
							"@example Update existing value",
							"CACHE DO PUT WIT \"counter\" AN WIT \"1\"",
							"CACHE DO PUT WIT \"counter\" AN WIT \"2\" BTW Overwrites previous value",
							"@example Store configuration",
							"CACHE DO PUT WIT \"db_host\" AN WIT \"localhost\"",
							"CACHE DO PUT WIT \"db_port\" AN WIT \"5432\"",
							"@note Overwrites existing keys without warning",
							"@note Implementation depends on cache type (LRU vs TTL behavior)",
							"@note Thread-safe in concrete implementations",
							"@see GET, CONTAINS",
							"@category cache-operations",
						},
						ReturnType: "NOTHIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
							{Name: "value", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "PUT not implemented"}
						},
					},
					"GET": {
						Name: "GET",
						Documentation: []string{
							"Retrieves a value by key from the cache.",
							"Returns NOTHIN if key not found or expired (in TTL caches).",
							"",
							"@syntax <cache> DO GET WIT <key>",
							"@param {STRIN} key - The cache key to retrieve",
							"@returns {STRIN|NOTHIN} The cached value or NOTHIN if not found",
							"@example Retrieve a value",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"CACHE DO PUT WIT \"user\" AN WIT \"alice\"",
							"I HAS A VARIABLE USER TEH STRIN ITZ CACHE DO GET WIT \"user\"",
							"SAYZ WIT \"User: \" MOAR USER",
							"@example Handle missing keys",
							"I HAS A VARIABLE VALUE TEH STRIN ITZ CACHE DO GET WIT \"nonexistent\"",
							"IZ VALUE SAEM AS NOTHIN?",
							"    SAYZ WIT \"Key not found\"",
							"KTHX",
							"@example Safe retrieval pattern",
							"IZ CACHE DO CONTAINS WIT \"config\"?",
							"    I HAS A VARIABLE CONFIG TEH STRIN ITZ CACHE DO GET WIT \"config\"",
							"    BTW ... use config ...",
							"KTHX",
							"@note Returns NOTHIN for missing keys",
							"@note In TTL caches, also returns NOTHIN for expired entries",
							"@note May affect cache ordering (LRU caches move accessed items to front)",
							"@see PUT, CONTAINS",
							"@category cache-operations",
						},
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "GET not implemented"}
						},
					},
					"CONTAINS": {
						Name: "CONTAINS",
						Documentation: []string{
							"Checks if a key exists in the cache.",
							"Returns YEZ if exists and not expired (in TTL caches), NO otherwise.",
							"",
							"@syntax <cache> DO CONTAINS WIT <key>",
							"@param {STRIN} key - The cache key to check",
							"@returns {BOOL} YEZ if key exists, NO otherwise",
							"@example Check key existence",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"CACHE DO PUT WIT \"session\" AN WIT \"abc123\"",
							"IZ CACHE DO CONTAINS WIT \"session\"?",
							"    SAYZ WIT \"Session exists\"",
							"KTHX",
							"@example Safe operations",
							"IZ CACHE DO CONTAINS WIT \"user_data\"?",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ CACHE DO GET WIT \"user_data\"",
							"    BTW ... process data ...",
							"NOPE",
							"    SAYZ WIT \"No cached user data\"",
							"KTHX",
							"@example Cache statistics",
							"I HAS A VARIABLE TOTAL_KEYS TEH NUMBR ITZ 0",
							"BTW ... iterate through possible keys ...",
							"IZ CACHE DO CONTAINS WIT \"key1\"?",
							"    TOTAL_KEYS ITZ TOTAL_KEYS MOAR 1",
							"KTHX",
							"@note Does not affect cache ordering (unlike GET)",
							"@note In TTL caches, automatically cleans up expired entries",
							"@note Thread-safe in concrete implementations",
							"@see GET, PUT",
							"@category cache-operations",
						},
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NO, runtime.Exception{Message: "CONTAINS not implemented"}
						},
					},
					"DELETE": {
						Name: "DELETE",
						Documentation: []string{
							"Removes a key-value pair from the cache.",
							"Returns YEZ if deleted, NO if key not found.",
							"",
							"@syntax <cache> DO DELETE WIT <key>",
							"@param {STRIN} key - The cache key to remove",
							"@returns {BOOL} YEZ if key was deleted, NO if key not found",
							"@example Remove a cached item",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"CACHE DO PUT WIT \"temp\" AN WIT \"temporary_data\"",
							"I HAS A VARIABLE DELETED TEH BOOL ITZ CACHE DO DELETE WIT \"temp\"",
							"IZ DELETED?",
							"    SAYZ WIT \"Item removed successfully\"",
							"KTHX",
							"@example Cleanup expired sessions",
							"I HAS A VARIABLE SESSION_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 3600",
							"BTW ... store session data ...",
							"SESSION_CACHE DO DELETE WIT \"session_123\" BTW Force logout",
							"@example Batch cleanup",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"BTW ... store multiple items ...",
							"CACHE DO DELETE WIT \"item1\"",
							"CACHE DO DELETE WIT \"item2\"",
							"CACHE DO DELETE WIT \"item3\"",
							"@note Returns NO if key doesn't exist (not an error)",
							"@note In TTL caches, works regardless of expiration status",
							"@note Thread-safe in concrete implementations",
							"@see CLEAR, PUT",
							"@category cache-operations",
						},
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NO, runtime.Exception{Message: "DELETE not implemented"}
						},
					},
					"CLEAR": {
						Name: "CLEAR",
						Documentation: []string{
							"Removes all items from the cache.",
							"Resets the cache to empty state.",
							"",
							"@syntax <cache> DO CLEAR",
							"@returns {NOTHIN}",
							"@example Clear entire cache",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"BTW ... add many items ...",
							"CACHE DO CLEAR BTW Remove all items",
							"SAYZ WIT \"Cache cleared, size: \" MOAR CACHE SIZ",
							"@example Reset between tests",
							"I HAS A VARIABLE TEST_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
							"BTW ... run test that adds items ...",
							"TEST_CACHE DO CLEAR BTW Clean up for next test",
							"@example Emergency cleanup",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000",
							"BTW ... cache grows too large ...",
							"CACHE DO CLEAR BTW Emergency memory cleanup",
							"@note Completely empties the cache",
							"@note Cannot be undone - use with caution",
							"@note Thread-safe in concrete implementations",
							"@note Size becomes 0 after clearing",
							"@see DELETE, SIZ",
							"@category cache-operations",
						},
						ReturnType: "NOTHIN",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "CLEAR not implemented"}
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"SIZ": {
						Variable: environment.Variable{
							Name: "SIZ",
							Type: "INTEGR",
							Documentation: []string{
								"Read-only property that returns the current number of items in the cache.",
								"",
								"@var {INTEGR} SIZ",
								"@readonly",
								"@example Check cache size",
								"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
								"CACHE DO PUT WIT \"key1\" AN WIT \"value1\"",
								"CACHE DO PUT WIT \"key2\" AN WIT \"value2\"",
								"SAYZ WIT \"Cache size: \" MOAR CACHE SIZ",
								"@example Monitor cache growth",
								"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
								"BTW ... add items over time ...",
								"IZ CACHE SIZ BIGGR THAN 50?",
								"    SAYZ WIT \"Cache getting large, consider cleanup\"",
								"KTHX",
								"@example Empty cache check",
								"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
								"IZ CACHE SIZ SAEM AS 0?",
								"    SAYZ WIT \"Cache is empty\"",
								"KTHX",
								"@note For TTL caches, this may trigger cleanup of expired entries",
								"@note Always reflects current state (may change between calls)",
								"@note Thread-safe in concrete implementations",
								"@see CLEAR",
								"@category cache-properties",
							},
							// Default to 0, overridden in subclasses
							Value:    environment.IntegerValue(0),
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
			"MEMSTASH": {
				Name: "MEMSTASH",
				Documentation: []string{
					"An LRU (Least Recently Used) cache with fixed capacity.",
					"When the cache reaches capacity, the least recently used item is evicted automatically.",
					"Thread-safe implementation with proper concurrency support.",
					"",
					"@class MEMSTASH",
					"@inherits STASH",
					"@example Basic LRU cache usage",
					"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
					"CACHE DO PUT WIT \"user1\" AN WIT \"John\"",
					"CACHE DO PUT WIT \"user2\" AN WIT \"Jane\"",
					"I HAS A VARIABLE USER TEH STRIN ITZ CACHE DO GET WIT \"user1\" BTW Moves to front",
					"SAYZ WIT \"User: \" MOAR USER",
					"@example LRU eviction behavior",
					"I HAS A VARIABLE SMALL_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 2",
					"SMALL_CACHE DO PUT WIT \"first\" AN WIT \"value1\"",
					"SMALL_CACHE DO PUT WIT \"second\" AN WIT \"value2\"",
					"SMALL_CACHE DO PUT WIT \"third\" AN WIT \"value3\" BTW Evicts 'first'",
					"I HAS A VARIABLE MISSING TEH STRIN ITZ SMALL_CACHE DO GET WIT \"first\"",
					"IZ MISSING SAEM AS NOTHIN?",
					"    SAYZ WIT \"First item was evicted\"",
					"KTHX",
					"@example Access pattern optimization",
					"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000",
					"BTW Frequently accessed items stay in cache",
					"WHILE YEZ",
					"    CACHE DO GET WIT \"hot_data\" BTW Keeps this item fresh",
					"    CACHE DO PUT WIT \"new_data\" AN WIT \"value\"",
					"    IZ CACHE SIZ BIGGR THAN 900?",
					"        GTFO BTW Cache getting full",
					"    KTHX",
					"KTHX",
					"@note Fixed capacity - oldest accessed items are evicted first",
					"@note Thread-safe for concurrent access",
					"@note GET operations move items to most recently used position",
					"@note PUT operations on existing keys also move them to front",
					"@see TIMESTASH, STASH",
					"@category lru-cache",
				},
				QualifiedName: "stdlib:CACHE.MEMSTASH",
				ModulePath:    "stdlib:CACHE",
				ParentClasses: []string{"stdlib:CACHE.STASH"},
				MRO:           []string{"stdlib:CACHE.MEMSTASH", "stdlib:CACHE.STASH"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"MEMSTASH": {
						Name: "MEMSTASH",
						Documentation: []string{
							"Initializes a MEMSTASH LRU cache with the specified capacity.",
							"Capacity must be positive integer representing maximum number of items.",
							"",
							"@syntax NEW MEMSTASH WIT <capacity>",
							"@param {INTEGR} capacity - Maximum number of items the cache can hold",
							"@returns {MEMSTASH} A new MEMSTASH instance",
							"@throws Exception if capacity is not positive",
							"@example Create small cache",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 10",
							"@example Create large cache",
							"I HAS A VARIABLE BIG_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 10000",
							"@example Handle invalid capacity",
							"BTW This would throw an exception:",
							"BTW I HAS A VARIABLE BAD_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 0",
							"@example Memory-conscious caching",
							"I HAS A VARIABLE MEM_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 50",
							"BTW Small cache for memory-limited environments",
							"@note Capacity cannot be changed after creation",
							"@note Capacity must be positive (> 0)",
							"@note Memory usage scales with capacity",
							"@note Consider available memory when setting capacity",
							"@see TIMESTASH",
							"@category lru-cache",
						},
						Parameters: []environment.Parameter{
							{Name: "capacity", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							capacity := args[0]
							capacityVal, ok := capacity.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("MEMSTASH constructor expects INTEGR capacity, got %s", capacity.Type())}
							}

							if int(capacityVal) <= 0 {
								return environment.NOTHIN, runtime.Exception{Message: "MEMSTASH capacity must be positive"}
							}

							memData := &MemStashData{
								cache:    make(map[string]*LRUNode),
								capacity: int(capacityVal),
								size:     0,
								mutex:    sync.RWMutex{},
							}

							// Initialize doubly-linked list with dummy head/tail
							memData.head = &LRUNode{}
							memData.tail = &LRUNode{}
							memData.head.next = memData.tail
							memData.tail.prev = memData.head

							this.NativeData = memData
							return environment.NOTHIN, nil
						},
					},
					"PUT": {
						Name: "PUT",
						Documentation: []string{
							"Stores a key-value pair in the MEMSTASH. Updates existing keys and moves them to most recently used position.",
							"If at capacity, removes the least recently used item before adding the new one.",
							"",
							"@syntax <memstash> DO PUT WIT <key> AN WIT <value>",
							"@param {STRIN} key - The cache key to store",
							"@param {STRIN} value - The value to associate with the key",
							"@returns {NOTHIN}",
							"@example Basic storage",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"CACHE DO PUT WIT \"config\" AN WIT \"{\\\"host\\\":\\\"localhost\\\"}\"",
							"@example Update with LRU movement",
							"CACHE DO PUT WIT \"user\" AN WIT \"alice\"",
							"CACHE DO PUT WIT \"user\" AN WIT \"bob\" BTW Moves to front, no eviction",
							"@example Capacity management",
							"I HAS A VARIABLE SMALL_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 2",
							"SMALL_CACHE DO PUT WIT \"a\" AN WIT \"1\"",
							"SMALL_CACHE DO PUT WIT \"b\" AN WIT \"2\"",
							"SMALL_CACHE DO PUT WIT \"c\" AN WIT \"3\" BTW Evicts 'a'",
							"IZ SMALL_CACHE DO CONTAINS WIT \"a\"?",
							"    SAYZ WIT \"'a' still exists\"",
							"NOPE",
							"    SAYZ WIT \"'a' was evicted\"",
							"KTHX",
							"@note Updates existing keys without eviction",
							"@note Moves updated items to most recently used position",
							"@note May trigger eviction if at capacity",
							"@note Thread-safe for concurrent operations",
							"@see GET, CONTAINS",
							"@category lru-cache",
						},
						ReturnType: "NOTHIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
							{Name: "value", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "PUT: invalid context"}
							}

							key := args[0]
							value := args[1]

							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("PUT expects STRIN key, got %s", key.Type())}
							}

							valueVal, ok := value.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("PUT expects STRIN value, got %s", value.Type())}
							}

							memData.mutex.Lock()
							defer memData.mutex.Unlock()

							keyStr := string(keyVal)
							valueStr := string(valueVal)

							// If key exists, update and move to front
							if existingNode, exists := memData.cache[keyStr]; exists {
								existingNode.value = valueStr
								memData.moveToFront(existingNode)
								return environment.NOTHIN, nil
							}

							// If at capacity, remove LRU
							if memData.size >= memData.capacity {
								memData.removeLRU()
							}

							// Add new node to front
							newNode := &LRUNode{
								key:   keyStr,
								value: valueStr,
							}
							memData.cache[keyStr] = newNode
							memData.addToFront(newNode)
							memData.size++

							return environment.NOTHIN, nil
						},
					},
					"GET": {
						Name: "GET",
						Documentation: []string{
							"Retrieves a value by key and marks it as recently used.",
							"Moves the accessed item to the front (most recently used position).",
							"Returns NOTHIN if key not found.",
							"",
							"@syntax <memstash> DO GET WIT <key>",
							"@param {STRIN} key - The cache key to retrieve",
							"@returns {STRIN|NOTHIN} The cached value or NOTHIN if not found",
							"@example Access with LRU update",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"CACHE DO PUT WIT \"data\" AN WIT \"value\"",
							"I HAS A VARIABLE VALUE TEH STRIN ITZ CACHE DO GET WIT \"data\" BTW Moves to front",
							"SAYZ WIT \"Got: \" MOAR VALUE",
							"@example LRU behavior demonstration",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 3",
							"CACHE DO PUT WIT \"a\" AN WIT \"1\"",
							"CACHE DO PUT WIT \"b\" AN WIT \"2\"",
							"CACHE DO PUT WIT \"c\" AN WIT \"3\"",
							"CACHE DO GET WIT \"a\" BTW 'a' now most recently used",
							"CACHE DO PUT WIT \"d\" AN WIT \"4\" BTW Evicts 'b' (least recently used)",
							"IZ CACHE DO CONTAINS WIT \"b\"?",
							"    SAYZ WIT \"'b' still exists\"",
							"NOPE",
							"    SAYZ WIT \"'b' was evicted, 'a' was saved by access\"",
							"KTHX",
							"@example Cache warming",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000",
							"BTW Pre-load frequently accessed items",
							"CACHE DO PUT WIT \"hot1\" AN WIT \"value1\"",
							"CACHE DO PUT WIT \"hot2\" AN WIT \"value2\"",
							"BTW ... application runs, accessing items ...",
							"CACHE DO GET WIT \"hot1\" BTW Keep this item fresh",
							"@note Moves accessed items to most recently used position",
							"@note Affects which items get evicted when capacity is reached",
							"@note Thread-safe for concurrent access",
							"@note Use CONTAINS to check existence without affecting LRU order",
							"@see PUT, CONTAINS",
							"@category lru-cache",
						},
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "GET: invalid context"}
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("GET expects STRIN key, got %s", key.Type())}
							}

							memData.mutex.Lock()
							defer memData.mutex.Unlock()

							keyStr := string(keyVal)
							if node, exists := memData.cache[keyStr]; exists {
								memData.moveToFront(node)
								return environment.StringValue(node.value), nil
							}

							return environment.NOTHIN, nil
						},
					},
					"CONTAINS": {
						Name: "CONTAINS",
						Documentation: []string{
							"Checks if a key exists in the MEMSTASH without affecting its position.",
							"Does not move the item in the LRU ordering (unlike GET).",
							"Returns YEZ if exists, NO otherwise.",
							"",
							"@syntax <memstash> DO CONTAINS WIT <key>",
							"@param {STRIN} key - The cache key to check",
							"@returns {BOOL} YEZ if key exists, NO otherwise",
							"@example Check without affecting LRU",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"CACHE DO PUT WIT \"data\" AN WIT \"value\"",
							"IZ CACHE DO CONTAINS WIT \"data\"?",
							"    SAYZ WIT \"Key exists (not moved to front)\"",
							"KTHX",
							"@example Safe access pattern",
							"IZ CACHE DO CONTAINS WIT \"user_session\"?",
							"    I HAS A VARIABLE SESSION TEH STRIN ITZ CACHE DO GET WIT \"user_session\"",
							"    BTW ... use session ...",
							"NOPE",
							"    SAYZ WIT \"Session not cached\"",
							"KTHX",
							"@example Cache statistics",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"BTW ... add various items ...",
							"I HAS A VARIABLE HIT_COUNT TEH NUMBR ITZ 0",
							"I HAS A VARIABLE MISS_COUNT TEH NUMBR ITZ 0",
							"IZ CACHE DO CONTAINS WIT \"key1\"?",
							"    HIT_COUNT ITZ HIT_COUNT MOAR 1",
							"NOPE",
							"    MISS_COUNT ITZ MISS_COUNT MOAR 1",
							"KTHX",
							"@note Does not affect LRU ordering",
							"@note Use GET if you want to access and update LRU position",
							"@note Thread-safe for concurrent operations",
							"@note Faster than GET for existence checks only",
							"@see GET, PUT",
							"@category lru-cache",
						},
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NO, runtime.Exception{Message: "CONTAINS: invalid context"}
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NO, runtime.Exception{Message: fmt.Sprintf("CONTAINS expects STRIN key, got %s", key.Type())}
							}

							memData.mutex.RLock()
							defer memData.mutex.RUnlock()

							keyStr := string(keyVal)
							_, exists := memData.cache[keyStr]
							if exists {
								return environment.YEZ, nil
							}
							return environment.NO, nil
						},
					},
					"DELETE": {
						Name: "DELETE",
						Documentation: []string{
							"Removes a key-value pair from the MEMSTASH.",
							"Returns YEZ if deleted, NO if key not found.",
							"",
							"@syntax <memstash> DO DELETE WIT <key>",
							"@param {STRIN} key - The cache key to remove",
							"@returns {BOOL} YEZ if key was deleted, NO if key not found",
							"@example Remove specific item",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"CACHE DO PUT WIT \"temp\" AN WIT \"temporary_data\"",
							"I HAS A VARIABLE REMOVED TEH BOOL ITZ CACHE DO DELETE WIT \"temp\"",
							"IZ REMOVED?",
							"    SAYZ WIT \"Item removed successfully\"",
							"KTHX",
							"@example Cleanup expired sessions",
							"I HAS A VARIABLE SESSIONS TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000",
							"BTW ... store session data ...",
							"SESSIONS DO DELETE WIT \"expired_session_123\"",
							"@example Selective cache clearing",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"BTW ... cache has various data types ...",
							"CACHE DO DELETE WIT \"debug_info\" BTW Remove debug data",
							"CACHE DO DELETE WIT \"temp_calc\" BTW Remove temporary results",
							"@example Handle deletion result",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"I HAS A VARIABLE WAS_DELETED TEH BOOL ITZ CACHE DO DELETE WIT \"nonexistent\"",
							"IZ WAS_DELETED?",
							"    SAYZ WIT \"Item was deleted\"",
							"NOPE",
							"    SAYZ WIT \"Item didn't exist\"",
							"KTHX",
							"@note Returns NO if key doesn't exist (not an error)",
							"@note Thread-safe for concurrent operations",
							"@note Freed space can be used for new items",
							"@note Does not affect LRU ordering of remaining items",
							"@see CLEAR, PUT",
							"@category lru-cache",
						},
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NO, runtime.Exception{Message: "DELETE: invalid context"}
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NO, runtime.Exception{Message: fmt.Sprintf("DELETE expects STRIN key, got %s", key.Type())}
							}

							memData.mutex.Lock()
							defer memData.mutex.Unlock()

							keyStr := string(keyVal)
							if node, exists := memData.cache[keyStr]; exists {
								memData.removeNode(node)
								delete(memData.cache, keyStr)
								memData.size--
								return environment.YEZ, nil
							}
							return environment.NO, nil
						},
					},
					"CLEAR": {
						Name: "CLEAR",
						Documentation: []string{
							"Removes all items from the MEMSTASH.",
							"Resets the cache to empty state while preserving capacity.",
							"",
							"@syntax <memstash> DO CLEAR",
							"@returns {NOTHIN}",
							"@example Clear entire cache",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"BTW ... add many items ...",
							"SAYZ WIT \"Size before clear: \" MOAR CACHE SIZ",
							"CACHE DO CLEAR",
							"SAYZ WIT \"Size after clear: \" MOAR CACHE SIZ",
							"@example Reset between operations",
							"I HAS A VARIABLE TEMP_CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 50",
							"BTW ... use for temporary calculations ...",
							"TEMP_CACHE DO CLEAR BTW Clean slate for next operation",
							"@example Emergency cleanup",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000",
							"BTW ... cache grew too large or corrupted ...",
							"CACHE DO CLEAR BTW Emergency reset",
							"@example Memory management",
							"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
							"BTW ... long-running application ...",
							"IZ CACHE SIZ BIGGR THAN 80?",
							"    CACHE DO CLEAR BTW Prevent memory bloat",
							"    SAYZ WIT \"Cache cleared to prevent memory issues\"",
							"KTHX",
							"@note Completely empties the cache",
							"@note Preserves the original capacity",
							"@note Cannot be undone - use with caution",
							"@note Thread-safe for concurrent operations",
							"@note Size becomes 0 after clearing",
							"@see DELETE, SIZ",
							"@category lru-cache",
						},
						ReturnType: "NOTHIN",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CLEAR: invalid context"}
							}

							memData.mutex.Lock()
							defer memData.mutex.Unlock()

							memData.cache = make(map[string]*LRUNode)
							memData.head.next = memData.tail
							memData.tail.prev = memData.head
							memData.size = 0

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"SIZ": {
						Variable: environment.Variable{
							Name:     "SIZ",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property that returns the current number of items in the MEMSTASH.",
								"",
								"@var {INTEGR} SIZ",
								"@readonly",
								"@example Monitor cache usage",
								"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
								"CACHE DO PUT WIT \"item1\" AN WIT \"value1\"",
								"CACHE DO PUT WIT \"item2\" AN WIT \"value2\"",
								"SAYZ WIT \"Current size: \" MOAR CACHE SIZ MOAR \"/100\"",
								"@example Capacity management",
								"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 50",
								"WHILE YEZ",
								"    BTW ... add items ...",
								"    IZ CACHE SIZ BIGGR THAN 45?",
								"        GTFO BTW Near capacity",
								"    KTHX",
								"KTHX",
								"@example Empty cache detection",
								"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 100",
								"IZ CACHE SIZ SAEM AS 0?",
								"    SAYZ WIT \"Cache is empty\"",
								"KTHX",
								"@example Cache efficiency metrics",
								"I HAS A VARIABLE CACHE TEH MEMSTASH ITZ NEW MEMSTASH WIT 1000",
								"I HAS A VARIABLE HITS TEH NUMBR ITZ 0",
								"I HAS A VARIABLE MISSES TEH NUMBR ITZ 0",
								"BTW ... process requests ...",
								"IZ CACHE DO CONTAINS WIT \"request_key\"?",
								"    HITS ITZ HITS MOAR 1",
								"    I HAS A VARIABLE DATA TEH STRIN ITZ CACHE DO GET WIT \"request_key\"",
								"NOPE",
								"    MISSES ITZ MISSES MOAR 1",
								"    BTW ... fetch from source ...",
								"KTHX",
								"I HAS A VARIABLE HIT_RATE TEH NUMBR ITZ HITS / (HITS + MISSES)",
								"@note Always reflects current item count",
								"@note Maximum value is the cache capacity",
								"@note Thread-safe for concurrent access",
								"@note Use for monitoring and capacity planning",
								"@see CLEAR",
								"@category cache-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if memData, ok := this.NativeData.(*MemStashData); ok {
								memData.mutex.RLock()
								defer memData.mutex.RUnlock()
								return environment.IntegerValue(memData.size), nil
							}
							return environment.IntegerValue(0), runtime.Exception{Message: "SIZ: invalid context"}
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"TIMESTASH": {
				Name: "TIMESTASH",
				Documentation: []string{
					"A TTL (Time To Live) cache where items expire after a specified duration.",
					"Items are automatically removed when accessed after their expiration time.",
					"Thread-safe implementation with automatic cleanup of expired entries.",
					"",
					"@class TIMESTASH",
					"@inherits STASH",
					"@example Basic TTL cache usage",
					"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
					"CACHE DO PUT WIT \"session\" AN WIT \"user123\"",
					"I HAS A VARIABLE USER TEH STRIN ITZ CACHE DO GET WIT \"session\"",
					"SAYZ WIT \"User: \" MOAR USER",
					"BTW Wait 5+ minutes and session will expire",
					"@example Session management",
					"I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800",
					"SESSIONS DO PUT WIT \"user_abc\" AN WIT \"{\\\"login\\\":\\\"2024-01-01\\\",\\\"role\\\":\\\"admin\\\"}\"",
					"BTW Session automatically expires in 30 minutes",
					"@example Short-lived data caching",
					"I HAS A VARIABLE API_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
					"API_CACHE DO PUT WIT \"weather_nyc\" AN WIT \"{\\\"temp\\\":72,\\\"condition\\\":\\\"sunny\\\"}\"",
					"BTW Data expires in 1 minute",
					"@example Automatic cleanup",
					"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 10",
					"CACHE DO PUT WIT \"temp1\" AN WIT \"value1\"",
					"CACHE DO PUT WIT \"temp2\" AN WIT \"value2\"",
					"BTW ... wait for expiration ...",
					"I HAS A VARIABLE SIZE TEH NUMBR ITZ CACHE SIZ BTW Triggers cleanup",
					"SAYZ WIT \"Active items: \" MOAR SIZE",
					"@note Items expire based on time, not access patterns",
					"@note Expired items are cleaned up automatically on access",
					"@note Thread-safe for concurrent access",
					"@note TTL is set at cache creation and applies to all items",
					"@see MEMSTASH, STASH",
					"@category ttl-cache",
				},
				QualifiedName: "stdlib:CACHE.TIMESTASH",
				ModulePath:    "stdlib:CACHE",
				ParentClasses: []string{"stdlib:CACHE.STASH"},
				MRO:           []string{"stdlib:CACHE.TIMESTASH", "stdlib:CACHE.STASH"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"TIMESTASH": {
						Name: "TIMESTASH",
						Documentation: []string{
							"Initializes a TIMESTASH TTL cache with the specified expiration time.",
							"TTL must be positive integer representing seconds until expiration.",
							"",
							"@syntax NEW TIMESTASH WIT <ttl_seconds>",
							"@param {INTEGR} ttl_seconds - Time-to-live in seconds for all cache items",
							"@returns {TIMESTASH} A new TIMESTASH instance",
							"@throws Exception if ttl_seconds is not positive",
							"@example Create 5-minute cache",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
							"@example Create session cache (30 minutes)",
							"I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800",
							"@example Create short-lived cache (1 minute)",
							"I HAS A VARIABLE TEMP_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
							"@example Handle invalid TTL",
							"BTW This would throw an exception:",
							"BTW I HAS A VARIABLE BAD_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 0",
							"@example API response caching",
							"I HAS A VARIABLE API_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600",
							"BTW Cache API responses for 10 minutes",
							"@note TTL applies to all items in this cache instance",
							"@note TTL cannot be changed after creation",
							"@note TTL is in seconds (not milliseconds)",
							"@note Consider data freshness requirements when setting TTL",
							"@see MEMSTASH",
							"@category ttl-cache",
						},
						Parameters: []environment.Parameter{
							{Name: "ttl_seconds", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							ttlSeconds := args[0]
							ttlVal, ok := ttlSeconds.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("TIMESTASH constructor expects INTEGR ttl_seconds, got %s", ttlSeconds.Type())}
							}

							if int(ttlVal) <= 0 {
								return environment.NOTHIN, runtime.Exception{Message: "TIMESTASH TTL must be positive"}
							}

							timeData := &TimeStashData{
								cache: make(map[string]*TTLNode),
								ttl:   time.Duration(int(ttlVal)) * time.Second,
								mutex: sync.RWMutex{},
							}
							this.NativeData = timeData
							return environment.NOTHIN, nil
						},
					},
					"PUT": {
						Name: "PUT",
						Documentation: []string{
							"Stores a key-value pair in the TIMESTASH with TTL expiration.",
							"Overwrites existing keys with new expiration time.",
							"",
							"@syntax <timestash> DO PUT WIT <key> AN WIT <value>",
							"@param {STRIN} key - The cache key to store",
							"@param {STRIN} value - The value to associate with the key",
							"@returns {NOTHIN}",
							"@example Store with TTL",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
							"CACHE DO PUT WIT \"user_data\" AN WIT \"{\\\"name\\\":\\\"John\\\",\\\"id\\\":123}\"",
							"@example Update existing key",
							"CACHE DO PUT WIT \"config\" AN WIT \"version:1\"",
							"CACHE DO PUT WIT \"config\" AN WIT \"version:2\" BTW Resets expiration",
							"@example Session storage",
							"I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800",
							"SESSIONS DO PUT WIT \"session_abc123\" AN WIT \"user_id:456\"",
							"@example API response caching",
							"I HAS A VARIABLE API_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600",
							"API_CACHE DO PUT WIT \"/api/users\" AN WIT \"[{\\\"id\\\":1,\\\"name\\\":\\\"Alice\\\"}]\"",
							"@example Temporary data storage",
							"I HAS A VARIABLE TEMP TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
							"TEMP DO PUT WIT \"calculation\" AN WIT \"result:42\"",
							"@note Each PUT resets the expiration timer for that key",
							"@note Overwrites existing keys with new TTL",
							"@note Thread-safe for concurrent operations",
							"@note TTL applies from time of PUT, not from last access",
							"@see GET, CONTAINS",
							"@category ttl-cache",
						},
						ReturnType: "NOTHIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
							{Name: "value", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "PUT: invalid context"}
							}

							key := args[0]
							value := args[1]

							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("PUT expects STRIN key, got %s", key.Type())}
							}

							valueVal, ok := value.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("PUT expects STRIN value, got %s", value.Type())}
							}

							timeData.mutex.Lock()
							defer timeData.mutex.Unlock()

							keyStr := string(keyVal)
							valueStr := string(valueVal)
							expiresAt := time.Now().Add(timeData.ttl)

							timeData.cache[keyStr] = &TTLNode{
								value:     valueStr,
								expiresAt: expiresAt,
							}

							return environment.NOTHIN, nil
						},
					},
					"GET": {
						Name: "GET",
						Documentation: []string{
							"Retrieves a non-expired value by key.",
							"Returns NOTHIN if key not found or expired. Automatically removes expired entries.",
							"",
							"@syntax <timestash> DO GET WIT <key>",
							"@param {STRIN} key - The cache key to retrieve",
							"@returns {STRIN|NOTHIN} The cached value or NOTHIN if not found/expired",
							"@example Retrieve fresh data",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
							"CACHE DO PUT WIT \"data\" AN WIT \"fresh_value\"",
							"I HAS A VARIABLE VALUE TEH STRIN ITZ CACHE DO GET WIT \"data\"",
							"IZ VALUE SAEM AS NOTHIN?",
							"    SAYZ WIT \"Data not found or expired\"",
							"NOPE",
							"    SAYZ WIT \"Got: \" MOAR VALUE",
							"KTHX",
							"@example Handle expiration",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 1",
							"CACHE DO PUT WIT \"temp\" AN WIT \"temporary\"",
							"BTW ... wait 2 seconds ...",
							"I HAS A VARIABLE EXPIRED TEH STRIN ITZ CACHE DO GET WIT \"temp\"",
							"IZ EXPIRED SAEM AS NOTHIN?",
							"    SAYZ WIT \"Data expired and was cleaned up\"",
							"KTHX",
							"@example Session validation",
							"I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800",
							"I HAS A VARIABLE SESSION_DATA TEH STRIN ITZ SESSIONS DO GET WIT \"user_session\"",
							"IZ SESSION_DATA SAEM AS NOTHIN?",
							"    SAYZ WIT \"Session expired, please login again\"",
							"KTHX",
							"@example Cache hit/miss tracking",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
							"I HAS A VARIABLE RESULT TEH STRIN ITZ CACHE DO GET WIT \"api_data\"",
							"IZ RESULT SAEM AS NOTHIN?",
							"    SAYZ WIT \"Cache miss - fetching from API\"",
							"    BTW ... fetch from API ...",
							"    CACHE DO PUT WIT \"api_data\" AN WIT \"fetched_data\"",
							"NOPE",
							"    SAYZ WIT \"Cache hit - using cached data\"",
							"KTHX",
							"@note Automatically cleans up expired entries on access",
							"@note Returns NOTHIN for both missing keys and expired entries",
							"@note Thread-safe for concurrent operations",
							"@note Does not extend TTL on access (unlike some cache implementations)",
							"@see PUT, CONTAINS",
							"@category ttl-cache",
						},
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "GET: invalid context"}
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("GET expects STRIN key, got %s", key.Type())}
							}

							timeData.mutex.Lock()
							defer timeData.mutex.Unlock()

							keyStr := string(keyVal)
							if node, exists := timeData.cache[keyStr]; exists {
								// Check if expired
								if time.Now().After(node.expiresAt) {
									delete(timeData.cache, keyStr)
									return environment.NOTHIN, nil
								}
								return environment.StringValue(node.value), nil
							}

							return environment.NOTHIN, nil
						},
					},
					"CONTAINS": {
						Name: "CONTAINS",
						Documentation: []string{
							"Checks if a non-expired key exists in the TIMESTASH.",
							"Returns YEZ if exists and not expired, NO otherwise. Automatically removes expired entries.",
							"",
							"@syntax <timestash> DO CONTAINS WIT <key>",
							"@param {STRIN} key - The cache key to check",
							"@returns {BOOL} YEZ if key exists and not expired, NO otherwise",
							"@example Check existence without retrieval",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
							"CACHE DO PUT WIT \"data\" AN WIT \"value\"",
							"IZ CACHE DO CONTAINS WIT \"data\"?",
							"    SAYZ WIT \"Data exists and is fresh\"",
							"KTHX",
							"@example Safe access pattern",
							"IZ CACHE DO CONTAINS WIT \"user_session\"?",
							"    I HAS A VARIABLE SESSION TEH STRIN ITZ CACHE DO GET WIT \"user_session\"",
							"    BTW ... use session ...",
							"NOPE",
							"    SAYZ WIT \"Session not available\"",
							"KTHX",
							"@example Expiration detection",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 1",
							"CACHE DO PUT WIT \"temp\" AN WIT \"value\"",
							"BTW ... wait 2 seconds ...",
							"IZ CACHE DO CONTAINS WIT \"temp\"?",
							"    SAYZ WIT \"Still fresh\"",
							"NOPE",
							"    SAYZ WIT \"Expired and cleaned up\"",
							"KTHX",
							"@example Cache statistics",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
							"I HAS A VARIABLE FRESH_COUNT TEH NUMBR ITZ 0",
							"BTW ... check multiple keys ...",
							"IZ CACHE DO CONTAINS WIT \"key1\"?",
							"    FRESH_COUNT ITZ FRESH_COUNT MOAR 1",
							"KTHX",
							"SAYZ WIT \"Fresh items: \" MOAR FRESH_COUNT",
							"@note Automatically cleans up expired entries during check",
							"@note Does not extend TTL (unlike GET in some cache systems)",
							"@note Thread-safe for concurrent operations",
							"@note Use GET if you need the actual value",
							"@see GET, PUT",
							"@category ttl-cache",
						},
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NO, runtime.Exception{Message: "CONTAINS: invalid context"}
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NO, runtime.Exception{Message: fmt.Sprintf("CONTAINS expects STRIN key, got %s", key.Type())}
							}

							timeData.mutex.Lock()
							defer timeData.mutex.Unlock()

							keyStr := string(keyVal)
							if node, exists := timeData.cache[keyStr]; exists {
								// Check if expired
								if time.Now().After(node.expiresAt) {
									delete(timeData.cache, keyStr)
									return environment.NO, nil
								}
								return environment.YEZ, nil
							}

							return environment.NO, nil
						},
					},
					"DELETE": {
						Name: "DELETE",
						Documentation: []string{
							"Removes a key-value pair from the TIMESTASH regardless of expiration.",
							"Returns YEZ if deleted, NO if key not found.",
							"",
							"@syntax <timestash> DO DELETE WIT <key>",
							"@param {STRIN} key - The cache key to remove",
							"@returns {BOOL} YEZ if key was deleted, NO if key not found",
							"@example Remove specific item",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
							"CACHE DO PUT WIT \"temp\" AN WIT \"temporary_data\"",
							"I HAS A VARIABLE REMOVED TEH BOOL ITZ CACHE DO DELETE WIT \"temp\"",
							"IZ REMOVED?",
							"    SAYZ WIT \"Item removed successfully\"",
							"KTHX",
							"@example Force logout",
							"I HAS A VARIABLE SESSIONS TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800",
							"SESSIONS DO DELETE WIT \"compromised_session_123\"",
							"@example Cleanup expired data",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
							"BTW ... some items may have expired ...",
							"CACHE DO DELETE WIT \"expired_item\" BTW Remove even if expired",
							"@example Handle deletion result",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
							"I HAS A VARIABLE WAS_DELETED TEH BOOL ITZ CACHE DO DELETE WIT \"nonexistent\"",
							"IZ WAS_DELETED?",
							"    SAYZ WIT \"Item was deleted\"",
							"NOPE",
							"    SAYZ WIT \"Item didn't exist\"",
							"KTHX",
							"@example Selective cache management",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600",
							"BTW ... cache has various data types ...",
							"CACHE DO DELETE WIT \"debug_info\" BTW Remove debug data",
							"CACHE DO DELETE WIT \"temp_calc\" BTW Remove temporary results",
							"@note Works regardless of expiration status",
							"@note Returns NO if key doesn't exist (not an error)",
							"@note Thread-safe for concurrent operations",
							"@note Useful for manual cache management",
							"@see CLEAR, PUT",
							"@category ttl-cache",
						},
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NO, runtime.Exception{Message: "DELETE: invalid context"}
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NO, runtime.Exception{Message: fmt.Sprintf("DELETE expects STRIN key, got %s", key.Type())}
							}

							timeData.mutex.Lock()
							defer timeData.mutex.Unlock()

							keyStr := string(keyVal)
							if _, exists := timeData.cache[keyStr]; exists {
								delete(timeData.cache, keyStr)
								return environment.YEZ, nil
							}
							return environment.NO, nil
						},
					},
					"CLEAR": {
						Name: "CLEAR",
						Documentation: []string{
							"Removes all items from the TIMESTASH regardless of expiration.",
							"Resets the cache to empty state while preserving TTL setting.",
							"",
							"@syntax <timestash> DO CLEAR",
							"@returns {NOTHIN}",
							"@example Clear entire cache",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
							"BTW ... add many items ...",
							"SAYZ WIT \"Size before clear: \" MOAR CACHE SIZ",
							"CACHE DO CLEAR",
							"SAYZ WIT \"Size after clear: \" MOAR CACHE SIZ",
							"@example Reset between operations",
							"I HAS A VARIABLE TEMP_CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
							"BTW ... use for temporary batch processing ...",
							"TEMP_CACHE DO CLEAR BTW Clean slate for next batch",
							"@example Emergency cleanup",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600",
							"BTW ... cache corrupted or needs refresh ...",
							"CACHE DO CLEAR BTW Emergency reset",
							"@example Memory management",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
							"BTW ... long-running application ...",
							"IZ CACHE SIZ BIGGR THAN 1000?",
							"    CACHE DO CLEAR BTW Prevent memory bloat",
							"    SAYZ WIT \"Cache cleared to manage memory\"",
							"KTHX",
							"@example Application restart simulation",
							"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800",
							"BTW ... simulate application restart ...",
							"CACHE DO CLEAR BTW Clear all sessions",
							"SAYZ WIT \"All sessions cleared - fresh start\"",
							"@note Completely empties the cache",
							"@note Preserves the original TTL setting",
							"@note Cannot be undone - use with caution",
							"@note Thread-safe for concurrent operations",
							"@note Size becomes 0 after clearing",
							"@see DELETE, SIZ",
							"@category ttl-cache",
						},
						ReturnType: "NOTHIN",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CLEAR: invalid context"}
							}

							timeData.mutex.Lock()
							defer timeData.mutex.Unlock()

							timeData.cache = make(map[string]*TTLNode)
							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"SIZ": {
						Variable: environment.Variable{
							Name:     "SIZ",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property that returns the current number of non-expired items in the TIMESTASH.",
								"Automatically cleans up expired entries before returning the count.",
								"",
								"@var {INTEGR} SIZ",
								"@readonly",
								"@example Monitor active items",
								"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
								"CACHE DO PUT WIT \"item1\" AN WIT \"value1\"",
								"CACHE DO PUT WIT \"item2\" AN WIT \"value2\"",
								"SAYZ WIT \"Active items: \" MOAR CACHE SIZ",
								"@example Automatic cleanup trigger",
								"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 60",
								"BTW ... items expire over time ...",
								"I HAS A VARIABLE CURRENT_SIZE TEH NUMBR ITZ CACHE SIZ",
								"SAYZ WIT \"Fresh items after cleanup: \" MOAR CURRENT_SIZE",
								"@example Empty cache detection",
								"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 300",
								"IZ CACHE SIZ SAEM AS 0?",
								"    SAYZ WIT \"Cache is empty\"",
								"KTHX",
								"@example Cache efficiency monitoring",
								"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 600",
								"I HAS A VARIABLE INITIAL_SIZE TEH NUMBR ITZ CACHE SIZ",
								"BTW ... application runs and items expire ...",
								"I HAS A VARIABLE FINAL_SIZE TEH NUMBR ITZ CACHE SIZ",
								"I HAS A VARIABLE EXPIRED_COUNT TEH NUMBR ITZ INITIAL_SIZE - FINAL_SIZE",
								"SAYZ WIT \"Items expired: \" MOAR EXPIRED_COUNT",
								"@example Capacity planning",
								"I HAS A VARIABLE CACHE TEH TIMESTASH ITZ NEW TIMESTASH WIT 1800",
								"WHILE YEZ",
								"    BTW ... monitor cache usage ...",
								"    I HAS A VARIABLE USAGE TEH NUMBR ITZ CACHE SIZ",
								"    IZ USAGE BIGGR THAN 1000?",
								"        SAYZ WIT \"Cache usage high: \" MOAR USAGE",
								"    KTHX",
								"    I HAS A VARIABLE DELAY TEH NUMBR ITZ 60",
								"    SLEEPZ WIT DELAY",
								"KTHX",
								"@note Triggers automatic cleanup of expired entries",
								"@note Only counts non-expired items",
								"@note Thread-safe for concurrent access",
								"@note May return different values between calls due to expiration",
								"@note Use for monitoring cache health and planning",
								"@see CLEAR",
								"@category cache-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if timeData, ok := this.NativeData.(*TimeStashData); ok {
								timeData.mutex.Lock()
								defer timeData.mutex.Unlock()

								// Clean expired entries before returning size
								now := time.Now()
								for key, node := range timeData.cache {
									if now.After(node.expiresAt) {
										delete(timeData.cache, key)
									}
								}

								return environment.IntegerValue(len(timeData.cache)), nil
							}
							return environment.IntegerValue(0), runtime.Exception{Message: "SIZ: invalid context"}
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
	return cacheClasses
}

// LRU helper methods for MemStashData
func (m *MemStashData) addToFront(node *LRUNode) {
	node.prev = m.head
	node.next = m.head.next
	m.head.next.prev = node
	m.head.next = node
}

func (m *MemStashData) removeNode(node *LRUNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (m *MemStashData) moveToFront(node *LRUNode) {
	m.removeNode(node)
	m.addToFront(node)
}

func (m *MemStashData) removeLRU() {
	lru := m.tail.prev
	m.removeNode(lru)
	delete(m.cache, lru.key)
	m.size--
}

// RegisterCACHEInEnv registers CACHE classes in the given environment
// declarations: empty slice means import all, otherwise import only specified classes
func RegisterCACHEInEnv(env *environment.Environment, declarations ...string) error {
	cacheClasses := getCacheClasses()

	// If declarations is empty, import all classes
	if len(declarations) == 0 {
		for _, class := range cacheClasses {
			env.DefineClass(class)
		}
		return nil
	}

	// Otherwise, import only specified classes
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := cacheClasses[declUpper]; exists {
			env.DefineClass(class)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown CACHE declaration: %s", decl)}
		}
	}

	return nil
}
