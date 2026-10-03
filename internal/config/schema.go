package config

import (
	"fmt"
	"reflect"
	"strings"
)

// Field documents one key of the config file.
type Field struct {
	Key  string `json:"key"`  // dotted path, e.g. "diff.only_changed"; "[]" marks a list item
	Type string `json:"type"` // string, boolean, number, range, list of strings, list, object
	Doc  string `json:"doc"`
	// the value when the key is absent; empty when there is none
	Default string `json:"default,omitempty"`
	// allowed values, when the key takes one of a fixed set
	Values []string `json:"values,omitempty"`
	// the command-line flag that sets the same key, without the dash
	Flag string `json:"flag,omitempty"`
	// valid only in the root config file, not in the config of a folder
	RootOnly bool    `json:"rootOnly,omitempty"`
	Fields   []Field `json:"fields,omitempty"`
}

// MetricDoc documents one metric a config can list.
type MetricDoc struct {
	Name          string `json:"name"`
	Label         string `json:"label"`
	Doc           string `json:"doc"`
	Kind          string `json:"kind"` // "percentage" or "value"
	LowerIsBetter bool   `json:"lowerIsBetter,omitempty"`
}

// OutputFormatDoc documents one value of report_types.
type OutputFormatDoc struct {
	Name    string `json:"name"`
	Writes  string `json:"writes"`
	Doc     string `json:"doc"`
	Default bool   `json:"default,omitempty"` // written when report_types is absent
}

// SchemaDoc is the whole config reference. It is generated from the config
// struct and the metric tables, and feeds `nanovision config docs` and the
// website.
type SchemaDoc struct {
	Fields        []Field                `json:"fields"`
	Metrics       map[string][]MetricDoc `json:"metrics"` // "files" and "methods"
	OutputFormats []OutputFormatDoc      `json:"outputFormats"`
}

// Schema builds the config reference from the doc tags of AppConfig.
func Schema() SchemaDoc {
	scoped := yamlKeys(ScopedConfig{})
	fields := fieldsOf(reflect.TypeOf(AppConfig{}), reflect.ValueOf(*GetDefaultConfig()), "")
	for i := range fields {
		rootOnly := true
		for _, k := range scoped {
			rootOnly = rootOnly && fields[i].Key != k
		}
		fields[i].RootOnly = rootOnly
		if fields[i].Key == "report_types" {
			fields[i].Values = outputFormatNames()
		}
	}

	defaults := GetDefaultConfig().ReportTypes
	formats := make([]OutputFormatDoc, len(OutputFormats))
	for i, f := range OutputFormats {
		formats[i] = OutputFormatDoc{Name: f.Name, Writes: f.Writes, Doc: f.Doc}
		for _, d := range defaults {
			formats[i].Default = formats[i].Default || d == f.Name
		}
	}
	return SchemaDoc{
		OutputFormats: formats,
		Fields:        fields,
		Metrics: map[string][]MetricDoc{
			"files":   metricDocs(FileMetricDefs),
			"methods": metricDocs(MethodMetricDefs),
		},
	}
}

func metricDocs(defs []MetricDef) []MetricDoc {
	docs := make([]MetricDoc, len(defs))
	for i, d := range defs {
		kind := "percentage"
		if d.Value {
			kind = "value"
		}
		docs[i] = MetricDoc{Name: d.Name, Label: d.Label, Doc: d.Doc, Kind: kind, LowerIsBetter: d.LowerIsBetter}
	}
	return docs
}

// fieldsOf lists the yaml keys of a struct type. def is the default value of
// the struct, or the zero Value when there is none (inside a list).
func fieldsOf(t reflect.Type, def reflect.Value, prefix string) []Field {
	var out []Field
	for i := range t.NumField() {
		f := t.Field(i)
		tag := f.Tag.Get("yaml")
		name, opts, _ := strings.Cut(tag, ",")
		var fieldDef reflect.Value
		if def.IsValid() {
			fieldDef = def.Field(i)
		}
		if f.Anonymous && strings.Contains(opts, "inline") {
			out = append(out, fieldsOf(f.Type, fieldDef, prefix)...)
			continue
		}
		if name == "" || name == "-" || !f.IsExported() {
			continue
		}

		field := Field{Key: prefix + name, Doc: f.Tag.Get("doc"), Flag: f.Tag.Get("flag")}
		if v := f.Tag.Get("values"); v != "" {
			field.Values = strings.Split(v, ",")
		}

		ft := f.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
		switch {
		case ft == reflect.TypeOf(Band{}):
			field.Type = "range"
		case ft.Kind() == reflect.Struct:
			field.Type = "object"
			field.Fields = fieldsOf(ft, fieldDef, field.Key+".")
		case ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Struct:
			field.Type = "list"
			field.Fields = fieldsOf(ft.Elem(), reflect.Value{}, field.Key+"[].")
		case ft.Kind() == reflect.Slice:
			field.Type = "list of strings"
			field.Default = defaultString(fieldDef)
		case ft.Kind() == reflect.Bool:
			field.Type = "boolean"
			field.Default = defaultString(fieldDef)
		case ft.Kind() == reflect.String:
			field.Type = "string"
			field.Default = defaultString(fieldDef)
		default:
			field.Type = "number"
			field.Default = defaultString(fieldDef)
		}
		out = append(out, field)
	}
	return out
}

// defaultString renders a default for the docs; a zero value has none.
func defaultString(v reflect.Value) string {
	if !v.IsValid() || v.IsZero() {
		return ""
	}
	if v.Kind() == reflect.Slice {
		parts := make([]string, v.Len())
		for i := range parts {
			parts[i] = fmt.Sprint(v.Index(i).Interface())
		}
		return strings.Join(parts, ", ")
	}
	return fmt.Sprint(v.Interface())
}

// yamlKeys lists the top-level yaml keys of a config struct.
func yamlKeys(v any) []string {
	var keys []string
	for _, f := range fieldsOf(reflect.TypeOf(v), reflect.Value{}, "") {
		keys = append(keys, f.Key)
	}
	return keys
}

// Flatten lists every field of the schema, nested ones included.
func Flatten(fields []Field) []Field {
	var out []Field
	for _, f := range fields {
		out = append(out, f)
		out = append(out, Flatten(f.Fields)...)
	}
	return out
}
