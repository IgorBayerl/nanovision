/**
 * NOTE: This is a developer data fixture for testing purposes, loaded only by
 * the Vite dev server (pnpm dev). In production, the Go reporter writes a
 * real data.js next to index.html.
 *
 * Generated from the nanovision self-coverage reports (full merged) with a
 * git diff applied, plus the HtmlReview `review` block (gate verdict,
 * changelist stats, risk hotspots) so the review header renders in dev.
 *
 * Regenerate with: python scripts/gen_ui_fixture.py
 * Schema version: 1 (flat node list)
 */
window.__NANOVISION_SUMMARY__ = {
 "schemaVersion": 1,
 "generatedAt": "2026-07-06T22:11:05Z",
 "title": "nanovision Self-Coverage (dev fixture)",
 "totals": {
  "statement_coverage": {
   "covered": 2776,
   "uncovered": 843,
   "coverable": 3619,
   "total": 3619,
   "percentage": 76.7
  },
  "branch_coverage": {
   "covered": 39,
   "total": 58,
   "percentage": 67.24
  },
  "methods_hit": {
   "covered": 380,
   "total": 477,
   "percentage": 79.66
  },
  "methods_fully_covered": {
   "covered": 244,
   "total": 477,
   "percentage": 51.15
  },
  "max_cyclomatic_complexity": {
   "value": 31
  },
  "patch_statement_coverage": {
   "covered": 640,
   "uncovered": 325,
   "coverable": 965,
   "total": 965,
   "percentage": 66.32
  },
  "patch_methods_hit": {
   "covered": 187,
   "total": 260,
   "percentage": 71.92
  },
  "files": 87,
  "folders": 48,
  "statuses": {
   "branch_coverage": "warning",
   "methods_fully_covered": "warning",
   "patch_methods_hit": "danger",
   "patch_statement_coverage": "danger",
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
     "value": 28
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 13,
     "percentage": 46.15
    },
    "methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 6,
     "percentage": 66.66
    },
    "patch_statement_coverage": {
     "covered": 37,
     "uncovered": 45,
     "coverable": 82,
     "total": 82,
     "percentage": 45.12
    },
    "statement_coverage": {
     "covered": 179,
     "uncovered": 112,
     "coverable": 291,
     "total": 291,
     "percentage": 61.51
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   }
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
     "value": 28
    },
    "methods_fully_covered": {
     "covered": 6,
     "total": 13,
     "percentage": 46.15
    },
    "methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 6,
     "percentage": 66.66
    },
    "patch_statement_coverage": {
     "covered": 37,
     "uncovered": 45,
     "coverable": 82,
     "total": 82,
     "percentage": 45.12
    },
    "statement_coverage": {
     "covered": 179,
     "uncovered": 112,
     "coverable": 291,
     "total": 291,
     "percentage": 61.51
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "cmd_main.go.html",
   "diffStatus": "modified"
  },
  {
   "id": "demo_projects",
   "name": "demo_projects",
   "type": "folder",
   "path": "demo_projects",
   "depth": 0,
   "metrics": {
    "branch_coverage": {
     "covered": 39,
     "total": 58,
     "percentage": 67.24
    },
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
    "patch_methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    },
    "patch_statement_coverage": {
     "covered": 20,
     "uncovered": 6,
     "coverable": 26,
     "total": 26,
     "percentage": 76.92
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
    "branch_coverage": "warning",
    "methods_fully_covered": "warning",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
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
    "branch_coverage": {
     "covered": 21,
     "total": 34,
     "percentage": 61.76
    },
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
    "branch_coverage": "warning",
    "methods_fully_covered": "danger",
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
    "branch_coverage": {
     "covered": 21,
     "total": 34,
     "percentage": 61.76
    },
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
    "branch_coverage": "warning",
    "methods_fully_covered": "danger",
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
    "branch_coverage": {
     "covered": 21,
     "total": 34,
     "percentage": 61.76
    },
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
    "branch_coverage": "warning",
    "methods_fully_covered": "danger",
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
    "branch_coverage": {
     "covered": 9,
     "total": 14,
     "percentage": 64.28
    },
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
    "branch_coverage": "warning",
    "methods_fully_covered": "danger",
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
    "branch_coverage": {
     "covered": 9,
     "total": 14,
     "percentage": 64.28
    },
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
    "branch_coverage": "warning",
    "methods_fully_covered": "danger",
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
    "branch_coverage": {
     "covered": 6,
     "total": 12,
     "percentage": 50
    },
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
    "branch_coverage": "warning",
    "methods_fully_covered": "danger",
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
    "branch_coverage": {
     "covered": 6,
     "total": 8,
     "percentage": 75
    },
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
    "branch_coverage": "safe",
    "methods_fully_covered": "warning",
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
    "branch_coverage": {
     "covered": 18,
     "total": 24,
     "percentage": 75
    },
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "statuses": {
    "branch_coverage": "safe"
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
    "branch_coverage": {
     "covered": 18,
     "total": 24,
     "percentage": 75
    },
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "statuses": {
    "branch_coverage": "safe"
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
    "branch_coverage": {
     "covered": 18,
     "total": 24,
     "percentage": 75
    },
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "statuses": {
    "branch_coverage": "safe"
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
    "branch_coverage": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "statuses": {
    "branch_coverage": "safe"
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
    "branch_coverage": {
     "covered": 2,
     "total": 4,
     "percentage": 50
    },
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "statuses": {
    "branch_coverage": "warning"
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
    "branch_coverage": {
     "covered": 4,
     "total": 8,
     "percentage": 50
    },
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "statuses": {
    "branch_coverage": "warning"
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
    "branch_coverage": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "max_cyclomatic_complexity": {
     "value": 0
    }
   },
   "statuses": {
    "branch_coverage": "safe"
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
    "patch_methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    },
    "patch_statement_coverage": {
     "covered": 20,
     "uncovered": 6,
     "coverable": 26,
     "total": 26,
     "percentage": 76.92
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
    "methods_fully_covered": "warning",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
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
    "patch_methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    },
    "patch_statement_coverage": {
     "covered": 20,
     "uncovered": 6,
     "coverable": 26,
     "total": 26,
     "percentage": 76.92
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
    "methods_fully_covered": "warning",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
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
    "patch_methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    },
    "patch_statement_coverage": {
     "covered": 20,
     "uncovered": 6,
     "coverable": 26,
     "total": 26,
     "percentage": 76.92
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
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
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
    "patch_methods_hit": {
     "covered": 6,
     "total": 7,
     "percentage": 85.71
    },
    "patch_statement_coverage": {
     "covered": 20,
     "uncovered": 6,
     "coverable": 26,
     "total": 26,
     "percentage": 76.92
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
    "methods_fully_covered": "warning",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "demo_projects_go_project_calculator_calculator.go.html",
   "diffStatus": "added"
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
    "methods_fully_covered": "safe",
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
    "methods_fully_covered": "danger",
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
    "methods_fully_covered": "warning",
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
    "methods_fully_covered": "danger",
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
     "value": 31
    },
    "methods_fully_covered": {
     "covered": 224,
     "total": 431,
     "percentage": 51.97
    },
    "methods_hit": {
     "covered": 343,
     "total": 431,
     "percentage": 79.58
    },
    "patch_methods_hit": {
     "covered": 177,
     "total": 247,
     "percentage": 71.65
    },
    "patch_statement_coverage": {
     "covered": 583,
     "uncovered": 274,
     "coverable": 857,
     "total": 857,
     "percentage": 68.02
    },
    "statement_coverage": {
     "covered": 2521,
     "uncovered": 699,
     "coverable": 3220,
     "total": 3220,
     "percentage": 78.29
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
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
     "covered": 9,
     "total": 13,
     "percentage": 69.23
    },
    "methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 5,
     "uncovered": 0,
     "coverable": 5,
     "total": 5,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 160,
     "uncovered": 10,
     "coverable": 170,
     "total": 170,
     "percentage": 94.11
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
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
     "covered": 9,
     "total": 13,
     "percentage": 69.23
    },
    "methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 5,
     "uncovered": 0,
     "coverable": 5,
     "total": 5,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 160,
     "uncovered": 10,
     "coverable": 170,
     "total": 170,
     "percentage": 94.11
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_aggregator_aggrgator.go.html",
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
    "statement_coverage": {
     "covered": 212,
     "uncovered": 27,
     "coverable": 239,
     "total": 239,
     "percentage": 88.7
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
    "statement_coverage": {
     "covered": 81,
     "uncovered": 9,
     "coverable": 90,
     "total": 90,
     "percentage": 90
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
    "statement_coverage": {
     "covered": 81,
     "uncovered": 9,
     "coverable": 90,
     "total": 90,
     "percentage": 90
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_analyzer_cpp_analyzer.go.html"
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
    "statement_coverage": {
     "covered": 60,
     "uncovered": 9,
     "coverable": 69,
     "total": 69,
     "percentage": 86.95
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
    "statement_coverage": {
     "covered": 60,
     "uncovered": 9,
     "coverable": 69,
     "total": 69,
     "percentage": 86.95
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_analyzer_gdscript_analyzer.go.html"
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
    "statement_coverage": {
     "covered": 71,
     "uncovered": 9,
     "coverable": 80,
     "total": 80,
     "percentage": 88.75
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
    "statement_coverage": {
     "covered": 71,
     "uncovered": 9,
     "coverable": 80,
     "total": 80,
     "percentage": 88.75
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_analyzer_go_analyzer.go.html"
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
     "covered": 28,
     "uncovered": 1,
     "coverable": 29,
     "total": 29,
     "percentage": 96.55
    },
    "statement_coverage": {
     "covered": 28,
     "uncovered": 1,
     "coverable": 29,
     "total": 29,
     "percentage": 96.55
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
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
     "covered": 28,
     "uncovered": 1,
     "coverable": 29,
     "total": 29,
     "percentage": 96.55
    },
    "statement_coverage": {
     "covered": 28,
     "uncovered": 1,
     "coverable": 29,
     "total": 29,
     "percentage": 96.55
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_bootlog_bootlog.go.html",
   "diffStatus": "added"
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
     "covered": 4,
     "total": 11,
     "percentage": 36.36
    },
    "methods_hit": {
     "covered": 8,
     "total": 11,
     "percentage": 72.72
    },
    "statement_coverage": {
     "covered": 80,
     "uncovered": 14,
     "coverable": 94,
     "total": 94,
     "percentage": 85.1
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
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
     "covered": 4,
     "total": 8,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 8,
     "total": 8,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 80,
     "uncovered": 11,
     "coverable": 91,
     "total": 91,
     "percentage": 87.91
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_cache_cache.go.html",
   "diffStatus": "modified"
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
     "uncovered": 3,
     "coverable": 3,
     "total": 3,
     "percentage": 0
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_cache_cache_validator.go.html",
   "diffStatus": "modified"
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
     "covered": 26,
     "total": 61,
     "percentage": 42.62
    },
    "methods_hit": {
     "covered": 36,
     "total": 61,
     "percentage": 59.01
    },
    "patch_methods_hit": {
     "covered": 36,
     "total": 61,
     "percentage": 59.01
    },
    "patch_statement_coverage": {
     "covered": 140,
     "uncovered": 65,
     "coverable": 205,
     "total": 205,
     "percentage": 68.29
    },
    "statement_coverage": {
     "covered": 140,
     "uncovered": 65,
     "coverable": 205,
     "total": 205,
     "percentage": 68.29
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
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
     "value": 10
    },
    "methods_fully_covered": {
     "covered": 22,
     "total": 57,
     "percentage": 38.59
    },
    "methods_hit": {
     "covered": 32,
     "total": 57,
     "percentage": 56.14
    },
    "patch_methods_hit": {
     "covered": 32,
     "total": 57,
     "percentage": 56.14
    },
    "patch_statement_coverage": {
     "covered": 77,
     "uncovered": 65,
     "coverable": 142,
     "total": 142,
     "percentage": 54.22
    },
    "statement_coverage": {
     "covered": 77,
     "uncovered": 65,
     "coverable": 142,
     "total": 142,
     "percentage": 54.22
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_calculator_calculators.go.html",
   "diffStatus": "added"
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
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 63,
     "uncovered": 0,
     "coverable": 63,
     "total": 63,
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
    "methods_fully_covered": "safe",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_calculator_engine.go.html",
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
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 32,
     "uncovered": 3,
     "coverable": 35,
     "total": 35,
     "percentage": 91.42
    },
    "statement_coverage": {
     "covered": 124,
     "uncovered": 42,
     "coverable": 166,
     "total": 166,
     "percentage": 74.69
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
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
     "covered": 3,
     "total": 9,
     "percentage": 33.33
    },
    "methods_hit": {
     "covered": 9,
     "total": 9,
     "percentage": 100
    },
    "patch_methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 32,
     "uncovered": 3,
     "coverable": 35,
     "total": 35,
     "percentage": 91.42
    },
    "statement_coverage": {
     "covered": 124,
     "uncovered": 42,
     "coverable": 166,
     "total": 166,
     "percentage": 74.69
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "warning"
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
     "covered": 3,
     "total": 11,
     "percentage": 27.27
    },
    "methods_hit": {
     "covered": 9,
     "total": 11,
     "percentage": 81.81
    },
    "patch_methods_hit": {
     "covered": 9,
     "total": 11,
     "percentage": 81.81
    },
    "patch_statement_coverage": {
     "covered": 55,
     "uncovered": 25,
     "coverable": 80,
     "total": 80,
     "percentage": 68.75
    },
    "statement_coverage": {
     "covered": 55,
     "uncovered": 25,
     "coverable": 80,
     "total": 80,
     "percentage": 68.75
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
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
     "covered": 3,
     "total": 11,
     "percentage": 27.27
    },
    "methods_hit": {
     "covered": 9,
     "total": 11,
     "percentage": 81.81
    },
    "patch_methods_hit": {
     "covered": 9,
     "total": 11,
     "percentage": 81.81
    },
    "patch_statement_coverage": {
     "covered": 55,
     "uncovered": 25,
     "coverable": 80,
     "total": 80,
     "percentage": 68.75
    },
    "statement_coverage": {
     "covered": 55,
     "uncovered": 25,
     "coverable": 80,
     "total": 80,
     "percentage": 68.75
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_diagnostics_diagnostics.go.html",
   "diffStatus": "added"
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
     "covered": 6,
     "total": 13,
     "percentage": 46.15
    },
    "methods_hit": {
     "covered": 13,
     "total": 13,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 122,
     "uncovered": 11,
     "coverable": 133,
     "total": 133,
     "percentage": 91.72
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
     "covered": 6,
     "total": 12,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 12,
     "total": 12,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 113,
     "uncovered": 10,
     "coverable": 123,
     "total": 123,
     "percentage": 91.86
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_diff_parser.go.html"
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
    "statement_coverage": {
     "covered": 9,
     "uncovered": 1,
     "coverable": 10,
     "total": 10,
     "percentage": 90
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_diff_path.go.html"
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
     "value": 15
    },
    "methods_fully_covered": {
     "covered": 3,
     "total": 7,
     "percentage": 42.85
    },
    "methods_hit": {
     "covered": 7,
     "total": 7,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 148,
     "uncovered": 12,
     "coverable": 160,
     "total": 160,
     "percentage": 92.5
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
     "value": 13
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
     "covered": 29,
     "uncovered": 4,
     "coverable": 33,
     "total": 33,
     "percentage": 87.87
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_diffapply_apply.go.html"
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
     "covered": 119,
     "uncovered": 8,
     "coverable": 127,
     "total": 127,
     "percentage": 93.7
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_diffapply_resolver.go.html"
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
    "statement_coverage": {
     "covered": 94,
     "uncovered": 18,
     "coverable": 112,
     "total": 112,
     "percentage": 83.92
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
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
    "statement_coverage": {
     "covered": 94,
     "uncovered": 18,
     "coverable": 112,
     "total": 112,
     "percentage": 83.92
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_enricher_enricher.go.html"
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
     "total": 8,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 7,
     "total": 8,
     "percentage": 87.5
    },
    "statement_coverage": {
     "covered": 54,
     "uncovered": 7,
     "coverable": 61,
     "total": 61,
     "percentage": 88.52
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
    "methods_fully_covered": "safe",
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
     "total": 3,
     "percentage": 0
    },
    "methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 50,
     "uncovered": 6,
     "coverable": 56,
     "total": 56,
     "percentage": 89.28
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_filereader_filereader.go.html"
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
    "methods_fully_covered": "warning",
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
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_filtering_filter.go.html"
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
     "covered": 5,
     "total": 10,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 7,
     "total": 10,
     "percentage": 70
    },
    "statement_coverage": {
     "covered": 53,
     "uncovered": 10,
     "coverable": 63,
     "total": 63,
     "percentage": 84.12
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
     "covered": 5,
     "total": 10,
     "percentage": 50
    },
    "methods_hit": {
     "covered": 7,
     "total": 10,
     "percentage": 70
    },
    "statement_coverage": {
     "covered": 53,
     "uncovered": 10,
     "coverable": 63,
     "total": 63,
     "percentage": 84.12
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
    "methods_fully_covered": "safe",
    "statement_coverage": "safe"
   }
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
    "methods_fully_covered": "safe",
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
     "value": 22
    },
    "methods_fully_covered": {
     "covered": 19,
     "total": 37,
     "percentage": 51.35
    },
    "methods_hit": {
     "covered": 34,
     "total": 37,
     "percentage": 91.89
    },
    "statement_coverage": {
     "covered": 265,
     "uncovered": 56,
     "coverable": 321,
     "total": 321,
     "percentage": 82.55
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
     "total": 10,
     "percentage": 30
    },
    "methods_hit": {
     "covered": 10,
     "total": 10,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 77,
     "uncovered": 23,
     "coverable": 100,
     "total": 100,
     "percentage": 77
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
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
    "methods_fully_covered": "danger",
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
     "total": 4,
     "percentage": 25
    },
    "methods_hit": {
     "covered": 4,
     "total": 4,
     "percentage": 100
    },
    "statement_coverage": {
     "covered": 40,
     "uncovered": 11,
     "coverable": 51,
     "total": 51,
     "percentage": 78.43
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_cobertura_processing.go.html"
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
     "value": 14
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
     "covered": 54,
     "uncovered": 4,
     "coverable": 58,
     "total": 58,
     "percentage": 93.1
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
    "methods_fully_covered": "warning",
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
     "value": 14
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
     "covered": 35,
     "uncovered": 1,
     "coverable": 36,
     "total": 36,
     "percentage": 97.22
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_gcov_processing.go.html"
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
    "statement_coverage": {
     "covered": 63,
     "uncovered": 6,
     "coverable": 69,
     "total": 69,
     "percentage": 91.3
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
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
    "methods_fully_covered": "warning",
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
    "statement_coverage": {
     "covered": 29,
     "uncovered": 0,
     "coverable": 29,
     "total": 29,
     "percentage": 100
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_gocover_processing.go.html"
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
     "value": 22
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
    "statement_coverage": {
     "covered": 60,
     "uncovered": 19,
     "coverable": 79,
     "total": 79,
     "percentage": 75.94
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
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
    "methods_fully_covered": "warning",
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
     "value": 22
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
     "covered": 52,
     "uncovered": 4,
     "coverable": 56,
     "total": 56,
     "percentage": 92.85
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_parsers_parser_lcov_processing.go.html"
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
    "methods_fully_covered": "safe",
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
    "methods_fully_covered": "safe",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_parsers_parser_config.go.html"
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
     "value": 31
    },
    "methods_fully_covered": {
     "covered": 18,
     "total": 65,
     "percentage": 27.69
    },
    "methods_hit": {
     "covered": 41,
     "total": 65,
     "percentage": 63.07
    },
    "patch_methods_hit": {
     "covered": 14,
     "total": 32,
     "percentage": 43.75
    },
    "patch_statement_coverage": {
     "covered": 127,
     "uncovered": 125,
     "coverable": 252,
     "total": 252,
     "percentage": 50.39
    },
    "statement_coverage": {
     "covered": 490,
     "uncovered": 311,
     "coverable": 801,
     "total": 801,
     "percentage": 61.17
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
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
    "patch_methods_hit": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 32,
     "coverable": 32,
     "total": 32,
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
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
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
    "patch_methods_hit": {
     "covered": 0,
     "total": 5,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 32,
     "coverable": 32,
     "total": 32,
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
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_annotations_reporter.go.html",
   "diffStatus": "added"
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
     "value": 31
    },
    "methods_fully_covered": {
     "covered": 14,
     "total": 40,
     "percentage": 35
    },
    "methods_hit": {
     "covered": 32,
     "total": 40,
     "percentage": 80
    },
    "patch_methods_hit": {
     "covered": 11,
     "total": 16,
     "percentage": 68.75
    },
    "patch_statement_coverage": {
     "covered": 112,
     "uncovered": 43,
     "coverable": 155,
     "total": 155,
     "percentage": 72.25
    },
    "statement_coverage": {
     "covered": 387,
     "uncovered": 193,
     "coverable": 580,
     "total": 580,
     "percentage": 66.72
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
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
     "value": 20
    },
    "methods_fully_covered": {
     "covered": 4,
     "total": 17,
     "percentage": 23.52
    },
    "methods_hit": {
     "covered": 12,
     "total": 17,
     "percentage": 70.58
    },
    "patch_methods_hit": {
     "covered": 8,
     "total": 13,
     "percentage": 61.53
    },
    "patch_statement_coverage": {
     "covered": 88,
     "uncovered": 24,
     "coverable": 112,
     "total": 112,
     "percentage": 78.57
    },
    "statement_coverage": {
     "covered": 171,
     "uncovered": 63,
     "coverable": 234,
     "total": 234,
     "percentage": 73.07
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
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
     "value": 31
    },
    "methods_fully_covered": {
     "covered": 9,
     "total": 17,
     "percentage": 52.94
    },
    "methods_hit": {
     "covered": 16,
     "total": 17,
     "percentage": 94.11
    },
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 24,
     "uncovered": 19,
     "coverable": 43,
     "total": 43,
     "percentage": 55.81
    },
    "statement_coverage": {
     "covered": 167,
     "uncovered": 42,
     "coverable": 209,
     "total": 209,
     "percentage": 79.9
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "danger",
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
    "methods_fully_covered": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_htmlreact_embed.go.html"
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
    "methods_fully_covered": "danger",
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
    "methods_fully_covered": "danger",
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
    "methods_fully_covered": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_htmlreact_generator_single.go.html"
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
     "value": 22
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
     "covered": 69,
     "uncovered": 15,
     "coverable": 84,
     "total": 84,
     "percentage": 82.14
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
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
     "value": 22
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
     "covered": 69,
     "uncovered": 15,
     "coverable": 84,
     "total": 84,
     "percentage": 82.14
    }
   },
   "statuses": {
    "methods_fully_covered": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_reporter_lcov_reporter.go.html"
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
    "methods_fully_covered": "danger",
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
    "methods_fully_covered": "danger",
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
    "patch_methods_hit": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 33,
     "coverable": 33,
     "total": 33,
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
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
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
    "patch_methods_hit": {
     "covered": 0,
     "total": 6,
     "percentage": 0
    },
    "patch_statement_coverage": {
     "covered": 0,
     "uncovered": 33,
     "coverable": 33,
     "total": 33,
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
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_sarif_reporter.go.html",
   "diffStatus": "added"
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
     "value": 9
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
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "patch_statement_coverage": {
     "covered": 15,
     "uncovered": 17,
     "coverable": 32,
     "total": 32,
     "percentage": 46.87
    },
    "statement_coverage": {
     "covered": 26,
     "uncovered": 35,
     "coverable": 61,
     "total": 61,
     "percentage": 42.62
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   }
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
     "value": 9
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
     "covered": 3,
     "total": 5,
     "percentage": 60
    },
    "patch_statement_coverage": {
     "covered": 15,
     "uncovered": 17,
     "coverable": 32,
     "total": 32,
     "percentage": 46.87
    },
    "statement_coverage": {
     "covered": 26,
     "uncovered": 35,
     "coverable": 61,
     "total": 61,
     "percentage": 42.62
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "danger"
   },
   "targetUrl": "internal_reporter_textsummary_reporter.go.html",
   "diffStatus": "modified"
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
     "covered": 61,
     "uncovered": 13,
     "coverable": 74,
     "total": 74,
     "percentage": 82.43
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
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
     "covered": 61,
     "uncovered": 13,
     "coverable": 74,
     "total": 74,
     "percentage": 82.43
    }
   },
   "statuses": {
    "methods_fully_covered": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_review_review.go.html"
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
     "value": 6
    },
    "methods_fully_covered": {
     "covered": 96,
     "total": 135,
     "percentage": 71.11
    },
    "methods_hit": {
     "covered": 110,
     "total": 135,
     "percentage": 81.48
    },
    "patch_methods_hit": {
     "covered": 109,
     "total": 134,
     "percentage": 81.34
    },
    "patch_statement_coverage": {
     "covered": 196,
     "uncovered": 55,
     "coverable": 251,
     "total": 251,
     "percentage": 78.08
    },
    "statement_coverage": {
     "covered": 211,
     "uncovered": 55,
     "coverable": 266,
     "total": 266,
     "percentage": 79.32
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
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
     "covered": 89,
     "total": 126,
     "percentage": 70.63
    },
    "methods_hit": {
     "covered": 101,
     "total": 126,
     "percentage": 80.15
    },
    "patch_methods_hit": {
     "covered": 101,
     "total": 126,
     "percentage": 80.15
    },
    "patch_statement_coverage": {
     "covered": 157,
     "uncovered": 53,
     "coverable": 210,
     "total": 210,
     "percentage": 74.76
    },
    "statement_coverage": {
     "covered": 157,
     "uncovered": 53,
     "coverable": 210,
     "total": 210,
     "percentage": 74.76
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   }
  },
  {
   "id": "internal/status/evaluators/branch_coverage.go",
   "name": "branch_coverage.go",
   "type": "file",
   "path": "internal/status/evaluators/branch_coverage.go",
   "parentId": "internal/status/evaluators",
   "depth": 3,
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
    "patch_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 9,
     "uncovered": 1,
     "coverable": 10,
     "total": 10,
     "percentage": 90
    },
    "statement_coverage": {
     "covered": 9,
     "uncovered": 1,
     "coverable": 10,
     "total": 10,
     "percentage": 90
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_branch_coverage.go.html",
   "diffStatus": "added"
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
     "covered": 25,
     "total": 36,
     "percentage": 69.44
    },
    "methods_hit": {
     "covered": 30,
     "total": 36,
     "percentage": 83.33
    },
    "patch_methods_hit": {
     "covered": 30,
     "total": 36,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 47,
     "uncovered": 13,
     "coverable": 60,
     "total": 60,
     "percentage": 78.33
    },
    "statement_coverage": {
     "covered": 47,
     "uncovered": 13,
     "coverable": 60,
     "total": 60,
     "percentage": 78.33
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_complexity.go.html",
   "diffStatus": "added"
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
     "covered": 8,
     "total": 12,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 8,
     "total": 12,
     "percentage": 66.66
    },
    "patch_methods_hit": {
     "covered": 8,
     "total": 12,
     "percentage": 66.66
    },
    "patch_statement_coverage": {
     "covered": 12,
     "uncovered": 8,
     "coverable": 20,
     "total": 20,
     "percentage": 60
    },
    "statement_coverage": {
     "covered": 12,
     "uncovered": 8,
     "coverable": 20,
     "total": 20,
     "percentage": 60
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_line_coverage.go.html",
   "diffStatus": "added"
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
     "covered": 9,
     "uncovered": 1,
     "coverable": 10,
     "total": 10,
     "percentage": 90
    },
    "statement_coverage": {
     "covered": 9,
     "uncovered": 1,
     "coverable": 10,
     "total": 10,
     "percentage": 90
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_methods_fully_covered.go.html",
   "diffStatus": "added"
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
     "covered": 9,
     "uncovered": 1,
     "coverable": 10,
     "total": 10,
     "percentage": 90
    },
    "statement_coverage": {
     "covered": 9,
     "uncovered": 1,
     "coverable": 10,
     "total": 10,
     "percentage": 90
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_methods_hit.go.html",
   "diffStatus": "added"
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
     "covered": 8,
     "total": 12,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 8,
     "total": 12,
     "percentage": 66.66
    },
    "patch_methods_hit": {
     "covered": 8,
     "total": 12,
     "percentage": 66.66
    },
    "patch_statement_coverage": {
     "covered": 12,
     "uncovered": 8,
     "coverable": 20,
     "total": 20,
     "percentage": 60
    },
    "statement_coverage": {
     "covered": 12,
     "uncovered": 8,
     "coverable": 20,
     "total": 20,
     "percentage": 60
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "danger",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_patch_line_coverage.go.html",
   "diffStatus": "added"
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
     "covered": 4,
     "total": 6,
     "percentage": 66.66
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
     "covered": 7,
     "uncovered": 3,
     "coverable": 10,
     "total": 10,
     "percentage": 70
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
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_patch_methods_hit.go.html",
   "diffStatus": "added"
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
     "covered": 8,
     "total": 12,
     "percentage": 66.66
    },
    "methods_hit": {
     "covered": 10,
     "total": 12,
     "percentage": 83.33
    },
    "patch_methods_hit": {
     "covered": 10,
     "total": 12,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 14,
     "uncovered": 6,
     "coverable": 20,
     "total": 20,
     "percentage": 70
    },
    "statement_coverage": {
     "covered": 14,
     "uncovered": 6,
     "coverable": 20,
     "total": 20,
     "percentage": 70
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_patch_statement_coverage.go.html",
   "diffStatus": "added"
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
    "patch_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 7,
     "uncovered": 3,
     "coverable": 10,
     "total": 10,
     "percentage": 70
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
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_patch_statement_methods_hit.go.html",
   "diffStatus": "added"
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
     "covered": 9,
     "total": 12,
     "percentage": 75
    },
    "methods_hit": {
     "covered": 10,
     "total": 12,
     "percentage": 83.33
    },
    "patch_methods_hit": {
     "covered": 10,
     "total": 12,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 17,
     "uncovered": 3,
     "coverable": 20,
     "total": 20,
     "percentage": 85
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
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "warning",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_evaluators_statement_coverage.go.html",
   "diffStatus": "added"
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
    "patch_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 7,
     "uncovered": 3,
     "coverable": 10,
     "total": 10,
     "percentage": 70
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
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_statement_methods_fully_covered.go.html",
   "diffStatus": "added"
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
    "patch_methods_hit": {
     "covered": 5,
     "total": 6,
     "percentage": 83.33
    },
    "patch_statement_coverage": {
     "covered": 7,
     "uncovered": 3,
     "coverable": 10,
     "total": 10,
     "percentage": 70
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
    "methods_fully_covered": "safe",
    "patch_methods_hit": "warning",
    "patch_statement_coverage": "danger",
    "statement_coverage": "warning"
   },
   "targetUrl": "internal_status_evaluators_statement_methods_hit.go.html",
   "diffStatus": "added"
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
     "covered": 5,
     "total": 5,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 29,
     "uncovered": 2,
     "coverable": 31,
     "total": 31,
     "percentage": 93.54
    },
    "statement_coverage": {
     "covered": 39,
     "uncovered": 2,
     "coverable": 41,
     "total": 41,
     "percentage": 95.12
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_annotate.go.html",
   "diffStatus": "modified"
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
    "patch_methods_hit": {
     "covered": 3,
     "total": 3,
     "percentage": 100
    },
    "patch_statement_coverage": {
     "covered": 10,
     "uncovered": 0,
     "coverable": 10,
     "total": 10,
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
    "methods_fully_covered": "safe",
    "patch_methods_hit": "safe",
    "patch_statement_coverage": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_status_classifier.go.html",
   "diffStatus": "modified"
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
    "statement_coverage": {
     "covered": 87,
     "uncovered": 6,
     "coverable": 93,
     "total": 93,
     "percentage": 93.54
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
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
    "statement_coverage": {
     "covered": 87,
     "uncovered": 6,
     "coverable": 93,
     "total": 93,
     "percentage": 93.54
    }
   },
   "statuses": {
    "methods_fully_covered": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_tree_builder.go.html"
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
    "methods_fully_covered": "danger",
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
    "methods_fully_covered": "safe",
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
    "methods_fully_covered": "safe",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_utils_line_sorter.go.html"
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
    "methods_fully_covered": "danger",
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
    "methods_fully_covered": "danger",
    "statement_coverage": "safe"
   },
   "targetUrl": "internal_utils_paths.go.html"
  }
 ],
 "metricDefinitions": {
  "a_statement_coverage": {
   "label": "Statements",
   "shortLabel": "Statements",
   "subMetrics": [
    {
     "id": "total",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "branch_coverage": {
   "label": "Branches",
   "shortLabel": "Branches",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Covered",
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
  "c_patch_statement_coverage": {
   "label": "Patch Statements",
   "shortLabel": "Patch Stmts",
   "subMetrics": [
    {
     "id": "total",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "e_branch_coverage": {
   "label": "Branches",
   "shortLabel": "Branches",
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
   "subMetrics": [
    {
     "id": "total",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "h_patch_crap_score": {
   "label": "Patch CRAP Score",
   "shortLabel": "PCRAP",
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
   "subMetrics": [
    {
     "id": "total",
     "label": "Value",
     "width": 100
    }
   ]
  },
  "j_defect_probability": {
   "label": "Defect Probability",
   "shortLabel": "DPI",
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
   "kind": "value",
   "subMetrics": [
    {
     "id": "value",
     "label": "Value",
     "width": 140
    }
   ]
  },
  "method_branch_coverage": {
   "label": "Method Branches",
   "shortLabel": "Method Branches",
   "subMetrics": [
    {
     "id": "covered",
     "label": "Covered",
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
  "methods_fully_covered": {
   "label": "Methods Fully Covered",
   "shortLabel": "Fully Covered",
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
 "metadata": [
  {
   "label": "Generated At",
   "value": "2026-07-06 22:11:05"
  },
  {
   "label": "Parser",
   "value": "Cobertura | GCov | GoCover"
  }
 ],
 "diagnostics": [
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "cmd/main.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 46.1% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "error",
   "file": "cmd/main.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 66.7% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "cmd/main.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 45.1% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "cmd/main.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 61.5% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "branch_coverage",
   "ruleName": "Branch Coverage",
   "severity": "warning",
   "file": "demo_projects/cpp/project/src/advanced_calculator.cpp",
   "startLine": 1,
   "endLine": 1,
   "message": "Branch Coverage 50.0% (threshold 50..70)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "demo_projects/cpp/project/src/advanced_calculator.cpp",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "demo_projects/cpp/project/src/advanced_calculator.cpp",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 75.0% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "demo_projects/cpp/project/src/calculator.cpp",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 60.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "branch_coverage",
   "ruleName": "Branch Coverage",
   "severity": "warning",
   "file": "demo_projects/cpp/project/src/utils/math_utils.cpp",
   "startLine": 1,
   "endLine": 1,
   "message": "Branch Coverage 64.3% (threshold 50..70)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "demo_projects/cpp/project/src/utils/math_utils.cpp",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "branch_coverage",
   "ruleName": "Branch Coverage",
   "severity": "warning",
   "file": "demo_projects/csharp/project/Test/PartialClass.cs",
   "startLine": 1,
   "endLine": 1,
   "message": "Branch Coverage 50.0% (threshold 50..70)",
   "scope": "file"
  },
  {
   "ruleId": "branch_coverage",
   "ruleName": "Branch Coverage",
   "severity": "warning",
   "file": "demo_projects/csharp/project/Test/TestClass.cs",
   "startLine": 1,
   "endLine": 1,
   "message": "Branch Coverage 50.0% (threshold 50..70)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "demo_projects/go/project/calculator/calculator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 42.9% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "demo_projects/go/project/calculator/calculator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 85.7% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "demo_projects/go/project/calculator/calculator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 76.9% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "demo_projects/go/project/calculator/entities.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 75.0% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "demo_projects/go/project/calculator_2/calculator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 40.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "demo_projects/go/project/calculator_2/calculator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 66.7% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "demo_projects/go/project/calculator_2/entities.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 16.7% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "demo_projects/go/project/calculator_2/entities.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 25.0% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/analyzer/cpp/analyzer.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/analyzer/gdscript/analyzer.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 40.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/analyzer/go/analyzer.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 60.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/bootlog/bootlog.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/cache/cache.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/cache/cache_validator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "internal/cache/cache_validator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 0.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/calculator/calculators.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 38.6% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "error",
   "file": "internal/calculator/calculators.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 56.1% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/calculator/calculators.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 54.2% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "internal/calculator/calculators.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 54.2% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/config/config.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 33.3% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/config/config.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 74.7% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/diagnostics/diagnostics.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 27.3% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/diagnostics/diagnostics.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 81.8% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/diagnostics/diagnostics.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 68.8% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/diagnostics/diagnostics.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 68.8% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/diff/parser.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/diff/path.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/diffapply/apply.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/diffapply/resolver.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 40.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/filereader/filereader.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/filtering/filter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/logging/logging.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/parsers/factory.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 75.0% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/parsers/parser_cobertura/parser.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 33.3% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/parsers/parser_cobertura/processing.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 25.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/parsers/parser_config.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 66.7% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/parsers/parser_gcov/parser.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/parsers/parser_gcov/processing.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/parsers/parser_gocover/parser.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 40.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/parsers/parser_lcov/parser.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "internal/parsers/parser_lcov/parser.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 34.8% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/parsers/parser_lcov/processing.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 50.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/reporter/annotations/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "error",
   "file": "internal/reporter/annotations/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 0.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/annotations/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 0.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/annotations/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 0.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/reporter/htmlreact/builder.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 23.5% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "error",
   "file": "internal/reporter/htmlreact/builder.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 61.5% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/htmlreact/builder.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 78.6% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/reporter/htmlreact/builder.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 73.1% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/reporter/htmlreact/details_generator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 52.9% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/htmlreact/details_generator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 55.8% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/reporter/htmlreact/emit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/reporter/htmlreact/generator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/reporter/htmlreact/generator.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 68.0% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/reporter/htmlreact/generator_single.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/htmlreact/generator_single.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 0.0% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "warning",
   "file": "internal/reporter/lcov/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 40.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/reporter/reporter_rawjson/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 33.3% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/reporter/reporter_rawjson/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 72.7% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/reporter/sarif/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "error",
   "file": "internal/reporter/sarif/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 0.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/sarif/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 0.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/sarif/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 0.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/reporter/textsummary/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 16.7% (threshold 40..60)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "error",
   "file": "internal/reporter/textsummary/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 60.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/textsummary/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 46.9% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "internal/reporter/textsummary/reporter.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 42.6% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "modified"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/review/review.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 16.7% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/branch_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/branch_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 90.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/complexity.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/status/evaluators/complexity.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 78.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "error",
   "file": "internal/status/evaluators/line_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 66.7% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/status/evaluators/line_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 60.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/line_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 60.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/methods_fully_covered.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/methods_fully_covered.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 90.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 90.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "error",
   "file": "internal/status/evaluators/patch_line_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 66.7% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/status/evaluators/patch_line_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 60.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/patch_line_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 60.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/patch_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/status/evaluators/patch_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 70.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/patch_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 70.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/patch_statement_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/status/evaluators/patch_statement_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 70.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/patch_statement_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 70.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/patch_statement_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/status/evaluators/patch_statement_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 70.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/patch_statement_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 70.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/statement_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/statement_coverage.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 85.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/statement_methods_fully_covered.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/status/evaluators/statement_methods_fully_covered.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 70.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/statement_methods_fully_covered.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 70.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_methods_hit",
   "ruleName": "Patch Methods Hit",
   "severity": "warning",
   "file": "internal/status/evaluators/statement_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Methods Hit 83.3% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "patch_statement_coverage",
   "ruleName": "Patch Statement Coverage",
   "severity": "error",
   "file": "internal/status/evaluators/statement_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Patch Statement Coverage 70.0% (threshold 80..90)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "warning",
   "file": "internal/status/evaluators/statement_methods_hit.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 70.0% (threshold 60..75)",
   "scope": "file",
   "diffStatus": "added"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/utils/math.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  },
  {
   "ruleId": "statement_coverage",
   "ruleName": "Statement Coverage",
   "severity": "error",
   "file": "internal/utils/math.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Statement Coverage 50.0% (threshold 60..75)",
   "scope": "file"
  },
  {
   "ruleId": "methods_fully_covered",
   "ruleName": "Methods Fully Covered",
   "severity": "error",
   "file": "internal/utils/paths.go",
   "startLine": 1,
   "endLine": 1,
   "message": "Methods Fully Covered 0.0% (threshold 40..60)",
   "scope": "file"
  }
 ],
 "review": {
  "passed": false,
  "checks": [
   {
    "key": "patch_statement_coverage",
    "label": "Patch statement coverage",
    "value": 66.32124352331606,
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
   "changedFiles": 29,
   "methodsAdded": 228,
   "methodsModified": 32,
   "untestedChangedMethods": 74,
   "patchStatementsValid": 965,
   "patchStatementsCovered": 640,
   "maxChangedComplexity": 31
  },
  "hotspots": [
   {
    "file": "cmd/main.go",
    "method": "main",
    "startLine": 392,
    "diffStatus": "modified",
    "complexity": 28,
    "patchCoverage": 50,
    "risk": 13.740740740740742
   },
   {
    "file": "internal/reporter/htmlreact/details_generator.go",
    "method": "(*HtmlReactReportBuilder).buildMethodDetails",
    "startLine": 221,
    "diffStatus": "modified",
    "complexity": 31,
    "patchCoverage": 50,
    "risk": 12.020408163265305
   },
   {
    "file": "cmd/main.go",
    "method": "evaluateReviewGate",
    "startLine": 340,
    "diffStatus": "added",
    "complexity": 11,
    "patchCoverage": 17.647058823529413,
    "risk": 9.058823529411764
   },
   {
    "file": "internal/calculator/calculators.go",
    "method": "(MethodDefectProbabilityCalculator).Calculate",
    "startLine": 383,
    "diffStatus": "added",
    "complexity": 10,
    "patchCoverage": 27.77777777777778,
    "risk": 7.222222222222222
   },
   {
    "file": "internal/reporter/htmlreact/builder.go",
    "method": "(*HtmlReactReportBuilder).createSingleFileReport",
    "startLine": 94,
    "diffStatus": "modified",
    "complexity": 7,
    "patchCoverage": 0,
    "risk": 7
   },
   {
    "file": "internal/config/config.go",
    "method": "(*AppConfig).mergeCliOverrides",
    "startLine": 290,
    "diffStatus": "modified",
    "complexity": 28,
    "patchCoverage": 85.71428571428571,
    "risk": 6.3859649122807
   },
   {
    "file": "internal/reporter/annotations/reporter.go",
    "method": "(*builder).CreateReport",
    "startLine": 34,
    "diffStatus": "added",
    "complexity": 6,
    "patchCoverage": 0,
    "risk": 6
   },
   {
    "file": "internal/reporter/textsummary/reporter.go",
    "method": "(*TextReportBuilder).buildNodeParts",
    "startLine": 98,
    "diffStatus": "added",
    "complexity": 6,
    "patchCoverage": 0,
    "risk": 6
   },
   {
    "file": "internal/reporter/textsummary/reporter.go",
    "method": "(*TextReportBuilder).printNode",
    "startLine": 117,
    "diffStatus": "modified",
    "complexity": 5,
    "patchCoverage": 0,
    "risk": 5
   },
   {
    "file": "cmd/main.go",
    "method": "setupCacheManager",
    "startLine": 232,
    "diffStatus": "modified",
    "complexity": 6,
    "patchCoverage": 0,
    "risk": 4.9411764705882355
   }
  ]
 }
}
