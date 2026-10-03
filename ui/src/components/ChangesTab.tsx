import { AlertTriangle } from 'lucide-react'
import { useMemo } from 'react'
import DiffStatusBadge from '@/components/DiffStatusBadge'
import ReviewSummary from '@/components/ReviewSummary'
import { cn } from '@/lib/utils'
import type { ChangeSet, Comparison, FileDelta, FileNode, Review } from '@/lib/validation'
import type { CoverageDetail } from '@/types/summary'

/** "+0.46" / "-1.14" / "±0.00", in percentage points */
function signed(value: number, digits = 2): string {
    const rounded = Number(value.toFixed(digits))
    if (rounded === 0) return `±${(0).toFixed(digits)}`
    return `${rounded > 0 ? '+' : ''}${rounded.toFixed(digits)}`
}

const signedCount = (n: number) => `${n > 0 ? '+' : ''}${n.toLocaleString('en-US')}`

const deltaTone = (value: number, digits = 2) => {
    const rounded = Number(value.toFixed(digits))
    return rounded > 0 ? 'text-covered' : rounded < 0 ? 'text-uncovered' : 'text-muted-foreground'
}

/** "CL 118432" for Perforce, the first ten characters of a git commit. */
function shortRevision(revision: string): string {
    const at = revision.lastIndexOf('@')
    if (at >= 0 && /^\d+$/.test(revision.slice(at + 1))) return `CL ${revision.slice(at + 1)}`
    return /^[0-9a-f]{12,}$/i.test(revision) ? revision.slice(0, 10) : revision
}

const plural = (n: number, one: string, many: string) => `${n} ${n === 1 ? one : many}`

/** How well the base run matches the revision the change starts from. */
function baseMatch(comparison: Comparison): string {
    if (comparison.exact) return 'exact'
    if (comparison.distance > 0)
        return `${plural(comparison.distance, 'revision', 'revisions')} before the base revision`
    return 'older than the base revision'
}

function changeLine(change: ChangeSet | undefined): string | null {
    if (!change) return null
    const groups: [number, string][] = [
        [change.measured?.length ?? 0, 'measured'],
        [change.ignored?.length ?? 0, 'ignored by the coverage filters'],
        [change.notInReports?.length ?? 0, 'not in the coverage reports'],
        [change.deleted?.length ?? 0, 'deleted'],
    ]
    const total = groups.reduce((sum, [n]) => sum + n, 0)
    if (total === 0) return null
    const parts = groups.filter(([n]) => n > 0).map(([n, what]) => `${n} ${what}`)
    return `The change has ${plural(total, 'file', 'files')}: ${parts.join(', ')}.`
}

/** The answer the report exists for: did coverage go up or down against the base run. */
function DeltaSummary({ comparison }: { comparison: Comparison }) {
    const headline = comparison.metrics.find((m) => m.key === comparison.headline) ?? comparison.metrics[0]
    const others = comparison.metrics.filter((m) => m !== headline)
    const base = comparison.base
    const baseName = base.revision ? shortRevision(base.revision) : `run ${base.id ?? '?'}`
    // only a server has a page for the base run
    const baseHref = window.__NANOVISION_MODE__ === 'server' && base.id ? `/runs/${base.id}` : undefined
    const change = comparison.change
    const testOnly =
        !!change && !change.measured?.length && !change.deleted?.length && (change.ignored?.length ?? 0) > 0
    const changed = changeLine(change)

    return (
        <section className="rounded-md border border-border bg-card p-4">
            {headline && (
                <>
                    <div className="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                        <span
                            className={cn('font-semibold text-3xl tabular-nums', deltaTone(headline.delta))}
                            title="Change in percentage points"
                        >
                            {signed(headline.delta)}
                        </span>
                        <span className="font-medium text-sm">{headline.label}</span>
                        <span className="text-muted-foreground text-sm tabular-nums">
                            {headline.base.toFixed(2)}% → {headline.current.toFixed(2)}%
                        </span>
                    </div>
                    <p className="mt-1 text-muted-foreground text-xs">
                        {signedCount(headline.currentCovered - headline.baseCovered)} covered,{' '}
                        {signedCount(headline.currentTotal - headline.baseTotal)} {headline.unit} against{' '}
                        {baseHref ? (
                            <a href={baseHref} className="text-primary hover:underline" title={base.revision}>
                                {baseName}
                            </a>
                        ) : (
                            <span className="font-mono text-foreground" title={base.revision}>
                                {baseName}
                            </span>
                        )}{' '}
                        ({baseMatch(comparison)})
                    </p>
                </>
            )}

            {others.length > 0 && (
                <dl className="mt-3 flex flex-wrap gap-x-6 gap-y-1 border-border border-t pt-3 text-sm">
                    {others.map((m) => (
                        <div key={m.key} className="flex items-baseline gap-2">
                            <dt className="text-muted-foreground">{m.label}</dt>
                            <dd className={cn('font-medium tabular-nums', deltaTone(m.delta))}>{signed(m.delta)}</dd>
                        </div>
                    ))}
                </dl>
            )}

            {(changed || testOnly) && (
                <p className="mt-3 text-muted-foreground text-xs">
                    {changed}
                    {testOnly && ' It has no measured code, so only this delta shows what it is worth.'}
                </p>
            )}

            {comparison.warnings?.map((warning) => (
                <p key={warning} className="mt-2 flex items-start gap-1.5 text-partial text-xs">
                    <AlertTriangle className="mt-0.5 h-3.5 w-3.5 shrink-0" />
                    {warning}
                </p>
            ))}
        </section>
    )
}

interface Row {
    path: string
    url?: string
    change?: 'added' | 'modified'
    patch?: CoverageDetail
    delta?: FileDelta
}

// rows drawn at once; a submit run far from its base run can move thousands of files
const MAX_ROWS = 200

const patchOf = (node: FileNode): CoverageDetail | undefined => {
    const patch = node.metrics?.patch_statement_coverage ?? node.metrics?.patch_line_coverage
    return patch && 'percentage' in patch ? patch : undefined
}

/** The changed files, and the files whose coverage moved, largest movement first. */
function ChangedFiles({ nodes, files }: { nodes: FileNode[]; files: FileDelta[] }) {
    const rows = useMemo(() => {
        const byPath = new Map<string, Row>()
        const urls = new Map<string, string>()
        for (const node of nodes) {
            if (node.type !== 'file') continue
            if (node.targetUrl) urls.set(node.path, node.targetUrl)
            if (node.diffStatus === 'added' || node.diffStatus === 'modified') {
                byPath.set(node.path, { path: node.path, change: node.diffStatus, patch: patchOf(node) })
            }
        }
        for (const file of files) {
            const row = byPath.get(file.path) ?? { path: file.path }
            row.delta = file
            byPath.set(file.path, row)
        }
        const uncovered = (row: Row) => (row.patch ? row.patch.total - row.patch.covered : 0)
        return [...byPath.values()]
            .map((row) => ({ ...row, url: urls.get(row.path) }))
            .sort(
                (a, b) =>
                    Math.abs(b.delta?.delta ?? 0) - Math.abs(a.delta?.delta ?? 0) ||
                    uncovered(b) - uncovered(a) ||
                    a.path.localeCompare(b.path),
            )
    }, [nodes, files])

    if (rows.length === 0) {
        return (
            <p className="rounded-md border border-border bg-card p-4 text-muted-foreground text-sm">
                No measured file changed, and no file's coverage moved.
            </p>
        )
    }

    const hasDelta = files.length > 0
    return (
        <section className="overflow-hidden rounded-md border border-border bg-card">
            <div className="flex items-baseline gap-2 border-border border-b bg-muted/40 px-3 py-2">
                <h2 className="font-semibold text-sm">Files</h2>
                <span className="text-muted-foreground text-xs">
                    {hasDelta ? 'changed, or with coverage that moved; largest movement first' : 'changed by this run'}
                </span>
            </div>
            <div className="overflow-x-auto">
                <table className="w-full text-sm">
                    <thead>
                        <tr className="text-left text-muted-foreground text-xs">
                            <th className="px-3 py-1.5 font-medium">File</th>
                            <th className="px-2 py-1.5 text-center font-medium">Change</th>
                            <th className="px-2 py-1.5 text-right font-medium">Patch coverage</th>
                            {hasDelta && (
                                <>
                                    <th className="px-2 py-1.5 text-right font-medium">Base</th>
                                    <th className="px-2 py-1.5 text-right font-medium">Now</th>
                                    <th className="px-3 py-1.5 text-right font-medium">Delta</th>
                                </>
                            )}
                        </tr>
                    </thead>
                    <tbody>
                        {rows.slice(0, MAX_ROWS).map((row) => (
                            <tr key={row.path} className="border-border border-t hover:bg-muted/40">
                                <td className="max-w-[32rem] truncate px-3 py-1.5 font-mono text-xs" title={row.path}>
                                    {row.url ? (
                                        <a href={row.url} className="text-primary hover:underline">
                                            {row.path}
                                        </a>
                                    ) : (
                                        row.path
                                    )}
                                </td>
                                <td className="px-2 py-1.5 text-center">
                                    <DiffStatusBadge status={row.change} />
                                </td>
                                <td className="px-2 py-1.5 text-right tabular-nums">
                                    {row.patch ? (
                                        <>
                                            <span className="text-muted-foreground text-xs">
                                                {row.patch.covered} / {row.patch.total}
                                            </span>{' '}
                                            <span
                                                className={cn(
                                                    row.patch.percentage < 50 && 'text-uncovered',
                                                    row.patch.percentage >= 50 &&
                                                        row.patch.percentage < 80 &&
                                                        'text-partial',
                                                )}
                                            >
                                                {row.patch.percentage.toFixed(1)}%
                                            </span>
                                        </>
                                    ) : (
                                        <span className="text-muted-foreground">—</span>
                                    )}
                                </td>
                                {hasDelta && (
                                    <>
                                        <td className="px-2 py-1.5 text-right text-muted-foreground tabular-nums">
                                            {row.delta ? `${row.delta.base.toFixed(1)}%` : '—'}
                                        </td>
                                        <td className="px-2 py-1.5 text-right tabular-nums">
                                            {row.delta ? `${row.delta.current.toFixed(1)}%` : '—'}
                                        </td>
                                        <td
                                            className={cn(
                                                'px-3 py-1.5 text-right font-semibold tabular-nums',
                                                row.delta ? deltaTone(row.delta.delta, 1) : 'text-muted-foreground',
                                            )}
                                        >
                                            {row.delta ? signed(row.delta.delta, 1) : '—'}
                                        </td>
                                    </>
                                )}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            {rows.length > MAX_ROWS && (
                <p className="border-border border-t px-3 py-2 text-muted-foreground text-xs">
                    {rows.length - MAX_ROWS} more files are not listed. The Files tab has them all.
                </p>
            )}
        </section>
    )
}

interface ChangesTabProps {
    /** The delta against the base run, when the run has one. */
    comparison?: Comparison
    /** The verdict on the changed code, when the run has a diff. */
    review?: Review
    nodes: FileNode[]
}

/**
 * What this run changed: the coverage delta against its base run, the verdict
 * on the changed code, and the files behind both.
 */
export default function ChangesTab({ comparison, review, nodes }: ChangesTabProps) {
    return (
        <div className="flex flex-col gap-4">
            <p className="rounded-md border border-partial/40 bg-partial/10 px-4 py-2 text-sm">
                Work in progress: this view is being developed and will change.
            </p>
            {comparison ? (
                <DeltaSummary comparison={comparison} />
            ) : (
                <p className="rounded-md border border-border bg-card p-4 text-muted-foreground text-sm">
                    This run has no base run, so the change in overall coverage is not known. Patch coverage below is
                    measured on the changed lines alone.
                </p>
            )}
            {/* a change of tests alone has no changed code to judge */}
            {review && review.stats.changedFiles > 0 && <ReviewSummary review={review} />}
            <ChangedFiles nodes={nodes} files={comparison?.files ?? []} />
        </div>
    )
}
