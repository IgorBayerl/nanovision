
# Analyzer Implementation Guide

This directory contains the static analysis engines for NanoVision. Each analyzer uses [Tree-sitter](https://tree-sitter.github.io/) to parse source code and extract metrics.

## Core Philosophy: Atomic Logic Counting

NanoVision uses a strict **Atomic/Leaf Coverage** strategy. We do not count "Container Nodes" (like function bodies or block wrappers) as executable statements. We only count the specific, executable "leaf" instructions.

### Why? (The "Blast Radius" Problem)

If we count containers (like `if_statement`), modifying a single line deep inside a nested block marks the entire parent chain as "modified."

**Incorrect "Container" Counting:**
```go
func main() {           // Counted (Container)
    if true {           // Counted (Container)
        print("hi")     // Counted (Leaf)
    }
}

```

*Result:* Modifying `print("hi")` marks 3 statements as modified (The print, the if, and the function). This inflates patch coverage metrics falsely.

**Correct "Atomic" Counting:**

```go
func main() {           // Ignored
    if true {           // Counted (Condition Logic Only)
        print("hi")     // Counted (Action)
    }
}

```

*Result:* Modifying `print("hi")` marks **only 1** statement as modified.

---

## How to Add a New Language

To add a new language (e.g., Python), create a new directory `internal/analyzer/python/` and implement the `Analyzer` interface.

### Step 1: Define the Atomic Query

Your Tree-sitter query must explicitly select **Actions** and **Control Flow Headers**. Do NOT select parent containers.

**Incorrect (Selects Containers):**

```scheme
(if_statement) @stmt
(for_statement) @stmt

```

**Correct (Selects Logic/Actions):**

```scheme
; --- ATOMIC ACTIONS ---
(expression_statement) @stmt
(assignment) @stmt
(return_statement) @stmt
(raise_statement) @stmt

; --- CONTROL FLOW HEADERS ---
; Only capture the logic that DECIDES the flow, not the whole block.

; IF: Capture the condition
(if_statement 
    condition: (_) @stmt)

; FOR: Capture the iterable
(for_statement 
    right: (_) @stmt)  ; In Python: "for x in ITERABLE"

; WHILE: Capture the condition
(while_statement 
    condition: (_) @stmt)

```

### Step 2: Handle "Anonymous" Children

Some Tree-sitter grammars do not name their children (e.g., Go's `for` loop clauses). You must use **Child Types** or **Anchors** instead of field names.

**Example (Go):**

```scheme
; The grammar has no field name for the clause "i:=0; i<10; i++"
; So we select the child node by its type.
(for_statement
    (for_clause) @stmt)

```

**Example (GDScript):**

```scheme
; The grammar uses 'loop_variable' as the anchor for the loop structure.
(for_statement
    loop_variable: (_) @stmt)

```

### Step 3: Implement the Interface

Your Go code simply wraps the query execution.

```go
const statementQuery = `... paste your query here ...`

func (a *PythonAnalyzer) Analyze(source []byte) (Result, error) {
    // Standard boilerplate:
    // 1. Parse code with Tree-sitter
    // 2. Execute statementQuery
    // 3. Map captures to StatementMetric structs
    // 4. Return result
}

```

---

## How to Create Tests for an Analyzer

When adding a new analyzer, you must create a corresponding `_test.go` file to verify that your queries are strictly selecting atomic statements and not containers.

### Test Strategy

1. **Create a Sample Snippet:** Write a small piece of code in the target language that includes:
* A function definition (should NOT be counted).
* A variable assignment (Atomic).
* A control structure (If/While) with a body.
* A return statement.


2. **Assert Count:** Verify the total number of statements matches the number of *lines* of logic, not the number of syntax nodes.
3. **Assert Precision:** Verify that the `StartLine` and `EndLine` for a control structure match **only** the condition line.

### Example Test Template (Go)

Create `internal/analyzer/python/analyzer_test.go`:

```go
package python

import (
	"testing"
	"[github.com/stretchr/testify/assert](https://github.com/stretchr/testify/assert)"
)

func TestPythonAnalyzer_Statements(t *testing.T) {
	// Sample Python Code
	code := `
def my_func():         # Line 2: Ignored (Container)
    x = 10             # Line 3: Assignment (Counted)
    if x > 5:          # Line 4: If Condition (Counted - Atomic)
        print("hi")    # Line 5: Call (Counted)
        return         # Line 6: Return (Counted)
`
	analyzer := New()
	result, err := analyzer.Analyze([]byte(code))

	assert.NoError(t, err)

	// We expect exactly 4 statements (Lines 3, 4, 5, 6).
	// We do NOT expect Line 2 (def) or the Block of the If.
	assert.Equal(t, 4, len(result.Statements))

	// Verify the 'if' statement is atomic (Line 4 only)
	var ifStmt *StatementMetric
	for _, s := range result.Statements {
		if s.StartLine == 4 {
			ifStmt = &s
			break
		}
	}

	assert.NotNil(t, ifStmt, "Should find the if condition")
	assert.Equal(t, 4, ifStmt.EndLine, "If statement should strictly cover only the condition line")
}

```

### Checklist for New Analyzers

1. [ ] **No Containers:** Ensure `{ }` blocks, function bodies, and class definitions are **ignored**.
2. [ ] **Atomic Conditions:** For `if`, `while`, `switch`, ensure you strictly capture the **condition expression**, not the parent statement.
3. [ ] **Verify Precision:** Run the test and confirm that `EndLine` equals `StartLine` for control structures (unless the condition itself spans multiple lines).
