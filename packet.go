package transfer

// Field with database mapping
// source database field name
// target database filed name
// target_type target database field data type
// TODO: converter transfer script template
type Field struct {
	Source     string `mapstructure:"source"`
	Target     string `mapstructure:"target"`
	TargetType string `mapstructure:"target_type"`
	Converter  string `mapstructure:"converter"`
}

// Mapping for database fields map
type Mapping []Field

// Map return field mapping
func (mapping Mapping) Map() M {
	m := make(M, len(mapping))
	for _, f := range mapping {
		m[f.Target] = f
	}
	return m
}

// FieldMap return select field mapping
func (mapping Mapping) FieldMap(join ...string) M {
	var prefix, suffix string

	if len(join) > 0 {
		prefix = join[0]
	}

	if len(join) > 1 {
		suffix = join[1]
	}

	m := make(M, len(mapping))
	for _, f := range mapping {
		m[f.Target] = prefix + f.Source + suffix
	}
	return m
}

// Fields return select fields array
func (mapping Mapping) Fields(join ...string) (fields []string) {
	var prefix, suffix string

	if len(join) > 0 {
		prefix = join[0]
	}

	if len(join) > 1 {
		suffix = join[1]
	}

	for _, field := range mapping {
		fields = append(fields, prefix+field.Target+suffix)
	}
	return
}

// M transfer data type
type M map[string]interface{}

// Primitive return json Marshal then Unmarshal result
func (m M) Primitive() (d M, err error) {
	byteData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(byteData), &d)
	return
}

// Packet for transfer data
type Packet []M
