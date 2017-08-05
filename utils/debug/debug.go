package debug

import (
	"runtime"

	"github.com/NyaaPantsu/nyaa/utils/log"
)

func LogCaller(parent int) {
	if parent <= 0 {
		parent = 1
	}
	parent++
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(parent, pc)
	frames := runtime.CallersFrames(pc)
	for {
		frame, ok := frames.Next()
		if !ok {
			return
		}
		log.Infof("called from %s in %s#%d\n", frame.Func.Name(), frame.File, frame.Line)
	}
}
