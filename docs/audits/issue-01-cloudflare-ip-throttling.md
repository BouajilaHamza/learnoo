# Issue: Cloudflare Reverse Proxy IP Spoofing and DoS on Login Rate Limiting

## Status
- **Severity**: High (Authentication / Rate Limiting DoS)
- **Affects**: `internal/auth/http.go` (`clientAddr`), `internal/auth/service.go` (`Login`)
- **Found**: 2026-11-24 Adversarial Audit

---

## Summary
When deployed behind Cloudflare (documented as Learnoo's planned reverse proxy in `docs/CAHIER_DES_CHARGES.md §1`), `net.SplitHostPort(r.RemoteAddr)` resolves to Cloudflare's egress IP address rather than the real client IP.

## Vulnerability Mechanics
In `internal/auth/http.go`:
```go
func clientAddr(r *http.Request) string {
    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return r.RemoteAddr
    }
    return host
}
```

In `internal/auth/service.go`:
```go
AddressLimit: 30, // 30 attempts per client address per 15-minute window
```

1. **Denial of Service (Collateral Lockout)**:
   All Tunisian students routed through the same Cloudflare Point of Presence (POP) share the same upstream Cloudflare IP. If an attacker triggers 30 failed login attempts, the entire Cloudflare POP IP is rate-limited, locking out legitimate users across the region.
2. **Rate Limit Bypass**:
   Distributed credential stuffing or brute-force attacks can easily rotate across different Cloudflare edge IPs, effectively bypassing the 30-attempt per-address limit.

## Proposed Remediation
Add trusted proxy configuration (e.g. `LEARNOO_TRUSTED_PROXIES` or `LEARNOO_BEHIND_CLOUDFLARE=true`):
1. Check that `r.RemoteAddr` originates from verified Cloudflare CIDRs.
2. If verified, extract the client IP securely from `CF-Connecting-IP`.
3. If unverified or running in local dev, fall back to `r.RemoteAddr`.
