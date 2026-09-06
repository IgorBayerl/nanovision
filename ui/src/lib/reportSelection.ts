/**
 * Recomputes coverage numbers for a subset of the parsed reports.
 *
 * The backend ships a compressed index instead of one full copy of the tree per
 * report combination: for every file and metric it groups countable entities
 * (lines, statements, methods) into buckets that share a coverage requirement.
 * A bucket's entities are covered when every report bitmask in it intersects the
 * active selection — one mask for "hit by any of these reports", several for
 * "fully covered", where each part must be reached.
 *
 * Totals never move with the selection. Which lines and methods exist is decided
 * by the union of all reports, which is the same rule the per-line view on the
 * details page already follows.
 */
import type { ReportBucket, ReportIndex, StatusBand } from '@/lib/validation'
import type { CoverageDetail, Metrics, RiskLevel, ScoreDetail, Statuses, Totals } from '@/types/summary'

/** Bitmask of the selected reports; bit i is report i. */
export type ReportSelection = number

export const ALL_REPORTS = -1 >>> 0

export function selectionFromIndices(indices: Iterable<number>): ReportSelection {
    let mask = 0
    for (const index of indices) {
        mask |= 1 << index
    }
    return mask >>> 0
}

export function indicesFromSelection(selection: ReportSelection, reportCount: number): number[] {
    const indices: number[] = []
    for (let i = 0; i < reportCount; i++) {
        if (selection & (1 << i)) indices.push(i)
    }
    return indices
}

export function isReportSelected(selection: ReportSelection, index: number): boolean {
    return (selection & (1 << index)) !== 0
}

export function toggleReport(selection: ReportSelection, index: number): ReportSelection {
    return (selection ^ (1 << index)) >>> 0
}

/** True when the selection covers every report, i.e. the report as generated. */
export function isFullSelection(selection: ReportSelection, reportCount: number): boolean {
    const full = reportCount >= 32 ? ALL_REPORTS : (1 << reportCount) - 1
    return (selection & full) === full
}

export type BucketTally = { covered: number; total: number }

export function resolveBuckets(buckets: ReportBucket[], selection: ReportSelection): BucketTally {
    let covered = 0
    let total = 0

    for (const bucket of buckets) {
        total += bucket.n
        if (bucket.m.every((mask) => (mask & selection) !== 0)) {
            covered += bucket.n
        }
    }

    return { covered, total }
}

// Mirrors utils.CalculatePercentage: nothing to cover counts as fully covered,
// so a file with no changed statements does not read as 0% patch coverage.
function percentage(covered: number, total: number): number {
    if (total <= 0) return 100
    return Math.round((covered / total) * 10000) / 100
}

/**
 * Rebuilds one metric from its buckets, keeping the sub-metric fields the
 * generated value carried. `total` stays as generated because some metrics
 * (line coverage) report physical lines there rather than coverable ones.
 */
function applyTally(previous: CoverageDetail | ScoreDetail | undefined, tally: BucketTally): CoverageDetail {
    const base = previous && 'covered' in previous ? previous : undefined
    const detail: CoverageDetail = {
        ...base,
        covered: tally.covered,
        uncovered: tally.total - tally.covered,
        total: base?.total ?? tally.total,
        percentage: percentage(tally.covered, tally.total),
    }

    if (base?.coverable !== undefined) detail.coverable = tally.total
    return detail
}

function classify(value: number, band: StatusBand | undefined): RiskLevel | undefined {
    if (!band) return undefined
    if (value < band.min) return 'danger'
    if (value <= band.max) return 'warning'
    return 'safe'
}

type Bands = Record<string, StatusBand> | undefined

/** Sums buckets per metric, used to roll files up into folders and totals. */
type TallyMap = Map<string, BucketTally>

function addTallies(into: TallyMap, index: ReportIndex, selection: ReportSelection): void {
    for (const [metricId, buckets] of Object.entries(index)) {
        const tally = resolveBuckets(buckets, selection)
        const running = into.get(metricId)
        if (running) {
            running.covered += tally.covered
            running.total += tally.total
        } else {
            into.set(metricId, { ...tally })
        }
    }
}

function mergeTallies(into: TallyMap, from: TallyMap): void {
    for (const [metricId, tally] of from) {
        const running = into.get(metricId)
        if (running) {
            running.covered += tally.covered
            running.total += tally.total
        } else {
            into.set(metricId, { ...tally })
        }
    }
}

type Measured = { metrics?: Metrics; statuses?: Statuses }

/**
 * Writes recomputed metrics and their risk levels into a copy of `metrics` and
 * `statuses`. Nodes keep both in named fields; the report totals keep the
 * metrics at the top level, so the two callers below differ only in shape.
 */
function applyTallies(
    metrics: Metrics | undefined,
    statuses: Statuses | undefined,
    tallies: TallyMap,
    bands: Bands,
): { metrics: Metrics; statuses: Statuses } {
    const nextMetrics: Metrics = { ...(metrics ?? {}) }
    const nextStatuses: Statuses = { ...(statuses ?? {}) }

    for (const [metricId, tally] of tallies) {
        const detail = applyTally(metrics?.[metricId], tally)
        nextMetrics[metricId] = detail

        const risk = classify(detail.percentage, bands?.[metricId])
        if (risk) {
            nextStatuses[metricId] = risk
        } else {
            delete nextStatuses[metricId]
        }
    }

    return { metrics: nextMetrics, statuses: nextStatuses }
}

function applyTalliesToNode<T extends Measured>(node: T, tallies: TallyMap, bands: Bands): T {
    if (tallies.size === 0) return node
    return { ...node, ...applyTallies(node.metrics, node.statuses, tallies, bands) }
}

function applyTalliesToTotals<T extends Totals>(totals: T, tallies: TallyMap, bands: Bands): T {
    if (tallies.size === 0) return totals

    const { metrics, statuses } = applyTallies(totals as Metrics, totals.statuses, tallies, bands)
    return { ...totals, ...metrics, statuses }
}

/**
 * The node fields this module needs. Both the hand-written FileNode and the
 * schema-inferred one satisfy it, and the caller gets its own type back.
 */
export type SelectableNode = Measured & {
    id: string
    type: 'file' | 'folder'
    path: string
    parentId?: string
}

/**
 * Rebuilds the whole node list and the report totals for a report selection.
 *
 * Files come straight from their own index; folders and the totals are summed
 * from their descendants. The node list is depth-first pre-ordered, so one
 * reverse pass hands every folder the tallies of the subtree below it.
 */
export function applyReportSelection<N extends SelectableNode, T extends Totals>(
    nodes: N[],
    totals: T,
    indexes: Record<string, ReportIndex> | undefined,
    bands: Bands,
    selection: ReportSelection,
): { nodes: N[]; totals: T } {
    if (!indexes) return { nodes, totals }

    const subtree = new Map<string, TallyMap>()
    const rootTallies: TallyMap = new Map()

    // Tallies waiting for the folder that owns them; nodes without a parent
    // roll straight into the report totals.
    const bucketFor = (parentId: string | undefined): TallyMap => {
        if (!parentId) return rootTallies
        let pending = subtree.get(parentId)
        if (!pending) {
            pending = new Map()
            subtree.set(parentId, pending)
        }
        return pending
    }

    const nextNodes = new Array<N>(nodes.length)

    for (let i = nodes.length - 1; i >= 0; i--) {
        const node = nodes[i]

        if (node.type === 'folder') {
            const own = subtree.get(node.id) ?? new Map()
            subtree.delete(node.id)
            nextNodes[i] = applyTalliesToNode(node, own, bands)
            mergeTallies(bucketFor(node.parentId), own)
            continue
        }

        const index = indexes[node.path]
        if (!index) {
            nextNodes[i] = node
            continue
        }

        const own: TallyMap = new Map()
        addTallies(own, index, selection)
        nextNodes[i] = applyTalliesToNode(node, own, bands)
        mergeTallies(bucketFor(node.parentId), own)
    }

    return {
        nodes: nextNodes,
        totals: applyTalliesToTotals(totals, rootTallies, bands),
    }
}

/** The details-page equivalent: one file's totals card block. */
export function applyReportSelectionToTotals<T extends Totals>(
    totals: T,
    index: ReportIndex | undefined,
    bands: Bands,
    selection: ReportSelection,
): T {
    if (!index) return totals

    const tallies: TallyMap = new Map()
    addTallies(tallies, index, selection)
    return applyTalliesToTotals(totals, tallies, bands)
}

/**
 * Metrics the selection cannot move, so the UI can say so instead of lying.
 * A metric counts as filterable when any file indexed it: a single file may
 * miss one simply by having no entities to count, such as patch coverage on an
 * unchanged file.
 */
export function unfilterableMetrics(metricIds: string[], indexes: Record<string, ReportIndex> | undefined): string[] {
    if (!indexes) return []

    const indexed = new Set<string>()
    for (const index of Object.values(indexes)) {
        for (const metricId of Object.keys(index)) indexed.add(metricId)
    }

    return metricIds.filter((id) => !indexed.has(id))
}
