## Change

Describe the problem and resulting behavior.

## Upgrade impact

- [ ] No action is required from lesson maintainers (explain why), or an upgrade note
  is included in `content/docs/upgrades/`.
- [ ] For changes requiring action, the note explains affected lessons, the
  consequence of doing nothing, the edits to make, and how to verify them.
- [ ] Any removal has a documented transition period, or the reason for immediate
  removal is explained. Breaking changes are identified in the commit message.

Create a note with `hugo new content --kind upgrade-note docs/upgrades/short-name.md`,
fill in the four sections, and run `python3 scripts/upgrade-report.py preview`.
Leave its version `unreleased` until release-PR preparation. Explicit breaking-change
markers require a new breaking note in CI.

## Validation

Describe the relevant checks and their results. For upgrade changes, include an
example generated report and the affected/migrated cases tested. The **Check upgrade
guidance** workflow also provides a preview in its summary and downloadable artifact.
