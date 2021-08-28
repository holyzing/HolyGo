package advanced

import (
	"fmt"
	"runtime"
	"testing"
)

// NOTE 逃逸
func TestMemoryEscape(t *testing.T) {
	fmt.Println(runtime.NumCPU())
}
