---
name: Feature request
about: Suggest an idea for this project
title: "[FEATURE] - Feature Title Here"
labels: ''
assignees: ''

---

<!--These comments do not to be deleted and will not be displayed in the final output-->

## Description
<!--Required-->
This is a high level explanation of the feature request. Try to keep it around 2-3 paragraphs in length, but more is better than less. Attempt to keep this section as specific as possible without going into in-depth implementation details.

## Acceptance Criteria
<!--Required-->
This is a list of sub features that must be implemented for the overall request to be considered complete. This is where in-depth implementation considerations, suggestions, etc. should also go.

- Must support GET + POST requests
  - Return FileNotFound status code if file does not exist
- Must have full unit coverage
 - The database dependency should be mocked
- Must create interface for future dependency alternatives
 - Should implement the methods `Upload`, `Delete`, `Retrieve`

## Things to Consider
<!--Optional-->
While this section is not required, it should be used as often as possible. Potential conflicts with the existing codebase, required refactorization, concerns for maintaining should be mentioned here.
