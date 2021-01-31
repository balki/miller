package transformers

import (
	"container/list"
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameGroupBy = "group-by"

var GroupBySetup = transforming.TransformerSetup{
	Verb:         verbNameGroupBy,
	ParseCLIFunc: transformerGroupByParseCLI,
	UsageFunc:    transformerGroupByUsage,
	IgnoresInput: false,
}

func transformerGroupByParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerGroupByUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerGroupByUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	// Get the group-by field names from the command line
	if argi >= argc {
		transformerGroupByUsage(os.Stderr, true, 1)
	}
	groupByFieldNames := lib.SplitString(args[argi], ",")
	argi += 1

	transformer, _ := NewTransformerGroupBy(
		groupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerGroupByUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s {comma-separated field names}]\n", os.Args[0], verbNameGroupBy)
	fmt.Fprint(o,
		`Outputs records in batches having identical values at specified field names.
`)

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerGroupBy struct {
	// input
	groupByFieldNames []string

	// state
	// map from string to *list.List
	recordListsByGroup *lib.OrderedMap
}

func NewTransformerGroupBy(
	groupByFieldNames []string,
) (*TransformerGroupBy, error) {

	this := &TransformerGroupBy{
		groupByFieldNames: groupByFieldNames,

		recordListsByGroup: lib.NewOrderedMap(),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerGroupBy) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
		if !ok {
			return
		}

		recordListForGroup := this.recordListsByGroup.Get(groupingKey)
		if recordListForGroup == nil {
			recordListForGroup = list.New()
			this.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		recordListForGroup.(*list.List).PushBack(inrecAndContext)

	} else {
		for outer := this.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value.(*list.List)
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				outputChannel <- inner.Value.(*types.RecordAndContext)
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
