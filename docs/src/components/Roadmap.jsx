import React from 'react';
import { CheckCircle, Circle } from 'lucide-react';

const Roadmap = () => (
    <section id="roadmap" className="py-24 bg-background relative">
        <div className="max-w-7xl mx-auto px-6">
            <div className="text-center mb-16">
                <h2 className="text-3xl font-bold text-foreground mb-4">Mission Roadmap</h2>
                <p className="text-muted-foreground">The future of Nanovision involves deeper language support and tighter CI integration.</p>
            </div>

            <div className="grid md:grid-cols-3 gap-6">
                {/* Q1 */}
                <div className="bg-card border border-border p-6 relative overflow-hidden group rounded-xl shadow-sm">
                    <div className="absolute top-0 right-0 p-10 bg-primary/5 blur-2xl group-hover:bg-primary/10 transition-colors duration-500"></div>
                    <div className="text-xs font-mono text-primary mb-4 font-bold">PHASE 1: EXPANSION</div>
                    <ul className="space-y-3 text-sm text-muted-foreground">
                        <li className="flex items-center gap-2"><CheckCircle size={14} className="text-green-500" /> C# Static Analysis</li>
                        <li className="flex items-center gap-2"><Circle size={14} className="" /> TS/JS Static Analysis</li>
                        <li className="flex items-center gap-2"><Circle size={14} className="" /> GitLab CI Native Integration</li>
                    </ul>
                </div>

                {/* Q2 */}
                <div className="bg-card border border-border p-6 relative overflow-hidden group rounded-xl shadow-sm">
                    <div className="absolute top-0 right-0 p-10 bg-secondary/10 blur-2xl group-hover:bg-secondary/20 transition-colors duration-500"></div>
                    <div className="text-xs font-mono text-foreground mb-4 font-bold">PHASE 2: INTELLIGENCE</div>
                    <ul className="space-y-3 text-sm text-muted-foreground">
                        <li className="flex items-center gap-2"><Circle size={14} className="" /> Trend Analysis (History)</li>
                        <li className="flex items-center gap-2"><Circle size={14} className="" /> Flaky Test Detection</li>
                        <li className="flex items-center gap-2"><Circle size={14} className="" /> Risk Heatmaps</li>
                    </ul>
                </div>

                {/* Q3 */}
                <div className="bg-card border border-border p-6 relative overflow-hidden group rounded-xl shadow-sm">
                    <div className="absolute top-0 right-0 p-10 bg-accent/5 blur-2xl group-hover:bg-accent/10 transition-colors duration-500"></div>
                    <div className="text-xs font-mono text-foreground mb-4 font-bold">PHASE 3: ECOSYSTEM</div>
                    <ul className="space-y-3 text-sm text-muted-foreground">
                        <li className="flex items-center gap-2"><Circle size={14} className="" /> VS Code Extension</li>
                        <li className="flex items-center gap-2"><Circle size={14} className="" /> IntelliJ Plugin</li>
                        <li className="flex items-center gap-2"><Circle size={14} className="" /> Cloud Dashboard</li>
                    </ul>
                </div>
            </div>
        </div>
    </section>
);

export default Roadmap;
