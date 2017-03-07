package version

import (
	"fmt"
	"runtime"
)

//VERSION 版本号
const VERSION = "1.0.0"

//String 返回当前版本号描述
func String(app string) string {
	return fmt.Sprintf("%s v%s (built w/%s)", app, VERSION, runtime.Version())
}
