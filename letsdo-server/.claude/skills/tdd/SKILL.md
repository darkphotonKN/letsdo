---
name: tdd
description: Test-driven development with red-green-refactor loop. Use when user says /tdd or asks to implement with TDD.
---

# Test-Driven Development

Implement features using the RED → GREEN → REFACTOR cycle.

## Philosophy

- Tests verify **behavior** through public interfaces
- Tests describe **WHAT**, not HOW
- One behavior at a time — not all tests first

## Anti-Pattern: Horizontal Slices

WRONG (horizontal):
   Write test1, test2, test3, test4
   Then implement all

RIGHT (vertical):
   test1 → impl1 → test2 → impl2 → test3 → impl3

## Process

### 1. Identify the task

| User says | Action |
|-----------|--------|
| /tdd #44 | Fetch gh issue view 44 |
| /tdd (no number) | Check conversation for current task context. If none, run gh issue list --limit 10 and ask which to work on |

### 2. Plan first (REQUIRED)

Before writing any code, output a plan:

**Task:** [one-line summary]

**Approach:**
1. [First thing to do]
2. [Second thing]
3. [Third thing]

**Tricky parts:**
- [potential gotcha]
- [edge case to handle]

**Files to touch:**
- [file1]
- [file2]

**First test:** [what the first RED test will be]

Ask: "Does this plan look right?" Wait for approval before coding.

### 3. Preparation

- [ ] Read the task issue description
- [ ] Check docs/schema/*.md for relevant tables
- [ ] Check CLAUDE.md for project patterns
- [ ] Identify test file conventions (look at existing tests)
- [ ] List behaviors to implement (from acceptance criteria)

### 4. TDD Loop

For EACH behavior:

**RED**
Write ONE failing test
Run tests → FAIL

**GREEN**
Write MINIMAL code to pass
Run tests → PASS

Repeat for next behavior.

### 5. Refactor

After all behaviors pass:
- Extract duplication
- Simplify interfaces
- Improve naming
- Run tests after each change → must stay GREEN

## Go Patterns

### Test Naming
Test[Function]_[Scenario]_[Expected]

### Table-Driven Tests
```go
func TestCreateStore_Validation(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateStoreInput
        wantErr bool
    }{
        {"valid", CreateStoreInput{Name: "Shop"}, false},
        {"empty name", CreateStoreInput{Name: ""}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := NewService(mockRepo)
            _, err := svc.Create(ctx, tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Commands
```bash
make test
go test ./... -v
go test ./path/to/package -run TestName -v
```

## Commit Pattern

After completing the issue:
```bash
git commit -m "feat(scope): description (#issue-number)"
```
