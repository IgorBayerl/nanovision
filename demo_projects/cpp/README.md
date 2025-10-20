# C++ Example: Simple Calculator

A minimal, self‑contained calculator library written in modern C++ (C++17).
It demonstrates a **clean CMake layout, Google Test integration, and multiple
coverage‑report formats**.


---

## Requirements

| Platform          | Required tools                                                                                                      | Notes                                                                                              |
| ----------------- | ------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| **Windows 10/11** | PowerShell 5.1 ×64, [MinGW](https://www.mingw-w64.org/), [CMake ≥ 3.14](https://cmake.org/), **Python 3** + `gcovr` | `setup.ps1` installs MinGW, CMake, Python and gcovr via Chocolatey/pip (run **as Administrator**). |
| **Linux / macOS** | -| TODO: Create `setup.sh`                                                       |

---

## Quick Start (Windows)

```powershell
# one‑time tool‑chain prep (CMake, MinGW, Python, gcovr)
PS> ./setup.ps1    # run from an elevated x64 PowerShell session

# pick a coverage script
PS> ./gen-gcov.ps1           
PS> ./gen-html.ps1            
PS> ./gen-cobertura.ps1   
```


---

## Coverage Reports

Three helper scripts wrap *build → test → report* for different audiences:

| Script                                             | Format                                                                  
| -----------------------  | ----------------------------------- | 
| `gen-gcov`          | raw **.gcov** text (basic, branch‑counts, probabilities, unconditional) 
| `gen-html`          | **HTML** site                                          
| `gen-cobertura`     | **Cobertura XML**                                                       

### `.gcov` Breakdown

`gen-gcov` creates four granular sub‑folders:

| Sub‑folder                | gcov flags | What it shows             |
| ------------------------- | ---------- | ------------------------- |
| `basic/`                  | *(none)*   | line execution counts     |
| `branch‑probabilities/`   | `-b`       | line & branch hit ratios  |
| `branch‑counts/`          | `-b -c`    | branch hit *counts*       |
| `unconditional‑branches/` | `-b -c -u` | unconditional branch hits |

---

## Project Layout (abridged)

```
project/
 ├── src/              # library sources
 ├── tests/            # Google Test suite
 └── CMakeLists.txt    # build script
report/
 ├── gcov/             # raw .gcov outputs (gen-gcov)
 ├── html/             # HTML site (gen-html)
 └── cobertura/        # cobertura.xml (gen-cobertura-xml)
setup.ps1 / setup.sh   # environment boot‑strap
gen-*.ps1 / .sh        # build + coverage helpers
```

---

## Next Steps

* Implement .sh scripts for linux setup