# The Tagging Strategy Explained

From the [CI configuration](../.github/workflows/ci.yml):

```
tags: |
  type=ref,event=branch              # Branch name (e.g., 'main', 'develop')
  type=ref,event=pr                  # PR number (e.g., 'pr-123')
  type=semver,pattern={{version}}    # Git tags like v1.0.0
  type=semver,pattern={{major}}.{{minor}}  # 1.0
  type=semver,pattern={{major}}      # 1
  type=sha,prefix={{branch}}-        # Branch + short SHA
  type=raw,value=latest,enable={{is_default_branch}}  # 'latest' only on main!
  ```
So:

- Pushing to main → gets latest, main, main-abc1234 tags
- Pushing to develop → gets develop, develop-xyz5678 tags
- Creating a tag like v1.0.0 → gets 1.0.0, 1.0, 1, latest tags
- Other branches → get branch name and sha tags only