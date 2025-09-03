/**
 * NOTE: This is a developer data fixture for testing purposes.
 * This file provides a larger, more complex dataset to test rendering,
 * filtering, and performance of the coverage report UI.
 * Schema version: 1
 */
window.__ADLERCOV_SUMMARY__ = {
    schemaVersion: 1,
    generatedAt: '2025-09-01T16:49:31Z',
    title: 'Coverage Report',
    totals: {
        lineCoverage: { covered: 75, uncovered: 39, coverable: 114, total: 218, percentage: 65.78 },
        methodsCovered: { covered: 15, total: 22, percentage: 68.18 },
        methodsFullyCovered: { covered: 10, total: 22, percentage: 45.45 },
        files: 4,
        folders: 3,
        statuses: { lineCoverage: 'warning', methodsCovered: 'warning', methodsFullyCovered: 'danger' },
    },
    tree: [
        {
            id: 'test_project_go',
            name: 'test_project_go',
            type: 'folder',
            path: 'test_project_go',
            children: [
                {
                    id: 'test_project_go/calculator',
                    name: 'calculator',
                    type: 'folder',
                    path: 'test_project_go/calculator',
                    children: [
                        {
                            id: 'test_project_go/calculator/calculator.go',
                            name: 'calculator.go',
                            type: 'file',
                            path: 'test_project_go/calculator/calculator.go',
                            metrics: {
                                lineCoverage: {
                                    covered: 27,
                                    uncovered: 5,
                                    coverable: 32,
                                    total: 53,
                                    percentage: 84.37,
                                },
                                methodsCovered: { covered: 5, total: 5, percentage: 100 },
                                methodsFullyCovered: { covered: 2, total: 5, percentage: 40 },
                            },
                            statuses: { lineCoverage: 'safe', methodsCovered: 'safe', methodsFullyCovered: 'danger' },
                            targetUrl: 'test_project_go_calculator_calculator.go.html',
                        },
                        {
                            id: 'test_project_go/calculator/entities.go',
                            name: 'entities.go',
                            type: 'file',
                            path: 'test_project_go/calculator/entities.go',
                            metrics: {
                                lineCoverage: { covered: 20, uncovered: 5, coverable: 25, total: 56, percentage: 80 },
                                methodsCovered: { covered: 5, total: 6, percentage: 83.33 },
                                methodsFullyCovered: { covered: 5, total: 6, percentage: 83.33 },
                            },
                            statuses: { lineCoverage: 'safe', methodsCovered: 'safe', methodsFullyCovered: 'safe' },
                            targetUrl: 'test_project_go_calculator_entities.go.html',
                        },
                    ],
                    metrics: {
                        lineCoverage: { covered: 47, uncovered: 10, coverable: 57, total: 109, percentage: 82.45 },
                        methodsCovered: { covered: 10, total: 11, percentage: 90.9 },
                        methodsFullyCovered: { covered: 7, total: 11, percentage: 63.63 },
                    },
                    statuses: { lineCoverage: 'safe', methodsCovered: 'safe', methodsFullyCovered: 'warning' },
                },
                {
                    id: 'test_project_go/calculator_2',
                    name: 'calculator_2',
                    type: 'folder',
                    path: 'test_project_go/calculator_2',
                    children: [
                        {
                            id: 'test_project_go/calculator_2/calculator.go',
                            name: 'calculator.go',
                            type: 'file',
                            path: 'test_project_go/calculator_2/calculator.go',
                            metrics: {
                                lineCoverage: {
                                    covered: 23,
                                    uncovered: 9,
                                    coverable: 32,
                                    total: 53,
                                    percentage: 71.87,
                                },
                                methodsCovered: { covered: 4, total: 5, percentage: 80 },
                                methodsFullyCovered: { covered: 2, total: 5, percentage: 40 },
                            },
                            statuses: {
                                lineCoverage: 'warning',
                                methodsCovered: 'safe',
                                methodsFullyCovered: 'danger',
                            },
                            targetUrl: 'test_project_go_calculator_2_calculator.go.html',
                        },
                        {
                            id: 'test_project_go/calculator_2/entities.go',
                            name: 'entities.go',
                            type: 'file',
                            path: 'test_project_go/calculator_2/entities.go',
                            metrics: {
                                lineCoverage: { covered: 5, uncovered: 20, coverable: 25, total: 56, percentage: 20 },
                                methodsCovered: { covered: 1, total: 6, percentage: 16.66 },
                                methodsFullyCovered: { covered: 1, total: 6, percentage: 16.66 },
                            },
                            statuses: {
                                lineCoverage: 'danger',
                                methodsCovered: 'danger',
                                methodsFullyCovered: 'danger',
                            },
                            targetUrl: 'test_project_go_calculator_2_entities.go.html',
                        },
                    ],
                    metrics: {
                        lineCoverage: { covered: 28, uncovered: 29, coverable: 57, total: 109, percentage: 49.12 },
                        methodsCovered: { covered: 5, total: 11, percentage: 45.45 },
                        methodsFullyCovered: { covered: 3, total: 11, percentage: 27.27 },
                    },
                    statuses: { lineCoverage: 'danger', methodsCovered: 'danger', methodsFullyCovered: 'danger' },
                },
            ],
            metrics: {
                lineCoverage: { covered: 75, uncovered: 39, coverable: 114, total: 218, percentage: 65.78 },
                methodsCovered: { covered: 15, total: 22, percentage: 68.18 },
                methodsFullyCovered: { covered: 10, total: 22, percentage: 45.45 },
            },
            statuses: { lineCoverage: 'warning', methodsCovered: 'warning', methodsFullyCovered: 'danger' },
        },
    ],
    metricDefinitions: {
        branchCoverage: {
            label: 'Branches',
            shortLabel: 'Branches',
            subMetrics: [
                { id: 'covered', label: 'Covered', width: 100 },
                { id: 'total', label: 'Total', width: 80 },
                { id: 'percentage', label: 'Percentage %', width: 160 },
            ],
        },
        lineCoverage: {
            label: 'Lines',
            shortLabel: 'Lines',
            subMetrics: [
                { id: 'covered', label: 'Covered', width: 100 },
                { id: 'uncovered', label: 'Uncovered', width: 100 },
                { id: 'coverable', label: 'Coverable', width: 100 },
                { id: 'total', label: 'Total', width: 80 },
                { id: 'percentage', label: 'Percentage %', width: 160 },
            ],
        },
        methodsCovered: {
            label: 'Methods Covered',
            shortLabel: 'Methods Cov.',
            subMetrics: [
                { id: 'covered', label: 'Covered', width: 80 },
                { id: 'total', label: 'Total', width: 80 },
                { id: 'percentage', label: 'Percentage %', width: 160 },
            ],
        },
        methodsFullyCovered: {
            label: 'Methods Fully Covered',
            shortLabel: 'Methods Full Cov.',
            subMetrics: [
                { id: 'covered', label: 'Covered', width: 80 },
                { id: 'total', label: 'Total', width: 80 },
                { id: 'percentage', label: 'Percentage %', width: 160 },
            ],
        },
    },
    metadata: [
        { label: 'Generated At', value: '2025-09-01 16:49:31' },
        { label: 'Parser', value: 'GoCover' },
    ],
}
