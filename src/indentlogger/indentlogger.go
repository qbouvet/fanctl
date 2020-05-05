package indentlogger

import "fmt"


type IndentLogger struct {
	indent string
}

func New () *IndentLogger {
	return &IndentLogger{""}
}


func (L *IndentLogger) Log (s ...interface{}) {
	fmt.Printf(L.indent)
	fmt.Println(s...)
}

func (L *IndentLogger) Logf (f string, a ...interface{}) {
	fmt.Printf(L.indent)
	fmt.Printf(f, a...)
}


func (L *IndentLogger) Indent () {
	L.indent = L.indent + "  "
}

func (L *IndentLogger) UnIndent () {
	if len(L.indent) < 2 {
		return
	}
	L.indent = L.indent[0:len(L.indent)-2]
}