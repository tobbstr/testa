package assert

import (
	"bytes"
	"runtime"
	"strings"
	"text/template"
)

const (
	messageTmpl = `

Assertion failed!
	Description: {{.Msg}}
	Expected: {{.Want}}
	Observed:  {{.Got}}

Call stack:
{{range .Entries}}
	{{.Filename}}.{{.Line}}: {{.FuncName}}{{end}}
	
	 `
)

var preparsedMessageTmpl *template.Template = template.Must(template.New("message").Parse(messageTmpl))

type callStackEntry struct {
	Filename string
	FuncName string
	Line     int
}

func rawCallStack() []callStackEntry {
	var callStack []callStackEntry

	errorfCallFound := false
	for i := 0; ; i++ {
		programCounter, filename, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		runtimeFn := runtime.FuncForPC(programCounter)
		funcName := runtimeFn.Name()

		if !errorfCallFound {
			if strings.Contains(funcName, "(*asserter).errorf") {
				errorfCallFound = true
			}

			continue
		}

		if funcName == "testing.tRunner" {
			break
		}

		callStack = append(callStack, callStackEntry{
			Filename: filename,
			FuncName: funcName,
			Line:     line,
		})
	}
	return callStack
}

func callStack() []callStackEntry {
	var formattedCallStack []callStackEntry
	for _, rawEntry := range rawCallStack() {
		formattedStackEntry := callStackEntry{}

		filenameLastSlashIdx := strings.LastIndex(rawEntry.Filename, "/")
		formattedStackEntry.Filename = rawEntry.Filename[filenameLastSlashIdx+1:]

		funcLastSlashIdx := strings.LastIndex(rawEntry.FuncName, "/")
		formattedStackEntry.FuncName = rawEntry.FuncName[funcLastSlashIdx+1:]

		formattedStackEntry.Line = rawEntry.Line
		formattedCallStack = append(formattedCallStack, formattedStackEntry)
	}
	return formattedCallStack
}

func errorMsg(msg string, want, got interface{}, assertHasWantParam bool) string {
	var buf bytes.Buffer

	type messageValues struct {
		Msg     string
		Want    interface{}
		Got     interface{}
		Entries []callStackEntry
	}

	msgValues := messageValues{
		Msg:     msg,
		Want:    want,
		Got:     got,
		Entries: callStack(),
	}

	if !assertHasWantParam {
		msgValues.Want = "N/A"
	}

	if isFunc(want) {
		msgValues.Want = "function"
	} else if isChan(want) {
		msgValues.Want = "chan"
	}
	if isFunc(got) {
		msgValues.Got = "function"
	} else if isChan(got) {
		msgValues.Got = "chan"
	}

	if err := preparsedMessageTmpl.Execute(&buf, msgValues); err != nil {
		panic(err)
	}

	return buf.String()
}
