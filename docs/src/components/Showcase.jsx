import React, { useState, useEffect } from 'react';
import { LayoutDashboard, Microscope, GitPullRequest } from 'lucide-react';

const Showcase = () => {
    const [activeTab, setActiveTab] = useState(0);
    const [isHovered, setIsHovered] = useState(false);

    useEffect(() => {
        if (isHovered) return;
        const interval = setInterval(() => {
            setActiveTab((prev) => (prev + 1) % 3);
        }, 5000);
        return () => clearInterval(interval);
    }, [isHovered]);

    return (
        <section id="showcase" className="py-24 bg-background relative">
            <div className="max-w-7xl mx-auto px-6">
                <div className="grid lg:grid-cols-12 gap-12 items-center">
                    {/* Left Controls */}
                    <div
                        className="lg:col-span-4 space-y-2"
                        onMouseEnter={() => setIsHovered(true)}
                        onMouseLeave={() => setIsHovered(false)}
                    >
                        {/* Tab 1: Dashboard */}
                        <button
                            onClick={() => setActiveTab(0)}
                            className={`w-full text-left p-6 rounded-xl border relative overflow-hidden group cursor-pointer transition-all duration-300 ${activeTab === 0
                                ? 'bg-card border-primary/50 shadow-md'
                                : 'bg-transparent border-transparent hover:bg-secondary text-muted-foreground hover:text-foreground'
                                }`}
                        >
                            {activeTab === 0 && !isHovered && (
                                <div className="absolute bottom-0 left-0 h-1 bg-primary animate-fill-bar w-full origin-left"></div>
                            )}

                            <h3 className={`font-bold text-lg mb-2 flex items-center gap-3 ${activeTab === 0 ? 'text-foreground' : ''}`}>
                                <LayoutDashboard size={20} className={activeTab === 0 ? 'text-primary' : 'text-muted-foreground'} />
                                Fully Featured HTML Report
                            </h3>
                            <p className="text-sm leading-relaxed opacity-90 pl-8">Generates Static HTML reports that can be distributed and visualized easily.</p>
                        </button>

                        {/* Tab 2: Static Analysis */}
                        <button
                            onClick={() => setActiveTab(1)}
                            className={`w-full text-left p-6 rounded-xl border relative overflow-hidden group cursor-pointer transition-all duration-300 ${activeTab === 1
                                ? 'bg-card border-primary/50 shadow-md'
                                : 'bg-transparent border-transparent hover:bg-secondary text-muted-foreground hover:text-foreground'
                                }`}
                        >
                            {activeTab === 1 && !isHovered && (
                                <div className="absolute bottom-0 left-0 h-1 bg-primary animate-fill-bar w-full origin-left"></div>
                            )}

                            <h3 className={`font-bold text-lg mb-2 flex items-center gap-3 ${activeTab === 1 ? 'text-foreground' : ''}`}>
                                <Microscope size={20} className={activeTab === 1 ? 'text-primary' : 'text-muted-foreground'} />
                                Static Analysis
                            </h3>
                            <p className="text-sm leading-relaxed opacity-90 pl-8">The tool understands your code. It shows complexity alongside coverage.</p>
                        </button>

                        {/* Tab 3: Patch Coverage */}
                        <button
                            onClick={() => setActiveTab(2)}
                            className={`w-full text-left p-6 rounded-xl border relative overflow-hidden group cursor-pointer transition-all duration-300 ${activeTab === 2
                                ? 'bg-card border-primary/50 shadow-md'
                                : 'bg-transparent border-transparent hover:bg-secondary text-muted-foreground hover:text-foreground'
                                }`}
                        >
                            {activeTab === 2 && !isHovered && (
                                <div className="absolute bottom-0 left-0 h-1 bg-primary animate-fill-bar w-full origin-left"></div>
                            )}

                            <h3 className={`font-bold text-lg mb-2 flex items-center gap-3 ${activeTab === 2 ? 'text-foreground' : ''}`}>
                                <GitPullRequest size={20} className={activeTab === 2 ? 'text-primary' : 'text-muted-foreground'} />
                                Patch Coverage
                            </h3>
                            <p className="text-sm leading-relaxed opacity-90 pl-8">Focus on what matters. See coverage only for lines changed in your PR.</p>
                        </button>
                    </div>

                    {/* Right Preview (Image Area) */}
                    <div
                        className="lg:col-span-8 relative h-[500px] lg:h-[600px]"
                        onMouseEnter={() => setIsHovered(true)}
                        onMouseLeave={() => setIsHovered(false)}
                    >
                        <div className="absolute inset-0 bg-primary/5 blur-3xl -z-10"></div>
                        <div className="border border-border h-full w-full bg-card p-2 relative rounded-xl overflow-hidden shadow-2xl">

                            {/* Image 1 */}
                            <div
                                className={`absolute inset-2 flex items-center justify-center bg-background rounded-lg overflow-hidden transition-opacity duration-500 ${activeTab === 0 ? 'opacity-100 z-10' : 'opacity-0 z-0'
                                    }`}
                            >
                                <img
                                    src="/nanovision_srummary_screen_2.png"
                                    alt="Fully Featured HTML Report"
                                    className="w-full h-full object-cover"
                                />
                            </div>

                            {/* Image 2 */}
                            <div
                                className={`absolute inset-2 flex items-center justify-center bg-background rounded-lg overflow-hidden transition-opacity duration-500 ${activeTab === 1 ? 'opacity-100 z-10' : 'opacity-0 z-0'
                                    }`}
                            >
                                <img
                                    src="/nanovision_details.png"
                                    alt="Static Analysis"
                                    className="w-full h-full object-cover"
                                />
                            </div>

                            {/* Image 3 */}
                            <div
                                className={`absolute inset-2 flex items-center justify-center bg-background rounded-lg overflow-hidden transition-opacity duration-500 ${activeTab === 2 ? 'opacity-100 z-10' : 'opacity-0 z-0'
                                    }`}
                            >
                                <img
                                    src="/nanovision_patch_coverage.png"
                                    alt="Patch Coverage"
                                    className="w-full h-full object-cover"
                                />
                            </div>

                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default Showcase;
