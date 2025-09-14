// /**
//  * NOTE: This is a developer data fixture for testing purposes.
//  * This file provides a larger, more complex dataset to test rendering,
//  * filtering, and performance of the coverage report UI.
//  * Schema version: 1
//  */
// window.__ADLERCOV_SUMMARY__ = {
//     schemaVersion: 1,
//     generatedAt: '2025-09-01T16:49:31Z',
//     title: 'Coverage Report',
//     totals: {
//         lineCoverage: { covered: 75, uncovered: 39, coverable: 114, total: 218, percentage: 65.78 },
//         methodsCovered: { covered: 15, total: 22, percentage: 68.18 },
//         methodsFullyCovered: { covered: 10, total: 22, percentage: 45.45 },
//         files: 4,
//         folders: 3,
//         statuses: { lineCoverage: 'warning', methodsCovered: 'warning', methodsFullyCovered: 'danger' },
//     },
//     tree: [
//         {
//             id: 'test_project_go',
//             name: 'test_project_go',
//             type: 'folder',
//             path: 'test_project_go',
//             children: [
//                 {
//                     id: 'test_project_go/calculator',
//                     name: 'calculator',
//                     type: 'folder',
//                     path: 'test_project_go/calculator',
//                     children: [
//                         {
//                             id: 'test_project_go/calculator/calculator.go',
//                             name: 'calculator.go',
//                             type: 'file',
//                             path: 'test_project_go/calculator/calculator.go',
//                             metrics: {
//                                 lineCoverage: {
//                                     covered: 27,
//                                     uncovered: 5,
//                                     coverable: 32,
//                                     total: 53,
//                                     percentage: 84.37,
//                                 },
//                                 methodsCovered: { covered: 5, total: 5, percentage: 100 },
//                                 methodsFullyCovered: { covered: 2, total: 5, percentage: 40 },
//                             },
//                             statuses: { lineCoverage: 'safe', methodsCovered: 'safe', methodsFullyCovered: 'danger' },
//                             targetUrl: 'test_project_go_calculator_calculator.go.html',
//                         },
//                         {
//                             id: 'test_project_go/calculator/entities.go',
//                             name: 'entities.go',
//                             type: 'file',
//                             path: 'test_project_go/calculator/entities.go',
//                             metrics: {
//                                 lineCoverage: { covered: 20, uncovered: 5, coverable: 25, total: 56, percentage: 80 },
//                                 methodsCovered: { covered: 5, total: 6, percentage: 83.33 },
//                                 methodsFullyCovered: { covered: 5, total: 6, percentage: 83.33 },
//                             },
//                             statuses: { lineCoverage: 'safe', methodsCovered: 'safe', methodsFullyCovered: 'safe' },
//                             targetUrl: 'test_project_go_calculator_entities.go.html',
//                         },
//                     ],
//                     metrics: {
//                         lineCoverage: { covered: 47, uncovered: 10, coverable: 57, total: 109, percentage: 82.45 },
//                         methodsCovered: { covered: 10, total: 11, percentage: 90.9 },
//                         methodsFullyCovered: { covered: 7, total: 11, percentage: 63.63 },
//                     },
//                     statuses: { lineCoverage: 'safe', methodsCovered: 'safe', methodsFullyCovered: 'warning' },
//                 },
//                 {
//                     id: 'test_project_go/calculator_2',
//                     name: 'calculator_2',
//                     type: 'folder',
//                     path: 'test_project_go/calculator_2',
//                     children: [
//                         {
//                             id: 'test_project_go/calculator_2/calculator.go',
//                             name: 'calculator.go',
//                             type: 'file',
//                             path: 'test_project_go/calculator_2/calculator.go',
//                             metrics: {
//                                 lineCoverage: {
//                                     covered: 23,
//                                     uncovered: 9,
//                                     coverable: 32,
//                                     total: 53,
//                                     percentage: 71.87,
//                                 },
//                                 methodsCovered: { covered: 4, total: 5, percentage: 80 },
//                                 methodsFullyCovered: { covered: 2, total: 5, percentage: 40 },
//                             },
//                             statuses: {
//                                 lineCoverage: 'warning',
//                                 methodsCovered: 'safe',
//                                 methodsFullyCovered: 'danger',
//                             },
//                             targetUrl: 'test_project_go_calculator_2_calculator.go.html',
//                         },
//                         {
//                             id: 'test_project_go/calculator_2/entities.go',
//                             name: 'entities.go',
//                             type: 'file',
//                             path: 'test_project_go/calculator_2/entities.go',
//                             metrics: {
//                                 lineCoverage: { covered: 5, uncovered: 20, coverable: 25, total: 56, percentage: 20 },
//                                 methodsCovered: { covered: 1, total: 6, percentage: 16.66 },
//                                 methodsFullyCovered: { covered: 1, total: 6, percentage: 16.66 },
//                             },
//                             statuses: {
//                                 lineCoverage: 'danger',
//                                 methodsCovered: 'danger',
//                                 methodsFullyCovered: 'danger',
//                             },
//                             targetUrl: 'test_project_go_calculator_2_entities.go.html',
//                         },
//                     ],
//                     metrics: {
//                         lineCoverage: { covered: 28, uncovered: 29, coverable: 57, total: 109, percentage: 49.12 },
//                         methodsCovered: { covered: 5, total: 11, percentage: 45.45 },
//                         methodsFullyCovered: { covered: 3, total: 11, percentage: 27.27 },
//                     },
//                     statuses: { lineCoverage: 'danger', methodsCovered: 'danger', methodsFullyCovered: 'danger' },
//                 },
//             ],
//             metrics: {
//                 lineCoverage: { covered: 75, uncovered: 39, coverable: 114, total: 218, percentage: 65.78 },
//                 methodsCovered: { covered: 15, total: 22, percentage: 68.18 },
//                 methodsFullyCovered: { covered: 10, total: 22, percentage: 45.45 },
//             },
//             statuses: { lineCoverage: 'warning', methodsCovered: 'warning', methodsFullyCovered: 'danger' },
//         },
//     ],
//     metricDefinitions: {
//         branchCoverage: {
//             label: 'Branches',
//             shortLabel: 'Branches',
//             subMetrics: [
//                 { id: 'covered', label: 'Covered', width: 100 },
//                 { id: 'total', label: 'Total', width: 80 },
//                 { id: 'percentage', label: 'Percentage %', width: 160 },
//             ],
//         },
//         lineCoverage: {
//             label: 'Lines',
//             shortLabel: 'Lines',
//             subMetrics: [
//                 { id: 'covered', label: 'Covered', width: 100 },
//                 { id: 'uncovered', label: 'Uncovered', width: 100 },
//                 { id: 'coverable', label: 'Coverable', width: 100 },
//                 { id: 'total', label: 'Total', width: 80 },
//                 { id: 'percentage', label: 'Percentage %', width: 160 },
//             ],
//         },
//         methodsCovered: {
//             label: 'Methods Covered',
//             shortLabel: 'Methods Cov.',
//             subMetrics: [
//                 { id: 'covered', label: 'Covered', width: 80 },
//                 { id: 'total', label: 'Total', width: 80 },
//                 { id: 'percentage', label: 'Percentage %', width: 160 },
//             ],
//         },
//         methodsFullyCovered: {
//             label: 'Methods Fully Covered',
//             shortLabel: 'Methods Full Cov.',
//             subMetrics: [
//                 { id: 'covered', label: 'Covered', width: 80 },
//                 { id: 'total', label: 'Total', width: 80 },
//                 { id: 'percentage', label: 'Percentage %', width: 160 },
//             ],
//         },
//     },
//     metadata: [
//         { label: 'Generated At', value: '2025-09-01 16:49:31' },
//         { label: 'Parser', value: 'GoCover' },
//     ],
// }

window.__ADLERCOV_SUMMARY__ = {
    schemaVersion: 1,
    generatedAt: '2025-09-13T10:54:29Z',
    title: 'Coverage Report',
    totals: {
        lineCoverage: { covered: 2141, uncovered: 1875, coverable: 4016, total: 7062, percentage: 53.31 },
        methodsCovered: { covered: 162, total: 354, percentage: 45.76 },
        methodsFullyCovered: { covered: 94, total: 354, percentage: 26.55 },
        files: 52,
        folders: 37,
        statuses: { lineCoverage: 'danger', methodsCovered: 'danger', methodsFullyCovered: 'danger' },
    },
    tree: [
        {
            id: 'C:',
            name: 'C:',
            type: 'folder',
            path: 'C:',
            children: [
                {
                    id: 'C:/www',
                    name: 'www',
                    type: 'folder',
                    path: 'C:/www',
                    children: [
                        {
                            id: 'C:/www/AdlerCov',
                            name: 'AdlerCov',
                            type: 'folder',
                            path: 'C:/www/AdlerCov',
                            children: [
                                {
                                    id: 'C:/www/AdlerCov/cmd',
                                    name: 'cmd',
                                    type: 'folder',
                                    path: 'C:/www/AdlerCov/cmd',
                                    children: [
                                        {
                                            id: 'C:/www/AdlerCov/cmd/main.go',
                                            name: 'main.go',
                                            type: 'file',
                                            path: 'C:/www/AdlerCov/cmd/main.go',
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 130,
                                                    uncovered: 51,
                                                    coverable: 181,
                                                    total: 264,
                                                    percentage: 71.82,
                                                },
                                                methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                                methodsFullyCovered: { covered: 2, total: 6, percentage: 33.33 },
                                            },
                                            statuses: {
                                                lineCoverage: 'warning',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'danger',
                                            },
                                            targetUrl: 'C:_www_AdlerCov_cmd_main.go.html',
                                        },
                                    ],
                                    metrics: {
                                        lineCoverage: {
                                            covered: 130,
                                            uncovered: 51,
                                            coverable: 181,
                                            total: 264,
                                            percentage: 71.82,
                                        },
                                        methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                        methodsFullyCovered: { covered: 2, total: 6, percentage: 33.33 },
                                    },
                                    statuses: {
                                        lineCoverage: 'warning',
                                        methodsCovered: 'safe',
                                        methodsFullyCovered: 'danger',
                                    },
                                },
                            ],
                            metrics: {
                                lineCoverage: {
                                    covered: 130,
                                    uncovered: 51,
                                    coverable: 181,
                                    total: 264,
                                    percentage: 71.82,
                                },
                                methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                methodsFullyCovered: { covered: 2, total: 6, percentage: 33.33 },
                            },
                            statuses: {
                                lineCoverage: 'warning',
                                methodsCovered: 'safe',
                                methodsFullyCovered: 'danger',
                            },
                        },
                    ],
                    metrics: {
                        lineCoverage: { covered: 130, uncovered: 51, coverable: 181, total: 264, percentage: 71.82 },
                        methodsCovered: { covered: 6, total: 6, percentage: 100 },
                        methodsFullyCovered: { covered: 2, total: 6, percentage: 33.33 },
                    },
                    statuses: { lineCoverage: 'warning', methodsCovered: 'safe', methodsFullyCovered: 'danger' },
                },
            ],
            metrics: {
                lineCoverage: { covered: 130, uncovered: 51, coverable: 181, total: 264, percentage: 71.82 },
                methodsCovered: { covered: 6, total: 6, percentage: 100 },
                methodsFullyCovered: { covered: 2, total: 6, percentage: 33.33 },
            },
            statuses: { lineCoverage: 'warning', methodsCovered: 'safe', methodsFullyCovered: 'danger' },
        },
        {
            id: 'github.com',
            name: 'github.com',
            type: 'folder',
            path: 'github.com',
            children: [
                {
                    id: 'github.com/IgorBayerl',
                    name: 'IgorBayerl',
                    type: 'folder',
                    path: 'github.com/IgorBayerl',
                    children: [
                        {
                            id: 'github.com/IgorBayerl/AdlerCov',
                            name: 'AdlerCov',
                            type: 'folder',
                            path: 'github.com/IgorBayerl/AdlerCov',
                            children: [
                                {
                                    id: 'github.com/IgorBayerl/AdlerCov/analyzer',
                                    name: 'analyzer',
                                    type: 'folder',
                                    path: 'github.com/IgorBayerl/AdlerCov/analyzer',
                                    children: [
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/analyzer/cpp',
                                            name: 'cpp',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/analyzer/cpp',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/analyzer/cpp/analyzer.go',
                                                    name: 'analyzer.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/analyzer/cpp/analyzer.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 97,
                                                            uncovered: 13,
                                                            coverable: 110,
                                                            total: 194,
                                                            percentage: 88.18,
                                                        },
                                                        methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                                        methodsFullyCovered: { covered: 3, total: 6, percentage: 50 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_analyzer_cpp_analyzer.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 97,
                                                    uncovered: 13,
                                                    coverable: 110,
                                                    total: 194,
                                                    percentage: 88.18,
                                                },
                                                methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                                methodsFullyCovered: { covered: 3, total: 6, percentage: 50 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'danger',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/analyzer/go',
                                            name: 'go',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/analyzer/go',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/analyzer/go/analyzer.go',
                                                    name: 'analyzer.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/analyzer/go/analyzer.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 82,
                                                            uncovered: 13,
                                                            coverable: 95,
                                                            total: 164,
                                                            percentage: 86.31,
                                                        },
                                                        methodsCovered: { covered: 5, total: 5, percentage: 100 },
                                                        methodsFullyCovered: { covered: 3, total: 5, percentage: 60 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'warning',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_analyzer_go_analyzer.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 82,
                                                    uncovered: 13,
                                                    coverable: 95,
                                                    total: 164,
                                                    percentage: 86.31,
                                                },
                                                methodsCovered: { covered: 5, total: 5, percentage: 100 },
                                                methodsFullyCovered: { covered: 3, total: 5, percentage: 60 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'warning',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/analyzer/analyzer.go',
                                            name: 'analyzer.go',
                                            type: 'file',
                                            path: 'github.com/IgorBayerl/AdlerCov/analyzer/analyzer.go',
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 0,
                                                    uncovered: 3,
                                                    coverable: 3,
                                                    total: 36,
                                                    percentage: 0,
                                                },
                                                methodsCovered: { covered: 0, total: 1, percentage: 0 },
                                                methodsFullyCovered: { covered: 0, total: 1, percentage: 0 },
                                            },
                                            statuses: {
                                                lineCoverage: 'danger',
                                                methodsCovered: 'danger',
                                                methodsFullyCovered: 'danger',
                                            },
                                            targetUrl: 'github.com_IgorBayerl_AdlerCov_analyzer_analyzer.go.html',
                                        },
                                    ],
                                    metrics: {
                                        lineCoverage: {
                                            covered: 179,
                                            uncovered: 29,
                                            coverable: 208,
                                            total: 394,
                                            percentage: 86.05,
                                        },
                                        methodsCovered: { covered: 11, total: 12, percentage: 91.66 },
                                        methodsFullyCovered: { covered: 6, total: 12, percentage: 50 },
                                    },
                                    statuses: {
                                        lineCoverage: 'safe',
                                        methodsCovered: 'safe',
                                        methodsFullyCovered: 'danger',
                                    },
                                },
                                {
                                    id: 'github.com/IgorBayerl/AdlerCov/filereader',
                                    name: 'filereader',
                                    type: 'folder',
                                    path: 'github.com/IgorBayerl/AdlerCov/filereader',
                                    children: [
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/filereader/default_reader.go',
                                            name: 'default_reader.go',
                                            type: 'file',
                                            path: 'github.com/IgorBayerl/AdlerCov/filereader/default_reader.go',
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 15,
                                                    uncovered: 0,
                                                    coverable: 15,
                                                    total: 28,
                                                    percentage: 100,
                                                },
                                                methodsCovered: { covered: 5, total: 5, percentage: 100 },
                                                methodsFullyCovered: { covered: 5, total: 5, percentage: 100 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'safe',
                                            },
                                            targetUrl:
                                                'github.com_IgorBayerl_AdlerCov_filereader_default_reader.go.html',
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/filereader/filereader.go',
                                            name: 'filereader.go',
                                            type: 'file',
                                            path: 'github.com/IgorBayerl/AdlerCov/filereader/filereader.go',
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 47,
                                                    uncovered: 16,
                                                    coverable: 63,
                                                    total: 95,
                                                    percentage: 74.6,
                                                },
                                                methodsCovered: { covered: 3, total: 3, percentage: 100 },
                                                methodsFullyCovered: { covered: 0, total: 3, percentage: 0 },
                                            },
                                            statuses: {
                                                lineCoverage: 'warning',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'danger',
                                            },
                                            targetUrl: 'github.com_IgorBayerl_AdlerCov_filereader_filereader.go.html',
                                        },
                                    ],
                                    metrics: {
                                        lineCoverage: {
                                            covered: 62,
                                            uncovered: 16,
                                            coverable: 78,
                                            total: 123,
                                            percentage: 79.48,
                                        },
                                        methodsCovered: { covered: 8, total: 8, percentage: 100 },
                                        methodsFullyCovered: { covered: 5, total: 8, percentage: 62.5 },
                                    },
                                    statuses: {
                                        lineCoverage: 'warning',
                                        methodsCovered: 'safe',
                                        methodsFullyCovered: 'warning',
                                    },
                                },
                                {
                                    id: 'github.com/IgorBayerl/AdlerCov/filesystem',
                                    name: 'filesystem',
                                    type: 'folder',
                                    path: 'github.com/IgorBayerl/AdlerCov/filesystem',
                                    children: [
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/filesystem/filesystem.go',
                                            name: 'filesystem.go',
                                            type: 'file',
                                            path: 'github.com/IgorBayerl/AdlerCov/filesystem/filesystem.go',
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 0,
                                                    uncovered: 11,
                                                    coverable: 11,
                                                    total: 102,
                                                    percentage: 0,
                                                },
                                                methodsCovered: { covered: 0, total: 9, percentage: 0 },
                                                methodsFullyCovered: { covered: 0, total: 9, percentage: 0 },
                                            },
                                            statuses: {
                                                lineCoverage: 'danger',
                                                methodsCovered: 'danger',
                                                methodsFullyCovered: 'danger',
                                            },
                                            targetUrl: 'github.com_IgorBayerl_AdlerCov_filesystem_filesystem.go.html',
                                        },
                                    ],
                                    metrics: {
                                        lineCoverage: {
                                            covered: 0,
                                            uncovered: 11,
                                            coverable: 11,
                                            total: 102,
                                            percentage: 0,
                                        },
                                        methodsCovered: { covered: 0, total: 9, percentage: 0 },
                                        methodsFullyCovered: { covered: 0, total: 9, percentage: 0 },
                                    },
                                    statuses: {
                                        lineCoverage: 'danger',
                                        methodsCovered: 'danger',
                                        methodsFullyCovered: 'danger',
                                    },
                                },
                                {
                                    id: 'github.com/IgorBayerl/AdlerCov/filtering',
                                    name: 'filtering',
                                    type: 'folder',
                                    path: 'github.com/IgorBayerl/AdlerCov/filtering',
                                    children: [
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/filtering/filter.go',
                                            name: 'filter.go',
                                            type: 'file',
                                            path: 'github.com/IgorBayerl/AdlerCov/filtering/filter.go',
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 46,
                                                    uncovered: 33,
                                                    coverable: 79,
                                                    total: 166,
                                                    percentage: 58.22,
                                                },
                                                methodsCovered: { covered: 3, total: 4, percentage: 75 },
                                                methodsFullyCovered: { covered: 0, total: 4, percentage: 0 },
                                            },
                                            statuses: {
                                                lineCoverage: 'danger',
                                                methodsCovered: 'warning',
                                                methodsFullyCovered: 'danger',
                                            },
                                            targetUrl: 'github.com_IgorBayerl_AdlerCov_filtering_filter.go.html',
                                        },
                                    ],
                                    metrics: {
                                        lineCoverage: {
                                            covered: 46,
                                            uncovered: 33,
                                            coverable: 79,
                                            total: 166,
                                            percentage: 58.22,
                                        },
                                        methodsCovered: { covered: 3, total: 4, percentage: 75 },
                                        methodsFullyCovered: { covered: 0, total: 4, percentage: 0 },
                                    },
                                    statuses: {
                                        lineCoverage: 'danger',
                                        methodsCovered: 'warning',
                                        methodsFullyCovered: 'danger',
                                    },
                                },
                                {
                                    id: 'github.com/IgorBayerl/AdlerCov/internal',
                                    name: 'internal',
                                    type: 'folder',
                                    path: 'github.com/IgorBayerl/AdlerCov/internal',
                                    children: [
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/internal/aggregator',
                                            name: 'aggregator',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/internal/aggregator',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/aggregator/aggrgator.go',
                                                    name: 'aggrgator.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/aggregator/aggrgator.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 52,
                                                            uncovered: 0,
                                                            coverable: 52,
                                                            total: 77,
                                                            percentage: 100,
                                                        },
                                                        methodsCovered: { covered: 4, total: 4, percentage: 100 },
                                                        methodsFullyCovered: { covered: 4, total: 4, percentage: 100 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'safe',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_aggregator_aggrgator.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 52,
                                                    uncovered: 0,
                                                    coverable: 52,
                                                    total: 77,
                                                    percentage: 100,
                                                },
                                                methodsCovered: { covered: 4, total: 4, percentage: 100 },
                                                methodsFullyCovered: { covered: 4, total: 4, percentage: 100 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'safe',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/internal/config',
                                            name: 'config',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/internal/config',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/config/config.go',
                                                    name: 'config.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/config/config.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 44,
                                                            uncovered: 8,
                                                            coverable: 52,
                                                            total: 130,
                                                            percentage: 84.61,
                                                        },
                                                        methodsCovered: { covered: 2, total: 2, percentage: 100 },
                                                        methodsFullyCovered: { covered: 1, total: 2, percentage: 50 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_config_config.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 44,
                                                    uncovered: 8,
                                                    coverable: 52,
                                                    total: 130,
                                                    percentage: 84.61,
                                                },
                                                methodsCovered: { covered: 2, total: 2, percentage: 100 },
                                                methodsFullyCovered: { covered: 1, total: 2, percentage: 50 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'danger',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/internal/enricher',
                                            name: 'enricher',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/internal/enricher',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/enricher/enricher.go',
                                                    name: 'enricher.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/enricher/enricher.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 79,
                                                            uncovered: 5,
                                                            coverable: 84,
                                                            total: 168,
                                                            percentage: 94.04,
                                                        },
                                                        methodsCovered: { covered: 7, total: 7, percentage: 100 },
                                                        methodsFullyCovered: {
                                                            covered: 6,
                                                            total: 7,
                                                            percentage: 85.71,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'safe',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_enricher_enricher.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 79,
                                                    uncovered: 5,
                                                    coverable: 84,
                                                    total: 168,
                                                    percentage: 94.04,
                                                },
                                                methodsCovered: { covered: 7, total: 7, percentage: 100 },
                                                methodsFullyCovered: { covered: 6, total: 7, percentage: 85.71 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'safe',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/internal/parsers',
                                            name: 'parsers',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/internal/parsers',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_cobertura',
                                                    name: 'parser_cobertura',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_cobertura',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_cobertura/parser.go',
                                                            name: 'parser.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_cobertura/parser.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 62,
                                                                    uncovered: 23,
                                                                    coverable: 85,
                                                                    total: 133,
                                                                    percentage: 72.94,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 6,
                                                                    total: 6,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 2,
                                                                    total: 6,
                                                                    percentage: 33.33,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'warning',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_parsers_parser_cobertura_parser.go.html',
                                                        },
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_cobertura/processing.go',
                                                            name: 'processing.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_cobertura/processing.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 64,
                                                                    uncovered: 13,
                                                                    coverable: 77,
                                                                    total: 141,
                                                                    percentage: 83.11,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 4,
                                                                    total: 4,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 1,
                                                                    total: 4,
                                                                    percentage: 25,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_parsers_parser_cobertura_processing.go.html',
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 126,
                                                            uncovered: 36,
                                                            coverable: 162,
                                                            total: 274,
                                                            percentage: 77.77,
                                                        },
                                                        methodsCovered: { covered: 10, total: 10, percentage: 100 },
                                                        methodsFullyCovered: { covered: 3, total: 10, percentage: 30 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'warning',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gcov',
                                                    name: 'parser_gcov',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gcov',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gcov/parser.go',
                                                            name: 'parser.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gcov/parser.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 37,
                                                                    uncovered: 8,
                                                                    coverable: 45,
                                                                    total: 79,
                                                                    percentage: 82.22,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 4,
                                                                    total: 4,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 2,
                                                                    total: 4,
                                                                    percentage: 50,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_parsers_parser_gcov_parser.go.html',
                                                        },
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gcov/processing.go',
                                                            name: 'processing.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gcov/processing.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 56,
                                                                    uncovered: 11,
                                                                    coverable: 67,
                                                                    total: 108,
                                                                    percentage: 83.58,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 2,
                                                                    total: 2,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 1,
                                                                    total: 2,
                                                                    percentage: 50,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_parsers_parser_gcov_processing.go.html',
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 93,
                                                            uncovered: 19,
                                                            coverable: 112,
                                                            total: 187,
                                                            percentage: 83.03,
                                                        },
                                                        methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                                        methodsFullyCovered: { covered: 3, total: 6, percentage: 50 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover',
                                                    name: 'parser_gocover',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover/parser.go',
                                                            name: 'parser.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover/parser.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 66,
                                                                    uncovered: 12,
                                                                    coverable: 78,
                                                                    total: 124,
                                                                    percentage: 84.61,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 5,
                                                                    total: 5,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 2,
                                                                    total: 5,
                                                                    percentage: 40,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_parsers_parser_gocover_parser.go.html',
                                                        },
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover/processing.go',
                                                            name: 'processing.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_gocover/processing.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 53,
                                                                    uncovered: 6,
                                                                    coverable: 59,
                                                                    total: 103,
                                                                    percentage: 89.83,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 4,
                                                                    total: 4,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 2,
                                                                    total: 4,
                                                                    percentage: 50,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_parsers_parser_gocover_processing.go.html',
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 119,
                                                            uncovered: 18,
                                                            coverable: 137,
                                                            total: 227,
                                                            percentage: 86.86,
                                                        },
                                                        methodsCovered: { covered: 9, total: 9, percentage: 100 },
                                                        methodsFullyCovered: {
                                                            covered: 4,
                                                            total: 9,
                                                            percentage: 44.44,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/factory.go',
                                                    name: 'factory.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/factory.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 13,
                                                            uncovered: 1,
                                                            coverable: 14,
                                                            total: 25,
                                                            percentage: 92.85,
                                                        },
                                                        methodsCovered: { covered: 2, total: 2, percentage: 100 },
                                                        methodsFullyCovered: { covered: 1, total: 2, percentage: 50 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_parsers_factory.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_config.go',
                                                    name: 'parser_config.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/parsers/parser_config.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 3,
                                                            uncovered: 0,
                                                            coverable: 3,
                                                            total: 44,
                                                            percentage: 100,
                                                        },
                                                        methodsCovered: { covered: 3, total: 3, percentage: 100 },
                                                        methodsFullyCovered: { covered: 3, total: 3, percentage: 100 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'safe',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_parsers_parser_config.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 354,
                                                    uncovered: 74,
                                                    coverable: 428,
                                                    total: 757,
                                                    percentage: 82.71,
                                                },
                                                methodsCovered: { covered: 30, total: 30, percentage: 100 },
                                                methodsFullyCovered: { covered: 14, total: 30, percentage: 46.66 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'danger',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter',
                                            name: 'reporter',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact',
                                                    name: 'htmlreact',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/builder.go',
                                                            name: 'builder.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/builder.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 231,
                                                                    uncovered: 28,
                                                                    coverable: 259,
                                                                    total: 323,
                                                                    percentage: 89.18,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 11,
                                                                    total: 12,
                                                                    percentage: 91.66,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 7,
                                                                    total: 12,
                                                                    percentage: 58.33,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_reporter_htmlreact_builder.go.html',
                                                        },
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/details_generator.go',
                                                            name: 'details_generator.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/details_generator.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 173,
                                                                    uncovered: 24,
                                                                    coverable: 197,
                                                                    total: 258,
                                                                    percentage: 87.81,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 6,
                                                                    total: 6,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 1,
                                                                    total: 6,
                                                                    percentage: 16.66,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_reporter_htmlreact_details_generator.go.html',
                                                        },
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/embed.go',
                                                            name: 'embed.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/embed.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 3,
                                                                    uncovered: 0,
                                                                    coverable: 3,
                                                                    total: 16,
                                                                    percentage: 100,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 1,
                                                                    total: 1,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 1,
                                                                    total: 1,
                                                                    percentage: 100,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'safe',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_reporter_htmlreact_embed.go.html',
                                                        },
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/emit.go',
                                                            name: 'emit.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/emit.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 15,
                                                                    uncovered: 7,
                                                                    coverable: 22,
                                                                    total: 34,
                                                                    percentage: 68.18,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 1,
                                                                    total: 1,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 0,
                                                                    total: 1,
                                                                    percentage: 0,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'warning',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_reporter_htmlreact_emit.go.html',
                                                        },
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/generator.go',
                                                            name: 'generator.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/htmlreact/generator.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 32,
                                                                    uncovered: 29,
                                                                    coverable: 61,
                                                                    total: 90,
                                                                    percentage: 52.45,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 2,
                                                                    total: 2,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 0,
                                                                    total: 2,
                                                                    percentage: 0,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'danger',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_reporter_htmlreact_generator.go.html',
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 454,
                                                            uncovered: 88,
                                                            coverable: 542,
                                                            total: 721,
                                                            percentage: 83.76,
                                                        },
                                                        methodsCovered: { covered: 21, total: 22, percentage: 95.45 },
                                                        methodsFullyCovered: {
                                                            covered: 9,
                                                            total: 22,
                                                            percentage: 40.9,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/lcov',
                                                    name: 'lcov',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/lcov',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/lcov/reporter.go',
                                                            name: 'reporter.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/lcov/reporter.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 72,
                                                                    uncovered: 31,
                                                                    coverable: 103,
                                                                    total: 143,
                                                                    percentage: 69.9,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 4,
                                                                    total: 5,
                                                                    percentage: 80,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 2,
                                                                    total: 5,
                                                                    percentage: 40,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'warning',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_reporter_lcov_reporter.go.html',
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 72,
                                                            uncovered: 31,
                                                            coverable: 103,
                                                            total: 143,
                                                            percentage: 69.9,
                                                        },
                                                        methodsCovered: { covered: 4, total: 5, percentage: 80 },
                                                        methodsFullyCovered: { covered: 2, total: 5, percentage: 40 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'warning',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/reporter_rawjson',
                                                    name: 'reporter_rawjson',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/reporter_rawjson',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/reporter_rawjson/reporter.go',
                                                            name: 'reporter.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/reporter_rawjson/reporter.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 15,
                                                                    uncovered: 7,
                                                                    coverable: 22,
                                                                    total: 44,
                                                                    percentage: 68.18,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 2,
                                                                    total: 3,
                                                                    percentage: 66.66,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 1,
                                                                    total: 3,
                                                                    percentage: 33.33,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'warning',
                                                                methodsCovered: 'warning',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_reporter_reporter_rawjson_reporter.go.html',
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 15,
                                                            uncovered: 7,
                                                            coverable: 22,
                                                            total: 44,
                                                            percentage: 68.18,
                                                        },
                                                        methodsCovered: { covered: 2, total: 3, percentage: 66.66 },
                                                        methodsFullyCovered: {
                                                            covered: 1,
                                                            total: 3,
                                                            percentage: 33.33,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'warning',
                                                        methodsCovered: 'warning',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/textsummary',
                                                    name: 'textsummary',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/textsummary',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/textsummary/reporter.go',
                                                            name: 'reporter.go',
                                                            type: 'file',
                                                            path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/textsummary/reporter.go',
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 65,
                                                                    uncovered: 7,
                                                                    coverable: 72,
                                                                    total: 109,
                                                                    percentage: 90.27,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 3,
                                                                    total: 4,
                                                                    percentage: 75,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 2,
                                                                    total: 4,
                                                                    percentage: 50,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'warning',
                                                                methodsFullyCovered: 'danger',
                                                            },
                                                            targetUrl:
                                                                'github.com_IgorBayerl_AdlerCov_internal_reporter_textsummary_reporter.go.html',
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 65,
                                                            uncovered: 7,
                                                            coverable: 72,
                                                            total: 109,
                                                            percentage: 90.27,
                                                        },
                                                        methodsCovered: { covered: 3, total: 4, percentage: 75 },
                                                        methodsFullyCovered: { covered: 2, total: 4, percentage: 50 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'warning',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/reporter/context.go',
                                                    name: 'context.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/reporter/context.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 10,
                                                            coverable: 10,
                                                            total: 35,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 3, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 3, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_reporter_context.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 606,
                                                    uncovered: 143,
                                                    coverable: 749,
                                                    total: 1052,
                                                    percentage: 80.9,
                                                },
                                                methodsCovered: { covered: 30, total: 37, percentage: 81.08 },
                                                methodsFullyCovered: { covered: 14, total: 37, percentage: 37.83 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'danger',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/internal/tree',
                                            name: 'tree',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/internal/tree',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/tree/builder.go',
                                                    name: 'builder.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/tree/builder.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 95,
                                                            uncovered: 2,
                                                            coverable: 97,
                                                            total: 134,
                                                            percentage: 97.93,
                                                        },
                                                        methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                                        methodsFullyCovered: {
                                                            covered: 5,
                                                            total: 6,
                                                            percentage: 83.33,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'safe',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_tree_builder.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 95,
                                                    uncovered: 2,
                                                    coverable: 97,
                                                    total: 134,
                                                    percentage: 97.93,
                                                },
                                                methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                                methodsFullyCovered: { covered: 5, total: 6, percentage: 83.33 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'safe',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/internal/utils',
                                            name: 'utils',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/internal/utils',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/utils/analyzer.go',
                                                    name: 'analyzer.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/utils/analyzer.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 5,
                                                            uncovered: 12,
                                                            coverable: 17,
                                                            total: 33,
                                                            percentage: 29.41,
                                                        },
                                                        methodsCovered: { covered: 1, total: 3, percentage: 33.33 },
                                                        methodsFullyCovered: {
                                                            covered: 1,
                                                            total: 3,
                                                            percentage: 33.33,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_utils_analyzer.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/utils/brace_finder.go',
                                                    name: 'brace_finder.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/utils/brace_finder.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 45,
                                                            coverable: 45,
                                                            total: 66,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 1, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 1, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_utils_brace_finder.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/utils/collections.go',
                                                    name: 'collections.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/utils/collections.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 13,
                                                            coverable: 13,
                                                            total: 20,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 1, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 1, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_utils_collections.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/utils/encoding.go',
                                                    name: 'encoding.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/utils/encoding.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 26,
                                                            coverable: 26,
                                                            total: 44,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 1, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 1, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_utils_encoding.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/utils/line_sorter.go',
                                                    name: 'line_sorter.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/utils/line_sorter.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 19,
                                                            coverable: 19,
                                                            total: 40,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 1, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 1, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_utils_line_sorter.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/utils/math.go',
                                                    name: 'math.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/utils/math.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 37,
                                                            uncovered: 9,
                                                            coverable: 46,
                                                            total: 65,
                                                            percentage: 80.43,
                                                        },
                                                        methodsCovered: { covered: 2, total: 2, percentage: 100 },
                                                        methodsFullyCovered: { covered: 0, total: 2, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_utils_math.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/utils/paths.go',
                                                    name: 'paths.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/utils/paths.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 63,
                                                            uncovered: 9,
                                                            coverable: 72,
                                                            total: 109,
                                                            percentage: 87.5,
                                                        },
                                                        methodsCovered: { covered: 2, total: 2, percentage: 100 },
                                                        methodsFullyCovered: { covered: 0, total: 2, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_utils_paths.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/internal/utils/stringutils.go',
                                                    name: 'stringutils.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/internal/utils/stringutils.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 85,
                                                            coverable: 85,
                                                            total: 137,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 4, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 4, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_internal_utils_stringutils.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 105,
                                                    uncovered: 218,
                                                    coverable: 323,
                                                    total: 514,
                                                    percentage: 32.5,
                                                },
                                                methodsCovered: { covered: 5, total: 15, percentage: 33.33 },
                                                methodsFullyCovered: { covered: 1, total: 15, percentage: 6.66 },
                                            },
                                            statuses: {
                                                lineCoverage: 'danger',
                                                methodsCovered: 'danger',
                                                methodsFullyCovered: 'danger',
                                            },
                                        },
                                    ],
                                    metrics: {
                                        lineCoverage: {
                                            covered: 1335,
                                            uncovered: 450,
                                            coverable: 1785,
                                            total: 2832,
                                            percentage: 74.78,
                                        },
                                        methodsCovered: { covered: 84, total: 101, percentage: 83.16 },
                                        methodsFullyCovered: { covered: 45, total: 101, percentage: 44.55 },
                                    },
                                    statuses: {
                                        lineCoverage: 'warning',
                                        methodsCovered: 'safe',
                                        methodsFullyCovered: 'danger',
                                    },
                                },
                                {
                                    id: 'github.com/IgorBayerl/AdlerCov/logging',
                                    name: 'logging',
                                    type: 'folder',
                                    path: 'github.com/IgorBayerl/AdlerCov/logging',
                                    children: [
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/logging/logging.go',
                                            name: 'logging.go',
                                            type: 'file',
                                            path: 'github.com/IgorBayerl/AdlerCov/logging/logging.go',
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 49,
                                                    uncovered: 48,
                                                    coverable: 97,
                                                    total: 178,
                                                    percentage: 50.51,
                                                },
                                                methodsCovered: { covered: 7, total: 10, percentage: 70 },
                                                methodsFullyCovered: { covered: 4, total: 10, percentage: 40 },
                                            },
                                            statuses: {
                                                lineCoverage: 'danger',
                                                methodsCovered: 'warning',
                                                methodsFullyCovered: 'danger',
                                            },
                                            targetUrl: 'github.com_IgorBayerl_AdlerCov_logging_logging.go.html',
                                        },
                                    ],
                                    metrics: {
                                        lineCoverage: {
                                            covered: 49,
                                            uncovered: 48,
                                            coverable: 97,
                                            total: 178,
                                            percentage: 50.51,
                                        },
                                        methodsCovered: { covered: 7, total: 10, percentage: 70 },
                                        methodsFullyCovered: { covered: 4, total: 10, percentage: 40 },
                                    },
                                    statuses: {
                                        lineCoverage: 'danger',
                                        methodsCovered: 'warning',
                                        methodsFullyCovered: 'danger',
                                    },
                                },
                                {
                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter',
                                    name: 'tree-sitter',
                                    type: 'folder',
                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter',
                                    children: [
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter',
                                            name: 'go-tree-sitter',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/allocator.go',
                                                    name: 'allocator.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/allocator.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 45,
                                                            uncovered: 20,
                                                            coverable: 65,
                                                            total: 111,
                                                            percentage: 69.23,
                                                        },
                                                        methodsCovered: { covered: 6, total: 6, percentage: 100 },
                                                        methodsFullyCovered: {
                                                            covered: 4,
                                                            total: 6,
                                                            percentage: 66.66,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'warning',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'warning',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_allocator.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/dup_windows.go',
                                                    name: 'dup_windows.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/dup_windows.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 4,
                                                            coverable: 4,
                                                            total: 16,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 1, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 1, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_dup_windows.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/edit.go',
                                                    name: 'edit.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/edit.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 10,
                                                            coverable: 10,
                                                            total: 27,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 1, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 1, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_edit.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/language.go',
                                                    name: 'language.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/language.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 6,
                                                            uncovered: 55,
                                                            coverable: 61,
                                                            total: 157,
                                                            percentage: 9.83,
                                                        },
                                                        methodsCovered: { covered: 2, total: 17, percentage: 11.76 },
                                                        methodsFullyCovered: {
                                                            covered: 2,
                                                            total: 17,
                                                            percentage: 11.76,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_language.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/lookahead_iterator.go',
                                                    name: 'lookahead_iterator.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/lookahead_iterator.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 33,
                                                            coverable: 33,
                                                            total: 71,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 9, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 9, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_lookahead_iterator.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/node.go',
                                                    name: 'node.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/node.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 36,
                                                            uncovered: 170,
                                                            coverable: 206,
                                                            total: 392,
                                                            percentage: 17.47,
                                                        },
                                                        methodsCovered: { covered: 10, total: 51, percentage: 19.6 },
                                                        methodsFullyCovered: {
                                                            covered: 9,
                                                            total: 51,
                                                            percentage: 17.64,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_node.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/parser.go',
                                                    name: 'parser.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/parser.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 74,
                                                            uncovered: 298,
                                                            coverable: 372,
                                                            total: 754,
                                                            percentage: 19.89,
                                                        },
                                                        methodsCovered: { covered: 8, total: 35, percentage: 22.85 },
                                                        methodsFullyCovered: {
                                                            covered: 5,
                                                            total: 35,
                                                            percentage: 14.28,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_parser.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/point.go',
                                                    name: 'point.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/point.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 4,
                                                            uncovered: 9,
                                                            coverable: 13,
                                                            total: 31,
                                                            percentage: 30.76,
                                                        },
                                                        methodsCovered: { covered: 1, total: 3, percentage: 33.33 },
                                                        methodsFullyCovered: {
                                                            covered: 1,
                                                            total: 3,
                                                            percentage: 33.33,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_point.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/query.go',
                                                    name: 'query.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/query.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 159,
                                                            uncovered: 511,
                                                            coverable: 670,
                                                            total: 1092,
                                                            percentage: 23.73,
                                                        },
                                                        methodsCovered: { covered: 11, total: 47, percentage: 23.4 },
                                                        methodsFullyCovered: {
                                                            covered: 6,
                                                            total: 47,
                                                            percentage: 12.76,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_query.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/ranges.go',
                                                    name: 'ranges.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/ranges.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 17,
                                                            coverable: 17,
                                                            total: 42,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 3, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 3, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_ranges.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/tree.go',
                                                    name: 'tree.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/tree.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 10,
                                                            uncovered: 48,
                                                            coverable: 58,
                                                            total: 121,
                                                            percentage: 17.24,
                                                        },
                                                        methodsCovered: { covered: 3, total: 11, percentage: 27.27 },
                                                        methodsFullyCovered: {
                                                            covered: 3,
                                                            total: 11,
                                                            percentage: 27.27,
                                                        },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_tree.go.html',
                                                },
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/tree_cursor.go',
                                                    name: 'tree_cursor.go',
                                                    type: 'file',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/go-tree-sitter/tree_cursor.go',
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 0,
                                                            uncovered: 62,
                                                            coverable: 62,
                                                            total: 161,
                                                            percentage: 0,
                                                        },
                                                        methodsCovered: { covered: 0, total: 18, percentage: 0 },
                                                        methodsFullyCovered: { covered: 0, total: 18, percentage: 0 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'danger',
                                                        methodsCovered: 'danger',
                                                        methodsFullyCovered: 'danger',
                                                    },
                                                    targetUrl:
                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_go-tree-sitter_tree_cursor.go.html',
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 334,
                                                    uncovered: 1237,
                                                    coverable: 1571,
                                                    total: 2975,
                                                    percentage: 21.26,
                                                },
                                                methodsCovered: { covered: 41, total: 202, percentage: 20.29 },
                                                methodsFullyCovered: { covered: 30, total: 202, percentage: 14.85 },
                                            },
                                            statuses: {
                                                lineCoverage: 'danger',
                                                methodsCovered: 'danger',
                                                methodsFullyCovered: 'danger',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-cpp',
                                            name: 'tree-sitter-cpp',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-cpp',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-cpp/bindings',
                                                    name: 'bindings',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-cpp/bindings',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-cpp/bindings/go',
                                                            name: 'go',
                                                            type: 'folder',
                                                            path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-cpp/bindings/go',
                                                            children: [
                                                                {
                                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-cpp/bindings/go/binding.go',
                                                                    name: 'binding.go',
                                                                    type: 'file',
                                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-cpp/bindings/go/binding.go',
                                                                    metrics: {
                                                                        lineCoverage: {
                                                                            covered: 3,
                                                                            uncovered: 0,
                                                                            coverable: 3,
                                                                            total: 13,
                                                                            percentage: 100,
                                                                        },
                                                                        methodsCovered: {
                                                                            covered: 1,
                                                                            total: 1,
                                                                            percentage: 100,
                                                                        },
                                                                        methodsFullyCovered: {
                                                                            covered: 1,
                                                                            total: 1,
                                                                            percentage: 100,
                                                                        },
                                                                    },
                                                                    statuses: {
                                                                        lineCoverage: 'safe',
                                                                        methodsCovered: 'safe',
                                                                        methodsFullyCovered: 'safe',
                                                                    },
                                                                    targetUrl:
                                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_tree-sitter-cpp_bindings_go_binding.go.html',
                                                                },
                                                            ],
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 3,
                                                                    uncovered: 0,
                                                                    coverable: 3,
                                                                    total: 13,
                                                                    percentage: 100,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 1,
                                                                    total: 1,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 1,
                                                                    total: 1,
                                                                    percentage: 100,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'safe',
                                                            },
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 3,
                                                            uncovered: 0,
                                                            coverable: 3,
                                                            total: 13,
                                                            percentage: 100,
                                                        },
                                                        methodsCovered: { covered: 1, total: 1, percentage: 100 },
                                                        methodsFullyCovered: { covered: 1, total: 1, percentage: 100 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'safe',
                                                    },
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 3,
                                                    uncovered: 0,
                                                    coverable: 3,
                                                    total: 13,
                                                    percentage: 100,
                                                },
                                                methodsCovered: { covered: 1, total: 1, percentage: 100 },
                                                methodsFullyCovered: { covered: 1, total: 1, percentage: 100 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'safe',
                                            },
                                        },
                                        {
                                            id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-go',
                                            name: 'tree-sitter-go',
                                            type: 'folder',
                                            path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-go',
                                            children: [
                                                {
                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-go/bindings',
                                                    name: 'bindings',
                                                    type: 'folder',
                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-go/bindings',
                                                    children: [
                                                        {
                                                            id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-go/bindings/go',
                                                            name: 'go',
                                                            type: 'folder',
                                                            path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-go/bindings/go',
                                                            children: [
                                                                {
                                                                    id: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-go/bindings/go/binding.go',
                                                                    name: 'binding.go',
                                                                    type: 'file',
                                                                    path: 'github.com/IgorBayerl/AdlerCov/tree-sitter/tree-sitter-go/bindings/go/binding.go',
                                                                    metrics: {
                                                                        lineCoverage: {
                                                                            covered: 3,
                                                                            uncovered: 0,
                                                                            coverable: 3,
                                                                            total: 15,
                                                                            percentage: 100,
                                                                        },
                                                                        methodsCovered: {
                                                                            covered: 1,
                                                                            total: 1,
                                                                            percentage: 100,
                                                                        },
                                                                        methodsFullyCovered: {
                                                                            covered: 1,
                                                                            total: 1,
                                                                            percentage: 100,
                                                                        },
                                                                    },
                                                                    statuses: {
                                                                        lineCoverage: 'safe',
                                                                        methodsCovered: 'safe',
                                                                        methodsFullyCovered: 'safe',
                                                                    },
                                                                    targetUrl:
                                                                        'github.com_IgorBayerl_AdlerCov_tree-sitter_tree-sitter-go_bindings_go_binding.go.html',
                                                                },
                                                            ],
                                                            metrics: {
                                                                lineCoverage: {
                                                                    covered: 3,
                                                                    uncovered: 0,
                                                                    coverable: 3,
                                                                    total: 15,
                                                                    percentage: 100,
                                                                },
                                                                methodsCovered: {
                                                                    covered: 1,
                                                                    total: 1,
                                                                    percentage: 100,
                                                                },
                                                                methodsFullyCovered: {
                                                                    covered: 1,
                                                                    total: 1,
                                                                    percentage: 100,
                                                                },
                                                            },
                                                            statuses: {
                                                                lineCoverage: 'safe',
                                                                methodsCovered: 'safe',
                                                                methodsFullyCovered: 'safe',
                                                            },
                                                        },
                                                    ],
                                                    metrics: {
                                                        lineCoverage: {
                                                            covered: 3,
                                                            uncovered: 0,
                                                            coverable: 3,
                                                            total: 15,
                                                            percentage: 100,
                                                        },
                                                        methodsCovered: { covered: 1, total: 1, percentage: 100 },
                                                        methodsFullyCovered: { covered: 1, total: 1, percentage: 100 },
                                                    },
                                                    statuses: {
                                                        lineCoverage: 'safe',
                                                        methodsCovered: 'safe',
                                                        methodsFullyCovered: 'safe',
                                                    },
                                                },
                                            ],
                                            metrics: {
                                                lineCoverage: {
                                                    covered: 3,
                                                    uncovered: 0,
                                                    coverable: 3,
                                                    total: 15,
                                                    percentage: 100,
                                                },
                                                methodsCovered: { covered: 1, total: 1, percentage: 100 },
                                                methodsFullyCovered: { covered: 1, total: 1, percentage: 100 },
                                            },
                                            statuses: {
                                                lineCoverage: 'safe',
                                                methodsCovered: 'safe',
                                                methodsFullyCovered: 'safe',
                                            },
                                        },
                                    ],
                                    metrics: {
                                        lineCoverage: {
                                            covered: 340,
                                            uncovered: 1237,
                                            coverable: 1577,
                                            total: 3003,
                                            percentage: 21.55,
                                        },
                                        methodsCovered: { covered: 43, total: 204, percentage: 21.07 },
                                        methodsFullyCovered: { covered: 32, total: 204, percentage: 15.68 },
                                    },
                                    statuses: {
                                        lineCoverage: 'danger',
                                        methodsCovered: 'danger',
                                        methodsFullyCovered: 'danger',
                                    },
                                },
                            ],
                            metrics: {
                                lineCoverage: {
                                    covered: 2011,
                                    uncovered: 1824,
                                    coverable: 3835,
                                    total: 6798,
                                    percentage: 52.43,
                                },
                                methodsCovered: { covered: 156, total: 348, percentage: 44.82 },
                                methodsFullyCovered: { covered: 92, total: 348, percentage: 26.43 },
                            },
                            statuses: {
                                lineCoverage: 'danger',
                                methodsCovered: 'danger',
                                methodsFullyCovered: 'danger',
                            },
                        },
                    ],
                    metrics: {
                        lineCoverage: {
                            covered: 2011,
                            uncovered: 1824,
                            coverable: 3835,
                            total: 6798,
                            percentage: 52.43,
                        },
                        methodsCovered: { covered: 156, total: 348, percentage: 44.82 },
                        methodsFullyCovered: { covered: 92, total: 348, percentage: 26.43 },
                    },
                    statuses: { lineCoverage: 'danger', methodsCovered: 'danger', methodsFullyCovered: 'danger' },
                },
            ],
            metrics: {
                lineCoverage: { covered: 2011, uncovered: 1824, coverable: 3835, total: 6798, percentage: 52.43 },
                methodsCovered: { covered: 156, total: 348, percentage: 44.82 },
                methodsFullyCovered: { covered: 92, total: 348, percentage: 26.43 },
            },
            statuses: { lineCoverage: 'danger', methodsCovered: 'danger', methodsFullyCovered: 'danger' },
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
        maxCyclomaticComplexity: {
            label: 'Max Cyclomatic Complexity',
            shortLabel: 'Max Complexity',
            subMetrics: [{ id: 'total', label: 'Value', width: 100 }],
        },
        methodBranchCoverage: {
            label: 'Method Branches',
            shortLabel: 'Method Branches',
            subMetrics: [
                { id: 'covered', label: 'Covered', width: 100 },
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
        { label: 'Generated At', value: '2025-09-13 10:54:29' },
        { label: 'Parser', value: 'GoCover' },
    ],
}
