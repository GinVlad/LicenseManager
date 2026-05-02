# Skill: Agentic Workflow

## Purpose
Coordinate tasks through Brainstorm -> Plan -> Cook -> Verify -> Fix stages automatically.

## Usage
```
/agentic [feature description]      # Full flow
/agentic brainstorm [feature]       # Single stage
/agentic plan [feature]
/agentic cook [feature]
/agentic verify [feature]
/agentic fix [feature]
/agentic resume                     # Continue from session.md
```

## When Invoked

### 1. Check Session State
```
Read: .claude/rules/session.md
- What phase are we in?
- Any in-progress features?
- Any blockers?
```

### 2. Determine Stage
- If new feature -> Start at BRAINSTORM
- If resuming -> Continue from last stage
- If specific stage requested -> Go there

### 3. Execute Stage

#### BRAINSTORM
```
Agents: CTO (lead)

1. Parse user request in LicenseManager context
2. Identify scope: backend? frontend? schema? SDK?
3. Load relevant rules ONLY
4. Generate 2-3 approaches
5. Recommend one

Output: .claude/plans/[feature]-brainstorm.md
Ask: "Does this approach work?"
```

#### PLAN
```
Agents: CTO + Backend-dev + Frontend-dev + Database

1. Read brainstorm
2. Break into tasks:
   - Task ID, description
   - Files to create/modify (exact paths)
   - Dependencies
   - S/M/L complexity
3. Order by dependencies

Output: .claude/plans/[feature]-plan.md
Ask: "Ready to implement?"
```

#### COOK
```
Agents: Backend-dev + Frontend-dev (parallel if possible)

1. Read plan
2. Implement task by task
3. Mark done in plan
4. Update session.md

Output: Working code
Gate: All tasks checked
```

#### VERIFY
```
Agents: Security (lead)

1. Check all new SQL uses params (no concat)
2. Check new routes have correct middleware
3. Check error handling covers all cases
4. Check frontend shows error states

Output: .claude/plans/[feature]-verify.md
Gate: No critical issues
```

#### FIX
```
1. Address each issue from verify
2. Re-check fixed items
3. Loop until clean
```

### 4. Update Session
After each stage update `.claude/rules/session.md`.
