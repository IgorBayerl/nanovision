package lang_csharp

import (
	"strings"
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helpers
type triple struct {
	Name      string
	StartLine int
	EndLine   int
}

func analyze(t *testing.T, src string) []model.MethodMetrics {
	t.Helper()
	p := NewCSharpProcessor()
	lines := strings.Split(src, "\n")
	methods, err := p.AnalyzeFile("test.cs", lines)
	require.NoError(t, err)
	require.NotNil(t, methods)
	return methods
}

func bag(methods []model.MethodMetrics) map[triple]int {
	m := map[triple]int{}
	for _, mm := range methods {
		k := triple{mm.Name, mm.StartLine, mm.EndLine}
		m[k]++
	}
	return m
}

func assertMatched(t *testing.T, got []model.MethodMetrics, want []model.MethodMetrics) {
	gb := bag(got)
	wb := bag(want)
	assert.Equalf(t, wb, gb, "method bag mismatch\nGot:  %#v\nWant: %#v", got, want)
}

// Detect
func Test_CSharpProcessor_Detect(t *testing.T) {
	p := NewCSharpProcessor()
	tests := []struct {
		name string
		path string
		ok   bool
	}{
		{"plain .cs", "Program.cs", true},
		{"path .cs", "src/app/Service.cs", true},
		{"uppercase", "FILE.CS", true},
		{"empty", "", false},
		{"just extension", ".cs", false},
		{"almost .cs", "Class.cs.bak", false},
		{"other ext", "Class.java", false},
	}
	for _, tt := range tests {
		t.Run(strings.ReplaceAll(tt.name, " ", "_"), func(t *testing.T) {
			assert.Equal(t, tt.ok, p.Detect(tt.path))
		})
	}
}

// Core extraction
func Test_Methods_Basic(t *testing.T) {
	src := `
public class Calc {
    public int Add(int a, int b) {
        return a + b;
    }

    private void Touch() {
    }
}
`
	got := analyze(t, src)
	want := []model.MethodMetrics{
		{Name: "Calc.Add", StartLine: 3, EndLine: 5},
		{Name: "Calc.Touch", StartLine: 7, EndLine: 8},
	}
	assertMatched(t, got, want)
}

func Test_Constructors_And_Destructors(t *testing.T) {
	src := `
public class Resource {
    public Resource() {
    }

    ~Resource() {
        // cleanup
    }

    static Resource() {
        // type initializer
    }
}
`
	got := analyze(t, src)
	want := []model.MethodMetrics{
		{Name: "Resource.Resource", StartLine: 3, EndLine: 4},   // instance ctor
		{Name: "Resource.~Resource", StartLine: 6, EndLine: 8},  // dtor
		{Name: "Resource.Resource", StartLine: 10, EndLine: 12}, // static ctor
	}
	assertMatched(t, got, want)
}

func Test_Async_Methods(t *testing.T) {
	src := `
using System.Threading.Tasks;
public class Svc {
    public async System.Threading.Tasks.Task<string> GetAsync(int id) {
        return await System.Threading.Tasks.Task.FromResult(id.ToString());
    }
}
`
	got := analyze(t, src)
	want := []model.MethodMetrics{
		{Name: "Svc.GetAsync", StartLine: 4, EndLine: 6},
	}
	assertMatched(t, got, want)
}

func Test_Operators_And_Indexers(t *testing.T) {
	src := `
public class Vector {
    public static Vector operator +(Vector a, Vector b) {
        return new Vector();
    }

    public int this[int i] {
        get { return 0; }
        set { }
    }
}
`
	got := analyze(t, src)
	// We intentionally IGNORE operators for now; only the indexer getter is tracked.
	want := []model.MethodMetrics{
		{Name: "Vector.get_this", StartLine: 8, EndLine: 9},
	}
	assertMatched(t, got, want)
}

func Test_Properties_Auto_ExpressionBodied_And_FullBody(t *testing.T) {
	t.Run("auto-property_emits_get_/set__with_single-line_span", func(t *testing.T) {
		src := `
public class Person {
    public string Name { get; set; }
}
`
		got := analyze(t, src)
		want := []model.MethodMetrics{
			{Name: "Person.get_Name", StartLine: 3, EndLine: 3},
			{Name: "Person.set_Name", StartLine: 3, EndLine: 3},
		}
		assertMatched(t, got, want)
	})

	t.Run("expression-bodied_property_emits_Class.get_", func(t *testing.T) {
		src := `
public class Circle {
    public double Radius { get; set; }
    public double Area => System.Math.PI * Radius * Radius;
}
`
		got := analyze(t, src)
		want := []model.MethodMetrics{
			{Name: "Circle.get_Radius", StartLine: 3, EndLine: 3},
			{Name: "Circle.set_Radius", StartLine: 3, EndLine: 3},
			{Name: "Circle.get_Area", StartLine: 4, EndLine: 4},
		}
		assertMatched(t, got, want)
	})

	t.Run("full-body_property_emits_single_get__spanning_accessor_lines", func(t *testing.T) {
		src := `
public class Temperature {
    private double c;
    public double Celsius {
        get { return c; }
        set { c = value; }
    }
}
`
		got := analyze(t, src)
		want := []model.MethodMetrics{
			{Name: "Temperature.get_Celsius", StartLine: 5, EndLine: 6},
		}
		assertMatched(t, got, want)
	})
}

func Test_NestedTypes_And_Namespaces(t *testing.T) {
	src := `
namespace MyApp.Core {
    public class Outer {
        public void A() {}
        public class Inner {
            public void B() { }
        }
    }
}
`
	got := analyze(t, src)
	want := []model.MethodMetrics{
		{Name: "MyApp.Outer.A", StartLine: 4, EndLine: 4},
		{Name: "MyApp.Outer.B", StartLine: 6, EndLine: 6},
	}
	assertMatched(t, got, want)
}

func Test_Multiline_Signatures_And_SingleLine_Bodies(t *testing.T) {
	src := `
public class Weird {
    public int Sum(
        int a,
        int b
    )
    {
        return a + b;
    }

    public void OneLiner() { DoThing(); }
}
`
	got := analyze(t, src)
	want := []model.MethodMetrics{
		{Name: "Weird.Sum", StartLine: 3, EndLine: 9},
		{Name: "Weird.OneLiner", StartLine: 11, EndLine: 11},
	}
	assertMatched(t, got, want)
}

func Test_Comments_Noise_And_Using_Are_Ignored(t *testing.T) {
	src := `
using System;
// line comment
public class T {
    /* block
       comment */
    /// xml doc
    public void M() {
        // inside
    }
}
`
	got := analyze(t, src)
	want := []model.MethodMetrics{
		{Name: "T.M", StartLine: 8, EndLine: 10},
	}
	assertMatched(t, got, want)
}

func Test_Interface_And_Abstract_Declarations_Are_Ignored(t *testing.T) {
	src := `
public interface I {
    void Absent();
}

public abstract class B {
    public abstract void AbstractGone();
    public virtual void V() {}
}
`
	got := analyze(t, src)
	want := []model.MethodMetrics{
		{Name: "B.V", StartLine: 8, EndLine: 8},
	}
	assertMatched(t, got, want)
}

func Test_Generics_And_Overloads(t *testing.T) {
	src := `
public class G {
    public T Id<T>(T x) { return x; }
    public int Id(int x) { return x; }
}
`
	got := analyze(t, src)
	// overloads merged: min start 3, max end 4
	want := []model.MethodMetrics{
		{Name: "G.Id", StartLine: 3, EndLine: 4},
	}
	assertMatched(t, got, want)
}

func Test_Records(t *testing.T) {
	t.Run("positional_record_(no_body)_yields_no_methods", func(t *testing.T) {
		src := `public record Person(string Name, int Age);`
		got := analyze(t, src)
		assert.Empty(t, got)
	})

	t.Run("record_with_body_yields_methods", func(t *testing.T) {
		src := `
public record Person(string Name, int Age) {
    public void Greet() { }
}
`
		got := analyze(t, src)
		want := []model.MethodMetrics{
			{Name: "Person.Greet", StartLine: 3, EndLine: 3},
		}
		assertMatched(t, got, want)
	})
}

func Test_LocalFunctions_Are_Ignored(t *testing.T) {
	src := `
public class C {
    public void M() {
        int Add(int x, int y) { return x+y; } // local function
    }
}
`
	got := analyze(t, src)
	want := []model.MethodMetrics{
		{Name: "C.M", StartLine: 3, EndLine: 5},
	}
	assertMatched(t, got, want)
}

// This reproduces your PartialClass.cs and asserts we:
// - detect the two normal methods
// - emit get_/set_ accessors with correct spans
// - DO NOT invent a fake method named "if"
func Test_PartialClass_Property_GetSet_NoFakeIf(t *testing.T) {
	src := `using System;

namespace Test
{
    partial class PartialClass
    {
        public void ExecutedMethod_1()
        {
            Console.WriteLine("Test");
        }

        public void UnExecutedMethod_1()
        {
            Console.WriteLine("Test");
        }

        private int someProperty;

        public int SomeProperty
        {
            get { return this.someProperty; }

            set
            {
                if (value < 0)
                {
                    this.someProperty = 0;
                }
                else
                {
                    this.someProperty = value;
                }
            }
        }
    }
}
`
	lines := strings.Split(src, "\n")
	p := NewCSharpProcessor()
	got, err := p.AnalyzeFile("PartialClass.cs", lines)
	require.NoError(t, err)

	want := []model.MethodMetrics{
		{Name: "Test.PartialClass.ExecutedMethod_1", StartLine: 7, EndLine: 10},
		{Name: "Test.PartialClass.UnExecutedMethod_1", StartLine: 12, EndLine: 15},
		{Name: "Test.PartialClass.get_SomeProperty", StartLine: 21, EndLine: 21},
		{Name: "Test.PartialClass.set_SomeProperty", StartLine: 23, EndLine: 33},
	}

	// Compare as a bag for stability.
	gb := toBag(got)
	wb := toBag(want)
	if !equalBags(wb, gb) {
		t.Fatalf("method bag mismatch\nGot:  %#v\nWant: %#v", got, want)
	}

	// Extra safety: ensure nothing named "...if" slipped in.
	for _, m := range got {
		if strings.HasSuffix(m.Name, ".if") {
			t.Fatalf("spurious method detected: %q", m.Name)
		}
	}
}

func toBag(list []model.MethodMetrics) map[string]int {
	m := map[string]int{}
	for _, mm := range list {
		key := mm.Name
		key += "|" + itoa(mm.StartLine) + "|" + itoa(mm.EndLine)
		m[key]++
	}
	return m
}

func equalBags(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, va := range a {
		if b[k] != va {
			return false
		}
	}
	return true
}

func itoa(i int) string {
	// small fast int->string without fmt to keep test lightweight
	const digits = "0123456789"
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var buf [20]byte
	n := len(buf)
	for i > 0 {
		n--
		buf[n] = digits[i%10]
		i /= 10
	}
	if neg {
		n--
		buf[n] = '-'
	}
	return string(buf[n:])
}
