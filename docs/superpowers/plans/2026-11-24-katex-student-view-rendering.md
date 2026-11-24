# KaTeX and Go Student View Rendering Implementation Plan (LRN-004)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement automated unit and end-to-end browser test suites verifying Go server HTML generation and client-side KaTeX rendering across mobile viewports, bidirectional text (Arabic/French), formulas, and media.

**Architecture:** Combine Go integration testing (`internal/rendering/render_test.go`) with Playwright browser testing (`tests/browser/renderer.spec.mjs`), validating HTML DOM output, KaTeX formula formatting, CSS layout stability at 375px viewport, and payload asset sizes.

**Tech Stack:** Go, HTML5/CSS3, KaTeX, Playwright (Node.js).

**Spec:** [`https://github.com/hadden-me/learnoo/issues/4`](https://github.com/hadden-me/learnoo/issues/4), [`docs/RENDERER_EXPERIMENT.md`](file:///home/hamza/firstmate/projects/learnoo/docs/RENDERER_EXPERIMENT.md).

## Global Constraints
- Zero external CDN dependencies in tests: KaTeX and fonts must be served locally from `internal/rendering/static/`.
- Strict CSP: script-src, style-src, and frame-src policies enforced.
- Mobile viewport: verified on 375x667 (iPhone SE / standard Android screen).
- Zero AI attribution trailers in commits (`scripts/check-attribution.sh`).

## Review Focus
1. Mixed direction: Arabic text with embedded inline KaTeX formulas must not reverse formula operands.
2. Display math (`\operatorname{Re}(z)=a \qquad \operatorname{Im}(z)=b`) must center without overflowing 375px screen width.
3. YouTube video container must maintain 16:9 aspect ratio and render placeholder/consent button.
4. Exercise callout blocks must expand/collapse solutions without layout jump.
5. Malformed math syntax must degrade gracefully without crashing the Go server or freezing browser DOM.

---

### Task 1: Go Server-Side Render Verification Tests

**Files:**
- Modify: `internal/rendering/render_test.go`
- Test fixtures: `internal/rendering/fixtures/complex-numbers.json`

- [ ] **Step 1: Write test for bidirectional formula wrapping**
  Add unit test ensuring `render.go` outputs appropriate `dir="auto"` or `dir="rtl"` wrapping for Arabic paragraphs containing `<span class="math">` tags.
- [ ] **Step 2: Run test and verify Go output**
  Run `go test -v ./internal/rendering -run TestRenderFormulas`.
- [ ] **Step 3: Benchmark rendering latency and memory**
  Run `go test -bench=BenchmarkRender ./internal/rendering` and assert rendering is under 5ms per document.

---

### Task 2: Playwright Mobile Browser & KaTeX Verification

**Files:**
- Modify: `tests/browser/renderer.spec.mjs`
- Modify: `playwright.config.mjs`

- [ ] **Step 1: Write Playwright test for mobile viewport rendering (375px)**
  Add test suite in `tests/browser/renderer.spec.mjs`:
  - Set viewport to `{ width: 375, height: 667 }`.
  - Navigate to `/preview` with `complex-numbers.json`.
  - Assert `.katex-display` elements do not cause horizontal scrolling (`scrollWidth <= clientWidth`).
  - Assert KaTeX rendered DOM contains `.katex-html` with math symbols ($\mathbb{R}$, $i^2 = -1$).
- [ ] **Step 2: Measure asset sizes**
  Assert total downloaded bundle size (HTML + CSS + local KaTeX fonts) is within acceptable bounds (< 300 KB gzipped).
- [ ] **Step 3: Run Playwright test suite**
  Run: `npx playwright test tests/browser/renderer.spec.mjs`
  Confirm all assertions pass on Chromium, Firefox, and WebKit mobile viewports.

---

### Task 3: Document Findings and Evidence

**Files:**
- Modify: `docs/RENDERER_EXPERIMENT.md`

- [ ] **Step 1: Update experiment report**
  Record measured latency, asset sizes, and mobile rendering fidelity notes.
- [ ] **Step 2: Commit**
  `git commit -m "test(renderer): add mobile viewport and KaTeX rendering verification"`
