#!/bin/sh -l

export GITHUB_USER=$1
export GITHUB_TOKEN=$2
export FILENAME=${3:-README.md}
export OUTPUT_FORMAT=${4:-list}
export IGNORE_REPOS=$5
export WITH_POC=${6-true}
export WITH_LICENSE=${7-true}
export WITH_STARS=${8-true}

stargazer -o "/github/workspace/$FILENAME"
