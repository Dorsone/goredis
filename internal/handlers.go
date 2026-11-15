package internal

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"START":   start,
	"COMMAND": command,
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

func start(args []Value) Value {
	return Value{typ: BULK, bulk: "FINISH"}
}
