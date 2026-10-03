import { CheckCircle2, XCircle } from 'lucide-react'
import InfoTooltip from '@/components/InfoTooltip'
import { cn } from '@/lib/utils'
import type { Review } from '@/lib/validation'

interface ReviewSummaryProps {
    review: Review
}

const formatPct = (value: number) => `${value.toFixed(1)}%`

/**
 * The verdict on the changed code: gate banner and changelist stat cards. A
 * report carries a `review` block when its run was measured with a diff.
 */
export default function ReviewSummary({ review }: ReviewSummaryProps) {
    const { stats, checks = [] } = review
    const failedChecks = checks.filter((c) => !c.passed)

    const patchPct =
        stats.patchStatementsValid > 0 ? (100 * stats.patchStatementsCovered) / stats.patchStatementsValid : null
    const patchCheck = checks.find((c) => c.key === 'patch_statement_coverage')
    const complexityCheck = checks.find((c) => c.key === 'max_changed_method_complexity')
    const methodsTouched = stats.methodsAdded + stats.methodsModified

    return (
        <div className="flex flex-col gap-4">
            {checks.length > 0 && (
                <div
                    className={cn(
                        'flex items-start gap-2 rounded-md border p-3 text-sm',
                        review.passed
                            ? 'border-covered/40 bg-covered/10 text-covered'
                            : 'border-uncovered/40 bg-uncovered/10 text-uncovered',
                    )}
                >
                    {review.passed ? (
                        <CheckCircle2 className="mt-0.5 h-4 w-4 shrink-0" />
                    ) : (
                        <XCircle className="mt-0.5 h-4 w-4 shrink-0" />
                    )}
                    <div>
                        <div className="flex items-center gap-1.5 font-semibold">
                            {review.passed ? 'Review gate passed' : 'Review gate failed'}
                            <InfoTooltip label="How the gate is evaluated">
                                The gate and the changelist numbers are evaluated against all reports. The report
                                selection does not change them.
                            </InfoTooltip>
                        </div>
                        {failedChecks.length > 0 && (
                            <ul className="mt-1 list-inside list-disc">
                                {failedChecks.map((check) => (
                                    <li key={check.key}>
                                        {check.label} is {check.value.toFixed(1)} (limit {check.threshold.toFixed(1)})
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            )}

            <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
                <StatCard
                    label="Patch statement coverage"
                    value={patchPct === null ? 'n/a' : formatPct(patchPct)}
                    hint={
                        stats.patchStatementsValid > 0
                            ? `${stats.patchStatementsCovered} / ${stats.patchStatementsValid} statements`
                            : 'no coverable changed statements'
                    }
                    tone={patchCheck ? (patchCheck.passed ? 'good' : 'bad') : 'neutral'}
                />
                <StatCard label="Changed files" value={String(stats.changedFiles)} hint="in this changelist" />
                <StatCard
                    label="Methods touched"
                    value={String(methodsTouched)}
                    hint={`${stats.methodsAdded} added, ${stats.methodsModified} modified`}
                />
                <StatCard
                    label="Untested changed methods"
                    value={String(stats.untestedChangedMethods)}
                    hint="changed code, zero coverage"
                    tone={stats.untestedChangedMethods > 0 ? 'bad' : 'good'}
                />
                <StatCard
                    label="Max changed complexity"
                    value={String(stats.maxChangedComplexity)}
                    hint="highest CC among changed methods"
                    tone={complexityCheck ? (complexityCheck.passed ? 'good' : 'bad') : 'neutral'}
                />
            </div>
        </div>
    )
}

function StatCard({
    label,
    value,
    hint,
    tone = 'neutral',
}: {
    label: string
    value: string
    hint?: string
    tone?: 'neutral' | 'good' | 'bad'
}) {
    return (
        <div className="rounded-md border border-border bg-card p-3">
            <div className="text-muted-foreground text-xs">{label}</div>
            <div
                className={cn(
                    'font-semibold text-xl tabular-nums',
                    tone === 'good' && 'text-covered',
                    tone === 'bad' && 'text-uncovered',
                )}
            >
                {value}
            </div>
            {hint && <div className="mt-0.5 text-[11px] text-muted-foreground">{hint}</div>}
        </div>
    )
}
