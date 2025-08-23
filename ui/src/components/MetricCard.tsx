import { AlertCircle, AlertTriangle, ShieldCheck } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { CoverageDetail, MetricDefinition, RiskLevel } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

const DetailRow = ({ label, value }: { label: string; value: number | string }) => (
    <div className="group flex items-baseline justify-between px-2 text-sm hover:bg-muted/50">
        <span className="text-muted-foreground group-hover:text-foreground">{label}</span>
        <span className="font-medium font-mono text-foreground">{value}</span>
    </div>
)

const StatusIcon = ({ status }: { status: RiskLevel }) => {
    if (status === 'danger') return <AlertCircle className="h-4 w-4 text-destructive" />
    if (status === 'warning') return <AlertTriangle className="h-4 w-4 text-accent" />
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

    const statusBarColor = status === 'danger' ? 'bg-destructive' : status === 'warning' ? 'bg-accent' : 'bg-primary'

    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="text-lg">{label}</CardTitle>
                {status && <StatusIcon status={status} />}
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-4">
                <div className="flex flex-col items-center justify-center">
                    <div className="font-extrabold text-5xl text-foreground tabular-nums tracking-tight">
                        {pct !== undefined ? `${pct}%` : 'N/A'}
                    </div>
                    <div className="relative mt-2 h-2 w-full max-w-[100px] overflow-hidden bg-muted">
                        <div
                            className={cn('absolute inset-y-0 left-0', statusBarColor)}
                            style={{ width: pct !== undefined ? `${pct}%` : '0%' }}
                        />
                    </div>
                </div>
                <div className="flex flex-col divide-y">
                    {details && definition ? (
                        definition.subMetrics
                            .filter((sub) => sub.id !== 'percentage') // The percentage is already shown on the left
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
