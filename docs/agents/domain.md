# Domain Docs

## Before exploring, read these

- `GLOSSARY.md` at the repo root, or `GLOSSARY-MAP.md` if present; read relevant context glossaries.
- `docs/adr/`; read ADRs that touch the area. For multi-context repos, also check `src/<context>/docs/adr/`.

If these files do not exist, proceed silently. Do not suggest creating them upfront; create them when domain terms or decisions are resolved.

## Layout

This repo uses single-context docs: root `GLOSSARY.md` and `docs/adr/`.

## Use glossary vocabulary

Use defined domain terms in issues, refactor proposals, hypotheses, and tests. If a needed term is missing, reconsider the term or note the gap for domain modeling.

## Flag ADR conflicts

Surface conflicts with existing ADRs rather than silently overriding them.
