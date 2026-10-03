package blob

import "math"

const analysisVersion = 1

// Function is one method or function found by an analyzer.
type Function struct {
	Name      string `json:"name"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	// nil when the analyzer gives no complexity
	Complexity *int `json:"complexity,omitempty"`
}

type Statement struct {
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Type      string `json:"type,omitempty"`
}

// Analysis is what the static analysis found in one version of a file. The
// order of the functions and statements is the order of the analyzer.
type Analysis struct {
	TotalLines int         `json:"totalLines"`
	Functions  []Function  `json:"functions,omitempty"`
	Statements []Statement `json:"statements,omitempty"`
}

func EncodeAnalysis(a Analysis) []byte {
	e := newEncoder(KindAnalysis, analysisVersion, 16*len(a.Functions)+3*len(a.Statements)+8)
	e.count(a.TotalLines)

	e.count(len(a.Functions))
	for _, f := range a.Functions {
		e.string(f.Name)
		e.varint(int64(f.StartLine))
		e.varint(int64(f.EndLine) - int64(f.StartLine))
		if f.Complexity == nil || *f.Complexity < 0 {
			e.uvarint(0)
		} else {
			e.uvarint(uint64(*f.Complexity) + 1)
		}
	}

	// statement types repeat, so they are written once and referenced by index
	typeIndex := make(map[string]int)
	var types []string
	for _, s := range a.Statements {
		if _, ok := typeIndex[s.Type]; !ok {
			typeIndex[s.Type] = len(types)
			types = append(types, s.Type)
		}
	}
	e.count(len(types))
	for _, t := range types {
		e.string(t)
	}

	e.count(len(a.Statements))
	prev := 0
	for _, s := range a.Statements {
		e.varint(int64(s.StartLine) - int64(prev))
		e.varint(int64(s.EndLine) - int64(s.StartLine))
		e.count(typeIndex[s.Type])
		prev = s.StartLine
	}
	return e.buf
}

func DecodeAnalysis(data []byte) (Analysis, error) {
	d, _, err := newDecoder(data, KindAnalysis, analysisVersion)
	if err != nil {
		return Analysis{}, err
	}
	var a Analysis
	if a.TotalLines, err = d.count(); err != nil {
		return Analysis{}, err
	}

	nFuncs, err := d.length(4)
	if err != nil {
		return Analysis{}, err
	}
	if nFuncs > 0 {
		a.Functions = make([]Function, nFuncs)
	}
	for i := range nFuncs {
		f := &a.Functions[i]
		if f.Name, err = d.string(); err != nil {
			return Analysis{}, err
		}
		if f.StartLine, err = d.int32(); err != nil {
			return Analysis{}, err
		}
		span, err := d.int32()
		if err != nil {
			return Analysis{}, err
		}
		if f.EndLine, err = addInt32(f.StartLine, span); err != nil {
			return Analysis{}, err
		}
		cc, err := d.uvarint()
		if err != nil {
			return Analysis{}, err
		}
		if cc > 0 {
			if cc > math.MaxInt32 {
				return Analysis{}, corrupt("complexity %d out of range", cc)
			}
			v := int(cc) - 1
			f.Complexity = &v
		}
	}

	nTypes, err := d.length(1)
	if err != nil {
		return Analysis{}, err
	}
	types := make([]string, nTypes)
	seen := make(map[string]bool, nTypes)
	for i := range nTypes {
		if types[i], err = d.string(); err != nil {
			return Analysis{}, err
		}
		if seen[types[i]] {
			return Analysis{}, corrupt("statement type %q listed twice", types[i])
		}
		seen[types[i]] = true
	}

	nStmts, err := d.length(3)
	if err != nil {
		return Analysis{}, err
	}
	if nStmts > 0 {
		a.Statements = make([]Statement, nStmts)
	}
	used := make([]bool, nTypes)
	prev := 0
	for i := range nStmts {
		s := &a.Statements[i]
		delta, err := d.int32()
		if err != nil {
			return Analysis{}, err
		}
		if s.StartLine, err = addInt32(prev, delta); err != nil {
			return Analysis{}, err
		}
		span, err := d.int32()
		if err != nil {
			return Analysis{}, err
		}
		if s.EndLine, err = addInt32(s.StartLine, span); err != nil {
			return Analysis{}, err
		}
		idx, err := d.count()
		if err != nil {
			return Analysis{}, err
		}
		if idx >= nTypes {
			return Analysis{}, corrupt("statement type %d of %d", idx, nTypes)
		}
		// the encoder numbers types in order of first use
		if !used[idx] && (idx > 0 && !used[idx-1]) {
			return Analysis{}, corrupt("statement type %d used before %d", idx, idx-1)
		}
		used[idx] = true
		s.Type = types[idx]
		prev = s.StartLine
	}
	for i, u := range used {
		if !u {
			return Analysis{}, corrupt("statement type %d is never used", i)
		}
	}
	return a, d.done()
}

// int32 reads a signed varint that fits in 32 bits, the range of a line number.
func (d *decoder) int32() (int, error) {
	v, err := d.varint()
	if err != nil {
		return 0, err
	}
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, corrupt("value %d out of range", v)
	}
	return int(v), nil
}

func addInt32(a, b int) (int, error) {
	sum := int64(a) + int64(b)
	if sum < math.MinInt32 || sum > math.MaxInt32 {
		return 0, corrupt("line %d out of range", sum)
	}
	return int(sum), nil
}
