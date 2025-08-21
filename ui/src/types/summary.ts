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
    // This index signature allows for dynamic metric keys while keeping `files` and `folders` correctly typed.
    [key: string]: CoverageDetail | number
}

export type MetadataItem = {
    label: string
    value: string | string[]
    /** optional: override sizing */
    sizeHint?: 'small' | 'medium' | 'large'
}

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
    metadata?: MetadataItem[]
}

export type RiskFilter = 'all' | 'danger' | 'warning' | 'safe'

export type MetricKey = string

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
