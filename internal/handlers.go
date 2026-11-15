package internal

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"COMMAND": command,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func command(args []Value) Value {
	return Value{typ: BULK, bulk: "OK"}
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: BULK, bulk: "PONG"}
	}

	return Value{typ: BULK, bulk: args[0].bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) < 2 {
		return Value{typ: ERROR, str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	values := args[1:]

	SETsMu.Lock()
	SETs[key] = JoinValues(values, " ")
	SETsMu.Unlock()

	return Value{typ: STRING, str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: ERROR, str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: NULL}
	}

	return Value{typ: BULK, bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: ERROR, str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: STRING, str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: ERROR, str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: NULL}
	}

	return Value{typ: BULK, bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: ERROR, str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	values, ok := HSETs[hash]
	HSETsMu.RUnlock()
	fmt.Println(values)
	if !ok {
		return Value{typ: NULL}
	}

	result := make([]Value, len(values)*2)

	for k, v := range values {
		fmt.Println("KV: ", k, v)
		result = append(result, Value{typ: BULK, bulk: k})
		result = append(result, Value{typ: BULK, bulk: v})
	}

	return Value{typ: ARRAY, array: result}
}
