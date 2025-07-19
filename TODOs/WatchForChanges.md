# Feature Plan: Live Reloading with Watch Mode

This document outlines the implementation plan for a `--watch` flag. This feature will enable a live-reloading development workflow, where the tool automatically monitors input coverage reports for changes and regenerates the output reports immediately, providing instant feedback to developers.

## 1. Feature Overview

The core of this feature is to introduce a new command-line flag, `--watch`, that shifts the application from a single-run execution model to a persistent, file-watching process. When active, the tool will:

1.  Perform an initial full run of the parse, merge, and report pipeline.
2.  Identify all unique directories containing the input report files (resolved from the `--report` glob patterns).
3.  Begin monitoring these directories for any file changes (creations, writes, or deletions).
4.  Upon detecting a change, it will automatically re-trigger the entire pipeline to generate a fresh set of reports.

This creates a "set it and forget it" experience for developers during a coding and testing session.

## 2. User Impact & Benefits

*   **Instant Feedback Loop**: Developers no longer need to manually re-run the report generator after their tests complete. Changes to coverage are reflected in the output reports in real-time, dramatically speeding up the test-code-review cycle.
*   **Enhanced IDE Integration**: For formats like LCOV, which are consumed by IDE extensions (e.g., in VS Code), this feature creates a dynamic development environment. As a developer saves a file and tests re-run in the background, the coverage highlighting in their editor will update automatically, providing a seamless and interactive experience.
*   **Improved Developer Workflow**: This feature removes a tedious manual step, allowing developers to stay focused on writing code and tests. It turns the report generator from an explicit command into a background service that continuously provides value.
*   **Efficient Debugging of Coverage**: When trying to understand why a specific line isn't covered, a developer can make a change to a test, save it, and see the impact on the report instantly without ever leaving their editor or running another command.

## 3. High-Level Architectural Approach

The implementation will be integrated into the existing application entrypoint (`cmd/main.go`) by conditionally altering the program flow after the initial setup.

1.  **Dependency**: A robust, cross-platform file-watching library will be introduced. The industry standard in the Go ecosystem is `fsnotify`.
2.  **Flag & Control Flow**: A new `--watch` flag will be added. The `run()` function in `cmd/main.go` will be modified. If the flag is present, after the initial report generation, the application will enter a persistent watch loop instead of exiting.
3.  **Refactoring for Reusability**: The core logic for parsing, merging, and reporting, which currently resides within the `run()` function, will be extracted into a separate, reusable function (e.g., `executePipeline()`). This allows both the initial run and subsequent re-runs triggered by the watcher to call the exact same logic, ensuring consistency and avoiding code duplication.
4.  **Watcher Initialization**: The watcher will be configured to monitor the directories of the files resolved from the `--report` glob patterns. This is a crucial detail: to handle new files matching a glob, we must watch directories, not just the initial set of files.
5.  **Debouncing**: To prevent redundant executions from rapid file writes (common with some testing tools), a simple debouncing mechanism will be implemented. After a change is detected, the system will wait for a short, quiet period before triggering the regeneration.

## 4. Detailed Implementation Plan

### Phase 1: Core Setup and Logic Extraction

**Goal**: Add the necessary dependency, introduce the `--watch` flag, and refactor the main application flow to be reusable.

1.  **Add `fsnotify` Dependency:**
    *   **Action:** Add the file-watching library to the project.
        ```bash
        go get github.com/fsnotify/fsnotify
        ```

2.  **Add the `--watch` Flag:**
    *   **Action:** Modify `cmd/main.go`.
    *   **Change:** Add the new flag to the `cliFlags` struct and the `parseFlags()` function.
        ```go
        // in: cmd/main.go -> cliFlags struct
        type cliFlags struct {
            // ... existing flags ...
            watch             *bool
        }

        // in: cmd/main.go -> parseFlags()
        func parseFlags() (*cliFlags, error) {
            f := &cliFlags{
                // ... existing flag definitions ...
                watch: flag.Bool("watch", false, "Enable watch mode to automatically regenerate reports on file changes"),
            }
            // ...
            return f, nil
        }
        ```

3.  **Refactor Main Logic into a Reusable Function:**
    *   **Action:** Modify `cmd/main.go`.
    *   **Change:** Create a new function `executePipeline` that encapsulates the core work of a single run. The existing logic from `run()` will be moved here.

        ```go
        // in: cmd/main.go
        
        // New function to encapsulate a single report generation run
        func executePipeline(logger *slog.Logger, flags *cliFlags, langFactory *language.ProcessorFactory, parserFactory *parsers.ParserFactory) error {
            actualReportFiles, invalidPatterns, err := resolveAndValidateInputs(logger, flags)
            if err != nil {
                // In watch mode, a pattern not matching any files isn't a fatal error, just a state.
                if *flags.watch {
                    logger.Warn("No report files found matching patterns, waiting for changes...", "patterns", *flags.reportsPatterns)
                    return nil // Return nil to allow the watcher to continue.
                }
                return fmt.Errorf("no valid report files found: %w", err)
            }

            reportConfig, err := createReportConfiguration(flags, logging.Info, actualReportFiles, invalidPatterns, langFactory, logger) // Assuming verbosity is handled
            if err != nil {
                return err
            }

            summaryResult, err := parseAndMergeReports(logger, reportConfig, parserFactory)
            if err != nil {
                // Also not fatal in watch mode
                logger.Error("Failed to parse and merge reports", "error", err)
                return nil 
            }
            
            logger.Info("Report data processed, starting report generation...")
            reportCtx := reporter.NewBuilderContext(reportConfig, settings.NewSettings(), logger)
            err = generateReports(reportCtx, summaryResult)
            if err == nil {
                logger.Info("Successfully generated reports.", "outputDir", reportConfig.TargetDirectory())
            }
            return err
        }

        // The existing run() function will be modified in the next phase.
        ```

### Phase 2: Implementing the Watch Loop

**Goal**: Modify the `run()` function to perform the initial run and then enter the file-watching loop if `--watch` is enabled.

1.  **Modify the `run()` Function:**
    *   **Action:** Update `cmd/main.go`.
    *   **Change:** The `run()` function will now be the main controller, deciding whether to run once or to start the watch loop.

    ```go
    // in: cmd/main.go
    
    // (Import fsnotify)
    import (
        // ...
        "github.com/fsnotify/fsnotify"
    )

    func run() error {
        // --- (Initial setup from existing run() function remains the same) ---
        flags, err := parseFlags()
        // ... buildLogger, create factories ...
        
        // --- Initial Run ---
        logger.Info("Performing initial report generation...")
        if err := executePipeline(logger, flags, langFactory, parserFactory); err != nil {
            // A hard error on the first run should still exit the program.
            return fmt.Errorf("initial run failed: %w", err)
        }

        // --- Conditional Watch Loop ---
        if !*flags.watch {
            return nil // Exit after single run if not in watch mode
        }

        logger.Info("Watch mode enabled. Waiting for file changes...", "patterns", *flags.reportsPatterns)

        watcher, err := fsnotify.NewWatcher()
        if err != nil {
            return fmt.Errorf("failed to create file watcher: %w", err)
        }
        defer watcher.Close()

        // Goroutine to handle events
        go func() {
            var timer *time.Timer
            for {
                select {
                case event, ok := <-watcher.Events:
                    if !ok { return }
                    // We only care about writes, creations, and removals that could affect glob results
                    if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
                        logger.Debug("Change detected", "file", event.Name, "op", event.Op)
                        // Debounce: Reset the timer on each new event
                        if timer != nil {
                            timer.Stop()
                        }
                        timer = time.AfterFunc(500*time.Millisecond, func() {
                            logger.Info("File change detected. Regenerating report...")
                            if err := executePipeline(logger, flags, langFactory, parserFactory); err != nil {
                                logger.Error("Error during report regeneration", "error", err)
                            }
                        })
                    }
                case err, ok := <-watcher.Errors:
                    if !ok { return }
                    logger.Error("File watcher error", "error", err)
                }
            }
        }()

        // --- Add Directories to Watcher ---
        // We watch the directories of the glob patterns to detect new files.
        reportPatterns := strings.Split(*flags.reportsPatterns, ";")
        watchedDirs := make(map[string]struct{})
        for _, pattern := range reportPatterns {
            // This is a simplification. A robust implementation would walk up from the first wildcard.
            // For `src/**/*.xml`, it would watch `src`. For `/tmp/reports/*.xml`, it would watch `/tmp/reports`.
            dir := filepath.Dir(pattern)
            if strings.Contains(pattern, "**") {
                dir = strings.Split(pattern, "**")[0]
            }
             if _, err := os.Stat(dir); err == nil {
                if _, watched := watchedDirs[dir]; !watched {
                    logger.Info("Watching directory for changes", "dir", dir)
                    if err := watcher.Add(dir); err != nil {
                        return fmt.Errorf("failed to watch directory %s: %w", dir, err)
                    }
                    watchedDirs[dir] = struct{}{}
                }
            }
        }
        if len(watchedDirs) == 0 {
            logger.Warn("Could not determine valid directories to watch from the provided patterns.")
        }

        // Block main goroutine until interrupted
        <-make(chan struct{})

        return nil
    }
    ```

### Phase 3: User Experience Enhancements

**Goal**: Make the watch mode more user-friendly and informative.

1.  **Clear Console on Re-run (Optional but Recommended):**
    *   **Action:** Add a utility function to clear the console screen.
    *   **Change:** Before calling `executePipeline` inside the debounced timer, add a call to clear the screen. This makes each new report generation feel clean and intentional.

    ```go
    // in internal/utils/term.go (new file)
    package utils

    import (
        "os"
        "os/exec"
        "runtime"
    )

    var clear map[string]func() //create a map for storing clear funcs

    func init() {
        clear = make(map[string]func()) //Initialize it
        clear["linux"] = func() { 
            cmd := exec.Command("clear") //Linux example, its tested
            cmd.Stdout = os.Stdout
            cmd.Run()
        }
        clear["windows"] = func() {
            cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested 
            cmd.Stdout = os.Stdout
            cmd.Run()
        }
    }

    func CallClear() {
        value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
        if ok { //if we defined a clear func for that platform:
            value()  //we execute it
        }
    }
    ```

2.  **Graceful Shutdown:**
    *   **Action:** Modify the watch loop in `cmd/main.go` to listen for an interrupt signal (Ctrl+C).
    *   **Change:** Use `os.Signal` to gracefully shut down the watcher and exit.

    ```go
    // in: cmd/main.go -> run()
    
    // ... after watcher setup ...

    // Use a channel to listen for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    
    logger.Info("Press Ctrl+C to exit watch mode.")

    // Block until a signal is received
    <-quit
    logger.Info("Shutdown signal received, exiting.")
    
    // Replace the empty channel block: <-make(chan struct{})
    ```

## 5. Future Possibilities & Enhancements

*   **Performance Optimization (Intelligent Re-parsing)**: The current plan re-runs the entire pipeline. A more advanced version could cache the `ParserResult` of each input file. When a single file changes, the system would only re-parse that one file and then re-run the `analyzer.MergeParserResults` with the updated set of cached results. This would provide a near-instantaneous update for projects with many large, unchanged report files.
*   **Configuration File Watching**: Extend the watch functionality to monitor a project-specific configuration file (e.g., `.reportgenerator.yml`). If the configuration changes (e.g., a filter is updated), the tool would automatically reload the settings and regenerate the report.
*   **Browser Live-Reload**: Integrate a lightweight web server and WebSocket connection to automatically refresh the HTML report in the browser whenever it's regenerated, completing the live feedback loop.