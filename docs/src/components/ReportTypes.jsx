import React from 'react';
import { Layout, Terminal, FileCode2, FileJson, ArrowUpRight } from 'lucide-react';

const ReportTypes = () => {
    const formats = [
        {
            id: "Html",
            title: "HTML Report",
            icon: Layout,
            desc: "Interactive, searchable web dashboard. Visualize complexity, file trees, and coverage heatmaps.",
            useCase: "Human analysis, distribution to stakeholders.",
            exampleLink: "https://igorbayerl.github.io/nanovision/reports/"
        },
        {
            id: "TextSummary",
            title: "Text Summary",
            icon: Terminal,
            desc: "Compact ASCII table output printed directly to stdout.",
            useCase: "CI pipeline logs, quick sanity checks in terminal.",
            exampleLink: "https://igorbayerl.github.io/nanovision/reports/Summary.txt"
        },
        {
            id: "Lcov",
            title: "LCOV (.info)",
            icon: FileCode2,
            desc: "Standard tracefile format compatible with legacy tools.",
            useCase: "Uploading to Coveralls, Codecov, or IDE extensions.",
            exampleLink: "https://igorbayerl.github.io/nanovision/reports/lcov.info"
        },
        {
            id: "RawJson",
            title: "Raw JSON",
            icon: FileJson,
            desc: "Full dump of the internal metrics tree structure.",
            useCase: "Custom post-processing, archiving metrics for historical analysis.",
            exampleLink: "https://igorbayerl.github.io/nanovision/reports/RawJson.json"
        }
    ];

    return (
        <section id="reports" className="py-24 border-t border-border bg-secondary/20">
            <div className="max-w-7xl mx-auto px-6">
                <div className="text-center mb-16">
                    <h2 className="text-3xl font-bold text-foreground mb-4">Multiple Output Formats</h2>
                    <p className="text-muted-foreground max-w-2xl mx-auto">
                        Ingest reports from multiple instrumenters and output standardized data for any pipeline.
                    </p>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    {formats.map(f => {
                        const Icon = f.icon;
                        return (
                            <div key={f.id} className="bg-card border border-border p-6 rounded-xl hover:border-primary/50 group h-full flex flex-col">
                                <div className="w-12 h-12 rounded-lg bg-secondary flex items-center justify-center text-primary mb-4 group-hover:scale-110 transition-transform">
                                    <Icon size={24} />
                                </div>
                                <h3 className="font-bold text-foreground mb-2">{f.title}</h3>
                                <p className="text-sm text-muted-foreground mb-4 flex-1 leading-relaxed">{f.desc}</p>
                                <div className="mt-auto pt-4 border-t border-border">
                                    <p className="text-xs font-semibold text-primary uppercase mb-1">Best For</p>
                                    <p className="text-xs text-muted-foreground mb-3">{f.useCase}</p>
                                    <a href={f.exampleLink} onClick={(e) => e.preventDefault()} className="inline-flex items-center gap-1 text-xs font-bold text-foreground hover:text-primary">
                                        View Example <ArrowUpRight size={12} />
                                    </a>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>
        </section>
    );
};

export default ReportTypes;
