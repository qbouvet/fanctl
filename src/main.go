package main

import "fmt"
import "time"

import "sensor"
import "sensor/reader"
import "curve"
import "actuator"
import "loop"



const (
	SAMPLEPERIOD 	= 1500	// ms	Main loop period
	HYSTERESIS 		= 2000	// m°C	Skip downward changes smaller than this
)

var (
	sensors		= make(map[string]*sensor.Sensor)
	curves 		= make(map[string]curve.Curve)
	actuators 	= make(map[string]*actuator.Actuator)
	ctrlLoops	= make(map[string]*loop.ControlLoop)
)



func main() {
	initialize()
	//pprint()
	for {
		log("+iteration+")
		for _,l := range ctrlLoops {
			log("+loop iteration+")
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

	globalCurve := curve.ClampedLinear(50000, 30, 70000, 100)

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
		SAMPLEPERIOD, 2, 8,	// Max speed in 10s
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

	cpuL := loop.New("cpuL", cpuT, cpuP, cpuA, HYSTERESIS)
	gpuL := loop.New("gpuL", gpuT, gpuP, gpuA, HYSTERESIS)

	ctrlLoops["cpuL"] = cpuL
	ctrlLoops["gpuL"] = gpuL

	// Initialize everything to 0 to force refresh at first iteration
	// predictedTemperatures  = []float32{0,0}
	// predictedTemperatures  = []float32{0,0}
	// hysteresisTemperatures = []float32{0,0}
	// currentActuations      = []float32{0,0}
	// nextActuations         = []float32{0,0}
}
















/* 
// 	Sample data from sensors
func sample() {
	log("Sample")
	l.Indent()

	sampledTemperatures = []float32{
		sensors[0].Sample(),
		sensors[1].Sample(),
	}

	sampledPowers =	 []float32{
		sensors[2].Sample(),
		sensors[3].Sample(),
	}
	
	l.UnIndent() 
}


// 	Predict temperature for the next loop
func predict() {
	log("Predict")
	l.Indent()

	predictedTemperatures = make([]float32, len(sampledTemperatures))
	copy(predictedTemperatures, sampledTemperatures)

	l.UnIndent() 
}


//	Compute values for the actuators
func compute() {
	log("Compute")
	l.Indent()

	// Selectively update temperatures if HYSTERESIS is crossed
	for i,t := range predictedTemperatures {
		if t<=hysteresisTemperatures[i] && t>=hysteresisTemperatures[i]-HYSTERESIS {
			continue
		}
		nextActuations[i] = fancurves[i].Query(t)
		hysteresisTemperatures[i] = t
	}

	l.UnIndent() 
}


// 	Apply the computed values to the actuators
func apply() {
	log("Apply")
	l.Indent()

	for i,a := range actuators {
		if nextActuations[i] == currentActuations[i] {
			continue
		}
		actuationStepRange  := float32(math.Abs(float64(nextActuations[i]-currentActuations[i])))
		logf("Need to cover an actuation range of %f\n", actuationStepRange)
		actuationStepSign   := int((nextActuations[i]-currentActuations[i])/actuationStepRange)
		logf("Steps sign: %d\n", actuationStepSign)
		actuationStepsCount := actuationStepRange / fan_step
		logf("Requires %f steps of size %f\n", actuationStepsCount, fan_step)
		if actuationStepsCount > steps_per_period {
			actuationStepsCount = steps_per_period
			logf("WARNING: can't cover step range (%f) with given step size (%d), steps per sample period (%d), and sample period time duration (%d)\n", 
				float32(actuationStepSign) * actuationStepRange, fan_step, steps_per_period, SAMPLEPERIOD)
		}
		actuationPeriod   := float32(SAMPLEPERIOD/actuationStepsCount)
		actuationStepSize := actuationStepSign * fan_step
		logf("Will carry out steps of %f at interval of %s\n", actuationStepSize, actuationPeriod)
		actuationHandler  := func (actuator actuator.Actuator, currentActuation float32, actuationStepSize float32, actuationPeriod float32, until time.Time) {
			for time.Now().Before(until) {
				currentActuation := currentActuation + actuationStepSize
				actuator.Set(currentActuation)
				time.Sleep(time.Duration(int64(time.Millisecond)* int64(math.Ceil(float64(actuationPeriod)))))
			}
		}
		logf("Carrying out %f actuation steps of size %d \n", actuationStepsCount, actuationStepSize)
		go actuationHandler(a, currentActuations[i], float32(actuationStepSize), actuationPeriod, 
			time.Now().Add(time.Duration(int64(time.Millisecond)*int64(SAMPLEPERIOD))))
		currentActuations[i] = currentActuations[i] + float32(actuationStepSize)*actuationStepsCount
	}

	l.UnIndent() 
}
 */





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

func log(s ...interface{}) {
	fmt.Println(s...)
}

func logf(f string, a ...interface{}) {
	fmt.Printf(f, a...)
}
