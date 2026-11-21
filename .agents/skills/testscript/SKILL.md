---
name: testscript
description: >-
  Choose, write, or review Layerctl end-to-end CLI tests in
  cmd/layerctl/testdata/script, including workflows that persist Git state
  across commands and user-facing failure or recovery behavior.
---

# Layerctl test scripts

Use a script when the promise depends on the built CLI, such as its
arguments, exit status, output, or state observed by a later invocation.
Use Go package tests when a package API is the relevant boundary or when
the case needs precise control of internal state.

Give each script one user-visible promise.
Keep the manifests, patch inputs, and expected results alongside the commands
that exercise them.
Use real Git and local repositories for sources.
Keep a workflow together when the result depends on earlier invocations.

Assert the full value or state that the promise covers.
Compare a complete output stream when its format is part of the contract.
For Git state, query the relevant tree, ref, or worktree with Git and compare
the result with an independent expectation.
For a rejected command, check the failure and any state the command promises
to preserve.
When stderr includes Git diagnostics or progress, anchor the Layerctl message
you care about and allow variation in unrelated text.
Match the expected values of paths, refs, and other fields that matter to the
scenario.
Use a regex for fields whose literal value can vary within the contract.

Read [Script format](references/format.md) when authoring or running a script
for fixture syntax, assertions, and focused test commands.
