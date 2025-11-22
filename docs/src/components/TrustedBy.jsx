import React from 'react';

const TrustedBy = () => (
    <section className="py-12 border-y border-border bg-secondary/30">
        <div className="max-w-7xl mx-auto px-6 text-center">
            <p className="text-xs font-mono text-muted-foreground uppercase tracking-widest mb-6">Trusted By Engineering Teams</p>
            <div className="flex justify-center items-center opacity-70 hover:opacity-100">
                <div className="flex items-center gap-3 text-2xl font-bold text-foreground tracking-wider">
                    <div className="w-10 h-10 bg-foreground text-background flex items-center justify-center font-black rounded transform -skew-x-6">C</div>
                    CRYTEK
                </div>
            </div>
        </div>
    </section>
);

export default TrustedBy;
