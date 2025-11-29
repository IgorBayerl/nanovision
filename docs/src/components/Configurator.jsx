import { useState } from 'react';
import { Layers, FileOutput, Tag, GitCompare, Settings, Github, Plus, Trash2, Check, Copy, Info, ChevronDown, Terminal, Command, FileCode, FileCode2, Activity } from 'lucide-react';

const REPORT_DESCRIPTIONS = {
    Html: {
        desc: "Interactive web report with source code navigation, search, and detailed metrics.",
        link: "https://igorbayerl.github.io/nanovision/reports/"
    },
    TextSummary: {
        desc: "Compact summary table printed to the console. Perfect for CI build logs.",
        link: "https://igorbayerl.github.io/nanovision/reports/Summary.txt"
    },
    Lcov: {
        desc: "Standard lcov.info format. Useful for uploading data to Coveralls or Codecov.",
        link: "https://igorbayerl.github.io/nanovision/reports/lcov.info"
    },
    RawJson: {
        desc: "Full internal data tree exported as JSON. Use this for custom post-processing tools.",
        link: "https://igorbayerl.github.io/nanovision/reports/RawJson.json"
    }
};

const METRICS = [
    { id: 'line_coverage', label: 'Line Coverage' },
    { id: 'branch_coverage', label: 'Branch Coverage' },
    { id: 'methods_covered', label: 'Methods Covered' },
    { id: 'methods_fully_covered', label: 'Methods Fully Covered' },
    { id: 'method_branch_coverage', label: 'Method Branch Coverage' },
    { id: 'patch_line_coverage', label: 'Patch Line Coverage' },
    { id: 'patch_methods_covered', label: 'Patch Methods Covered' },
];

// --- UI Components ---

const SectionHeader = ({ title, icon: Icon }) => (
    <div className="flex items-center gap-2 text-foreground mb-5 pb-2 border-b border-border">
        <div className="text-primary"><Icon size={18} /></div>
        <h3 className="font-semibold text-sm uppercase tracking-wide">{title}</h3>
    </div>
);

const InputField = ({ label, value, onChange, placeholder, help, error, type = "text", labelRight }) => (
    <div className="mb-5">
        <label className="block text-sm font-medium text-muted-foreground mb-1.5 flex items-center justify-between gap-2">
            <span>{label}</span>
            {labelRight && (
                <span className="text-xs">
                    {labelRight}
                </span>
            )}
        </label>
        <input
            type={type}
            className={`w-full px-3 py-2 bg-card border rounded-md focus:ring-2 focus:ring-offset-0 outline-none text-foreground placeholder-muted-foreground text-sm ${error
                ? 'border-red-500/50 focus:border-red-500 focus:ring-red-500/20'
                : 'border-border focus:border-primary focus:ring-primary/20 hover:border-primary/50'
                }`}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            placeholder={placeholder}
        />
        {help && !error && <p className="text-xs text-muted-foreground mt-1.5 leading-relaxed">{help}</p>}
        {error && <p className="text-xs text-red-400 mt-1.5">This field is required.</p>}
    </div>
);


const Toggle = ({ label, checked, onChange, help, className = "mb-5" }) => (
    <div className={`flex items-start justify-between group cursor-pointer ${className}`} onClick={() => onChange(!checked)}>
        <div className="mr-4">
            {label && <label className="block text-sm font-medium text-muted-foreground cursor-pointer">{label}</label>}
            {help && <p className="text-xs text-muted-foreground mt-1 leading-relaxed">{help}</p>}
        </div>
        <div
            className={`w-11 h-6 rounded-full p-1 relative flex-shrink-0 transition-colors ${checked ? 'bg-primary' : 'bg-secondary'}`}
        >
            <div className={`bg-white w-4 h-4 rounded-full shadow-sm transform transition-transform duration-200 ${checked ? 'translate-x-5' : ''}`}></div>
        </div>
    </div>
);

const Select = ({ label, value, onChange, options, help }) => (
    <div className="mb-5">
        <label className="block text-sm font-medium text-muted-foreground mb-1.5">{label}</label>
        <div className="relative">
            <select
                value={value}
                onChange={(e) => onChange(e.target.value)}
                className="w-full px-3 py-2 bg-card border border-border rounded-md focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none text-foreground text-sm appearance-none cursor-pointer hover:border-primary/50"
            >
                {options.map(opt => <option key={opt} value={opt}>{opt}</option>)}
            </select>
            <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-muted-foreground">
                <ChevronDown size={14} />
            </div>
        </div>
        {help && <p className="text-xs text-muted-foreground mt-1.5 leading-relaxed">{help}</p>}
    </div>
);

const ReportOption = ({ id, label, description, link, selected, onChange }) => (
    <div
        onClick={() => onChange(id)}
        className={`p-3 rounded-lg border cursor-pointer mb-3 flex items-start gap-3 ${selected
            ? 'bg-primary/10 border-primary/50'
            : 'bg-card border-border hover:border-primary/50 hover:bg-secondary/50'
            }`}
    >
        <div className={`mt-0.5 w-4 h-4 rounded border flex items-center justify-center shrink-0 ${selected ? 'bg-primary border-primary' : 'border-muted-foreground bg-card'
            }`}>
            {selected && <Check size={12} className="text-white" />}
        </div>
        <div className="flex-1">
            <div className="flex justify-between items-center mb-1">
                <span className={`text-sm font-medium ${selected ? 'text-primary' : 'text-foreground'}`}>{label}</span>
                <a target="_blank" href={link} onClick={(e) => e.stopPropagation()} className="text-[10px] text-primary hover:text-primary/80 hover:underline">Example</a>
            </div>
            <p className="text-xs text-muted-foreground leading-snug">{description}</p>
        </div>
    </div>
);

const CopyButton = ({ text }) => {
    const [copied, setCopied] = useState(false);

    const handleCopy = () => {
        navigator.clipboard.writeText(text).then(() => {
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        }).catch(err => {
            console.error('Copy failed', err);
        });
    };

    return (
        <button
            onClick={handleCopy}
            className="flex items-center gap-1.5 px-2 py-1 rounded hover:bg-secondary text-xs text-muted-foreground hover:text-foreground cursor-pointer"
        >
            {copied ? <Check size={14} className="text-green-500" /> : <Copy size={14} />}
            {copied ? "Copied" : "Copy"}
        </button>
    );
};

const CodePreview = ({ code, language, title }) => {
    return (
        <div className="rounded-xl border border-border bg-card overflow-hidden shadow-sm group">
            <div className="flex items-center justify-between px-4 py-2 border-b border-border bg-secondary/30">
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">{title || language}</span>
                <CopyButton text={code} />
            </div>
            <div className="p-4 overflow-x-auto overflow-y-auto max-h-[calc(100vh-16rem)] custom-scrollbar">
                <pre className="code-preview font-mono text-foreground whitespace-pre text-sm"><code>{code}</code></pre>
            </div>
        </div>
    );
};

// --- Main Configurator Component ---

const Configurator = () => {
    const [config, setConfig] = useState({
        inputs: [{ id: Date.now(), report: "coverage.out", source: "." }],
        outputDir: "coverage-report",
        reportTypes: { Html: true, TextSummary: true, Lcov: false, RawJson: false },
        title: "",
        tag: "",
        verbosity: "Info",
        logFormat: "text",
        filters: [{ id: 1, value: "-vendor/**" }, { id: 2, value: "-**/*_test.go" }],
        useDiff: false,
        diffFile: "",
        diffStrip: "auto",
        thresholds: {
            line_coverage: { min: 60, max: 80, active: false },
            branch_coverage: { min: 50, max: 70, active: false },
            methods_covered: { min: 50, max: 70, active: false },
            methods_fully_covered: { min: 50, max: 70, active: false },
            method_branch_coverage: { min: 50, max: 70, active: false },
            patch_line_coverage: { min: 50, max: 70, active: false },
            patch_methods_covered: { min: 50, max: 70, active: false },
        }
    });

    const [activeTab, setActiveTab] = useState("cli");

    // --- Handlers ---
    const handleChange = (key, value) => setConfig(prev => ({ ...prev, [key]: value }));

    const handleInputChange = (id, field, value) => {
        setConfig(prev => ({
            ...prev,
            inputs: prev.inputs.map(item => item.id === id ? { ...item, [field]: value } : item)
        }));
    };

    const handleThresholdChange = (metricId, field, value) => {
        setConfig(prev => ({
            ...prev,
            thresholds: {
                ...prev.thresholds,
                [metricId]: { ...prev.thresholds[metricId], [field]: value }
            }
        }));
    };

    const toggleThreshold = (metricId) => {
        setConfig(prev => ({
            ...prev,
            thresholds: {
                ...prev.thresholds,
                [metricId]: { ...prev.thresholds[metricId], active: !prev.thresholds[metricId].active }
            }
        }));
    };

    const addInputRow = () => setConfig(prev => ({ ...prev, inputs: [...prev.inputs, { id: Date.now(), report: "", source: "." }] }));
    const removeInputRow = (id) => { if (config.inputs.length > 1) setConfig(prev => ({ ...prev, inputs: prev.inputs.filter(item => item.id !== id) })); };

    const addFilter = () => setConfig(prev => ({ ...prev, filters: [...prev.filters, { id: Date.now(), value: "" }] }));
    const updateFilter = (id, val) => setConfig(prev => ({ ...prev, filters: prev.filters.map(f => f.id === id ? { ...f, value: val } : f) }));
    const removeFilter = (id) => setConfig(prev => ({ ...prev, filters: prev.filters.filter(f => f.id !== id) }));

    const toggleReportType = (type) => setConfig(prev => ({ ...prev, reportTypes: { ...prev.reportTypes, [type]: !prev.reportTypes[type] } }));

    // --- Generators ---
    const getReportsArray = () => config.inputs.map(i => i.report).filter(r => r.trim() !== "");
    const getSourcesArray = () => config.inputs.map(i => i.source).filter(s => s.trim() !== "");
    const getFiltersString = () => config.filters.map(f => f.value).filter(v => v.trim() !== "").join(";");
    const getSelectedTypes = () => Object.keys(config.reportTypes).filter(k => config.reportTypes[k]).join(",");

    const getThresholdsArray = () => {
        const parts = [];
        Object.entries(config.thresholds).forEach(([key, t]) => {
            if (t.active) {
                parts.push(`${key}=${t.min}..${t.max}`);
            }
        });
        return parts;
    };

    const generateBash = () => {
        const parts = ["nanovision"];
        const reports = getReportsArray().join(";");
        const sources = getSourcesArray().join(";");
        const filters = getFiltersString();
        const thresholds = getThresholdsArray();

        if (reports) parts.push(`-report="${reports}"`);
        if (sources) parts.push(`-sourcedirs="${sources}"`);
        if (config.outputDir) parts.push(`-output="${config.outputDir}"`);
        const types = getSelectedTypes();
        if (types && types !== "TextSummary,Html") parts.push(`-reporttypes="${types}"`);
        if (config.title) parts.push(`-title="${config.title}"`);
        if (config.tag) parts.push(`-tag="${config.tag}"`);
        if (config.verbosity !== "Info") parts.push(`-verbosity="${config.verbosity}"`);
        if (config.logFormat !== "text") parts.push(`-logformat="${config.logFormat}"`);
        if (filters) parts.push(`-filefilters="${filters}"`);
        if (config.useDiff) {
            if (config.diffFile) parts.push(`-diff="${config.diffFile}"`);
            if (config.diffStrip !== "auto") parts.push(`-diff-strip="${config.diffStrip}"`);
        }
        thresholds.forEach(t => parts.push(`-threshold="${t}"`));

        return parts.join(" \\\n  ");
    };

    const generatePowerShell = () => {
        const parts = ["nanovision"];
        const reports = getReportsArray().join(";");
        const sources = getSourcesArray().join(";");
        const filters = getFiltersString();
        const thresholds = getThresholdsArray();

        if (reports) parts.push(`-report "${reports}"`);
        if (sources) parts.push(`-sourcedirs "${sources}"`);
        if (config.outputDir) parts.push(`-output "${config.outputDir}"`);
        const types = getSelectedTypes();
        if (types && types !== "TextSummary,Html") parts.push(`-reporttypes "${types}"`);
        if (config.title) parts.push(`-title "${config.title}"`);
        if (config.tag) parts.push(`-tag "${config.tag}"`);
        if (config.verbosity !== "Info") parts.push(`-verbosity "${config.verbosity}"`);
        if (config.logFormat !== "text") parts.push(`-logformat "${config.logFormat}"`);
        if (filters) parts.push(`-filefilters "${filters}"`);
        if (config.useDiff) {
            if (config.diffFile) parts.push(`-diff "${config.diffFile}"`);
            if (config.diffStrip !== "auto") parts.push(`-diff-strip "${config.diffStrip}"`);
        }
        thresholds.forEach(t => parts.push(`-threshold "${t}"`));

        return parts.join(" `\n  ");
    };

    const generateYaml = () => {
        const typesList = Object.keys(config.reportTypes).filter(k => config.reportTypes[k]).map(t => `  - "${t}"`).join("\n");
        const reportsList = getReportsArray().map(r => `  - "${r}"`).join("\n");
        const sourcesList = getSourcesArray().map(s => `  - "${s}"`).join("\n");
        const filtersList = config.filters.map(f => f.value).filter(v => v.trim() !== "").map(f => `  - "${f}"`).join("\n");

        let yaml = `# nanovision.yaml Configuration
reports:
${reportsList || '  - ""'}

source_dirs:
${sourcesList || '  - "."'}

output_dir: "${config.outputDir}"

report_types:
${typesList || '  - "Html"'}

verbosity: "${config.verbosity}"
log_format: "${config.logFormat}"
`;
        if (config.title) yaml += `title: "${config.title}"\n`;
        if (config.tag) yaml += `tag: "${config.tag}"\n`;
        if (filtersList) yaml += `\nignore_files:\n${filtersList}\n`;
        if (config.useDiff) {
            yaml += `\ndiff:
  file: "${config.diffFile || 'my_changes.diff'}"
  strip: "${config.diffStrip}"
`;
        }

        const activeThresholds = Object.entries(config.thresholds).filter(([_, t]) => t.active);
        if (activeThresholds.length > 0) {
            yaml += `\n# Thresholds for status colors\nstatus_bands:\n`;
            activeThresholds.forEach(([key, t]) => {
                yaml += `  ${key}: "${t.min}..${t.max}"\n`;
            });
        }
        return yaml;
    };

    const generateGithubActions = () => {
        const bashCmd = generateBash().replace(/nanovision/g, './nanovision').replace(/\\\n/g, '\\\n        ');
        let diffStep = "";
        if (config.useDiff) {
            diffStep = `
      - name: Fetch history for Diff Coverage
        run: git fetch origin main:main --depth=50

      - name: Generate Diff
        run: git diff main > ${config.diffFile || 'changes.diff'}
`;
        }
        return `name: Coverage Report
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Nanovision
        run: curl -sL https://raw.githubusercontent.com/IgorBayerl/nanovision/main/scripts/install.sh | bash
${diffStep}
      - name: Generate Report
        run: |
          ${bashCmd}

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: ${config.outputDir}
`;
    };

    const getCode = () => {
        switch (activeTab) {
            case 'cli': return generateBash();
            case 'ps': return generatePowerShell();
            case 'yaml': return generateYaml();
            case 'github': return generateGithubActions();
            default: return "";
        }
    };

    return (
        <div className="min-h-screen flex flex-col bg-background text-foreground">
            {/* Main Content */}
            <main className="flex-1 py-7 px-4 sm:px-8">
                <div className="max-w-7xl mx-auto">

                    <div className="grid lg:grid-cols-12 gap-10">

                        {/* LEFT PANEL: CONFIG FORM */}
                        <div className="lg:col-span-5 space-y-10">

                            <section>
                                <SectionHeader title="Core Settings" icon={Layers} />

                                <div className="bg-primary/5 border border-primary/20 rounded-lg p-4 mb-6 text-sm text-primary leading-relaxed">
                                    <p className="mb-2"><strong className="font-bold">Report & Source Pairs:</strong> This is where you define your inputs.</p>
                                    <p>Add the coverage file generated by your tool (e.g., <code>coverage.out</code> from Go, <code>cobertura.xml</code>, or <code>gcov</code> files) and map it to the folder containing the actual source code.</p>
                                </div>

                                <div className="space-y-4 mb-6">
                                    <div className="flex items-center justify-between">
                                        <label className="block text-sm font-medium text-muted-foreground">Report & Source Pairs</label>
                                        <button onClick={addInputRow} className="text-xs bg-primary hover:bg-primary/90 text-primary-foreground px-2.5 py-1.5 rounded-md flex items-center gap-1.5 font-medium shadow-sm cursor-pointer">
                                            <Plus size={14} /> Add Pair
                                        </button>
                                    </div>

                                    <div className="space-y-3">
                                        {config.inputs.map((item, index) => (
                                            <div key={item.id} className="bg-card p-4 rounded-lg border border-border relative group hover:border-primary/50">
                                                <div className="grid gap-4">
                                                    <div>
                                                        <label className="text-[10px] text-muted-foreground uppercase tracking-wider font-bold mb-1.5 block">Report Pattern</label>
                                                        <input
                                                            type="text"
                                                            value={item.report}
                                                            onChange={(e) => handleInputChange(item.id, 'report', e.target.value)}
                                                            placeholder="e.g. coverage.out"
                                                            className={`w-full px-2.5 py-2 bg-background border rounded-md text-sm text-foreground focus:ring-2 outline-none ${!item.report.trim() ? 'border-red-500/50 focus:border-red-500 focus:ring-red-500/20' : 'border-border focus:border-primary focus:ring-primary/20'}`}
                                                        />
                                                    </div>
                                                    <div>
                                                        <label className="text-[10px] text-muted-foreground uppercase tracking-wider font-bold mb-1.5 block">Source Directory</label>
                                                        <input
                                                            type="text"
                                                            value={item.source}
                                                            onChange={(e) => handleInputChange(item.id, 'source', e.target.value)}
                                                            placeholder="e.g. ."
                                                            className={`w-full px-2.5 py-2 bg-background border rounded-md text-sm text-foreground focus:ring-2 outline-none ${!item.source.trim() ? 'border-red-500/50 focus:border-red-500 focus:ring-red-500/20' : 'border-border focus:border-primary focus:ring-primary/20'}`}
                                                        />
                                                    </div>
                                                </div>
                                                {config.inputs.length > 1 && (
                                                    <button onClick={() => removeInputRow(item.id)} className="absolute top-3 right-3 text-muted-foreground hover:text-red-500 p-1 cursor-pointer">
                                                        <Trash2 size={14} />
                                                    </button>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                </div>

                                <InputField
                                    label="Output Directory"
                                    value={config.outputDir}
                                    onChange={(v) => handleChange('outputDir', v)}
                                    placeholder="coverage-report"
                                    help="Directory where the report files will be created. Required."
                                    error={!config.outputDir.trim()}
                                />
                            </section>

                            <section>
                                <SectionHeader title="Output Formats" icon={FileOutput} />
                                <div className="space-y-2">
                                    {Object.keys(REPORT_DESCRIPTIONS).map(type => (
                                        <ReportOption
                                            key={type}
                                            id={type}
                                            label={type}
                                            description={REPORT_DESCRIPTIONS[type].desc}
                                            link={REPORT_DESCRIPTIONS[type].link}
                                            selected={config.reportTypes[type]}
                                            onChange={toggleReportType}
                                        />
                                    ))}
                                </div>
                            </section>

                            <section>
                                <SectionHeader title="Metadata" icon={Tag} />
                                <div className="space-y-4">
                                    <InputField
                                        label="Report Title"
                                        value={config.title}
                                        onChange={(v) => handleChange('title', v)}
                                        placeholder="My Project Coverage"
                                        help="Displayed in the header of the HTML report. Optional."
                                    />
                                    <InputField
                                        label="Tag / Build ID"
                                        value={config.tag}
                                        onChange={(v) => handleChange('tag', v)}
                                        placeholder="v1.0.0"
                                        help="Build version or tag displayed in the report. Optional."
                                    />
                                </div>
                            </section>

                            <section>
                                <SectionHeader title="Patch Coverage" icon={GitCompare} />
                                <div className="bg-card p-5 rounded-lg border border-border">
                                    <Toggle
                                        label="Enable Diff Analysis"
                                        checked={config.useDiff}
                                        onChange={(v) => handleChange('useDiff', v)}
                                        help="Limit coverage analysis to lines changed in a git diff."
                                    />

                                    {config.useDiff && (
                                        <div className="mt-5 space-y-5 pt-5 border-t border-border animate-in fade-in slide-in-from-top-2">
                                            <InputField
                                                label="Diff File Path"
                                                labelRight={
                                                    <a
                                                        href={`${import.meta.env.BASE_URL}about-coverage#diff-support`}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="text-primary hover:text-primary/80 hover:underline"
                                                        onClick={(e) => e.stopPropagation()}
                                                    >
                                                        How to generate diff file
                                                    </a>
                                                }
                                                value={config.diffFile}
                                                onChange={(v) => handleChange('diffFile', v)}
                                                placeholder="changes.diff"
                                                help="Path to a unified diff file."
                                            />
                                            <Select
                                                label="Strip Path Components"
                                                value={config.diffStrip}
                                                onChange={(v) => handleChange('diffStrip', v)}
                                                options={["auto", "0", "1", "2", "3"]}
                                                help="Removes leading folders from diff paths to match source files. 'Auto' works best."
                                            />
                                        </div>
                                    )}
                                </div>
                            </section>

                            <section>
                                <SectionHeader title="Advanced" icon={Settings} />
                                <div className="mb-6">
                                    <div className="flex items-center justify-between mb-2">
                                        <label className="block text-sm font-medium text-muted-foreground">File Filters</label>
                                        <button onClick={addFilter} className="text-xs text-primary hover:text-primary/80 flex items-center gap-1 cursor-pointer"><Plus size={12} /> Add</button>
                                    </div>

                                    <p className="text-xs text-muted-foreground mb-3 leading-relaxed">
                                        Include specific files or exclude libraries. Supports <b>glob patterns</b> (e.g., <code>**</code> for recursive folders).<br />
                                        Examples: <code>+src/**/*.go</code> (include), <code>-vendor/**</code> (exclude).
                                    </p>

                                    <div className="space-y-2">
                                        {config.filters.map((f) => (
                                            <div key={f.id} className="relative group">
                                                <input
                                                    type="text"
                                                    value={f.value}
                                                    onChange={(e) => updateFilter(f.id, e.target.value)}
                                                    className="w-full pl-3 pr-9 py-2 bg-card border border-border rounded-md text-sm text-foreground focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none"
                                                    placeholder="+include or -exclude"
                                                />
                                                <button
                                                    onClick={() => removeFilter(f.id)}
                                                    className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-red-500 p-1 rounded-md hover:bg-secondary cursor-pointer"
                                                    title="Remove filter"
                                                >
                                                    <Trash2 size={14} />
                                                </button>
                                            </div>
                                        ))}
                                    </div>
                                </div>

                                <div className="grid grid-cols-2 gap-4 mb-8">
                                    <Select label="Verbosity" value={config.verbosity} onChange={(v) => handleChange('verbosity', v)} options={["Verbose", "Info", "Warning", "Error", "Off"]} help="Log level details." />
                                    <Select label="Log Format" value={config.logFormat} onChange={(v) => handleChange('logFormat', v)} options={["text", "json"]} help="Console output format." />
                                </div>

                                {/* Thresholds Section moved inside Advanced */}
                                <div className="pt-6 border-t border-border">
                                    <div className="flex items-center gap-2 text-foreground mb-4">
                                        <div className="text-primary"><Activity size={16} /></div>
                                        <h4 className="font-semibold text-xs uppercase tracking-wide">Metric Thresholds</h4>
                                    </div>

                                    <div className="bg-card p-5 rounded-lg border border-border">
                                        <p className="text-xs text-muted-foreground mb-5 leading-relaxed">
                                            Define thresholds to color-code your report. <br />
                                            <span className="text-red-400 font-medium">Below Min = Danger</span> •
                                            <span className="text-yellow-500 font-medium mx-1">Min to Max = Warning</span> •
                                            <span className="text-green-500 font-medium">Above Max = Safe</span>
                                        </p>

                                        <div className="space-y-3">
                                            {METRICS.map(metric => {
                                                const t = config.thresholds[metric.id];
                                                return (
                                                    <div key={metric.id} className={`rounded-lg border transition-all duration-200 ${t.active ? 'bg-card border-primary/30 shadow-sm p-4' : 'bg-secondary/20 border-border opacity-60 p-3 hover:opacity-100'}`}>
                                                        <div className={`flex items-center justify-between ${t.active ? 'mb-4' : ''}`}>
                                                            <label className={`text-sm font-medium transition-colors ${t.active ? 'text-foreground' : 'text-muted-foreground'}`}>{metric.label}</label>
                                                            <Toggle
                                                                label=""
                                                                checked={t.active}
                                                                onChange={() => toggleThreshold(metric.id)}
                                                                className="mb-0"
                                                            />
                                                        </div>

                                                        {t.active && (
                                                            <div className="grid grid-cols-2 gap-4 animate-in fade-in slide-in-from-top-1">
                                                                <div>
                                                                    <label className="block text-[10px] uppercase tracking-wider font-bold text-muted-foreground mb-1.5">Min (Warning)</label>
                                                                    <input
                                                                        type="number"
                                                                        min="0"
                                                                        max="100"
                                                                        value={t.min}
                                                                        onChange={(e) => handleThresholdChange(metric.id, 'min', e.target.value)}
                                                                        className="w-full px-2.5 py-1.5 bg-background border border-border rounded text-sm focus:ring-1 focus:ring-primary outline-none"
                                                                    />
                                                                </div>
                                                                <div>
                                                                    <label className="block text-[10px] uppercase tracking-wider font-bold text-muted-foreground mb-1.5">Max (Success)</label>
                                                                    <input
                                                                        type="number"
                                                                        min="0"
                                                                        max="100"
                                                                        value={t.max}
                                                                        onChange={(e) => handleThresholdChange(metric.id, 'max', e.target.value)}
                                                                        className="w-full px-2.5 py-1.5 bg-background border border-border rounded text-sm focus:ring-1 focus:ring-primary outline-none"
                                                                    />
                                                                </div>
                                                            </div>
                                                        )}
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    </div>
                                </div>
                            </section>
                        </div>

                        {/* RIGHT PANEL: STICKY PREVIEW */}
                        <div className="lg:col-span-7 relative">
                            <div className="sticky top-24 space-y-4">
                                <div className="flex gap-1 p-1 bg-card border border-border rounded-lg mb-2 overflow-x-auto">
                                    {[
                                        { id: 'cli', label: 'Terminal (Bash)', icon: Terminal },
                                        { id: 'ps', label: 'PowerShell', icon: Command },
                                        { id: 'yaml', label: 'nanovision.yaml', icon: FileCode },
                                        { id: 'github', label: 'GitHub Actions', icon: Github },
                                    ].map(tab => {
                                        const Icon = tab.icon;
                                        return (
                                            <button
                                                key={tab.id}
                                                onClick={() => setActiveTab(tab.id)}
                                                className={`flex-1 flex items-center justify-center gap-2 px-3 py-2 rounded-md text-xs font-medium whitespace-nowrap cursor-pointer ${activeTab === tab.id
                                                    ? 'bg-secondary text-foreground shadow-sm'
                                                    : 'text-muted-foreground hover:text-foreground hover:bg-secondary/50'
                                                    }`}
                                            >
                                                <Icon size={14} />
                                                {tab.label}
                                            </button>
                                        );
                                    })}
                                </div>

                                <CodePreview
                                    code={getCode()}
                                    language={activeTab === 'yaml' ? 'yaml' : 'bash'}
                                    title={activeTab === 'yaml' ? 'Configuration File' : activeTab === 'github' ? 'Workflow YAML' : 'Command Line'}
                                />

                                {activeTab === 'github' && (
                                    <div className="bg-primary/5 border border-primary/20 p-4 rounded-lg flex gap-3 items-start">
                                        <div className="text-primary mt-0.5"><Info size={18} /></div>
                                        <div className="text-sm text-primary/80 leading-relaxed">
                                            This workflow automatically installs Nanovision. Ensure your tests run and generate the report file (e.g., <code>coverage.out</code>) <em>before</em> the Nanovision step executes.
                                        </div>
                                    </div>
                                )}

                                {activeTab === 'yaml' && (
                                    <div className="text-sm text-muted-foreground flex gap-2 items-center px-2">
                                        <FileCode2 size={16} />
                                        Save this as <code>nanovision.yaml</code> in your project root.
                                    </div>
                                )}
                            </div>
                        </div>

                    </div>
                </div>
            </main>
        </div>
    );
};

export default Configurator;
