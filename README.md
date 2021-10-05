# stargazer
[![Docker](https://github.com/rverst/stargazer/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/rverst/stargazer/actions/workflows/docker-publish.yml)

*stargazer* creates a sorted list of your stared GitHub repositories.
Like an [![Awesome](https://awesome.re/badge.svg)](https://awesome.re)
list, but personal. Automated with GitHub-Actions.

See [rverst/stars](https://github.com/rverst/stars) for an example. You can use
that repository as a template, the README.md will get overwritten with your own
list if you run the stargazer-action (runs daily at 02:42).

## Usage

Probably the easiest way to get your own stargazer list is to have a
repository with a workflow that uses the GitHub action.
All you need to do is create a new repository and create the following workflow.

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
        uses: rverst/stargazer@v1
        with:
          github-user: ${{ github.actor }}
          github-token: ${{ secrets.GITHUB_TOKEN }}
          list-file: "README.md"

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
|------|------|----------|-------------|
| github-user | string | true | GitHub user whose stars are fetched |
| github-token | string | true | Access token for the GitHub API |
| list-file | string | false | Filename of the stargazer list (default: README.md) |
| format | string | false | Format of the stargazer list [list, table, \<custom\>] (default: list) |
| ignored-repositories | string | false | Comma separated list of repositories (user/repo) to ignore |
| with-toc | bool | false | Print table of contents (default: true) |
| with-license | bool | false | Print license of repositories (default: true) |
| with-stars | bool | false | Print starcount of repositories (default: true) |
| with-back-to-top | bool | false | Generate 'back to top' links for each language (default: false) |

## Custom templates

You can put your own templates in the repository and give its name as `format`. Have a look at
the included templates to get an understanding of the template model. Use `{{ printf "%#v" . }}`
to print the underlying struct.  
If you use a custom template, please be so kind and credit this repository, thanks a lot!

## Inspiration

*stargazer* is inspired by [starred](https://github.com/gmolveau/starred),
which is very similar and written in python. I created *stargazer* because I
wanted to try out the [GitHub GraphQL API v4](https://docs.github.com/en/graphql).
