#!/bin/sh -l

export GITHUB_USER=$1
export GITHUB_TOKEN=$2
export FILENAME=${3:-README.md}
export IGNORE_REPOS=$4

stargazer -o "/github/workspace/$FILENAME"
