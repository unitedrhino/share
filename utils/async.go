package utils

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

func Recover(ctx context.Context) {
	if p := recover(); p != nil {
		HandleThrow(ctx, p)
	}
}

var setPanicNotify func(in string)

func SetPanicNotify(f func(string)) {
	setPanicNotify = f
}

func HandleThrow(ctx context.Context, p any) {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	msg := fmt.Sprintf("HandleThrow|func=%s|error=%#v|stack=%s\n", f, p, string(debug.Stack()))
	logx.WithContext(ctx).Error(msg)
	if setPanicNotify != nil {
		setPanicNotify(msg)
	}
	//os.Exit(-1)
}

func Go(ctx context.Context, f func()) {
	go func() {
		defer Recover(ctx)
		f()
	}()
}

var sDIr string

func init() {
	_, file, _, _ := runtime.Caller(0)
	// compatible solution to get gorm source directory with various operating systems
	sDIr = sourceDir(file)
}

func sourceDir(file string) string {
	dir := filepath.Dir(file)
	dir = filepath.Dir(dir)

	s := filepath.Dir(dir)
	base := filepath.Base(s)
	if base != "share" {
		s = dir
	}
	return filepath.ToSlash(s) + "/"
}

// FileWithLineNum return the file name and line number of the current file
func FileWithLineNum() string {
	pcs := [13]uintptr{}
	// the third caller usually from gorm internal
	len := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:len])
	for i := 0; i < len; i++ {
		// second return value is "more", not "ok"
		frame, _ := frames.Next()
		if (!strings.HasPrefix(frame.File, sDIr) ||
			strings.HasSuffix(frame.File, "_test.go")) && !strings.HasPrefix(frame.Function, "gorm.io") && !strings.HasSuffix(frame.File, ".gen.go") {
			return prettyCaller(frame.File, frame.Line)
		}
	}

	return ""
}

func prettyCaller(file string, line int) string {
	idx := strings.LastIndexByte(file, '/')
	if idx < 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}
	for i := 0; i < 4; i++ {
		idx = strings.LastIndexByte(file[:idx], '/')
		if idx < 0 {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	return fmt.Sprintf("%s:%d", file[idx+1:], line)
}
