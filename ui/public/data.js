// /**
//  * NOTE: This is a developer data fixture for testing purposes.
//  * This file provides a larger, more complex dataset to test rendering,
//  * filtering, and performance of the coverage report UI.
//  * Schema version: 1
//  */
window.__NANOVISION_SUMMARY__={
  "schemaVersion": 1,
  "generatedAt": "2025-10-28T18:39:29Z",
  "title": "Coverage Report",
  "totals": {
    "lineCoverage": { "covered": 2257, "uncovered": 779, "coverable": 3036, "total": 5333, "percentage": 74.34 },
    "branchCoverage": { "covered": 39, "total": 58, "percentage": 67.24 },
    "methodsCovered": { "covered": 150, "total": 216, "percentage": 69.44 },
    "methodsFullyCovered": { "covered": 83, "total": 216, "percentage": 38.42 },
    "files": 64,
    "folders": 36,
    "statuses": {
      "branchCoverage": "warning",
      "lineCoverage": "warning",
      "methodsCovered": "warning",
      "methodsFullyCovered": "danger"
    }
  },
  "tree": [
    {
      "id": "analyzer",
      "name": "analyzer",
      "type": "folder",
      "path": "analyzer",
      "children": [
        {
          "id": "analyzer/cpp",
          "name": "cpp",
          "type": "folder",
          "path": "analyzer/cpp",
          "children": [
            {
              "id": "analyzer/cpp/analyzer.go",
              "name": "analyzer.go",
              "type": "file",
              "path": "analyzer/cpp/analyzer.go",
              "metrics": {
                "lineCoverage": { "covered": 97, "uncovered": 13, "coverable": 110, "total": 194, "percentage": 88.18 },
                "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
                "methodsFullyCovered": { "covered": 3, "total": 6, "percentage": 50 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
              "targetUrl": "analyzer_cpp_analyzer.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 97, "uncovered": 13, "coverable": 110, "total": 194, "percentage": 88.18 },
            "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
            "methodsFullyCovered": { "covered": 3, "total": 6, "percentage": 50 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
        },
        {
          "id": "analyzer/go",
          "name": "go",
          "type": "folder",
          "path": "analyzer/go",
          "children": [
            {
              "id": "analyzer/go/analyzer.go",
              "name": "analyzer.go",
              "type": "file",
              "path": "analyzer/go/analyzer.go",
              "metrics": {
                "lineCoverage": { "covered": 82, "uncovered": 13, "coverable": 95, "total": 164, "percentage": 86.31 },
                "methodsCovered": { "covered": 5, "total": 5, "percentage": 100 },
                "methodsFullyCovered": { "covered": 3, "total": 5, "percentage": 60 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "warning" },
              "targetUrl": "analyzer_go_analyzer.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 82, "uncovered": 13, "coverable": 95, "total": 164, "percentage": 86.31 },
            "methodsCovered": { "covered": 5, "total": 5, "percentage": 100 },
            "methodsFullyCovered": { "covered": 3, "total": 5, "percentage": 60 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "warning" }
        },
        {
          "id": "analyzer/analyzer.go",
          "name": "analyzer.go",
          "type": "file",
          "path": "analyzer/analyzer.go",
          "metrics": {
            "lineCoverage": { "covered": 0, "uncovered": 3, "coverable": 3, "total": 36, "percentage": 0 },
            "methodsCovered": { "covered": 0, "total": 1, "percentage": 0 },
            "methodsFullyCovered": { "covered": 0, "total": 1, "percentage": 0 }
          },
          "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
          "targetUrl": "analyzer_analyzer.go.html"
        }
      ],
      "metrics": {
        "lineCoverage": { "covered": 179, "uncovered": 29, "coverable": 208, "total": 394, "percentage": 86.05 },
        "methodsCovered": { "covered": 11, "total": 12, "percentage": 91.66 },
        "methodsFullyCovered": { "covered": 6, "total": 12, "percentage": 50 }
      },
      "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
    },
    {
      "id": "cmd",
      "name": "cmd",
      "type": "folder",
      "path": "cmd",
      "children": [
        {
          "id": "cmd/main.go",
          "name": "main.go",
          "type": "file",
          "path": "cmd/main.go",
          "metrics": {
            "lineCoverage": { "covered": 143, "uncovered": 58, "coverable": 201, "total": 287, "percentage": 71.14 },
            "methodsCovered": { "covered": 7, "total": 7, "percentage": 100 },
            "methodsFullyCovered": { "covered": 2, "total": 7, "percentage": 28.57 }
          },
          "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
          "targetUrl": "cmd_main.go.html"
        }
      ],
      "metrics": {
        "lineCoverage": { "covered": 143, "uncovered": 58, "coverable": 201, "total": 287, "percentage": 71.14 },
        "methodsCovered": { "covered": 7, "total": 7, "percentage": 100 },
        "methodsFullyCovered": { "covered": 2, "total": 7, "percentage": 28.57 }
      },
      "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
    },
    {
      "id": "demo_projects",
      "name": "demo_projects",
      "type": "folder",
      "path": "demo_projects",
      "children": [
        {
          "id": "demo_projects/cpp",
          "name": "cpp",
          "type": "folder",
          "path": "demo_projects/cpp",
          "children": [
            {
              "id": "demo_projects/cpp/project",
              "name": "project",
              "type": "folder",
              "path": "demo_projects/cpp/project",
              "children": [
                {
                  "id": "demo_projects/cpp/project/src",
                  "name": "src",
                  "type": "folder",
                  "path": "demo_projects/cpp/project/src",
                  "children": [
                    {
                      "id": "demo_projects/cpp/project/src/utils",
                      "name": "utils",
                      "type": "folder",
                      "path": "demo_projects/cpp/project/src/utils",
                      "children": [
                        {
                          "id": "demo_projects/cpp/project/src/utils/math_utils.cpp",
                          "name": "math_utils.cpp",
                          "type": "file",
                          "path": "demo_projects/cpp/project/src/utils/math_utils.cpp",
                          "metrics": {
                            "branchCoverage": { "covered": 9, "total": 14, "percentage": 64.28 },
                            "lineCoverage": {
                              "covered": 13,
                              "uncovered": 3,
                              "coverable": 16,
                              "total": 33,
                              "percentage": 81.25
                            },
                            "methodsCovered": { "covered": 2, "total": 2, "percentage": 100 },
                            "methodsFullyCovered": { "covered": 0, "total": 2, "percentage": 0 }
                          },
                          "statuses": {
                            "branchCoverage": "warning",
                            "lineCoverage": "safe",
                            "methodsCovered": "safe",
                            "methodsFullyCovered": "danger"
                          },
                          "targetUrl": "demo_projects_cpp_project_src_utils_math_utils.cpp.html"
                        }
                      ],
                      "metrics": {
                        "branchCoverage": { "covered": 9, "total": 14, "percentage": 64.28 },
                        "lineCoverage": {
                          "covered": 13,
                          "uncovered": 3,
                          "coverable": 16,
                          "total": 33,
                          "percentage": 81.25
                        },
                        "methodsCovered": { "covered": 2, "total": 2, "percentage": 100 },
                        "methodsFullyCovered": { "covered": 0, "total": 2, "percentage": 0 }
                      },
                      "statuses": {
                        "branchCoverage": "warning",
                        "lineCoverage": "safe",
                        "methodsCovered": "safe",
                        "methodsFullyCovered": "danger"
                      }
                    },
                    {
                      "id": "demo_projects/cpp/project/src/advanced_calculator.cpp",
                      "name": "advanced_calculator.cpp",
                      "type": "file",
                      "path": "demo_projects/cpp/project/src/advanced_calculator.cpp",
                      "metrics": {
                        "branchCoverage": { "covered": 6, "total": 12, "percentage": 50 },
                        "lineCoverage": {
                          "covered": 13,
                          "uncovered": 4,
                          "coverable": 17,
                          "total": 32,
                          "percentage": 76.47
                        },
                        "methodsCovered": { "covered": 2, "total": 2, "percentage": 100 },
                        "methodsFullyCovered": { "covered": 0, "total": 2, "percentage": 0 }
                      },
                      "statuses": {
                        "branchCoverage": "danger",
                        "lineCoverage": "warning",
                        "methodsCovered": "safe",
                        "methodsFullyCovered": "danger"
                      },
                      "targetUrl": "demo_projects_cpp_project_src_advanced_calculator.cpp.html"
                    },
                    {
                      "id": "demo_projects/cpp/project/src/calculator.cpp",
                      "name": "calculator.cpp",
                      "type": "file",
                      "path": "demo_projects/cpp/project/src/calculator.cpp",
                      "metrics": {
                        "branchCoverage": { "covered": 6, "total": 8, "percentage": 75 },
                        "lineCoverage": {
                          "covered": 13,
                          "uncovered": 3,
                          "coverable": 16,
                          "total": 32,
                          "percentage": 81.25
                        },
                        "methodsCovered": { "covered": 4, "total": 5, "percentage": 80 },
                        "methodsFullyCovered": { "covered": 3, "total": 5, "percentage": 60 }
                      },
                      "statuses": {
                        "branchCoverage": "warning",
                        "lineCoverage": "safe",
                        "methodsCovered": "safe",
                        "methodsFullyCovered": "warning"
                      },
                      "targetUrl": "demo_projects_cpp_project_src_calculator.cpp.html"
                    }
                  ],
                  "metrics": {
                    "branchCoverage": { "covered": 21, "total": 34, "percentage": 61.76 },
                    "lineCoverage": {
                      "covered": 39,
                      "uncovered": 10,
                      "coverable": 49,
                      "total": 97,
                      "percentage": 79.59
                    },
                    "methodsCovered": { "covered": 8, "total": 9, "percentage": 88.88 },
                    "methodsFullyCovered": { "covered": 3, "total": 9, "percentage": 33.33 }
                  },
                  "statuses": {
                    "branchCoverage": "warning",
                    "lineCoverage": "warning",
                    "methodsCovered": "safe",
                    "methodsFullyCovered": "danger"
                  }
                }
              ],
              "metrics": {
                "branchCoverage": { "covered": 21, "total": 34, "percentage": 61.76 },
                "lineCoverage": { "covered": 39, "uncovered": 10, "coverable": 49, "total": 97, "percentage": 79.59 },
                "methodsCovered": { "covered": 8, "total": 9, "percentage": 88.88 },
                "methodsFullyCovered": { "covered": 3, "total": 9, "percentage": 33.33 }
              },
              "statuses": {
                "branchCoverage": "warning",
                "lineCoverage": "warning",
                "methodsCovered": "safe",
                "methodsFullyCovered": "danger"
              }
            }
          ],
          "metrics": {
            "branchCoverage": { "covered": 21, "total": 34, "percentage": 61.76 },
            "lineCoverage": { "covered": 39, "uncovered": 10, "coverable": 49, "total": 97, "percentage": 79.59 },
            "methodsCovered": { "covered": 8, "total": 9, "percentage": 88.88 },
            "methodsFullyCovered": { "covered": 3, "total": 9, "percentage": 33.33 }
          },
          "statuses": {
            "branchCoverage": "warning",
            "lineCoverage": "warning",
            "methodsCovered": "safe",
            "methodsFullyCovered": "danger"
          }
        },
        {
          "id": "demo_projects/csharp",
          "name": "csharp",
          "type": "folder",
          "path": "demo_projects/csharp",
          "children": [
            {
              "id": "demo_projects/csharp/project",
              "name": "project",
              "type": "folder",
              "path": "demo_projects/csharp/project",
              "children": [
                {
                  "id": "demo_projects/csharp/project/Test",
                  "name": "Test",
                  "type": "folder",
                  "path": "demo_projects/csharp/project/Test",
                  "children": [
                    {
                      "id": "demo_projects/csharp/project/Test/AbstractClass.cs",
                      "name": "AbstractClass.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/AbstractClass.cs",
                      "metrics": {
                        "lineCoverage": {
                          "covered": 9,
                          "uncovered": 4,
                          "coverable": 13,
                          "total": 54,
                          "percentage": 69.23
                        }
                      },
                      "statuses": { "lineCoverage": "warning" },
                      "targetUrl": "demo_projects_csharp_project_Test_AbstractClass.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/AnalyzerTestClass.cs",
                      "name": "AnalyzerTestClass.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/AnalyzerTestClass.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 0, "uncovered": 7, "coverable": 7, "total": 46, "percentage": 0 }
                      },
                      "statuses": { "lineCoverage": "danger" },
                      "targetUrl": "demo_projects_csharp_project_Test_AnalyzerTestClass.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/AsyncClass.cs",
                      "name": "AsyncClass.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/AsyncClass.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 3, "uncovered": 0, "coverable": 3, "total": 15, "percentage": 100 }
                      },
                      "statuses": { "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_AsyncClass.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/ClassWithExcludes.cs",
                      "name": "ClassWithExcludes.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/ClassWithExcludes.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 6, "uncovered": 0, "coverable": 6, "total": 22, "percentage": 100 }
                      },
                      "statuses": { "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_ClassWithExcludes.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/ClassWithLocalFunctions.cs",
                      "name": "ClassWithLocalFunctions.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/ClassWithLocalFunctions.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 7, "uncovered": 0, "coverable": 7, "total": 26, "percentage": 100 }
                      },
                      "statuses": { "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_ClassWithLocalFunctions.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/CodeContract_Contract.cs",
                      "name": "CodeContract_Contract.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/CodeContract_Contract.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 0, "uncovered": 2, "coverable": 2, "total": 16, "percentage": 0 }
                      },
                      "statuses": { "lineCoverage": "danger" },
                      "targetUrl": "demo_projects_csharp_project_Test_CodeContract_Contract.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/CodeContract_Target.cs",
                      "name": "CodeContract_Target.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/CodeContract_Target.cs",
                      "metrics": {
                        "branchCoverage": { "covered": 4, "total": 4, "percentage": 100 },
                        "lineCoverage": { "covered": 3, "uncovered": 0, "coverable": 3, "total": 18, "percentage": 100 }
                      },
                      "statuses": { "branchCoverage": "safe", "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_CodeContract_Target.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/GenericAsyncClass.cs",
                      "name": "GenericAsyncClass.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/GenericAsyncClass.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 2, "uncovered": 0, "coverable": 2, "total": 13, "percentage": 100 }
                      },
                      "statuses": { "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_GenericAsyncClass.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/GenericClass.cs",
                      "name": "GenericClass.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/GenericClass.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 4, "uncovered": 1, "coverable": 5, "total": 49, "percentage": 80 }
                      },
                      "statuses": { "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_GenericClass.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/NotMatchingFileName.cs",
                      "name": "NotMatchingFileName.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/NotMatchingFileName.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 1, "uncovered": 0, "coverable": 1, "total": 8, "percentage": 100 }
                      },
                      "statuses": { "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_NotMatchingFileName.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/PartialClass.cs",
                      "name": "PartialClass.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/PartialClass.cs",
                      "metrics": {
                        "branchCoverage": { "covered": 2, "total": 4, "percentage": 50 },
                        "lineCoverage": {
                          "covered": 4,
                          "uncovered": 5,
                          "coverable": 9,
                          "total": 36,
                          "percentage": 44.44
                        }
                      },
                      "statuses": { "branchCoverage": "danger", "lineCoverage": "danger" },
                      "targetUrl": "demo_projects_csharp_project_Test_PartialClass.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/PartialClass2.cs",
                      "name": "PartialClass2.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/PartialClass2.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 2, "uncovered": 2, "coverable": 4, "total": 17, "percentage": 50 }
                      },
                      "statuses": { "lineCoverage": "danger" },
                      "targetUrl": "demo_projects_csharp_project_Test_PartialClass2.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/PartialClassWithAutoProperties.cs",
                      "name": "PartialClassWithAutoProperties.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/PartialClassWithAutoProperties.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 1, "uncovered": 0, "coverable": 1, "total": 8, "percentage": 100 }
                      },
                      "statuses": { "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_PartialClassWithAutoProperties.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/PartialClassWithAutoProperties2.cs",
                      "name": "PartialClassWithAutoProperties2.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/PartialClassWithAutoProperties2.cs",
                      "metrics": {
                        "lineCoverage": { "covered": 1, "uncovered": 1, "coverable": 2, "total": 10, "percentage": 50 }
                      },
                      "statuses": { "lineCoverage": "danger" },
                      "targetUrl": "demo_projects_csharp_project_Test_PartialClassWithAutoProperties2.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/Program.cs",
                      "name": "Program.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/Program.cs",
                      "metrics": {
                        "lineCoverage": {
                          "covered": 33,
                          "uncovered": 2,
                          "coverable": 35,
                          "total": 75,
                          "percentage": 94.28
                        }
                      },
                      "statuses": { "lineCoverage": "safe" },
                      "targetUrl": "demo_projects_csharp_project_Test_Program.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/TestClass.cs",
                      "name": "TestClass.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/TestClass.cs",
                      "metrics": {
                        "branchCoverage": { "covered": 4, "total": 8, "percentage": 50 },
                        "lineCoverage": {
                          "covered": 16,
                          "uncovered": 7,
                          "coverable": 23,
                          "total": 60,
                          "percentage": 69.56
                        }
                      },
                      "statuses": { "branchCoverage": "danger", "lineCoverage": "warning" },
                      "targetUrl": "demo_projects_csharp_project_Test_TestClass.cs.html"
                    },
                    {
                      "id": "demo_projects/csharp/project/Test/TestClass2.cs",
                      "name": "TestClass2.cs",
                      "type": "file",
                      "path": "demo_projects/csharp/project/Test/TestClass2.cs",
                      "metrics": {
                        "branchCoverage": { "covered": 8, "total": 8, "percentage": 100 },
                        "lineCoverage": {
                          "covered": 17,
                          "uncovered": 10,
                          "coverable": 27,
                          "total": 85,
                          "percentage": 62.96
                        }
                      },
                      "statuses": { "branchCoverage": "safe", "lineCoverage": "warning" },
                      "targetUrl": "demo_projects_csharp_project_Test_TestClass2.cs.html"
                    }
                  ],
                  "metrics": {
                    "branchCoverage": { "covered": 18, "total": 24, "percentage": 75 },
                    "lineCoverage": {
                      "covered": 109,
                      "uncovered": 41,
                      "coverable": 150,
                      "total": 558,
                      "percentage": 72.66
                    }
                  },
                  "statuses": { "branchCoverage": "warning", "lineCoverage": "warning" }
                }
              ],
              "metrics": {
                "branchCoverage": { "covered": 18, "total": 24, "percentage": 75 },
                "lineCoverage": { "covered": 109, "uncovered": 41, "coverable": 150, "total": 558, "percentage": 72.66 }
              },
              "statuses": { "branchCoverage": "warning", "lineCoverage": "warning" }
            }
          ],
          "metrics": {
            "branchCoverage": { "covered": 18, "total": 24, "percentage": 75 },
            "lineCoverage": { "covered": 109, "uncovered": 41, "coverable": 150, "total": 558, "percentage": 72.66 }
          },
          "statuses": { "branchCoverage": "warning", "lineCoverage": "warning" }
        },
        {
          "id": "demo_projects/go",
          "name": "go",
          "type": "folder",
          "path": "demo_projects/go",
          "children": [
            {
              "id": "demo_projects/go/project",
              "name": "project",
              "type": "folder",
              "path": "demo_projects/go/project",
              "children": [
                {
                  "id": "demo_projects/go/project/calculator",
                  "name": "calculator",
                  "type": "folder",
                  "path": "demo_projects/go/project/calculator",
                  "children": [
                    {
                      "id": "demo_projects/go/project/calculator/calculator.go",
                      "name": "calculator.go",
                      "type": "file",
                      "path": "demo_projects/go/project/calculator/calculator.go",
                      "metrics": {
                        "lineCoverage": {
                          "covered": 27,
                          "uncovered": 5,
                          "coverable": 32,
                          "total": 53,
                          "percentage": 84.37
                        },
                        "methodsCovered": { "covered": 5, "total": 5, "percentage": 100 },
                        "methodsFullyCovered": { "covered": 2, "total": 5, "percentage": 40 }
                      },
                      "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                      "targetUrl": "demo_projects_go_project_calculator_calculator.go.html"
                    },
                    {
                      "id": "demo_projects/go/project/calculator/entities.go",
                      "name": "entities.go",
                      "type": "file",
                      "path": "demo_projects/go/project/calculator/entities.go",
                      "metrics": {
                        "lineCoverage": {
                          "covered": 20,
                          "uncovered": 5,
                          "coverable": 25,
                          "total": 56,
                          "percentage": 80
                        },
                        "methodsCovered": { "covered": 5, "total": 6, "percentage": 83.33 },
                        "methodsFullyCovered": { "covered": 5, "total": 6, "percentage": 83.33 }
                      },
                      "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
                      "targetUrl": "demo_projects_go_project_calculator_entities.go.html"
                    }
                  ],
                  "metrics": {
                    "lineCoverage": {
                      "covered": 47,
                      "uncovered": 10,
                      "coverable": 57,
                      "total": 109,
                      "percentage": 82.45
                    },
                    "methodsCovered": { "covered": 10, "total": 11, "percentage": 90.9 },
                    "methodsFullyCovered": { "covered": 7, "total": 11, "percentage": 63.63 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "warning" }
                },
                {
                  "id": "demo_projects/go/project/calculator_2",
                  "name": "calculator_2",
                  "type": "folder",
                  "path": "demo_projects/go/project/calculator_2",
                  "children": [
                    {
                      "id": "demo_projects/go/project/calculator_2/calculator.go",
                      "name": "calculator.go",
                      "type": "file",
                      "path": "demo_projects/go/project/calculator_2/calculator.go",
                      "metrics": {
                        "lineCoverage": {
                          "covered": 23,
                          "uncovered": 9,
                          "coverable": 32,
                          "total": 53,
                          "percentage": 71.87
                        },
                        "methodsCovered": { "covered": 4, "total": 5, "percentage": 80 },
                        "methodsFullyCovered": { "covered": 2, "total": 5, "percentage": 40 }
                      },
                      "statuses": {
                        "lineCoverage": "warning",
                        "methodsCovered": "safe",
                        "methodsFullyCovered": "danger"
                      },
                      "targetUrl": "demo_projects_go_project_calculator_2_calculator.go.html"
                    },
                    {
                      "id": "demo_projects/go/project/calculator_2/entities.go",
                      "name": "entities.go",
                      "type": "file",
                      "path": "demo_projects/go/project/calculator_2/entities.go",
                      "metrics": {
                        "lineCoverage": {
                          "covered": 5,
                          "uncovered": 20,
                          "coverable": 25,
                          "total": 56,
                          "percentage": 20
                        },
                        "methodsCovered": { "covered": 1, "total": 6, "percentage": 16.66 },
                        "methodsFullyCovered": { "covered": 1, "total": 6, "percentage": 16.66 }
                      },
                      "statuses": {
                        "lineCoverage": "danger",
                        "methodsCovered": "danger",
                        "methodsFullyCovered": "danger"
                      },
                      "targetUrl": "demo_projects_go_project_calculator_2_entities.go.html"
                    }
                  ],
                  "metrics": {
                    "lineCoverage": {
                      "covered": 28,
                      "uncovered": 29,
                      "coverable": 57,
                      "total": 109,
                      "percentage": 49.12
                    },
                    "methodsCovered": { "covered": 5, "total": 11, "percentage": 45.45 },
                    "methodsFullyCovered": { "covered": 3, "total": 11, "percentage": 27.27 }
                  },
                  "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" }
                }
              ],
              "metrics": {
                "lineCoverage": { "covered": 75, "uncovered": 39, "coverable": 114, "total": 218, "percentage": 65.78 },
                "methodsCovered": { "covered": 15, "total": 22, "percentage": 68.18 },
                "methodsFullyCovered": { "covered": 10, "total": 22, "percentage": 45.45 }
              },
              "statuses": { "lineCoverage": "warning", "methodsCovered": "warning", "methodsFullyCovered": "danger" }
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 75, "uncovered": 39, "coverable": 114, "total": 218, "percentage": 65.78 },
            "methodsCovered": { "covered": 15, "total": 22, "percentage": 68.18 },
            "methodsFullyCovered": { "covered": 10, "total": 22, "percentage": 45.45 }
          },
          "statuses": { "lineCoverage": "warning", "methodsCovered": "warning", "methodsFullyCovered": "danger" }
        }
      ],
      "metrics": {
        "branchCoverage": { "covered": 39, "total": 58, "percentage": 67.24 },
        "lineCoverage": { "covered": 223, "uncovered": 90, "coverable": 313, "total": 873, "percentage": 71.24 },
        "methodsCovered": { "covered": 23, "total": 31, "percentage": 74.19 },
        "methodsFullyCovered": { "covered": 13, "total": 31, "percentage": 41.93 }
      },
      "statuses": {
        "branchCoverage": "warning",
        "lineCoverage": "warning",
        "methodsCovered": "warning",
        "methodsFullyCovered": "danger"
      }
    },
    {
      "id": "filereader",
      "name": "filereader",
      "type": "folder",
      "path": "filereader",
      "children": [
        {
          "id": "filereader/default_reader.go",
          "name": "default_reader.go",
          "type": "file",
          "path": "filereader/default_reader.go",
          "metrics": {
            "lineCoverage": { "covered": 15, "uncovered": 0, "coverable": 15, "total": 28, "percentage": 100 },
            "methodsCovered": { "covered": 5, "total": 5, "percentage": 100 },
            "methodsFullyCovered": { "covered": 5, "total": 5, "percentage": 100 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
          "targetUrl": "filereader_default_reader.go.html"
        },
        {
          "id": "filereader/filereader.go",
          "name": "filereader.go",
          "type": "file",
          "path": "filereader/filereader.go",
          "metrics": {
            "lineCoverage": { "covered": 47, "uncovered": 16, "coverable": 63, "total": 95, "percentage": 74.6 },
            "methodsCovered": { "covered": 3, "total": 3, "percentage": 100 },
            "methodsFullyCovered": { "covered": 0, "total": 3, "percentage": 0 }
          },
          "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
          "targetUrl": "filereader_filereader.go.html"
        }
      ],
      "metrics": {
        "lineCoverage": { "covered": 62, "uncovered": 16, "coverable": 78, "total": 123, "percentage": 79.48 },
        "methodsCovered": { "covered": 8, "total": 8, "percentage": 100 },
        "methodsFullyCovered": { "covered": 5, "total": 8, "percentage": 62.5 }
      },
      "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "warning" }
    },
    {
      "id": "filesystem",
      "name": "filesystem",
      "type": "folder",
      "path": "filesystem",
      "children": [
        {
          "id": "filesystem/filesystem.go",
          "name": "filesystem.go",
          "type": "file",
          "path": "filesystem/filesystem.go",
          "metrics": {
            "lineCoverage": { "covered": 0, "uncovered": 11, "coverable": 11, "total": 102, "percentage": 0 },
            "methodsCovered": { "covered": 0, "total": 9, "percentage": 0 },
            "methodsFullyCovered": { "covered": 0, "total": 9, "percentage": 0 }
          },
          "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
          "targetUrl": "filesystem_filesystem.go.html"
        }
      ],
      "metrics": {
        "lineCoverage": { "covered": 0, "uncovered": 11, "coverable": 11, "total": 102, "percentage": 0 },
        "methodsCovered": { "covered": 0, "total": 9, "percentage": 0 },
        "methodsFullyCovered": { "covered": 0, "total": 9, "percentage": 0 }
      },
      "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" }
    },
    {
      "id": "filtering",
      "name": "filtering",
      "type": "folder",
      "path": "filtering",
      "children": [
        {
          "id": "filtering/filter.go",
          "name": "filter.go",
          "type": "file",
          "path": "filtering/filter.go",
          "metrics": {
            "lineCoverage": { "covered": 75, "uncovered": 4, "coverable": 79, "total": 166, "percentage": 94.93 },
            "methodsCovered": { "covered": 4, "total": 4, "percentage": 100 },
            "methodsFullyCovered": { "covered": 2, "total": 4, "percentage": 50 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
          "targetUrl": "filtering_filter.go.html"
        }
      ],
      "metrics": {
        "lineCoverage": { "covered": 75, "uncovered": 4, "coverable": 79, "total": 166, "percentage": 94.93 },
        "methodsCovered": { "covered": 4, "total": 4, "percentage": 100 },
        "methodsFullyCovered": { "covered": 2, "total": 4, "percentage": 50 }
      },
      "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
    },
    {
      "id": "internal",
      "name": "internal",
      "type": "folder",
      "path": "internal",
      "children": [
        {
          "id": "internal/aggregator",
          "name": "aggregator",
          "type": "folder",
          "path": "internal/aggregator",
          "children": [
            {
              "id": "internal/aggregator/aggrgator.go",
              "name": "aggrgator.go",
              "type": "file",
              "path": "internal/aggregator/aggrgator.go",
              "metrics": {
                "lineCoverage": { "covered": 52, "uncovered": 0, "coverable": 52, "total": 77, "percentage": 100 },
                "methodsCovered": { "covered": 4, "total": 4, "percentage": 100 },
                "methodsFullyCovered": { "covered": 4, "total": 4, "percentage": 100 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
              "targetUrl": "internal_aggregator_aggrgator.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 52, "uncovered": 0, "coverable": 52, "total": 77, "percentage": 100 },
            "methodsCovered": { "covered": 4, "total": 4, "percentage": 100 },
            "methodsFullyCovered": { "covered": 4, "total": 4, "percentage": 100 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" }
        },
        {
          "id": "internal/config",
          "name": "config",
          "type": "folder",
          "path": "internal/config",
          "children": [
            {
              "id": "internal/config/config.go",
              "name": "config.go",
              "type": "file",
              "path": "internal/config/config.go",
              "metrics": {
                "lineCoverage": { "covered": 89, "uncovered": 29, "coverable": 118, "total": 195, "percentage": 75.42 },
                "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
                "methodsFullyCovered": { "covered": 2, "total": 6, "percentage": 33.33 }
              },
              "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_config_config.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 89, "uncovered": 29, "coverable": 118, "total": 195, "percentage": 75.42 },
            "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
            "methodsFullyCovered": { "covered": 2, "total": 6, "percentage": 33.33 }
          },
          "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
        },
        {
          "id": "internal/enricher",
          "name": "enricher",
          "type": "folder",
          "path": "internal/enricher",
          "children": [
            {
              "id": "internal/enricher/enricher.go",
              "name": "enricher.go",
              "type": "file",
              "path": "internal/enricher/enricher.go",
              "metrics": {
                "lineCoverage": { "covered": 91, "uncovered": 13, "coverable": 104, "total": 195, "percentage": 87.5 },
                "methodsCovered": { "covered": 8, "total": 8, "percentage": 100 },
                "methodsFullyCovered": { "covered": 6, "total": 8, "percentage": 75 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "warning" },
              "targetUrl": "internal_enricher_enricher.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 91, "uncovered": 13, "coverable": 104, "total": 195, "percentage": 87.5 },
            "methodsCovered": { "covered": 8, "total": 8, "percentage": 100 },
            "methodsFullyCovered": { "covered": 6, "total": 8, "percentage": 75 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "warning" }
        },
        {
          "id": "internal/parsers",
          "name": "parsers",
          "type": "folder",
          "path": "internal/parsers",
          "children": [
            {
              "id": "internal/parsers/parser_cobertura",
              "name": "parser_cobertura",
              "type": "folder",
              "path": "internal/parsers/parser_cobertura",
              "children": [
                {
                  "id": "internal/parsers/parser_cobertura/parser.go",
                  "name": "parser.go",
                  "type": "file",
                  "path": "internal/parsers/parser_cobertura/parser.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 64,
                      "uncovered": 21,
                      "coverable": 85,
                      "total": 133,
                      "percentage": 75.29
                    },
                    "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 2, "total": 6, "percentage": 33.33 }
                  },
                  "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_parsers_parser_cobertura_parser.go.html"
                },
                {
                  "id": "internal/parsers/parser_cobertura/processing.go",
                  "name": "processing.go",
                  "type": "file",
                  "path": "internal/parsers/parser_cobertura/processing.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 64,
                      "uncovered": 13,
                      "coverable": 77,
                      "total": 141,
                      "percentage": 83.11
                    },
                    "methodsCovered": { "covered": 4, "total": 4, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 1, "total": 4, "percentage": 25 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_parsers_parser_cobertura_processing.go.html"
                }
              ],
              "metrics": {
                "lineCoverage": {
                  "covered": 128,
                  "uncovered": 34,
                  "coverable": 162,
                  "total": 274,
                  "percentage": 79.01
                },
                "methodsCovered": { "covered": 10, "total": 10, "percentage": 100 },
                "methodsFullyCovered": { "covered": 3, "total": 10, "percentage": 30 }
              },
              "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
            },
            {
              "id": "internal/parsers/parser_gcov",
              "name": "parser_gcov",
              "type": "folder",
              "path": "internal/parsers/parser_gcov",
              "children": [
                {
                  "id": "internal/parsers/parser_gcov/parser.go",
                  "name": "parser.go",
                  "type": "file",
                  "path": "internal/parsers/parser_gcov/parser.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 39,
                      "uncovered": 6,
                      "coverable": 45,
                      "total": 79,
                      "percentage": 86.66
                    },
                    "methodsCovered": { "covered": 4, "total": 4, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 2, "total": 4, "percentage": 50 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_parsers_parser_gcov_parser.go.html"
                },
                {
                  "id": "internal/parsers/parser_gcov/processing.go",
                  "name": "processing.go",
                  "type": "file",
                  "path": "internal/parsers/parser_gcov/processing.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 54,
                      "uncovered": 2,
                      "coverable": 56,
                      "total": 93,
                      "percentage": 96.42
                    },
                    "methodsCovered": { "covered": 2, "total": 2, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 1, "total": 2, "percentage": 50 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_parsers_parser_gcov_processing.go.html"
                }
              ],
              "metrics": {
                "lineCoverage": { "covered": 93, "uncovered": 8, "coverable": 101, "total": 172, "percentage": 92.07 },
                "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
                "methodsFullyCovered": { "covered": 3, "total": 6, "percentage": 50 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
            },
            {
              "id": "internal/parsers/parser_gocover",
              "name": "parser_gocover",
              "type": "folder",
              "path": "internal/parsers/parser_gocover",
              "children": [
                {
                  "id": "internal/parsers/parser_gocover/parser.go",
                  "name": "parser.go",
                  "type": "file",
                  "path": "internal/parsers/parser_gocover/parser.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 66,
                      "uncovered": 12,
                      "coverable": 78,
                      "total": 124,
                      "percentage": 84.61
                    },
                    "methodsCovered": { "covered": 5, "total": 5, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 2, "total": 5, "percentage": 40 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_parsers_parser_gocover_parser.go.html"
                },
                {
                  "id": "internal/parsers/parser_gocover/processing.go",
                  "name": "processing.go",
                  "type": "file",
                  "path": "internal/parsers/parser_gocover/processing.go",
                  "metrics": {
                    "lineCoverage": { "covered": 61, "uncovered": 0, "coverable": 61, "total": 101, "percentage": 100 },
                    "methodsCovered": { "covered": 4, "total": 4, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 4, "total": 4, "percentage": 100 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
                  "targetUrl": "internal_parsers_parser_gocover_processing.go.html"
                }
              ],
              "metrics": {
                "lineCoverage": {
                  "covered": 127,
                  "uncovered": 12,
                  "coverable": 139,
                  "total": 225,
                  "percentage": 91.36
                },
                "methodsCovered": { "covered": 9, "total": 9, "percentage": 100 },
                "methodsFullyCovered": { "covered": 6, "total": 9, "percentage": 66.66 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "warning" }
            },
            {
              "id": "internal/parsers/factory.go",
              "name": "factory.go",
              "type": "file",
              "path": "internal/parsers/factory.go",
              "metrics": {
                "lineCoverage": { "covered": 14, "uncovered": 0, "coverable": 14, "total": 25, "percentage": 100 },
                "methodsCovered": { "covered": 2, "total": 2, "percentage": 100 },
                "methodsFullyCovered": { "covered": 2, "total": 2, "percentage": 100 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
              "targetUrl": "internal_parsers_factory.go.html"
            },
            {
              "id": "internal/parsers/parser_config.go",
              "name": "parser_config.go",
              "type": "file",
              "path": "internal/parsers/parser_config.go",
              "metrics": {
                "lineCoverage": { "covered": 2, "uncovered": 1, "coverable": 3, "total": 44, "percentage": 66.66 },
                "methodsCovered": { "covered": 2, "total": 3, "percentage": 66.66 },
                "methodsFullyCovered": { "covered": 2, "total": 3, "percentage": 66.66 }
              },
              "statuses": { "lineCoverage": "warning", "methodsCovered": "warning", "methodsFullyCovered": "warning" },
              "targetUrl": "internal_parsers_parser_config.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 364, "uncovered": 55, "coverable": 419, "total": 740, "percentage": 86.87 },
            "methodsCovered": { "covered": 29, "total": 30, "percentage": 96.66 },
            "methodsFullyCovered": { "covered": 16, "total": 30, "percentage": 53.33 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
        },
        {
          "id": "internal/reporter",
          "name": "reporter",
          "type": "folder",
          "path": "internal/reporter",
          "children": [
            {
              "id": "internal/reporter/htmlreact",
              "name": "htmlreact",
              "type": "folder",
              "path": "internal/reporter/htmlreact",
              "children": [
                {
                  "id": "internal/reporter/htmlreact/builder.go",
                  "name": "builder.go",
                  "type": "file",
                  "path": "internal/reporter/htmlreact/builder.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 234,
                      "uncovered": 28,
                      "coverable": 262,
                      "total": 326,
                      "percentage": 89.31
                    },
                    "methodsCovered": { "covered": 11, "total": 12, "percentage": 91.66 },
                    "methodsFullyCovered": { "covered": 7, "total": 12, "percentage": 58.33 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_reporter_htmlreact_builder.go.html"
                },
                {
                  "id": "internal/reporter/htmlreact/details_generator.go",
                  "name": "details_generator.go",
                  "type": "file",
                  "path": "internal/reporter/htmlreact/details_generator.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 170,
                      "uncovered": 27,
                      "coverable": 197,
                      "total": 258,
                      "percentage": 86.29
                    },
                    "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 1, "total": 6, "percentage": 16.66 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_reporter_htmlreact_details_generator.go.html"
                },
                {
                  "id": "internal/reporter/htmlreact/embed.go",
                  "name": "embed.go",
                  "type": "file",
                  "path": "internal/reporter/htmlreact/embed.go",
                  "metrics": {
                    "lineCoverage": { "covered": 3, "uncovered": 0, "coverable": 3, "total": 16, "percentage": 100 },
                    "methodsCovered": { "covered": 1, "total": 1, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 1, "total": 1, "percentage": 100 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
                  "targetUrl": "internal_reporter_htmlreact_embed.go.html"
                },
                {
                  "id": "internal/reporter/htmlreact/emit.go",
                  "name": "emit.go",
                  "type": "file",
                  "path": "internal/reporter/htmlreact/emit.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 15,
                      "uncovered": 7,
                      "coverable": 22,
                      "total": 34,
                      "percentage": 68.18
                    },
                    "methodsCovered": { "covered": 1, "total": 1, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 0, "total": 1, "percentage": 0 }
                  },
                  "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_reporter_htmlreact_emit.go.html"
                },
                {
                  "id": "internal/reporter/htmlreact/generator.go",
                  "name": "generator.go",
                  "type": "file",
                  "path": "internal/reporter/htmlreact/generator.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 32,
                      "uncovered": 29,
                      "coverable": 61,
                      "total": 90,
                      "percentage": 52.45
                    },
                    "methodsCovered": { "covered": 2, "total": 2, "percentage": 100 },
                    "methodsFullyCovered": { "covered": 0, "total": 2, "percentage": 0 }
                  },
                  "statuses": { "lineCoverage": "danger", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_reporter_htmlreact_generator.go.html"
                }
              ],
              "metrics": {
                "lineCoverage": { "covered": 454, "uncovered": 91, "coverable": 545, "total": 724, "percentage": 83.3 },
                "methodsCovered": { "covered": 21, "total": 22, "percentage": 95.45 },
                "methodsFullyCovered": { "covered": 9, "total": 22, "percentage": 40.9 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
            },
            {
              "id": "internal/reporter/lcov",
              "name": "lcov",
              "type": "folder",
              "path": "internal/reporter/lcov",
              "children": [
                {
                  "id": "internal/reporter/lcov/reporter.go",
                  "name": "reporter.go",
                  "type": "file",
                  "path": "internal/reporter/lcov/reporter.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 72,
                      "uncovered": 31,
                      "coverable": 103,
                      "total": 143,
                      "percentage": 69.9
                    },
                    "methodsCovered": { "covered": 4, "total": 5, "percentage": 80 },
                    "methodsFullyCovered": { "covered": 2, "total": 5, "percentage": 40 }
                  },
                  "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_reporter_lcov_reporter.go.html"
                }
              ],
              "metrics": {
                "lineCoverage": { "covered": 72, "uncovered": 31, "coverable": 103, "total": 143, "percentage": 69.9 },
                "methodsCovered": { "covered": 4, "total": 5, "percentage": 80 },
                "methodsFullyCovered": { "covered": 2, "total": 5, "percentage": 40 }
              },
              "statuses": { "lineCoverage": "warning", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
            },
            {
              "id": "internal/reporter/reporter_rawjson",
              "name": "reporter_rawjson",
              "type": "folder",
              "path": "internal/reporter/reporter_rawjson",
              "children": [
                {
                  "id": "internal/reporter/reporter_rawjson/reporter.go",
                  "name": "reporter.go",
                  "type": "file",
                  "path": "internal/reporter/reporter_rawjson/reporter.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 15,
                      "uncovered": 7,
                      "coverable": 22,
                      "total": 44,
                      "percentage": 68.18
                    },
                    "methodsCovered": { "covered": 2, "total": 3, "percentage": 66.66 },
                    "methodsFullyCovered": { "covered": 1, "total": 3, "percentage": 33.33 }
                  },
                  "statuses": {
                    "lineCoverage": "warning",
                    "methodsCovered": "warning",
                    "methodsFullyCovered": "danger"
                  },
                  "targetUrl": "internal_reporter_reporter_rawjson_reporter.go.html"
                }
              ],
              "metrics": {
                "lineCoverage": { "covered": 15, "uncovered": 7, "coverable": 22, "total": 44, "percentage": 68.18 },
                "methodsCovered": { "covered": 2, "total": 3, "percentage": 66.66 },
                "methodsFullyCovered": { "covered": 1, "total": 3, "percentage": 33.33 }
              },
              "statuses": { "lineCoverage": "warning", "methodsCovered": "warning", "methodsFullyCovered": "danger" }
            },
            {
              "id": "internal/reporter/textsummary",
              "name": "textsummary",
              "type": "folder",
              "path": "internal/reporter/textsummary",
              "children": [
                {
                  "id": "internal/reporter/textsummary/reporter.go",
                  "name": "reporter.go",
                  "type": "file",
                  "path": "internal/reporter/textsummary/reporter.go",
                  "metrics": {
                    "lineCoverage": {
                      "covered": 66,
                      "uncovered": 7,
                      "coverable": 73,
                      "total": 111,
                      "percentage": 90.41
                    },
                    "methodsCovered": { "covered": 3, "total": 4, "percentage": 75 },
                    "methodsFullyCovered": { "covered": 2, "total": 4, "percentage": 50 }
                  },
                  "statuses": { "lineCoverage": "safe", "methodsCovered": "warning", "methodsFullyCovered": "danger" },
                  "targetUrl": "internal_reporter_textsummary_reporter.go.html"
                }
              ],
              "metrics": {
                "lineCoverage": { "covered": 66, "uncovered": 7, "coverable": 73, "total": 111, "percentage": 90.41 },
                "methodsCovered": { "covered": 3, "total": 4, "percentage": 75 },
                "methodsFullyCovered": { "covered": 2, "total": 4, "percentage": 50 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "warning", "methodsFullyCovered": "danger" }
            },
            {
              "id": "internal/reporter/context.go",
              "name": "context.go",
              "type": "file",
              "path": "internal/reporter/context.go",
              "metrics": {
                "lineCoverage": { "covered": 0, "uncovered": 10, "coverable": 10, "total": 35, "percentage": 0 },
                "methodsCovered": { "covered": 0, "total": 3, "percentage": 0 },
                "methodsFullyCovered": { "covered": 0, "total": 3, "percentage": 0 }
              },
              "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_reporter_context.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 607, "uncovered": 146, "coverable": 753, "total": 1057, "percentage": 80.61 },
            "methodsCovered": { "covered": 30, "total": 37, "percentage": 81.08 },
            "methodsFullyCovered": { "covered": 14, "total": 37, "percentage": 37.83 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" }
        },
        {
          "id": "internal/testutil",
          "name": "testutil",
          "type": "folder",
          "path": "internal/testutil",
          "children": [
            {
              "id": "internal/testutil/mock_filesystem.go",
              "name": "mock_filesystem.go",
              "type": "file",
              "path": "internal/testutil/mock_filesystem.go",
              "metrics": {
                "lineCoverage": { "covered": 0, "uncovered": 137, "coverable": 137, "total": 209, "percentage": 0 },
                "methodsCovered": { "covered": 0, "total": 25, "percentage": 0 },
                "methodsFullyCovered": { "covered": 0, "total": 25, "percentage": 0 }
              },
              "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_testutil_mock_filesystem.go.html"
            },
            {
              "id": "internal/testutil/parser_helper.go",
              "name": "parser_helper.go",
              "type": "file",
              "path": "internal/testutil/parser_helper.go",
              "metrics": {
                "lineCoverage": { "covered": 0, "uncovered": 12, "coverable": 12, "total": 31, "percentage": 0 },
                "methodsCovered": { "covered": 0, "total": 4, "percentage": 0 },
                "methodsFullyCovered": { "covered": 0, "total": 4, "percentage": 0 }
              },
              "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_testutil_parser_helper.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 0, "uncovered": 149, "coverable": 149, "total": 240, "percentage": 0 },
            "methodsCovered": { "covered": 0, "total": 29, "percentage": 0 },
            "methodsFullyCovered": { "covered": 0, "total": 29, "percentage": 0 }
          },
          "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" }
        },
        {
          "id": "internal/tree",
          "name": "tree",
          "type": "folder",
          "path": "internal/tree",
          "children": [
            {
              "id": "internal/tree/builder.go",
              "name": "builder.go",
              "type": "file",
              "path": "internal/tree/builder.go",
              "metrics": {
                "lineCoverage": { "covered": 122, "uncovered": 7, "coverable": 129, "total": 186, "percentage": 94.57 },
                "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
                "methodsFullyCovered": { "covered": 5, "total": 6, "percentage": 83.33 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
              "targetUrl": "internal_tree_builder.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 122, "uncovered": 7, "coverable": 129, "total": 186, "percentage": 94.57 },
            "methodsCovered": { "covered": 6, "total": 6, "percentage": 100 },
            "methodsFullyCovered": { "covered": 5, "total": 6, "percentage": 83.33 }
          },
          "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" }
        },
        {
          "id": "internal/utils",
          "name": "utils",
          "type": "folder",
          "path": "internal/utils",
          "children": [
            {
              "id": "internal/utils/analyzer.go",
              "name": "analyzer.go",
              "type": "file",
              "path": "internal/utils/analyzer.go",
              "metrics": {
                "lineCoverage": { "covered": 5, "uncovered": 12, "coverable": 17, "total": 33, "percentage": 29.41 },
                "methodsCovered": { "covered": 1, "total": 3, "percentage": 33.33 },
                "methodsFullyCovered": { "covered": 1, "total": 3, "percentage": 33.33 }
              },
              "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_utils_analyzer.go.html"
            },
            {
              "id": "internal/utils/brace_finder.go",
              "name": "brace_finder.go",
              "type": "file",
              "path": "internal/utils/brace_finder.go",
              "metrics": {
                "lineCoverage": { "covered": 45, "uncovered": 0, "coverable": 45, "total": 66, "percentage": 100 },
                "methodsCovered": { "covered": 1, "total": 1, "percentage": 100 },
                "methodsFullyCovered": { "covered": 1, "total": 1, "percentage": 100 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
              "targetUrl": "internal_utils_brace_finder.go.html"
            },
            {
              "id": "internal/utils/collections.go",
              "name": "collections.go",
              "type": "file",
              "path": "internal/utils/collections.go",
              "metrics": {
                "lineCoverage": { "covered": 0, "uncovered": 13, "coverable": 13, "total": 20, "percentage": 0 },
                "methodsCovered": { "covered": 0, "total": 1, "percentage": 0 },
                "methodsFullyCovered": { "covered": 0, "total": 1, "percentage": 0 }
              },
              "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_utils_collections.go.html"
            },
            {
              "id": "internal/utils/encoding.go",
              "name": "encoding.go",
              "type": "file",
              "path": "internal/utils/encoding.go",
              "metrics": {
                "lineCoverage": { "covered": 0, "uncovered": 26, "coverable": 26, "total": 44, "percentage": 0 },
                "methodsCovered": { "covered": 0, "total": 1, "percentage": 0 },
                "methodsFullyCovered": { "covered": 0, "total": 1, "percentage": 0 }
              },
              "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_utils_encoding.go.html"
            },
            {
              "id": "internal/utils/line_sorter.go",
              "name": "line_sorter.go",
              "type": "file",
              "path": "internal/utils/line_sorter.go",
              "metrics": {
                "lineCoverage": { "covered": 19, "uncovered": 0, "coverable": 19, "total": 40, "percentage": 100 },
                "methodsCovered": { "covered": 1, "total": 1, "percentage": 100 },
                "methodsFullyCovered": { "covered": 1, "total": 1, "percentage": 100 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "safe" },
              "targetUrl": "internal_utils_line_sorter.go.html"
            },
            {
              "id": "internal/utils/math.go",
              "name": "math.go",
              "type": "file",
              "path": "internal/utils/math.go",
              "metrics": {
                "lineCoverage": { "covered": 37, "uncovered": 9, "coverable": 46, "total": 65, "percentage": 80.43 },
                "methodsCovered": { "covered": 2, "total": 2, "percentage": 100 },
                "methodsFullyCovered": { "covered": 0, "total": 2, "percentage": 0 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_utils_math.go.html"
            },
            {
              "id": "internal/utils/paths.go",
              "name": "paths.go",
              "type": "file",
              "path": "internal/utils/paths.go",
              "metrics": {
                "lineCoverage": { "covered": 67, "uncovered": 7, "coverable": 74, "total": 115, "percentage": 90.54 },
                "methodsCovered": { "covered": 2, "total": 2, "percentage": 100 },
                "methodsFullyCovered": { "covered": 0, "total": 2, "percentage": 0 }
              },
              "statuses": { "lineCoverage": "safe", "methodsCovered": "safe", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_utils_paths.go.html"
            },
            {
              "id": "internal/utils/stringutils.go",
              "name": "stringutils.go",
              "type": "file",
              "path": "internal/utils/stringutils.go",
              "metrics": {
                "lineCoverage": { "covered": 0, "uncovered": 85, "coverable": 85, "total": 137, "percentage": 0 },
                "methodsCovered": { "covered": 0, "total": 4, "percentage": 0 },
                "methodsFullyCovered": { "covered": 0, "total": 4, "percentage": 0 }
              },
              "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" },
              "targetUrl": "internal_utils_stringutils.go.html"
            }
          ],
          "metrics": {
            "lineCoverage": { "covered": 173, "uncovered": 152, "coverable": 325, "total": 520, "percentage": 53.23 },
            "methodsCovered": { "covered": 7, "total": 15, "percentage": 46.66 },
            "methodsFullyCovered": { "covered": 3, "total": 15, "percentage": 20 }
          },
          "statuses": { "lineCoverage": "danger", "methodsCovered": "danger", "methodsFullyCovered": "danger" }
        }
      ],
      "metrics": {
        "lineCoverage": { "covered": 1498, "uncovered": 551, "coverable": 2049, "total": 3210, "percentage": 73.1 },
        "methodsCovered": { "covered": 90, "total": 135, "percentage": 66.66 },
        "methodsFullyCovered": { "covered": 50, "total": 135, "percentage": 37.03 }
      },
      "statuses": { "lineCoverage": "warning", "methodsCovered": "warning", "methodsFullyCovered": "danger" }
    },
    {
      "id": "logging",
      "name": "logging",
      "type": "folder",
      "path": "logging",
      "children": [
        {
          "id": "logging/logging.go",
          "name": "logging.go",
          "type": "file",
          "path": "logging/logging.go",
          "metrics": {
            "lineCoverage": { "covered": 77, "uncovered": 20, "coverable": 97, "total": 178, "percentage": 79.38 },
            "methodsCovered": { "covered": 7, "total": 10, "percentage": 70 },
            "methodsFullyCovered": { "covered": 5, "total": 10, "percentage": 50 }
          },
          "statuses": { "lineCoverage": "warning", "methodsCovered": "warning", "methodsFullyCovered": "danger" },
          "targetUrl": "logging_logging.go.html"
        }
      ],
      "metrics": {
        "lineCoverage": { "covered": 77, "uncovered": 20, "coverable": 97, "total": 178, "percentage": 79.38 },
        "methodsCovered": { "covered": 7, "total": 10, "percentage": 70 },
        "methodsFullyCovered": { "covered": 5, "total": 10, "percentage": 50 }
      },
      "statuses": { "lineCoverage": "warning", "methodsCovered": "warning", "methodsFullyCovered": "danger" }
    }
  ],
  "metricDefinitions": {
    "branchCoverage": {
      "label": "Branches",
      "shortLabel": "Branches",
      "subMetrics": [
        { "id": "covered", "label": "Covered", "width": 100 },
        { "id": "total", "label": "Total", "width": 80 },
        { "id": "percentage", "label": "Percentage %", "width": 160 }
      ]
    },
    "lineCoverage": {
      "label": "Lines",
      "shortLabel": "Lines",
      "subMetrics": [
        { "id": "covered", "label": "Covered", "width": 100 },
        { "id": "uncovered", "label": "Uncovered", "width": 100 },
        { "id": "coverable", "label": "Coverable", "width": 100 },
        { "id": "total", "label": "Total", "width": 80 },
        { "id": "percentage", "label": "Percentage %", "width": 160 }
      ]
    },
    "maxCyclomaticComplexity": {
      "label": "Max Cyclomatic Complexity",
      "shortLabel": "Max Complexity",
      "subMetrics": [{ "id": "total", "label": "Value", "width": 100 }]
    },
    "methodBranchCoverage": {
      "label": "Method Branches",
      "shortLabel": "Method Branches",
      "subMetrics": [
        { "id": "covered", "label": "Covered", "width": 100 },
        { "id": "total", "label": "Total", "width": 80 },
        { "id": "percentage", "label": "Percentage %", "width": 160 }
      ]
    },
    "methodsCovered": {
      "label": "Methods Covered",
      "shortLabel": "Methods Cov.",
      "subMetrics": [
        { "id": "covered", "label": "Covered", "width": 80 },
        { "id": "total", "label": "Total", "width": 80 },
        { "id": "percentage", "label": "Percentage %", "width": 160 }
      ]
    },
    "methodsFullyCovered": {
      "label": "Methods Fully Covered",
      "shortLabel": "Methods Full Cov.",
      "subMetrics": [
        { "id": "covered", "label": "Covered", "width": 80 },
        { "id": "total", "label": "Total", "width": 80 },
        { "id": "percentage", "label": "Percentage %", "width": 160 }
      ]
    }
  },
  "metadata": [
    { "label": "Generated At", "value": "2025-10-28 18:39:29" },
    { "label": "Parser", "value": "Cobertura | GCov | GoCover" }
  ]
}
