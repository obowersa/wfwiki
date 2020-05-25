package lua

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

//TODO: Improve testing to validate against structs
func TestParseTable(t *testing.T) {
	testString := lua.LString("Test")
	test2 := lua.LValue(testString)

	res := LuaMachine.ParseTable(&test2, "returnJson")
	if res != "\"Test\"" {
		t.Errorf("Result was incorrect, expected %s got %s", "Test", res)
	}

}
