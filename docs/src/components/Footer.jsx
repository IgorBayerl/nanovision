import React from 'react';
import { Hexagon, Github } from 'lucide-react';

const Footer = () => {
    const base = import.meta.env.BASE_URL.endsWith('/') ? import.meta.env.BASE_URL.slice(0, -1) : import.meta.env.BASE_URL;

    return (
        <footer className="border-t border-border bg-background py-12">
            <div className="max-w-7xl mx-auto px-6 flex flex-col md:flex-row justify-between items-center gap-6">
                <div className="flex items-center gap-2">
                    <img src={`${base}/nanovision.png`} alt="Nanovision Logo" className="w-6 h-6" />
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
};

export default Footer;
