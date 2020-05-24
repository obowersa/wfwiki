package lua

import (
	"github.com/yuin/gopher-lua"
)

//TODO: Maybe move this into a new file
const (
	tableToJson = `
	json = { _version = "0.1.2" }

	local encode
	
	local escape_char_map = {
	[ "\\" ] = "\\",
	[ "\"" ] = "\"",
	[ "\b" ] = "b",
	[ "\f" ] = "f",
	[ "\n" ] = "n",
	[ "\r" ] = "r",
	[ "\t" ] = "t",
	}
	
	local escape_char_map_inv = { [ "/" ] = "/" }
	for k, v in pairs(escape_char_map) do
	escape_char_map_inv[v] = k
	end
	
	
	local function escape_char(c)
	return "\\" .. (escape_char_map[c] or string.format("u%04x", c:byte()))
	end
	
	
	local function encode_nil(val)
	return "null"
	end
	
	
	local function encode_table(val, stack)
	local res = {}
	stack = stack or {}
	
	-- Circular reference?
	if stack[val] then error("circular reference") end
	
	stack[val] = true
	
	if rawget(val, 1) ~= nil or next(val) == nil then
	-- Treat as array -- check keys are valid and it is not sparse
	local n = 0
	for k in pairs(val) do
	if type(k) ~= "number" then
	error("invalid table: mixed or invalid key types")
	end
	n = n + 1
	end
	if n ~= #val then
	error("invalid table: sparse array")
	end
	-- Encode
	for i, v in ipairs(val) do
	table.insert(res, encode(v, stack))
	end
	stack[val] = nil
	return "[" .. table.concat(res, ",") .. "]"
	
	else
	-- Treat as an object
	for k, v in pairs(val) do
	if type(k) ~= "string" then
	error("invalid table: mixed or invalid key types")
	end
	table.insert(res, encode(k, stack) .. ":" .. encode(v, stack))
	end
	stack[val] = nil
	return "{" .. table.concat(res, ",") .. "}"
	end
	end
	
	
	local function encode_string(val)
	return '"' .. val:gsub('[%z\1-\31\\"]', escape_char) .. '"'
	end
	
	
	local function encode_number(val)
	-- Check for NaN, -inf and inf
	if val ~= val or val <= -math.huge or val >= math.huge then
	error("unexpected number value '" .. tostring(val) .. "'")
	end
	return string.format("%.14g", val)
	end
	
	
	local type_func_map = {
	[ "nil"     ] = encode_nil,
	[ "table"   ] = encode_table,
	[ "string"  ] = encode_string,
	[ "number"  ] = encode_number,
	[ "boolean" ] = tostring,
	}
	
	
	encode = function(val, stack)
	local t = type(val)
	local f = type_func_map[t]
	if f then
		return f(val, stack)
	end
		error("unexpected type '" .. t .. "'")
	end
	
	
	function returnJson(val)
		return ( encode(val) )
	end
	
	return json
	`
)

// TODO: Implement VM pool to save on memory consumption
type VM struct {
	State *lua.LState
}

var LuaMachine VM

func (v VM) ParseTable(table *lua.LValue, lfunc string) string {

	if err := v.State.CallByParam(lua.P{
		Fn:      v.State.GetGlobal(lfunc),
		NRet:    1,
		Protect: true,
	}, *table); err != nil {
		//TODO: Handle this better
		panic(err)
	}

	//TODO: Error checking before return
	return v.State.Get(-1).String()
}

func (v VM) LoadModule(lfunc string) {
	if err := v.State.DoString(lfunc); err != nil {
		panic(err)
	}
}

//TODO: This is temporary. Need to figure out a better way to handle multiple return types
func (v VM) GetTable() lua.LValue {
	defer v.State.Pop(1)
	return v.State.Get(-1)
}

//TODO: Temporary, implement a goroutine for this with a channel
func init() {
	LuaMachine = VM{
		State: lua.NewState(),
	}

	LuaMachine.LoadModule(tableToJson)
}
