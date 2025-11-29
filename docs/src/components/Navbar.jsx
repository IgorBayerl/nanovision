import { useEffect, useState } from 'react';
import { Sun, Moon, Github } from 'lucide-react';

const Navbar = ({ currentPath = "/" }) => {
    const [isDark, setIsDark] = useState(false);
    const [stars, setStars] = useState(null);
    const base = import.meta.env.BASE_URL.endsWith('/')
        ? import.meta.env.BASE_URL.slice(0, -1)
        : import.meta.env.BASE_URL;

    useEffect(() => {
        const isDarkMode = document.documentElement.classList.contains('dark');
        setIsDark(isDarkMode);

        fetch('https://api.github.com/repos/IgorBayerl/nanovision')
            .then((res) => (res.ok ? res.json() : null))
            .then((data) => {
                if (data && typeof data.stargazers_count === 'number') {
                    setStars(data.stargazers_count);
                }
            })
            .catch(() => {
                // ignore errors
            });
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
        if (path === '/') return currentPath === base || currentPath === base + '/';
        return currentPath.startsWith(`${base}${path}`);
    };

    const getLink = (path) => {
        if (path === '/') return `${base}/`;
        return `${base}${path}`;
    };

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
                        className={isActive('/') ? 'text-primary' : 'text-muted-foreground hover:text-foreground'}
                    >
                        Home
                    </a>
                    <a
                        href={getLink('/about-coverage')}
                        className={
                            isActive('/about-coverage')
                                ? 'text-primary'
                                : 'text-muted-foreground hover:text-foreground'
                        }
                    >
                        About Coverage
                    </a>
                    <a
                        href={getLink('/getting-started')}
                        className={
                            isActive('/getting-started')
                                ? 'text-primary'
                                : 'text-muted-foreground hover:text-foreground'
                        }
                    >
                        Getting Started
                    </a>
                    <a
                        href={getLink('/configurator')}
                        className={
                            isActive('/configurator')
                                ? 'text-primary'
                                : 'text-muted-foreground hover:text-foreground'
                        }
                    >
                        Configurator
                    </a>
                </div>

                <div className="flex items-center gap-3">

                    <button
                        onClick={toggleTheme}
                        className="p-2 rounded-md text-muted-foreground hover:bg-secondary cursor-pointer"
                        title="Toggle Theme"
                    >
                        {isDark ? <Sun size={20} /> : <Moon size={20} />}
                    </button>

                    <a
                        href="https://github.com/IgorBayerl/nanovision"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-2 px-3 py-1.5 rounded-md border border-border bg-card text-xs font-mono text-muted-foreground hover:text-foreground hover:border-primary/60 hover:bg-secondary/50 cursor-pointer"
                    >
                        <Github size={14} className="opacity-80" />
                        <span className="hidden sm:inline">GitHub</span>
                        <span className="text-yellow-400">★</span>
                        <span>Stars</span>
                        <span className="text-[10px] px-1.5 py-0.5 rounded bg-secondary/80 text-muted-foreground tabular-nums w-12 text-center">
                            {stars !== null ? stars.toLocaleString() : '-'}
                        </span>
                    </a>

                </div>
            </div>
        </nav>
    );
};

export default Navbar;
