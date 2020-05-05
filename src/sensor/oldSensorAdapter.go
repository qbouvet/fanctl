package sensor

/*
 *	Adapater to integrate new sensors into old code
 */

 type OldSensor interface {
	Id() string
	Sample() float32
}

type OldSensorAdapter struct {
	S *Sensor 
}

func (O *OldSensorAdapter) Id() string {
	return O.S.Id()
}

func (O *OldSensorAdapter) Sample() float32 {
	return float32(O.S.Sample())
}