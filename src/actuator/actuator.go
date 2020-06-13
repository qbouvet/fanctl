package actuator

import "fmt"
import "math"
import "time"

import "sensor"
import "actuator/writer"


/*		Actuator:
 *		- Determines the appropriate amount of actuation depending on  
 *		  sampled values
 *		- Smoothes out applies the changes over time
 */


type Actuator struct {
	id 				string

	actuationCurve	Curve
	min, max 		int 

	sample_period	int // ms
	step_size		int	// step size
	max_step_rate	int	// steps.s⁻¹

	writer 			*writer.Writer
	feedback 		*sensor.Sensor
}

func New(id string, c Curve, min, max int, 
		 samplePeriod, stepSize, maxStepRate int, 
		 path string, feedback *sensor.Sensor) *Actuator {
	if feedback.Unit() != sensor.Natural {
		panic("Invalid argument")
	}
	A := &Actuator{
		id: id, 
		actuationCurve: c.MapY(0, 100, min, max), min: min, max: max, 
		sample_period: samplePeriod, 
		step_size: stepSize, max_step_rate: maxStepRate, 
		writer: writer.LinuxPwmWriter(path), feedback: feedback,
	}
	return A
} 

func (A *Actuator) Actuate(temperature int) {

	target  	:= A.actuationCurve.Lookup(temperature)
	current 	:= A.feedback.Sample()
	delta 		:= int(math.Abs(float64(target-current)))
	deltaSign	:= int(float64(target-current)/math.Abs(float64(target-current)))
	fmt.Printf("%s going from %d to %d: delta=%d sign=%d\n", A.id, current, target, delta, deltaSign)

	req_nb_steps		:= (delta / A.step_size)				  // #steps
	req_step_rate 		:=1000 * req_nb_steps / A.sample_period  // #steps.s⁻¹
	achieved_step_rate 	:= int(math.Min(
		float64(req_step_rate), float64(A.max_step_rate),
	))
	if achieved_step_rate ==0 {
		achieved_step_rate = 1
	}
	fmt.Printf("%s achieved step rate of %d steps.s⁻¹\n", 
		A.id, achieved_step_rate, 
	)

	achieved_step_period := 1000 / achieved_step_rate		// ms
	until := time.Now().Add(time.Millisecond*time.Duration(A.sample_period))
	fmt.Printf("%a equivalent to 1 step of %d every %d ms\n", 
		A.id, A.step_size, achieved_step_period,
	)

	go func() {
		for time.Now().Before(until) {
			current = current + deltaSign*A.step_size
			A.Write(current)
			time.Sleep(time.Millisecond*time.Duration(achieved_step_period))
		}
		return 
	} () 
}

func (A *Actuator) Write(value int) {
	_value := value
	if value >A.max {
		_value = A.max
	}
	if value <A.min {
		_value = A.min
	}
	A.writer.Write(_value)
}