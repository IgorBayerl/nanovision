import React from 'react';
import { Hexagon, Github } from 'lucide-react';

const Footer = () => (
    <footer className="border-t border-border bg-background py-12">
        <div className="max-w-7xl mx-auto px-6 flex flex-col md:flex-row justify-between items-center gap-6">
            <div className="flex items-center gap-2">
                <div className="text-primary">
                    <Hexagon size={20} />
                </div>
                <span className="font-bold text-foreground font-mono tracking-wider">NANOVISION</span>
            </div>
            <div className="text-sm text-muted-foreground">
                &copy; {new Date().getFullYear()} Open Source under Apache 2.0.
            </div>
            <div className="flex items-center gap-6">
                <a href="https://github.com/IgorBayerl/nanovision" className="text-muted-foreground hover:text-foreground">
                    <Github size={20} />
                </a>
            </div>
        </div>
    </footer>
);

export default Footer;
