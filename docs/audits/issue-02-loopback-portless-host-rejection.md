# Issue: Loopback Guard Fails on Portless Host Headers (HTTP 403)

## Status
- **Severity**: Medium (Availability / Reverse Proxy Incompatibility)
- **Affects**: `cmd/renderer/main.go` (line 79)
- **Found**: 2026-11-24 Adversarial Audit

---

## Summary
The loopback guard in `cmd/renderer/main.go` uses `net.SplitHostPort(r.Host)` without handling hosts that lack an explicit port. Any client or reverse proxy sending `Host: localhost` or `Host: 127.0.0.1` (standard for port 80/443 termination) receives an immediate HTTP 403 Forbidden.

## Vulnerability Mechanics
In Go standard library, `net.SplitHostPort("localhost")` returns an error (`missing port in address`).
In `cmd/renderer/main.go`:
```go
host, _, hostErr := net.SplitHostPort(r.Host)
if hostErr != nil || (host != "127.0.0.1" && host != "localhost") {
    http.Error(w, "Loopback host required", 403)
    return
}
```
Because `hostErr != nil`, valid requests from `curl http://localhost/` or local proxies without port numbers are rejected with 403.

## Proposed Remediation
Handle portless addresses:
```go
host := r.Host
if h, _, err := net.SplitHostPort(r.Host); err == nil {
    host = h
}
if host != "127.0.0.1" && host != "localhost" {
    http.Error(w, "Loopback host required", 403)
    return
}
```
