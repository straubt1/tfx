package output

import (
	"encoding/json"
	"fmt"

	"github.com/logrusorgru/aurora"
)

// OutputType represents how to return info
type OutputType string

// List of output types
const (
	DefaultOutput OutputType = "default"
	JsonOutput    OutputType = "json"
)

type Message struct {
	Description string
	Value       interface{}
}

// Store OutputType for reference when posting messages
type Output struct {
	OutputType OutputType
	messages   []*Message
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

// Add display information to show progress to terminal users
// This will be printed immediately for DefaultOutput
func (o Output) AddMessageInfo(m ...string) {
	// only output for default
	if o.OutputType == DefaultOutput {
		fmt.Println(m)
	}
}

// // Add json output, should only be called once
// func (o Output) AddJsonString(m string) {
// 	fmt.Println(m)
// }

// Add display information to show progress to terminal users
// This will be printed immediately for DefaultOutput
func (o Output) AddMessageUserProvided(description string, value string) {
	// only output for default
	if o.OutputType != DefaultOutput {
		return
	}

	fmt.Println(description, aurora.Green(value))
}

// Adds a message that will not print immediate.
// Call Close() to print and align
func (o *Output) AddDeferredMessageRead(description string, value interface{}) {
	o.messages = append(o.messages, &Message{description, value})
}

func (o *Output) Close() {
	if o.OutputType == DefaultOutput {
		maxLength := 0
		for _, k := range o.messages {
			if len(k.Description) > maxLength {
				maxLength = len(k.Description)
				// fmt.Println(k)
			}
		}

		for _, k := range o.messages {
			fmt.Println(fmt.Sprintf("%-*s", maxLength+1, aurora.Bold(k.Description+":")), aurora.Blue(k.Value))

		}
	} else {
		// json
		tempMap := make(map[string]interface{}, len(o.messages))
		for _, v := range o.messages {
			tempMap[v.Description] = v.Value
		}

		b, _ := json.Marshal(tempMap)
		fmt.Println(string(b))
	}

	//reset incase this is used again (outputting mid run)
	o.messages = o.messages[:0]
}
