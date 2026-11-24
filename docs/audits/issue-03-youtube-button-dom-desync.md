# Issue: Video Element DOM Leak on YouTube Embed Replacement

## Status
- **Severity**: Low (UX / Visual Defect)
- **Affects**: `internal/rendering/static/math.js`, `internal/rendering/page.html`
- **Found**: 2026-11-24 Adversarial Audit

---

## Summary
When clicking the "Charger la vidéo YouTube" button on a lesson page, only the `<button>` is replaced with the `<iframe>`. The sibling `<p>La vidéo ne se charge qu’après votre choix.</p>` remains permanently visible below the active playing video.

## Mechanics
In `page.html`:
```html
<section class="video">
  <h3>{{.Title}}</h3>
  {{if .VideoID}}
    <button type="button" class="video-load" data-video-id="{{.VideoID}}">Charger la vidéo YouTube</button>
    <p>La vidéo ne se charge qu’après votre choix.</p>
  {{end}}
</section>
```

In `math.js`:
```javascript
button.addEventListener('click', () => {
    ...
    button.replaceWith(iframe);
}, { once: true });
```
The `<p>` disclaimer is never removed or hidden upon iframe insertion.

## Proposed Remediation
Wrap both button and text in a container or remove the sibling notice inside `math.js`:
```javascript
const notice = button.nextElementSibling;
if (notice && notice.tagName === 'P') notice.remove();
button.replaceWith(iframe);
```
