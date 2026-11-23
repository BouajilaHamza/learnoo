# Contributing to Learnoo

The project is preparing its first prototype. Start with [README](README.md), [requirements](docs/REQUIREMENTS.md), and [the backlog](docs/ROADMAP.md).

Planning feedback, issue discussions, and expressions of interest are welcome. Acceptance of external code and educational contributions is deferred until complete [licensing and inbound contribution terms](docs/LICENSING.md) are established. Do not submit third-party teaching material without permission or publish student information.

## Development workflow

Use focused branches and pull requests. Describe the problem, resulting behavior, relevant acceptance criteria, and validation. The maintainer reviews merges. CODEOWNERS requests review; repository protection settings must separately enforce it.

Tests should cover meaningful risks: authorization, submitted-content immutability, publication races, lesson withdrawal, rendering safety, and recovery. Avoid broad scaffolding or unrelated dependencies. Document runnable setup when the prototype exists.

## Educational workflow

Teachers will use Studio, not GitHub. Content must have a curriculum objective, clear explanation, sources/rights information, and a different qualified reviewer. AI-generated material is subject to the same review requirements.

## Collaboration

Be respectful, specific, and constructive. Discuss the work rather than attacking contributors. Do not post secrets, personal student data, or exploit details publicly.

## Branches and pull requests

Fork this repository and create a feature branch from `develop`. Open a pull request targeting `develop`; explain the change and validation. Do not request direct pushes to `main` or `develop`. The sole code owner is @hadden-me; their approval is required for contributor changes. Release pull requests go from `develop` to `main`.

GitHub does not allow authors to approve their own pull requests. Maintainer-authored changes need a separately documented workflow; do not silently add an approval bypass. Branch rules must be configured on GitHub in addition to CODEOWNERS.
