import React, { useState, useEffect, useMemo } from 'react';
import { Play, FileText } from 'lucide-react';
import confetti from 'canvas-confetti';

const DynamicCoverageCard = ({ title, percentage, colorClass, subtext }) => {
    return (
        <div className="bg-card border border-border rounded-xl p-5 shadow-sm transition-shadow duration-500">
            <div className="flex justify-between items-end mb-3">
                <span className="text-sm font-medium text-muted-foreground">{title}</span>
                <span className="text-3xl font-bold text-foreground">{percentage}%</span>
            </div>
            <div className="w-full bg-secondary h-3 rounded-full overflow-hidden">
                <div
                    className={`h-full rounded-full ${colorClass} transition-all duration-1000 ease-out`}
                    style={{ width: `${percentage}%` }}
                ></div>
            </div>
            {subtext && (
                <div className="mt-2 text-[10px] text-muted-foreground font-mono flex justify-between">
                    {subtext}
                </div>
            )}
        </div>
    );
};

const Hero = () => {
    const [version, setVersion] = useState('Latest Release');

    const titleText = 'Fancy coverage reports ;)';
    const patchText =
        'It takes your coverage files (LCOV, Cobertura, GoCover, GCov), crunches the numbers, do some static analysis, and spits out many report formats. Simple as that.';

    const titleWords = useMemo(() => titleText.split(' ').filter((w) => w.length > 0), []);
    const patchWords = useMemo(() => patchText.split(' ').filter((w) => w.length > 0), []);

    const allWords = [...titleWords, ...patchWords];

    const titleWordCount = titleWords.length;
    const patchWordCount = patchWords.length;
    const totalWords = allWords.length;

    const [wordStates, setWordStates] = useState({});
    const resetTimerRef = React.useRef(null);

    const handleHover = (globalIndex) => {
        if (resetTimerRef.current) {
            clearTimeout(resetTimerRef.current);
        }

        setWordStates((prev) => {
            const currentState = prev[globalIndex];
            let newState;

            if (!currentState) {
                newState = 'covered';
            } else if (currentState === 'covered') {
                newState = 'uncovered';
            } else {
                newState = 'covered';
            }

            return { ...prev, [globalIndex]: newState };
        });

        resetTimerRef.current = setTimeout(() => {
            setWordStates({});
        }, 5000);
    };

    const coveredCount = Object.values(wordStates).filter((s) => s === 'covered').length;
    const lineCoverage =
        totalWords > 0 ? Math.min(100, Math.round((coveredCount / totalWords) * 100)) : 0;

    const isTitleFullyCovered = titleWords.every((_, i) => wordStates[i] === 'covered');
    const isDescFullyCovered = patchWords.every(
        (_, i) => wordStates[titleWordCount + i] === 'covered',
    );

    let methodsCovered = 0;
    if (isTitleFullyCovered) methodsCovered++;
    if (isDescFullyCovered) methodsCovered++;
    const methodCoverage = (methodsCovered / 2) * 100;

    const patchStart = titleWordCount;
    const patchEnd = titleWordCount + patchWordCount;
    let patchWordsCovered = 0;
    for (let i = patchStart; i < patchEnd; i++) {
        if (wordStates[i] === 'covered') {
            patchWordsCovered++;
        }
    }
    const patchCoverage =
        patchWordCount > 0 ? Math.min(100, Math.round((patchWordsCovered / patchWordCount) * 100)) : 0;

    useEffect(() => {
        fetch('https://api.github.com/repos/IgorBayerl/nanovision/releases/latest')
            .then((res) => res.json())
            .then((data) => {
                if (data.tag_name) setVersion(data.tag_name);
            })
            .catch(() => {
                setVersion('v1.0.0');
            });

        return () => {
            if (resetTimerRef.current) clearTimeout(resetTimerRef.current);
        };
    }, []);

    useEffect(() => {
        if (lineCoverage === 100) {
            confetti({
                particleCount: 100,
                spread: 70,
                origin: { y: 0.6 },
            });
        }
    }, [lineCoverage]);

    const renderInteractiveWords = (words, offset) => {
        return words.map((word, i) => {
            const globalIndex = offset + i;
            const state = wordStates[globalIndex];

            return (
                <span
                    key={globalIndex}
                    onMouseEnter={() => handleHover(globalIndex)}
                    className={`hover-word ${state || ''}`}
                >
                    {word}
                    {i < words.length - 1 ? ' ' : ''}
                </span>
            );
        });
    };

    return (
        <section className="relative pt-32 pb-24 overflow-hidden">
            <div className="max-w-7xl mx-auto px-6 grid lg:grid-cols-2 gap-16 items-center">
                <div className="text-center lg:text-left space-y-8">
                    <a
                        href="https://github.com/IgorBayerl/nanovision/releases"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-3 px-3 py-1 rounded-full border border-primary/20 bg-primary/5 text-primary text-xs font-mono mb-4 hover:bg-primary/10"
                    >
                        <span className="relative flex h-2 w-2">
                            <span
                                className={`animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${lineCoverage === 100 ? 'bg-green-400' : 'bg-primary'
                                    }`}
                            ></span>
                            <span
                                className={`relative inline-flex rounded-full h-2 w-2 ${lineCoverage === 100 ? 'bg-green-500' : 'bg-primary'
                                    }`}
                            ></span>
                        </span>
                        Latest Release: {version}
                    </a>

                    <h1 className="text-5xl lg:text-6xl font-bold leading-tight tracking-tighter text-foreground select-none">
                        {renderInteractiveWords(titleWords, 0)}
                    </h1>

                    <div className="text-lg text-foreground/90 leading-relaxed max-w-xl mx-auto lg:mx-0 font-normal select-none relative">
                        <div className="absolute left-[-18px] top-0 text-primary font-mono opacity-50 select-none pointer-events-none">
                            +
                        </div>
                        <p className="mb-2 flex flex-wrap gap-x-1">
                            {renderInteractiveWords(patchWords, titleWordCount)}
                        </p>
                    </div>

                    <div className="flex flex-col sm:flex-row items-center gap-4 justify-center lg:justify-start pt-4">
                        <a
                            href={`${import.meta.env.BASE_URL}getting-started`}
                            className="w-full sm:w-auto bg-primary text-primary-foreground px-8 py-4 rounded-lg font-bold flex items-center justify-center gap-2 hover:opacity-90 shadow-lg shadow-primary/20"
                        >
                            <Play size={20} />
                            Getting Started
                        </a>

                        <a
                            href="https://your-domain.com/example-report/index.html"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="w-full sm:w-auto border border-border bg-background text-foreground px-8 py-4 rounded-lg font-bold flex items-center justify-center gap-2 hover:bg-secondary transition-colors"
                        >
                            <FileText size={20} />
                            View Example Report
                        </a>
                    </div>

                </div>

                <div className="relative perspective-1000">
                    <div className="relative z-10 space-y-4 animate-float">
                        <div className="bg-card border border-border rounded-2xl p-6 shadow-2xl">
                            <div className="flex items-center justify-between border-b border-border pb-4 mb-6">
                                <div className="flex items-center gap-2">
                                    <div className="w-3 h-3 rounded-full bg-destructive"></div>
                                    <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                                    <div className="w-3 h-3 rounded-full bg-green-500"></div>
                                </div>
                                <div className="text-xs font-mono text-muted-foreground">
                                    coverage-report/index.html
                                </div>
                            </div>

                            <div className="grid gap-4">
                                <DynamicCoverageCard
                                    title="Line Coverage"
                                    percentage={lineCoverage}
                                    colorClass="bg-covered"
                                    subtext={<span>{coveredCount} / {totalWords} Words</span>}
                                />

                                <div className="grid grid-cols-2 gap-4">
                                    <DynamicCoverageCard
                                        title="Methods"
                                        percentage={methodCoverage}
                                        colorClass="bg-partial"
                                        subtext={<span>{methodsCovered} / 2</span>}
                                    />
                                    <DynamicCoverageCard
                                        title="Patch"
                                        percentage={patchCoverage}
                                        colorClass="bg-partial"
                                    />
                                </div>
                            </div>
                        </div>

                        <div className="absolute -right-8 -bottom-12 bg-popover border border-border p-4 rounded-lg shadow-xl hidden lg:block animate-pulse">
                            <pre className="text-xs font-mono text-muted-foreground whitespace-nowrap overflow-x-auto">
                                <span className="text-primary">nanovision</span> -report="coverage.out"
                            </pre>
                        </div>
                    </div>

                    <div className="absolute inset-0 bg-primary/5 blur-3xl -z-10 rounded-full transform scale-110"></div>
                </div>
            </div>
        </section>
    );
};

export default Hero;
