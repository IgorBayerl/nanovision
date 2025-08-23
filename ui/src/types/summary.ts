export type RiskLevel = 'safe' | 'warning' | 'danger'

export interface CoverageDetail {
    covered: number
    uncovered: number
    coverable: number
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
    targetUrl?: string
}

export interface Totals {
    files: number
    folders: number
    statuses?: Statuses // Added statuses to totals
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
    metricDefinitions: MetricDefinitions // Added metric definitions
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
