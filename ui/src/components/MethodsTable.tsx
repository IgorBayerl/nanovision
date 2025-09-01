import { Target } from 'lucide-react'
import { useMemo, useState } from 'react'
import InlineCoverage from '@/components/InlineCoverage'
import { camelCaseToTitleCase, cn } from '@/lib/utils'
import type { Method, MetricDefinitions, SortDir } from '@/types/summary'
import { Button } from '@/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

type MethodSortKey = 'name' | { metric: string }

// Helper to get a comparable value for sorting remains the same
function getSortValue(method: Method, key: MethodSortKey): string | number {
    if (key === 'name') {
        return method.name
    }
    const metric = method.metrics[key.metric]
    if (typeof metric === 'number') {
        return metric
    }
    if (typeof metric === 'object' && metric.percentage !== undefined) {
        return metric.percentage
    }
    return -1
}

export default function MethodsTable({
    methods,
    metricDefinitions,
}: {
    methods: Method[]
    metricDefinitions: MetricDefinitions
}) {
    const [sortKey, setSortKey] = useState<MethodSortKey>('name')
    const [sortDir, setSortDir] = useState<SortDir>('asc')

    const handleGoToLine = (lineNumber: number) => {
        const selector = `[data-line-number="${lineNumber}"]`
        const lineElement = document.querySelector(selector) as HTMLElement

        if (lineElement) {
            lineElement.scrollIntoView({
                behavior: 'smooth',
                block: 'center',
            })
            lineElement.classList.add('animate-pulse-bg')
            setTimeout(() => {
                lineElement.classList.remove('animate-pulse-bg')
            }, 1500)
        }
    }

    const metricConfigs = useMemo(() => {
        if (!methods || methods.length === 0) return []
        return Object.keys(methods[0].metrics).map((id) => {
            const def = metricDefinitions[id]
            return {
                id,
                label: def?.label ?? camelCaseToTitleCase(id),
                shortLabel: def?.shortLabel ?? camelCaseToTitleCase(id),
            }
        })
    }, [methods, metricDefinitions])

    const sortedMethods = useMemo(() => {
        return [...methods].sort((a, b) => {
            const valA = getSortValue(a, sortKey)
            const valB = getSortValue(b, sortKey)
            const dir = sortDir === 'asc' ? 1 : -1
            if (typeof valA === 'string' && typeof valB === 'string') {
                return valA.localeCompare(valB) * dir
            }
            return (valA as number) - (valB as number)
        })
    }, [methods, sortKey, sortDir])

    const handleHeaderClick = (key: MethodSortKey) => {
        if (
            sortKey === key ||
            (typeof sortKey === 'object' && typeof key === 'object' && sortKey.metric === key.metric)
        ) {
            setSortDir(sortDir === 'asc' ? 'desc' : 'asc')
        } else {
            setSortKey(key)
            setSortDir('asc')
        }
    }

    if (!methods || methods.length === 0) {
        return null
    }

    // --- THIS IS THE FIX ---
    // Create the CSS value for the grid-template-columns property directly.
    const gridTemplateColumns = `minmax(300px, 2fr) repeat(${metricConfigs.length}, minmax(150px, 1fr)) 60px`

    return (
        <Card>
            <CardHeader>
                <CardTitle>Method Coverage</CardTitle>
            </CardHeader>
            <CardContent className="overflow-x-auto p-0">
                <div className="w-full min-w-max">
                    {/* Table Header */}
                    <div
                        className="grid border-border border-b bg-subtle/50 font-semibold text-xs"
                        // Apply the dynamic style using the 'style' prop
                        style={{ gridTemplateColumns }}
                    >
                        <button
                            type="button"
                            onClick={() => handleHeaderClick('name')}
                            className="sticky left-0 z-10 flex items-center gap-1 bg-subtle/50 px-4 py-3 text-left hover:text-foreground"
                        >
                            Method
                            <span>{sortKey === 'name' ? (sortDir === 'asc' ? '▲' : '▼') : ''}</span>
                        </button>
                        {metricConfigs.map((mc) => (
                            <button
                                type="button"
                                key={mc.id}
                                onClick={() => handleHeaderClick({ metric: mc.id })}
                                className="flex items-center justify-end gap-1 px-4 py-3 text-right hover:text-foreground"
                            >
                                <span>
                                    {typeof sortKey === 'object' && sortKey.metric === mc.id
                                        ? sortDir === 'asc'
                                            ? '▲'
                                            : '▼'
                                        : ''}
                                </span>
                                {mc.shortLabel}
                            </button>
                        ))}
                        <div className="px-4 py-3 text-right" />
                    </div>

                    {/* Table Body */}
                    <div>
                        {sortedMethods.map((method, index) => (
                            <div
                                key={method.name}
                                className={cn(
                                    'group grid items-center',
                                    index % 2 === 0 ? 'bg-background' : 'bg-subtle/30',
                                )}
                                // Apply the same dynamic style to the rows
                                style={{ gridTemplateColumns }}
                            >
                                <div
                                    className="sticky left-0 z-10 truncate bg-inherit px-4 py-2 font-mono text-sm group-hover:bg-muted"
                                    title={method.name}
                                >
                                    {method.name}
                                </div>
                                {metricConfigs.map((mc) => {
                                    const metric = method.metrics[mc.id]
                                    const status = method.statuses?.[mc.id]
                                    return (
                                        <div
                                            key={mc.id}
                                            className="px-4 py-2 text-right font-mono text-xs group-hover:bg-muted"
                                        >
                                            {typeof metric === 'number' ? (
                                                metric.toFixed(2)
                                            ) : metric ? (
                                                <InlineCoverage
                                                    percentage={metric.percentage}
                                                    risk={status ?? 'safe'}
                                                />
                                            ) : (
                                                '-'
                                            )}
                                        </div>
                                    )
                                })}
                                <div className="flex items-center justify-center bg-inherit px-4 py-2 group-hover:bg-muted">
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="h-6 w-6 p-0 text-muted-foreground hover:text-foreground"
                                        onClick={() => handleGoToLine(method.startLine)}
                                        title={`Go to line ${method.startLine}`}
                                    >
                                        <Target className="h-4 w-4" />
                                    </Button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
