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
        uses: rverst/stargazer@v1.2.6
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