import React, { useEffect, useState } from 'react';
import { Hexagon, Sun, Moon } from 'lucide-react';

const Navbar = ({ currentPath = "/" }) => {
    const [isDark, setIsDark] = useState(false);

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
        if (path === '/') return currentPath === '/';
        return currentPath.startsWith(path);
    };

    return (
        <nav className="sticky top-0 w-full z-50 border-b border-border bg-background/80 backdrop-blur-md">
            <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
                <a href="/" className="flex items-center gap-3 hover:opacity-80">
                    <div className="text-primary">
                        <Hexagon size={28} className="fill-current opacity-20" />
                    </div>
                    <span className="font-bold text-xl tracking-tight font-mono text-foreground">NANOVISION</span>
                </a>

                <div className="hidden md:flex items-center gap-8 text-sm font-medium">
                    <a
                        href="/"
                        className={`${isActive('/') && currentPath === '/' ? 'text-primary font-bold' : 'text-muted-foreground hover:text-foreground'}`}
                    >
                        Home
                    </a>
                    <a
                        href="/getting-started"
                        className={`${isActive('/getting-started') ? 'text-primary font-bold' : 'text-muted-foreground hover:text-foreground'}`}
                    >
                        Getting Started
                    </a>
                    <a
                        href="/usage"
                        className={`${isActive('/usage') ? 'text-primary font-bold' : 'text-muted-foreground hover:text-foreground'}`}
                    >
                        Usage
                    </a>
                    <a
                        href="/configurator"
                        className={`${isActive('/configurator') ? 'text-primary font-bold' : 'text-muted-foreground hover:text-foreground'}`}
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
                    <a href="/getting-started" className="bg-primary text-primary-foreground px-5 py-2 text-sm font-semibold hover:opacity-90 rounded-md shadow-sm">
                        Get Started
                    </a>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
