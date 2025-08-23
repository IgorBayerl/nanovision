import { cn } from '@/lib/utils'
import type { MetadataItem } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

const InfoRow = ({ label, children }: { label: string; children: React.ReactNode }) => (
    <div className="flex flex-col gap-1 border-border border-t sm:flex-row sm:justify-between sm:gap-4">
        <dt className="font-semibold text-foreground">{label}:</dt>
        <dd className="text-muted-foreground sm:text-right">{children}</dd>
    </div>
)

const ValueDisplay = ({ value }: { value: MetadataItem['value'] }) => {
    if (value === undefined || value === null || (Array.isArray(value) && value.length === 0)) {
        return <span>-</span>
    }

    if (Array.isArray(value)) {
        const fullString = value.join(', ')
        return (
            <span className="font-mono" title={fullString}>
                {fullString}
            </span>
        )
    }

    return <span className="font-mono">{String(value)}</span>
}

interface InfoCardProps {
    title: string
    items: MetadataItem[]
}

export default function InfoCard({ title, items }: InfoCardProps) {
    return (
        <Card className="flex h-full flex-col rounded-md">
            <CardHeader>
                <CardTitle className="text-lg">{title}</CardTitle>
            </CardHeader>
            <CardContent className="flex-grow">
                <div className="flex flex-col flex-wrap content-start gap-x-6">
                    {items.map((item, _index) => (
                        <div key={item.label} className={cn('group hover:bg-muted/50')}>
                            <InfoRow label={item.label}>
                                <ValueDisplay value={item.value} />
                            </InfoRow>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    )
}
