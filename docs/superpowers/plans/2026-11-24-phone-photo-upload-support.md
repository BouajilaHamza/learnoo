# Phone-Photo Upload Support Implementation Plan (Issue #42)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Safely accept phone photos (Android Ultra HDR, Motion Photos, CIPA MPF previews, EXIF-oriented photos) in the `internal/media` package, auto-orienting them upright and stripping trailing metadata/video streams while strictly rejecting polyglots and arbitrary garbage.

**Architecture:** Extend JPEG decoding in `internal/media`:
1. Parse JPEG markers to find EOI (`FF D9`).
2. If bytes exist after EOI, structurally inspect trailing data: allow known phone markers (CIPA MPF secondary JPEGs, `ftyp` MP4 atoms for Motion Photos, zero padding), stripping them safely. Reject arbitrary bytes as `ErrTrailingData`.
3. Extract EXIF Orientation (tags 1-8) from the APP1 segment before re-encoding, and apply standard affine pixel rotation/flipping so stored output is always oriented upright.
4. Re-encode to clean JPEG (quality 90) or PNG without any EXIF/GPS/device metadata.

**Tech Stack:** Go standard library (`image`, `image/jpeg`, `image/png`), `golang.org/x/image`.

**Spec:** [`https://github.com/hadden-me/learnoo/issues/42`](https://github.com/hadden-me/learnoo/issues/42).

## Global Constraints
- Maximum upload limit remains 5 MiB (`MaxUploadBytes`).
- Dimensions remain bounded by 4096 px per side and 16 megapixels (`MaxPixels`).
- Memory bound: concurrency semaphore (`MaxConcurrentDecodes = 2`) preserved.
- Polyglot rejection: arbitrary HTML, JS, or executable bytes appended to images must be rejected (`ErrTrailingData`).
- Output storage: stored assets are strictly content-addressed by sha256 (`<sha256>.jpg` / `<sha256>.png`).
- Zero AI attribution trailers in commits (`scripts/check-attribution.sh`).

## Review Focus
1. Photo taken in portrait mode with EXIF orientation 6 (90° CW) must be rotated upright in stored pixels.
2. Android Ultra HDR JPEG (with embedded gain map) must be accepted, stripped, and re-encoded upright.
3. Samsung/Pixel Motion Photo (with embedded trailing MP4 container) must be accepted with video discarded safely.
4. Trailing zero padding (common in camera flash write buffers) must be accepted and stripped.
5. Malicious polyglot (`<script>` or zip payload after EOI) must continue returning `ErrTrailingData`.

---

### Task 1: Add EXIF Orientation Normalization

**Files:**
- Modify: `internal/media/media.go`
- Test: `internal/media/media_test.go`

- [ ] **Step 1: Write failing test with all 8 EXIF orientation fixtures**
  Generate synthetic 100x200 asymmetric test JPEGs tagged with EXIF orientation 1 through 8 in `internal/media/media_test.go`. Verify that after `Process(ctx, r, "image/jpeg")`, decoded dimensions and pixel positions correspond to orientation 1 (upright).
- [ ] **Step 2: Run test to observe failure**
  Run: `go test -v ./internal/media -run TestExifOrientation` (fails because current code ignores orientation).
- [ ] **Step 3: Implement EXIF orientation parser & pixel rotator**
  Parse APP1 Exif header for TIFF tag `0x0112`. Implement `applyOrientation(img image.Image, orientation int) image.Image` handling cases 1-8 (rotate 90°, 180°, 270°, and flips).
- [ ] **Step 4: Verify test passes**
  Run: `go test -v ./internal/media -run TestExifOrientation`.

---

### Task 2: Structural Trailing-Data Parser for Phone Photos

**Files:**
- Modify: `internal/media/media.go`
- Test: `internal/media/media_test.go`

- [ ] **Step 1: Write failing tests for permitted phone structures and forbidden polyglots**
  Create test cases:
  - JPEG + CIPA MPF preview JPEG (starts with `FF D8`) -> want Success.
  - JPEG + Motion photo MP4 (starts with `....ftypmp42` or `ftypisom`) -> want Success.
  - JPEG + trailing zero padding (`\x00` bytes) -> want Success.
  - JPEG + arbitrary ASCII/garbage -> want `ErrTrailingData`.
- [ ] **Step 2: Run test to observe failure**
  Run: `go test -v ./internal/media -run TestTrailingPhoneData`.
- [ ] **Step 3: Implement structural trailing parser**
  Replace naive `len(trailing) > 0` check with `parsePermittedTrailingData(trailing []byte) bool`:
  Verify trailing bytes belong exclusively to an allowed signature (MPF JPEG, MP4 container, or all zeros). If so, log/drop trailing slice and proceed to decode primary image.
- [ ] **Step 4: Verify tests pass**
  Run: `go test -v ./internal/media`.

---

### Task 3: Regression, Concurrency, and Memory Verification

**Files:**
- Test: `internal/media/media_test.go`

- [ ] **Step 1: Run race detector and memory checks**
  Run `go test -race ./internal/media` to ensure no race conditions under concurrent decodes.
- [ ] **Step 2: Verify polyglot security assertions**
  Ensure existing `TestPolyglot` and `TestOversized` tests still pass without regression.
- [ ] **Step 3: Commit**
  `git commit -m "feat(media): support phone photo uploads with EXIF orientation and trailing data stripping"`
