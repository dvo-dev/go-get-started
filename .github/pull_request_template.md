<!--Note these comments do not need to be deleted, and will not show up in the final PR-->

# Summary

## Related Issues
<!--Required-->

- Completes #69
- Partial fix #420

## Contributors
<!--Required-->
- @github: implementer
- @deeznutz: SQL consultation

## Affected Modules
<!--Required-->

- `pkg/database`
- `cmd/webapp`

## Description
<!--Required-->
This is where you should describe the contents of this pull request at a **high level**, aim for 2-3 paragraphs of content but more is always better than less.

Try to describe the changes in a "did X by Y", e.g. improved logging by adding multiple levels vs improved logging.

# Dependency Interactions
<!--Optional-->

If altering configurations/usage of external depedencies such as caches, databases, etc. describe it in this area.

## Redis
Adding usage for session caches.

## Elasticsearch
Change lookup parameters

# Testing (Unit/Integration/Performance)
<!--Optional if no tests altered/added-->
Describe tests modified or added for coverage of the PR's contents.

If something is not testable, explain here as well.

# Points of Attention
<!--Optional, use as often as possible-->
Call out anything you notice that may have an affect on others work, past or future.

This includes potential areas of refactorization, unavoidable code stinks, improved implementation methodology learned, etc.
