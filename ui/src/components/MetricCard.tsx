import { AlertCircle, AlertTriangle, ShieldCheck } from 'lucide-react'
import type { CoverageDetail, MetricDefinition, RiskLevel } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'
import { Progress } from '@/ui/progress'

const DetailRow = ({ label, value }: { label: string; value: number | string }) => (
    <div className="group flex items-baseline justify-between px-2 text-sm hover:bg-accent/50">
        <span className="text-muted-foreground group-hover:text-foreground">{label}</span>
        <span className="font-medium font-mono text-foreground">{value}</span>
    </div>
)

export const StatusIcon = ({ status }: { status: RiskLevel }) => {
    if (status === 'danger') return <AlertCircle className="h-4 w-4 text-uncovered" />
    if (status === 'warning') return <AlertTriangle className="h-4 w-4 text-partial" />
    return <ShieldCheck className="h-4 w-4 text-primary" />
}

export default function MetricCard({
    label,
    details,
    status,
    definition,
}: {
    label: string
    details: CoverageDetail | undefined
    status?: RiskLevel
    definition?: MetricDefinition
}) {
    const pct = details ? Math.max(0, Math.min(100, Math.round(details.percentage))) : undefined

    return (
        <Card className="flex h-full w-full flex-col rounded-md">
            <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="text-lg">{label}</CardTitle>
                {status && <StatusIcon status={status} />}
            </CardHeader>
            <CardContent className="flex flex-grow gap-4">
                <div className="flex flex-col items-center">
                    <div className="font-bold text-5xl text-foreground tabular-nums tracking-tight">
                        {pct !== undefined ? `${pct}%` : 'N/A'}
                    </div>
                    <div className="mt-2 w-full max-w-[100px]">
                        <Progress value={pct} indicatorClassName="bg-primary" />
                    </div>
                </div>
                <div className="flex flex-1 flex-col divide-y">
                    {details && definition ? (
                        definition.subMetrics
                            .filter((sub) => sub.id !== 'percentage')
                            .map((subMetric) => {
                                const value = details[subMetric.id as keyof typeof details]
                                return <DetailRow key={subMetric.id} label={subMetric.label} value={value ?? '-'} />
                            })
                    ) : (
                        <div className="flex h-full items-center justify-center text-muted-foreground">No data</div>
                    )}
                </div>
            </CardContent>
        </Card>
    )
}
