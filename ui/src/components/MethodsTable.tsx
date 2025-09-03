import { useMemo } from 'react'
import { camelCaseToTitleCase } from '@/lib/utils'
import type { Method, MetricDefinitions } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'
import { StatusIcon } from './MetricCard'

export default function MethodsTable({
    methods,
    metricDefinitions,
}: {
    methods: Method[]
    metricDefinitions: MetricDefinitions
}) {
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

    if (!methods || methods.length === 0) {
        return null
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Method Coverage</CardTitle>
            </CardHeader>
            <CardContent className="overflow-x-auto p-0">
                {/* --- 1. REMOVED 'table-fixed' to use the browser's auto-sizing algorithm --- */}
                <table className="w-full text-sm">
                    <thead>
                        <tr className="border-border border-b bg-subtle/50 font-semibold text-xs">
                            <th className="text-nowrap px-4 py-2 text-right text-muted-foreground">Line #</th>
                            {/* --- 2. THE KEY: Tell this column to expand --- */}
                            <th className="w-full px-4 py-2 text-left text-muted-foreground">Method</th>
                            {metricConfigs.map((mc) => (
                                <th
                                    key={mc.id}
                                    className="whitespace-nowrap px-4 py-2 text-right text-muted-foreground"
                                >
                                    {mc.shortLabel}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {methods.map((method) => (
                            <tr key={method.name} className="group border-border/50 border-b hover:bg-accent/50">
                                <td className="whitespace-nowrap px-4 py-1.5 text-right font-mono text-muted-foreground text-xs">
                                    {method.startLine}
                                </td>
                                <td className="px-4 py-1.5 text-left font-mono">
                                    <button
                                        type="button"
                                        className="truncate text-left hover:text-primary hover:underline"
                                        title={method.name}
                                        onClick={() => handleGoToLine(method.startLine)}
                                    >
                                        {method.name}
                                    </button>
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
