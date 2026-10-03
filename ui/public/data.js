/**
 * NOTE: This is a developer data fixture for testing purposes, loaded only by
 * the Vite dev server (pnpm dev). In production, the Go reporter writes a
 * real data.js next to index.html.
 *
 * Generated from the nanovision self-coverage reports (full merged) with a
 * git diff applied, plus the `review` block (gate verdict,
 * changelist stats, risk hotspots) so the review header renders in dev.
 *
 * Regenerate with: python scripts/gen_ui_fixture.py
 * Schema version: 1 (flat node list)
 */
window.__NANOVISION_SUMMARY__ = {
 "schemaVersion": 1,
 "generatedAt": "2026-10-03T16:49:35Z",
 "title": "nanovision Self-Coverage (dev fixture)",
 "totals": {
  "files": 117,
  "folders": 56,
  "max_cyclomatic_complexity": {
   "value": 31
  },
  "methods_fully_covered": {
   "covered": 291,
   "total": 633,
   "percentage": 45.97
  },
  "methods_hit": {
   "covered": 563,
   "total": 633,
   "percentage": 88.94
  },
  "patch_methods_hit": {
   "covered": 328,
   "total": 353,
   "percentage": 92.91
  },
  "patch_statement_coverage": {
   "covered": 2760,
   "uncovered": 701,
   "coverable": 3461,
   "total": 3461,
   "percentage": 79.74
  },
  "patch_statement_methods_hit": {
   "covered": 331,
   "total": 353,
   "percentage": 93.76
  },
  "statement_coverage": {
   "covered": 5272,
   "uncovered": 1367,
   "coverable": 6639,
   "total": 6639,
   "percentage": 79.4
  },
  "statement_methods_fully_covered": {
   "covered": 291,
   "total": 633,
   "percentage": 45.97
  },
  "statement_methods_hit": {
   "covered": 563,
   "total": 633,
   "percentage": 88.94
  },
  "statuses": {
   "patch_methods_hit": "safe",
   "patch_statement_coverage": "warning",
   "statement_coverage": "safe"
  }
 },
 "nodes": [
  {
   "id": "cmd",
   "name": "cmd",
   "type": "folder",
   "path": "cmd",
   "depth": 0,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 27
    },
    "methods_fully_covered": {
     "covered": 10,
     "total": 35,
     "percentage": 28.57
    },
    "methods_hit": {
     "covered": 27,
     "total": 35,
     "percentage": 77.14
    },
    "patch_methods_hit": {
     "covered": 21,
     "total": 24,
     "percentage": 87.5
    },
    "patch_statement_coverage": {
     "covered": 202,
     "uncovered": 180,
     "coverable": 382,
     "total": 382,
     "percentage": 52.87
    },
    "patch_statement_methods_hit": {
     "covered": 21,
     "total": 24,
     "percentage": 87.5
    },
    "statement_coverage": {
     "covered": 292,
     "uncovered": 309,
     "coverable": 601,
     "total": 601,
     "percentage": 48.58
    },
    "statement_methods_fully_covered": {
     "covered": 10,
     "total": 35,
     "percentage": 28.57
    },
    "statement_methods_hit": {
     "covered": 27,
     "total": 35,
     "percentage": 77.14
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "cmd/configcmd.go",
   "name": "configcmd.go",
   "type": "file",
   "path": "cmd/configcmd.go",
   "parentId": "cmd",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 15
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 7,
     "percentage": 28.57
    },
    "methods_hit": {
     "covered": 2,
     "total": 7,
     "percentage": 28.57
    },
    "statement_coverage": {
     "covered": 28,
     "uncovered": 83,
     "coverable": 111,
     "total": 111,
     "percentage": 25.22
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 7,
     "percentage": 28.57
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 7,
     "percentage": 28.57
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "cmd_configcmd.go.html"
  },
  {
   "id": "cmd/history.go",
   "name": "history.go",
   "type": "file",
   "path": "cmd/history.go",
   "parentId": "cmd",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 21
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 12,
     "percentage": 8.33
    },
    "methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 118,
     "uncovered": 38,
     "coverable": 156,
     "total": 156,
     "percentage": 75.64
    },
    "patch_statement_methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 118,
     "uncovered": 38,
     "coverable": 156,
     "total": 156,
     "percentage": 75.64
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 12,
     "percentage": 8.33
    },
    "statement_methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "cmd_history.go.html",
   "diffStatus": "added"
  },
  {
   "id": "cmd/main.go",
   "name": "main.go",
   "type": "file",
   "path": "cmd/main.go",
   "parentId": "cmd",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 27
    },
    "methods_fully_covered": {
     "covered": 7,
     "total": 11,
     "percentage": 63.63
    },
    "methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 60,
     "uncovered": 16,
     "coverable": 76,
     "total": 76,
     "percentage": 78.94
    },
    "patch_statement_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 122,
     "uncovered": 62,
     "coverable": 184,
     "total": 184,
     "percentage": 66.3
    },
    "statement_methods_fully_covered": {
     "covered": 7,
     "total": 11,
     "percentage": 63.63
    },
    "statement_methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "cmd_main.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "cmd/serve.go",
   "name": "serve.go",
   "type": "file",
   "path": "cmd/serve.go",
   "parentId": "cmd",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 16
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 90,
     "coverable": 90,
     "total": 90,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 0,
     "uncovered": 90,
     "coverable": 90,
     "total": 90,
     "percentage": 0
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "cmd_serve.go.html",
   "diffStatus": "added"
  },
  {
   "id": "cmd/storecmd.go",
   "name": "storecmd.go",
   "type": "file",
   "path": "cmd/storecmd.go",
   "parentId": "cmd",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 24,
     "uncovered": 36,
     "coverable": 60,
     "total": 60,
     "percentage": 40
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 24,
     "uncovered": 36,
     "coverable": 60,
     "total": 60,
     "percentage": 40
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "cmd_storecmd.go.html",
   "diffStatus": "added"
  },
  {
   "id": "demo_projects",
   "name": "demo_projects",
   "type": "folder",
   "path": "demo_projects",
   "depth": 0,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 14,
     "total": 33,
     "percentage": 42.42
    },
    "methods_hit": {
     "covered": 24,
     "total": 33,
     "percentage": 72.72
    },
    "statement_coverage": {
     "covered": 76,
     "uncovered": 32,
     "coverable": 108,
     "total": 108,
     "percentage": 70.37
    },
    "statement_methods_fully_covered": {
     "covered": 14,
     "total": 33,
     "percentage": 42.42
    },
    "statement_methods_hit": {
     "covered": 24,
     "total": 33,
     "percentage": 72.72
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   }
  },
  {
   "id": "demo_projects/cpp",
   "name": "cpp",
   "type": "folder",
   "path": "demo_projects/cpp",
   "parentId": "demo_projects",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 8,
     "total": 9,
     "percentage": 88.88
    },
    "statement_coverage": {
     "covered": 34,
     "uncovered": 9,
     "coverable": 43,
     "total": 43,
     "percentage": 79.06
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 9,
     "percentage": 88.88
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "demo_projects/cpp/project",
   "name": "project",
   "type": "folder",
   "path": "demo_projects/cpp/project",
   "parentId": "demo_projects/cpp",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 8,
     "total": 9,
     "percentage": 88.88
    },
    "statement_coverage": {
     "covered": 34,
     "uncovered": 9,
     "coverable": 43,
     "total": 43,
     "percentage": 79.06
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 9,
     "percentage": 88.88
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "demo_projects/cpp/project/src",
   "name": "src",
   "type": "folder",
   "path": "demo_projects/cpp/project/src",
   "parentId": "demo_projects/cpp/project",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 8,
     "total": 9,
     "percentage": 88.88
    },
    "statement_coverage": {
     "covered": 34,
     "uncovered": 9,
     "coverable": 43,
     "total": 43,
     "percentage": 79.06
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 9,
     "percentage": 88.88
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "demo_projects/cpp/project/src/utils",
   "name": "utils",
   "type": "folder",
   "path": "demo_projects/cpp/project/src/utils",
   "parentId": "demo_projects/cpp/project/src",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 13,
     "uncovered": 3,
     "coverable": 16,
     "total": 16,
     "percentage": 81.25
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "demo_projects/cpp/project/src/utils/math_utils.cpp",
   "name": "math_utils.cpp",
   "type": "file",
   "path": "demo_projects/cpp/project/src/utils/math_utils.cpp",
   "parentId": "demo_projects/cpp/project/src/utils",
   "depth": 5,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 13,
     "uncovered": 3,
     "coverable": 16,
     "total": 16,
     "percentage": 81.25
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "demo_projects_cpp_project_src_utils_math_utils.cpp.html"
  },
  {
   "id": "demo_projects/cpp/project/src/advanced_calculator.cpp",
   "name": "advanced_calculator.cpp",
   "type": "file",
   "path": "demo_projects/cpp/project/src/advanced_calculator.cpp",
   "parentId": "demo_projects/cpp/project/src",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 12,
     "uncovered": 4,
     "coverable": 16,
     "total": 16,
     "percentage": 75
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "demo_projects_cpp_project_src_advanced_calculator.cpp.html"
  },
  {
   "id": "demo_projects/cpp/project/src/calculator.cpp",
   "name": "calculator.cpp",
   "type": "file",
   "path": "demo_projects/cpp/project/src/calculator.cpp",
   "parentId": "demo_projects/cpp/project/src",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 3
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "statement_coverage": {
     "covered": 9,
     "uncovered": 2,
     "coverable": 11,
     "total": 11,
     "percentage": 81.81
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "demo_projects_cpp_project_src_calculator.cpp.html"
  },
  {
   "id": "demo_projects/cpp/nanovision.yaml",
   "name": "nanovision.yaml",
   "type": "file",
   "path": "demo_projects/cpp/nanovision.yaml",
   "parentId": "demo_projects/cpp",
   "depth": 2,
   "config": true
  },
  {
   "id": "demo_projects/csharp",
   "name": "csharp",
   "type": "folder",
   "path": "demo_projects/csharp",
   "parentId": "demo_projects",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   }
  },
  {
   "id": "demo_projects/csharp/project",
   "name": "project",
   "type": "folder",
   "path": "demo_projects/csharp/project",
   "parentId": "demo_projects/csharp",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   }
  },
  {
   "id": "demo_projects/csharp/project/Test",
   "name": "Test",
   "type": "folder",
   "path": "demo_projects/csharp/project/Test",
   "parentId": "demo_projects/csharp/project",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   }
  },
  {
   "id": "demo_projects/csharp/project/Test/AbstractClass.cs",
   "name": "AbstractClass.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/AbstractClass.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_AbstractClass.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/AnalyzerTestClass.cs",
   "name": "AnalyzerTestClass.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/AnalyzerTestClass.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_AnalyzerTestClass.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/AsyncClass.cs",
   "name": "AsyncClass.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/AsyncClass.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_AsyncClass.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/ClassWithExcludes.cs",
   "name": "ClassWithExcludes.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/ClassWithExcludes.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_ClassWithExcludes.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/ClassWithLocalFunctions.cs",
   "name": "ClassWithLocalFunctions.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/ClassWithLocalFunctions.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_ClassWithLocalFunctions.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/CodeContract_Contract.cs",
   "name": "CodeContract_Contract.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/CodeContract_Contract.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_CodeContract_Contract.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/CodeContract_Target.cs",
   "name": "CodeContract_Target.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/CodeContract_Target.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_CodeContract_Target.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/GenericAsyncClass.cs",
   "name": "GenericAsyncClass.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/GenericAsyncClass.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_GenericAsyncClass.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/GenericClass.cs",
   "name": "GenericClass.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/GenericClass.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_GenericClass.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/NotMatchingFileName.cs",
   "name": "NotMatchingFileName.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/NotMatchingFileName.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_NotMatchingFileName.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/PartialClass.cs",
   "name": "PartialClass.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/PartialClass.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_PartialClass.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/PartialClass2.cs",
   "name": "PartialClass2.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/PartialClass2.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_PartialClass2.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/PartialClassWithAutoProperties.cs",
   "name": "PartialClassWithAutoProperties.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/PartialClassWithAutoProperties.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_PartialClassWithAutoProperties.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/PartialClassWithAutoProperties2.cs",
   "name": "PartialClassWithAutoProperties2.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/PartialClassWithAutoProperties2.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_PartialClassWithAutoProperties2.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/Program.cs",
   "name": "Program.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/Program.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_Program.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/TestClass.cs",
   "name": "TestClass.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/TestClass.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_TestClass.cs.html"
  },
  {
   "id": "demo_projects/csharp/project/Test/TestClass2.cs",
   "name": "TestClass2.cs",
   "type": "file",
   "path": "demo_projects/csharp/project/Test/TestClass2.cs",
   "parentId": "demo_projects/csharp/project/Test",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "targetUrl": "demo_projects_csharp_project_Test_TestClass2.cs.html"
  },
  {
   "id": "demo_projects/csharp/nanovision.yaml",
   "name": "nanovision.yaml",
   "type": "file",
   "path": "demo_projects/csharp/nanovision.yaml",
   "parentId": "demo_projects/csharp",
   "depth": 2,
   "config": true
  },
  {
   "id": "demo_projects/go",
   "name": "go",
   "type": "folder",
   "path": "demo_projects/go",
   "parentId": "demo_projects",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 11,
     "total": 24,
     "percentage": 45.83
    },
    "methods_hit": {
     "covered": 16,
     "total": 24,
     "percentage": 66.66
    },
    "statement_coverage": {
     "covered": 42,
     "uncovered": 23,
     "coverable": 65,
     "total": 65,
     "percentage": 64.61
    },
    "statement_methods_fully_covered": {
     "covered": 11,
     "total": 24,
     "percentage": 45.83
    },
    "statement_methods_hit": {
     "covered": 16,
     "total": 24,
     "percentage": 66.66
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   }
  },
  {
   "id": "demo_projects/go/project",
   "name": "project",
   "type": "folder",
   "path": "demo_projects/go/project",
   "parentId": "demo_projects/go",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 11,
     "total": 24,
     "percentage": 45.83
    },
    "methods_hit": {
     "covered": 16,
     "total": 24,
     "percentage": 66.66
    },
    "statement_coverage": {
     "covered": 42,
     "uncovered": 23,
     "coverable": 65,
     "total": 65,
     "percentage": 64.61
    },
    "statement_methods_fully_covered": {
     "covered": 11,
     "total": 24,
     "percentage": 45.83
    },
    "statement_methods_hit": {
     "covered": 16,
     "total": 24,
     "percentage": 66.66
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   }
  },
  {
   "id": "demo_projects/go/project/calculator",
   "name": "calculator",
   "type": "folder",
   "path": "demo_projects/go/project/calculator",
   "parentId": "demo_projects/go/project",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 8,
     "total": 13,
     "percentage": 61.53
    },
    "methods_hit": {
     "covered": 11,
     "total": 13,
     "percentage": 84.61
    },
    "statement_coverage": {
     "covered": 29,
     "uncovered": 9,
     "coverable": 38,
     "total": 38,
     "percentage": 76.31
    },
    "statement_methods_fully_covered": {
     "covered": 8,
     "total": 13,
     "percentage": 61.53
    },
    "statement_methods_hit": {
     "covered": 11,
     "total": 13,
     "percentage": 84.61
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "demo_projects/go/project/calculator/calculator.go",
   "name": "calculator.go",
   "type": "file",
   "path": "demo_projects/go/project/calculator/calculator.go",
   "parentId": "demo_projects/go/project/calculator",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 7,
     "percentage": 42.85
    },
    "methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    },
    "statement_coverage": {
     "covered": 20,
     "uncovered": 6,
     "coverable": 26,
     "total": 26,
     "percentage": 76.92
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 7,
     "percentage": 42.85
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "demo_projects_go_project_calculator_calculator.go.html"
  },
  {
   "id": "demo_projects/go/project/calculator/entities.go",
   "name": "entities.go",
   "type": "file",
   "path": "demo_projects/go/project/calculator/entities.go",
   "parentId": "demo_projects/go/project/calculator",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 9,
     "uncovered": 3,
     "coverable": 12,
     "total": 12,
     "percentage": 75
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "demo_projects_go_project_calculator_entities.go.html"
  },
  {
   "id": "demo_projects/go/project/calculator_2",
   "name": "calculator_2",
   "type": "folder",
   "path": "demo_projects/go/project/calculator_2",
   "parentId": "demo_projects/go/project",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 11,
     "percentage": 27.27
    },
    "methods_hit": {
     "covered": 5,
     "total": 11,
     "percentage": 45.45
    },
    "statement_coverage": {
     "covered": 13,
     "uncovered": 14,
     "coverable": 27,
     "total": 27,
     "percentage": 48.14
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 11,
     "percentage": 27.27
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 11,
     "percentage": 45.45
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   }
  },
  {
   "id": "demo_projects/go/project/calculator_2/calculator.go",
   "name": "calculator.go",
   "type": "file",
   "path": "demo_projects/go/project/calculator_2/calculator.go",
   "parentId": "demo_projects/go/project/calculator_2",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "statement_coverage": {
     "covered": 10,
     "uncovered": 5,
     "coverable": 15,
     "total": 15,
     "percentage": 66.66
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "demo_projects_go_project_calculator_2_calculator.go.html"
  },
  {
   "id": "demo_projects/go/project/calculator_2/entities.go",
   "name": "entities.go",
   "type": "file",
   "path": "demo_projects/go/project/calculator_2/entities.go",
   "parentId": "demo_projects/go/project/calculator_2",
   "depth": 4,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "methods_hit": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "statement_coverage": {
     "covered": 3,
     "uncovered": 9,
     "coverable": 12,
     "total": 12,
     "percentage": 25
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "demo_projects_go_project_calculator_2_entities.go.html"
  },
  {
   "id": "demo_projects/go/nanovision.yaml",
   "name": "nanovision.yaml",
   "type": "file",
   "path": "demo_projects/go/nanovision.yaml",
   "parentId": "demo_projects/go",
   "depth": 2,
   "config": true
  },
  {
   "id": "internal",
   "name": "internal",
   "type": "folder",
   "path": "internal",
   "depth": 0,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 31
    },
    "methods_fully_covered": {
     "covered": 267,
     "total": 565,
     "percentage": 47.25
    },
    "methods_hit": {
     "covered": 512,
     "total": 565,
     "percentage": 90.61
    },
    "patch_methods_hit": {
     "covered": 307,
     "total": 329,
     "percentage": 93.31
    },
    "patch_statement_coverage": {
     "covered": 2558,
     "uncovered": 521,
     "coverable": 3079,
     "total": 3079,
     "percentage": 83.07
    },
    "patch_statement_methods_hit": {
     "covered": 310,
     "total": 329,
     "percentage": 94.22
    },
    "statement_coverage": {
     "covered": 4904,
     "uncovered": 1026,
     "coverable": 5930,
     "total": 5930,
     "percentage": 82.69
    },
    "statement_methods_fully_covered": {
     "covered": 267,
     "total": 565,
     "percentage": 47.25
    },
    "statement_methods_hit": {
     "covered": 512,
     "total": 565,
     "percentage": 90.61
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/aggregator",
   "name": "aggregator",
   "type": "folder",
   "path": "internal/aggregator",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 15
    },
    "methods_fully_covered": {
     "covered": 23,
     "total": 31,
     "percentage": 74.19
    },
    "methods_hit": {
     "covered": 31,
     "total": 31,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 1,
     "uncovered": 0,
     "coverable": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 328,
     "uncovered": 16,
     "coverable": 344,
     "total": 344,
     "percentage": 95.34
    },
    "statement_methods_fully_covered": {
     "covered": 23,
     "total": 31,
     "percentage": 74.19
    },
    "statement_methods_hit": {
     "covered": 31,
     "total": 31,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/aggregator/aggrgator.go",
   "name": "aggrgator.go",
   "type": "file",
   "path": "internal/aggregator/aggrgator.go",
   "parentId": "internal/aggregator",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 15
    },
    "methods_fully_covered": {
     "covered": 12,
     "total": 13,
     "percentage": 92.3
    },
    "methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 163,
     "uncovered": 5,
     "coverable": 168,
     "total": 168,
     "percentage": 97.02
    },
    "statement_methods_fully_covered": {
     "covered": 12,
     "total": 13,
     "percentage": 92.3
    },
    "statement_methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_aggregator_aggrgator.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/aggregator/report_index.go",
   "name": "report_index.go",
   "type": "file",
   "path": "internal/aggregator/report_index.go",
   "parentId": "internal/aggregator",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 10
    },
    "methods_fully_covered": {
     "covered": 11,
     "total": 18,
     "percentage": 61.11
    },
    "methods_hit": {
     "covered": 18,
     "total": 18,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 1,
     "uncovered": 0,
     "coverable": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 165,
     "uncovered": 11,
     "coverable": 176,
     "total": 176,
     "percentage": 93.75
    },
    "statement_methods_fully_covered": {
     "covered": 11,
     "total": 18,
     "percentage": 61.11
    },
    "statement_methods_hit": {
     "covered": 18,
     "total": 18,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_aggregator_report_index.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/analyzer",
   "name": "analyzer",
   "type": "folder",
   "path": "internal/analyzer",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 18
    },
    "methods_fully_covered": {
     "covered": 8,
     "total": 16,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 15,
     "total": 16,
     "percentage": 93.75
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 6,
     "coverable": 6,
     "total": 6,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 212,
     "uncovered": 27,
     "coverable": 239,
     "total": 239,
     "percentage": 88.7
    },
    "statement_methods_fully_covered": {
     "covered": 8,
     "total": 16,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 15,
     "total": 16,
     "percentage": 93.75
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/analyzer/cpp",
   "name": "cpp",
   "type": "folder",
   "path": "internal/analyzer/cpp",
   "parentId": "internal/analyzer",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 18
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 2,
     "coverable": 2,
     "total": 2,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 81,
     "uncovered": 9,
     "coverable": 90,
     "total": 90,
     "percentage": 90
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/analyzer/cpp/analyzer.go",
   "name": "analyzer.go",
   "type": "file",
   "path": "internal/analyzer/cpp/analyzer.go",
   "parentId": "internal/analyzer/cpp",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 18
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 2,
     "coverable": 2,
     "total": 2,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 81,
     "uncovered": 9,
     "coverable": 90,
     "total": 90,
     "percentage": 90
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_analyzer_cpp_analyzer.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/analyzer/gdscript",
   "name": "gdscript",
   "type": "folder",
   "path": "internal/analyzer/gdscript",
   "parentId": "internal/analyzer",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 11
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 2,
     "coverable": 2,
     "total": 2,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 60,
     "uncovered": 9,
     "coverable": 69,
     "total": 69,
     "percentage": 86.95
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/analyzer/gdscript/analyzer.go",
   "name": "analyzer.go",
   "type": "file",
   "path": "internal/analyzer/gdscript/analyzer.go",
   "parentId": "internal/analyzer/gdscript",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 11
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 2,
     "coverable": 2,
     "total": 2,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 60,
     "uncovered": 9,
     "coverable": 69,
     "total": 69,
     "percentage": 86.95
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_analyzer_gdscript_analyzer.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/analyzer/go",
   "name": "go",
   "type": "folder",
   "path": "internal/analyzer/go",
   "parentId": "internal/analyzer",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 2,
     "coverable": 2,
     "total": 2,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 71,
     "uncovered": 9,
     "coverable": 80,
     "total": 80,
     "percentage": 88.75
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/analyzer/go/analyzer.go",
   "name": "analyzer.go",
   "type": "file",
   "path": "internal/analyzer/go/analyzer.go",
   "parentId": "internal/analyzer/go",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 2,
     "coverable": 2,
     "total": 2,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 71,
     "uncovered": 9,
     "coverable": 80,
     "total": 80,
     "percentage": 88.75
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_analyzer_go_analyzer.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/bootlog",
   "name": "bootlog",
   "type": "folder",
   "path": "internal/bootlog",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 29,
     "uncovered": 0,
     "coverable": 29,
     "total": 29,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/bootlog/bootlog.go",
   "name": "bootlog.go",
   "type": "file",
   "path": "internal/bootlog/bootlog.go",
   "parentId": "internal/bootlog",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 29,
     "uncovered": 0,
     "coverable": 29,
     "total": 29,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_bootlog_bootlog.go.html"
  },
  {
   "id": "internal/cache",
   "name": "cache",
   "type": "folder",
   "path": "internal/cache",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 11,
     "percentage": 54.54
    },
    "methods_hit": {
     "covered": 9,
     "total": 11,
     "percentage": 81.81
    },
    "statement_coverage": {
     "covered": 83,
     "uncovered": 11,
     "coverable": 94,
     "total": 94,
     "percentage": 88.29
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 11,
     "percentage": 54.54
    },
    "statement_methods_hit": {
     "covered": 9,
     "total": 11,
     "percentage": 81.81
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/cache/cache.go",
   "name": "cache.go",
   "type": "file",
   "path": "internal/cache/cache.go",
   "parentId": "internal/cache",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 8,
     "percentage": 62.5
    },
    "methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 82,
     "uncovered": 9,
     "coverable": 91,
     "total": 91,
     "percentage": 90.1
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 8,
     "percentage": 62.5
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_cache_cache.go.html"
  },
  {
   "id": "internal/cache/cache_validator.go",
   "name": "cache_validator.go",
   "type": "file",
   "path": "internal/cache/cache_validator.go",
   "parentId": "internal/cache",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "statement_coverage": {
     "covered": 1,
     "uncovered": 2,
     "coverable": 3,
     "total": 3,
     "percentage": 33.33
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_cache_cache_validator.go.html"
  },
  {
   "id": "internal/calculator",
   "name": "calculator",
   "type": "folder",
   "path": "internal/calculator",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 4,
     "percentage": 75
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 51,
     "uncovered": 3,
     "coverable": 54,
     "total": 54,
     "percentage": 94.44
    },
    "patch_statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 71,
     "uncovered": 3,
     "coverable": 74,
     "total": 74,
     "percentage": 95.94
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 4,
     "percentage": 75
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/calculator/engine.go",
   "name": "engine.go",
   "type": "file",
   "path": "internal/calculator/engine.go",
   "parentId": "internal/calculator",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 20,
     "uncovered": 0,
     "coverable": 20,
     "total": 20,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 40,
     "uncovered": 0,
     "coverable": 40,
     "total": 40,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_calculator_engine.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/calculator/registry.go",
   "name": "registry.go",
   "type": "file",
   "path": "internal/calculator/registry.go",
   "parentId": "internal/calculator",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 3
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 31,
     "uncovered": 3,
     "coverable": 34,
     "total": 34,
     "percentage": 91.17
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 31,
     "uncovered": 3,
     "coverable": 34,
     "total": 34,
     "percentage": 91.17
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_calculator_registry.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/client",
   "name": "client",
   "type": "folder",
   "path": "internal/client",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 10
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 10,
     "percentage": 30
    },
    "methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 50,
     "uncovered": 16,
     "coverable": 66,
     "total": 66,
     "percentage": 75.75
    },
    "patch_statement_methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 50,
     "uncovered": 16,
     "coverable": 66,
     "total": 66,
     "percentage": 75.75
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 10,
     "percentage": 30
    },
    "statement_methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/client/client.go",
   "name": "client.go",
   "type": "file",
   "path": "internal/client/client.go",
   "parentId": "internal/client",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 10
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 10,
     "percentage": 30
    },
    "methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 50,
     "uncovered": 16,
     "coverable": 66,
     "total": 66,
     "percentage": 75.75
    },
    "patch_statement_methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 50,
     "uncovered": 16,
     "coverable": 66,
     "total": 66,
     "percentage": 75.75
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 10,
     "percentage": 30
    },
    "statement_methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_client_client.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/compare",
   "name": "compare",
   "type": "folder",
   "path": "internal/compare",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 27
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 9,
     "percentage": 55.55
    },
    "methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 139,
     "uncovered": 9,
     "coverable": 148,
     "total": 148,
     "percentage": 93.91
    },
    "patch_statement_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 139,
     "uncovered": 9,
     "coverable": 148,
     "total": 148,
     "percentage": 93.91
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 9,
     "percentage": 55.55
    },
    "statement_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/compare/compare.go",
   "name": "compare.go",
   "type": "file",
   "path": "internal/compare/compare.go",
   "parentId": "internal/compare",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 27
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 9,
     "percentage": 55.55
    },
    "methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 139,
     "uncovered": 9,
     "coverable": 148,
     "total": 148,
     "percentage": 93.91
    },
    "patch_statement_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 139,
     "uncovered": 9,
     "coverable": 148,
     "total": 148,
     "percentage": 93.91
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 9,
     "percentage": 55.55
    },
    "statement_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_compare_compare.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/config",
   "name": "config",
   "type": "folder",
   "path": "internal/config",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 28
    },
    "methods_fully_covered": {
     "covered": 24,
     "total": 47,
     "percentage": 51.06
    },
    "methods_hit": {
     "covered": 45,
     "total": 47,
     "percentage": 95.74
    },
    "patch_methods_hit": {
     "covered": 25,
     "total": 26,
     "percentage": 96.15
    },
    "patch_statement_coverage": {
     "covered": 277,
     "uncovered": 25,
     "coverable": 302,
     "total": 302,
     "percentage": 91.72
    },
    "patch_statement_methods_hit": {
     "covered": 25,
     "total": 26,
     "percentage": 96.15
    },
    "statement_coverage": {
     "covered": 518,
     "uncovered": 64,
     "coverable": 582,
     "total": 582,
     "percentage": 89
    },
    "statement_methods_fully_covered": {
     "covered": 24,
     "total": 47,
     "percentage": 51.06
    },
    "statement_methods_hit": {
     "covered": 45,
     "total": 47,
     "percentage": 95.74
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/config/config.go",
   "name": "config.go",
   "type": "file",
   "path": "internal/config/config.go",
   "parentId": "internal/config",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 28
    },
    "methods_fully_covered": {
     "covered": 11,
     "total": 26,
     "percentage": 42.3
    },
    "methods_hit": {
     "covered": 25,
     "total": 26,
     "percentage": 96.15
    },
    "patch_methods_hit": {
     "covered": 25,
     "total": 26,
     "percentage": 96.15
    },
    "patch_statement_coverage": {
     "covered": 277,
     "uncovered": 25,
     "coverable": 302,
     "total": 302,
     "percentage": 91.72
    },
    "patch_statement_methods_hit": {
     "covered": 25,
     "total": 26,
     "percentage": 96.15
    },
    "statement_coverage": {
     "covered": 325,
     "uncovered": 45,
     "coverable": 370,
     "total": 370,
     "percentage": 87.83
    },
    "statement_methods_fully_covered": {
     "covered": 11,
     "total": 26,
     "percentage": 42.3
    },
    "statement_methods_hit": {
     "covered": 25,
     "total": 26,
     "percentage": 96.15
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_config_config.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/config/metrics.go",
   "name": "metrics.go",
   "type": "file",
   "path": "internal/config/metrics.go",
   "parentId": "internal/config",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 7,
     "percentage": 71.42
    },
    "methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    },
    "statement_coverage": {
     "covered": 22,
     "uncovered": 5,
     "coverable": 27,
     "total": 27,
     "percentage": 81.48
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 7,
     "percentage": 71.42
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_config_metrics.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/config/outputs.go",
   "name": "outputs.go",
   "type": "file",
   "path": "internal/config/outputs.go",
   "parentId": "internal/config",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 13,
     "uncovered": 2,
     "coverable": 15,
     "total": 15,
     "percentage": 86.66
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_config_outputs.go.html"
  },
  {
   "id": "internal/config/schema.go",
   "name": "schema.go",
   "type": "file",
   "path": "internal/config/schema.go",
   "parentId": "internal/config",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 17
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 71,
     "uncovered": 0,
     "coverable": 71,
     "total": 71,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_config_schema.go.html"
  },
  {
   "id": "internal/config/scoped.go",
   "name": "scoped.go",
   "type": "file",
   "path": "internal/config/scoped.go",
   "parentId": "internal/config",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 87,
     "uncovered": 12,
     "coverable": 99,
     "total": 99,
     "percentage": 87.87
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_config_scoped.go.html"
  },
  {
   "id": "internal/diagnostics",
   "name": "diagnostics",
   "type": "folder",
   "path": "internal/diagnostics",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 10
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 11,
     "percentage": 36.36
    },
    "methods_hit": {
     "covered": 10,
     "total": 11,
     "percentage": 90.9
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 4,
     "uncovered": 0,
     "coverable": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 55,
     "uncovered": 21,
     "coverable": 76,
     "total": 76,
     "percentage": 72.36
    },
    "statement_methods_fully_covered": {
     "covered": 4,
     "total": 11,
     "percentage": 36.36
    },
    "statement_methods_hit": {
     "covered": 10,
     "total": 11,
     "percentage": 90.9
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/diagnostics/diagnostics.go",
   "name": "diagnostics.go",
   "type": "file",
   "path": "internal/diagnostics/diagnostics.go",
   "parentId": "internal/diagnostics",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 10
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 11,
     "percentage": 36.36
    },
    "methods_hit": {
     "covered": 10,
     "total": 11,
     "percentage": 90.9
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 4,
     "uncovered": 0,
     "coverable": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 55,
     "uncovered": 21,
     "coverable": 76,
     "total": 76,
     "percentage": 72.36
    },
    "statement_methods_fully_covered": {
     "covered": 4,
     "total": 11,
     "percentage": 36.36
    },
    "statement_methods_hit": {
     "covered": 10,
     "total": 11,
     "percentage": 90.9
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_diagnostics_diagnostics.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/diff",
   "name": "diff",
   "type": "folder",
   "path": "internal/diff",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 12
    },
    "methods_fully_covered": {
     "covered": 8,
     "total": 14,
     "percentage": 57.14
    },
    "methods_hit": {
     "covered": 14,
     "total": 14,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 13,
     "uncovered": 0,
     "coverable": 13,
     "total": 13,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 132,
     "uncovered": 10,
     "coverable": 142,
     "total": 142,
     "percentage": 92.95
    },
    "statement_methods_fully_covered": {
     "covered": 8,
     "total": 14,
     "percentage": 57.14
    },
    "statement_methods_hit": {
     "covered": 14,
     "total": 14,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/diff/parser.go",
   "name": "parser.go",
   "type": "file",
   "path": "internal/diff/parser.go",
   "parentId": "internal/diff",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 12
    },
    "methods_fully_covered": {
     "covered": 8,
     "total": 13,
     "percentage": 61.53
    },
    "methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 7,
     "uncovered": 0,
     "coverable": 7,
     "total": 7,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 121,
     "uncovered": 9,
     "coverable": 130,
     "total": 130,
     "percentage": 93.07
    },
    "statement_methods_fully_covered": {
     "covered": 8,
     "total": 13,
     "percentage": 61.53
    },
    "statement_methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_diff_parser.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/diff/path.go",
   "name": "path.go",
   "type": "file",
   "path": "internal/diff/path.go",
   "parentId": "internal/diff",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 6,
     "uncovered": 0,
     "coverable": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 11,
     "uncovered": 1,
     "coverable": 12,
     "total": 12,
     "percentage": 91.66
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_diff_path.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/diffapply",
   "name": "diffapply",
   "type": "folder",
   "path": "internal/diffapply",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 19
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 8,
     "percentage": 25
    },
    "methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 50,
     "uncovered": 8,
     "coverable": 58,
     "total": 58,
     "percentage": 86.2
    },
    "patch_statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 194,
     "uncovered": 17,
     "coverable": 211,
     "total": 211,
     "percentage": 91.94
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 8,
     "percentage": 25
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/diffapply/apply.go",
   "name": "apply.go",
   "type": "file",
   "path": "internal/diffapply/apply.go",
   "parentId": "internal/diffapply",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 19
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 23,
     "uncovered": 1,
     "coverable": 24,
     "total": 24,
     "percentage": 95.83
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 48,
     "uncovered": 3,
     "coverable": 51,
     "total": 51,
     "percentage": 94.11
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_diffapply_apply.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/diffapply/resolver.go",
   "name": "resolver.go",
   "type": "file",
   "path": "internal/diffapply/resolver.go",
   "parentId": "internal/diffapply",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 15
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 27,
     "uncovered": 7,
     "coverable": 34,
     "total": 34,
     "percentage": 79.41
    },
    "patch_statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 146,
     "uncovered": 14,
     "coverable": 160,
     "total": 160,
     "percentage": 91.25
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_diffapply_resolver.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/enricher",
   "name": "enricher",
   "type": "folder",
   "path": "internal/enricher",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 11
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 8,
     "percentage": 75
    },
    "methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 3,
     "uncovered": 1,
     "coverable": 4,
     "total": 4,
     "percentage": 75
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 97,
     "uncovered": 12,
     "coverable": 109,
     "total": 109,
     "percentage": 88.99
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 8,
     "percentage": 75
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/enricher/enricher.go",
   "name": "enricher.go",
   "type": "file",
   "path": "internal/enricher/enricher.go",
   "parentId": "internal/enricher",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 11
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 8,
     "percentage": 75
    },
    "methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 3,
     "uncovered": 1,
     "coverable": 4,
     "total": 4,
     "percentage": 75
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 97,
     "uncovered": 12,
     "coverable": 109,
     "total": 109,
     "percentage": 88.99
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 8,
     "percentage": 75
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_enricher_enricher.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/filereader",
   "name": "filereader",
   "type": "folder",
   "path": "internal/filereader",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 10,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 9,
     "total": 10,
     "percentage": 90
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 16,
     "uncovered": 4,
     "coverable": 20,
     "total": 20,
     "percentage": 80
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 71,
     "uncovered": 10,
     "coverable": 81,
     "total": 81,
     "percentage": 87.65
    },
    "statement_methods_fully_covered": {
     "covered": 4,
     "total": 10,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 9,
     "total": 10,
     "percentage": 90
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/filereader/default_reader.go",
   "name": "default_reader.go",
   "type": "file",
   "path": "internal/filereader/default_reader.go",
   "parentId": "internal/filereader",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 1
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "statement_coverage": {
     "covered": 4,
     "uncovered": 1,
     "coverable": 5,
     "total": 5,
     "percentage": 80
    },
    "statement_methods_fully_covered": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_filereader_default_reader.go.html"
  },
  {
   "id": "internal/filereader/filereader.go",
   "name": "filereader.go",
   "type": "file",
   "path": "internal/filereader/filereader.go",
   "parentId": "internal/filereader",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 16,
     "uncovered": 4,
     "coverable": 20,
     "total": 20,
     "percentage": 80
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 67,
     "uncovered": 9,
     "coverable": 76,
     "total": 76,
     "percentage": 88.15
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_filereader_filereader.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/filtering",
   "name": "filtering",
   "type": "folder",
   "path": "internal/filtering",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 11
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 45,
     "uncovered": 3,
     "coverable": 48,
     "total": 48,
     "percentage": 93.75
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/filtering/filter.go",
   "name": "filter.go",
   "type": "file",
   "path": "internal/filtering/filter.go",
   "parentId": "internal/filtering",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 11
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 45,
     "uncovered": 3,
     "coverable": 48,
     "total": 48,
     "percentage": 93.75
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_filtering_filter.go.html"
  },
  {
   "id": "internal/history",
   "name": "history",
   "type": "folder",
   "path": "internal/history",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 18
    },
    "methods_fully_covered": {
     "covered": 12,
     "total": 24,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 24,
     "total": 24,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 24,
     "total": 24,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 283,
     "uncovered": 40,
     "coverable": 323,
     "total": 323,
     "percentage": 87.61
    },
    "patch_statement_methods_hit": {
     "covered": 24,
     "total": 24,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 283,
     "uncovered": 40,
     "coverable": 323,
     "total": 323,
     "percentage": 87.61
    },
    "statement_methods_fully_covered": {
     "covered": 12,
     "total": 24,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 24,
     "total": 24,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/history/capture.go",
   "name": "capture.go",
   "type": "file",
   "path": "internal/history/capture.go",
   "parentId": "internal/history",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 11
    },
    "methods_fully_covered": {
     "covered": 7,
     "total": 9,
     "percentage": 77.77
    },
    "methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 67,
     "uncovered": 2,
     "coverable": 69,
     "total": 69,
     "percentage": 97.1
    },
    "patch_statement_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 67,
     "uncovered": 2,
     "coverable": 69,
     "total": 69,
     "percentage": 97.1
    },
    "statement_methods_fully_covered": {
     "covered": 7,
     "total": 9,
     "percentage": 77.77
    },
    "statement_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_history_capture.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/history/meta.go",
   "name": "meta.go",
   "type": "file",
   "path": "internal/history/meta.go",
   "parentId": "internal/history",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 47,
     "uncovered": 11,
     "coverable": 58,
     "total": 58,
     "percentage": 81.03
    },
    "patch_statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 47,
     "uncovered": 11,
     "coverable": 58,
     "total": 58,
     "percentage": 81.03
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_history_meta.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/history/rebuild.go",
   "name": "rebuild.go",
   "type": "file",
   "path": "internal/history/rebuild.go",
   "parentId": "internal/history",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 18
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 72,
     "uncovered": 11,
     "coverable": 83,
     "total": 83,
     "percentage": 86.74
    },
    "patch_statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 72,
     "uncovered": 11,
     "coverable": 83,
     "total": 83,
     "percentage": 86.74
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_history_rebuild.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/history/record.go",
   "name": "record.go",
   "type": "file",
   "path": "internal/history/record.go",
   "parentId": "internal/history",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 15
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 7,
     "percentage": 28.57
    },
    "methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 97,
     "uncovered": 16,
     "coverable": 113,
     "total": 113,
     "percentage": 85.84
    },
    "patch_statement_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 97,
     "uncovered": 16,
     "coverable": 113,
     "total": 113,
     "percentage": 85.84
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 7,
     "percentage": 28.57
    },
    "statement_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_history_record.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/logging",
   "name": "logging",
   "type": "folder",
   "path": "internal/logging",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 10,
     "percentage": 60
    },
    "methods_hit": {
     "covered": 7,
     "total": 10,
     "percentage": 70
    },
    "statement_coverage": {
     "covered": 54,
     "uncovered": 9,
     "coverable": 63,
     "total": 63,
     "percentage": 85.71
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 10,
     "percentage": 60
    },
    "statement_methods_hit": {
     "covered": 7,
     "total": 10,
     "percentage": 70
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/logging/logging.go",
   "name": "logging.go",
   "type": "file",
   "path": "internal/logging/logging.go",
   "parentId": "internal/logging",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 10,
     "percentage": 60
    },
    "methods_hit": {
     "covered": 7,
     "total": 10,
     "percentage": 70
    },
    "statement_coverage": {
     "covered": 54,
     "uncovered": 9,
     "coverable": 63,
     "total": 63,
     "percentage": 85.71
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 10,
     "percentage": 60
    },
    "statement_methods_hit": {
     "covered": 7,
     "total": 10,
     "percentage": 70
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_logging_logging.go.html"
  },
  {
   "id": "internal/model",
   "name": "model",
   "type": "folder",
   "path": "internal/model",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 8,
     "uncovered": 2,
     "coverable": 10,
     "total": 10,
     "percentage": 80
    },
    "patch_statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 19,
     "uncovered": 2,
     "coverable": 21,
     "total": 21,
     "percentage": 90.47
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/model/comparison.go",
   "name": "comparison.go",
   "type": "file",
   "path": "internal/model/comparison.go",
   "parentId": "internal/model",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 8,
     "uncovered": 2,
     "coverable": 10,
     "total": 10,
     "percentage": 80
    },
    "patch_statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 8,
     "uncovered": 2,
     "coverable": 10,
     "total": 10,
     "percentage": 80
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_model_comparison.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/model/diff.go",
   "name": "diff.go",
   "type": "file",
   "path": "internal/model/diff.go",
   "parentId": "internal/model",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 3
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 11,
     "uncovered": 0,
     "coverable": 11,
     "total": 11,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_model_diff.go.html"
  },
  {
   "id": "internal/parsers",
   "name": "parsers",
   "type": "folder",
   "path": "internal/parsers",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 16
    },
    "methods_fully_covered": {
     "covered": 19,
     "total": 36,
     "percentage": 52.77
    },
    "methods_hit": {
     "covered": 33,
     "total": 36,
     "percentage": 91.66
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 4,
     "uncovered": 0,
     "coverable": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 236,
     "uncovered": 48,
     "coverable": 284,
     "total": 284,
     "percentage": 83.09
    },
    "statement_methods_fully_covered": {
     "covered": 19,
     "total": 36,
     "percentage": 52.77
    },
    "statement_methods_hit": {
     "covered": 33,
     "total": 36,
     "percentage": 91.66
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/parsers/parser_cobertura",
   "name": "parser_cobertura",
   "type": "folder",
   "path": "internal/parsers/parser_cobertura",
   "parentId": "internal/parsers",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 68,
     "uncovered": 15,
     "coverable": 83,
     "total": 83,
     "percentage": 81.92
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/parsers/parser_cobertura/parser.go",
   "name": "parser.go",
   "type": "file",
   "path": "internal/parsers/parser_cobertura/parser.go",
   "parentId": "internal/parsers/parser_cobertura",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 6,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 37,
     "uncovered": 12,
     "coverable": 49,
     "total": 49,
     "percentage": 75.51
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 6,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_cobertura_parser.go.html"
  },
  {
   "id": "internal/parsers/parser_cobertura/processing.go",
   "name": "processing.go",
   "type": "file",
   "path": "internal/parsers/parser_cobertura/processing.go",
   "parentId": "internal/parsers/parser_cobertura",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 31,
     "uncovered": 3,
     "coverable": 34,
     "total": 34,
     "percentage": 91.17
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_cobertura_processing.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/parsers/parser_gcov",
   "name": "parser_gcov",
   "type": "folder",
   "path": "internal/parsers/parser_gcov",
   "parentId": "internal/parsers",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 45,
     "uncovered": 4,
     "coverable": 49,
     "total": 49,
     "percentage": 91.83
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/parsers/parser_gcov/parser.go",
   "name": "parser.go",
   "type": "file",
   "path": "internal/parsers/parser_gcov/parser.go",
   "parentId": "internal/parsers/parser_gcov",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 19,
     "uncovered": 3,
     "coverable": 22,
     "total": 22,
     "percentage": 86.36
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_gcov_parser.go.html"
  },
  {
   "id": "internal/parsers/parser_gcov/processing.go",
   "name": "processing.go",
   "type": "file",
   "path": "internal/parsers/parser_gcov/processing.go",
   "parentId": "internal/parsers/parser_gcov",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 26,
     "uncovered": 1,
     "coverable": 27,
     "total": 27,
     "percentage": 96.29
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_gcov_processing.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/parsers/parser_gocover",
   "name": "parser_gocover",
   "type": "folder",
   "path": "internal/parsers/parser_gocover",
   "parentId": "internal/parsers",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 9,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 1,
     "uncovered": 0,
     "coverable": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 63,
     "uncovered": 6,
     "coverable": 69,
     "total": 69,
     "percentage": 91.3
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 9,
     "percentage": 66.66
    },
    "statement_methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/parsers/parser_gocover/parser.go",
   "name": "parser.go",
   "type": "file",
   "path": "internal/parsers/parser_gocover/parser.go",
   "parentId": "internal/parsers/parser_gocover",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 34,
     "uncovered": 6,
     "coverable": 40,
     "total": 40,
     "percentage": 85
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_gocover_parser.go.html"
  },
  {
   "id": "internal/parsers/parser_gocover/processing.go",
   "name": "processing.go",
   "type": "file",
   "path": "internal/parsers/parser_gocover/processing.go",
   "parentId": "internal/parsers/parser_gocover",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 5
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 1,
     "uncovered": 0,
     "coverable": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 29,
     "uncovered": 0,
     "coverable": 29,
     "total": 29,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_gocover_processing.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/parsers/parser_lcov",
   "name": "parser_lcov",
   "type": "folder",
   "path": "internal/parsers/parser_lcov",
   "parentId": "internal/parsers",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 16
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 3,
     "uncovered": 0,
     "coverable": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 49,
     "uncovered": 19,
     "coverable": 68,
     "total": 68,
     "percentage": 72.05
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/parsers/parser_lcov/parser.go",
   "name": "parser.go",
   "type": "file",
   "path": "internal/parsers/parser_lcov/parser.go",
   "parentId": "internal/parsers/parser_lcov",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 3,
     "total": 4,
     "percentage": 75
    },
    "statement_coverage": {
     "covered": 8,
     "uncovered": 15,
     "coverable": 23,
     "total": 23,
     "percentage": 34.78
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 4,
     "percentage": 75
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_parsers_parser_lcov_parser.go.html"
  },
  {
   "id": "internal/parsers/parser_lcov/processing.go",
   "name": "processing.go",
   "type": "file",
   "path": "internal/parsers/parser_lcov/processing.go",
   "parentId": "internal/parsers/parser_lcov",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 16
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 3,
     "uncovered": 0,
     "coverable": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 41,
     "uncovered": 4,
     "coverable": 45,
     "total": 45,
     "percentage": 91.11
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_lcov_processing.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/parsers/factory.go",
   "name": "factory.go",
   "type": "file",
   "path": "internal/parsers/factory.go",
   "parentId": "internal/parsers",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 3
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "statement_coverage": {
     "covered": 9,
     "uncovered": 3,
     "coverable": 12,
     "total": 12,
     "percentage": 75
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_parsers_factory.go.html"
  },
  {
   "id": "internal/parsers/parser_config.go",
   "name": "parser_config.go",
   "type": "file",
   "path": "internal/parsers/parser_config.go",
   "parentId": "internal/parsers",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 1
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "statement_coverage": {
     "covered": 2,
     "uncovered": 1,
     "coverable": 3,
     "total": 3,
     "percentage": 66.66
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_parsers_parser_config.go.html"
  },
  {
   "id": "internal/pipeline",
   "name": "pipeline",
   "type": "folder",
   "path": "internal/pipeline",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 67,
     "uncovered": 19,
     "coverable": 86,
     "total": 86,
     "percentage": 77.9
    },
    "patch_statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 67,
     "uncovered": 19,
     "coverable": 86,
     "total": 86,
     "percentage": 77.9
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/pipeline/pipeline.go",
   "name": "pipeline.go",
   "type": "file",
   "path": "internal/pipeline/pipeline.go",
   "parentId": "internal/pipeline",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 67,
     "uncovered": 19,
     "coverable": 86,
     "total": 86,
     "percentage": 77.9
    },
    "patch_statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 67,
     "uncovered": 19,
     "coverable": 86,
     "total": 86,
     "percentage": 77.9
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_pipeline_pipeline.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/reporter",
   "name": "reporter",
   "type": "folder",
   "path": "internal/reporter",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 29
    },
    "methods_fully_covered": {
     "covered": 35,
     "total": 97,
     "percentage": 36.08
    },
    "methods_hit": {
     "covered": 70,
     "total": 97,
     "percentage": 72.16
    },
    "patch_methods_hit": {
     "covered": 46,
     "total": 52,
     "percentage": 88.46
    },
    "patch_statement_coverage": {
     "covered": 277,
     "uncovered": 58,
     "coverable": 335,
     "total": 335,
     "percentage": 82.68
    },
    "patch_statement_methods_hit": {
     "covered": 46,
     "total": 52,
     "percentage": 88.46
    },
    "statement_coverage": {
     "covered": 637,
     "uncovered": 326,
     "coverable": 963,
     "total": 963,
     "percentage": 66.14
    },
    "statement_methods_fully_covered": {
     "covered": 35,
     "total": 97,
     "percentage": 36.08
    },
    "statement_methods_hit": {
     "covered": 70,
     "total": 97,
     "percentage": 72.16
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/reporter/annotations",
   "name": "annotations",
   "type": "folder",
   "path": "internal/reporter/annotations",
   "parentId": "internal/reporter",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 0,
     "uncovered": 32,
     "coverable": 32,
     "total": 32,
     "percentage": 0
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   }
  },
  {
   "id": "internal/reporter/annotations/reporter.go",
   "name": "reporter.go",
   "type": "file",
   "path": "internal/reporter/annotations/reporter.go",
   "parentId": "internal/reporter/annotations",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 0,
     "uncovered": 32,
     "coverable": 32,
     "total": 32,
     "percentage": 0
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_annotations_reporter.go.html"
  },
  {
   "id": "internal/reporter/htmlreact",
   "name": "htmlreact",
   "type": "folder",
   "path": "internal/reporter/htmlreact",
   "parentId": "internal/reporter",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 16
    },
    "methods_fully_covered": {
     "covered": 27,
     "total": 61,
     "percentage": 44.26
    },
    "methods_hit": {
     "covered": 50,
     "total": 61,
     "percentage": 81.96
    },
    "patch_methods_hit": {
     "covered": 34,
     "total": 40,
     "percentage": 85
    },
    "patch_statement_coverage": {
     "covered": 152,
     "uncovered": 35,
     "coverable": 187,
     "total": 187,
     "percentage": 81.28
    },
    "patch_statement_methods_hit": {
     "covered": 34,
     "total": 40,
     "percentage": 85
    },
    "statement_coverage": {
     "covered": 426,
     "uncovered": 188,
     "coverable": 614,
     "total": 614,
     "percentage": 69.38
    },
    "statement_methods_fully_covered": {
     "covered": 27,
     "total": 61,
     "percentage": 44.26
    },
    "statement_methods_hit": {
     "covered": 50,
     "total": 61,
     "percentage": 81.96
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/reporter/htmlreact/builder.go",
   "name": "builder.go",
   "type": "file",
   "path": "internal/reporter/htmlreact/builder.go",
   "parentId": "internal/reporter/htmlreact",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 16
    },
    "methods_fully_covered": {
     "covered": 19,
     "total": 34,
     "percentage": 55.88
    },
    "methods_hit": {
     "covered": 31,
     "total": 34,
     "percentage": 91.17
    },
    "patch_methods_hit": {
     "covered": 27,
     "total": 28,
     "percentage": 96.42
    },
    "patch_statement_coverage": {
     "covered": 121,
     "uncovered": 11,
     "coverable": 132,
     "total": 132,
     "percentage": 91.66
    },
    "patch_statement_methods_hit": {
     "covered": 27,
     "total": 28,
     "percentage": 96.42
    },
    "statement_coverage": {
     "covered": 244,
     "uncovered": 53,
     "coverable": 297,
     "total": 297,
     "percentage": 82.15
    },
    "statement_methods_fully_covered": {
     "covered": 19,
     "total": 34,
     "percentage": 55.88
    },
    "statement_methods_hit": {
     "covered": 31,
     "total": 34,
     "percentage": 91.17
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_htmlreact_builder.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/reporter/htmlreact/details_generator.go",
   "name": "details_generator.go",
   "type": "file",
   "path": "internal/reporter/htmlreact/details_generator.go",
   "parentId": "internal/reporter/htmlreact",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 11
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 15,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 14,
     "total": 15,
     "percentage": 93.33
    },
    "patch_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 25,
     "uncovered": 8,
     "coverable": 33,
     "total": 33,
     "percentage": 75.75
    },
    "patch_statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 127,
     "uncovered": 31,
     "coverable": 158,
     "total": 158,
     "percentage": 80.37
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 15,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 14,
     "total": 15,
     "percentage": 93.33
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_htmlreact_details_generator.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/reporter/htmlreact/embed.go",
   "name": "embed.go",
   "type": "file",
   "path": "internal/reporter/htmlreact/embed.go",
   "parentId": "internal/reporter/htmlreact",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 1
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 1,
     "coverable": 1,
     "total": 1,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 1,
     "uncovered": 1,
     "coverable": 2,
     "total": 2,
     "percentage": 50
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_htmlreact_embed.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/reporter/htmlreact/emit.go",
   "name": "emit.go",
   "type": "file",
   "path": "internal/reporter/htmlreact/emit.go",
   "parentId": "internal/reporter/htmlreact",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 14,
     "uncovered": 4,
     "coverable": 18,
     "total": 18,
     "percentage": 77.77
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 1,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_htmlreact_emit.go.html"
  },
  {
   "id": "internal/reporter/htmlreact/generator.go",
   "name": "generator.go",
   "type": "file",
   "path": "internal/reporter/htmlreact/generator.go",
   "parentId": "internal/reporter/htmlreact",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 34,
     "uncovered": 16,
     "coverable": 50,
     "total": 50,
     "percentage": 68
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_reporter_htmlreact_generator.go.html"
  },
  {
   "id": "internal/reporter/htmlreact/generator_single.go",
   "name": "generator_single.go",
   "type": "file",
   "path": "internal/reporter/htmlreact/generator_single.go",
   "parentId": "internal/reporter/htmlreact",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 0,
     "uncovered": 68,
     "coverable": 68,
     "total": 68,
     "percentage": 0
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_htmlreact_generator_single.go.html"
  },
  {
   "id": "internal/reporter/htmlreact/schema.go",
   "name": "schema.go",
   "type": "file",
   "path": "internal/reporter/htmlreact/schema.go",
   "parentId": "internal/reporter/htmlreact",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 6,
     "uncovered": 0,
     "coverable": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 6,
     "uncovered": 0,
     "coverable": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_htmlreact_schema.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/reporter/htmlreact/views.go",
   "name": "views.go",
   "type": "file",
   "path": "internal/reporter/htmlreact/views.go",
   "parentId": "internal/reporter/htmlreact",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 4,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 0,
     "total": 4,
     "percentage": 0
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 4,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 15,
     "coverable": 15,
     "total": 15,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 0,
     "total": 4,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 0,
     "uncovered": 15,
     "coverable": 15,
     "total": 15,
     "percentage": 0
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 4,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 0,
     "total": 4,
     "percentage": 0
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_htmlreact_views.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/reporter/lcov",
   "name": "lcov",
   "type": "folder",
   "path": "internal/reporter/lcov",
   "parentId": "internal/reporter",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "statement_coverage": {
     "covered": 52,
     "uncovered": 12,
     "coverable": 64,
     "total": 64,
     "percentage": 81.25
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/reporter/lcov/reporter.go",
   "name": "reporter.go",
   "type": "file",
   "path": "internal/reporter/lcov/reporter.go",
   "parentId": "internal/reporter/lcov",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    },
    "statement_coverage": {
     "covered": 52,
     "uncovered": 12,
     "coverable": 64,
     "total": 64,
     "percentage": 81.25
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 5,
     "percentage": 80
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_lcov_reporter.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/reporter/reporter_rawjson",
   "name": "reporter_rawjson",
   "type": "folder",
   "path": "internal/reporter/reporter_rawjson",
   "parentId": "internal/reporter",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 3
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "statement_coverage": {
     "covered": 8,
     "uncovered": 3,
     "coverable": 11,
     "total": 11,
     "percentage": 72.72
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/reporter/reporter_rawjson/reporter.go",
   "name": "reporter.go",
   "type": "file",
   "path": "internal/reporter/reporter_rawjson/reporter.go",
   "parentId": "internal/reporter/reporter_rawjson",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 3
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    },
    "statement_coverage": {
     "covered": 8,
     "uncovered": 3,
     "coverable": 11,
     "total": 11,
     "percentage": 72.72
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 3,
     "percentage": 66.66
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_reporter_reporter_rawjson_reporter.go.html"
  },
  {
   "id": "internal/reporter/sarif",
   "name": "sarif",
   "type": "folder",
   "path": "internal/reporter/sarif",
   "parentId": "internal/reporter",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 3
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 0,
     "uncovered": 33,
     "coverable": 33,
     "total": 33,
     "percentage": 0
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   }
  },
  {
   "id": "internal/reporter/sarif/reporter.go",
   "name": "reporter.go",
   "type": "file",
   "path": "internal/reporter/sarif/reporter.go",
   "parentId": "internal/reporter/sarif",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 3
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 0,
     "uncovered": 33,
     "coverable": 33,
     "total": 33,
     "percentage": 0
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_sarif_reporter.go.html"
  },
  {
   "id": "internal/reporter/textsummary",
   "name": "textsummary",
   "type": "folder",
   "path": "internal/reporter/textsummary",
   "parentId": "internal/reporter",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 29
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 17,
     "percentage": 29.41
    },
    "methods_hit": {
     "covered": 14,
     "total": 17,
     "percentage": 82.35
    },
    "patch_methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 125,
     "uncovered": 23,
     "coverable": 148,
     "total": 148,
     "percentage": 84.45
    },
    "patch_statement_methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 151,
     "uncovered": 58,
     "coverable": 209,
     "total": 209,
     "percentage": 72.24
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 17,
     "percentage": 29.41
    },
    "statement_methods_hit": {
     "covered": 14,
     "total": 17,
     "percentage": 82.35
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/reporter/textsummary/comparison.go",
   "name": "comparison.go",
   "type": "file",
   "path": "internal/reporter/textsummary/comparison.go",
   "parentId": "internal/reporter/textsummary",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 17
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 8,
     "percentage": 37.5
    },
    "methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 65,
     "uncovered": 10,
     "coverable": 75,
     "total": 75,
     "percentage": 86.66
    },
    "patch_statement_methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 65,
     "uncovered": 10,
     "coverable": 75,
     "total": 75,
     "percentage": 86.66
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 8,
     "percentage": 37.5
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_textsummary_comparison.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/reporter/textsummary/reporter.go",
   "name": "reporter.go",
   "type": "file",
   "path": "internal/reporter/textsummary/reporter.go",
   "parentId": "internal/reporter/textsummary",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 10
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "methods_hit": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 3,
     "uncovered": 0,
     "coverable": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 29,
     "uncovered": 35,
     "coverable": 64,
     "total": 64,
     "percentage": 45.31
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 6,
     "percentage": 50
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_textsummary_reporter.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/reporter/textsummary/terminal.go",
   "name": "terminal.go",
   "type": "file",
   "path": "internal/reporter/textsummary/terminal.go",
   "parentId": "internal/reporter/textsummary",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 29
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 57,
     "uncovered": 13,
     "coverable": 70,
     "total": 70,
     "percentage": 81.42
    },
    "patch_statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 57,
     "uncovered": 13,
     "coverable": 70,
     "total": 70,
     "percentage": 81.42
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 3,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_textsummary_terminal.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/review",
   "name": "review",
   "type": "folder",
   "path": "internal/review",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 19
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 6,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 62,
     "uncovered": 12,
     "coverable": 74,
     "total": 74,
     "percentage": 83.78
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 6,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/review/review.go",
   "name": "review.go",
   "type": "file",
   "path": "internal/review/review.go",
   "parentId": "internal/review",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 19
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 6,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 62,
     "uncovered": 12,
     "coverable": 74,
     "total": 74,
     "percentage": 83.78
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 6,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_review_review.go.html"
  },
  {
   "id": "internal/server",
   "name": "server",
   "type": "folder",
   "path": "internal/server",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 12
    },
    "methods_fully_covered": {
     "covered": 12,
     "total": 40,
     "percentage": 30
    },
    "methods_hit": {
     "covered": 37,
     "total": 40,
     "percentage": 92.5
    },
    "patch_methods_hit": {
     "covered": 37,
     "total": 40,
     "percentage": 92.5
    },
    "patch_statement_coverage": {
     "covered": 315,
     "uncovered": 125,
     "coverable": 440,
     "total": 440,
     "percentage": 71.59
    },
    "patch_statement_methods_hit": {
     "covered": 37,
     "total": 40,
     "percentage": 92.5
    },
    "statement_coverage": {
     "covered": 315,
     "uncovered": 125,
     "coverable": 440,
     "total": 440,
     "percentage": 71.59
    },
    "statement_methods_fully_covered": {
     "covered": 12,
     "total": 40,
     "percentage": 30
    },
    "statement_methods_hit": {
     "covered": 37,
     "total": 40,
     "percentage": 92.5
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/server/api.go",
   "name": "api.go",
   "type": "file",
   "path": "internal/server/api.go",
   "parentId": "internal/server",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 14,
     "percentage": 14.28
    },
    "methods_hit": {
     "covered": 13,
     "total": 14,
     "percentage": 92.85
    },
    "patch_methods_hit": {
     "covered": 13,
     "total": 14,
     "percentage": 92.85
    },
    "patch_statement_coverage": {
     "covered": 129,
     "uncovered": 43,
     "coverable": 172,
     "total": 172,
     "percentage": 75
    },
    "patch_statement_methods_hit": {
     "covered": 13,
     "total": 14,
     "percentage": 92.85
    },
    "statement_coverage": {
     "covered": 129,
     "uncovered": 43,
     "coverable": 172,
     "total": 172,
     "percentage": 75
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 14,
     "percentage": 14.28
    },
    "statement_methods_hit": {
     "covered": 13,
     "total": 14,
     "percentage": 92.85
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_server_api.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/server/cache.go",
   "name": "cache.go",
   "type": "file",
   "path": "internal/server/cache.go",
   "parentId": "internal/server",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 5
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 38,
     "uncovered": 9,
     "coverable": 47,
     "total": 47,
     "percentage": 80.85
    },
    "patch_statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 38,
     "uncovered": 9,
     "coverable": 47,
     "total": 47,
     "percentage": 80.85
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_server_cache.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/server/maintenance.go",
   "name": "maintenance.go",
   "type": "file",
   "path": "internal/server/maintenance.go",
   "parentId": "internal/server",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 8
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "patch_methods_hit": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 24,
     "coverable": 24,
     "total": 24,
     "percentage": 0
    },
    "patch_statement_methods_hit": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_coverage": {
     "covered": 0,
     "uncovered": 24,
     "coverable": 24,
     "total": 24,
     "percentage": 0
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    }
   },
   "statuses": {
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_server_maintenance.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/server/server.go",
   "name": "server.go",
   "type": "file",
   "path": "internal/server/server.go",
   "parentId": "internal/server",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 7,
     "total": 11,
     "percentage": 63.63
    },
    "methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 55,
     "uncovered": 7,
     "coverable": 62,
     "total": 62,
     "percentage": 88.7
    },
    "patch_statement_methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 55,
     "uncovered": 7,
     "coverable": 62,
     "total": 62,
     "percentage": 88.7
    },
    "statement_methods_fully_covered": {
     "covered": 7,
     "total": 11,
     "percentage": 63.63
    },
    "statement_methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_server_server.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/server/ui.go",
   "name": "ui.go",
   "type": "file",
   "path": "internal/server/ui.go",
   "parentId": "internal/server",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 12
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 24,
     "uncovered": 17,
     "coverable": 41,
     "total": 41,
     "percentage": 58.53
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 24,
     "uncovered": 17,
     "coverable": 41,
     "total": 41,
     "percentage": 58.53
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_server_ui.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/server/upload.go",
   "name": "upload.go",
   "type": "file",
   "path": "internal/server/upload.go",
   "parentId": "internal/server",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 34,
     "uncovered": 14,
     "coverable": 48,
     "total": 48,
     "percentage": 70.83
    },
    "patch_statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 34,
     "uncovered": 14,
     "coverable": 48,
     "total": 48,
     "percentage": 70.83
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 3,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_server_upload.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/server/wire.go",
   "name": "wire.go",
   "type": "file",
   "path": "internal/server/wire.go",
   "parentId": "internal/server",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 4,
     "percentage": 25
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 35,
     "uncovered": 11,
     "coverable": 46,
     "total": 46,
     "percentage": 76.08
    },
    "patch_statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 35,
     "uncovered": 11,
     "coverable": 46,
     "total": 46,
     "percentage": 76.08
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 4,
     "percentage": 25
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_server_wire.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/status",
   "name": "status",
   "type": "folder",
   "path": "internal/status",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 13,
     "total": 16,
     "percentage": 81.25
    },
    "methods_hit": {
     "covered": 15,
     "total": 16,
     "percentage": 93.75
    },
    "patch_methods_hit": {
     "covered": 7,
     "total": 8,
     "percentage": 87.5
    },
    "patch_statement_coverage": {
     "covered": 37,
     "uncovered": 1,
     "coverable": 38,
     "total": 38,
     "percentage": 97.36
    },
    "patch_statement_methods_hit": {
     "covered": 7,
     "total": 8,
     "percentage": 87.5
    },
    "statement_coverage": {
     "covered": 86,
     "uncovered": 4,
     "coverable": 90,
     "total": 90,
     "percentage": 95.55
    },
    "statement_methods_fully_covered": {
     "covered": 13,
     "total": 16,
     "percentage": 81.25
    },
    "statement_methods_hit": {
     "covered": 15,
     "total": 16,
     "percentage": 93.75
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/status/evaluators",
   "name": "evaluators",
   "type": "folder",
   "path": "internal/status/evaluators",
   "parentId": "internal/status",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 19,
     "uncovered": 1,
     "coverable": 20,
     "total": 20,
     "percentage": 95
    },
    "patch_statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 19,
     "uncovered": 1,
     "coverable": 20,
     "total": 20,
     "percentage": 95
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/status/evaluators/registry.go",
   "name": "registry.go",
   "type": "file",
   "path": "internal/status/evaluators/registry.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 19,
     "uncovered": 1,
     "coverable": 20,
     "total": 20,
     "percentage": 95
    },
    "patch_statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 19,
     "uncovered": 1,
     "coverable": 20,
     "total": 20,
     "percentage": 95
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_registry.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/status/annotate.go",
   "name": "annotate.go",
   "type": "file",
   "path": "internal/status/annotate.go",
   "parentId": "internal/status",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 6,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 5,
     "uncovered": 0,
     "coverable": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 39,
     "uncovered": 3,
     "coverable": 42,
     "total": 42,
     "percentage": 92.85
    },
    "statement_methods_fully_covered": {
     "covered": 4,
     "total": 6,
     "percentage": 66.66
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_annotate.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/status/capabilities.go",
   "name": "capabilities.go",
   "type": "file",
   "path": "internal/status/capabilities.go",
   "parentId": "internal/status",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 13,
     "uncovered": 0,
     "coverable": 13,
     "total": 13,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 13,
     "uncovered": 0,
     "coverable": 13,
     "total": 13,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_capabilities.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/status/classifier.go",
   "name": "classifier.go",
   "type": "file",
   "path": "internal/status/classifier.go",
   "parentId": "internal/status",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 4
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 15,
     "uncovered": 0,
     "coverable": 15,
     "total": 15,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_classifier.go.html"
  },
  {
   "id": "internal/store",
   "name": "store",
   "type": "folder",
   "path": "internal/store",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 31
    },
    "methods_fully_covered": {
     "covered": 44,
     "total": 96,
     "percentage": 45.83
    },
    "methods_hit": {
     "covered": 90,
     "total": 96,
     "percentage": 93.75
    },
    "patch_methods_hit": {
     "covered": 90,
     "total": 96,
     "percentage": 93.75
    },
    "patch_statement_coverage": {
     "covered": 763,
     "uncovered": 169,
     "coverable": 932,
     "total": 932,
     "percentage": 81.86
    },
    "patch_statement_methods_hit": {
     "covered": 90,
     "total": 96,
     "percentage": 93.75
    },
    "statement_coverage": {
     "covered": 763,
     "uncovered": 169,
     "coverable": 932,
     "total": 932,
     "percentage": 81.86
    },
    "statement_methods_fully_covered": {
     "covered": 44,
     "total": 96,
     "percentage": 45.83
    },
    "statement_methods_hit": {
     "covered": 90,
     "total": 96,
     "percentage": 93.75
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/store/blob",
   "name": "blob",
   "type": "folder",
   "path": "internal/store/blob",
   "parentId": "internal/store",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 31
    },
    "methods_fully_covered": {
     "covered": 33,
     "total": 62,
     "percentage": 53.22
    },
    "methods_hit": {
     "covered": 59,
     "total": 62,
     "percentage": 95.16
    },
    "patch_methods_hit": {
     "covered": 59,
     "total": 62,
     "percentage": 95.16
    },
    "patch_statement_coverage": {
     "covered": 478,
     "uncovered": 88,
     "coverable": 566,
     "total": 566,
     "percentage": 84.45
    },
    "patch_statement_methods_hit": {
     "covered": 59,
     "total": 62,
     "percentage": 95.16
    },
    "statement_coverage": {
     "covered": 478,
     "uncovered": 88,
     "coverable": 566,
     "total": 566,
     "percentage": 84.45
    },
    "statement_methods_fully_covered": {
     "covered": 33,
     "total": 62,
     "percentage": 53.22
    },
    "statement_methods_hit": {
     "covered": 59,
     "total": 62,
     "percentage": 95.16
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/store/blob/analysis.go",
   "name": "analysis.go",
   "type": "file",
   "path": "internal/store/blob/analysis.go",
   "parentId": "internal/store/blob",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 31
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 4,
     "percentage": 25
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 94,
     "uncovered": 23,
     "coverable": 117,
     "total": 117,
     "percentage": 80.34
    },
    "patch_statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 94,
     "uncovered": 23,
     "coverable": 117,
     "total": 117,
     "percentage": 80.34
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 4,
     "percentage": 25
    },
    "statement_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_blob_analysis.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/store/blob/blob.go",
   "name": "blob.go",
   "type": "file",
   "path": "internal/store/blob/blob.go",
   "parentId": "internal/store/blob",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 20,
     "total": 31,
     "percentage": 64.51
    },
    "methods_hit": {
     "covered": 30,
     "total": 31,
     "percentage": 96.77
    },
    "patch_methods_hit": {
     "covered": 30,
     "total": 31,
     "percentage": 96.77
    },
    "patch_statement_coverage": {
     "covered": 103,
     "uncovered": 18,
     "coverable": 121,
     "total": 121,
     "percentage": 85.12
    },
    "patch_statement_methods_hit": {
     "covered": 30,
     "total": 31,
     "percentage": 96.77
    },
    "statement_coverage": {
     "covered": 103,
     "uncovered": 18,
     "coverable": 121,
     "total": 121,
     "percentage": 85.12
    },
    "statement_methods_fully_covered": {
     "covered": 20,
     "total": 31,
     "percentage": 64.51
    },
    "statement_methods_hit": {
     "covered": 30,
     "total": 31,
     "percentage": 96.77
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_blob_blob.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/store/blob/coverage.go",
   "name": "coverage.go",
   "type": "file",
   "path": "internal/store/blob/coverage.go",
   "parentId": "internal/store/blob",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 17
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 7,
     "percentage": 14.28
    },
    "methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 104,
     "uncovered": 17,
     "coverable": 121,
     "total": 121,
     "percentage": 85.95
    },
    "patch_statement_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 104,
     "uncovered": 17,
     "coverable": 121,
     "total": 121,
     "percentage": 85.95
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 7,
     "percentage": 14.28
    },
    "statement_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_blob_coverage.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/store/blob/diff.go",
   "name": "diff.go",
   "type": "file",
   "path": "internal/store/blob/diff.go",
   "parentId": "internal/store/blob",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 12
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 77,
     "uncovered": 14,
     "coverable": 91,
     "total": 91,
     "percentage": 84.61
    },
    "patch_statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 77,
     "uncovered": 14,
     "coverable": 91,
     "total": 91,
     "percentage": 84.61
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 5,
     "percentage": 40
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_blob_diff.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/store/blob/manifest.go",
   "name": "manifest.go",
   "type": "file",
   "path": "internal/store/blob/manifest.go",
   "parentId": "internal/store/blob",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 9,
     "total": 15,
     "percentage": 60
    },
    "methods_hit": {
     "covered": 13,
     "total": 15,
     "percentage": 86.66
    },
    "patch_methods_hit": {
     "covered": 13,
     "total": 15,
     "percentage": 86.66
    },
    "patch_statement_coverage": {
     "covered": 100,
     "uncovered": 16,
     "coverable": 116,
     "total": 116,
     "percentage": 86.2
    },
    "patch_statement_methods_hit": {
     "covered": 13,
     "total": 15,
     "percentage": 86.66
    },
    "statement_coverage": {
     "covered": 100,
     "uncovered": 16,
     "coverable": 116,
     "total": 116,
     "percentage": 86.2
    },
    "statement_methods_fully_covered": {
     "covered": 9,
     "total": 15,
     "percentage": 60
    },
    "statement_methods_hit": {
     "covered": 13,
     "total": 15,
     "percentage": 86.66
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_blob_manifest.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/store/nanovision.yaml",
   "name": "nanovision.yaml",
   "type": "file",
   "path": "internal/store/nanovision.yaml",
   "parentId": "internal/store",
   "depth": 2,
   "config": true
  },
  {
   "id": "internal/store/blobs.go",
   "name": "blobs.go",
   "type": "file",
   "path": "internal/store/blobs.go",
   "parentId": "internal/store",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 9
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 10,
     "percentage": 20
    },
    "methods_hit": {
     "covered": 8,
     "total": 10,
     "percentage": 80
    },
    "patch_methods_hit": {
     "covered": 8,
     "total": 10,
     "percentage": 80
    },
    "patch_statement_coverage": {
     "covered": 68,
     "uncovered": 34,
     "coverable": 102,
     "total": 102,
     "percentage": 66.66
    },
    "patch_statement_methods_hit": {
     "covered": 8,
     "total": 10,
     "percentage": 80
    },
    "statement_coverage": {
     "covered": 68,
     "uncovered": 34,
     "coverable": 102,
     "total": 102,
     "percentage": 66.66
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 10,
     "percentage": 20
    },
    "statement_methods_hit": {
     "covered": 8,
     "total": 10,
     "percentage": 80
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_store_blobs.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/store/maintenance.go",
   "name": "maintenance.go",
   "type": "file",
   "path": "internal/store/maintenance.go",
   "parentId": "internal/store",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 7,
     "percentage": 28.57
    },
    "methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 79,
     "uncovered": 19,
     "coverable": 98,
     "total": 98,
     "percentage": 80.61
    },
    "patch_statement_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 79,
     "uncovered": 19,
     "coverable": 98,
     "total": 98,
     "percentage": 80.61
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 7,
     "percentage": 28.57
    },
    "statement_methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_maintenance.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/store/runs.go",
   "name": "runs.go",
   "type": "file",
   "path": "internal/store/runs.go",
   "parentId": "internal/store",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 11,
     "percentage": 54.54
    },
    "methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 101,
     "uncovered": 12,
     "coverable": 113,
     "total": 113,
     "percentage": 89.38
    },
    "patch_statement_methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 101,
     "uncovered": 12,
     "coverable": 113,
     "total": 113,
     "percentage": 89.38
    },
    "statement_methods_fully_covered": {
     "covered": 6,
     "total": 11,
     "percentage": 54.54
    },
    "statement_methods_hit": {
     "covered": 11,
     "total": 11,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_runs.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/store/store.go",
   "name": "store.go",
   "type": "file",
   "path": "internal/store/store.go",
   "parentId": "internal/store",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 8
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 37,
     "uncovered": 16,
     "coverable": 53,
     "total": 53,
     "percentage": 69.81
    },
    "patch_statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 37,
     "uncovered": 16,
     "coverable": 53,
     "total": 53,
     "percentage": 69.81
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 6,
     "percentage": 16.66
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_store_store.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/tree",
   "name": "tree",
   "type": "folder",
   "path": "internal/tree",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 2,
     "uncovered": 0,
     "coverable": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 79,
     "uncovered": 5,
     "coverable": 84,
     "total": 84,
     "percentage": 94.04
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/tree/builder.go",
   "name": "builder.go",
   "type": "file",
   "path": "internal/tree/builder.go",
   "parentId": "internal/tree",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 13
    },
    "methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 2,
     "uncovered": 0,
     "coverable": 2,
     "total": 2,
     "percentage": 100
    },
    "patch_statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 79,
     "uncovered": 5,
     "coverable": 84,
     "total": 84,
     "percentage": 94.04
    },
    "statement_methods_fully_covered": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_tree_builder.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/utils",
   "name": "utils",
   "type": "folder",
   "path": "internal/utils",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 2,
     "total": 6,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 81,
     "uncovered": 13,
     "coverable": 94,
     "total": 94,
     "percentage": 86.17
    },
    "statement_methods_fully_covered": {
     "covered": 2,
     "total": 6,
     "percentage": 33.33
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/utils/analyzer.go",
   "name": "analyzer.go",
   "type": "file",
   "path": "internal/utils/analyzer.go",
   "parentId": "internal/utils",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 1,
     "uncovered": 0,
     "coverable": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_utils_analyzer.go.html"
  },
  {
   "id": "internal/utils/line_sorter.go",
   "name": "line_sorter.go",
   "type": "file",
   "path": "internal/utils/line_sorter.go",
   "parentId": "internal/utils",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 12,
     "uncovered": 0,
     "coverable": 12,
     "total": 12,
     "percentage": 100
    },
    "statement_methods_fully_covered": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 1,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_utils_line_sorter.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/utils/math.go",
   "name": "math.go",
   "type": "file",
   "path": "internal/utils/math.go",
   "parentId": "internal/utils",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    },
    "statement_coverage": {
     "covered": 9,
     "uncovered": 9,
     "coverable": 18,
     "total": 18,
     "percentage": 50
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 1,
     "total": 2,
     "percentage": 50
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_utils_math.go.html"
  },
  {
   "id": "internal/utils/paths.go",
   "name": "paths.go",
   "type": "file",
   "path": "internal/utils/paths.go",
   "parentId": "internal/utils",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 59,
     "uncovered": 4,
     "coverable": 63,
     "total": 63,
     "percentage": 93.65
    },
    "statement_methods_fully_covered": {
     "covered": 0,
     "total": 2,
     "percentage": 0
    },
    "statement_methods_hit": {
     "covered": 2,
     "total": 2,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_utils_paths.go.html"
  },
  {
   "id": "internal/vcs",
   "name": "vcs",
   "type": "folder",
   "path": "internal/vcs",
   "parentId": "internal",
   "depth": 1,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 10,
     "total": 32,
     "percentage": 31.25
    },
    "methods_hit": {
     "covered": 30,
     "total": 32,
     "percentage": 93.75
    },
    "patch_methods_hit": {
     "covered": 30,
     "total": 32,
     "percentage": 93.75
    },
    "patch_statement_coverage": {
     "covered": 198,
     "uncovered": 35,
     "coverable": 233,
     "total": 233,
     "percentage": 84.97
    },
    "patch_statement_methods_hit": {
     "covered": 30,
     "total": 32,
     "percentage": 93.75
    },
    "statement_coverage": {
     "covered": 198,
     "uncovered": 35,
     "coverable": 233,
     "total": 233,
     "percentage": 84.97
    },
    "statement_methods_fully_covered": {
     "covered": 10,
     "total": 32,
     "percentage": 31.25
    },
    "statement_methods_hit": {
     "covered": 30,
     "total": 32,
     "percentage": 93.75
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/vcs/git.go",
   "name": "git.go",
   "type": "file",
   "path": "internal/vcs/git.go",
   "parentId": "internal/vcs",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 7
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 12,
     "percentage": 25
    },
    "methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 65,
     "uncovered": 14,
     "coverable": 79,
     "total": 79,
     "percentage": 82.27
    },
    "patch_statement_methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 65,
     "uncovered": 14,
     "coverable": 79,
     "total": 79,
     "percentage": 82.27
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 12,
     "percentage": 25
    },
    "statement_methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_vcs_git.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/vcs/perforce.go",
   "name": "perforce.go",
   "type": "file",
   "path": "internal/vcs/perforce.go",
   "parentId": "internal/vcs",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 14
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 15,
     "percentage": 26.66
    },
    "methods_hit": {
     "covered": 13,
     "total": 15,
     "percentage": 86.66
    },
    "patch_methods_hit": {
     "covered": 13,
     "total": 15,
     "percentage": 86.66
    },
    "patch_statement_coverage": {
     "covered": 106,
     "uncovered": 19,
     "coverable": 125,
     "total": 125,
     "percentage": 84.8
    },
    "patch_statement_methods_hit": {
     "covered": 13,
     "total": 15,
     "percentage": 86.66
    },
    "statement_coverage": {
     "covered": 106,
     "uncovered": 19,
     "coverable": 125,
     "total": 125,
     "percentage": 84.8
    },
    "statement_methods_fully_covered": {
     "covered": 4,
     "total": 15,
     "percentage": 26.66
    },
    "statement_methods_hit": {
     "covered": 13,
     "total": 15,
     "percentage": 86.66
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_vcs_perforce.go.html",
   "diffStatus": "added"
  },
  {
   "id": "internal/vcs/vcs.go",
   "name": "vcs.go",
   "type": "file",
   "path": "internal/vcs/vcs.go",
   "parentId": "internal/vcs",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 5
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 27,
     "uncovered": 2,
     "coverable": 29,
     "total": 29,
     "percentage": 93.1
    },
    "patch_statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 27,
     "uncovered": 2,
     "coverable": 29,
     "total": 29,
     "percentage": 93.1
    },
    "statement_methods_fully_covered": {
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "statement_methods_hit": {
     "covered": 5,
     "total": 5,
     "percentage": 100
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_vcs_vcs.go.html",
   "diffStatus": "added"
  },
  {
   "id": "nanovision.yaml",
   "name": "nanovision.yaml",
   "type": "file",
   "path": "nanovision.yaml",
   "depth": 0,
   "config": true
  }
 ],
 "metricDefinitions": {
  "a_statement_coverage": {
   "label": "Statements",
   "shortLabel": "Statements",
   "description": "Statements run by tests.",
   "subMetrics": [
    {
     "id": "value",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "b_complexity": {
   "label": "Cyclomatic Complexity",
   "shortLabel": "Complexity",
   "description": "Paths through the method. Lower is better.",
   "kind": "value",
   "subMetrics": [
    {
     "id": "value",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "c_patch_statement_coverage": {
   "label": "Patch Statements",
   "shortLabel": "Patch Stmts",
   "description": "Changed statements run by tests. Needs a diff.",
   "subMetrics": [
    {
     "id": "value",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "d_crap_score": {
   "label": "CRAP Score",
   "shortLabel": "CRAP",
   "description": "Complexity weighted by missing coverage. Lower is better.",
   "kind": "value",
   "subMetrics": [
    {
     "id": "value",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "max_cyclomatic_complexity": {
   "label": "Max Cyclomatic Complexity",
   "shortLabel": "Max Complexity",
   "description": "Highest method complexity. Lower is better.",
   "kind": "value",
   "subMetrics": [
    {
     "id": "value",
     "label": "Value",
     "width": 140
    }
   ]
  },
  "methods_fully_covered": {
   "label": "Methods Fully Covered",
   "shortLabel": "Fully Covered",
   "description": "Methods with every line run.",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Covered",
     "width": 80
    },
    {
     "id": "total",
     "label": "Total",
     "width": 80
    },
    {
     "id": "percentage",
     "label": "Percentage %",
     "width": 160
    }
   ]
  },
  "methods_hit": {
   "label": "Methods Hit",
   "shortLabel": "Methods Hit",
   "description": "Methods entered by tests.",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Hit",
     "width": 80
    },
    {
     "id": "total",
     "label": "Total",
     "width": 80
    },
    {
     "id": "percentage",
     "label": "Percentage %",
     "width": 160
    }
   ]
  },
  "patch_methods_hit": {
   "label": "Patch Methods Hit",
   "shortLabel": "Patch Methods Hit",
   "description": "Changed methods entered by tests. Needs a diff.",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Hit",
     "width": 80
    },
    {
     "id": "total",
     "label": "Total",
     "width": 80
    },
    {
     "id": "percentage",
     "label": "Percentage %",
     "width": 160
    }
   ]
  },
  "patch_statement_coverage": {
   "label": "Patch Statements",
   "shortLabel": "Patch Statements",
   "description": "Changed statements run by tests. Needs a diff.",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Covered",
     "width": 100
    },
    {
     "id": "uncovered",
     "label": "Uncovered",
     "width": 100
    },
    {
     "id": "total",
     "label": "Total",
     "width": 80
    },
    {
     "id": "percentage",
     "label": "Percentage %",
     "width": 160
    }
   ]
  },
  "patch_statement_methods_hit": {
   "label": "Patch Statement Methods Hit",
   "shortLabel": "Patch Stmt Methods Hit",
   "description": "Changed methods with a statement run. Needs a diff.",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Hit",
     "width": 80
    },
    {
     "id": "total",
     "label": "Total",
     "width": 80
    },
    {
     "id": "percentage",
     "label": "Percentage %",
     "width": 160
    }
   ]
  },
  "statement_coverage": {
   "label": "Statements",
   "shortLabel": "Statements",
   "description": "Statements run by tests.",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Covered",
     "width": 100
    },
    {
     "id": "uncovered",
     "label": "Uncovered",
     "width": 100
    },
    {
     "id": "total",
     "label": "Total",
     "width": 80
    },
    {
     "id": "percentage",
     "label": "Percentage %",
     "width": 160
    }
   ]
  },
  "statement_methods_fully_covered": {
   "label": "Statement Methods Fully Covered",
   "shortLabel": "Stmt Fully Covered",
   "description": "Methods with every statement run.",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Covered",
     "width": 80
    },
    {
     "id": "total",
     "label": "Total",
     "width": 80
    },
    {
     "id": "percentage",
     "label": "Percentage %",
     "width": 160
    }
   ]
  },
  "statement_methods_hit": {
   "label": "Statement Methods Hit",
   "shortLabel": "Stmt Methods Hit",
   "description": "Methods with a statement run.",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Hit",
     "width": 80
    },
    {
     "id": "total",
     "label": "Total",
     "width": 80
    },
    {
     "id": "percentage",
     "label": "Percentage %",
     "width": 160
    }
   ]
  }
 },
 "metricOrder": [
  "statement_coverage",
  "methods_hit",
  "methods_fully_covered",
  "patch_methods_hit",
  "patch_statement_coverage",
  "statement_methods_hit",
  "statement_methods_fully_covered",
  "patch_statement_methods_hit",
  "max_cyclomatic_complexity"
 ],
 "metadata": [
  {
   "label": "Generated At",
   "value": "2026-10-03 16:49:35"
  },
  {
   "label": "Parser",
   "value": "Cobertura | GCov | GoCover"
  }
 ],
 "review": {
  "passed": false,
  "checks": [
   {
    "key": "patch_statement_coverage",
    "label": "Patch statement coverage",
    "value": 79.74573822594625,
    "threshold": 80,
    "passed": false
   },
   {
    "key": "max_changed_method_complexity",
    "label": "Max changed-method complexity",
    "value": 31,
    "threshold": 15,
    "passed": false
   }
  ],
  "stats": {
   "changedFiles": 65,
   "methodsAdded": 294,
   "methodsModified": 62,
   "untestedChangedMethods": 25,
   "patchStatementsValid": 3461,
   "patchStatementsCovered": 2760,
   "maxChangedComplexity": 31
  },
  "hotspots": [
   {
    "file": "cmd/serve.go",
    "method": "runServe",
    "startLine": 26,
    "diffStatus": "added",
    "complexity": 16,
    "patchCoverage": 0,
    "risk": 16
   },
   {
    "file": "cmd/main.go",
    "method": "main",
    "startLine": 215,
    "diffStatus": "modified",
    "complexity": 27,
    "patchCoverage": 67.6470588235294,
    "risk": 11.903225806451614
   },
   {
    "file": "cmd/storecmd.go",
    "method": "resolveStoreDir",
    "startLine": 84,
    "diffStatus": "added",
    "complexity": 9,
    "patchCoverage": 10.526315789473685,
    "risk": 8.052631578947368
   },
   {
    "file": "internal/server/maintenance.go",
    "method": "(*Server).maintain",
    "startLine": 33,
    "diffStatus": "added",
    "complexity": 8,
    "patchCoverage": 0,
    "risk": 8
   },
   {
    "file": "internal/config/config.go",
    "method": "(*AppConfig).mergeCliOverrides",
    "startLine": 428,
    "diffStatus": "modified",
    "complexity": 28,
    "patchCoverage": 86.20689655172414,
    "risk": 7.924528301886793
   },
   {
    "file": "internal/store/blob/analysis.go",
    "method": "DecodeAnalysis",
    "startLine": 71,
    "diffStatus": "added",
    "complexity": 31,
    "patchCoverage": 75,
    "risk": 7.75
   },
   {
    "file": "internal/reporter/htmlreact/builder.go",
    "method": "(*HtmlReactReportBuilder).createSingleFileReport",
    "startLine": 77,
    "diffStatus": "modified",
    "complexity": 7,
    "patchCoverage": 0,
    "risk": 7
   },
   {
    "file": "internal/server/ui.go",
    "method": "(*Server).handleUI",
    "startLine": 44,
    "diffStatus": "added",
    "complexity": 12,
    "patchCoverage": 44.44444444444444,
    "risk": 6.666666666666667
   },
   {
    "file": "cmd/storecmd.go",
    "method": "runStoreCommand",
    "startLine": 28,
    "diffStatus": "added",
    "complexity": 13,
    "patchCoverage": 53.65853658536585,
    "risk": 6.024390243902438
   },
   {
    "file": "internal/reporter/textsummary/terminal.go",
    "method": "WriteTerminal",
    "startLine": 34,
    "diffStatus": "added",
    "complexity": 29,
    "patchCoverage": 80,
    "risk": 5.799999999999999
   }
  ]
 },
 "comparing": [
  {
   "label": "Current",
   "value": "working copy"
  },
  {
   "label": "Changed files",
   "value": "317 from fixture.diff"
  }
 ],
 "reports": [
  {
   "name": "unit tests",
   "path": "reports/nanovision_self_coverage/coverage-unit.out"
  },
  {
   "name": "e2e tests",
   "path": "reports/nanovision_self_coverage/coverage-integration.out"
  },
  {
   "name": "cpp demo (gcov)",
   "path": "demo_projects/cpp/report/gcov/branch-probabilities/*.gcov"
  },
  {
   "name": "csharp demo (cobertura)",
   "path": "demo_projects/csharp/report/cobertura/cobertura.xml"
  },
  {
   "name": "go demo (gocover)",
   "path": "demo_projects/go/report/gocover/coverage.out"
  }
 ],
 "reportIndexes": {
  "cmd/configcmd.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 83
    },
    {
     "m": [
      1
     ],
     "n": 28
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ]
  },
  "cmd/history.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 11
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 12
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 12
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 118
    },
    {
     "m": [
      0
     ],
     "n": 38
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 12
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 118
    },
    {
     "m": [
      0
     ],
     "n": 38
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 11
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 12
    }
   ]
  },
  "cmd/main.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 9
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 6
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 42
    },
    {
     "m": [
      3
     ],
     "n": 18
    },
    {
     "m": [
      0
     ],
     "n": 16
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 6
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 87
    },
    {
     "m": [
      0
     ],
     "n": 62
    },
    {
     "m": [
      3
     ],
     "n": 35
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 9
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ]
  },
  "cmd/serve.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 90
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 90
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 3
    }
   ]
  },
  "cmd/storecmd.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 36
    },
    {
     "m": [
      2
     ],
     "n": 24
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 36
    },
    {
     "m": [
      2
     ],
     "n": 24
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ]
  },
  "demo_projects/cpp/project/src/advanced_calculator.cpp": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      4
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      4
     ],
     "n": 12
    },
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      4
     ],
     "n": 2
    }
   ]
  },
  "demo_projects/cpp/project/src/calculator.cpp": {
   "methods_fully_covered": [
    {
     "m": [
      4
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      4
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      4
     ],
     "n": 9
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      4
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      4
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "demo_projects/cpp/project/src/utils/math_utils.cpp": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      4
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      4
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      4
     ],
     "n": 2
    }
   ]
  },
  "demo_projects/go/project/calculator/calculator.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      16
     ],
     "n": 3
    }
   ],
   "methods_hit": [
    {
     "m": [
      16
     ],
     "n": 6
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      16
     ],
     "n": 20
    },
    {
     "m": [
      0
     ],
     "n": 6
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      16
     ],
     "n": 3
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      16
     ],
     "n": 6
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "demo_projects/go/project/calculator/entities.go": {
   "methods_fully_covered": [
    {
     "m": [
      16
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      16
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      16
     ],
     "n": 9
    },
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      16
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      16
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "demo_projects/go/project/calculator_2/calculator.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      16
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      16
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      16
     ],
     "n": 10
    },
    {
     "m": [
      0
     ],
     "n": 5
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      16
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      16
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "demo_projects/go/project/calculator_2/entities.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      16
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      16
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      16
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      16
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      16
     ],
     "n": 1
    }
   ]
  },
  "internal/aggregator/aggrgator.go": {
   "methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      1,
      3
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 13
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 143
    },
    {
     "m": [
      1
     ],
     "n": 17
    },
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      1,
      3
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 13
    }
   ]
  },
  "internal/aggregator/report_index.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 7
    },
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      1,
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 17
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 124
    },
    {
     "m": [
      1
     ],
     "n": 41
    },
    {
     "m": [
      0
     ],
     "n": 11
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 7
    },
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      1,
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 17
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/analyzer/cpp/analyzer.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 72
    },
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 9
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ]
  },
  "internal/analyzer/gdscript/analyzer.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 58
    },
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ]
  },
  "internal/analyzer/go/analyzer.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 65
    },
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 6
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ]
  },
  "internal/bootlog/bootlog.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 29
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ]
  },
  "internal/cache/cache.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1,
      3
     ],
     "n": 3
    },
    {
     "m": [
      1,
      2,
      3
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 65
    },
    {
     "m": [
      1
     ],
     "n": 15
    },
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1,
      3
     ],
     "n": 3
    },
    {
     "m": [
      1,
      2,
      3
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    }
   ]
  },
  "internal/cache/cache_validator.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/calculator/engine.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 20
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 40
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ]
  },
  "internal/calculator/registry.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 31
    },
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 31
    },
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ]
  },
  "internal/client/client.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 7
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 10
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 10
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 50
    },
    {
     "m": [
      0
     ],
     "n": 16
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 10
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 50
    },
    {
     "m": [
      0
     ],
     "n": 16
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 7
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 10
    }
   ]
  },
  "internal/compare/compare.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      1,
      3
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 9
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 9
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 116
    },
    {
     "m": [
      1
     ],
     "n": 20
    },
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 9
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 116
    },
    {
     "m": [
      1
     ],
     "n": 20
    },
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      1,
      3
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 9
    }
   ]
  },
  "internal/config/config.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 15
    },
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      1,
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 23
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 23
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 207
    },
    {
     "m": [
      1
     ],
     "n": 60
    },
    {
     "m": [
      0
     ],
     "n": 25
    },
    {
     "m": [
      2
     ],
     "n": 10
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 23
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 248
    },
    {
     "m": [
      1
     ],
     "n": 62
    },
    {
     "m": [
      0
     ],
     "n": 45
    },
    {
     "m": [
      2
     ],
     "n": 15
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 15
    },
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      1,
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 23
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/config/metrics.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 10
    },
    {
     "m": [
      3
     ],
     "n": 7
    },
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 5
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/config/outputs.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    }
   ]
  },
  "internal/config/schema.go": {
   "methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 71
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ]
  },
  "internal/config/scoped.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 53
    },
    {
     "m": [
      1
     ],
     "n": 33
    },
    {
     "m": [
      0
     ],
     "n": 12
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/diagnostics/diagnostics.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 7
    },
    {
     "m": [
      2
     ],
     "n": 4
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 10
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 4
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 55
    },
    {
     "m": [
      0
     ],
     "n": 21
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 7
    },
    {
     "m": [
      2
     ],
     "n": 4
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 10
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/diff/parser.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 2
    },
    {
     "m": [
      1,
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 13
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 7
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 95
    },
    {
     "m": [
      1
     ],
     "n": 25
    },
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 2
    },
    {
     "m": [
      1,
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 13
    }
   ]
  },
  "internal/diff/path.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 9
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ]
  },
  "internal/diffapply/apply.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 17
    },
    {
     "m": [
      1
     ],
     "n": 6
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 40
    },
    {
     "m": [
      1
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    }
   ]
  },
  "internal/diffapply/resolver.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 26
    },
    {
     "m": [
      0
     ],
     "n": 7
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 97
    },
    {
     "m": [
      1
     ],
     "n": 49
    },
    {
     "m": [
      0
     ],
     "n": 14
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ]
  },
  "internal/enricher/enricher.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 3
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 54
    },
    {
     "m": [
      3
     ],
     "n": 43
    },
    {
     "m": [
      0
     ],
     "n": 12
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 3
    }
   ]
  },
  "internal/filereader/default_reader.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/filereader/filereader.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 16
    },
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 36
    },
    {
     "m": [
      3
     ],
     "n": 30
    },
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ]
  },
  "internal/filtering/filter.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 34
    },
    {
     "m": [
      1
     ],
     "n": 11
    },
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/history/capture.go": {
   "methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 61
    },
    {
     "m": [
      1
     ],
     "n": 6
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 61
    },
    {
     "m": [
      1
     ],
     "n": 6
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/history/meta.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 29
    },
    {
     "m": [
      3
     ],
     "n": 16
    },
    {
     "m": [
      0
     ],
     "n": 11
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 29
    },
    {
     "m": [
      3
     ],
     "n": 16
    },
    {
     "m": [
      0
     ],
     "n": 11
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ]
  },
  "internal/history/rebuild.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 72
    },
    {
     "m": [
      0
     ],
     "n": 11
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 72
    },
    {
     "m": [
      0
     ],
     "n": 11
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ]
  },
  "internal/history/record.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 77
    },
    {
     "m": [
      1
     ],
     "n": 18
    },
    {
     "m": [
      0
     ],
     "n": 16
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 77
    },
    {
     "m": [
      1
     ],
     "n": 18
    },
    {
     "m": [
      0
     ],
     "n": 16
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/logging/logging.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 23
    },
    {
     "m": [
      3
     ],
     "n": 17
    },
    {
     "m": [
      1
     ],
     "n": 14
    },
    {
     "m": [
      0
     ],
     "n": 9
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ]
  },
  "internal/model/comparison.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    }
   ]
  },
  "internal/model/diff.go": {
   "methods_fully_covered": [
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 9
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ]
  },
  "internal/parsers/factory.go": {
   "methods_fully_covered": [
    {
     "m": [
      1,
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 7
    },
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      1,
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/parsers/parser_cobertura/parser.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 24
    },
    {
     "m": [
      0
     ],
     "n": 12
    },
    {
     "m": [
      2
     ],
     "n": 12
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/parsers/parser_cobertura/processing.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 30
    },
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    }
   ]
  },
  "internal/parsers/parser_config.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/parsers/parser_gcov/parser.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      2
     ],
     "n": 7
    },
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/parsers/parser_gcov/processing.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 23
    },
    {
     "m": [
      1
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    }
   ]
  },
  "internal/parsers/parser_gocover/parser.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 27
    },
    {
     "m": [
      2
     ],
     "n": 7
    },
    {
     "m": [
      0
     ],
     "n": 6
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/parsers/parser_gocover/processing.go": {
   "methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 26
    },
    {
     "m": [
      1
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    }
   ]
  },
  "internal/parsers/parser_lcov/parser.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 15
    },
    {
     "m": [
      1
     ],
     "n": 7
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ]
  },
  "internal/parsers/parser_lcov/processing.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 41
    },
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    }
   ]
  },
  "internal/pipeline/pipeline.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 5
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 5
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 67
    },
    {
     "m": [
      0
     ],
     "n": 19
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 5
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 67
    },
    {
     "m": [
      0
     ],
     "n": 19
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 5
    }
   ]
  },
  "internal/reporter/annotations/reporter.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 5
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 32
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 5
    }
   ]
  },
  "internal/reporter/htmlreact/builder.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 15
    },
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      1,
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 21
    },
    {
     "m": [
      2
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 18
    },
    {
     "m": [
      2
     ],
     "n": 7
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 69
    },
    {
     "m": [
      1
     ],
     "n": 26
    },
    {
     "m": [
      2
     ],
     "n": 26
    },
    {
     "m": [
      0
     ],
     "n": 11
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 18
    },
    {
     "m": [
      2
     ],
     "n": 7
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 119
    },
    {
     "m": [
      2
     ],
     "n": 98
    },
    {
     "m": [
      0
     ],
     "n": 53
    },
    {
     "m": [
      1
     ],
     "n": 27
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 15
    },
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      1,
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 21
    },
    {
     "m": [
      2
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ]
  },
  "internal/reporter/htmlreact/details_generator.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      1,
      3
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 10
    },
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 18
    },
    {
     "m": [
      0
     ],
     "n": 8
    },
    {
     "m": [
      3
     ],
     "n": 6
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 89
    },
    {
     "m": [
      3
     ],
     "n": 33
    },
    {
     "m": [
      0
     ],
     "n": 31
    },
    {
     "m": [
      1
     ],
     "n": 5
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      1,
      3
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 10
    },
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/reporter/htmlreact/embed.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/reporter/htmlreact/emit.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 14
    },
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/reporter/htmlreact/generator.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 34
    },
    {
     "m": [
      0
     ],
     "n": 16
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ]
  },
  "internal/reporter/htmlreact/generator_single.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 68
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ]
  },
  "internal/reporter/htmlreact/schema.go": {
   "methods_fully_covered": [
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ]
  },
  "internal/reporter/htmlreact/views.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 15
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 15
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 4
    }
   ]
  },
  "internal/reporter/lcov/reporter.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 52
    },
    {
     "m": [
      0
     ],
     "n": 12
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/reporter/reporter_rawjson/reporter.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/reporter/sarif/reporter.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 6
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 6
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 33
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 6
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 6
    }
   ]
  },
  "internal/reporter/textsummary/comparison.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 8
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 8
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 65
    },
    {
     "m": [
      0
     ],
     "n": 10
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 8
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 65
    },
    {
     "m": [
      0
     ],
     "n": 10
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 8
    }
   ]
  },
  "internal/reporter/textsummary/reporter.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 3
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 35
    },
    {
     "m": [
      3
     ],
     "n": 25
    },
    {
     "m": [
      2
     ],
     "n": 4
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 3
    }
   ]
  },
  "internal/reporter/textsummary/terminal.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 57
    },
    {
     "m": [
      0
     ],
     "n": 13
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 57
    },
    {
     "m": [
      0
     ],
     "n": 13
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    }
   ]
  },
  "internal/review/review.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 49
    },
    {
     "m": [
      0
     ],
     "n": 12
    },
    {
     "m": [
      1
     ],
     "n": 12
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ]
  },
  "internal/server/api.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 12
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 129
    },
    {
     "m": [
      0
     ],
     "n": 43
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 129
    },
    {
     "m": [
      0
     ],
     "n": 43
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 12
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/server/cache.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 4
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 4
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 38
    },
    {
     "m": [
      0
     ],
     "n": 9
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 4
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 38
    },
    {
     "m": [
      0
     ],
     "n": 9
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 4
    }
   ]
  },
  "internal/server/maintenance.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 24
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 24
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ]
  },
  "internal/server/server.go": {
   "methods_fully_covered": [
    {
     "m": [
      1
     ],
     "n": 7
    },
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 11
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 11
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 55
    },
    {
     "m": [
      0
     ],
     "n": 7
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 11
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 55
    },
    {
     "m": [
      0
     ],
     "n": 7
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      1
     ],
     "n": 7
    },
    {
     "m": [
      0
     ],
     "n": 4
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 11
    }
   ]
  },
  "internal/server/ui.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 24
    },
    {
     "m": [
      0
     ],
     "n": 17
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 24
    },
    {
     "m": [
      0
     ],
     "n": 17
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 2
    }
   ]
  },
  "internal/server/upload.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 34
    },
    {
     "m": [
      0
     ],
     "n": 14
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 34
    },
    {
     "m": [
      0
     ],
     "n": 14
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 3
    }
   ]
  },
  "internal/server/wire.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 20
    },
    {
     "m": [
      1
     ],
     "n": 15
    },
    {
     "m": [
      0
     ],
     "n": 11
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 20
    },
    {
     "m": [
      1
     ],
     "n": 15
    },
    {
     "m": [
      0
     ],
     "n": 11
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/status/annotate.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 5
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 35
    },
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 6
    }
   ]
  },
  "internal/status/capabilities.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 13
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 13
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/status/classifier.go": {
   "methods_fully_covered": [
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 9
    },
    {
     "m": [
      1
     ],
     "n": 6
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/status/evaluators/registry.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 14
    },
    {
     "m": [
      2
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 14
    },
    {
     "m": [
      2
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/store/blob/analysis.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 92
    },
    {
     "m": [
      0
     ],
     "n": 23
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 92
    },
    {
     "m": [
      0
     ],
     "n": 23
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    }
   ]
  },
  "internal/store/blob/blob.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 11
    },
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      1,
      3
     ],
     "n": 6
    },
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 24
    },
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 24
    },
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 78
    },
    {
     "m": [
      0
     ],
     "n": 18
    },
    {
     "m": [
      1
     ],
     "n": 17
    },
    {
     "m": [
      2
     ],
     "n": 8
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 24
    },
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 78
    },
    {
     "m": [
      0
     ],
     "n": 18
    },
    {
     "m": [
      1
     ],
     "n": 17
    },
    {
     "m": [
      2
     ],
     "n": 8
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 11
    },
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      1,
      3
     ],
     "n": 6
    },
    {
     "m": [
      2
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 24
    },
    {
     "m": [
      2
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/store/blob/coverage.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 6
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 65
    },
    {
     "m": [
      1
     ],
     "n": 39
    },
    {
     "m": [
      0
     ],
     "n": 17
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 65
    },
    {
     "m": [
      1
     ],
     "n": 39
    },
    {
     "m": [
      0
     ],
     "n": 17
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 6
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ]
  },
  "internal/store/blob/diff.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 72
    },
    {
     "m": [
      0
     ],
     "n": 14
    },
    {
     "m": [
      1
     ],
     "n": 5
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 72
    },
    {
     "m": [
      0
     ],
     "n": 14
    },
    {
     "m": [
      1
     ],
     "n": 5
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 3
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/store/blob/manifest.go": {
   "methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 6
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 87
    },
    {
     "m": [
      0
     ],
     "n": 16
    },
    {
     "m": [
      1
     ],
     "n": 10
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 87
    },
    {
     "m": [
      0
     ],
     "n": 16
    },
    {
     "m": [
      1
     ],
     "n": 10
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 6
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/store/blobs.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 8
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 64
    },
    {
     "m": [
      0
     ],
     "n": 34
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 64
    },
    {
     "m": [
      0
     ],
     "n": 34
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 8
    },
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ]
  },
  "internal/store/maintenance.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 61
    },
    {
     "m": [
      0
     ],
     "n": 19
    },
    {
     "m": [
      3
     ],
     "n": 17
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 61
    },
    {
     "m": [
      0
     ],
     "n": 19
    },
    {
     "m": [
      3
     ],
     "n": 17
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 4
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/store/runs.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      1,
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 9
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 9
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 75
    },
    {
     "m": [
      1
     ],
     "n": 23
    },
    {
     "m": [
      0
     ],
     "n": 12
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 9
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 75
    },
    {
     "m": [
      1
     ],
     "n": 23
    },
    {
     "m": [
      0
     ],
     "n": 12
    },
    {
     "m": [
      2
     ],
     "n": 3
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      1,
      3
     ],
     "n": 2
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 1
    },
    {
     "m": [
      1,
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 9
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ]
  },
  "internal/store/store.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 36
    },
    {
     "m": [
      0
     ],
     "n": 16
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 36
    },
    {
     "m": [
      0
     ],
     "n": 16
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ]
  },
  "internal/tree/builder.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 6
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 79
    },
    {
     "m": [
      0
     ],
     "n": 5
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 5
    },
    {
     "m": [
      0
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 6
    }
   ]
  },
  "internal/utils/analyzer.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/utils/line_sorter.go": {
   "methods_fully_covered": [
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 12
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/utils/math.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      2
     ],
     "n": 9
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      0
     ],
     "n": 1
    },
    {
     "m": [
      2
     ],
     "n": 1
    }
   ]
  },
  "internal/utils/paths.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 58
    },
    {
     "m": [
      0
     ],
     "n": 4
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 2
    }
   ]
  },
  "internal/vcs/git.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 44
    },
    {
     "m": [
      1
     ],
     "n": 19
    },
    {
     "m": [
      0
     ],
     "n": 14
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 44
    },
    {
     "m": [
      1
     ],
     "n": 19
    },
    {
     "m": [
      0
     ],
     "n": 14
    },
    {
     "m": [
      2
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 9
    },
    {
     "m": [
      3
     ],
     "n": 2
    },
    {
     "m": [
      2,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 11
    },
    {
     "m": [
      1
     ],
     "n": 1
    }
   ]
  },
  "internal/vcs/perforce.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 11
    },
    {
     "m": [
      1
     ],
     "n": 4
    }
   ],
   "methods_hit": [
    {
     "m": [
      1
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 106
    },
    {
     "m": [
      0
     ],
     "n": 19
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      1
     ],
     "n": 106
    },
    {
     "m": [
      0
     ],
     "n": 19
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 11
    },
    {
     "m": [
      1
     ],
     "n": 4
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      1
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ]
  },
  "internal/vcs/vcs.go": {
   "methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 18
    },
    {
     "m": [
      1
     ],
     "n": 9
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "patch_statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ],
   "statement_coverage": [
    {
     "m": [
      3
     ],
     "n": 18
    },
    {
     "m": [
      1
     ],
     "n": 9
    },
    {
     "m": [
      0
     ],
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 2
    },
    {
     "m": [
      1
     ],
     "n": 2
    },
    {
     "m": [
      1,
      3
     ],
     "n": 1
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      3
     ],
     "n": 3
    },
    {
     "m": [
      1
     ],
     "n": 2
    }
   ]
  }
 },
 "statusBands": {
  "patch_methods_hit": {
   "min": 80,
   "max": 90
  },
  "patch_statement_coverage": {
   "min": 70,
   "max": 80
  },
  "statement_coverage": {
   "min": 60,
   "max": 75
  }
 },
 "folderBands": [
  {
   "path": "cmd",
   "bands": {
    "statement_coverage": {
     "min": 40,
     "max": 60
    }
   }
  },
  {
   "path": "internal/store",
   "bands": {
    "cyclomatic_complexity": {
     "min": 10,
     "max": 15
    },
    "statement_coverage": {
     "min": 80,
     "max": 90
    }
   }
  }
 ],
 "configs": [
  {
   "path": "",
   "source": "nanovision.yaml"
  },
  {
   "path": "cmd",
   "source": "nanovision.yaml (overrides)"
  },
  {
   "path": "demo_projects/cpp",
   "source": "demo_projects/cpp/nanovision.yaml"
  },
  {
   "path": "demo_projects/csharp",
   "source": "demo_projects/csharp/nanovision.yaml"
  },
  {
   "path": "demo_projects/go",
   "source": "demo_projects/go/nanovision.yaml"
  },
  {
   "path": "internal/store",
   "source": "internal/store/nanovision.yaml"
  }
 ]
}
