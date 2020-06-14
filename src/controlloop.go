package main

import "logging"
import "sensor"
import "actuator"


/*
 *	Control loops implement the interactions between sensors, curves, 
 * 	actuators
 */
type ControlLoop struct {
	// Constructor arguments
	id 			string
	sourceT 	*sensor.Sensor		
	sourceP 	*sensor.Sensor
	actuator 	*actuator.Actuator
	// Global settings by default, can be overridden
	hysteresis	int
	// Loop data
	hystT 	int	
}

func NewControlLoop(id string, srcT, srcP *sensor.Sensor, a *actuator.Actuator, hysteresis  int) *ControlLoop{
	if srcT.Unit() != sensor.MilliCelsius {
		panic("Invalid argument")
	}
	if srcP.Unit() != sensor.MilliWatt {
		panic("Invalid argument")
	}
	return &ControlLoop {
		id: id, sourceT: srcT, sourceP: srcP, 
		actuator: a, hysteresis: hysteresis,
	}
}

func (L *ControlLoop) Id() string {
	return L.id
}

func (L *ControlLoop) Loop() {
	// Sample Phase
	sampledT := L.sourceT.Sample()
	logging.Info("%s sampled °C: %d\n", L.id, sampledT)
	//sampledP := L.sourceP.Sample()
	// Prediction Phase
	predT := sampledT
	// Actuation Phase
	if predT <= L.hystT && predT >= L.hystT-L.hysteresis {
		logging.Debug("Hysteresis not crossed, skipping actuation\n")
		return 
	}
	L.hystT = sampledT
	L.actuator.Actuate(predT)
}