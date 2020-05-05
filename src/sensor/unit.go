package sensor

/*
 *	Helpers for Sensor fields
 */

type Unit int 
const (
 	Celsius 		Unit = iota	// Temperature
 	MilliCelsius
	Watt						// Power
	MilliWatt
	Natural						// Utilization, DutyCycle						
)
 
func (U *Unit) Normalize () (mult int, u Unit) {
	switch *U {
	case Celsius: 		return 1000, MilliCelsius
	case MilliCelsius: 	return 1, MilliCelsius
	case Watt: 			return 1000, MilliWatt
	case MilliWatt:		return 1, MilliWatt
	case Natural: 		return 1, Natural
	default: 			panic("Unknown unit: "+string(u))
	}
}