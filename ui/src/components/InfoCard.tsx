import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'

const InfoRow = ({ label, children }: { label: string; children: React.ReactNode }) => (
    <div className="flex flex-col gap-1 border-border border-t pt-3 sm:flex-row sm:justify-between sm:gap-4">
        <dt className="font-semibold text-foreground">{label}</dt>
        <dd className="text-muted-foreground sm:text-right">{children}</dd>
    </div>
)

const ListValue = ({ items }: { items?: string[] }) => {
    if (!items || items.length === 0) return <span>-</span>
    return (
        <ul className="space-y-1">
            {items.map((item) => (
                <li key={item} className="truncate font-mono" title={item}>
                    {item}
                </li>
            ))}
        </ul>
    )
}

type InfoCardProps = {
    generatedAt: string
    files: number
    folders: number
    parsers?: string[]
    configFiles?: string[]
    importedReports?: string[]
}

export default function InfoCard({
    generatedAt,
    files,
    folders,
    parsers,
    configFiles,
    importedReports,
}: InfoCardProps) {
    const formattedDate = new Date(generatedAt).toLocaleString()

    return (
        <Card className="flex h-full flex-col rounded-md">
            <CardHeader>
                <CardTitle className="text-lg">Report Information</CardTitle>
            </CardHeader>
            <CardContent className="flex-grow">
                <dl className="flex flex-col gap-3 text-sm">
                    <InfoRow label="Generated At">{formattedDate}</InfoRow>
                    <InfoRow label="Total Files">{files}</InfoRow>
                    <InfoRow label="Total Folders">{folders}</InfoRow>
                    <InfoRow label="Parsers">
                        <ListValue items={parsers} />
                    </InfoRow>
                    <InfoRow label="Config Files">
                        <ListValue items={configFiles} />
                    </InfoRow>
                    <InfoRow label="Imported Reports">
                        <ListValue items={importedReports} />
                    </InfoRow>
                </dl>
            </CardContent>
        </Card>
    )
}
