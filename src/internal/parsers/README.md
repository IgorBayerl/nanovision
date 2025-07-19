# Creating a New Parser

This guide provides a comprehensive walkthrough for adding a new coverage report parser to the `go-report-generator` project.

## 1. The Role of a Parser

Before writing code, it's essential to understand the architecture. A parser's single responsibility is to act as a **translator**. It takes a specific coverage report format (like *Cobertura*, *GoCover*, *JaCoCo*, etc.) and translates it into the project's **universal data model**.

#### What a Parser Does:

*   **Reads** a specific file format (`.xml`, `.json`, `.out`, etc.).
*   **Understands** the structure and semantics of that format.
*   **Translates** the format's data into the `internal/model` structs.
*   **Applies** filters (`Assembly`, `Class`, `File`) provided by the configuration to exclude unwanted data *during* the parsing process.

#### What a Parser **Does Not** Do:

*   It **does not** generate HTML, text summaries, or any other output file. That is the job of a `Reporter`.
*   It **does not** merge multiple report files. That is the job of the `Analyzer`.
*   It **does not** contain any user-facing state; it is a pure, stateless data transformation component.

The application framework handles the rest:
1.  The **Parser Factory** (`factory.go`) auto-detects which parser to use for a given file.
2.  The **Analyzer** (`internal/analyzer/`) merges the results from one or more parsed files.
3.  The **Reporters** (`internal/reporter/`) take the final, unified data and generate the output files.

## 2. The `IParser` Interface: Your Contract

This is the most important part of the parser core. Your new parser **must** implement this interface, which is defined in `internal/parsers/parser_config.go`. It is the bridge between the application framework and your specific implementation.

| Method | Responsibility | Implementation Notes |
| :--- | :--- | :--- |
| **`Name() string`** | Return the unique, human-readable name of your parser (e.g., "GoCover", "JaCoCo"). | This name is used in logs and potentially in the UI, so make it descriptive. |
| **`SupportsFile(filePath string) bool`** | Quickly and efficiently determine if your parser can handle the given file. | **This check must be fast.** The factory calls this on every available parser for every input file. Do not read the entire file here. <br> • **For XML:** Read just enough to find the root element name. <br> • **For JSON:** Check for a unique top-level key. <br> • **For line-based formats:** Check if the first line contains a specific "magic string" (e.g., `mode: set` for Go coverage). |
| **`Parse(filePath string, config ParserConfig) (*ParserResult, error)`** | Read the entire report file and perform the full translation into the `parsers.ParserResult` struct. | This is where the main logic resides. You have access to filters (`config.FileFilters()`, etc.) to exclude data as you process it. This method should be **stateless**—all necessary context comes from the `filePath` and `config` arguments. This design allows the application to run multiple `Parse` operations in parallel in the future. |

## 3. Core Principles & Best Practices

Follow these principles to ensure your parser is robust, testable, and aligns with the project's architecture.

*   **Be Stateless and Parallelizable:** A parser instance should not hold any state related to a specific `Parse` operation. All data should be passed in via arguments (`filePath`, `config`) and returned in the `ParserResult`. This ensures that a single parser instance can be safely used to parse multiple files concurrently without interference.

*   **Encapsulate Your Logic:** All code and data structures specific to your parser should live within its own package (e.g., `internal/parsers/yourformat/`). This includes format-specific structs, processing logic, and tests.

*   **Filter Early, Filter Often:** Apply the filters provided in the `ParserConfig` as you process the data. For example, if an assembly or class is excluded, don't waste time processing its files and lines. This improves performance.

*   **Separate Raw Parsing from Translation:** It is a best practice to first unmarshal or parse the input file into structs that *exactly match the source format* (defined in your `input.go` file). Then, in a separate step, iterate over these raw structs and translate them into the universal `model` objects. This separation makes the code easier to read and debug.

*   **Handle Errors Gracefully:** Don't let your parser crash the whole application on a single malformed line or element.
    *   For **recoverable issues** (e.g., a source file not found on disk), log a warning with `slog.Warn` and continue processing. The report can still be generated, albeit without some source code.
    *   For **unrecoverable issues** (e.g., the report file is fundamentally invalid XML), return a descriptive `error` from the `Parse` method.

*   **Prioritize Testability:** Use interfaces to abstract away external dependencies like the file system. The existing `FileReader` interface in the Cobertura parser is a perfect template for this, allowing you to mock file reads during unit tests.

## 4. The Universal Data Model: The "Language" of the Application

Your parser's main job is to produce a `[]model.Assembly` populated with data. Understanding the abstract meaning of these structs is critical, especially for non-.NET formats.

| Struct in `internal/model/` | C# Origin | Abstract Meaning & How to Map Your Format |
| :--- | :--- | :--- |
| **`Assembly`** | `Assembly` | **The largest logical unit of code; a project or module.** The name `Assembly` is kept for consistency with the original .NET ReportGenerator. For non-.NET formats, map this concept appropriately: <br> • **Go:** The module name from `go.mod`. <br> • **Java:** The root package name or the JAR/WAR name. <br> • **JavaScript/TypeScript:** The project name from `package.json`. |
| **`Class`** | `Class` | **A logical grouping of source files within an Assembly.** This provides the primary grouping in the report's UI. <br> • **Go:** A Go package path (e.g., `internal/utils`). <br> • **Java:** A Java package path. <br> • **Python:** A Python module. |
| **`CodeFile`** | `CodeFile` | **A single source code file.** Your parser is responsible for finding the absolute path to this file on disk using the provided source directories. |
| **`Method`** | `Method` | **A function or method within a Class.** Some formats (like Cobertura) provide this data. Many (like `go cover`) do not. If your format doesn't provide method-level data, you can simply leave the `Methods` slice empty on the `Class` struct. |
| **`Line`** | `LineAnalysis` | **A single line in a source file.** This is the most granular piece of data, containing the hit count (`Hits`) and branch information. Your parser must populate this for every relevant line. |
| **`Branch`** | `Branch` | **A conditional branch on a line (e.g., if/else).** If your format does not support branch coverage, you can ignore these fields. The framework will automatically adapt the UI. |

## 5. Step-by-Step Guide: Creating a New Parser

Follow this structure to create a parser for a new format (e.g., "YourFormat").

#### Step 1: Create the Parser Package

Create a new, self-contained directory inside `internal/parsers/`.

```
internal/parsers/
└── yourformat/            #<-- New directory for your parser
    ├── parser.go          #<-- The public entrypoint (IParser implementation)
    ├── processing.go      #<-- The core translation logic
    ├── input.go           #<-- Structs for your raw input data
    ├── interfaces.go      #<-- (Optional) Interfaces for testing
    └── processing_test.go #<-- Tests for your logic
```

#### Step 2: Define the Input Structs (`input.go`)

In `internal/parsers/yourformat/input.go`, define Go structs that map directly to your input file format.

```go
// in: internal/parsers/yourformat/input.go
package yourformat

type YourFormatRoot struct { /* ... */ }
```

#### Step 3: Implement the `IParser` Interface (`parser_config.go`)

This file is the public entrypoint. It implements `IParser`.

```go
// in: internal/parsers/yourformat/parser.go
package yourformat

import (
	// ... your imports
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/parsers"
)

type YourFormatParser struct{
    // ... dependencies like a file reader
}

func NewYourFormatParser(/*...deps*/) parsers.IParser {
	return &YourFormatParser{/*...*/}
}

func (p *YourFormatParser) Name() string {
	return "YourFormat"
}

func (p *YourFormatParser) SupportsFile(filePath string) bool {
	// ... quick check logic
}

func (p *YourFormatParser) Parse(filePath string, config parsers.ParserConfig) (*parsers.ParserResult, error) {
	// ... main parsing and translation logic
}
```

#### Step 4: Implement the Core Logic (`processing.go`)

This file does the actual work of translating your format into the `model` objects.

```go
// in: internal/parsers/yourformat/processing.go
package yourformat

// ...
func (p *yourFormatProcessor) Process(report *YourFormatRoot) (assemblies []model.Assembly, /*...*/) {
    // Main translation logic here...
	return nil, nil, nil // Placeholder
}
```

#### Step 5: Register the New Parser

To make the application aware of your new parser, you must add it to the `ParserFactory` in the application's entrypoint.

Navigate to `cmd/main.go` and find the `run()` function. Inside, locate the `parsers.NewParserFactory` call and add an instance of your new parser to the list.

```go
// in: cmd/main.go

// ... imports
import (
    // ... other imports
    "github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/parsers/yourformat" // 1. Import your new package
)

func run() error {
    // ...
    // 2. Create all desired parsers and the factory that holds them.
    prodFileReader := filereader.NewDefaultReader()
    parserFactory := parsers.NewParserFactory(
        cobertura.NewCoberturaParser(prodFileReader),
        gocover.NewGoCoverParser(prodFileReader),
        yourformat.NewYourFormatParser(prodFileReader), // 2. Add your new parser here
    )
    // ...
}
```

#### Step 6: Write Tests

Create a `_test.go` file for your parsers. Add sample report files to a `testdata` directory and write tests that parse them and assert that the resulting `model` structs are populated correctly. Refer to `internal/parsers/cobertura/processing_test.go` for a comprehensive example.