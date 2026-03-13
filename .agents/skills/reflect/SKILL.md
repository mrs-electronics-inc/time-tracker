---
name: reflect
description: "Updates agent long-term memory in AGENTS.md and relevant skills with learnings and mistakes from conversations. Use when asked to reflect on a conversation or update memory."
---

# Reflect

Updates agent long-term memory with learnings and mistakes.

## Guidelines

### Always Do (no asking)

- Extract mistakes and learnings from the conversation
- Focus on high-level mistakes that warrant long-term documentation (patterns, documentation gaps, repeated issues), not one-off errors
- Present a refined list to the user for review
- Use markdown subsection headers to divide sections
- Prefer nested bullet lists over paragraphs
- Keep bullets to one or two sentences max
- Keep edits concise; avoid wasting context and tokens
- Integrate learnings into existing structure
- Use examples only when they add significant value

### Ask First (pause for approval)

- Update AGENTS.md
- Update skill files
- Create commits, open PRs, or push changes

### Never Do (hard stop)

- Modify files outside AGENTS.md and skill files
- Commit/push or publish changes without explicit user permission
- Update other projects' memory
- Add generic sections like "Mistakes to Avoid"
- Use numbers in section headers
- Use bold text to divide information in sections
- Use more than two sentences in a bullet point
- Use long paragraphs instead of bullet lists
- Add unhelpful examples

## Workflow

### Develop List of Mistakes and Learnings

Review the conversation and extract:

- **Mistakes**: What went wrong, what approach failed, what took too long
- **Learnings**: Patterns discovered, best practices identified, tool usage insights, project structure knowledge, user preferences

Format as bullet points. Each entry should capture one specific insight.

### Get User Feedback

Present the list to the user and ask them to:

- Add items you missed
- Edit items for clarity
- Remove items that aren't valuable
- Prioritize which learnings matter most

### Update Memory Files

#### AGENTS.md (project-level)

- Integrate learnings into existing sections
- Create new sections only when learnings don't fit existing structure
- Keep entries concise: one or two sentence per bullet point
- Link to relevant files when helpful

#### Skill Files

- Update skill descriptions or instructions if learnings reveal gaps
- Follow the same concise, bullet-point style
- Only update skills directly related to learnings
