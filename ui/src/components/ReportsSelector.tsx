import type { Report } from '@/types/summary'
import { Card } from '@/ui/card'
import { Checkbox } from '@/ui/checkbox'
import { Label } from '@/ui/label'

interface ReportsSelectorProps {
    reports: Report[]
    activeReportIndices: Set<number>
    onToggleReport: (index: number) => void
}

export default function ReportsSelector({ reports, activeReportIndices, onToggleReport }: ReportsSelectorProps) {
    return (
        <Card className="flex w-full flex-col gap-2 rounded-md px-3 py-2.5">
            <span className="font-medium text-muted-foreground text-xs uppercase tracking-wide">Reports</span>
            <div className="flex flex-col gap-1.5">
                {reports.map((report, index) => (
                    <div key={report.path} className="flex items-center space-x-2">
                        <Checkbox
                            id={`report-checkbox-${report.path}`}
                            checked={activeReportIndices.has(index)}
                            onCheckedChange={() => onToggleReport(index)}
                        />
                        <Label
                            htmlFor={`report-checkbox-${report.path}`}
                            title={report.path}
                            className="truncate text-xs"
                        >
                            {report.name}
                        </Label>
                    </div>
                ))}
            </div>
        </Card>
    )
}
