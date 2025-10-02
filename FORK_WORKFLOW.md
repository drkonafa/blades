# Fork Maintenance Workflow

This document explains how to maintain your fork of the go-kratos/blades repository while keeping your enhancements.

## Current Setup

- `upstream` = Original go-kratos/blades repository
- `origin` = Your fork (to be created)

## Initial Setup (One-time)

1. **Create your fork on GitHub:**
   - Go to https://github.com/go-kratos/blades
   - Click "Fork" button
   - This creates `https://github.com/YOUR_USERNAME/blades`

2. **Add your fork as origin:**
   ```bash
   git remote add origin https://github.com/YOUR_USERNAME/blades.git
   git push -u origin main
   ```

## Regular Workflow

### When you want to pull updates from the original repo:

```bash
# 1. Fetch latest changes from upstream
git fetch upstream

# 2. Switch to main branch
git checkout main

# 3. Merge upstream changes into your main
git merge upstream/main

# 4. Push updates to your fork
git push origin main
```

### When you want to make new changes:

```bash
# 1. Create a new branch for your feature
git checkout -b feature/your-feature-name

# 2. Make your changes
# ... edit files ...

# 3. Commit your changes
git add .
git commit -m "feat: Your feature description"

# 4. Push to your fork
git push origin feature/your-feature-name

# 5. Create a Pull Request on GitHub to upstream (optional)
```

### When you want to sync with upstream and resolve conflicts:

```bash
# 1. Fetch upstream changes
git fetch upstream

# 2. Switch to main
git checkout main

# 3. Try to merge (may have conflicts)
git merge upstream/main

# 4. If there are conflicts, resolve them:
# - Edit the conflicted files
# - Remove conflict markers (<<<<<<< ======= >>>>>>>)
# - git add the resolved files
# - git commit

# 5. Push the resolved changes
git push origin main
```

## Your Enhancements

Your fork includes these enhancements:

- **Enhanced Chain Visualization**: Automatic progress bars, colors, timing
- **Gemini Provider**: Full streaming support for Google's Gemini API
- **Chain Executor**: Verbose execution with beautiful output
- **Environment Configuration**: Automatic .env file loading
- **Agent Getters**: Public methods for name and instructions

## Best Practices

1. **Keep your main branch clean** - only merge from upstream or your feature branches
2. **Use feature branches** for new enhancements
3. **Regularly sync with upstream** to avoid large conflicts
4. **Test your changes** before merging to main
5. **Document your enhancements** in this file

## Commands Reference

```bash
# Check current remotes
git remote -v

# Check status
git status

# See commit history
git log --oneline

# See differences with upstream
git diff upstream/main

# Create a backup branch
git checkout -b backup-$(date +%Y%m%d)
```
