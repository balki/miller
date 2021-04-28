package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameFillEmpty = "fill-empty"
const defaultFillEmptyString = "N/A"

var FillEmptySetup = transforming.TransformerSetup{
	Verb:         verbNameFillEmpty,
	UsageFunc:    transformerFillEmptyUsage,
	ParseCLIFunc: transformerFillEmptyParseCLI,
	IgnoresInput: false,
}

func transformerFillEmptyUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameFillEmpty)
	fmt.Fprintf(o, "Fills empty-string fields with specified fill-value.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-v {string} Fill-value: defaults to \"%s\"\n", defaultFillEmptyString)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerFillEmptyParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	fillString := defaultFillEmptyString

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerFillEmptyUsage(os.Stdout, true, 0)

		} else if opt == "-v" {
			fillString = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerFillEmptyUsage(os.Stderr, true, 1)
		}
	}

	transformer, _ := NewTransformerFillEmpty(fillString)

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerFillEmpty struct {
	fillValue *types.Mlrval
}

func NewTransformerFillEmpty(fillString string) (*TransformerFillEmpty, error) {
	this := &TransformerFillEmpty{
		fillValue: types.MlrvalPointerFromString(fillString),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerFillEmpty) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if pe.Value.IsEmpty() {
				pe.Value = this.fillValue
			}
		}

		outputChannel <- inrecAndContext

	} else { // end of record stream
		outputChannel <- inrecAndContext // emit end-of-stream marker
	}
}
