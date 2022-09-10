package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/logrusorgru/aurora"
)

type Message struct {
	Description  string
	Value        interface{}
	ValueList    []interface{}
	ValueComplex [][]interface{}
}

// Store OutputType for reference when posting messages
type Output struct {
	// OutputType   OutputType
	jsonOutput   bool
	messages     []*Message
	tableHeaders []interface{}
	tableRows    [][]interface{}
}

func New(jsonOutput bool) *Output {
	o := &Output{
		jsonOutput: jsonOutput,
	}

	return o
}

// Add display information to show progress to terminal users
// This will be printed immediately for DefaultOutput
func (o *Output) AddMessageUserProvided(description string, value interface{}) {
	// only output for default
	if o.jsonOutput {
		return
	}

	fmt.Println(description, aurora.Green(value))
}

// Add FORMATTED display information to show progress to terminal users
// This will be printed immediately for DefaultOutput
func (o *Output) AddFormattedMessageUserProvided(description string, value interface{}) {
	// only output for default
	if o.jsonOutput {
		return
	}

	fmt.Printf(description+"\n", aurora.Yellow(value))
}

// Add display information to show progress to terminal users
// This will be printed immediately for DefaultOutput
func (o Output) AddMessageCalculated(description string, value interface{}) {
	// only output for default
	if o.jsonOutput {
		return
	}

	fmt.Println(description, aurora.Yellow(value))
}

// Add FORMATTED display information to show progress to terminal users
// This will be printed immediately for DefaultOutput
func (o Output) AddFormattedMessageCalculated(description string, value interface{}) {
	// only output for default
	if o.jsonOutput {
		return
	}

	fmt.Printf(description+"\n", aurora.Yellow(value))
}

// Adds a message that will not print immediate.
// Single primitive Value
// Call Close() to print and align
func (o *Output) AddDeferredMessageRead(description string, value interface{}) {
	o.messages = append(o.messages, &Message{description, value, nil, nil})
}

// Adds a message that will not print immediate.
// List primitive Value
// Call Close() to print and align
func (o *Output) AddDeferredListMessageRead(description string, value []interface{}) {
	o.messages = append(o.messages, &Message{description, nil, value, nil})
}

// Adds a message that will not print immediate.
// List complex Value
// Call Close() to print and align
func (o *Output) AddDeferredListComplexMessageRead(description string, value [][]interface{}) {
	o.messages = append(o.messages, &Message{description, nil, nil, value})
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
	// determine spacing based on largest message description (left justify)
	for _, k := range o.messages {
		if len(k.Description) > maxLength {
			maxLength = len(k.Description)
		}
	}

	for _, k := range o.messages {
		if k.Value != nil {
			fmt.Println(fmt.Sprintf("%-*s", maxLength+1, aurora.Bold(k.Description+":")), aurora.Blue(k.Value))
		}
		if k.ValueList != nil {
			fmt.Printf("%-*s\n", maxLength+1, aurora.Bold(k.Description+":"))
			for _, v := range k.ValueList {
				// fmt.Println(fmt.Sprintf("%-*s", maxLength+1, ""), aurora.Blue(v))
				fmt.Println(aurora.Blue(v))
			}
		}
		if k.ValueComplex != nil {
			fmt.Printf("%-*s\n", maxLength+1, aurora.Bold(k.Description+":"))
			for _, v := range k.ValueComplex {
				fmt.Println(aurora.Blue(v))
			}
		}
	}
}

func (o Output) closeMessagesJson() {
	tempMap := make(map[string]interface{}, len(o.messages))
	for _, k := range o.messages {
		if k.Value != nil {
			tempMap[k.Description] = k.Value
		}
		if k.ValueList != nil {
			tempMap[k.Description] = k.ValueList
		}
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
	if o == nil { // in the case this has not been initialized yet
		return
	}
	if len(o.messages) > 0 {
		if o.jsonOutput {
			o.closeMessagesJson()
		} else {
			o.closeMessagesDefault()
		}
		o.messages = o.messages[:0]
	}

	if len(o.tableHeaders) > 0 {
		if o.jsonOutput {
			o.closeTableJson()
		} else {
			o.closeTableDefault()
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
