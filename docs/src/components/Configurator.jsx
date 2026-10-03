import { useState } from 'react';
import { flushSync } from 'react-dom';
import { Layers, FileOutput, GitCompare, Settings, Github, Plus, Trash2, Check, Copy, Info, ChevronDown, ChevronUp, Terminal, Command, FileCode, FileCode2, Activity, BookOpen, RotateCcw } from 'lucide-react';

// The config reference of nanovision: keys, metrics and output formats.
// `go generate ./...` writes it from the Go code. Do not edit it.
// This page is built from it. A new key, metric or format shows up here
// without a change to this component.
import schema from '../generated/config.json';

// ---------------------------------------------------------------------------
// Schema helpers
// ---------------------------------------------------------------------------

const flatten = (fields) => fields.flatMap(f => [f, ...flatten(f.fields || [])]);
const ALL_FIELDS = flatten(schema.fields);
const fieldOf = (key) => ALL_FIELDS.find(f => f.key === key) || {};

// Keys with their own section on this page.
const CUSTOM_KEYS = ['reports', 'report_types', 'metrics'];

const isLeaf = (f) => f.type !== 'object' && f.type !== 'list';

// The leaves of a field that the generic form can edit. A list of objects
// (diff.path_maps) has no form; the reference table documents it.
const leavesOf = (f) => {
    if (f.type === 'list') return [];
    if (f.type === 'object') return f.fields.flatMap(leavesOf);
    return [f];
};

// The generic form: the plain top-level keys, then one group per object.
const GENERAL_FIELDS = schema.fields.filter(f => isLeaf(f) && !CUSTOM_KEYS.includes(f.key));
const GROUPS = schema.fields.filter(f => f.type === 'object' && !CUSTOM_KEYS.includes(f.key));
const FORM_LEAVES = [...GENERAL_FIELDS, ...GROUPS.flatMap(leavesOf)];

const SCOPES = [
    { id: 'files', title: 'Per file and folder', help: fieldOf('metrics.files').doc, flag: 'file-metric', thresholdPrefix: '' },
    { id: 'methods', title: 'Per method', help: fieldOf('metrics.methods').doc, flag: 'method-metric', thresholdPrefix: 'methods.' },
];

const DEFAULT_FORMATS = schema.outputFormats.filter(f => f.default).map(f => f.name);

// Optional extras for the output formats that have an example on this site.
const SITE = 'https://igorbayerl.github.io/nanovision/reports/';
const FORMAT_EXAMPLES = { Html: SITE, TextSummary: SITE + 'Summary.txt', Lcov: SITE + 'lcov.info', RawJson: SITE + 'RawJson.json' };

const initialMetrics = (scope) => schema.metrics[scope].map(m => ({
    name: m.name,
    shown: true, // without a metrics section the tool shows every metric
    ranged: false,
    min: m.kind === 'value' ? 10 : 60,
    max: m.kind === 'value' ? 15 : 80,
}));

// The state the page starts with, and the Reset button goes back to.
const INITIAL_REPORTS = [{ id: 1, path: 'coverage.out', source: '.', name: '' }];
const INITIAL_METRICS = { files: initialMetrics('files'), methods: initialMetrics('methods') };

const metricDoc = (scope, name) => schema.metrics[scope].find(m => m.name === name);

// ---------------------------------------------------------------------------
// Values
// ---------------------------------------------------------------------------

// The value of a leaf as the form holds it: boolean, or string. A list of
// strings is one string with an item per line.
const defaultValue = (f) => {
    if (f.type === 'boolean') return f.default === 'true';
    if (f.type === 'list of strings') return (f.default || '').split(', ').filter(Boolean).join('\n');
    return f.default || '';
};

const listItems = (value) => value.split('\n').map(s => s.trim()).filter(Boolean);

const rangeError = (scope, m) => {
    if (!m.shown || !m.ranged) return '';
    const min = Number(m.min), max = Number(m.max);
    if (m.min === '' || m.max === '' || Number.isNaN(min) || Number.isNaN(max)) return 'Both numbers are needed.';
    if (min < 0 || min > max) return 'Min must be 0 or more, and not above Max.';
    if (metricDoc(scope, m.name).kind === 'percentage' && max > 100) return 'A percentage cannot be above 100.';
    return '';
};

// ---------------------------------------------------------------------------
// UI components
// ---------------------------------------------------------------------------

// onReset adds a reset button at the end of the line; it is active when the
// section differs from its defaults.
const SectionHeader = ({ title, icon: Icon, onReset, changed }) => (
    <div className="flex items-center gap-2 text-foreground mb-5 pb-2 border-b border-border">
        <div className="text-primary"><Icon size={18} /></div>
        <h3 className="font-semibold text-sm uppercase tracking-wide">{title}</h3>
        {onReset && (
            <button
                type="button"
                onClick={onReset}
                disabled={!changed}
                aria-label={`Reset ${title} to defaults`}
                title={changed ? `Reset ${title} to defaults` : 'Nothing is changed'}
                className="ml-auto p-1 rounded text-muted-foreground enabled:hover:text-foreground enabled:hover:bg-secondary enabled:cursor-pointer disabled:opacity-30"
            >
                <RotateCcw size={15} />
            </button>
        )}
    </div>
);

const inputClass = "w-full px-3 py-2 bg-card border border-border rounded-md focus:ring-2 focus:ring-primary/20 focus:border-primary hover:border-primary/50 outline-none text-foreground placeholder-muted-foreground text-sm";

const Toggle = ({ checked, onChange, title }) => (
    <button
        type="button"
        role="switch"
        aria-checked={checked}
        title={title}
        onClick={() => onChange(!checked)}
        className={`w-11 h-6 rounded-full p-1 relative flex-shrink-0 transition-colors cursor-pointer ${checked ? 'bg-primary' : 'bg-secondary'}`}
    >
        <div className={`bg-white w-4 h-4 rounded-full shadow-sm transform transition-transform duration-200 ${checked ? 'translate-x-5' : ''}`}></div>
    </button>
);

const KeyLabel = ({ field }) => (
    <label className="block mb-1.5">
        <code className="text-sm font-medium text-foreground">{field.key}</code>
        {field.flag && <code className="ml-2 text-xs text-muted-foreground">-{field.flag}</code>}
    </label>
);

// One config key, edited by the control that fits its type.
const SettingField = ({ field, value, onChange }) => {
    const help = <p className="text-xs text-muted-foreground mt-1.5 leading-relaxed">{field.doc}</p>;

    if (field.type === 'boolean') {
        return (
            <div className="mb-5 flex items-start justify-between gap-4">
                <div>
                    <KeyLabel field={field} />
                    <p className="text-xs text-muted-foreground leading-relaxed">{field.doc}</p>
                </div>
                <Toggle checked={value} onChange={onChange} title={field.key} />
            </div>
        );
    }
    if (field.values) {
        return (
            <div className="mb-5">
                <KeyLabel field={field} />
                <div className="relative">
                    <select value={value} onChange={(e) => onChange(e.target.value)} className={`${inputClass} appearance-none cursor-pointer`}>
                        {!field.values.includes(value) && <option value={value}>{value || '(off)'}</option>}
                        {field.values.map(v => <option key={v} value={v}>{v}</option>)}
                    </select>
                    <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-muted-foreground"><ChevronDown size={14} /></div>
                </div>
                {help}
            </div>
        );
    }
    if (field.type === 'list of strings') {
        return (
            <div className="mb-5">
                <KeyLabel field={field} />
                <textarea
                    rows={Math.max(2, value.split('\n').length)}
                    value={value}
                    onChange={(e) => onChange(e.target.value)}
                    placeholder="one per line"
                    className={`${inputClass} font-mono`}
                />
                {help}
            </div>
        );
    }
    return (
        <div className="mb-5">
            <KeyLabel field={field} />
            <input
                type={field.type === 'number' ? 'number' : 'text'}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                placeholder={field.default || ''}
                className={inputClass}
            />
            {help}
        </div>
    );
};

const CopyButton = ({ text }) => {
    const [copied, setCopied] = useState(false);
    const handleCopy = () => {
        navigator.clipboard.writeText(text).then(() => {
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        }).catch(err => console.error('Copy failed', err));
    };
    return (
        <button onClick={handleCopy} className="flex items-center gap-1.5 px-2 py-1 rounded hover:bg-secondary text-xs text-muted-foreground hover:text-foreground cursor-pointer">
            {copied ? <Check size={14} className="text-green-500" /> : <Copy size={14} />}
            {copied ? "Copied" : "Copy"}
        </button>
    );
};

const CodePreview = ({ code, title }) => (
    <div className="rounded-xl border border-border bg-card overflow-hidden shadow-sm group">
        <div className="flex items-center justify-between px-4 py-2 border-b border-border bg-secondary/30">
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">{title}</span>
            <CopyButton text={code} />
        </div>
        <div className="p-4 overflow-x-auto overflow-y-auto max-h-[calc(100vh-16rem)] custom-scrollbar">
            <pre className="code-preview font-mono text-foreground whitespace-pre text-sm"><code>{code}</code></pre>
        </div>
    </div>
);

const MetricRow = ({ scope, metric, first, last, onChange, onMove }) => {
    const doc = metricDoc(scope, metric.name);
    const error = rangeError(scope, metric);
    // for a lower-is-better metric the range reads the other way
    const below = doc.lowerIsBetter ? 'Safe' : 'Danger';
    const above = doc.lowerIsBetter ? 'Danger' : 'Safe';
    const numberClass = `w-full px-2.5 py-1.5 bg-background border rounded text-sm focus:ring-1 focus:ring-primary outline-none ${error ? 'border-red-500/60' : 'border-border'}`;

    return (
        <div
            style={{ viewTransitionName: `metric-${scope}-${metric.name}` }}
            className={`rounded-lg border p-3 ${metric.shown ? 'bg-card border-border' : 'bg-secondary/20 border-border opacity-60'}`}
        >
            <div className="flex items-start gap-3">
                <div className="flex flex-col text-muted-foreground">
                    <button type="button" title="Move up" disabled={first} onClick={() => onMove(-1)} className="disabled:opacity-20 hover:text-foreground cursor-pointer"><ChevronUp size={14} /></button>
                    <button type="button" title="Move down" disabled={last} onClick={() => onMove(1)} className="disabled:opacity-20 hover:text-foreground cursor-pointer"><ChevronDown size={14} /></button>
                </div>
                <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium text-foreground">{doc.label} <code className="ml-1 text-xs font-normal text-muted-foreground">{metric.name}</code></div>
                    <p className="text-xs text-muted-foreground leading-relaxed mt-0.5">{doc.doc}</p>
                </div>
                <Toggle checked={metric.shown} onChange={(v) => onChange({ shown: v })} title={`Show ${doc.label}`} />
            </div>

            {metric.shown && (
                <div className="mt-3 pl-7">
                    <label className="flex items-center gap-2 text-xs text-muted-foreground cursor-pointer">
                        <input type="checkbox" checked={metric.ranged} onChange={(e) => onChange({ ranged: e.target.checked })} />
                        Set a warning range
                    </label>
                    <p className="text-xs text-muted-foreground mt-1">
                        {doc.lowerIsBetter ? 'Above' : 'Below'} the range is danger. Danger in changed code can fail the build (<code>review.fail_on</code>).
                    </p>
                    {metric.ranged && (
                        <div className="mt-2">
                            <div className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="block text-[10px] uppercase tracking-wider font-bold text-muted-foreground mb-1">Warning from</label>
                                    <input type="number" min="0" value={metric.min} onChange={(e) => onChange({ min: e.target.value })} className={numberClass} />
                                </div>
                                <div>
                                    <label className="block text-[10px] uppercase tracking-wider font-bold text-muted-foreground mb-1">Warning to</label>
                                    <input type="number" min="0" value={metric.max} onChange={(e) => onChange({ max: e.target.value })} className={numberClass} />
                                </div>
                            </div>
                            {error
                                ? <p className="text-xs text-red-400 mt-1.5">{error}</p>
                                : <p className="text-xs text-muted-foreground mt-1.5">
                                    Below {metric.min} = <span className={below === 'Danger' ? 'text-red-400' : 'text-green-500'}>{below}</span>,
                                    {' '}{metric.min} to {metric.max} = <span className="text-yellow-500">Warning</span>,
                                    {' '}above {metric.max} = <span className={above === 'Danger' ? 'text-red-400' : 'text-green-500'}>{above}</span>
                                </p>}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

// ---------------------------------------------------------------------------
// The configurator
// ---------------------------------------------------------------------------

const Configurator = () => {
    const [reports, setReports] = useState(INITIAL_REPORTS);
    const [formats, setFormats] = useState(DEFAULT_FORMATS);
    const [metrics, setMetrics] = useState(INITIAL_METRICS);
    // only the keys the user touched; every other key keeps its default
    const [values, setValues] = useState({});
    const [activeTab, setActiveTab] = useState('cli');

    const valueOf = (f) => (f.key in values ? values[f.key] : defaultValue(f));
    const setValue = (key, value) => setValues(prev => ({ ...prev, [key]: value }));
    // an emptied number field means "not set"
    const isChanged = (f) => valueOf(f) !== defaultValue(f) && !(f.type === 'number' && valueOf(f) === '');

    const updateReport = (id, patch) => setReports(prev => prev.map(r => r.id === id ? { ...r, ...patch } : r));
    const toggleFormat = (name) => setFormats(prev => prev.includes(name) ? prev.filter(n => n !== name) : [...prev, name]);
    const updateMetric = (scope, name, patch) => setMetrics(prev => ({ ...prev, [scope]: prev[scope].map(m => m.name === name ? { ...m, ...patch } : m) }));
    const moveMetric = (scope, index, by) => {
        const move = () => setMetrics(prev => {
            const list = [...prev[scope]];
            [list[index], list[index + by]] = [list[index + by], list[index]];
            return { ...prev, [scope]: list };
        });
        // the browser animates each row from its old place to its new place
        if (!document.startViewTransition) return move();
        // a hidden tab cancels the animation; the move itself still happens
        document.startViewTransition(() => flushSync(move)).ready.catch(() => {});
    };

    // Each section resets on its own. A key set back to its default by hand
    // counts as unchanged.
    const reportsChanged = JSON.stringify(reports) !== JSON.stringify(INITIAL_REPORTS);
    const formatsEdited = [...formats].sort().join() !== [...DEFAULT_FORMATS].sort().join();
    const metricsChanged = JSON.stringify(metrics) !== JSON.stringify(INITIAL_METRICS);
    const fieldsChanged = (fields) => fields.some(f => valueOf(f) !== defaultValue(f));
    const resetFields = (fields) => setValues(prev => {
        const next = { ...prev };
        for (const f of fields) delete next[f.key];
        return next;
    });

    // ---- what differs from the defaults ----

    const usedReports = reports.filter(r => r.path.trim() !== '');
    // output formats in the order of the reference, whatever the click order
    const chosenFormats = schema.outputFormats.map(f => f.name).filter(n => formats.includes(n));
    const formatsChanged = chosenFormats.join() !== DEFAULT_FORMATS.join();

    const shownMetrics = (scope) => metrics[scope].filter(m => m.shown);
    // the list is written when it is not "every metric, in the default order"
    const listChanged = (scope) => shownMetrics(scope).map(m => m.name).join() !== schema.metrics[scope].map(m => m.name).join();
    const rangedMetrics = (scope) => shownMetrics(scope).filter(m => m.ranged && !rangeError(scope, m));
    const scopeChanged = (scope) => listChanged(scope) || rangedMetrics(scope).length > 0;

    const changedLeaves = FORM_LEAVES.filter(isChanged);

    // ---- nanovision.yaml ----

    const quote = (s) => JSON.stringify(String(s));

    const yamlScalar = (f) => {
        const v = valueOf(f);
        if (f.type === 'boolean' || f.type === 'number') return String(v);
        return quote(v);
    };

    const yamlLeaf = (f, indent) => {
        const key = f.key.split('.').pop();
        if (f.type === 'list of strings') {
            return `${indent}${key}:\n` + listItems(valueOf(f)).map(item => `${indent}  - ${quote(item)}\n`).join('');
        }
        return `${indent}${key}: ${yamlScalar(f)}\n`;
    };

    // an object is written with only its changed keys
    const yamlObject = (f, indent) => {
        let body = '';
        for (const child of f.fields) {
            if (child.type === 'object') body += yamlObject(child, indent + '  ');
            else if (isLeaf(child) && isChanged(child)) body += yamlLeaf(child, indent + '  ');
        }
        return body ? `${indent}${f.key.split('.').pop()}:\n${body}` : '';
    };

    const yamlReports = () => 'reports:\n' + (usedReports.length ? usedReports : [{ path: '', source: '.', name: '' }]).map(r => {
        let item = `  - path: ${quote(r.path)}\n`;
        if (r.source.trim() && r.source.trim() !== '.') item += `    source: ${quote(r.source.trim())}\n`;
        if (r.name.trim()) item += `    name: ${quote(r.name.trim())}\n`;
        return item;
    }).join('');

    const yamlMetrics = () => {
        let body = '';
        for (const scope of SCOPES) {
            if (!scopeChanged(scope.id)) continue;
            body += `  ${scope.id}:\n`;
            for (const m of shownMetrics(scope.id)) {
                body += `    - name: ${m.name}\n`;
                if (m.ranged && !rangeError(scope.id, m)) body += `      warning: "${m.min}..${m.max}"\n`;
            }
        }
        return body ? `metrics:\n${body}` : '';
    };

    // The keys come out in the order of the reference.
    const generateYaml = () => {
        const blocks = [];
        for (const f of schema.fields) {
            let block = '';
            if (f.key === 'reports') block = yamlReports();
            else if (f.key === 'report_types') block = formatsChanged ? 'report_types:\n' + chosenFormats.map(n => `  - ${quote(n)}\n`).join('') : '';
            else if (f.key === 'metrics') block = yamlMetrics();
            else if (f.type === 'object') block = yamlObject(f, '');
            else if (isLeaf(f) && isChanged(f)) block = yamlLeaf(f, '');
            if (block) blocks.push(block);
        }
        return '# nanovision.yaml\n# Keys that are not here keep their default.\n\n' + blocks.join('\n');
    };

    // ---- command line ----

    // [flag, value] pairs, one for each line; a value of null is a flag without a value
    const flagPairs = () => {
        const pairs = [];
        for (const r of usedReports) {
            let value = r.path.trim();
            if (r.source.trim() && r.source.trim() !== '.') value += `,source=${r.source.trim()}`;
            if (r.name.trim()) value += `,name=${r.name.trim()}`;
            pairs.push(['report', value]);
        }
        if (formatsChanged) pairs.push([fieldOf('report_types').flag, chosenFormats.join(',')]);

        // a key with its own flag uses it; every other key uses -set key=value
        for (const f of changedLeaves) {
            const v = valueOf(f);
            if (f.type === 'list of strings') listItems(v).forEach(item => pairs.push(['set', `${f.key}=${item}`]));
            else if (!f.flag) pairs.push(['set', `${f.key}=${v}`]);
            else if (f.type === 'boolean') pairs.push(v ? [f.flag, null] : [`${f.flag}=false`, null]);
            else pairs.push([f.flag, v]);
        }

        for (const scope of SCOPES) {
            const ranged = (m) => m.ranged && !rangeError(scope.id, m);
            if (listChanged(scope.id)) {
                // one flag for each metric, in column order
                for (const m of shownMetrics(scope.id)) pairs.push([scope.flag, ranged(m) ? `${m.name}=${m.min}..${m.max}` : m.name]);
            } else {
                for (const m of rangedMetrics(scope.id)) pairs.push(['threshold', `${scope.thresholdPrefix}${m.name}=${m.min}..${m.max}`]);
            }
        }
        return pairs;
    };

    const bashQuote = (v) => '"' + String(v).replace(/(["\\$`])/g, '\\$1') + '"';
    const psQuote = (v) => '"' + String(v).replace(/(["`$])/g, '`$1') + '"';
    const generateBash = () => ['nanovision', ...flagPairs().map(([flag, v]) => v === null ? `-${flag}` : `-${flag}=${bashQuote(v)}`)].join(' \\\n  ');
    const generatePowerShell = () => ['nanovision', ...flagPairs().map(([flag, v]) => v === null ? `-${flag}` : `-${flag} ${psQuote(v)}`)].join(' `\n  ');

    const generateGithubActions = () => {
        const bashCmd = generateBash().replace(/^nanovision/, './nanovision').replace(/\\\n/g, '\\\n        ');
        const diffFile = valueOf(fieldOf('diff.file'));
        const diffStep = diffFile ? `
      - name: Fetch history for the diff
        run: git fetch origin main:main --depth=50

      - name: Make the diff
        run: git diff main > ${diffFile}
` : '';
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
          path: ${valueOf(fieldOf('output_dir'))}
`;
    };

    const TABS = [
        { id: 'cli', label: 'Terminal (Bash)', icon: Terminal, title: 'Command Line', code: generateBash },
        { id: 'ps', label: 'PowerShell', icon: Command, title: 'Command Line', code: generatePowerShell },
        { id: 'yaml', label: 'nanovision.yaml', icon: FileCode, title: 'Configuration File', code: generateYaml },
        { id: 'github', label: 'GitHub Actions', icon: Github, title: 'Workflow YAML', code: generateGithubActions },
    ];
    const tab = TABS.find(t => t.id === activeTab);

    return (
        <div className="min-h-screen flex flex-col bg-background text-foreground">
            <style>{'::view-transition-group(*) { animation-duration: 200ms; }'}</style>
            <main className="flex-1 py-7 px-4 sm:px-8">
                <div className="max-w-7xl mx-auto">
                    <div className="grid lg:grid-cols-12 gap-10">

                        {/* LEFT PANEL: CONFIG FORM */}
                        <div className="lg:col-span-5 space-y-10">

                            <section>
                                <SectionHeader title="Coverage reports" icon={Layers} changed={reportsChanged} onReset={() => setReports(INITIAL_REPORTS)} />
<p className="text-xs text-muted-foreground mb-4 leading-relaxed">{fieldOf('reports').doc}</p>
                                <div className="space-y-3">
                                    {reports.map(r => (
                                        <div key={r.id} className="bg-card p-4 rounded-lg border border-border relative">
                                            <div className="grid gap-3">
                                                {['path', 'source', 'name'].map(part => (
                                                    <div key={part}>
                                                        <label className="block mb-1"><code className="text-xs text-muted-foreground">reports[].{part}</code></label>
                                                        <input
                                                            type="text"
                                                            value={r[part]}
                                                            onChange={(e) => updateReport(r.id, { [part]: e.target.value })}
                                                            placeholder={{ path: 'e.g. coverage.out or build/**/*.gcov', source: '.', name: 'optional label' }[part]}
                                                            title={fieldOf(`reports[].${part}`).doc}
                                                            className={`${inputClass} ${part === 'path' && !r.path.trim() ? 'border-red-500/50' : ''}`}
                                                        />
                                                    </div>
                                                ))}
                                            </div>
                                            {reports.length > 1 && (
                                                <button onClick={() => setReports(prev => prev.filter(x => x.id !== r.id))} title="Remove report" className="absolute top-3 right-3 text-muted-foreground hover:text-red-500 p-1 cursor-pointer">
                                                    <Trash2 size={14} />
                                                </button>
                                            )}
                                        </div>
                                    ))}
                                </div>
                                <button onClick={() => setReports(prev => [...prev, { id: Date.now(), path: '', source: '.', name: '' }])} className="mt-3 text-xs bg-primary hover:bg-primary/90 text-primary-foreground px-2.5 py-1.5 rounded-md flex items-center gap-1.5 font-medium shadow-sm cursor-pointer">
                                    <Plus size={14} /> Add report
                                </button>
                            </section>

                            <section>
                                <SectionHeader title="Output formats" icon={FileOutput} changed={formatsEdited} onReset={() => setFormats(DEFAULT_FORMATS)} />
                                <div className="space-y-3">
                                    {schema.outputFormats.map(f => {
                                        const selected = formats.includes(f.name);
                                        return (
                                            <div
                                                key={f.name}
                                                onClick={() => toggleFormat(f.name)}
                                                className={`p-3 rounded-lg border cursor-pointer flex items-start gap-3 ${selected ? 'bg-primary/10 border-primary/50' : 'bg-card border-border hover:border-primary/50 hover:bg-secondary/50'}`}
                                            >
                                                <div className={`mt-0.5 w-4 h-4 rounded border flex items-center justify-center shrink-0 ${selected ? 'bg-primary border-primary' : 'border-muted-foreground bg-card'}`}>
                                                    {selected && <Check size={12} className="text-white" />}
                                                </div>
                                                <div className="flex-1">
                                                    <div className="flex justify-between items-center mb-1">
                                                        <span className={`text-sm font-medium ${selected ? 'text-primary' : 'text-foreground'}`}>{f.name}</span>
                                                        {FORMAT_EXAMPLES[f.name] && (
                                                            <a target="_blank" rel="noopener noreferrer" href={FORMAT_EXAMPLES[f.name]} onClick={(e) => e.stopPropagation()} className="text-[10px] text-primary hover:text-primary/80 hover:underline">Example</a>
                                                        )}
                                                    </div>
                                                    <p className="text-xs text-muted-foreground leading-snug">{f.doc} Writes <code>{f.writes}</code>.</p>
                                                </div>
                                            </div>
                                        );
                                    })}
                                </div>
                            </section>

                            <section>
                                <SectionHeader title="Metrics" icon={Activity} changed={metricsChanged} onReset={() => setMetrics(INITIAL_METRICS)} />
                                <p className="text-xs text-muted-foreground mb-5 leading-relaxed">{fieldOf('metrics').doc} The order is the column order.</p>
                                {SCOPES.map(scope => (
                                    <div key={scope.id} className="mb-6">
                                        <h4 className="font-semibold text-xs uppercase tracking-wide mb-1">{scope.title} <code className="ml-1 font-normal normal-case text-muted-foreground">metrics.{scope.id}</code></h4>
                                        <p className="text-xs text-muted-foreground mb-3">{scope.help}</p>
                                        <div className="space-y-2">
                                            {metrics[scope.id].map((m, i) => (
                                                <MetricRow
                                                    key={m.name}
                                                    scope={scope.id}
                                                    metric={m}
                                                    first={i === 0}
                                                    last={i === metrics[scope.id].length - 1}
                                                    onChange={(patch) => updateMetric(scope.id, m.name, patch)}
                                                    onMove={(by) => moveMetric(scope.id, i, by)}
                                                />
                                            ))}
                                        </div>
                                    </div>
                                ))}
                            </section>

                            <section>
                                <SectionHeader title="General" icon={Settings} changed={fieldsChanged(GENERAL_FIELDS)} onReset={() => resetFields(GENERAL_FIELDS)} />
                                {GENERAL_FIELDS.map(f => <SettingField key={f.key} field={f} value={valueOf(f)} onChange={(v) => setValue(f.key, v)} />)}
                            </section>

                            {GROUPS.map(group => (
                                <section key={group.key}>
                                    <SectionHeader title={group.key} icon={group.key === 'diff' ? GitCompare : Settings} changed={fieldsChanged(leavesOf(group))} onReset={() => resetFields(leavesOf(group))} />
                                    <p className="text-xs text-muted-foreground mb-4 leading-relaxed">
                                        {group.doc}
                                        {group.key === 'diff' && (
                                            <> <a href={`${import.meta.env.BASE_URL}about-coverage#diff-support`} target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">How to make a diff file</a></>
                                        )}
                                    </p>
                                    {leavesOf(group).map(f => <SettingField key={f.key} field={f} value={valueOf(f)} onChange={(v) => setValue(f.key, v)} />)}
                                </section>
                            ))}
                        </div>

                        {/* RIGHT PANEL: STICKY PREVIEW */}
                        <div className="lg:col-span-7 relative">
                            <div className="sticky top-24 space-y-4">
                                <div className="flex gap-1 p-1 bg-card border border-border rounded-lg mb-2 overflow-x-auto">
                                    {TABS.map(t => {
                                        const Icon = t.icon;
                                        return (
                                            <button
                                                key={t.id}
                                                onClick={() => setActiveTab(t.id)}
                                                className={`flex-1 flex items-center justify-center gap-2 px-3 py-2 rounded-md text-xs font-medium whitespace-nowrap cursor-pointer ${activeTab === t.id ? 'bg-secondary text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground hover:bg-secondary/50'}`}
                                            >
                                                <Icon size={14} />
                                                {t.label}
                                            </button>
                                        );
                                    })}
                                </div>

                                <CodePreview code={tab.code()} title={tab.title} />

                                {activeTab === 'github' && (
                                    <div className="bg-primary/5 border border-primary/20 p-4 rounded-lg flex gap-3 items-start">
                                        <div className="text-primary mt-0.5"><Info size={18} /></div>
                                        <div className="text-sm text-primary/80 leading-relaxed">
                                            This workflow installs Nanovision. Run your tests before this step, so the report file exists.
                                        </div>
                                    </div>
                                )}

                                {activeTab === 'yaml' && (
                                    <div className="text-sm text-muted-foreground flex gap-2 items-center px-2">
                                        <FileCode2 size={16} />
                                        Save it as <code>nanovision.yaml</code> in the project root. Check it with <code>nanovision config check</code>.
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* REFERENCE: every key, form or not */}
                    <section className="mt-16">
                        <SectionHeader title="Every key of nanovision.yaml" icon={BookOpen} />
                        <p className="text-xs text-muted-foreground mb-4 leading-relaxed">
                            Also printed by <code>nanovision config docs</code>. A <em>root</em> key is valid only in the root file.
                            A <code>nanovision.yaml</code> inside a folder can set the other keys for that folder. On the command line, <code>-set key=value</code> sets any key.
                        </p>
                        <div className="overflow-x-auto rounded-lg border border-border">
                            <table className="w-full text-sm">
                                <thead className="bg-secondary/30 text-left text-xs uppercase tracking-wide text-muted-foreground">
                                    <tr>
                                        <th className="px-3 py-2">Key</th>
                                        <th className="px-3 py-2">Type</th>
                                        <th className="px-3 py-2">Default</th>
                                        <th className="px-3 py-2">Flag</th>
                                        <th className="px-3 py-2">Description</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {schema.fields.flatMap(top => [top, ...flatten(top.fields || [])].map(f => (
                                        <tr key={f.key} className="border-t border-border align-top">
                                            <td className="px-3 py-2 whitespace-nowrap"><code>{f.key}</code>{top.rootOnly && <span className="ml-2 text-[10px] uppercase text-muted-foreground">root</span>}</td>
                                            <td className="px-3 py-2 text-muted-foreground">{f.values ? f.values.join(' | ') : f.type}</td>
                                            <td className="px-3 py-2 text-muted-foreground">{f.default || ''}</td>
                                            <td className="px-3 py-2 text-muted-foreground whitespace-nowrap">{f.flag ? <code>-{f.flag}</code> : ''}</td>
                                            <td className="px-3 py-2 text-muted-foreground">{f.doc}</td>
                                        </tr>
                                    )))}
                                </tbody>
                            </table>
                        </div>
                    </section>
                </div>
            </main>
        </div>
    );
};

export default Configurator;
