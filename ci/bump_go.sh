#!/bin/bash

REF_BRANCH=master
PR_KEYWORD="[bump-go]"
PR_LABEL="bump-go"

set -e

if [ ! -z "$(git status -s -uall)" ]; then
    echo "Working directory is not clean" 2>&1
    exit 1
fi

git show-ref --verify --quiet refs/heads/bump-go || git branch -c $REF_BRANCH bump-go

# The version in the GO_VERSION file
current_version=$(git show $REF_BRANCH:GO_VERSION 2>/dev/null || echo "")
echo "Current version: $current_version"

# The latest version installed
latest_version=$(go version | sed -E 's/.*go([0-9\.]*).*/\1/')
echo "Installed version: $latest_version"

# The version already open in a PR
bump_candidate_version=$(git show bump-go:GO_VERSION 2>/dev/null || echo "")
echo "Bump candidate version: $bump_candidate_version"

existing_prs=$(gh pr list --state open --base $REF_BRANCH --label "$PR_LABEL" --json "number" | jq -r 'map(.number) | .[]')

if [ "$current_version" == "$latest_version" ]; then
    echo "Go is up to date"
    exit 0
elif [ "$latest_version" == "$bump_candidate_version" ] && [ ! -z "$existing_prs" ]; then
    echo "A PR is already open"
    exit 0
fi

if [ -z "$GITHUB_TOKEN" ]; then
    echo "GITHUB_TOKEN not set, but this is what I would do:"
    echo "git checkout bump-go && git merge --ff-only $REF_BRANCH"
    echo "echo \"$latest_version\" > GO_VERSION"
    echo "git add GO_VERSION"
    echo "git commit -m \"$PR_KEYWORD Bump Go to $latest_version\""
    echo "git push origin bump-go"
    if [ -z "$existing_prs" ]; then
        echo "gh pr create --base $REF_BRANCH --head bump-go --title \"$PR_KEYWORD Bump Go to $latest_version\" --label \"$PR_LABEL\""
    else
        first_id=$(echo "$existing_prs" | head -n 1)
        echo "gh pr edit $first_id --title \"$PR_KEYWORD Bump Go to $latest_version\""
    fi
else
    git checkout bump-go && git merge --ff-only $REF_BRANCH
    echo "$latest_version" > GO_VERSION
    git add GO_VERSION
    if git diff --quiet --cached; then
        echo "No changes to commit"
    else
        git commit -m "$PR_KEYWORD Bump Go to $latest_version"
    fi
    git push origin bump-go
    if [ -z "$existing_prs" ]; then
        gh pr create --base $REF_BRANCH --head bump-go --title "$PR_KEYWORD Bump Go to $latest_version" --label "$PR_LABEL" \
            --body "This PR was automatically created by the bump-go action."
    else
        first_id=$(echo "$existing_prs" | head -n 1)
        gh pr edit $first_id --title "$PR_KEYWORD Bump Go to $latest_version"
    fi
fi
   