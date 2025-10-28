import InlineCoverage from '@/components/InlineCoverage'
import { cn } from '@/lib/utils'
import type { FileNode, MetricConfig, Metrics } from '@/types/summary'
import { ChevronDown, ChevronRight, File, Folder, FolderOpen } from 'lucide-react'

const NodeName = ({ node, viewMode }: { node: FileNode; viewMode: 'tree' | 'flat' }) => {
    const isFolder = node.type === 'folder' && viewMode === 'tree'
    const content = viewMode === 'flat' ? node.path : node.name
    const commonClasses = cn('truncate', isFolder ? 'font-semibold' : 'font-medium text-foreground/90')

    // If it's a file and has a targetUrl, render it as a link.
    if (!isFolder && node.targetUrl) {
        return (
            <a
                href={node.targetUrl}
                className={cn(commonClasses, 'hover:text-primary hover:underline')}
                title={node.path}
            >
                {content}
            </a>
        )
    }

    // Otherwise, render it as a plain span.
    return (
        <span className={commonClasses} title={node.path}>
            {content}
        </span>
    )
}

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
    enabledMetrics: MetricConfig[]
    isExpanded: boolean
    onToggleFolder: (id: string, event: React.MouseEvent | React.KeyboardEvent) => void
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
              onClick: (e: React.MouseEvent) => onToggleFolder(node.id, e),
              onKeyDown: (e: React.KeyboardEvent<HTMLDivElement>) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                      e.preventDefault()
                      onToggleFolder(node.id, e)
                  }
              },
          }
        : {}

    const totalMetricsWidth = enabledMetrics.reduce(
        (sum, metric) => sum + metric.definition.subMetrics.reduce((s, c) => s + c.width, 0),
        0,
    )

    return (
        <div
            className={cn(
                'group grid w-full items-center',
                isFolder ? 'text-foreground/70 font-bold' : '',
            )}
        >
            <div className="grid h-full" style={{ gridTemplateColumns: `minmax(300px, 1fr) ${totalMetricsWidth}px` }}>
                <div
                    className={cn(
                        'flex min-w-0 items-center gap-2 border-border border-r',
                        isPinned && 'sticky left-0 z-10',
                        isOdd ? 'bg-subtle' : 'bg-background',
                        'group-hover:bg-muted',
                        // This makes only this section clickable and focusable for folders
                        isFolder &&
                            'cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
                    )}
                    {...interactiveProps}
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
                        <File className="h-4 w-4 shrink-0 text-muted-foreground" />
                    )}
                    <NodeName node={node} viewMode={viewMode} />
                </div>

                <div
                    className={cn('grid items-center', isOdd ? 'bg-subtle' : 'bg-background', 'group-hover:bg-muted')}
                    style={{ gridTemplateColumns: `repeat(${enabledMetrics.length}, 1fr)` }}
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
                                style={{
                                    gridTemplateColumns: cfg.definition.subMetrics.map((c) => `${c.width}px`).join(' '),
                                }}
                            >
                                {cfg.definition.subMetrics.map((subMetric) => {
                                    if (subMetric.id === 'percentage') {
                                        return (
                                            <div key={subMetric.id} className="px-2 text-right text-xs tabular-nums">
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
                                        )
                                    }
                                    const value = metricData?.[subMetric.id as keyof typeof metricData]
                                    return (
                                        <div key={subMetric.id} className="px-2 text-right text-xs tabular-nums">
                                            {value ?? '-'}
                                        </div>
                                    )
                                })}
                            </div>
                        )
                    })}
                </div>
            </div>
        </div>
    )
}
