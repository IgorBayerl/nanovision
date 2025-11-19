# Unified Diff Parser

This package implements a parser for unified diff format (git/perforce style).

## Supported Syntax

- Git-style headers:
  - `diff --git a/oldpath b/newpath`
  - `--- oldpath` or `--- /dev/null` (new files)
  - `+++ newpath`
  - `new file mode XXXX`
  
- Unified diff hunks:
  - `@@ -oldStart,oldLen +newStart,newLen @@`
  - Context lines (space-prefixed)
  - Removed lines (`-` prefix)
  - Added lines (`+` prefix)

## Line Classification

Added lines (`+`) are classified as either:
- **Modified**: when they immediately follow removed lines (`-`) in the same hunk
- **Added**: when they appear without preceding removed lines

## Limitations

- Only tracks new-file line numbers (post-change state)
- Treats renames as modifications (uses NewPath only)
- Requires valid hunk headers
- No binary file support
- Expects UTF-8 encoded text files