---
name: reflection
description: This skill should be used when the user says "reflection", "reflect on the changes", "reflect on this session", or asks to capture lessons learned from the current conversation.
disable-model-invocation: true
---

# Reflection

The goal is to identify lessons from this session that should be permanently captured in CLAUDE.md — so future sessions benefit without repeating the same mistakes, clarifications, or decisions.

Review the **entire conversation history** — every message, correction, and preference — as the primary source. Code diffs are supplementary.

## Current Git State (supplementary context)

- Branch: !`git branch --show-current`
- Changes since last commit: !`git diff --stat HEAD`
- Staged changes: !`git diff --stat --cached`

## Process

1. **Read the full conversation from top to bottom.** Pay attention to:
   - Questions the user had to answer that should have been obvious from CLAUDE.md
   - Corrections the user made to your approach or output
   - Preferences or constraints the user stated (even casually)
   - Things you got wrong on the first attempt and had to revise
   - Decisions made about architecture, naming, tooling, or workflow
   - Anything the user explicitly said to always/never do

2. **Review the diffs** (supplementary) — Read modified files for context on what was built and why, but do not let this overshadow lessons from the conversation itself.

3. **Identify lessons** in these categories:
   - **Patterns & conventions** — Things that worked well and should be encoded as rules
   - **Gotchas & pitfalls** — Things that caused confusion, required retries, or were non-obvious
   - **Architecture decisions** — Choices made that future sessions should know about
   - **Workflow & communication preferences** — How the user prefers to work, communicate, or receive output
   - **Outdated/wrong memory** — Anything in CLAUDE.md or MEMORY.md that turned out to be incorrect or missing

4. **Read the current CLAUDE.md** to avoid duplicating what's already there and to find gaps.

5. **Propose edits to CLAUDE.md** — For each lesson worth keeping, suggest the specific text to add, change, or remove, and where it belongs.

6. **If no lessons are found**, explicitly state: "No CLAUDE.md updates needed from this session." This confirms the session was considered.

7. **Do not apply edits automatically.** Present proposals to the user and wait for approval.

8. **After applying approved edits**, print a brief summary in chat of what changed.

## Output Format

**When lessons are found:**
```
## Reflection

### [Category]
**Lesson**: <what was learned>
**Proposed CLAUDE.md change**: <exact text, with target section>

---
(one block per lesson)
```

**When no lessons are found:**
```
## Reflection

No CLAUDE.md updates needed from this session. The following were considered but already covered or not worth persisting:
- <item> — already in CLAUDE.md / too session-specific / etc.
```

**After applying approved changes:**
```
## CLAUDE.md updated

- Added: "<description>" under ## Section
- Modified: "<what changed>" in ## Section
- Removed: "<what was removed>"
```
