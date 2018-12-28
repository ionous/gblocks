package gblocks

const (
	opt_message0 = "message0"
	opt_args0    = "args0"
	opt_output   = "output"

	opt_previousStatement = "previousStatement"
	opt_nextStatement     = "nextStatement"
	opt_name              = "name"
	opt_type              = "type"
	opt_check             = "check"
	opt_options           = "options"
	opt_mutation          = "mutation" // custom for type mutation

	opt_value     = "value"
	opt_min       = "min"
	opt_max       = "max"
	opt_text      = "text"
	opt_precision = "precision"
	opt_readOnly  = "readOnly"

	input_statement = "input_statement"
	input_value     = "input_value"
	input_dummy     = "input_dummy"

	field_angle    = "field_angle"
	field_checkbox = "field_checkbox"
	field_colour   = "field_colour"
	field_date     = "field_date"
	field_dropdown = "field_dropdown"
	field_image    = "field_image"
	field_label    = "field_label"
	field_number   = "field_number" //options['value'], options['min'], options['max'], options['precision']
	field_input    = "field_input"  // text input options['spellcheck''],
	field_variable = "field_variable"

	PreviousStatementField = "PreviousStatement"
	NextStatementField     = "NextStatement"
	OutputMethod           = "Output"

	PREVIOUS_STATEMENT InputName = "PREVIOUS_STATEMENT"
	NEXT_STATEMENT     InputName = "NEXT_STATEMENT"

	// colour            = "colour"
	// helpUrl           = "helpUrl"
	// mutator: "controls_if_mutator",
	// extensions: ["controls_if_tooltip"]
)
