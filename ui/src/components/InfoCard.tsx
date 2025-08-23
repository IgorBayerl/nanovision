import type { MetadataItem } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

const InfoRow = ({ label, children }: { label: string; children: React.ReactNode }) => (
    <div className="group flex items-baseline  w-full justify-between text-sm hover:bg-accent/50">
        <span className="text-muted-foreground group-hover:text-foreground">{label}:</span>
        {children}
    </div>
)

const ValueDisplay = ({ value }: { value: MetadataItem['value'] }) => {
    if (value === undefined || value === null || (Array.isArray(value) && value.length === 0)) {
        return <span className="font-medium font-mono text-foreground">-</span>
    }

    const displayString = Array.isArray(value) ? value.join(', ') : String(value)

    return (
        <span className="font-medium font-mono text-foreground" title={displayString}>
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
        <Card className="flex h-full flex-col rounded-md">
            <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="text-lg">{title}</CardTitle>
            </CardHeader>
            <CardContent className="flex-grow">
                <div className="flex flex-col flex-wrap content-start gap-x-6 divide-y border-border">
                    {items.map((item) => (
                        <InfoRow key={item.label} label={item.label}>
                            <ValueDisplay value={item.value} />
                        </InfoRow>
                    ))}
                </div>
            </CardContent>
        </Card>
    )
}
