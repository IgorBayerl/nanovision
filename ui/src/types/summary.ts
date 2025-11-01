export type RiskLevel = 'safe' | 'warning' | 'danger'

export interface CoverageDetail {
    covered: number
    uncovered?: number
    coverable?: number
    total: number
    percentage: number
}

export type Metrics = Record<string, CoverageDetail>

export type Statuses = Partial<Record<string, RiskLevel>>

export interface FileNode {
    id: string
    name: string
    type: 'file' | 'folder'
    path: string
    children?: FileNode[]
    metrics?: Metrics
    statuses?: Statuses
    componentId?: string
    componentName?: string
    targetUrl?: string | null
}

export interface Totals {
    files: number
    folders: number
    statuses?: Statuses
    [key: string]: CoverageDetail | number | Statuses | undefined
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

export type MetricDefinition = {
    label: string
    shortLabel?: string
    subMetrics: SubMetric[]
}

export type MetricDefinitions = Record<string, MetricDefinition>

export interface SummaryV1 {
    schemaVersion: number
    generatedAt: string
    reportId?: string
    title: string
    totals: Totals
    tree: FileNode[]
    metricDefinitions: MetricDefinitions
    metadata?: MetadataItem[]
}

export type RiskFilter = 'all' | 'danger' | 'warning' | 'safe'
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

export type DiffStatus = 'added' | 'removed' | 'unchanged'

export interface Report {
    name: string
    path: string
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
}
