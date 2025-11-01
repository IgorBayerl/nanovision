import type { Report } from '@/types/summary'
import { Card, CardContent, CardHeader, CardTitle } from '@/ui/card'
import { Checkbox } from '@/ui/checkbox'
import { Label } from '@/ui/label'

interface ReportsSelectorProps {
    reports: Report[]
    activeReportIndices: Set<number>
    onToggleReport: (index: number) => void
}

export default function ReportsSelector({ reports, activeReportIndices, onToggleReport }: ReportsSelectorProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Reports</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-2">
                    {reports.map((report, index) => (
                        <div key={index} className="flex items-center space-x-2">
                            <Checkbox
                                id={`report-checkbox-${index}`}
                                checked={activeReportIndices.has(index)}
                                onCheckedChange={() => onToggleReport(index)}
                            />
                            <Label htmlFor={`report-checkbox-${index}`} title={report.path}>
                                {report.name}
                            </Label>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    )
}