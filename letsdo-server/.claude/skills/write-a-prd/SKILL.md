---
name: write-a-prd
description: Create a PRD through codebase exploration and submit as a GitHub issue. Use when user wants to write a PRD, plan a new feature, or says /write-a-prd.
---

# Write a PRD

Creates a Product Requirements Document for a feature, then submits it as a GitHub issue.

## Process

1. **Detect repo structure** - Check root directory:
   - Monorepo (has admin/, app/, server/, or similar top-level service dirs): Note which services are affected
   - Single repo: Look for internal/, src/, pkg/, or language-specific conventions

2. **Locate context** - Read SPECIFICATION.md if it exists. Check docs/schema/*.md for relevant tables. Read root CLAUDE.md and any sub-package CLAUDE.md files.

3. **Interview the user** - Ask for detailed problem description. Explore repo to verify assertions. Ask about alternatives considered. Be extremely detailed and thorough about scope.

4. **Hammer out scope** - What we build vs what we DON'T build.

5. **Identify affected areas** - Based on repo structure:
   - Monorepo: Which services? (e.g., "server + admin")
   - Single repo: Which domains/modules?

6. **Write the PRD** - Use the template below. Be exhaustive on user stories.

7. **Create GitHub issue** - Run gh issue create --title "PRD: [Feature Name]" --body "[content]". Do NOT ask for review first—just create it.

## PRD Template

## Problem Statement

[What problem does this feature solve? Write from user's perspective.]

## Solution

[High-level description. 2-3 sentences max.]

## User Stories

[EXHAUSTIVE numbered list. Cover happy paths, edge cases, error states, all actors.]

Format: "As a [actor], I want [capability], so that [benefit]"

1. As a ...
2. As a ...
(continue until all behaviors are captured — typically 15-30 stories)

## Implementation Decisions

- **Repo structure**: [monorepo / single-repo — as detected]
- **Affected services/packages**: [what you found]
- **Database tables**: [reference docs/schema/*.md if available]
- **API endpoints**: [list endpoints to add/modify]
- **State transitions**: [if applicable]
- **External services**: [if applicable]

Do NOT include file paths or code. Let implementer follow project conventions.

## Testing Approach

- Key behaviors to test
- Edge cases requiring coverage
- Integration points

## Out of Scope

[Explicitly list what this PRD does NOT cover]

## Open Questions

[Anything needing clarification before implementation]

## After Creation

Output:
1. The GitHub issue URL
2. One-line summary of what was created
3. Suggest: "Run /prd-to-issues #[number] to break this into task issues"
