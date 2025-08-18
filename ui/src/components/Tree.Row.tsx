import { ChevronDown, ChevronRight, File, Folder, FolderOpen } from 'lucide-react'
import InlineCoverage from '@/components/InlineCoverage'
import { SUB_METRIC_COLS } from '@/lib/consts'
import { cn } from '@/lib/utils'
import type { FileNode, MetricKey, Metrics } from '@/types/summary'

export function TreeRow({
    node,
    depth,
    enabledMetrics,
    isExpanded,
    onToggleFolder,
    metricsForNode,
    viewMode,
    index,
    isPinned,
}: {
    node: FileNode
    depth: number
    enabledMetrics: { id: MetricKey; label: string }[]
    isExpanded: boolean
    onToggleFolder: (id: string) => void
    metricsForNode: (n: FileNode) => Partial<Metrics> | undefined
    viewMode: 'tree' | 'flat'
    index: number
    isPinned: boolean
}) {
    const isFolder = node.type === 'folder' && viewMode === 'tree'
    const metrics = metricsForNode(node)
    const indentPx = viewMode === 'tree' ? depth * 20 : 0
    const isOdd = index % 2 !== 0

    const interactiveProps = isFolder
        ? {
              role: 'button',
              tabIndex: 0,
              onClick: () => onToggleFolder(node.id),
              onKeyDown: (e: React.KeyboardEvent<HTMLDivElement>) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                      e.preventDefault()
                      onToggleFolder(node.id)
                  }
              },
          }
        : {}

    const totalMetricsWidth = enabledMetrics.reduce((sum) => sum + SUB_METRIC_COLS.reduce((s, c) => s + c.width, 0), 0)

    return (
        <div
            className={cn(
                'group grid w-full items-center',
                'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
                isFolder ? 'cursor-pointer' : 'cursor-default',
            )}
            {...interactiveProps}
        >
            <div
                className="grid h-full"
                style={{
                    gridTemplateColumns: `minmax(300px, 1fr) ${totalMetricsWidth}px`,
                }}
            >
                <div
                    className={cn(
                        'flex min-w-0 items-center gap-2 border-border border-r',
                        isPinned && 'sticky left-0 z-10',
                        isOdd ? 'bg-subtle' : 'bg-background',
                        'group-hover:bg-muted',
                    )}
                >
                    <div style={{ paddingLeft: indentPx }}>
                        {isFolder ? (
                            isExpanded ? (
                                <ChevronDown className="h-4 w-4 shrink-0 text-muted-foreground" />
                            ) : (
                                <ChevronRight className="h-4 w-4 shrink-0 text-muted-foreground" />
                            )
                        ) : viewMode === 'tree' ? (
                            <div className="w-4 shrink-0" />
                        ) : null}
                    </div>

                    {isFolder ? (
                        isExpanded ? (
                            <FolderOpen className="h-4 w-4 shrink-0 text-primary" />
                        ) : (
                            <Folder className="h-4 w-4 shrink-0 text-primary" />
                        )
                    ) : (
                        <File className="h-4 w-4 text-muted-foreground" />
                    )}

                    <span
                        className={cn(
                            'truncate',
                            isFolder ? 'font-semibold text-foreground' : 'font-medium text-foreground/90',
                        )}
                        title={node.path}
                    >
                        {viewMode === 'flat' ? node.path : node.name}
                    </span>
                </div>

                <div
                    className={cn('grid items-center', isOdd ? 'bg-subtle' : 'bg-background', 'group-hover:bg-muted')}
                    style={{
                        gridTemplateColumns: `repeat(${enabledMetrics.length}, 1fr)`,
                    }}
                >
                    {enabledMetrics.map((cfg, index) => {
                        const metricData = metrics?.[cfg.id]
                        return (
                            <div
                                key={cfg.id}
                                className={cn(
                                    'grid h-full items-center',
                                    index < enabledMetrics.length - 1 && 'border-border border-r',
                                )}
                                style={{ gridTemplateColumns: SUB_METRIC_COLS.map((c) => `${c.width}px`).join(' ') }}
                            >
                                <div className="px-2 text-right text-xs tabular-nums">{metricData?.covered ?? '-'}</div>
                                <div className="px-2 text-right text-xs tabular-nums">
                                    {metricData?.uncovered ?? '-'}
                                </div>
                                <div className="px-2 text-right text-xs tabular-nums">
                                    {metricData?.coverable ?? '-'}
                                </div>
                                <div className="px-2 text-right text-xs tabular-nums">{metricData?.total ?? '-'}</div>
                                <div className="px-2 text-right text-xs tabular-nums">
                                    {metricData !== undefined ? (
                                        <InlineCoverage
                                            percentage={metricData.percentage}
                                            risk={node.statuses?.[cfg.id] ?? 'safe'}
                                            isFolder={isFolder}
                                        />
                                    ) : (
                                        <span className="text-muted-foreground text-xs">-</span>
                                    )}
                                </div>
                            </div>
                        )
                    })}
                </div>
            </div>
        </div>
    )
}
