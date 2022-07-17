package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
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
	OutputType   OutputType
	messages     []*Message
	tableHeaders []interface{}
	tableRows    [][]interface{}
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
func (o Output) AddMessageUserProvided(description string, value string) {
	// only output for default
	if o.OutputType != DefaultOutput {
		return
	}

	fmt.Println(description, aurora.Green(value))
}

// Add display information to show progress to terminal users
// This will be printed immediately for DefaultOutput
func (o Output) AddMessageCalculated(description string, value string) {
	// only output for default
	if o.OutputType != DefaultOutput {
		return
	}

	fmt.Println(description, aurora.Yellow(value))
}

// Adds a message that will not print immediate.
// Call Close() to print and align
func (o *Output) AddDeferredMessageRead(description string, value interface{}) {
	o.messages = append(o.messages, &Message{description, value})
}

// create a row from an array of interfaces, required since table.Row{} uses a variadic
// should work for any type
func createRow(items []interface{}) table.Row {
	h := make([]interface{}, len(items))
	copy(h, items)
	return h
}

// print
func (o Output) closeMessagesDefault() {
	maxLength := 0
	for _, k := range o.messages {
		if len(k.Description) > maxLength {
			maxLength = len(k.Description)
		}
	}

	for _, k := range o.messages {
		fmt.Println(fmt.Sprintf("%-*s", maxLength+1, aurora.Bold(k.Description+":")), aurora.Blue(k.Value))
	}
}

func (o Output) closeMessagesJson() {
	tempMap := make(map[string]interface{}, len(o.messages))
	for _, v := range o.messages {
		tempMap[v.Description] = v.Value
	}

	b, _ := json.Marshal(tempMap)
	fmt.Println(string(b))
}

func (o Output) closeTableDefault() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(createRow(o.tableHeaders))
	for _, i := range o.tableRows {
		t.AppendRow(createRow(i))
	}
	t.SetStyle(table.StyleRounded)
	t.Render()
}

// Craziness... Create an array of maps that can then be Marshalled to JSON
func (o Output) closeTableJson() {
	tempList := make([]map[string]interface{}, len(o.tableRows))
	for index1, v1 := range o.tableRows {
		tempMap := make(map[string]interface{}, len(o.tableHeaders))
		for index2, v2 := range o.tableHeaders {
			if index2 < len(v1) { // quick check for the case where there are more header values than elements in the row
				tempMap[v2.(string)] = v1[index2]
			}
		}
		tempList[index1] = tempMap
	}

	b, _ := json.Marshal(tempList)
	fmt.Println(string(b))
}

func (o *Output) Close() {
	if len(o.messages) > 0 {
		if o.OutputType == DefaultOutput {
			o.closeMessagesDefault()
		} else {
			o.closeMessagesJson()
		}
		o.messages = o.messages[:0]
	}

	if len(o.tableHeaders) > 0 {
		if o.OutputType == DefaultOutput {
			o.closeTableDefault()
		} else {
			o.closeTableJson()
		}
		o.tableHeaders = o.tableHeaders[:0]
		o.tableRows = o.tableRows[:0]
	}
}

// Add headers for table, used for a list of items
func (o *Output) AddTableHeader(headers ...interface{}) {
	o.tableHeaders = headers
}

// Add rows for at table, will be matched to headers set in AddTableHeaders, used for a list of items
func (o *Output) AddTableRows(rows ...interface{}) {
	o.tableRows = append(o.tableRows, rows)
}
