package lang_csharp

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/model"
)

// CSharpProcessor implements language.Processor for C#.
type CSharpProcessor struct{}

func NewCSharpProcessor() *CSharpProcessor { return &CSharpProcessor{} }

// Name identifies this processor.
func (p *CSharpProcessor) Name() string { return "csharp" }

// Detect reports whether the path looks like a C# source file.
func (p *CSharpProcessor) Detect(path string) bool {
	if path == "" {
		return false
	}
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".cs" {
		return false
	}
	base := filepath.Base(path)
	return strings.ToLower(base) != ".cs"
}

// AnalyzeFile extracts method-like members and their spans from C# source.
func (p *CSharpProcessor) AnalyzeFile(_ string, lines []string) ([]model.MethodMetrics, error) {
	if lines == nil {
		return []model.MethodMetrics{}, nil
	}
	s := &scanner{lines: lines}
	return s.scan(), nil
}

//
// ================= scanner =================
//

type scanner struct {
	lines []string
	i     int

	nsFirst   string // first segment of namespace
	rootClass string // outermost class/record/struct name only

	braceDepth int
	inMethod   bool

	agg map[string]*model.MethodMetrics // merge (non-ctor) overloads by name
	out []model.MethodMetrics           // constructors/destructors kept separate (not merged)
}

// Scan all lines and collect results.
func (s *scanner) scan() []model.MethodMetrics {
	s.agg = make(map[string]*model.MethodMetrics)

	for s.i = 0; s.i < len(s.lines); s.i++ {
		raw := s.lines[s.i]
		line := stripComments(raw)
		trim := strings.TrimSpace(line)

		// Always advance brace depth
		s.braceDepth += deltaBraces(line)

		if trim == "" {
			continue
		}

		// Namespace detection
		if s.matchNamespace(trim) {
			continue
		}

		// Type declarations (allow modifiers/attributes anywhere before keyword)
		if s.matchTypeDecl(trim) {
			continue
		}

		// Ignore local constructs while inside a method body
		if s.inMethod {
			continue
		}

		// Properties (auto, expression-bodied, full block inc. header on next line)
		if s.matchAutoProperty(trim) {
			continue
		}
		if s.matchExprProperty(trim) {
			continue
		}
		if s.matchPropertyBlockFlexible(trim) {
			continue
		}

		// Indexer (getter span only, spanning accessor lines)
		if s.matchIndexer(trim) {
			continue
		}

		// Methods / constructors / destructors; supports multiline signatures
		if s.matchMemberSignature() {
			continue
		}
	}

	// materialize: constructors/destructors + merged overloads
	result := make([]model.MethodMetrics, 0, len(s.out)+len(s.agg))
	result = append(result, s.out...)
	for _, mm := range s.agg {
		result = append(result, *mm)
	}
	return result
}

//
// ================= regexes & matchers =================
//

var (
	// Namespace, e.g. "namespace My.App.Core"
	reNamespace = regexp.MustCompile(`\bnamespace\s+([A-Za-z_][A-Za-z0-9_\.]*)`)

	// Type declaration with optional attributes/modifiers.
	// Works for: "public partial class Foo", "record Person", "struct S", "interface I", "enum E"
	reTypeDecl = regexp.MustCompile(`\b(class|record|struct|interface|enum)\s+([A-Za-z_][A-Za-z0-9_]*)\b`)

	// Auto property: public T Name { get; set; }
	reAutoProp = regexp.MustCompile(`^(?:\[.*\]\s*)*(?:public|private|protected|internal)(?:\s+static)?\s+([A-Za-z0-9_<>\[\],\.\s\?]+?)\s+([A-Za-z_][A-Za-z0-9_]*)\s*{\s*get;\s*set;\s*}`)

	// Expression-bodied property: public T Area => expr;
	reExprProp = regexp.MustCompile(`^(?:\[.*\]\s*)*(?:public|private|protected|internal)(?:\s+static)?\s+([A-Za-z0-9_<>\[\],\.\s\?]+?)\s+([A-Za-z_][A-Za-z0-9_]*)\s*=>`)

	// Full property block start (brace on same line)
	rePropertyBlock = regexp.MustCompile(`^(?:\[.*\]\s*)*(?:public|private|protected|internal)(?:\s+static)?\s+([A-Za-z0-9_<>\[\],\.\s\?]+?)\s+([A-Za-z_][A-Za-z0-9_]*)\s*{$`)

	// Property header without the opening brace (brace on next line)
	rePropertyHeader = regexp.MustCompile(`^(?:\[.*\]\s*)*(?:public|private|protected|internal)(?:\s+static)?\s+([A-Za-z0-9_<>\[\],\.\s\?]+?)\s+([A-Za-z_][A-Za-z0-9_]*)\s*$`)

	// Indexer start: "public T this["
	reIndexerStart = regexp.MustCompile(`^(?:\[.*\]\s*)*(?:public|private|protected|internal)(?:\s+static)?\s+[A-Za-z0-9_<>\[\],\.\s\?]+?\s+this\s*\[`)

	// Looks like it has a parameter list
	reHasParen = regexp.MustCompile(`\(`)

	// Record positional declaration; skip as member
	reRecordDecl = regexp.MustCompile(`\brecord\s+[A-Za-z_][A-Za-z0-9_]*\s*\(`)

	// Ignore abstract/interface declarations that end with ';'
	reAbstractOrIface = regexp.MustCompile(`\babstract\b|^\s*interface\b`)
	reEndsWithSemi    = regexp.MustCompile(`;\s*$`)

	// Destructor line: "~TypeName()"
	reDestructor = regexp.MustCompile(`^~\s*([A-Za-z_][A-Za-z0-9_]*)\s*\(`)
)

// --- structure ---

func (s *scanner) matchNamespace(trim string) bool {
	m := reNamespace.FindStringSubmatch(trim)
	if m == nil {
		return false
	}
	full := m[1]
	if dot := strings.Index(full, "."); dot > 0 {
		s.nsFirst = full[:dot]
	} else {
		s.nsFirst = full
	}
	return true
}

func (s *scanner) matchTypeDecl(trim string) bool {
	m := reTypeDecl.FindStringSubmatch(trim)
	if m == nil {
		return false
	}
	kind := m[1]
	name := m[2]
	// Only set rootClass for class/record/struct (NOT for interface/enum)
	switch kind {
	case "class", "record", "struct":
		if s.rootClass == "" {
			s.rootClass = name
		}
	}
	return true
}

// --- properties ---

func (s *scanner) matchAutoProperty(trim string) bool {
	m := reAutoProp.FindStringSubmatch(trim)
	if m == nil || s.rootClass == "" {
		return false
	}
	typ := strings.TrimSpace(m[1])
	if isTypeKeyword(typ) {
		return false
	}
	prop := m[2]
	// Expect "Class.get_Name", "Class.set_Name"
	s.emit(s.qual("get_"+prop, false), s.i+1, s.i+1, false)
	s.emit(s.qual("set_"+prop, false), s.i+1, s.i+1, false)
	return true
}

func (s *scanner) matchExprProperty(trim string) bool {
	m := reExprProp.FindStringSubmatch(trim)
	if m == nil || s.rootClass == "" {
		return false
	}
	typ := strings.TrimSpace(m[1])
	if isTypeKeyword(typ) {
		return false
	}
	prop := m[2]
	// "Class.get_Prop"
	s.emit(s.qual("get_"+prop, false), s.i+1, s.i+1, false)
	return true
}

// matchPropertyBlockFlexible handles property blocks where the opening brace
// is on the same line OR on the next non-empty line.
func (s *scanner) matchPropertyBlockFlexible(trim string) bool {
	var m []string
	idx := s.i

	// Case 1: brace on the same line
	if mm := rePropertyBlock.FindStringSubmatch(trim); mm != nil {
		m = mm
	} else if mm := rePropertyHeader.FindStringSubmatch(trim); mm != nil {
		// Case 2: header, brace on next non-empty line
		n := s.peekNextNonEmptyIndex(s.i)
		if n >= 0 {
			next := strings.TrimSpace(stripComments(s.lines[n]))
			if strings.HasPrefix(next, "{") {
				m = mm
				idx = n // start scanning block from the line that contains "{"
			}
		}
	}
	if m == nil || s.rootClass == "" {
		return false
	}

	typ := strings.TrimSpace(m[1])
	if isTypeKeyword(typ) {
		return false
	}
	prop := m[2]

	// From idx (the line with "{"), scan the whole property block; compute accessor spans
	depth := 0
	foundOpen := false
	propEndIdx := -1

	getLine := -1
	setLine := -1
	getEnd := -1
	setEnd := -1

	for j := idx; j < len(s.lines); j++ {
		l := stripComments(s.lines[j])
		t := strings.TrimSpace(l)

		if !foundOpen && strings.Contains(l, "{") {
			foundOpen = true
		}
		if foundOpen {
			// detect accessor starts
			if getLine == -1 && strings.HasPrefix(t, "get") {
				getLine = j + 1
				// single-line get { ... } ?
				if strings.Contains(t, "{") && strings.Contains(t, "}") {
					getEnd = getLine
				} else {
					getEnd = findBlockEndFrom(s.lines, j)
				}
			} else if setLine == -1 && strings.HasPrefix(t, "set") {
				setLine = j + 1
				// single-line set { ... } ?
				if strings.Contains(t, "{") && strings.Contains(t, "}") {
					setEnd = setLine
				} else {
					setEnd = findBlockEndFrom(s.lines, j)
				}
			}

			depth += deltaBraces(l)
			if depth <= 0 {
				propEndIdx = j
				break
			}
		}
	}

	// Emit accessors if found
	if getLine > 0 && setLine > 0 && getEnd == getLine && setEnd == setLine {
		// Both accessors are single-line: emit a single "get_" spanning both lines.
		start := getLine
		end := setLine
		if setLine < getLine {
			start, end = setLine, getLine
		}
		s.emit(s.qual("get_"+prop, false), start, end, false)
	} else {
		if getLine > 0 {
			if getEnd < getLine {
				getEnd = getLine
			}
			s.emit(s.qual("get_"+prop, false), getLine, getEnd, false)
		}
		if setLine > 0 {
			if setEnd < setLine {
				setEnd = setLine
			}
			s.emit(s.qual("set_"+prop, false), setLine, setEnd, false)
		}
	}

	// Advance cursor past the property so we don’t re-scan its inner lines (avoids "if" becoming a method).
	if propEndIdx >= 0 {
		s.i = propEndIdx
	}
	return true
}

// findBlockEndFrom starts reading at index j where a block opens on that same line,
// and returns the 1-based line where the block ends (inclusive).
func findBlockEndFrom(lines []string, j int) int {
	depth := 0
	foundOpen := false
	for m := j; m < len(lines); m++ {
		ll := stripComments(lines[m])
		if !foundOpen && strings.Contains(ll, "{") {
			foundOpen = true
		}
		if foundOpen {
			depth += deltaBraces(ll)
			if depth <= 0 {
				return m + 1
			}
		}
	}
	return len(lines)
}

// --- indexer ---

func (s *scanner) matchIndexer(trim string) bool {
	if !reIndexerStart.MatchString(trim) || s.rootClass == "" {
		return false
	}

	depth := 0
	foundOpen := false
	firstAccessor := -1
	lastAccessor := -1
	propEndIdx := -1

	for j := s.i; j < len(s.lines); j++ {
		l := stripComments(s.lines[j])
		t := strings.TrimSpace(l)
		if !foundOpen && strings.Contains(l, "{") {
			foundOpen = true
		}
		if foundOpen {
			if (strings.HasPrefix(t, "get") || strings.HasPrefix(t, "set") || strings.Contains(t, " get ") || strings.Contains(t, " set ")) && strings.Contains(t, "{") {
				ln := j + 1
				if firstAccessor == -1 {
					firstAccessor = ln
				}
				lastAccessor = ln
			}
			depth += deltaBraces(l)
			if depth <= 0 {
				propEndIdx = j
				break
			}
		}
	}

	if firstAccessor > 0 {
		if lastAccessor < firstAccessor {
			lastAccessor = firstAccessor
		}
		// Name format: "Class.get_this" spanning accessor lines (e.g., 8..9)
		s.emit(s.qual("get_this", false), firstAccessor, lastAccessor, false)
	}

	// Skip inner lines
	if propEndIdx >= 0 {
		s.i = propEndIdx
	}
	return firstAccessor > 0
}

// --- methods / ctors / dtors ---

func (s *scanner) matchMemberSignature() bool {
	line := stripComments(s.lines[s.i])
	if !reHasParen.MatchString(line) {
		return false
	}

	sigStart := s.i
	sig := line
	paren := deltaParens(line)
	j := s.i
	for paren > 0 && j+1 < len(s.lines) {
		j++
		add := stripComments(s.lines[j])
		sig += "\n" + add
		paren += deltaParens(add)
	}

	compact := strings.TrimSpace(sig)

	// Skip record positional declarations
	if reRecordDecl.MatchString(compact) {
		s.i = j
		return true
	}

	// Ignore abstract/interface declarations without body
	if reEndsWithSemi.MatchString(compact) && (reAbstractOrIface.MatchString(compact) || (!strings.Contains(compact, "=>") && !strings.Contains(compact, "{"))) {
		s.i = j
		return true
	}

	// Destructor "~Type()"
	if m := reDestructor.FindStringSubmatch(strings.TrimSpace(compact)); m != nil {
		name := m[1]
		if s.rootClass != "" {
			name = s.rootClass + ".~" + name
		}
		end := s.findBodyEnd(j)
		s.emit(name, sigStart+1, end, true)
		s.i = end - 1
		return true
	}

	// Constructor (identifier equals class name)
	if s.rootClass != "" && looksLikeCtor(compact, s.rootClass) {
		name := s.rootClass + "." + s.rootClass
		end := s.findBodyEnd(j)
		s.emit(name, sigStart+1, end, true)
		s.i = end - 1
		return true
	}

	// Expression-bodied method
	if strings.Contains(compact, "=>") {
		// skip operator overloads (e.g., "operator +")
		if containsOperatorKeyword(compact) {
			s.i = j
			return true
		}
		name := s.extractMethodName(compact)
		if name != "" {
			s.emit(s.qual(name, false), sigStart+1, j+1, false)
		}
		s.i = j
		return true
	}

	// Block-bodied method (either has "{" in compact or on the next non-empty line)
	if strings.Contains(compact, "{") || s.peekNextHasOpenBrace(j) {
		// skip operator overloads
		if containsOperatorKeyword(compact) {
			s.i = j
			return true
		}
		name := s.extractMethodName(compact)
		if name != "" {
			end := s.findBodyEnd(j)
			s.emit(s.qual(name, false), sigStart+1, end, false)
			s.i = end - 1
			return true
		}
	}

	s.i = j
	return true
}

//
// ================= utilities =================
//

func stripComments(s string) string {
	// Remove line comments
	if i := strings.Index(s, "//"); i >= 0 {
		s = s[:i]
	}
	// Remove /* ... */ ranges on same line (best-effort)
	for {
		start := strings.Index(s, "/*")
		if start < 0 {
			break
		}
		end := strings.Index(s, "*/")
		if end < 0 || end < start {
			s = s[:start]
			break
		}
		s = s[:start] + s[end+2:]
	}
	return s
}

func deltaBraces(s string) int { return strings.Count(s, "{") - strings.Count(s, "}") }
func deltaParens(s string) int { return strings.Count(s, "(") - strings.Count(s, ")") }

func (s *scanner) peekNextNonEmptyIndex(i int) int {
	k := i + 1
	for k < len(s.lines) {
		if strings.TrimSpace(stripComments(s.lines[k])) != "" {
			return k
		}
		k++
	}
	return -1
}

func (s *scanner) peekNextHasOpenBrace(j int) bool {
	k := j + 1
	for k < len(s.lines) {
		if strings.TrimSpace(stripComments(s.lines[k])) != "" {
			break
		}
		k++
	}
	if k < len(s.lines) {
		return strings.Contains(stripComments(s.lines[k]), "{")
	}
	return false
}

// findBodyEnd returns the 1-based line number where the block closes (inclusive).
func (s *scanner) findBodyEnd(startIdx int) int {
	depth := 0
	foundOpen := false
	for m := startIdx; m < len(s.lines); m++ {
		l := stripComments(s.lines[m])
		if !foundOpen && strings.Contains(l, "{") {
			foundOpen = true
			s.inMethod = true
		}
		if foundOpen {
			depth += deltaBraces(l)
			if depth <= 0 {
				s.inMethod = false
				return m + 1
			}
		}
	}
	s.inMethod = false
	return len(s.lines)
}

func (s *scanner) extractMethodName(sig string) string {
	// Take token before '(' as method name (handles generics Foo<T>(...))
	idx := strings.Index(sig, "(")
	if idx < 0 {
		return ""
	}
	before := strings.TrimSpace(sig[:idx])
	// skip entire operator signatures
	if containsOperatorKeyword(before) {
		return ""
	}
	parts := strings.Fields(before)
	if len(parts) == 0 {
		return ""
	}
	last := parts[len(parts)-1]
	if gt := strings.Index(last, "<"); gt >= 0 {
		last = last[:gt]
	}
	switch last {
	case "get", "set", "this", "if", "for", "foreach", "while", "switch", "catch", "when":
		return ""
	}
	return last
}

func containsOperatorKeyword(s string) bool {
	// Handle "operator +", "operator -", "operator ==" etc.
	return strings.Contains(s, " operator ")
}

func looksLikeCtor(sig string, className string) bool {
	idx := strings.Index(sig, "(")
	if idx < 0 {
		return false
	}
	before := strings.TrimSpace(sig[:idx])
	parts := strings.Fields(before)
	if len(parts) == 0 {
		return false
	}
	return parts[len(parts)-1] == className
}

func isTypeKeyword(s string) bool {
	switch strings.TrimSpace(s) {
	case "class", "struct", "interface", "record", "enum":
		return true
	default:
		return false
	}
}

func (s *scanner) qual(member string, _ bool) string {
	// Produce "NamespaceFirst.RootClass.Member" or "RootClass.Member"
	name := member
	// For normal methods: ensure "RootClass." prefix
	if s.rootClass != "" &&
		!strings.HasPrefix(member, s.rootClass+".") &&
		!strings.HasPrefix(member, "get_") &&
		!strings.HasPrefix(member, "set_") &&
		member != "get_this" {
		name = s.rootClass + "." + member
	}
	// For property helpers "get_X"/"set_X" and indexer "get_this": "RootClass.get_X" / "RootClass.get_this"
	if strings.HasPrefix(member, "get_") || strings.HasPrefix(member, "set_") || member == "get_this" {
		name = s.rootClass + "." + member
	}
	if s.nsFirst != "" {
		return s.nsFirst + "." + name
	}
	return name
}

func (s *scanner) emit(name string, start, end int, isCtor bool) {
	if isCtor {
		s.out = append(s.out, model.MethodMetrics{Name: name, StartLine: start, EndLine: end})
		return
	}
	if mm, ok := s.agg[name]; ok {
		if start < mm.StartLine {
			mm.StartLine = start
		}
		if end > mm.EndLine {
			mm.EndLine = end
		}
		return
	}
	c := model.MethodMetrics{Name: name, StartLine: start, EndLine: end}
	s.agg[name] = &c
}
