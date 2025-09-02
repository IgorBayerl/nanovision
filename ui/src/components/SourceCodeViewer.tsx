import CodeLine from '@/components/CodeLine'
import SourceCodeHeader from '@/components/SourceCodeHeader'
import type { DetailsV1 } from '@/lib/validation'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

interface SourceCodeViewerProps {
    fileName: string
    lines: DetailsV1['lines']
}

export default function SourceCodeViewer({ fileName, lines }: SourceCodeViewerProps) {
    // The virtualizer and its related logic have been completely removed.
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
                            {lines.map((line) => (
                                <CodeLine key={line.lineNumber} {...line} />
                            ))}
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
