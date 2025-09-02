const gridTemplateColumns = '1.5rem 4rem 4rem 1.5rem 1.5rem 1fr'

export default function SourceCodeHeader() {
    return (
        <div
            className="sticky top-0 z-20 grid select-none border-border border-b bg-subtle font-semibold text-muted-foreground text-xs"
            style={{ gridTemplateColumns }}
        >
            {/* Status Color */}
            <div className="py-2 text-center" />

            {/* Line Number */}
            <div className="border-border/50 border-r py-2 pr-4 text-right">#</div>

            {/* Hit Count */}
            <div className="border-border/50 border-r px-4 py-2 text-center">Hits</div>

            {/* Branch Indicator */}
            <div className="py-2 text-center" />

            {/* Diff Indicator */}
            <div className="border-border/50 border-r py-2" />

            {/* Source Code */}
            <div className="py-2 pl-4">Source</div>
        </div>
    )
}
