# stargazer

*stargazer* creates a sorted list of your stared github repositories.
Like an [![Awesome](https://awesome.re/badge.svg)](https://awesome.re)
list, but personal.

## Usage

Probably the easiest way to get your own stargazer list is to have a
repository with a workflow that uses the GitHub action.

```yaml
# This workflow builds a list of your starred repositories
name: Stargazer

on:
  schedule:
    - cron: '42 2 * * *'

  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      # Checks-out your repository under $GITHUB_WORKSPACE, so your job can access it
      - uses: actions/checkout@v2

      # Generate the list
      - name: Create star list
        id: stargazer
        uses: rverst/stargazer@5e231084ea229d6649fcaeaa04c4f3e336b57139
        with:
          github-user: ${{ github.actor }}
          github-token: ${{ secrets.GITHUB_TOKEN }}
          list-file: "README.md"
          # you can ignore repositories, that will not get on the list
          # ignored-repositories: ${{ secrets.IGNORE_REPO }}

      # Commit the changes
      - name: Commit files
        run: |
          git config --local user.email "actions@noreply.github.com"
          git config --local user.name "github-actions[bot]"
          git add .
          git commit -m "Update list"

      # Push the changes
      - name: Push
        uses: ad-m/github-push-action@master
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          branch: ${{ github.ref }}
```

## Inputs

| Name | Type | Required | Description |
|------|------|----------|---------|-------------|
| github-user | string | true GitHub user whose stars are fetched |
| github-token | string | true Access token for the GitHub API |
| list-file | string | false Filename of the stargazer list (default: README.md) |
| ignored-repositories | string | false | Comma separated list of repositories (user/repo) to ignore |

## Inspiration

*stargazer* is inspired by [starred](https://github.com/gmolveau/starred),
which is very similar and written in python. I created *stargazer* because I
wanted to try out the [GitHub GraphQL API v4](https://docs.github.com/en/graphql).
