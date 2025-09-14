import { Pin, PinOff } from 'lucide-react'
import HeaderRangeSlider from '@/components/HeaderRangeSlider'
import { cn } from '@/lib/utils'
import type { FilterRange, MetricConfig, MetricKey, SortDir, SortKey } from '@/types/summary'
import { Button } from '@/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/ui/tooltip'

interface HeaderProps {
    isNameColumnPinned: boolean
    onPinColumn: (pinned: boolean) => void
    enabledMetrics: MetricConfig[]
    sortKey: SortKey
    sortDir: SortDir
    onHeaderClick: (key: SortKey) => void
    filterRanges: Record<MetricKey, FilterRange>
    onRangeUpdate: (id: MetricKey, vals: [number, number]) => void
    totalMetricsWidth: number
}

export default function FileExplorerHeader({
    isNameColumnPinned,
    onPinColumn,
    enabledMetrics,
    sortKey,
    sortDir,
    onHeaderClick,
    filterRanges,
    onRangeUpdate,
    totalMetricsWidth,
}: HeaderProps) {
    return (
        <div className="sticky top-0 z-20 grid bg-background font-semibold text-xs">
            <div
                className="grid"
                style={{
                    gridTemplateColumns: `minmax(300px, 1fr) ${totalMetricsWidth}px`,
                }}
            >
                <div
                    className={cn(
                        'border-border border-r border-b py-2',
                        isNameColumnPinned && 'sticky left-0 z-10 bg-background',
                    )}
                >
                    <div className="flex items-center justify-between px-2">
                        <button
                            type="button"
                            className="flex items-center gap-1 text-left text-muted-foreground hover:text-foreground"
                            onClick={() => onHeaderClick('name')}
                        >
                            File / Folder
                            <span className="inline-flex h-3 w-3 items-center justify-center">
                                {sortKey === 'name' ? (sortDir === 'asc' ? '▲' : '▼') : null}
                            </span>
                        </button>
                        <TooltipProvider delayDuration={100}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="h-6 w-6 p-0 text-muted-foreground hover:text-foreground"
                                        onClick={() => onPinColumn(!isNameColumnPinned)}
                                        aria-label={isNameColumnPinned ? 'Unpin column' : 'Pin column'}
                                    >
                                        {isNameColumnPinned ? (
                                            <PinOff className="h-4 w-4" />
                                        ) : (
                                            <Pin className="h-4 w-4" />
                                        )}
                                    </Button>
                                </TooltipTrigger>
                                <TooltipContent>
                                    <p>{isNameColumnPinned ? 'Unpin column' : 'Pin column'}</p>
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>
                    </div>
                </div>
                <div className="grid" style={{ gridTemplateColumns: `repeat(${enabledMetrics.length}, 1fr)` }}>
                    {enabledMetrics.map((m, index) => (
                        <div
                            key={m.id}
                            className={cn(
                                'flex flex-col gap-2 border-border border-b py-2 text-left',
                                index < enabledMetrics.length - 1 && 'border-r',
                            )}
                        >
                            <div className="px-2 font-bold text-foreground">{m.shortLabel}</div>

                            <div className="px-2">
                                <HeaderRangeSlider
                                    range={filterRanges[m.id] ?? { min: 0, max: 100 }}
                                    onRangeUpdate={(vals) => onRangeUpdate(m.id, vals)}
                                />
                            </div>

                            <div
                                className="mt-1 grid"
                                style={{
                                    gridTemplateColumns: m.definition.subMetrics.map((c) => `${c.width}px`).join(' '),
                                }}
                            >
                                {m.definition.subMetrics.map((sub) => (
                                    <button
                                        key={sub.id}
                                        type="button"
                                        className="flex w-full items-center justify-end gap-1 px-2 text-left text-muted-foreground hover:text-foreground"
                                        onClick={() => onHeaderClick({ metric: m.id, subMetric: sub.id })}
                                        title={`Sort by ${m.label} ${sub.label}`}
                                    >
                                        <span className="inline-flex h-3 w-3 items-center justify-center">
                                            {typeof sortKey === 'object' &&
                                            sortKey.metric === m.id &&
                                            sortKey.subMetric === sub.id
                                                ? sortDir === 'asc'
                                                    ? '▲'
                                                    : '▼'
                                                : null}
                                        </span>
                                        {sub.label}
                                    </button>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    )
}
