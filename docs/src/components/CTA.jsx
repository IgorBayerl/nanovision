import React from 'react';
import { Sliders, ArrowUpRight, BookOpen } from 'lucide-react';

const CTA = () => (
    <section className="py-32 relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-transparent via-primary/5 to-transparent"></div>

        <div className="max-w-5xl mx-auto px-6 relative z-10">
            <h2 className="text-4xl font-bold mb-16 text-foreground text-center tracking-tight">READY TO DEPLOY?</h2>

            <div className="grid md:grid-cols-2 gap-8">
                <a href={`${import.meta.env.BASE_URL}configurator`} className="group bg-card border border-border hover:border-primary/50 p-8 hover:-translate-y-1 shadow-lg rounded-xl transition-all duration-300">
                    <div className="flex justify-between items-start mb-6">
                        <div className="text-primary"><Sliders size={32} /></div>
                        <ArrowUpRight size={20} className="text-muted-foreground group-hover:text-primary transition-colors duration-300" />
                    </div>
                    <h3 className="text-xl font-bold text-foreground mb-2">I already track coverage</h3>
                    <p className="text-muted-foreground mb-6 text-sm">Use the interactive configurator to generate your CLI command or YAML config file.</p>
                    <span className="text-primary text-sm font-bold tracking-wide uppercase group-hover:underline">Configure Now</span>
                </a>

                <a href={`${import.meta.env.BASE_URL}getting-started`} className="group bg-card border border-border hover:border-primary/50 p-8 hover:-translate-y-1 shadow-lg rounded-xl transition-all duration-300">
                    <div className="flex justify-between items-start mb-6">
                        <div className="text-foreground"><BookOpen size={32} /></div>
                        <ArrowUpRight size={20} className="text-muted-foreground group-hover:text-foreground transition-colors duration-300" />
                    </div>
                    <h3 className="text-xl font-bold text-foreground mb-2">I'm new to coverage</h3>
                    <p className="text-muted-foreground mb-6 text-sm">Learn how to instrument your Go, C++, or C# code to output formats Nanovision can read.</p>
                    <span className="text-foreground text-sm font-bold tracking-wide uppercase group-hover:underline">Read Setup Guide</span>
                </a>
            </div>
        </div>
    </section>
);

export default CTA;
