# Learnoo design foundations

Status: initial implementation approved as a starting point, adjustable through review. This is a styling foundation, not completion of LRN-002's wireframes or a production UI kit.

## Direction

Use the [Udemy style reference on Refero](https://styles.refero.design/style/c03afcbd-96ed-4b7f-8d0a-277fc0042ba7) as inspiration for compact hierarchy, clear navigation, modest rounded corners, thin borders and a 4px spacing rhythm. Learnoo retains its own infinity identity, forest-green actions and restrained orange educational accents. Prefer flat light surfaces and readable content over decoration.

Refero's export is an interpreted style reference, not Udemy's official component library. These files are our adaptation, not a verbatim export. We do not adopt its brand font, violet palette, shopping interactions, decorative 3D or conflicting generated rules. Keep [MOTION.md](MOTION.md) authoritative.

## One source of truth

Edit [web/design/tokens.json](web/design/tokens.json), then run `npm run design:build`. Commit the source and generated outputs together. `npm run design:check` detects drift.

| File | Role |
| --- | --- |
| [tokens.json](web/design/tokens.json) | Canonical grouped values; simple project JSON, not a claimed DTCG interchange schema |
| [tokens.css](web/design/tokens.css) | Generated `--learnoo-*` CSS variables |
| [tailwind.css](web/design/tailwind.css) | Generated Tailwind v4 `@theme inline` mappings |
| [Reader tokens](internal/rendering/static/tokens.css) | Generated identical copy embedded by Go |
| [Reader styles](internal/rendering/static/site.css) | Hand-authored component/layout rules consuming tokens |

See [integration instructions](web/design/README.md). Never independently tune the generated CSS or Tailwind values. Layout-specific sizes may remain in component CSS; don't tokenize every measurement.

## Palette and typography

White content surfaces sit on a light neutral canvas. Ink is near-black; muted text must remain readable. Green is for links, primary actions and focus of the learning flow; orange identifies teaching notes and the unreviewed notice. Red indicates errors. Always pair status colours with explicit text; green alone never means pedagogical approval.

Use the tokenized system sans stack initially: no font downloads or Udemy font assets. Body text is 17px at the default root size with 1.7 line-height; section headings 24px, lesson titles 32px. A 48rem outer reader width keeps explanations focused; the general page maximum is 75rem. Use rem-based sizes to respect browser settings. KaTeX owns mathematical typography.

Scientific lessons use French. Teachers may add Arabic or Tunisian-dialect explanations with scientific terms in French. Mark language and direction on the relevant block; use logical spacing properties. Arabic uses its own fallback stack, normal letter spacing and generous line-height. Isolate formulas LTR inside RTL paragraphs. Interface localization remains independent of lesson language.

## Components and states

- Buttons: 8px corners, minimum 44px target height, green primary action, outlined secondary actions, explicit hover/focus/disabled states. Group actions by priority; do not invent commerce controls.
- Inputs: visible label, distinct border, legible content and errors, visible keyboard focus. Explain failure in text and preserve entered work.
- Content: clear heading order; neutral surfaces with limited borders. Use separators for teaching notes, pale-green formula blocks, expandable worked solutions and an honest video placeholder when none is selected.
- Navigation: readable links, wrapping at narrow widths, subject/chapter context. Planned hierarchy is section → subject → chapter → lesson; the prototype navigation remains a technical demo.
- Review status: explicitly distinguish unreviewed, in review, approved and withdrawn. Styling never grants access or changes editorial state.

Target at least 4.5:1 contrast for normal text, 3:1 for large text and meaningful control/focus boundaries. Verify actual foreground/background pairs. Preserve keyboard navigation, visible focus, useful reading without video, and horizontal scrolling within long formulas instead of page overflow. No decorative motion is introduced; future motion must follow reduced-motion requirements in MOTION.md.

## Brand and delivery

The original brand reference is [the supplied banner](assets/brand/learnoo-banner.png). Preserve it. Its multicolour infinity mark is the reference, not the superseded cartoon. The large original must not be loaded on every lesson; optimized responsive brand assets remain separate work. Use numerical examples such as `z = 3 + 2i` in promotional mockups to avoid the unintended Tunisian reading identified by the maintainer.

CSS variables serve the Go reader directly; Tailwind is build-time tooling for the planned Studio and must not add a browser runtime. No external font/CDN requests are needed for these foundations.

## Remaining design work

LRN-002 still needs mobile reader, Studio editor, reviewer decision, registration and dashboard wireframes. Use these foundations for those screens; do not infer that unbuilt flows are implemented. Preview and the student reader must continue sharing the canonical lesson renderer.
