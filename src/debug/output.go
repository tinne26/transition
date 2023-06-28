package debug

import "fmt"
import "os"
import "io"
import "sync"

var pkgOutput io.Writer = os.Stdout
var pkgErrOutput io.Writer = os.Stderr
var outLock sync.Mutex

func SetOutput(out io.Writer) {
	outLock.Lock()
	pkgOutput = out
	outLock.Unlock()
}

func SetErrOutput(out io.Writer) {
	outLock.Lock()
	pkgErrOutput = out
	outLock.Unlock()
}

func Fatal(args ...any) {
	outLock.Lock()
	_, err := fmt.Fprint(pkgOutput, "[FATAL ERROR] ")
	if err == nil {
		_, err = fmt.Fprint(pkgOutput, args...)
	}
	fmt.Fprint(pkgOutput, "\n")
	if err != nil && pkgErrOutput != os.Stderr {
		fmt.Fprintf(os.Stderr, "debug.Fatal() failed: %s\n", err.Error())
		fmt.Fprint(os.Stderr, "[FATAL ERROR] ")
		fmt.Fprint(os.Stderr, args...)
		fmt.Fprint(os.Stderr, "\n")
	}
	outLock.Unlock()
	os.Exit(1)
}

func Print(args ...any) {
	outLock.Lock()
	_, err := fmt.Fprint(pkgOutput, args...)
	if err != nil && pkgOutput != os.Stdout {
		fmt.Fprintf(os.Stderr, "debug.Print() failed: %s\n", err.Error())
		fmt.Print(args...)
	}
	outLock.Unlock()
}

func Printf(format string, args ...any) {
	outLock.Lock()
	_, err := fmt.Fprintf(pkgOutput, format, args...)
	if err != nil && pkgOutput != os.Stdout {
		fmt.Fprintf(os.Stderr, "debug.Printf() failed: %s\n", err.Error())
		fmt.Printf(format, args...)
	}
	outLock.Unlock()
}

func Trace(args ...any) {
	if !debugTrace { return }
	Print("[TRACE] ")
	Print(args...)
}

func Tracef(format string, args ...any) {
	if !debugTrace { return }
	Print("[TRACE] ")
	Printf(format, args...)
}
