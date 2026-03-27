---
name: prd-to-issues
description: Break a PRD into independently-grabbable GitHub issues using vertical slices. Use when user says /prd-to-issues or wants to convert a PRD to implementation tasks.
---

# PRD to Issues

Breaks a PRD GitHub issue into task issues. Each task is a vertical slice (tracer bullet).

## What's a Vertical Slice?

A thin cut through ALL layers delivering ONE complete behavior:

VERTICAL (correct):
   migration → model → repo → service → handler → tests
   = ONE working endpoint, demoable

HORIZONTAL (wrong):
   "all migrations" then "all models" then "all handlers"
   = nothing works until everything's done

## Process

### 1. Fetch the PRD

Determine which PRD to process:

| User says | Action |
|-----------|--------|
| /prd-to-issues #42 | Fetch gh issue view 42 |
| /prd-to-issues (no number) | Check conversation history for recently created PRD. If found, confirm: "Use PRD #42 (Title)?" |
| /prd-to-issues (no context) | Run gh issue list --search "PRD:" --limit 5 and ask user to pick |

Always confirm the PRD title before proceeding.

### 2. Read project context

Check SPECIFICATION.md, CLAUDE.md, and relevant docs/schema/*.md files.

### 3. Detect repo structure

- Monorepo: Note which services each slice affects
- Single repo: Note which modules/domains

### 4. Draft vertical slices

Rules:
- Each slice = ONE complete behavior through all layers
- First slice should be the foundation (migrations, models)
- Each subsequent slice is independently completable after foundation
- Prefer thin slices — if in doubt, split further

### 5. Present breakdown to user

Format:

Proposed breakdown for PRD #[number]:

1. **[Title]**
   - Scope: [brief description]
   - Type: AFK / HITL
   - Blocked by: none / #issue
   - User stories: [numbers from PRD]

2. **[Title]**
   ...

Definitions:
- **AFK** (Away From Keyboard): Can be done autonomously
- **HITL** (Human In The Loop): Needs human decisions (payments, security, complex logic)

Ask: "Is this granularity right? Should I split or combine any?"

Wait for approval before creating issues.

### 6. Create GitHub issues

After approval, create issues in dependency order using gh issue create.

Issue template:

## Parent PRD

#[prd-issue-number]

## What to Build

[End-to-end behavior description. NOT layer-by-layer steps.]

## Acceptance Criteria

- [ ] [Specific testable criterion]
- [ ] [Another criterion]
- [ ] Tests pass for happy path
- [ ] Tests pass for error cases
- [ ] Follows project conventions (CLAUDE.md)

## Blocked By

#[issue-number] or "None — can start immediately"

## User Stories Addressed

From PRD #[number]: [list story numbers]

## TDD Approach

- RED: [first failing test to write]
- GREEN: [what makes it pass]

### 7. Output summary

After creating all issues:

Created [N] task issues:

#[num] [title] — blocked by: [deps]
#[num] [title] — blocked by: [deps]
...

Dependency graph:
#43 ─┬→ #44
     ├→ #45
     └→ #46 → #47

Recommended order: #43 → #44 → #45 → #46 → #47

Start with: #43 (no blockers)

Run /tdd #[number] when working on each issue.
