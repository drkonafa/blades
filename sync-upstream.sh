#!/bin/bash

# Sync with upstream repository
# This script fetches updates from the original go-kratos/blades repo
# and merges them into your fork while preserving your enhancements

set -e

echo "🔄 Syncing with upstream repository..."

# Check if we're on main branch
current_branch=$(git branch --show-current)
if [ "$current_branch" != "main" ]; then
    echo "⚠️  Warning: You're not on the main branch. Current branch: $current_branch"
    echo "   Switching to main branch..."
    git checkout main
fi

# Fetch latest changes from upstream
echo "📥 Fetching latest changes from upstream..."
git fetch upstream

# Check if there are any new commits
behind=$(git rev-list --count HEAD..upstream/main)
if [ "$behind" -eq 0 ]; then
    echo "✅ Your fork is already up to date with upstream!"
    exit 0
fi

echo "📊 Found $behind new commits from upstream"

# Show what's coming
echo "📋 Recent commits from upstream:"
git log --oneline HEAD..upstream/main | head -5

# Ask for confirmation
read -p "🤔 Do you want to merge these changes? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Sync cancelled"
    exit 1
fi

# Merge upstream changes
echo "🔄 Merging upstream changes..."
if git merge upstream/main; then
    echo "✅ Successfully merged upstream changes!"
    
    # Push to your fork
    echo "📤 Pushing changes to your fork..."
    git push origin main
    
    echo "🎉 Sync complete! Your fork is now up to date with upstream."
    echo "   Your enhancements are preserved."
else
    echo "❌ Merge failed! There are conflicts that need to be resolved manually."
    echo "   Run 'git status' to see which files have conflicts."
    echo "   Edit the conflicted files, then run:"
    echo "   git add <resolved-files>"
    echo "   git commit"
    echo "   git push origin main"
    exit 1
fi
