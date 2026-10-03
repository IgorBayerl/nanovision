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
 "generatedAt": "2026-10-03T15:00:30Z",
 "title": "nanovision Self-Coverage (dev fixture)",
 "totals": {
  "statement_coverage": {
   "covered": 5044,
   "uncovered": 1416,
   "coverable": 6460,
   "total": 6460,
   "percentage": 78.08
  },
  "methods_hit": {
   "covered": 643,
   "total": 727,
   "percentage": 88.44
  },
  "methods_fully_covered": {
   "covered": 375,
   "total": 727,
   "percentage": 51.58
  },
  "max_cyclomatic_complexity": {
   "value": 38
  },
  "patch_statement_coverage": {
   "covered": 215,
   "uncovered": 58,
   "coverable": 273,
   "total": 273,
   "percentage": 78.75
  },
  "patch_methods_hit": {
   "covered": 46,
   "total": 51,
   "percentage": 90.19
  },
  "files": 121,
  "folders": 56,
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
     "value": 38
    },
    "methods_fully_covered": {
     "covered": 7,
     "total": 28,
     "percentage": 25
    },
    "methods_hit": {
     "covered": 25,
     "total": 28,
     "percentage": 89.28
    },
    "patch_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 47,
     "uncovered": 10,
     "coverable": 57,
     "total": 57,
     "percentage": 82.45
    },
    "statement_coverage": {
     "covered": 286,
     "uncovered": 244,
     "coverable": 530,
     "total": 530,
     "percentage": 53.96
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "danger"
   }
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
    "statement_coverage": {
     "covered": 118,
     "uncovered": 38,
     "coverable": 156,
     "total": 156,
     "percentage": 75.64
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "cmd_history.go.html"
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
     "value": 38
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
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 47,
     "uncovered": 10,
     "coverable": 57,
     "total": 57,
     "percentage": 82.45
    },
    "statement_coverage": {
     "covered": 144,
     "uncovered": 80,
     "coverable": 224,
     "total": 224,
     "percentage": 64.28
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
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
    "statement_coverage": {
     "covered": 0,
     "uncovered": 90,
     "coverable": 90,
     "total": 90,
     "percentage": 0
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "cmd_serve.go.html"
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
    "statement_coverage": {
     "covered": 24,
     "uncovered": 36,
     "coverable": 60,
     "total": 60,
     "percentage": 40
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "cmd_storecmd.go.html"
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
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "demo_projects_cpp_project_src_calculator.cpp.html"
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
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "demo_projects_go_project_calculator_2_entities.go.html"
  },
  {
   "id": "internal",
   "name": "internal",
   "type": "folder",
   "path": "internal",
   "depth": 0,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 33
    },
    "methods_fully_covered": {
     "covered": 354,
     "total": 666,
     "percentage": 53.15
    },
    "methods_hit": {
     "covered": 594,
     "total": 666,
     "percentage": 89.18
    },
    "patch_methods_hit": {
     "covered": 40,
     "total": 45,
     "percentage": 88.88
    },
    "patch_statement_coverage": {
     "covered": 168,
     "uncovered": 48,
     "coverable": 216,
     "total": 216,
     "percentage": 77.77
    },
    "statement_coverage": {
     "covered": 4682,
     "uncovered": 1140,
     "coverable": 5822,
     "total": 5822,
     "percentage": 80.41
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "warning",
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
    "statement_coverage": {
     "covered": 328,
     "uncovered": 16,
     "coverable": 344,
     "total": 344,
     "percentage": 95.34
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
    "statement_coverage": {
     "covered": 165,
     "uncovered": 11,
     "coverable": 176,
     "total": 176,
     "percentage": 93.75
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
    "statement_coverage": {
     "covered": 212,
     "uncovered": 27,
     "coverable": 239,
     "total": 239,
     "percentage": 88.7
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
    "statement_coverage": {
     "covered": 81,
     "uncovered": 9,
     "coverable": 90,
     "total": 90,
     "percentage": 90
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
    "statement_coverage": {
     "covered": 81,
     "uncovered": 9,
     "coverable": 90,
     "total": 90,
     "percentage": 90
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
    "statement_coverage": {
     "covered": 60,
     "uncovered": 9,
     "coverable": 69,
     "total": 69,
     "percentage": 86.95
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
    "statement_coverage": {
     "covered": 60,
     "uncovered": 9,
     "coverable": 69,
     "total": 69,
     "percentage": 86.95
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
    "statement_coverage": {
     "covered": 71,
     "uncovered": 9,
     "coverable": 80,
     "total": 80,
     "percentage": 88.75
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
    "statement_coverage": {
     "covered": 71,
     "uncovered": 9,
     "coverable": 80,
     "total": 80,
     "percentage": 88.75
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
     "value": 21
    },
    "methods_fully_covered": {
     "covered": 28,
     "total": 49,
     "percentage": 57.14
    },
    "methods_hit": {
     "covered": 34,
     "total": 49,
     "percentage": 69.38
    },
    "statement_coverage": {
     "covered": 136,
     "uncovered": 22,
     "coverable": 158,
     "total": 158,
     "percentage": 86.07
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/calculator/calculators.go",
   "name": "calculators.go",
   "type": "file",
   "path": "internal/calculator/calculators.go",
   "parentId": "internal/calculator",
   "depth": 2,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 5
    },
    "methods_fully_covered": {
     "covered": 24,
     "total": 45,
     "percentage": 53.33
    },
    "methods_hit": {
     "covered": 30,
     "total": 45,
     "percentage": 66.66
    },
    "statement_coverage": {
     "covered": 73,
     "uncovered": 22,
     "coverable": 95,
     "total": 95,
     "percentage": 76.84
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_calculator_calculators.go.html",
   "diffStatus": "modified"
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
     "value": 21
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
    "statement_coverage": {
     "covered": 63,
     "uncovered": 0,
     "coverable": 63,
     "total": 63,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_calculator_engine.go.html"
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
    "statement_coverage": {
     "covered": 50,
     "uncovered": 16,
     "coverable": 66,
     "total": 66,
     "percentage": 75.75
    }
   },
   "statuses": {
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
    "statement_coverage": {
     "covered": 50,
     "uncovered": 16,
     "coverable": 66,
     "total": 66,
     "percentage": 75.75
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_client_client.go.html"
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
    "statement_coverage": {
     "covered": 137,
     "uncovered": 11,
     "coverable": 148,
     "total": 148,
     "percentage": 92.56
    }
   },
   "statuses": {
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
    "statement_coverage": {
     "covered": 137,
     "uncovered": 11,
     "coverable": 148,
     "total": 148,
     "percentage": 92.56
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_compare_compare.go.html"
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
     "value": 33
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 10,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 18,
     "uncovered": 13,
     "coverable": 31,
     "total": 31,
     "percentage": 58.06
    },
    "statement_coverage": {
     "covered": 103,
     "uncovered": 89,
     "coverable": 192,
     "total": 192,
     "percentage": 53.64
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
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
     "value": 33
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 10,
     "percentage": 40
    },
    "methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 6,
     "total": 6,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 18,
     "uncovered": 13,
     "coverable": 31,
     "total": 31,
     "percentage": 58.06
    },
    "statement_coverage": {
     "covered": 103,
     "uncovered": 89,
     "coverable": 192,
     "total": 192,
     "percentage": 53.64
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_config_config.go.html",
   "diffStatus": "modified"
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
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 2,
     "uncovered": 0,
     "coverable": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 55,
     "uncovered": 21,
     "coverable": 76,
     "total": 76,
     "percentage": 72.36
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
     "covered": 1,
     "total": 1,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 2,
     "uncovered": 0,
     "coverable": 2,
     "total": 2,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 55,
     "uncovered": 21,
     "coverable": 76,
     "total": 76,
     "percentage": 72.36
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
    "statement_coverage": {
     "covered": 132,
     "uncovered": 10,
     "coverable": 142,
     "total": 142,
     "percentage": 92.95
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
    "statement_coverage": {
     "covered": 121,
     "uncovered": 9,
     "coverable": 130,
     "total": 130,
     "percentage": 93.07
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
    "statement_coverage": {
     "covered": 11,
     "uncovered": 1,
     "coverable": 12,
     "total": 12,
     "percentage": 91.66
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
    "statement_coverage": {
     "covered": 194,
     "uncovered": 17,
     "coverable": 211,
     "total": 211,
     "percentage": 91.94
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
    "statement_coverage": {
     "covered": 48,
     "uncovered": 3,
     "coverable": 51,
     "total": 51,
     "percentage": 94.11
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
    "statement_coverage": {
     "covered": 146,
     "uncovered": 14,
     "coverable": 160,
     "total": 160,
     "percentage": 91.25
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
    "statement_coverage": {
     "covered": 97,
     "uncovered": 12,
     "coverable": 109,
     "total": 109,
     "percentage": 88.99
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
    "statement_coverage": {
     "covered": 97,
     "uncovered": 12,
     "coverable": 109,
     "total": 109,
     "percentage": 88.99
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
    "statement_coverage": {
     "covered": 71,
     "uncovered": 10,
     "coverable": 81,
     "total": 81,
     "percentage": 87.65
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
    "statement_coverage": {
     "covered": 67,
     "uncovered": 9,
     "coverable": 76,
     "total": 76,
     "percentage": 88.15
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
     "covered": 11,
     "total": 22,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 22,
     "total": 22,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 275,
     "uncovered": 36,
     "coverable": 311,
     "total": 311,
     "percentage": 88.42
    }
   },
   "statuses": {
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
    "statement_coverage": {
     "covered": 67,
     "uncovered": 2,
     "coverable": 69,
     "total": 69,
     "percentage": 97.1
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_history_capture.go.html"
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
     "covered": 39,
     "uncovered": 7,
     "coverable": 46,
     "total": 46,
     "percentage": 84.78
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_history_meta.go.html"
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
    "statement_coverage": {
     "covered": 72,
     "uncovered": 11,
     "coverable": 83,
     "total": 83,
     "percentage": 86.74
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_history_rebuild.go.html"
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
    "statement_coverage": {
     "covered": 97,
     "uncovered": 16,
     "coverable": 113,
     "total": 113,
     "percentage": 85.84
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_history_record.go.html"
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
    "statement_coverage": {
     "covered": 19,
     "uncovered": 2,
     "coverable": 21,
     "total": 21,
     "percentage": 90.47
    }
   },
   "statuses": {
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
    "statement_coverage": {
     "covered": 8,
     "uncovered": 2,
     "coverable": 10,
     "total": 10,
     "percentage": 80
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_model_comparison.go.html"
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
    "statement_coverage": {
     "covered": 236,
     "uncovered": 48,
     "coverable": 284,
     "total": 284,
     "percentage": 83.09
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
    "statement_coverage": {
     "covered": 63,
     "uncovered": 6,
     "coverable": 69,
     "total": 69,
     "percentage": 91.3
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
    "statement_coverage": {
     "covered": 29,
     "uncovered": 0,
     "coverable": 29,
     "total": 29,
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
    "statement_coverage": {
     "covered": 49,
     "uncovered": 19,
     "coverable": 68,
     "total": 68,
     "percentage": 72.05
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
    "statement_coverage": {
     "covered": 41,
     "uncovered": 4,
     "coverable": 45,
     "total": 45,
     "percentage": 91.11
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
    "statement_coverage": {
     "covered": 64,
     "uncovered": 22,
     "coverable": 86,
     "total": 86,
     "percentage": 74.41
    }
   },
   "statuses": {
    "statement_coverage": "warning"
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
    "statement_coverage": {
     "covered": 64,
     "uncovered": 22,
     "coverable": 86,
     "total": 86,
     "percentage": 74.41
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_pipeline_pipeline.go.html"
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
     "covered": 29,
     "total": 90,
     "percentage": 32.22
    },
    "methods_hit": {
     "covered": 63,
     "total": 90,
     "percentage": 70
    },
    "patch_methods_hit": {
     "covered": 18,
     "total": 20,
     "percentage": 90
    },
    "patch_statement_coverage": {
     "covered": 59,
     "uncovered": 16,
     "coverable": 75,
     "total": 75,
     "percentage": 78.66
    },
    "statement_coverage": {
     "covered": 632,
     "uncovered": 375,
     "coverable": 1007,
     "total": 1007,
     "percentage": 62.76
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "warning",
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
     "value": 22
    },
    "methods_fully_covered": {
     "covered": 21,
     "total": 54,
     "percentage": 38.88
    },
    "methods_hit": {
     "covered": 43,
     "total": 54,
     "percentage": 79.62
    },
    "patch_methods_hit": {
     "covered": 17,
     "total": 19,
     "percentage": 89.47
    },
    "patch_statement_coverage": {
     "covered": 56,
     "uncovered": 16,
     "coverable": 72,
     "total": 72,
     "percentage": 77.77
    },
    "statement_coverage": {
     "covered": 421,
     "uncovered": 237,
     "coverable": 658,
     "total": 658,
     "percentage": 63.98
    }
   },
   "statuses": {
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "warning",
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
     "value": 17
    },
    "methods_fully_covered": {
     "covered": 11,
     "total": 25,
     "percentage": 44
    },
    "methods_hit": {
     "covered": 22,
     "total": 25,
     "percentage": 88
    },
    "patch_methods_hit": {
     "covered": 13,
     "total": 14,
     "percentage": 92.85
    },
    "patch_statement_coverage": {
     "covered": 47,
     "uncovered": 8,
     "coverable": 55,
     "total": 55,
     "percentage": 85.45
    },
    "statement_coverage": {
     "covered": 247,
     "uncovered": 67,
     "coverable": 314,
     "total": 314,
     "percentage": 78.66
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
     "value": 22
    },
    "methods_fully_covered": {
     "covered": 9,
     "total": 18,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 17,
     "total": 18,
     "percentage": 94.44
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 9,
     "uncovered": 7,
     "coverable": 16,
     "total": 16,
     "percentage": 56.25
    },
    "statement_coverage": {
     "covered": 125,
     "uncovered": 66,
     "coverable": 191,
     "total": 191,
     "percentage": 65.44
    }
   },
   "statuses": {
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
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
    "statement_coverage": {
     "covered": 1,
     "uncovered": 1,
     "coverable": 2,
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
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_htmlreact_generator_single.go.html"
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
    "statement_coverage": {
     "covered": 0,
     "uncovered": 15,
     "coverable": 15,
     "total": 15,
     "percentage": 0
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_htmlreact_views.go.html"
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
    "statement_coverage": {
     "covered": 151,
     "uncovered": 58,
     "coverable": 209,
     "total": 209,
     "percentage": 72.24
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
    "statement_coverage": {
     "covered": 65,
     "uncovered": 10,
     "coverable": 75,
     "total": 75,
     "percentage": 86.66
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_textsummary_comparison.go.html"
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
    "statement_coverage": {
     "covered": 29,
     "uncovered": 35,
     "coverable": 64,
     "total": 64,
     "percentage": 45.31
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
    "statement_coverage": {
     "covered": 57,
     "uncovered": 13,
     "coverable": 70,
     "total": 70,
     "percentage": 81.42
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_textsummary_terminal.go.html"
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
    "statement_coverage": {
     "covered": 315,
     "uncovered": 125,
     "coverable": 440,
     "total": 440,
     "percentage": 71.59
    }
   },
   "statuses": {
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
    "statement_coverage": {
     "covered": 129,
     "uncovered": 43,
     "coverable": 172,
     "total": 172,
     "percentage": 75
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_server_api.go.html"
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
    "statement_coverage": {
     "covered": 38,
     "uncovered": 9,
     "coverable": 47,
     "total": 47,
     "percentage": 80.85
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_server_cache.go.html"
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
    "statement_coverage": {
     "covered": 0,
     "uncovered": 24,
     "coverable": 24,
     "total": 24,
     "percentage": 0
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_server_maintenance.go.html"
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
    "statement_coverage": {
     "covered": 55,
     "uncovered": 7,
     "coverable": 62,
     "total": 62,
     "percentage": 88.7
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_server_server.go.html"
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
    "statement_coverage": {
     "covered": 24,
     "uncovered": 17,
     "coverable": 41,
     "total": 41,
     "percentage": 58.53
    }
   },
   "statuses": {
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_server_ui.go.html"
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
    "statement_coverage": {
     "covered": 34,
     "uncovered": 14,
     "coverable": 48,
     "total": 48,
     "percentage": 70.83
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_server_upload.go.html"
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
    "statement_coverage": {
     "covered": 35,
     "uncovered": 11,
     "coverable": 46,
     "total": 46,
     "percentage": 76.08
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_server_wire.go.html"
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
     "covered": 102,
     "total": 118,
     "percentage": 86.44
    },
    "methods_hit": {
     "covered": 111,
     "total": 118,
     "percentage": 94.06
    },
    "statement_coverage": {
     "covered": 229,
     "uncovered": 20,
     "coverable": 249,
     "total": 249,
     "percentage": 91.96
    }
   },
   "statuses": {
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
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 94,
     "total": 108,
     "percentage": 87.03
    },
    "methods_hit": {
     "covered": 101,
     "total": 108,
     "percentage": 93.51
    },
    "statement_coverage": {
     "covered": 163,
     "uncovered": 17,
     "coverable": 180,
     "total": 180,
     "percentage": 90.55
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   }
  },
  {
   "id": "internal/status/evaluators/complexity.go",
   "name": "complexity.go",
   "type": "file",
   "path": "internal/status/evaluators/complexity.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 22,
     "total": 24,
     "percentage": 91.66
    },
    "methods_hit": {
     "covered": 24,
     "total": 24,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 38,
     "uncovered": 2,
     "coverable": 40,
     "total": 40,
     "percentage": 95
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_complexity.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "internal/status/evaluators/line_coverage.go",
   "name": "line_coverage.go",
   "type": "file",
   "path": "internal/status/evaluators/line_coverage.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 9,
     "total": 12,
     "percentage": 75
    },
    "methods_hit": {
     "covered": 10,
     "total": 12,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 17,
     "uncovered": 3,
     "coverable": 20,
     "total": 20,
     "percentage": 85
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_line_coverage.go.html"
  },
  {
   "id": "internal/status/evaluators/methods_fully_covered.go",
   "name": "methods_fully_covered.go",
   "type": "file",
   "path": "internal/status/evaluators/methods_fully_covered.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
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
     "covered": 10,
     "uncovered": 0,
     "coverable": 10,
     "total": 10,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_methods_fully_covered.go.html"
  },
  {
   "id": "internal/status/evaluators/methods_hit.go",
   "name": "methods_hit.go",
   "type": "file",
   "path": "internal/status/evaluators/methods_hit.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
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
     "covered": 10,
     "uncovered": 0,
     "coverable": 10,
     "total": 10,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_methods_hit.go.html"
  },
  {
   "id": "internal/status/evaluators/patch_line_coverage.go",
   "name": "patch_line_coverage.go",
   "type": "file",
   "path": "internal/status/evaluators/patch_line_coverage.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 10,
     "total": 12,
     "percentage": 83.33
    },
    "methods_hit": {
     "covered": 10,
     "total": 12,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 18,
     "uncovered": 2,
     "coverable": 20,
     "total": 20,
     "percentage": 90
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_patch_line_coverage.go.html"
  },
  {
   "id": "internal/status/evaluators/patch_methods_hit.go",
   "name": "patch_methods_hit.go",
   "type": "file",
   "path": "internal/status/evaluators/patch_methods_hit.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
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
     "covered": 10,
     "uncovered": 0,
     "coverable": 10,
     "total": 10,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_patch_methods_hit.go.html"
  },
  {
   "id": "internal/status/evaluators/patch_statement_coverage.go",
   "name": "patch_statement_coverage.go",
   "type": "file",
   "path": "internal/status/evaluators/patch_statement_coverage.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 20,
     "uncovered": 0,
     "coverable": 20,
     "total": 20,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_patch_statement_coverage.go.html"
  },
  {
   "id": "internal/status/evaluators/patch_statement_methods_hit.go",
   "name": "patch_statement_methods_hit.go",
   "type": "file",
   "path": "internal/status/evaluators/patch_statement_methods_hit.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 6,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 7,
     "uncovered": 3,
     "coverable": 10,
     "total": 10,
     "percentage": 70
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_patch_statement_methods_hit.go.html"
  },
  {
   "id": "internal/status/evaluators/statement_coverage.go",
   "name": "statement_coverage.go",
   "type": "file",
   "path": "internal/status/evaluators/statement_coverage.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 11,
     "total": 12,
     "percentage": 91.66
    },
    "methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 19,
     "uncovered": 1,
     "coverable": 20,
     "total": 20,
     "percentage": 95
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_statement_coverage.go.html"
  },
  {
   "id": "internal/status/evaluators/statement_methods_fully_covered.go",
   "name": "statement_methods_fully_covered.go",
   "type": "file",
   "path": "internal/status/evaluators/statement_methods_fully_covered.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 6,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 7,
     "uncovered": 3,
     "coverable": 10,
     "total": 10,
     "percentage": 70
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_statement_methods_fully_covered.go.html"
  },
  {
   "id": "internal/status/evaluators/statement_methods_hit.go",
   "name": "statement_methods_hit.go",
   "type": "file",
   "path": "internal/status/evaluators/statement_methods_hit.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
   "metrics": {
    "max_cyclomatic_complexity": {
     "value": 2
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 6,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "statement_coverage": {
     "covered": 7,
     "uncovered": 3,
     "coverable": 10,
     "total": 10,
     "percentage": 70
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_statement_methods_hit.go.html"
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
    "statement_coverage": {
     "covered": 38,
     "uncovered": 3,
     "coverable": 41,
     "total": 41,
     "percentage": 92.68
    }
   },
   "statuses": {
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
    "statement_coverage": {
     "covered": 13,
     "uncovered": 0,
     "coverable": 13,
     "total": 13,
     "percentage": 100
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_capabilities.go.html"
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
    "statement_coverage": {
     "covered": 763,
     "uncovered": 169,
     "coverable": 932,
     "total": 932,
     "percentage": 81.86
    }
   },
   "statuses": {
    "statement_coverage": "safe"
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
    "statement_coverage": {
     "covered": 478,
     "uncovered": 88,
     "coverable": 566,
     "total": 566,
     "percentage": 84.45
    }
   },
   "statuses": {
    "statement_coverage": "safe"
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
    "statement_coverage": {
     "covered": 94,
     "uncovered": 23,
     "coverable": 117,
     "total": 117,
     "percentage": 80.34
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_store_blob_analysis.go.html"
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
    "statement_coverage": {
     "covered": 103,
     "uncovered": 18,
     "coverable": 121,
     "total": 121,
     "percentage": 85.12
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_store_blob_blob.go.html"
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
    "statement_coverage": {
     "covered": 104,
     "uncovered": 17,
     "coverable": 121,
     "total": 121,
     "percentage": 85.95
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_store_blob_coverage.go.html"
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
    "statement_coverage": {
     "covered": 77,
     "uncovered": 14,
     "coverable": 91,
     "total": 91,
     "percentage": 84.61
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_store_blob_diff.go.html"
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
    "statement_coverage": {
     "covered": 100,
     "uncovered": 16,
     "coverable": 116,
     "total": 116,
     "percentage": 86.2
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_store_blob_manifest.go.html"
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
    "statement_coverage": {
     "covered": 68,
     "uncovered": 34,
     "coverable": 102,
     "total": 102,
     "percentage": 66.66
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_blobs.go.html"
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
    "statement_coverage": {
     "covered": 79,
     "uncovered": 19,
     "coverable": 98,
     "total": 98,
     "percentage": 80.61
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_store_maintenance.go.html"
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
    "statement_coverage": {
     "covered": 101,
     "uncovered": 12,
     "coverable": 113,
     "total": 113,
     "percentage": 89.38
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_store_runs.go.html"
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
    "statement_coverage": {
     "covered": 37,
     "uncovered": 16,
     "coverable": 53,
     "total": 53,
     "percentage": 69.81
    }
   },
   "statuses": {
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_store_store.go.html"
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
    "statement_coverage": {
     "covered": 79,
     "uncovered": 5,
     "coverable": 84,
     "total": 84,
     "percentage": 94.04
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
    "statement_coverage": {
     "covered": 79,
     "uncovered": 5,
     "coverable": 84,
     "total": 84,
     "percentage": 94.04
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
    "statement_coverage": {
     "covered": 201,
     "uncovered": 39,
     "coverable": 240,
     "total": 240,
     "percentage": 83.75
    }
   },
   "statuses": {
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
    "statement_coverage": {
     "covered": 65,
     "uncovered": 14,
     "coverable": 79,
     "total": 79,
     "percentage": 82.27
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_vcs_git.go.html"
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
    "statement_coverage": {
     "covered": 106,
     "uncovered": 19,
     "coverable": 125,
     "total": 125,
     "percentage": 84.8
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_vcs_perforce.go.html"
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
     "value": 8
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
    "statement_coverage": {
     "covered": 30,
     "uncovered": 6,
     "coverable": 36,
     "total": 36,
     "percentage": 83.33
    }
   },
   "statuses": {
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_vcs_vcs.go.html"
  }
 ],
 "metricDefinitions": {
  "a_statement_coverage": {
   "label": "Statements",
   "shortLabel": "Statements",
   "description": "Percentage of executed statements.",
   "subMetrics": [
    {
     "id": "total",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "c_patch_statement_coverage": {
   "label": "Patch Statements",
   "shortLabel": "Patch Stmts",
   "description": "Statement coverage of changed (patched) code only.",
   "subMetrics": [
    {
     "id": "total",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "f_cyclomatic_complexity": {
   "label": "Cyclomatic Complexity",
   "shortLabel": "Complexity",
   "description": "Cyclomatic complexity of a function (lower is better).",
   "kind": "value",
   "subMetrics": [
    {
     "id": "value",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "g_crap_score": {
   "label": "CRAP Score",
   "shortLabel": "CRAP",
   "description": "Change Risk Anti-Pattern (CRAP) score combining complexity and coverage (lower is better).",
   "subMetrics": [
    {
     "id": "total",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "i_exposed_risk": {
   "label": "Exposed Risk",
   "shortLabel": "Risk",
   "description": "Absolute volume of complexity that is unprotected by tests (lower is better).",
   "subMetrics": [
    {
     "id": "total",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "max_cyclomatic_complexity": {
   "label": "Max Cyclomatic Complexity",
   "shortLabel": "Max Complexity",
   "description": "Maximum cyclomatic complexity of a function (lower is better).",
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
   "description": "Percentage of methods with 100% line coverage.",
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
   "description": "Percentage of methods with at least one hit.",
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
   "description": "Percentage of patched methods with at least one hit.",
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
   "description": "Statement coverage of changed (patched) code only.",
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
  "statement_coverage": {
   "label": "Statements",
   "shortLabel": "Statements",
   "description": "Percentage of executed statements.",
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
   "value": "2026-10-03 15:00:30"
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
    "value": 78.75457875457876,
    "threshold": 80,
    "passed": false
   },
   {
    "key": "max_changed_method_complexity",
    "label": "Max changed-method complexity",
    "value": 38,
    "threshold": 15,
    "passed": false
   }
  ],
  "stats": {
   "changedFiles": 28,
   "methodsAdded": 11,
   "methodsModified": 45,
   "untestedChangedMethods": 8,
   "patchStatementsValid": 273,
   "patchStatementsCovered": 215,
   "maxChangedComplexity": 38
  },
  "hotspots": [
   {
    "file": "cmd/main.go",
    "method": "main",
    "startLine": 220,
    "diffStatus": "modified",
    "complexity": 38,
    "patchCoverage": 74.19354838709677,
    "risk": 17.08527131782946
   },
   {
    "file": "internal/config/config.go",
    "method": "(*AppConfig).mergeCliOverrides",
    "startLine": 347,
    "diffStatus": "modified",
    "complexity": 33,
    "patchCoverage": 63.63636363636363,
    "risk": 15.04411764705882
   },
   {
    "file": "internal/reporter/htmlreact/details_generator.go",
    "method": "(*HtmlReactReportBuilder).buildMethodDetails",
    "startLine": 229,
    "diffStatus": "modified",
    "complexity": 22,
    "patchCoverage": 0,
    "risk": 11.314285714285713
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
    "file": "internal/config/config.go",
    "method": "(*AppConfig).validate",
    "startLine": 460,
    "diffStatus": "modified",
    "complexity": 15,
    "patchCoverage": 75,
    "risk": 6.923076923076923
   },
   {
    "file": "internal/config/config.go",
    "method": "Load",
    "startLine": 291,
    "diffStatus": "modified",
    "complexity": 11,
    "patchCoverage": 50,
    "risk": 4.7142857142857135
   },
   {
    "file": "internal/reporter/htmlreact/builder.go",
    "method": "addMeta",
    "startLine": 298,
    "diffStatus": "modified",
    "complexity": 8,
    "patchCoverage": 100,
    "risk": 3.2
   },
   {
    "file": "internal/enricher/enricher.go",
    "method": "(*Enricher).enrichFileNode",
    "startLine": 135,
    "diffStatus": "modified",
    "complexity": 11,
    "patchCoverage": 66.66666666666667,
    "risk": 3.1842105263157894
   },
   {
    "file": "internal/reporter/htmlreact/builder.go",
    "method": "describeMetric",
    "startLine": 742,
    "diffStatus": "modified",
    "complexity": 3,
    "patchCoverage": 0,
    "risk": 3
   },
   {
    "file": "internal/config/config.go",
    "method": "(*AppConfig).computeDerivedFields",
    "startLine": 514,
    "diffStatus": "modified",
    "complexity": 11,
    "patchCoverage": 40,
    "risk": 2.933333333333334
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
   "value": "188 from fixture.diff"
  }
 ],
 "reports": [
  {
   "name": "coverage-unit.out",
   "path": "reports/nanovision_self_coverage/coverage-unit.out"
  },
  {
   "name": "coverage-integration.out",
   "path": "reports/nanovision_self_coverage/coverage-integration.out"
  },
  {
   "name": "*.gcov",
   "path": "demo_projects/cpp/report/gcov/branch-probabilities/*.gcov"
  },
  {
   "name": "cobertura.xml",
   "path": "demo_projects/csharp/report/cobertura/cobertura.xml"
  },
  {
   "name": "coverage.out",
   "path": "demo_projects/go/report/gocover/coverage.out"
  }
 ],
 "reportIndexes": {
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
     "n": 6
    },
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
      2
     ],
     "n": 11
    }
   ],
   "patch_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 6
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 47
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
     "n": 6
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 144
    },
    {
     "m": [
      0
     ],
     "n": 80
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 6
    },
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
      2
     ],
     "n": 11
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
     "n": 3
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
      3
     ],
     "n": 16
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
     "n": 118
    },
    {
     "m": [
      1
     ],
     "n": 47
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
     "n": 3
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
      3
     ],
     "n": 16
    },
    {
     "m": [
      1
     ],
     "n": 2
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
  "internal/calculator/calculators.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 24
    },
    {
     "m": [
      0
     ],
     "n": 21
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 30
    },
    {
     "m": [
      0
     ],
     "n": 15
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 73
    },
    {
     "m": [
      0
     ],
     "n": 22
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 24
    },
    {
     "m": [
      0
     ],
     "n": 21
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 30
    },
    {
     "m": [
      0
     ],
     "n": 15
    }
   ]
  },
  "internal/calculator/engine.go": {
   "methods_fully_covered": [
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
     "n": 4
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 63
    }
   ],
   "statement_methods_fully_covered": [
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
     "n": 4
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
     "n": 11
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
     "n": 6
    },
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
     "n": 4
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
      2
     ],
     "n": 1
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      0
     ],
     "n": 13
    },
    {
     "m": [
      3
     ],
     "n": 8
    },
    {
     "m": [
      2
     ],
     "n": 6
    },
    {
     "m": [
      1
     ],
     "n": 4
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
     "n": 89
    },
    {
     "m": [
      2
     ],
     "n": 45
    },
    {
     "m": [
      3
     ],
     "n": 43
    },
    {
     "m": [
      1
     ],
     "n": 15
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
     "n": 1
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
     "n": 1
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
      1
     ],
     "n": 26
    },
    {
     "m": [
      3
     ],
     "n": 12
    },
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
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 64
    },
    {
     "m": [
      0
     ],
     "n": 22
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
     "n": 14
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
     "n": 2
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
      2
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
     "n": 2
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
     "n": 5
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
      1
     ],
     "n": 18
    },
    {
     "m": [
      3
     ],
     "n": 16
    },
    {
     "m": [
      2
     ],
     "n": 13
    },
    {
     "m": [
      0
     ],
     "n": 8
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
     "n": 5
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
     "n": 122
    },
    {
     "m": [
      2
     ],
     "n": 93
    },
    {
     "m": [
      0
     ],
     "n": 67
    },
    {
     "m": [
      1
     ],
     "n": 32
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 17
    },
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
      2
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
     "n": 4
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
      3
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
      2
     ],
     "n": 2
    },
    {
     "m": [
      2,
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
     "n": 8
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
    }
   ],
   "patch_statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 9
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
      2
     ],
     "n": 4
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
     "n": 66
    },
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
     "n": 6
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      0
     ],
     "n": 13
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
      2
     ],
     "n": 9
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
      2,
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
     "n": 33
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
     "n": 2
    },
    {
     "m": [
      2,
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
  "internal/status/evaluators/complexity.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 19
    },
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
    }
   ],
   "methods_hit": [
    {
     "m": [
      2
     ],
     "n": 21
    },
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
      2
     ],
     "n": 31
    },
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
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 19
    },
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
    }
   ],
   "statement_methods_hit": [
    {
     "m": [
      2
     ],
     "n": 21
    },
    {
     "m": [
      3
     ],
     "n": 3
    }
   ]
  },
  "internal/status/evaluators/line_coverage.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 6
    },
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
      2
     ],
     "n": 7
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
     "n": 6
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
      2
     ],
     "n": 6
    },
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
      2
     ],
     "n": 7
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
   ]
  },
  "internal/status/evaluators/methods_fully_covered.go": {
   "methods_fully_covered": [
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
     "n": 1
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
     "n": 1
    }
   ],
   "statement_coverage": [
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
     "n": 1
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
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ]
  },
  "internal/status/evaluators/methods_hit.go": {
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
     "n": 3
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
      2
     ],
     "n": 3
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
     "n": 3
    }
   ]
  },
  "internal/status/evaluators/patch_line_coverage.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 7
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
     "n": 7
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
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 11
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
     "n": 2
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 7
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
     "n": 7
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
   ]
  },
  "internal/status/evaluators/patch_methods_hit.go": {
   "methods_fully_covered": [
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
     "n": 1
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
     "n": 1
    }
   ],
   "statement_coverage": [
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
     "n": 1
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
     "n": 5
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ]
  },
  "internal/status/evaluators/patch_statement_coverage.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 11
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
     "n": 11
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
     "n": 19
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
      2
     ],
     "n": 11
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
     "n": 11
    },
    {
     "m": [
      3
     ],
     "n": 1
    }
   ]
  },
  "internal/status/evaluators/patch_statement_methods_hit.go": {
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
     "n": 4
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
      2
     ],
     "n": 6
    },
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
     "n": 4
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
  "internal/status/evaluators/statement_coverage.go": {
   "methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 8
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
     "n": 1
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
     "n": 3
    }
   ],
   "statement_coverage": [
    {
     "m": [
      2
     ],
     "n": 12
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
     "n": 1
    }
   ],
   "statement_methods_fully_covered": [
    {
     "m": [
      2
     ],
     "n": 8
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
     "n": 1
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
     "n": 3
    }
   ]
  },
  "internal/status/evaluators/statement_methods_fully_covered.go": {
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
     "n": 4
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
      2
     ],
     "n": 6
    },
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
     "n": 4
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
  "internal/status/evaluators/statement_methods_hit.go": {
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
     "n": 4
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
      2
     ],
     "n": 6
    },
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
     "n": 4
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
     "n": 12
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
 }
}
