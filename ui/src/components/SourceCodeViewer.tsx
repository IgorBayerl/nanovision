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
}

export default function SourceCodeViewer({ fileName, lines, activeReportIndices }: SourceCodeViewerProps) {
    const processedLines = useMemo(() => {
        return lines.map((line) => {
            if (line.status === 'not-coverable') {
                return { ...line, hits: undefined }
            }

            const totalHits =
                line.hits?.reduce((sum, hitCount, index) => {
                    if (activeReportIndices.has(index)) {
                        return sum + hitCount
                    }
                    return sum
                }, 0) ?? 0

            let status: LineStatus
            if (totalHits > 0) {
                status = line.branchInfo ? 'partial' : 'covered'
            } else {
                status = 'uncovered'
            }

            return { ...line, hits: totalHits, status }
        })
    }, [lines, activeReportIndices])

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
