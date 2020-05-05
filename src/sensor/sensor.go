package sensor

//import "fmt"
import "math"

import "sensor/reader"


/*
 *	Sensor interfaces reader, which can return any float value, with the rest 
 *	of the program. In particular, it has knowledge of the unit of the value 
 *	read, and normalizes this value for the rest of the program. 
 *	Namely, The rest of the program expects: 
 * 		- Integer temperatures in MilliCelsius
 *		- Integer powers in MilliWatt
  * 	- Natural count
 */

type Sensor struct {
	id 	 		string			// Metadata
	outputUnit	Unit
	r 	 		reader.Reader	// Reader metadata
	inputUnit  	Unit
	unitMult	int				// Tranformation data
	affineMult 	int
	affineDiv 	int
	affineOfst 	int
}

func New(id string, r reader.Reader, iu Unit, affineCoeffs... int) *Sensor {
	umult, ou := iu.Normalize()
	if len(affineCoeffs) ==0 {
		return &Sensor{id, ou, r, iu, umult, 1, 1, 0}	
	}
	if len(affineCoeffs) ==3 {
		return &Sensor{id, ou, r, iu, umult, affineCoeffs[0], affineCoeffs[1], affineCoeffs[2]}	
	}
	panic("Illegal args count: expecting 0 or 3")
}

func (S *Sensor) Id() string {
	return S.id
}

func (S *Sensor) Unit() Unit {
	return S.outputUnit
}

func (S *Sensor) Sample() int {
	standardized := int(math.Ceil(S.r.Read()*float64(S.unitMult)))
	return standardized * S.affineMult / S.affineDiv + S.affineOfst
}


