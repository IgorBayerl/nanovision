import { z } from 'zod'

const riskLevelSchema = z.enum(['safe', 'warning', 'danger'])

// A schema for the Statuses object, allowing any string key with a RiskLevel value
const statusesSchema = z.record(z.string(), riskLevelSchema).optional()

// A base schema for individual coverage metrics.
const coverageDetailSchema = z.object({
    covered: z.number(),
    uncovered: z.number().optional(),
    coverable: z.number().optional(),
    total: z.number(),
    percentage: z.number(),
})

// Scalar (non-percentage) metric payload, e.g. max cyclomatic complexity.
const scoreDetailSchema = z.object({
    value: z.number(),
})

// A schema for the Metrics object, which can contain coverage details or scalar scores.
const metricsSchema = z.record(z.string(), coverageDetailSchema.or(scoreDetailSchema))

const diffStatusSchema = z.enum(['added', 'modified', 'unchanged', 'removed'])

// Compressed per-report coverage. A bucket holds every entity (line, statement
// or method) sharing one coverage requirement: `n` of them are covered when all
// of the report bitmasks in `m` intersect the active selection.
const reportBucketSchema = z.object({
    m: z.array(z.number().int().nonnegative()),
    n: z.number().int().nonnegative(),
})

const reportIndexSchema = z.record(z.string(), z.array(reportBucketSchema))

const statusBandSchema = z.object({
    min: z.number(),
    max: z.number(),
})

/** The warning ranges a folder sets for itself, over the ones of the report. */
const folderBandsSchema = z.object({
    path: z.string(),
    bands: z.record(z.string(), statusBandSchema),
})

/** One source of settings: the root config file (path ""), or the settings of a folder. */
const configFileSchema = z.object({
    path: z.string(),
    source: z.string(),
})

const reportSchema = z.object({
    name: z.string(),
    path: z.string(),
    // details pages only: this report contributed a hit to the file
    relevant: z.boolean().optional(),
})

// A single file or folder node in the flat node list.
export type FileNode = {
    id: string
    name: string
    type: 'file' | 'folder'
    path: string
    parentId?: string
    depth?: number
    metrics?: z.infer<typeof metricsSchema>
    statuses?: z.infer<typeof statusesSchema>
    targetUrl?: string | null
    diffStatus?: z.infer<typeof diffStatusSchema>
    config?: boolean
}

const fileNodeSchema: z.ZodType<FileNode> = z.object({
    id: z.string().min(1, 'Node ID cannot be empty.'),
    name: z.string().min(1, 'Node name cannot be empty.'),
    type: z.enum(['file', 'folder']),
    path: z.string().min(1, 'Node path cannot be empty.'),
    parentId: z.string().optional(),
    depth: z.number().optional(),
    metrics: metricsSchema.optional(),
    statuses: statusesSchema,
    targetUrl: z.string().nullable().optional(),
    diffStatus: diffStatusSchema.optional(),
    /** A config file of the run, listed in its folder. It has no metrics. */
    config: z.boolean().optional(),
})

// A schema for the overall totals section
const totalsSchema = z
    .object({
        files: z.number(),
        folders: z.number(),
        statuses: statusesSchema,
        // Allows other keys to be present, as long as they are valid coverage details or scores
    })
    .catchall(coverageDetailSchema.or(scoreDetailSchema).or(z.number()))

// A schema for a single metadata item
const metadataItemSchema = z.object({
    label: z.string(),
    value: z.union([z.string(), z.array(z.string())]),
    sizeHint: z.enum(['small', 'medium', 'large']).optional(),
})

// Schemas for defining how metrics should be displayed
const subMetricSchema = z.object({
    id: z.string(),
    label: z.string(),
    width: z.number(),
})

const metricDefinitionSchema = z.object({
    label: z.string(),
    shortLabel: z.string().optional(),
    /** One-line explanation of the metric, shown in the UI's hover tooltips. */
    description: z.string().optional(),
    kind: z.enum(['percentage', 'value']).optional(),
    subMetrics: z.array(subMetricSchema),
})

// gate verdict and changelist stats of a run measured with a diff;
// the Changes tab shows them.
const reviewGateCheckSchema = z.object({
    key: z.string(),
    label: z.string(),
    value: z.number(),
    threshold: z.number(),
    passed: z.boolean(),
})

const reviewStatsSchema = z.object({
    changedFiles: z.number(),
    methodsAdded: z.number(),
    methodsModified: z.number(),
    untestedChangedMethods: z.number(),
    patchStatementsValid: z.number(),
    patchStatementsCovered: z.number(),
    maxChangedComplexity: z.number(),
})

const reviewSchema = z.object({
    passed: z.boolean(),
    checks: z.array(reviewGateCheckSchema).optional(),
    stats: reviewStatsSchema,
})

export type Review = z.infer<typeof reviewSchema>

// the delta of a run against its base run
const runRefSchema = z.object({
    id: z.number().optional(),
    revision: z.string().optional(),
})

const metricDeltaSchema = z.object({
    key: z.string(),
    label: z.string(),
    /** What the counts count, e.g. "statements". */
    unit: z.string(),
    baseCovered: z.number(),
    baseTotal: z.number(),
    currentCovered: z.number(),
    currentTotal: z.number(),
    base: z.number(),
    current: z.number(),
    /** Current minus base, in percentage points. */
    delta: z.number(),
})

const fileDeltaSchema = z.object({
    path: z.string(),
    base: z.number(),
    current: z.number(),
    delta: z.number(),
})

const changeSetSchema = z.object({
    measured: z.array(z.string()).optional(),
    ignored: z.array(z.string()).optional(),
    notInReports: z.array(z.string()).optional(),
    deleted: z.array(z.string()).optional(),
})

const comparisonSchema = z.object({
    base: runRefSchema,
    /** The base run is at the base revision itself. */
    exact: z.boolean(),
    /** Revisions between the base revision and the base run, -1 when unknown. */
    distance: z.number(),
    /** The key of the metric that leads. */
    headline: z.string(),
    metrics: z.array(metricDeltaSchema),
    /** Files whose headline metric changed, largest change first. */
    files: z.array(fileDeltaSchema).optional(),
    /** Path of a file or folder -> metric -> change in percentage points. */
    deltas: z.record(z.string(), z.record(z.string(), z.number())).optional(),
    change: changeSetSchema.optional(),
    warnings: z.array(z.string()).optional(),
})

export type Comparison = z.infer<typeof comparisonSchema>
export type FileDelta = z.infer<typeof fileDeltaSchema>
export type ChangeSet = z.infer<typeof changeSetSchema>

export type ReportBucket = z.infer<typeof reportBucketSchema>
export type ReportIndex = z.infer<typeof reportIndexSchema>
export type StatusBand = z.infer<typeof statusBandSchema>
export type FolderBands = z.infer<typeof folderBandsSchema>
export type ConfigFile = z.infer<typeof configFileSchema>
export type ReportRef = z.infer<typeof reportSchema>

export const summaryV1Schema = z.object({
    schemaVersion: z.literal(1, { message: 'This report requires schemaVersion 1.' }),
    generatedAt: z
        .string()
        .refine((val) => !Number.isNaN(Date.parse(val)), { message: 'GeneratedAt must be a valid date string.' }),
    reportId: z.string().optional(),
    title: z.string().min(1, { message: 'Report title is missing or empty.' }),
    totals: totalsSchema,
    nodes: z.array(fileNodeSchema),
    metricDefinitions: z.record(z.string(), metricDefinitionSchema),
    metricOrder: z.array(z.string()).optional(),
    metadata: z.array(metadataItemSchema).optional(),
    defaultFilters: z.string().optional(),
    review: reviewSchema.optional(),
    comparison: comparisonSchema.optional(),
    /** The revisions being compared, as rows for the side panel. */
    comparing: z.array(metadataItemSchema).optional(),
    /** The nodes hold only the changed files; the totals cover the whole run. */
    onlyChanged: z.boolean().optional(),
    reports: z.array(reportSchema).optional(),
    reportIndexes: z.record(z.string(), reportIndexSchema).optional(),
    statusBands: z.record(z.string(), statusBandSchema).optional(),
    /** Folders with their own warning ranges, outer folders first. */
    folderBands: z.array(folderBandsSchema).optional(),
    /** The config files of the run; absent when no folder has its own settings. */
    configs: z.array(configFileSchema).optional(),
})

export type SummaryV1 = z.infer<typeof summaryV1Schema>

/**
 * Validates the entire summary data object against the schema.
 * @param data The unknown data, typically from window.__NANOVISION_SUMMARY__.
 * @returns A Zod SafeParseReturnType which indicates success or failure with detailed errors.
 */
export function validateSummaryData(data: unknown) {
    return summaryV1Schema.safeParse(data)
}

// Schema for a single line of code
const lineStatusSchema = z.enum(['covered', 'uncovered', 'not-coverable'])

const lineDetailsSchema = z.object({
    lineNumber: z.number().int().positive(),
    content: z.string(),
    status: lineStatusSchema,
    hits: z.array(z.number().int()).optional(),
    diffStatus: diffStatusSchema.optional(),
})

// A schema for a method's metrics
const methodMetricSchema = z.object({
    value: z.string(),
    status: riskLevelSchema.optional(),
})

// A schema for a single method/function in the file
const methodSchema = z.object({
    name: z.string(),
    startLine: z.number(),
    endLine: z.number(),
    metrics: z.record(z.string(), methodMetricSchema),
    diffStatus: diffStatusSchema.optional(),
})

// Schema for the entire details page data object
export const detailsV1Schema = z.object({
    schemaVersion: z.literal(1),
    generatedAt: z.string().refine((val) => !Number.isNaN(Date.parse(val))),
    title: z.string(),
    fileName: z.string(),
    totals: totalsSchema,
    metricDefinitions: z.record(z.string(), metricDefinitionSchema),
    metricOrder: z.array(z.string()).optional(),
    lines: z.array(lineDetailsSchema),
    metadata: z.array(metadataItemSchema).optional(),
    methods: z.array(methodSchema).optional(),
    reports: z.array(reportSchema).optional(),
    reportIndex: reportIndexSchema.optional(),
    statusBands: z.record(z.string(), statusBandSchema).optional(),
    defaultFilters: z.string().optional(),
})

export type DetailsV1 = z.infer<typeof detailsV1Schema>

/**
 * Validates the details page data object against the schema.
 * @param data The unknown data, typically from window.__NANOVISION_DETAILS__.
 * @returns A Zod SafeParseReturnType which indicates success or failure with detailed errors.
 */
export function validateDetailsData(data: unknown) {
    return detailsV1Schema.safeParse(data)
}
