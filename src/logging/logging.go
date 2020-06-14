package logging

import "log"
import "io/ioutil"
import "io"

/*	Global logging package.
 *	  - Log messages are sortedcategorized, and categories can be toggled
 * 		independantly
 *	  - Log messages can be written file(s) or stdout/err depending on 
 *		initialize() arguments
 *	  - Errors and warnings are written to the errWriter, other logs are 
 *		written to the outWriter.
 */

const prelude = log.Lshortfile

var (
	err  = log.New(ioutil.Discard,  "", 0)	// Cannot be silenced
	warn = log.New(ioutil.Discard,  "", 0)	// Mask = 1
	info = log.New(ioutil.Discard,  "", 0)	// Mask = 2
	trace = log.New(ioutil.Discard, "", 0)	// Mask = 4
	debug = log.New(ioutil.Discard, "", 0)	// Mask = 8
	// Can be extended
)

func Initialize(mask int, errWriter, outWriter io.Writer) {
	err = log.New(errWriter, "E| ", log.Lshortfile)
	if mask %2 == 1 {	
		warn = log.New(errWriter, "W| ", log.Lshortfile)
		mask -= 1
	}
	if mask %4 == 2 {
		info = log.New(outWriter, "I| ", 0)
		mask -= 2
	}
	if mask %8 == 4 {
		trace = log.New(outWriter, "T| ", 0)
		mask -= 4
	}
	if mask %16 == 8 {
		debug = log.New(outWriter, "D| ", log.Lshortfile)
		mask -= 8
	}
}

func Err(fmt string, args ...interface{}) {
	err.Printf(fmt, args...)
}
func Warn(fmt string, args ...interface{}) {
	warn.Printf(fmt, args...)
}
func Info(fmt string, args ...interface{}) {
	info.Printf(fmt, args...)
}
func Trace(fmt string, args ...interface{}) {
	trace.Printf(fmt, args...)
}
func Debug(fmt string, args ...interface{}) {
	debug.Printf(fmt, args...)
}