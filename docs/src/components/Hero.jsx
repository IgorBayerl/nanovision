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

// ===== Autoplay configuration (easy to tweak) =====
const AUTOPLAY_INTERVAL_MS = 2800; // time between patterns
const IDLE_TIMEOUT_MS = 8000;      // time without user interaction before autoplay resumes
const MAX_AUTO_COVERAGE_RATIO = 0.85; // never cover more than 85% automatically

// simple pattern generator that always leaves at least one word uncovered
function computeAutoStates(step, titleWordCount, patchWordCount, totalWords) {
    if (totalWords === 0) return {};

    const maxCoverable = Math.max(
        1,
        Math.min(totalWords - 1, Math.floor(totalWords * MAX_AUTO_COVERAGE_RATIO))
    );

    const indices = [];
    const pattern = step % 3;

    if (pattern === 0) {
        // cover first part of title
        const count = Math.min(maxCoverable, Math.ceil(titleWordCount * 0.6));
        for (let i = 0; i < count; i++) indices.push(i);
    } else if (pattern === 1) {
        // cover whole title + bit of patch
        const titleEnd = titleWordCount;
        for (let i = 0; i < titleEnd && indices.length < maxCoverable; i++) {
            indices.push(i);
        }
        let i = titleWordCount;
        while (i < titleWordCount + patchWordCount && indices.length < maxCoverable) {
            indices.push(i);
            i++;
        }
    } else {
        // cover some middle patch words
        const start = titleWordCount;
        const end = titleWordCount + patchWordCount;
        const span = Math.max(1, Math.min(end - start, maxCoverable));
        const offset = Math.max(0, Math.floor((end - start - span) / 2));
        for (let i = start + offset; i < start + offset + span; i++) {
            indices.push(i);
        }
    }

    // build state object with only "covered" entries
    const newStates = {};
    indices.forEach((idx) => {
        newStates[idx] = 'covered';
    });
    return newStates;
}

const Hero = () => {
    const [version, setVersion] = useState('Latest Release');

    const titleText = 'See your coverage clearly';
    const patchText =
        'It takes your coverage files (LCOV, Cobertura, GoCover, GCov), enriches them with static analysis, and outputs powerful reports in multiple formats..';

    const titleWords = useMemo(() => titleText.split(' ').filter((w) => w.length > 0), []);
    const patchWords = useMemo(() => patchText.split(' ').filter((w) => w.length > 0), []);

    const allWords = [...titleWords, ...patchWords];
    const titleWordCount = titleWords.length;
    const patchWordCount = patchWords.length;
    const totalWords = allWords.length;

    // state: index -> 'covered' | undefined (no 'uncovered' class anymore)
    const [wordStates, setWordStates] = useState({});
    const [isAuto, setIsAuto] = useState(true);
    const [autoStep, setAutoStep] = useState(0);

    const idleTimeoutRef = React.useRef(null);

    const handleHover = (globalIndex) => {
        // disable autoplay on user interaction
        setIsAuto(false);

        // restart autoplay after some idle time
        if (idleTimeoutRef.current) {
            window.clearTimeout(idleTimeoutRef.current);
        }
        idleTimeoutRef.current = window.setTimeout(() => {
            setIsAuto(true);
            setAutoStep((prev) => prev + 1); // go to next pattern
        }, IDLE_TIMEOUT_MS);

        // toggle covered/normal on hover
        setWordStates((prev) => {
            const currentState = prev[globalIndex];
            const nextState = currentState === 'covered' ? undefined : 'covered';
            return { ...prev, [globalIndex]: nextState };
        });
    };

    // derived metrics
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

    // fetch latest release
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
            if (idleTimeoutRef.current) window.clearTimeout(idleTimeoutRef.current);
        };
    }, []);

    // autoplay: increment step on interval
    useEffect(() => {
        if (!isAuto) return;

        const intervalId = window.setInterval(() => {
            setAutoStep((prev) => prev + 1);
        }, AUTOPLAY_INTERVAL_MS);

        return () => {
            window.clearInterval(intervalId);
        };
    }, [isAuto]);

    // derive wordStates from autoStep when autoplay is active
    useEffect(() => {
        if (!isAuto) return;

        const newStates = computeAutoStates(
            autoStep,
            titleWordCount,
            patchWordCount,
            totalWords
        );
        setWordStates(newStates);
    }, [autoStep, isAuto, titleWordCount, patchWordCount, totalWords]);

    // confetti only when user (not autoplay) reaches 100%
    useEffect(() => {
        if (!isAuto && lineCoverage === 100) {
            confetti({
                particleCount: 100,
                spread: 70,
                origin: { y: 0.6 },
            });
        }
    }, [lineCoverage, isAuto]);

    const renderInteractiveWords = (words, offset) => {
        return words.map((word, i) => {
            const globalIndex = offset + i;
            const state = wordStates[globalIndex];
            const isCovered = state === 'covered';

            return (
                <span
                    key={globalIndex}
                    onMouseEnter={() => handleHover(globalIndex)}
                    className={
                        'hover-word font-semibold ' +
                        (isCovered ? 'text-green-400' : 'text-foreground')
                    }
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
                            href="https://igorbayerl.github.io/nanovision/reports/"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="w-full sm:w-auto border border-border bg-background text-foreground px-8 py-4 rounded-lg font-bold flex items-center justify-center gap-2 hover:bg-secondary"
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
                                    subtext={
                                        <span>
                                            {coveredCount} / {totalWords} Words
                                        </span>
                                    }
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
