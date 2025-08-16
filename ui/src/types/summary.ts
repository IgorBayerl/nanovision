export type RiskLevel = 'safe' | 'warning' | 'danger'

export interface Metrics {
    lineCoverage: number
    branchCoverage: number
    methodCoverage: number
    statementCoverage: number
    functionCoverage: number
}
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

export interface Totals {
    files: number
    folders: number
    linesCoveredPct: number
    branchesCoveredPct: number
    methodsCoveredPct: number
    statementsCoveredPct: number
    functionsCoveredPct: number
}

export interface SummaryV1 {
    schemaVersion: number
    generatedAt: string
    reportId?: string
    title: string
    totals: Totals
    tree: FileNode[]
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

export type SortKey = 'name' | MetricKey
export type SortDir = 'asc' | 'desc'
