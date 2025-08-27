import InlineCoverage from '@/components/InlineCoverage'
import type { RiskLevel } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

export type IslandProps = Record<string, unknown>

interface FileDetailsProps {
    fileName?: string
    coverage?: number
    risk?: RiskLevel
}

export function FileDetailsIsland(props: IslandProps) {
    const { fileName, coverage, risk } = props as FileDetailsProps

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-lg">
                    Interactivity for: <span className="font-mono">{fileName ?? 'Unknown File'}</span>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <p className="text-sm text-muted-foreground">
                    This component is a React "island" that was hydrated on the client. The data below was passed via a{' '}
                    <code>data-props</code> attribute.
                </p>
                <div className="mt-4">
                    {typeof coverage === 'number' && typeof risk === 'string' ? (
                        <InlineCoverage percentage={coverage} risk={risk} />
                    ) : (
                        <p>No coverage data provided.</p>
                    )}
                </div>
            </CardContent>
        </Card>
    )
}
