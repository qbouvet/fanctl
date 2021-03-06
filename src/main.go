package main

import "time"
import "flag"
import "os"
import "os/exec"
import "strconv"

import "logging"
import "sensor"
import "actuator"


var (
	// Parameters
	loglevel    = 15
	samplep     = 1500
	hysteresis  = 2000
	// Runtime data structures
	sensors		= make(map[string]*sensor.Sensor)
	curves 		= make(map[string]actuator.Curve)
	actuators 	= make(map[string]*actuator.Actuator)
	ctrlLoops	= make(map[string]*ControlLoop)
)


func main() {

	flag.IntVar(&loglevel, 	 "loglevel", 3, "Bitmask - Warn:1, Info:2, Trace:4, Debug:8")
	flag.IntVar(&samplep,	 "samplep",  1500, "Sample period in ms")
	flag.IntVar(&hysteresis, "hyst", 	 2000, "hysteresis in milliCelsius")
	flag.Parse()

	logging.Initialize(loglevel, os.Stdout, os.Stdout)

	output, err := exec.Command("id", "-u").Output()
	if err!=nil {
		logging.Err("Failed to run $ id -u")
		os.Exit(1)
	}
	id, err := strconv.Atoi(string(output[:len(output)-1]))
	if err!=nil {
		logging.Err("Failed to run $ id -u")
		os.Exit(1)
	}
	if id!=0 { 
		logging.Err("fanctl must be run as root")
		os.Exit(1)
	}

	logging.Info("Running with parameters:\n")
	logging.Info("loglevel    %d", loglevel)
	logging.Info("samplep     %d", samplep)
	logging.Info("hysteresis  %d", hysteresis)

	loadConfigurationInto(sensors, curves, actuators, ctrlLoops)
	
	for {
		logging.Trace("Iterations += 1")
		for _,l := range ctrlLoops {
			logging.Trace("Iterating control loops")
			l.Loop()
			time.Sleep(time.Millisecond*20)	
		}
		logging.Trace("")
		time.Sleep(time.Millisecond*time.Duration(samplep))
	}
}

func pprintState () {
	logging.Warn("Not implemented")
}

/* func pprint() {
	fmt.Printf("\n")
	fmt.Printf("%18s |", "Sensors")
	for _,s := range sensors {
		fmt.Printf(" %8s |", s.Id())
	}
	fmt.Printf("\n%18s |", "Sampled")
	for _,t := range sampledTemperatures {
		fmt.Printf(" %8.3f |", t)
	}
	for _,p := range sampledPowers {
		fmt.Printf(" %8.3f |", p)
	}
	fmt.Printf("\n%18s |", "Predicted")
	for _,t := range predictedTemperatures {
		fmt.Printf(" %8.3f |", t)
	}
	fmt.Printf("          |          |")
	fmt.Printf("\n%18s |", "hysteresis")
	for _,t := range hysteresisTemperatures {
		fmt.Printf(" %8.3f |", t)
	}
	fmt.Printf("          |          |")
	fmt.Printf("\n%18s |", "Actuation (last)")
	for _,a := range currentActuations {
		fmt.Printf(" %8.3f |", a)
	}
	fmt.Printf("          |          |")
	fmt.Printf("\n%18s |", "Actuation (next)")
	for _,a := range nextActuations {
		fmt.Printf(" %8.3f |", a)
	}
	fmt.Printf("\n\n")
} */
