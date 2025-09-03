import CodeLine from '@/components/CodeLine'
import InlineCoverage from '@/components/InlineCoverage'
import type { IslandProps } from '@/islands/placeholder'
import type { LineDetails, RiskLevel } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

interface SourceCodeViewerProps {
    fileName?: string
    coverage?: number
    risk?: RiskLevel
    lines?: LineDetails[]
}

export function SourceCodeViewer(props: IslandProps) {
    const { fileName, coverage, risk, lines } = props as SourceCodeViewerProps

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-lg">
                    File: <span className="font-mono">{fileName ?? 'Unknown File'}</span>
                </CardTitle>
                <div className="mt-2 w-full max-w-sm">
                    {typeof coverage === 'number' && typeof risk === 'string' ? (
                        <InlineCoverage percentage={coverage} risk={risk} />
                    ) : (
                        <p className="text-muted-foreground text-sm">No coverage data provided.</p>
                    )}
                </div>
            </CardHeader>
            <CardContent className="overflow-x-auto rounded-md border border-border bg-subtle p-0">
                <div className="relative font-mono">
                    {lines && lines.length > 0 ? (
                        lines.map((line) => <CodeLine key={line.lineNumber} {...line} />)
                    ) : (
                        <div className="p-4 text-center text-muted-foreground">No source code to display.</div>
                    )}
                </div>
            </CardContent>
        </Card>
    )
}
