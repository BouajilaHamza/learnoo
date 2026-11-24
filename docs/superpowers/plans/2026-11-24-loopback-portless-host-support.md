# Loopback Portless Host Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Allow portless Host headers (`localhost`, `127.0.0.1`, `[::1]`) in the renderer loopback security check so developers, curl, and standard reverse proxies are not rejected with HTTP 403 Forbidden.

**Architecture:** Use `net.SplitHostPort` to extract the hostname when a port is present, but cleanly fall back to `r.Host` when `net.ErrMissingPort` occurs, verifying that the host is strictly an authorized loopback address (`127.0.0.1`, `localhost`, or `::1`) before allowing the request to proceed.

**Tech Stack:** Go (Standard library: `net`, `net/http`, `net/http/httptest`).

**Spec:** [`docs/audits/issue-02-loopback-portless-host-rejection.md`](file:///home/hamza/firstmate/projects/learnoo/docs/audits/issue-02-loopback-portless-host-rejection.md).

## Global Constraints
- Do not bypass DNS rebinding protections: non-loopback domains (e.g. `evil.example`, `attacker.com:8080`) must remain strictly forbidden (HTTP 403).
- Zero external dependencies: rely strictly on Go standard library.
- Zero AI attribution trailers or footers in git commit messages (`scripts/check-attribution.sh` must exit 0).

## Review Focus
1. `Host: localhost` (portless) must return 200 OK, not 403 Forbidden.
2. `Host: 127.0.0.1` (portless) must return 200 OK, not 403 Forbidden.
3. `Host: localhost:8080` (with port) must continue returning 200 OK.
4. `Host: evil.example` (portless rebind attack) must return 403 Forbidden.
5. `Host: [::1]:8080` and `Host: [::1]` (IPv6 loopback) behavior must be explicit and safe.

---

### Task 1: Add Failing Unit Tests for Portless Loopback Hosts

**Files:**
- Modify: `cmd/renderer/main_test.go`

- [ ] **Step 1: Write failing test cases in `cmd/renderer/main_test.go`**
  Add test cases to `TestPreviewMatchesReader` or a new test `TestLoopbackHostAllowed` verifying that `GET /lesson` with `Host: localhost` and `Host: 127.0.0.1` succeeds with HTTP 200:
  ```go
  func TestLoopbackHostAllowed(t *testing.T) {
      h := handler()
      for _, host := range []string{"127.0.0.1:8080", "localhost:8080", "127.0.0.1", "localhost"} {
          t.Run(host, func(t *testing.T) {
              r := request("GET", "/lesson", nil)
              r.Host = host
              w := httptest.NewRecorder()
              h.ServeHTTP(w, r)
              if w.Code != http.StatusOK {
                  t.Fatalf("host %q got %d, want 200", host, w.Code)
              }
          })
      }
  }
  ```
- [ ] **Step 2: Run test to verify failure on portless hosts**
  Run: `go test -v ./cmd/renderer -run TestLoopbackHostAllowed`
  Confirm failure on `"127.0.0.1"` and `"localhost"` with error `got 403, want 200`.

---

### Task 2: Implement Portless Loopback Extraction

**Files:**
- Modify: `cmd/renderer/main.go:76-85`

- [ ] **Step 1: Replace brittle `net.SplitHostPort` with safe host parser**
  Update the loopback check in `cmd/renderer/main.go`:
  ```go
  host := r.Host
  if h, _, err := net.SplitHostPort(r.Host); err == nil {
      host = h
  }
  if host != "127.0.0.1" && host != "localhost" && host != "::1" && host != "[::1]" {
      http.Error(w, "Loopback host required", http.StatusForbidden)
      return
  }
  ```
- [ ] **Step 2: Run unit tests to verify fix**
  Run: `go test -v ./cmd/renderer`
  Confirm all tests in `cmd/renderer` pass cleanly.
- [ ] **Step 3: Run regression checks**
  Run `go test -race ./cmd/renderer` and verify rebinding attack test (`evil.example:8080` and `evil.example`) still correctly returns 403.
- [ ] **Step 4: Commit changes**
  Commit with clean attribution:
  `git commit -m "fix(renderer): support portless loopback Host headers"`
  Verify with `bash scripts/check-attribution.sh upstream/develop HEAD`.
