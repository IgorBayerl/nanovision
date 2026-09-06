import { AlertCircle, AlertTriangle, ChevronRight, Info } from 'lucide-react'
import { useMemo, useState } from 'react'
import InfoTooltip from '@/components/InfoTooltip'
import { cn } from '@/lib/utils'
import type { Diagnostic, FileNode } from '@/lib/validation'

interface ProblemsPanelProps {
    diagnostics: Diagnostic[]
    /** Flat node list, used to resolve a diagnostic's file to its details page URL. */
    nodes: FileNode[]
}

const severityIcon = (severity: Diagnostic['severity']) => {
    if (severity === 'error') return <AlertCircle className="h-4 w-4 shrink-0 text-uncovered" />
    if (severity === 'warning') return <AlertTriangle className="h-4 w-4 shrink-0 text-partial" />
    return <Info className="h-4 w-4 shrink-0 text-muted-foreground" />
}

/**
 * A collapsible "Problems" panel, styled after an editor's problems tab. It
 * lists the coverage diagnostics (errors/warnings) emitted by the backend and,
 * when expanded, scrolls internally with a bounded max height.
 */
export default function ProblemsPanel({ diagnostics, nodes }: ProblemsPanelProps) {
    // Default to expanded when there are errors, collapsed otherwise.
    const errorCount = diagnostics.filter((d) => d.severity === 'error').length
    const warningCount = diagnostics.filter((d) => d.severity === 'warning').length
    const [open, setOpen] = useState(errorCount > 0)

    // Map file path -> details page URL so each problem can link to its file.
    const urlByPath = useMemo(() => {
        const map = new Map<string, string>()
        for (const node of nodes) {
            if (node.type === 'file' && node.targetUrl) map.set(node.path, node.targetUrl)
        }
        return map
    }, [nodes])

    if (diagnostics.length === 0) return null

    return (
        <div className="mb-4 overflow-hidden rounded-md border border-border bg-card">
            <button
                type="button"
                onClick={() => setOpen((v) => !v)}
                className="flex w-full items-center gap-2 px-3 py-2 text-left hover:bg-muted/50"
                aria-expanded={open}
            >
                <ChevronRight className={cn('h-4 w-4 shrink-0 transition-transform', open && 'rotate-90')} />
                <span className="font-semibold text-sm">Problems</span>
                <InfoTooltip label="How problems are evaluated" className="mt-0.5">
                    Coverage rules evaluated against all reports. The report selection does not change this list.
                </InfoTooltip>
                <span className="ml-1 flex items-center gap-3 text-muted-foreground text-xs">
                    <span className="flex items-center gap-1">
                        <AlertCircle className="h-3.5 w-3.5 text-uncovered" />
                        {errorCount}
                    </span>
                    <span className="flex items-center gap-1">
                        <AlertTriangle className="h-3.5 w-3.5 text-partial" />
                        {warningCount}
                    </span>
                </span>
            </button>

            {open && (
                <ul className="max-h-64 divide-y divide-border overflow-y-auto border-border border-t">
                    {diagnostics.map((d, i) => {
                        const url = urlByPath.get(d.file)
                        const location = (
                            <span className="shrink-0 font-mono text-muted-foreground text-xs">
                                {d.file}:{d.startLine}
                            </span>
                        )
                        return (
                            <li
                                key={`${d.file}-${d.startLine}-${d.ruleId}-${i}`}
                                className="flex items-center gap-2 px-3 py-1.5 text-sm hover:bg-muted/40"
                            >
                                {severityIcon(d.severity)}
                                <span className="min-w-0 flex-1 truncate" title={d.message}>
                                    {d.message}
                                </span>
                                <span className="hidden shrink-0 rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground sm:inline">
                                    {d.ruleId}
                                </span>
                                {url ? (
                                    <a href={url} className="shrink-0 font-mono text-primary text-xs hover:underline">
                                        {d.file}:{d.startLine}
                                    </a>
                                ) : (
                                    location
                                )}
                            </li>
                        )
                    })}
                </ul>
            )}
        </div>
    )
}
