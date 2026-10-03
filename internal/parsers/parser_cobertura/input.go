package parser_cobertura

import "encoding/xml"

// <coverage>
type CoberturaRoot struct {
	XMLName      xml.Name `xml:"coverage"`
	LineRate     string   `xml:"line-rate,attr"`
	LinesCovered string   `xml:"lines-covered,attr"`
	LinesValid   string   `xml:"lines-valid,attr"`
	Complexity   string   `xml:"complexity,attr"`
	Version      string   `xml:"version,attr"`
	Timestamp    string   `xml:"timestamp,attr"`
	Sources      Sources  `xml:"sources"`
	Packages     Packages `xml:"packages"`
}

// <sources>
type Sources struct {
	Source []string `xml:"source"`
}

// <packages>
type Packages struct {
	Package []PackageXML `xml:"package"`
}

// <package>
type PackageXML struct {
	Name       string     `xml:"name,attr"`
	LineRate   string     `xml:"line-rate,attr"`
	Complexity string     `xml:"complexity,attr"`
	Classes    ClassesXML `xml:"classes"`
}

// <classes>
type ClassesXML struct {
	Class []ClassXML `xml:"class"`
}

// <class>
type ClassXML struct {
	Name       string     `xml:"name,attr"`
	Filename   string     `xml:"filename,attr"`
	LineRate   string     `xml:"line-rate,attr"`
	Complexity string     `xml:"complexity,attr"`
	Methods    MethodsXML `xml:"methods"`
	Lines      LinesXML   `xml:"lines"`
}

// <methods>
type MethodsXML struct {
	Method []MethodXML `xml:"method"`
}

// <method>
type MethodXML struct {
	Name       string   `xml:"name,attr"`
	Signature  string   `xml:"signature,attr"`
	LineRate   string   `xml:"line-rate,attr"`
	Complexity string   `xml:"complexity,attr"`
	Lines      LinesXML `xml:"lines"` // Lines specific to this method
}

// <lines>
type LinesXML struct {
	Line []LineXML `xml:"line"`
}

// <line>
type LineXML struct {
	Number string `xml:"number,attr"`
	Hits   string `xml:"hits,attr"`
}
