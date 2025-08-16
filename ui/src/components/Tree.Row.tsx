import { ChevronDown, ChevronRight, Folder, FolderOpen } from 'lucide-react'
import { getFileIcon } from '@/components/fileIcons'
import InlineCoverage from '@/components/InlineCoverage'
import { cn } from '@/lib/utils'
import type { CSSWithVars } from '@/types/css'
import type { FileNode, MetricKey, Metrics, RiskLevel } from '@/types/summary'

export function TreeRow({
    node,
    depth,
    enabledMetrics,
    rowDensity,
    isExpanded,
    onToggleFolder,
    metricsForNode,
    riskForMetric,
    filesInFolderCount,
}: {
    node: FileNode
    depth: number
    enabledMetrics: { id: MetricKey; label: string }[]
    rowDensity: 'comfortable' | 'compact'
    isExpanded: boolean
    onToggleFolder: (id: string) => void
    metricsForNode: (n: FileNode) => Partial<Metrics> | undefined
    riskForMetric: (n: FileNode, k: MetricKey) => RiskLevel
    filesInFolderCount: (n: FileNode) => number
}) {
    const isFolder = node.type === 'folder'
    const metrics = metricsForNode(node)
    const indentPx = depth * 20
    const rowY = rowDensity === 'compact' ? 'py-0.5' : 'py-1.5'

    const onClickRow = () => {
        if (isFolder) onToggleFolder(node.id)
    }

    const onKeyDownRow: React.KeyboardEventHandler<HTMLTableRowElement> = (e) => {
        if (!isFolder) return
        if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault()
            onToggleFolder(node.id)
        }
    }

    // The indent lines are now positioned relative to the <tr>.
    // A wrapper div is needed for positioning context.
    return (
        <tr
            aria-expanded={isFolder ? isExpanded : undefined}
            tabIndex={isFolder ? 0 : -1}
            className={cn(
                'grid items-center bg-background px-2 pr-3 hover:bg-primary/20',
                'selection:bg-primary selection:text-primary-foreground',
                isFolder ? 'cursor-pointer' : 'cursor-default',
                'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
                'relative', // Needed for the indent lines
                rowY,
            )}
            style={{
                ...({
                    '--metric-count': String(enabledMetrics.length),
                    '--metric-col-width': '170px',
                } as CSSWithVars),
                gridTemplateColumns: '1fr repeat(var(--metric-count), var(--metric-col-width))',
            }}
            onClick={onClickRow}
            onKeyDown={onKeyDownRow}
        >
            {depth > 0 && (
                <td className="pointer-events-none absolute inset-y-0 left-0">
                    {Array.from({ length: depth }, (_, i) => (
                        <div
                            key={`${node.id}-indent-${i}`}
                            className="absolute top-0 bottom-0 w-px bg-border"
                            style={{ left: `${i * 20 + 18}px` }}
                        />
                    ))}
                </td>
            )}

            {/* Name cell */}
            <td className="flex min-w-0 items-center gap-2" style={{ paddingLeft: indentPx }}>
                {isFolder ? (
                    <>
                        {isExpanded ? (
                            <ChevronDown className="h-4 w-4 text-muted-foreground" />
                        ) : (
                            <ChevronRight className="h-4 w-4 text-muted-foreground" />
                        )}
                        {isExpanded ? (
                            <FolderOpen className="h-4 w-4 text-primary" />
                        ) : (
                            <Folder className="h-4 w-4 text-primary" />
                        )}
                    </>
                ) : (
                    <>
                        <div className="w-4" />
                        {getFileIcon(node.name)}
                    </>
                )}
                <span
                    className={cn(
                        'truncate',
                        isFolder ? 'font-semibold text-foreground' : 'font-medium text-foreground/90',
                    )}
                    title={node.path}
                >
                    {node.name}
                    {isFolder && (
                        <span className="ml-2 font-normal text-muted-foreground text-xs">
                            ({filesInFolderCount(node)} files)
                        </span>
                    )}
                </span>
            </td>

            {/* Metric cells */}
            {enabledMetrics.map((cfg) => {
                const value = metrics?.[cfg.id] as number | undefined
                const risk = riskForMetric(node, cfg.id)
                return (
                    <td key={cfg.id} className="flex items-center px-2">
                        {value !== undefined ? (
                            <InlineCoverage percentage={value} risk={risk} isFolder={isFolder} density={rowDensity} />
                        ) : (
                            <span className="text-muted-foreground text-xs">-</span>
                        )}
                    </td>
                )
            })}
        </tr>
    )
}
