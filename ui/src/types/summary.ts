export type RiskLevel = 'safe' | 'warning' | 'danger'

export interface CoverageDetail {
    covered: number
    uncovered: number
    coverable: number
    total: number
    percentage: number
}

export type Metrics = Record<
    'lineCoverage' | 'branchCoverage' | 'methodCoverage' | 'statementCoverage' | 'functionCoverage',
    CoverageDetail
>

export type Statuses = Partial<Record<keyof Metrics, RiskLevel>>

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

export interface Totals extends Metrics {
    files: number
    folders: number
}

// UPDATED: Added new optional metadata fields.
export interface SummaryV1 {
    schemaVersion: number
    generatedAt: string
    reportId?: string
    title: string
    totals: Totals
    tree: FileNode[]
    parsers?: string[]
    configFiles?: string[]
    importedReports?: string[]
}

export type RiskFilter = 'all' | 'danger' | 'warning' | 'safe'

export type MetricKey = keyof Metrics

export type MetricConfig = {
    id: MetricKey
    label: string
    shortLabel: string
    enabled: boolean
}

export type FilterRange = {
    min: number
    max: number
}

export type SortableSubMetricKey = keyof CoverageDetail

export type SortKey = 'name' | { metric: MetricKey; subMetric: SortableSubMetricKey }
export type SortDir = 'asc' | 'desc'
