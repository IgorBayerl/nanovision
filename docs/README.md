# 📝 Project Documentation

This folder contains the documentation site for the **nanovision** project, built with [MkDocs](https://www.mkdocs.org/) and the [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/) theme.

---

## 📦 Stack Overview

| Tool                                                                | Description                                                 |
|---------------------------------------------------------------------|-------------------------------------------------------------|
| [MkDocs](https://www.mkdocs.org/)                                   | Static site generator geared towards project documentation. |
| [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/) | A MkDocs theme.                                             |

---

## 🚀 Getting Started

### Install dependencies

From the project root or the `docs/` directory:

```bash
pip install -r requirements.txt
```

---

## 🛠 Usage

### Run local development server

```bash
python -m mkdocs serve
```

Open your browser and go to [http://localhost:8000](http://localhost:8000)

### Build static site

```bash
python -m mkdocs build --strict
```

The generated static site will be in the `site/` directory.

---

## 📚 References

* [MkDocs documentation](https://www.mkdocs.org/)
* [Material for MkDocs documentation](https://squidfunk.github.io/mkdocs-material/)

---

## 🧼 Git Ignore Note

This folder includes a `.gitignore` to exclude the generated `site/` directory, which is automatically created during local builds or CI deployments.

