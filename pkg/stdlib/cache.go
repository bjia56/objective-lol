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
				Name:          "STASH",
				QualifiedName: "stdlib:CACHE.STASH",
				ModulePath:    "stdlib:CACHE",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:CACHE.STASH"},
				PublicFunctions: map[string]*environment.Function{
					"PUT": {
						Name:       "PUT",
						ReturnType: "NOTHIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
							{Name: "value", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "Not implemented"}
						},
					},
					"GET": {
						Name:       "GET",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "Not implemented"}
						},
					},
					"CONTAINS": {
						Name:       "CONTAINS",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NO, runtime.Exception{Message: "Not implemented"}
						},
					},
					"DELETE": {
						Name:       "DELETE",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NO, runtime.Exception{Message: "Not implemented"}
						},
					},
					"CLEAR": {
						Name:       "CLEAR",
						ReturnType: "NOTHIN",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "Not implemented"}
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
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							return environment.IntegerValue(0), runtime.Exception{Message: "Not implemented"}
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"MEMSTASH": {
				Name:          "MEMSTASH",
				QualifiedName: "stdlib:CACHE.MEMSTASH",
				ModulePath:    "stdlib:CACHE",
				ParentClasses: []string{"stdlib:CACHE.STASH"},
				MRO:           []string{"stdlib:CACHE.MEMSTASH", "stdlib:CACHE.STASH"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"MEMSTASH": {
						Name: "MEMSTASH",
						Parameters: []environment.Parameter{
							{Name: "capacity", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							capacity := args[0]
							capacityVal, ok := capacity.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("MEMSTASH constructor expects INTEGR capacity, got %s", capacity.Type())
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
						Name:       "PUT",
						ReturnType: "NOTHIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
							{Name: "value", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("PUT: invalid context")
							}

							key := args[0]
							value := args[1]

							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("PUT expects STRIN key, got %s", key.Type())
							}

							valueVal, ok := value.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("PUT expects STRIN value, got %s", value.Type())
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
						Name:       "GET",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("GET: invalid context")
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("GET expects STRIN key, got %s", key.Type())
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
						Name:       "CONTAINS",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NO, fmt.Errorf("CONTAINS: invalid context")
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NO, fmt.Errorf("CONTAINS expects STRIN key, got %s", key.Type())
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
						Name:       "DELETE",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NO, fmt.Errorf("DELETE: invalid context")
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NO, fmt.Errorf("DELETE expects STRIN key, got %s", key.Type())
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
						Name:       "CLEAR",
						ReturnType: "NOTHIN",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							memData, ok := this.NativeData.(*MemStashData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CLEAR: invalid context")
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
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if memData, ok := this.NativeData.(*MemStashData); ok {
								memData.mutex.RLock()
								defer memData.mutex.RUnlock()
								return environment.IntegerValue(memData.size), nil
							}
							return environment.IntegerValue(0), fmt.Errorf("invalid context for SIZ")
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"TIMESTASH": {
				Name:          "TIMESTASH",
				QualifiedName: "stdlib:CACHE.TIMESTASH",
				ModulePath:    "stdlib:CACHE",
				ParentClasses: []string{"stdlib:CACHE.STASH"},
				MRO:           []string{"stdlib:CACHE.TIMESTASH", "stdlib:CACHE.STASH"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"TIMESTASH": {
						Name: "TIMESTASH",
						Parameters: []environment.Parameter{
							{Name: "ttl_seconds", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							ttlSeconds := args[0]
							ttlVal, ok := ttlSeconds.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("TIMESTASH constructor expects INTEGR ttl_seconds, got %s", ttlSeconds.Type())
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
						Name:       "PUT",
						ReturnType: "NOTHIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
							{Name: "value", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("PUT: invalid context")
							}

							key := args[0]
							value := args[1]

							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("PUT expects STRIN key, got %s", key.Type())
							}

							valueVal, ok := value.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("PUT expects STRIN value, got %s", value.Type())
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
						Name:       "GET",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("GET: invalid context")
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("GET expects STRIN key, got %s", key.Type())
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
						Name:       "CONTAINS",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NO, fmt.Errorf("CONTAINS: invalid context")
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NO, fmt.Errorf("CONTAINS expects STRIN key, got %s", key.Type())
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
						Name:       "DELETE",
						ReturnType: "BOOL",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NO, fmt.Errorf("DELETE: invalid context")
							}

							key := args[0]
							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NO, fmt.Errorf("DELETE expects STRIN key, got %s", key.Type())
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
						Name:       "CLEAR",
						ReturnType: "NOTHIN",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							timeData, ok := this.NativeData.(*TimeStashData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CLEAR: invalid context")
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
							return environment.IntegerValue(0), fmt.Errorf("invalid context for SIZ")
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
			return fmt.Errorf("unknown CACHE class: %s", decl)
		}
	}

	return nil
}