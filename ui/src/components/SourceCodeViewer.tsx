import { useMemo } from 'react'
import CodeLine from '@/components/CodeLine'
import SourceCodeHeader from '@/components/SourceCodeHeader'
import type { DetailsV1 } from '@/lib/validation'
import type { LineStatus } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

interface SourceCodeViewerProps {
    fileName: string
    lines: DetailsV1['lines']
    activeReportIndices: Set<number>
    reports?: DetailsV1['reports']
}

export default function SourceCodeViewer({ fileName, lines, activeReportIndices, reports }: SourceCodeViewerProps) {
    const processedLines = useMemo(() => {
        return lines.map((line) => {
            if (line.status === 'not-coverable') {
                return { ...line, hits: undefined, reportHits: [] }
            }

            const totalHits =
                line.hits?.reduce((sum, hitCount, index) => {
                    if (activeReportIndices.has(index)) {
                        return sum + hitCount
                    }
                    return sum
                }, 0) ?? 0

            // Per-report breakdown for the enabled reports that actually touched
            // this line (hits > 0), used for the hit-count tooltip.
            const reportHits =
                line.hits?.reduce<{ name: string; hits: number }[]>((acc, hitCount, index) => {
                    if (hitCount > 0 && activeReportIndices.has(index)) {
                        acc.push({ name: reports?.[index]?.name ?? `Report ${index + 1}`, hits: hitCount })
                    }
                    return acc
                }, []) ?? []

            let status: LineStatus
            if (totalHits > 0) {
                status = line.branchInfo ? 'partial' : 'covered'
            } else {
                status = 'uncovered'
            }

            return { ...line, hits: totalHits, status, reportHits }
        })
    }, [lines, activeReportIndices, reports])

    return (
        <Card>
            <CardHeader>
                <CardTitle className="font-mono text-lg">{fileName}</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
                <div className="w-full overflow-x-auto">
                    <div className="min-w-max">
                        <SourceCodeHeader />
                        <div>
                            {processedLines.map((line) => (
                                <CodeLine key={line.lineNumber} {...line} />
                            ))}
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
