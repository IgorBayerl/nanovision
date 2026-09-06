import { ArrowDown, ArrowUp, ChevronsUpDown } from 'lucide-react'
import { useMemo } from 'react'
import DiffStatusBadge from '@/components/DiffStatusBadge'
import InfoTooltip from '@/components/InfoTooltip'
import { useUrlState } from '@/hooks/useUrlState'
import { scrollToLine } from '@/lib/scrollToLine'
import { camelCaseToTitleCase, cn } from '@/lib/utils'
import type { Method, MetricDefinitions } from '@/types/summary'
import { Card, CardContent, CardHeader } from '@/ui/card'
import { Input } from '@/ui/input'
import { StatusIcon } from './MetricCard'

// Parses the leading number out of a method metric value ("12", "3 / 5", "80%").
const parseMetricNumber = (value?: string): number => {
    if (!value) return Number.NEGATIVE_INFINITY
    const match = value.match(/-?\d+(\.\d+)?/)
    return match ? Number.parseFloat(match[0]) : Number.NEGATIVE_INFINITY
}

const SortIcon = ({ active, dir }: { active: boolean; dir: 'asc' | 'desc' }) => {
    if (!active) return <ChevronsUpDown className="h-3 w-3 opacity-40" />
    return dir === 'asc' ? <ArrowUp className="h-3 w-3" /> : <ArrowDown className="h-3 w-3" />
}

const SortableTh = ({
    columnKey,
    label,
    description,
    align = 'right',
    sortKey,
    sortDir,
    onSort,
}: {
    columnKey: string
    label: string
    description?: string
    align?: 'left' | 'right'
    sortKey: string
    sortDir: 'asc' | 'desc'
    onSort: (key: string) => void
}) => (
    <th
        className={cn(
            'whitespace-nowrap px-4 py-2 text-muted-foreground',
            align === 'left' ? 'text-left' : 'text-right',
            columnKey === 'name' && 'w-full',
        )}
    >
        {/* The info icon sits beside the sort button rather than inside it, so
            hovering the explanation cannot re-sort the table. It always trails
            the label: only the button reverses, putting the sort arrow first. */}
        <div className="inline-flex items-center gap-1.5">
            <button
                type="button"
                onClick={() => onSort(columnKey)}
                className={cn(
                    'inline-flex items-center gap-1 hover:text-foreground',
                    align === 'right' && 'flex-row-reverse',
                )}
            >
                <span>{label}</span>
                <SortIcon active={sortKey === columnKey} dir={sortDir} />
            </button>
            <InfoTooltip label={`What ${label} means`}>{description}</InfoTooltip>
        </div>
    </th>
)

export default function MethodsTable({
    methods,
    metricDefinitions,
}: {
    methods: Method[]
    metricDefinitions: MetricDefinitions
}) {
    const [query, setQuery] = useUrlState('mq', '')
    const [sortKey, setSortKey] = useUrlState('msort', 'line')
    const [sortDir, setSortDir] = useUrlState<'asc' | 'desc'>('msortDir', 'asc')

    const metricConfigs = useMemo(() => {
        if (!methods || methods.length === 0) return []
        const keys = new Set<string>()
        for (const method of methods) {
            for (const key of Object.keys(method.metrics)) keys.add(key)
        }
        // Sorted; the backend prefixes them (a_, b_, ...) to enforce order.
        return Array.from(keys)
            .sort()
            .map((id) => {
                const def = metricDefinitions[id]
                return {
                    id,
                    label: def?.label ?? camelCaseToTitleCase(id),
                    shortLabel: def?.shortLabel ?? camelCaseToTitleCase(id),
                    description: def?.description,
                }
            })
    }, [methods, metricDefinitions])

    const handleSort = (key: string) => {
        if (sortKey === key) {
            setSortDir(sortDir === 'asc' ? 'desc' : 'asc')
        } else {
            setSortKey(key)
            setSortDir('asc')
        }
    }

    const visibleMethods = useMemo(() => {
        const q = query.trim().toLowerCase()
        let rows = methods
        if (q) rows = rows.filter((m) => m.name.toLowerCase().includes(q))

        const sorted = [...rows].sort((a, b) => {
            let cmp: number
            if (sortKey === 'name') {
                cmp = a.name.localeCompare(b.name)
            } else if (sortKey === 'line') {
                cmp = a.startLine - b.startLine
            } else {
                cmp = parseMetricNumber(a.metrics[sortKey]?.value) - parseMetricNumber(b.metrics[sortKey]?.value)
            }
            return sortDir === 'asc' ? cmp : -cmp
        })
        return sorted
    }, [methods, query, sortKey, sortDir])

    if (!methods || methods.length === 0) return null

    return (
        <Card className="rounded-md">
            <CardHeader className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <div className="font-semibold text-lg">Method Coverage</div>
                <div className="flex flex-wrap items-center gap-2">
                    <Input
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        placeholder="Filter methods…"
                        className="h-8 w-full md:w-64"
                    />
                </div>
            </CardHeader>
            <CardContent className="overflow-x-auto p-0">
                <table className="w-full text-sm">
                    <thead>
                        <tr className="border-border border-b bg-subtle/50 font-semibold text-xs">
                            <SortableTh
                                columnKey="line"
                                label="Line #"
                                align="right"
                                sortKey={sortKey}
                                sortDir={sortDir}
                                onSort={handleSort}
                            />
                            <SortableTh
                                columnKey="name"
                                label="Method"
                                align="left"
                                sortKey={sortKey}
                                sortDir={sortDir}
                                onSort={handleSort}
                            />
                            {metricConfigs.map((mc) => (
                                <SortableTh
                                    key={mc.id}
                                    columnKey={mc.id}
                                    label={mc.shortLabel}
                                    description={mc.description}
                                    align="right"
                                    sortKey={sortKey}
                                    sortDir={sortDir}
                                    onSort={handleSort}
                                />
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {visibleMethods.map((method) => (
                            <tr
                                key={`${method.name}-${method.startLine}`}
                                className="group border-border/50 border-b hover:bg-accent/50"
                            >
                                <td className="whitespace-nowrap px-4 py-1.5 text-right font-mono text-muted-foreground text-xs">
                                    {method.startLine}
                                </td>
                                <td className="px-4 py-1.5 text-left font-mono">
                                    <div className="flex min-w-0 items-center gap-2">
                                        <button
                                            type="button"
                                            className="truncate text-left hover:text-primary hover:underline"
                                            title={method.name}
                                            onClick={() => scrollToLine(method.startLine)}
                                        >
                                            {method.name}
                                        </button>
                                        <DiffStatusBadge status={method.diffStatus} className="ml-auto" />
                                    </div>
                                </td>
                                {metricConfigs.map((mc) => {
                                    const metric = method.metrics[mc.id]
                                    return (
                                        <td
                                            key={mc.id}
                                            className="whitespace-nowrap px-4 py-1.5 text-right font-mono text-xs"
                                        >
                                            <div className="flex items-center justify-end gap-2">
                                                {metric?.status && <StatusIcon status={metric?.status} />}
                                                <span>{metric?.value ?? '-'}</span>
                                            </div>
                                        </td>
                                    )
                                })}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </CardContent>
        </Card>
    )
}
