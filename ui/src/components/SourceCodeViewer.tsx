import CodeLine from '@/components/CodeLine'
import type { DetailsV1 } from '@/lib/validation'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

interface SourceCodeViewerProps {
    fileName: string
    lines: DetailsV1['lines']
}

export default function SourceCodeViewer({ fileName, lines }: SourceCodeViewerProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle className="font-mono text-lg">{fileName}</CardTitle>
            </CardHeader>
            <CardContent className="overflow-x-auto rounded-md border border-border bg-subtle p-0">
                <div className="relative font-mono">
                    {lines.length > 0 ? (
                        lines.map((line) => <CodeLine key={line.lineNumber} {...line} />)
                    ) : (
                        <div className="p-4 text-center text-muted-foreground">No source code to display.</div>
                    )}
                </div>
            </CardContent>
        </Card>
    )
}
