#!/usr/bin/env python3
"""
Regenerates the UI dev fixture (ui/public/data.js) from real nanovision
self-coverage data.

The fixture is what the Vite dev server (pnpm dev in ui/) renders; in
production the Go reporter writes a real data.js next to index.html. This
script produces a fixture that exercises everything the summary page can
show: the full merged self-coverage tree, patch coverage from a git diff,
the Problems panel, and the review header (gate verdict, changelist stats,
risk hotspots) from the HtmlReview report type.

Prerequisites: the self-coverage input files must exist (run
`python scripts/e2e_test.py --self-cover` first, or have
reports/nanovision_self_coverage/*.out from a previous run).

Usage:
    python scripts/gen_ui_fixture.py [--diff-target HEAD~5]
"""
import argparse
import json
import os
import subprocess
import sys
import tempfile

SCRIPT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
FIXTURE_PATH = os.path.join(SCRIPT_ROOT, "ui", "public", "data.js")

HEADER = """/**
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
"""


def run(cmd, **kwargs):
    print(f"--- Running: {' '.join(cmd)}")
    return subprocess.run(cmd, cwd=SCRIPT_ROOT, check=True, **kwargs)


def extract_summary_from_data_js(path):
    """data.js is `window.__NANOVISION_SUMMARY__=<json>`."""
    with open(path, encoding="utf-8") as f:
        content = f.read()
    return json.loads(content[content.index("=") + 1 :].rstrip().rstrip(";"))


def extract_review_block(review_index_html):
    """The single-file review report embeds `__NANOVISION_FULL_DATA__ = <json>;`."""
    with open(review_index_html, encoding="utf-8") as f:
        content = f.read()
    marker = "__NANOVISION_FULL_DATA__ = "
    start = content.index(marker) + len(marker)
    end = content.index(";</script>", start)
    full_data = json.loads(content[start:end])
    review = full_data["summary"].get("review")
    if not review:
        raise SystemExit("HtmlReview output has no 'review' block; aborting.")
    return review


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--diff-target", default="HEAD~5", help="Git ref to diff against (default: HEAD~5)")
    args = parser.parse_args()

    with tempfile.TemporaryDirectory(prefix="nanovision_fixture_") as tmp:
        diff_path = os.path.join(tmp, "fixture.diff")
        out_dir = os.path.join(tmp, "out")

        with open(diff_path, "w", encoding="utf-8") as f:
            subprocess.run(["git", "diff", args.diff_target], cwd=SCRIPT_ROOT, check=True, stdout=f, text=True)

        run(
            [
                "go", "run", "./cmd",
                "-config", "nanovision.yaml",
                "-reporttypes", "Html,HtmlReview",
                "-diff", diff_path,
                "-output", out_dir,
                "-verbosity", "Warning",
            ]
        )

        summary = extract_summary_from_data_js(os.path.join(out_dir, "data.js"))
        summary["review"] = extract_review_block(os.path.join(out_dir, "review", "index.html"))
        summary["title"] = "nanovision Self-Coverage (dev fixture)"

    with open(FIXTURE_PATH, "w", encoding="utf-8", newline="\n") as f:
        f.write(HEADER)
        f.write("window.__NANOVISION_SUMMARY__ = ")
        json.dump(summary, f, indent=1)
        f.write("\n")

    size_kb = os.path.getsize(FIXTURE_PATH) // 1024
    review = summary["review"]
    print(f"Wrote {FIXTURE_PATH} ({size_kb} KB)")
    print(f"review.passed={review['passed']} checks={len(review.get('checks', []))} hotspots={len(review.get('hotspots', []))}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
