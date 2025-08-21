import type { CoverageDetail } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

const DetailRow = ({ label, value }: { label: string; value: number | string }) => (
    <div className="group flex items-baseline justify-between px-2 text-sm hover:bg-muted/50">
        <span className="text-muted-foreground group-hover:text-foreground">{label}</span>
        <span className="font-medium font-mono text-foreground">{value}</span>
    </div>
)

export default function MetricCard({ label, details }: { label: string; details: CoverageDetail | undefined }) {
    const pct = details ? Math.max(0, Math.min(100, Math.round(details.percentage))) : undefined

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-lg">{label}</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-4">
                <div className="flex flex-col items-center justify-center">
                    <div className="font-extrabold text-5xl text-foreground tabular-nums tracking-tight">
                        {pct !== undefined ? `${pct}%` : 'N/A'}
                    </div>
                    <div className="relative mt-2 h-2 w-full max-w-[100px] overflow-hidden bg-muted">
                        <div
                            className="absolute inset-y-0 left-0 bg-primary"
                            style={{ width: pct !== undefined ? `${pct}%` : '0%' }}
                        />
                    </div>
                </div>
                <div className="flex flex-col divide-y">
                    {details ? (
                        <>
                            <DetailRow label="Covered" value={details.covered} />
                            <DetailRow label="Uncovered" value={details.uncovered} />
                            <DetailRow label="Coverable" value={details.coverable} />
                            <DetailRow label="Total" value={details.total} />
                        </>
                    ) : (
                        <div className="flex h-full items-center justify-center text-muted-foreground">No data</div>
                    )}
                </div>
            </CardContent>
        </Card>
    )
}
