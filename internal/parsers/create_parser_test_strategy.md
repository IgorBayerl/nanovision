### `internal/parsers/create_parser_test_strategy.md` (Updated)

Hello there, future contributor! 👋

Welcome to the AdlerCov parser development guide. We're thrilled you're interested in adding support for a new coverage report format. This document is designed to make that process as smooth and clear as possible.

The core philosophy of this project is a clean, four-stage pipeline: **Parse -> Build Tree -> Hydrate -> Report**. Your work will be in the very first stage: **Parsing**.

### The Goal of a Parser

A parser has one simple but crucial job: **Translate a specific report format into a standardized Go struct.**

That's it. A parser doesn't need to know about other reports, how to merge them, or how to analyze source code. Its only responsibility is to read a file like `jacoco.xml` or `my_custom_format.txt` and produce a `parsers.ParserResult` object.

### A Parser's Core Responsibilities

When you build a new parser, you'll need to implement the `parsers.IParser` interface. This involves handling a few key tasks:

1.  **Detecting the File Type:** Your parser needs a `SupportsFile()` method that can quickly and efficiently determine if a given file is in the format it understands (e.g., by checking the file extension or reading the first few lines).

2.  **Extracting File Paths:** The parser must read the report and identify the relative paths to the source code files that were covered (e.g., `src/main/java/com/mycompany/App.java`).

3.  **Parsing Line Coverage:** For each source file, the parser must extract the line-by-line coverage data. The output should be a map where the key is the line number and the value is a `model.LineMetrics` struct containing the hit count.

4.  **Parsing Branch Coverage (If Available):** If the report format supports branch coverage (like Cobertura), the parser should also extract the total number of branches on a line and how many of those were covered. If not, these values should be `0`.

5.  **Handling Unresolved Files:** The parser itself doesn't resolve source files against the user's `-sourcedirs`. However, the `processing` logic that calls your parser *will* try to find the files. Your parser should collect a list of any source files mentioned in the report that could not be found on disk and return them in the `UnresolvedSourceFiles` field of the `ParserResult`.

Your final output for a single report file will be a `*parsers.ParserResult` struct, which is essentially a collection of coverage data grouped by the source files mentioned in that report.

### The Parser Testing Strategy: Our Secret to Success

To make development fast and reliable, we use **hermetic unit tests**. This means every test is self-contained and doesn't rely on the actual file system.

*   **Why We Do This:**
    *   **Speed:** Tests run in milliseconds without slow disk I/O.
    *   **Reliability:** Tests are predictable and won't fail just because a file was moved or you're on a different OS.
    *   **Clarity:** Everything needed to understand a test—the input report content, the mock source files, and the expected outcome—is in one place.
    *   **Edge Cases:** It's incredibly easy to test for malformed reports or other error conditions.

#### The Test Pattern

We use a simple, table-driven test pattern. For each scenario you want to test, you'll define:

1.  `name`: A descriptive name for the test case.
2.  `reportContent`: A string containing the exact text of the report file you're testing.
3.  `sourceFiles`: A map simulating the file system, where keys are file paths and values are their content.
4.  `asserter`: A function that receives the result of your parser and contains all the `assert` and `require` calls to verify that the output is correct.

By following this pattern, you can build a robust suite of tests that fully validates your new parser, making the development process a pleasure instead of a chore.

We're excited to see what you build! If you have any questions, feel free to open an issue on the project's GitHub page. Happy coding