package main

import "logging"
import "sensor"
import "sensor/reader"
import "actuator"

/*	Device configuration is currently hard-coded here. 
 * 	That is: 
 * 	  - Fan curves information
 *	  - Sensors information (file path, command line, ...)
 *	  - Actuators informations (file path, command line, ...)
 */ 


func loadConfigurationInto(
	map[string]*sensor.Sensor,
	map[string]actuator.Curve,
	map[string]*actuator.Actuator,
	map[string]*ControlLoop,
) {
	logging.Trace("Loading configuration")

	globalCurve := actuator.ClampedLinear(50000, 30, 70000, 100)

	curves["global"] = globalCurve

	cpuT := sensor.New(
		"cpuT", 
		//reader.FromFile("/sys/devices/pci0000:00/0000:00:18.3/hwmon/hwmon0/temp1_input"),
		reader.FromFile("/sys/class/hwmon/hwmon1/temp2_input"),
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
		reader.FromFile("/sys/devices/platform/nct6775.656/hwmon/hwmon3/pwm2"),
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
		reader.FromFile("/sys/devices/platform/nct6775.656/hwmon/hwmon3/pwm1"),
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
		samplep, 2, 8,		    // Max speed in 10s
		"/sys/devices/platform/nct6775.656/hwmon/hwmon3/pwm2", 
		sensors["cpuPWM"],
	)

	gpuA := actuator.New(
		"gpuA",
		curves["global"],
		0, 255, 				// Full range here
		samplep, 2, 6, 	        // Max speed in 10s
		"/sys/devices/platform/nct6775.656/hwmon/hwmon3/pwm3", 
		sensors["gpuPWM"],
	)

	actuators["cpuA"] = cpuA
	actuators["gpuA"] = gpuA

	cpuL := NewControlLoop("cpuL", cpuT, cpuP, cpuA, hysteresis)
	gpuL := NewControlLoop("gpuL", gpuT, gpuP, gpuA, hysteresis)

	ctrlLoops["cpuL"] = cpuL
	ctrlLoops["gpuL"] = gpuL

}