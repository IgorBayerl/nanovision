import { AlertCircle } from 'lucide-react'
import type { ZodIssue } from 'zod'

interface ValidationAlertsProps {
    issues: ZodIssue[]
}

/**
 * A component to display data validation errors.
 */
export default function ValidationAlerts({ issues }: ValidationAlertsProps) {
    if (!issues || issues.length === 0) {
        return null
    }

    return (
        <div className="rounded-md border border-destructive/50 bg-destructive/10 p-4 text-destructive">
            <div className="flex items-start gap-3">
                <div className="flex-shrink-0">
                    <AlertCircle className="h-5 w-5" aria-hidden="true" />
                </div>
                <div className="flex-1">
                    <h3 className="font-bold text-destructive">Invalid Report Data</h3>
                    <div className="mt-2 text-sm">
                        <p>
                            There {issues.length === 1 ? 'is an issue' : 'are issues'} with the data structure of this
                            report. The content below may be incomplete or inaccurate.
                        </p>
                        <ul className="mt-3 list-disc space-y-1 pl-5">
                            {issues.slice(0, 5).map((issue) => (
                                <li key={`${issue.path.join('-')}-${issue.code}`}>
                                    <span className="font-semibold">{issue.path.join('.')}</span>: {issue.message}
                                </li>
                            ))}
                        </ul>
                        {issues.length > 5 && (
                            // Corrected: Sorted Tailwind CSS classes
                            <p className="mt-2 font-medium text-xs">...and {issues.length - 5} more issues.</p>
                        )}
                    </div>
                </div>
            </div>
        </div>
    )
}
