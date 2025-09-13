package stdlib

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

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
					"Abstract base class for all cache implementations.",
					"Provides common interface for cache operations.",
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
							"Overwrites existing keys.",
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
							"Retrieves a value by key.",
							"Returns NOTHIN if key not found.",
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
							"Returns YEZ if exists, NO otherwise.",
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
					"When the cache reaches capacity, the least recently used item is evicted.",
					"Thread-safe implementation with proper concurrency support.",
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
							"Capacity must be positive.",
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
							"If at capacity, removes the least recently used item.",
						},
						ReturnType: "NOTHIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
							{Name: "value", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("PUT: invalid context")}
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
							"Returns NOTHIN if key not found.",
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
							"Returns YEZ if exists, NO otherwise.",
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
					"Thread-safe implementation with automatic cleanup of expired entries.",
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
							"TTL must be positive (in seconds).",
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
