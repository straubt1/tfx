package output

import (
	"fmt"
)

// OutputType represents how to return info
type OutputType string

// List of output types
const (
	DefaultOutput OutputType = "default"
	JsonOutput    OutputType = "json"
)

// Store OutputType for reference when posting messages
type Output struct {
	OutputType OutputType
}

func New(outputType string) Output {
	ot := DefaultOutput
	if outputType == string(JsonOutput) {
		ot = JsonOutput
	}
	o := Output{
		OutputType: ot,
	}
	return o
}

// func New(outputType OutputType) Output {
// 	o := Output{
// 		OutputType: outputType,
// 	}
// 	return o
// }

// Add display information to show progress to terminal users
func (o Output) AddMessageInfo(m ...string) {
	// only output for default
	if o.OutputType == DefaultOutput {
		fmt.Println(m)
	}
}

// Add json output, should only be called once
func (o Output) AddJsonString(m string) {
	fmt.Println(m)
}
