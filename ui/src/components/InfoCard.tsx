import type { MetadataItem } from '@/types/summary'
import { Card } from '@/ui/card'

const InfoRow = ({ label, children }: { label: string; children: React.ReactNode }) => (
    <div className="flex w-full items-baseline justify-between gap-2 text-xs">
        <span className="shrink-0 text-muted-foreground">{label}</span>
        {children}
    </div>
)

const ValueDisplay = ({ value }: { value: MetadataItem['value'] }) => {
    if (value === undefined || value === null || (Array.isArray(value) && value.length === 0)) {
        return <span className="font-medium font-mono text-foreground">-</span>
    }

    const displayString = Array.isArray(value) ? value.join(', ') : String(value)

    return (
        <span className="truncate text-right font-medium font-mono text-foreground" title={displayString}>
            {displayString}
        </span>
    )
}

interface InfoCardProps {
    title: string
    items: MetadataItem[]
}

export default function InfoCard({ title, items }: InfoCardProps) {
    return (
        <Card className="flex w-full flex-col gap-2 rounded-md px-3 py-2.5">
            <span className="font-medium text-muted-foreground text-xs uppercase tracking-wide">{title}</span>
            <div className="flex flex-col gap-0.5">
                {items.map((item) => (
                    <InfoRow key={item.label} label={item.label}>
                        <ValueDisplay value={item.value} />
                    </InfoRow>
                ))}
            </div>
        </Card>
    )
}
