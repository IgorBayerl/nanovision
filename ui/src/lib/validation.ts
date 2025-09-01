import { z } from 'zod'

const riskLevelSchema = z.enum(['safe', 'warning', 'danger'])

// A schema for the Statuses object, allowing any string key with a RiskLevel value
const statusesSchema = z.record(z.string(), riskLevelSchema).optional()

// A base schema for individual coverage metrics.
// for metrics like `branchCoverage` which may not include them.
const coverageDetailSchema = z.object({
    covered: z.number(),
    uncovered: z.number().optional(),
    coverable: z.number().optional(),
    total: z.number(),
    percentage: z.number(),
})

// A schema for the Metrics object, which can contain any number of coverage details
const metricsSchema = z.record(z.string(), coverageDetailSchema)

// A recursive schema for a file or folder node in the tree.
export type FileNode = {
    id: string
    name: string
    type: 'file' | 'folder'
    path: string
    children?: FileNode[]
    metrics?: z.infer<typeof metricsSchema>
    statuses?: z.infer<typeof statusesSchema>
    targetUrl?: string | null
}

const fileNodeSchema: z.ZodType<FileNode> = z.lazy(() =>
    z.object({
        id: z.string().min(1, 'Node ID cannot be empty.'),
        name: z.string().min(1, 'Node name cannot be empty.'),
        type: z.enum(['file', 'folder']),
        path: z.string().min(1, 'Node path cannot be empty.'),
        children: z.array(fileNodeSchema).optional(),
        metrics: metricsSchema.optional(),
        statuses: statusesSchema,
        targetUrl: z.string().nullable().optional(),
    }),
)

// A schema for the overall totals section
const totalsSchema = z
    .object({
        files: z.number(),
        folders: z.number(),
        statuses: statusesSchema,
        // Allows other keys to be present, as long as they are valid coverage details
    })
    .catchall(coverageDetailSchema.or(z.number()))

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
    subMetrics: z.array(subMetricSchema),
})

export const summaryV1Schema = z.object({
    schemaVersion: z.literal(1, { message: 'This report requires schemaVersion 1.' }),
    generatedAt: z
        .string()
        .refine((val) => !Number.isNaN(Date.parse(val)), { message: 'GeneratedAt must be a valid date string.' }),
    reportId: z.string().optional(),
    title: z.string().min(1, { message: 'Report title is missing or empty.' }),
    totals: totalsSchema,
    tree: z.array(fileNodeSchema),
    metricDefinitions: z.record(z.string(), metricDefinitionSchema),
    metadata: z.array(metadataItemSchema).optional(),
})

export type SummaryV1 = z.infer<typeof summaryV1Schema>

/**
 * Validates the entire summary data object against the schema.
 * @param data The unknown data, typically from window.__ADLERCOV_SUMMARY__.
 * @returns A Zod SafeParseReturnType which indicates success or failure with detailed errors.
 */
export function validateSummaryData(data: unknown) {
    return summaryV1Schema.safeParse(data)
}

// Schema for a single line of code
const lineStatusSchema = z.enum(['covered', 'uncovered', 'not-coverable', 'partial'])

const lineDetailsSchema = z.object({
    lineNumber: z.number().int().positive(),
    content: z.string(),
    status: lineStatusSchema,
    hits: z.number().int().optional(),
    branchInfo: z
        .object({
            covered: z.number().int(),
            total: z.number().int(),
        })
        .optional(),
})

// A schema for a metric that can be a number OR a full coverage object
const methodMetricValueSchema = z.union([coverageDetailSchema, z.number()])

// A schema for a method's metrics
const methodMetricsSchema = z.record(z.string(), methodMetricValueSchema)

// A schema for a single method/function in the file
const methodSchema = z.object({
    name: z.string(),
    startLine: z.number(),
    endLine: z.number(),
    metrics: methodMetricsSchema,
    statuses: statusesSchema,
})

// Schema for the entire details page data object
export const detailsV1Schema = z.object({
    schemaVersion: z.literal(1),
    generatedAt: z.string().refine((val) => !Number.isNaN(Date.parse(val))),
    title: z.string(),
    fileName: z.string(),
    totals: totalsSchema,
    metricDefinitions: z.record(z.string(), metricDefinitionSchema),
    lines: z.array(lineDetailsSchema),
    metadata: z.array(metadataItemSchema).optional(),
    methods: z.array(methodSchema).optional(),
})

export type DetailsV1 = z.infer<typeof detailsV1Schema>

/**
 * Validates the details page data object against the schema.
 * @param data The unknown data, typically from window.__ADLERCOV_DETAILS__.
 * @returns A Zod SafeParseReturnType which indicates success or failure with detailed errors.
 */
export function validateDetailsData(data: unknown) {
    return detailsV1Schema.safeParse(data)
}
