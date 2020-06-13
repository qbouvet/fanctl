package main

import "fmt"
import "time"

import "sensor"
import "sensor/reader"
import "actuator"



const (
	SAMPLEPERIOD 	= 1500	// ms	Main loop period
	HYSTERESIS 		= 2000	// m°C	Skip downward changes smaller than this
)

var (
	sensors		= make(map[string]*sensor.Sensor)
	curves 		= make(map[string]actuator.Curve)
	actuators 	= make(map[string]*actuator.Actuator)
	ctrlLoops	= make(map[string]*ControlLoop)
)



func main() {
	initialize()
	//pprint()
	for {
		log("+iteration+")
		for _,l := range ctrlLoops {
			log("+control loop iteration+")
			l.Loop()
			time.Sleep(time.Millisecond*50)	
			log()
		}
		time.Sleep(time.Millisecond*SAMPLEPERIOD)	
		//pprint()
	}
}



func initialize() {
	log("Init")

	globalCurve := actuator.ClampedLinear(50000, 30, 70000, 100)

	curves["global"] = globalCurve

	cpuT := sensor.New(
		"cpuT", 
		//reader.FromFile("/sys/devices/pci0000:00/0000:00:18.3/hwmon/hwmon0/temp1_input"),
		reader.FromFile("/sys/class/hwmon/hwmon0/device/hwmon/hwmon0/temp1_input"),
		sensor.MilliCelsius,
	)
	
	// 	CPU Should draw ~90 watts under load at stock, ~150W overclocked
	//	https://bit-tech.net/reviews/tech/amd-ryzen-5-1600-review/6/
	// 	For idle, a wild guess would be 10-15W ? 
	// 	So, Draw[W] = 15[W] + (90[W]-15[W])*utilization[%]
	cpuP := sensor.New(
		"cpuP", 
		reader.FromCmd(
			reader.NewCmdLine("s-tui", "-t"),
			reader.NewCmdLine("sed", "s|.*Util: Avg: \\([0-9.]*\\).*|\\1|"),
		), 
		sensor.Watt, 
		75, 100, 15000,
	)

	cpuPWM := sensor.New(
		"cpuPWM", 
		//reader.FromFile("/sys/devices/pci0000:00/0000:00:18.3/hwmon/hwmon0/pwm2"),
		reader.FromFile("/sys/devices/platform/nct6775.656/hwmon/hwmon2/pwm2"),
		sensor.Natural,
	)

	gpuT := sensor.New(
		"gpuT", 
		reader.FromCmd(reader.NewCmdLine("nvidia-smi", "--format=csv,noheader", "--query-gpu=temperature.gpu")),
		sensor.Celsius,
	)

	gpuP := sensor.New(
		"gpuP", 
		reader.FromCmd(reader.NewCmdLine("nvidia-smi", "--format=csv,noheader,nounits", "--query-gpu=power.draw")), 
		sensor.Watt,
	)

	gpuPWM := sensor.New(
		"gpuPWM", 
		//reader.FromFile("/sys/devices/pci0000:00/0000:00:18.3/hwmon/hwmon0/pwm1"),
		reader.FromFile("/sys/devices/platform/nct6775.656/hwmon/hwmon2/pwm1"),
		sensor.Natural,
	)

	sensors["cpuT"]=cpuT
	sensors["cpuP"]=cpuP
	sensors["cpuPWM"]=cpuPWM
	sensors["gpuT"]=gpuT
	sensors["gpuP"]=gpuP
	sensors["gpuPWM"]=gpuPWM

	cpuA := actuator.New(
		"cpuA",
		curves["global"],
		0, 255, 				// Full range here
		SAMPLEPERIOD, 2, 8,		// Max speed in 10s
		"/sys/devices/platform/nct6775.656/hwmon/hwmon2/pwm2", 
		sensors["cpuPWM"],
	)

	gpuA := actuator.New(
		"gpuA",
		curves["global"],
		0, 255, 				// Full range here
		SAMPLEPERIOD, 2, 6, 	// Max speed in 10s
		"/sys/devices/platform/nct6775.656/hwmon/hwmon2/pwm3", 
		sensors["gpuPWM"],
	)

	actuators["cpuA"] = cpuA
	actuators["gpuA"] = gpuA

	cpuL := NewControlLoop("cpuL", cpuT, cpuP, cpuA, HYSTERESIS)
	gpuL := NewControlLoop("gpuL", gpuT, gpuP, gpuA, HYSTERESIS)

	ctrlLoops["cpuL"] = cpuL
	ctrlLoops["gpuL"] = gpuL

}

func log(s ...interface{}) {
	fmt.Println(s...)
}

func logf(f string, a ...interface{}) {
	fmt.Printf(f, a...)
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
	fmt.Printf("\n%18s |", "Hysteresis")
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
