# Product Requirements Document: libfizz-go

**Generated**: 2026-01-16
**Status**: Ready for AI Interview

## Initial Input

**Project Description**: a sane, simple, idomatic go library for the https://fizzy.do API

**I want to build libfizz-go that has the following features:**
1. follow all the warping rules about go code
2. implement all of the features of the https://fizzy.do HTTP API
3. favor simplicity over complexity
4. be idiomatic to go best practices
---

# Specification Generation

Agent workflow for creating project specifications via structured interview.

Legend (from RFC2119): !=MUST, ~=SHOULD, ≉=SHOULD NOT, ⊗=MUST NOT, ?=MAY.

## Input Template

```
I want to build libfizz-go that has the following features:
1. [feature]
2. [feature]
...
N. [feature]
```

## Interview Process

- ~ Use Claude AskInterviewQuestion when available (emulate it if not available)
- ! If Input Template fields are empty: ask overview, then features, then details
- ! Ask one focused, non-trivial question per step
- ~ Provide numbered answer options when appropriate
- ! Include "other" option for custom/unknown responses
- ! when you are done, append to the end of this file all questions asked and answers given.

**Question Areas:**
- ! Missing decisions (language, framework, deployment)
- ! Edge cases (errors, boundaries, failure modes)
- ! Implementation details (architecture, patterns, libraries)
- ! Requirements (performance, security, scalability)
- ! UX/constraints (users, timeline, compatibility)
- ! Tradeoffs (simplicity vs features, speed vs safety)

**Completion:**
- ! Continue until little ambiguity remains
- ! Ensure spec is comprehensive enough to implement

## Output Generation

- ! Generate as SPECIFICATION.md
- ! Break into phases, subphases, tasks
- ! Mark all dependencies explicitly: "Phase 2 (depends on: Phase 1)"
- ! Design for parallel work (multiple agents)
- ⊗ Write code (specification only)

**Structure:**
```markdown
# Project Name
## Overview
## Requirements
## Architecture
## Implementation Plan
### Phase 1: Foundation
#### Subphase 1.1: Setup
- Task 1.1.1: (description, dependencies, acceptance criteria)
#### Subphase 1.2: Core (depends on: 1.1)
### Phase 2: Features (depends on: Phase 1)
## Testing Strategy
## Deployment
```

## Best Practices

- ! Detailed enough to implement without guesswork
- ! Clear scope boundaries (in vs out)
- ! Include rationale for major decisions
- ~ Size tasks for 1-4 hours
- ! Minimize inter-task dependencies
- ! Define clear component interfaces

## Anti-Patterns

- ⊗ Multiple questions at once
- ⊗ Assumptions without clarifying
- ⊗ Vague requirements
- ⊗ Missing dependencies
- ⊗ Sequential tasks that could be parallel

---

## Interview Questions and Answers

**Q1: What is the Fizzy.do API?**
A: Fizzy is a modern kanban board application by Basecamp for tracking bugs, issues, ideas, and small projects. The API provides full REST access to all resources: Identity, Boards, Cards, Comments, Reactions, Steps, Tags, Columns, Users, Notifications, and File Uploads. API base URL: https://app.fizzy.do

**Q2: Authentication approach - Which authentication method should libfizz-go support?**
A: Option 1 - Personal Access Token only (simpler, for scripts/integrations)

**Q3: API coverage scope - Which resources should v1.0 include?**
A: Option 1 - All resources from day one (comprehensive but more work)

**Q4: Error handling strategy - How should libfizz-go handle errors?**
A: Option 2 - Typed errors + automatic retries (retry on 429 rate limits and 5xx with exponential backoff)
Rationale: Idiomatic Go with typed errors, automatic retries save callers from reimplementing logic, skip circuit breaker for simplicity

**Q5: HTTP client caching strategy - Should libfizz-go support ETag caching?**
A: Option 1 - Built-in ETag caching (automatic, enabled by default, configurable)

**Q6: Client configuration - Which configuration approach?**
A: Option 1 - Functional options pattern (most idiomatic Go, backwards compatible, self-documenting)
Example: `NewClient(token, accountSlug, WithTimeout(30*time.Second), WithCache(false))`

**Q7: API client structure - How should resources be organized?**
A: Option 1 - Namespaced services (most idiomatic for Go API clients, clear organization, scalable)
Example: `client.Boards.List()`, `client.Cards.Get(id)`, `client.Comments.Create(...)`

**Q8: Pagination handling - How should the library handle pagination?**
A: Option 3 - Helper + manual (return both results and NextPage; provide `ListAll()` helper that auto-fetches all pages)
Rationale: Flexibility for power users, convenience for scripters, explicit is better than implicit
