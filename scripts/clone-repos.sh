#!/bin/bash

# fail out if GITHUB_TOKEN isn't set
set -u

curl -s "https://api.github.com/orgs/samsung-cnct/repos?access_token=${GITHUB_TOKEN}&per_page=100" \
  | jq -r '.[] | "git clone \(.ssh_url) \(.name)"' \
  | while read -r line; do
      echo "running $line ... "
      $line
    done
