# Agent guidance

## Diff library versus terminal integration
- Separate line/character diff calculation from terminal colour, interactive paging, filesystem lookup and CLI process state. Test deterministic diff results with direct strings/readers when possible.
- Preserve existing public option types, call shapes and defaults. If the implementation uses type-switched variadic options, identify/document the accepted option types, repeated-option precedence and unsupported-type behaviour from code; do not change them without compatibility tests.
- Keep alignment, lookahead and end-of-line semantics observable in regression tests; verify empty input, insertions/deletions, repeated lines and Unicode, including errors.
- Do not make all unit tests run under a real pager or a fabricated terminal. Retain focused CLI/TTY/subprocess tests when terminal-specific behaviour is the contract.

## Review workflow
Check README examples and output guide when changing output. An agent must inspect implementation and tests and report gaps to issue #22 with exact file/line evidence, running relevant Go tests before requesting a separate implementation PR.
## Engineering reference and follow-up

- Adopt only relevant guidance from https://arran4.github.io/blog/post/2026/055-suggested-coding-standards-for-repositories-and-agents/ and https://arran4.github.io/blog/post/2026/020-optional-dependency-injection-via-type-switched-variadic-args/ .
- This documentation is not evidence of repository-wide compliance. A coding-agent pass must report checked files, observed gaps, tests/CI and existing issue overlap to https://github.com/arran4/golang-diff/issues/22 before any broad implementation.
- Do not merge implementation automatically; request review and follow repository-specific instructions.
