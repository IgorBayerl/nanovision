import React, { useEffect, useState } from 'react';
import { Hexagon, Sun, Moon } from 'lucide-react';

const Navbar = ({ currentPath = "/" }) => {
    const [isDark, setIsDark] = useState(false);
    const base = import.meta.env.BASE_URL.endsWith('/') ? import.meta.env.BASE_URL.slice(0, -1) : import.meta.env.BASE_URL;

    useEffect(() => {
        // Initialize state from DOM or localStorage
        const isDarkMode = document.documentElement.classList.contains('dark');
        setIsDark(isDarkMode);
    }, []);

    const toggleTheme = () => {
        const newIsDark = !isDark;
        setIsDark(newIsDark);

        if (newIsDark) {
            document.documentElement.classList.add('dark');
            localStorage.setItem('theme', 'dark');
        } else {
            document.documentElement.classList.remove('dark');
            localStorage.setItem('theme', 'light');
        }
    };

    const isActive = (path) => {
        // Normalize currentPath to remove base for comparison if needed, 
        // but simpler to just check if currentPath ends with the target path
        // or matches the full path.
        // Let's just use the path as is since we are in a static site.
        if (path === '/') return currentPath === base || currentPath === base + '/';
        return currentPath.startsWith(`${base}${path}`);
    };

    const getLink = (path) => {
        if (path === '/') return `${base}/`;
        return `${base}${path}`;
    }

    return (
        <nav className="sticky top-0 w-full z-50 border-b border-border bg-background/80 backdrop-blur-md">
            <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
                <a href={getLink('/')} className="flex items-center gap-3 hover:opacity-80">
                    <img src={`${base}/nanovision.png`} alt="Nanovision Logo" className="h-12 w-auto" />
                    <span className="font-bold text-xl tracking-tight font-mono text-foreground">Nanovision</span>
                </a>

                <div className="hidden md:flex items-center gap-8 text-sm font-medium">
                    <a
                        href={getLink('/')}
                        className={`${isActive('/') ? 'text-primary' : 'text-muted-foreground hover:text-foreground'}`}
                    >
                        Home
                    </a>
                    <a
                        href={getLink('/getting-started')}
                        className={`${isActive('/getting-started') ? 'text-primary' : 'text-muted-foreground hover:text-foreground'}`}
                    >
                        Getting Started
                    </a>
                    <a
                        href={getLink('/usage')}
                        className={`${isActive('/usage') ? 'text-primary' : 'text-muted-foreground hover:text-foreground'}`}
                    >
                        Usage
                    </a>
                    <a
                        href={getLink('/configurator')}
                        className={`${isActive('/configurator') ? 'text-primary' : 'text-muted-foreground hover:text-foreground'}`}
                    >
                        Configurator
                    </a>
                </div>

                <div className="flex items-center gap-4">
                    <button
                        onClick={toggleTheme}
                        className="p-2 rounded-md text-muted-foreground hover:bg-secondary cursor-pointer"
                        title="Toggle Theme"
                    >
                        {isDark ? <Sun size={20} /> : <Moon size={20} />}
                    </button>
                    <a href={getLink('/getting-started')} className="bg-primary text-primary-foreground px-5 py-2 text-sm font-semibold hover:opacity-90 rounded-md shadow-sm">
                        Get Started
                    </a>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
