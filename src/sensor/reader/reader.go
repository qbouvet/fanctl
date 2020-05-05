package reader


type Reader interface {
	Read() float64
}



/*
 * 	A decorator to integrate/test Readers with the existing codebase
 */

 type ReaderDecorator struct {
	R Reader
 }

 func (R *ReaderDecorator) Id() string {
	 return "Reader"
 }

 func (R *ReaderDecorator) Sample() float32 {
	 return float32(R.R.Read()) 
 }