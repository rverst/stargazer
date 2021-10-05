#!/bin/sh -l

export GITHUB_USER=$1
export GITHUB_TOKEN=$2
export FILENAME=${3:-README.md}
export OUTPUT_FORMAT=$4
export IGNORE_REPOS=$5
export WITH_TOC=$6
export WITH_LICENSE=$7
export WITH_STARS=$8
export WITH_BACK_TOP=$9

stargazer -o "/github/workspace/$FILENAME"
