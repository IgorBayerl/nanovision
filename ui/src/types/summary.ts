export type RiskLevel = 'safe' | 'warning' | 'danger'

export interface CoverageDetail {
    covered: number
    uncovered?: number
    coverable?: number
    total: number
    percentage: number
}

export type Metrics = Record<string, CoverageDetail | ScoreDetail>

export type Statuses = Partial<Record<string, RiskLevel>>

export type DiffStatus = 'added' | 'modified' | 'unchanged' | 'removed'

export interface FileNode {
    id: string
    name: string
    type: 'file' | 'folder'
    path: string
    /** Parent node id; undefined/"" for top-level nodes. */
    parentId?: string
    /** Structural depth from the backend; root-level nodes are 0. */
    depth?: number
    metrics?: Metrics
    statuses?: Statuses
    componentId?: string
    componentName?: string
    targetUrl?: string | null
    diffStatus?: DiffStatus
}

export interface Totals {
    files: number
    folders: number
    statuses?: Statuses
    line_coverage?: CoverageDetail
    branch_coverage?: CoverageDetail
    methods_covered?: CoverageDetail
    methods_fully_covered?: CoverageDetail
    method_branch_coverage?: CoverageDetail
    max_cyclomatic_complexity?: ScoreDetail
    [key: string]: CoverageDetail | ScoreDetail | number | Statuses | undefined
}

export type MetadataItem = {
    label: string
    value: string | string[]
    sizeHint?: 'small' | 'medium' | 'large'
}

export type SubMetric = {
    id: keyof CoverageDetail | string
    label: string
    width: number
}

/** How a metric renders: a percentage card/column, or a plain scalar value. */
export type MetricKind = 'percentage' | 'value'

export type MetricDefinition = {
    label: string
    shortLabel?: string
    /** One-line explanation of the metric, shown in the UI's hover tooltips. */
    description?: string
    kind?: MetricKind
    subMetrics: SubMetric[]
}

/** Scalar (non-percentage) metric payload, e.g. max cyclomatic complexity. */
export interface ScoreDetail {
    value: number
}

export type MetricDefinitions = Record<string, MetricDefinition>

/** One evaluated review-gate threshold (e.g. minimum patch coverage). */
export interface ReviewGateCheck {
    key: string
    label: string
    value: number
    threshold: number
    passed: boolean
}

/** A changed method ranked by exposed risk (complexity × uncovered ratio). */
export interface ReviewHotspot {
    file: string
    method: string
    startLine: number
    diffStatus: string
    complexity?: number
    patchCoverage?: number
    risk: number
}

export interface ReviewStats {
    changedFiles: number
    methodsAdded: number
    methodsModified: number
    untestedChangedMethods: number
    patchStatementsValid: number
    patchStatementsCovered: number
    maxChangedComplexity: number
}

/** Changelist evaluation emitted by the HtmlReview report type. */
export interface ReviewResult {
    passed: boolean
    checks?: ReviewGateCheck[]
    stats: ReviewStats
    hotspots?: ReviewHotspot[]
}

export interface SummaryV1 {
    schemaVersion: number
    generatedAt: string
    reportId?: string
    title: string
    totals: Totals
    /** Flat, pre-ordered (depth-first) list of all file and folder nodes. */
    nodes: FileNode[]
    metricDefinitions: MetricDefinitions
    metadata?: MetadataItem[]
    /** Raw URL query string auto-applied on first load (no existing query). */
    defaultFilters?: string
    /** Present only in review reports; switches the UI to the review layout. */
    review?: ReviewResult
    /** Every parsed report, indexed as the coverage masks address them. */
    reports?: Report[]
    /** File path -> compressed per-report coverage. */
    reportIndexes?: Record<string, ReportIndex>
    statusBands?: Record<string, StatusBand>
}

export type RiskFilter = 'all' | 'danger' | 'warning' | 'safe'
export type DiffFilter = 'all' | 'changed'
export type MetricKey = string

export type MetricConfig = {
    id: MetricKey
    label: string
    shortLabel: string
    enabled: boolean
    definition: MetricDefinition
}

export type FilterRange = {
    min: number
    max: number
}

export type SortableSubMetricKey = keyof CoverageDetail | string

export type SortKey = 'name' | { metric: MetricKey; subMetric: SortableSubMetricKey }
export type SortDir = 'asc' | 'desc'

export type LineStatus = 'covered' | 'uncovered' | 'not-coverable' | 'partial'

export interface Report {
    name: string
    path: string
    /** Details pages only: this report contributed a hit to the file. */
    relevant?: boolean
}

/**
 * One group of entities (lines, statements or methods) sharing a coverage
 * requirement. They are covered when every report bitmask in `m` intersects the
 * active selection. See lib/reportSelection.ts.
 */
export interface ReportBucket {
    m: number[]
    n: number
}

/** Compressed per-report coverage for one file, keyed by metric. */
export type ReportIndex = Record<string, ReportBucket[]>

/** Risk thresholds, so the UI can re-classify recomputed numbers. */
export interface StatusBand {
    min: number
    max: number
}

export interface LineDetails {
    lineNumber: number
    content: string
    status: LineStatus
    hits?: number[]
    branchInfo?: {
        covered: number
        total: number
    }
    diffStatus?: DiffStatus
}

export type MethodMetric = {
    value: string // The value to display (e.g., "100%", "1.00")
    status?: RiskLevel
}

export interface Method {
    name: string
    startLine: number
    endLine: number
    metrics: Record<string, MethodMetric>
    diffStatus?: DiffStatus
}

export interface DetailsV1 {
    schemaVersion: 1
    generatedAt: string
    title: string
    fileName: string
    totals: Totals
    metricDefinitions: MetricDefinitions
    lines: LineDetails[]
    metadata?: MetadataItem[]
    methods?: Method[]
    reports?: Report[]
    /** Compressed per-report coverage for this file, keyed by metric. */
    reportIndex?: ReportIndex
    statusBands?: Record<string, StatusBand>
    /** Raw URL query string auto-applied on first load (no existing query). */
    defaultFilters?: string
}
