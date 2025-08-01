# Git Optimization Guide for MBankingCore

This document explains the Git optimizations implemented for the MBankingCore Go project.

## 🚀 Quick Start

```bash
# Clean workspace
./clean.sh

# Build application
./build.sh

# Run application
./bin/mbankingcore
```

## 📁 Git Configuration Files

### `.gitignore` - Optimized for Go Projects
- **Go specific**: Binaries, test files, coverage reports, workspace files
- **Build artifacts**: `bin/`, `dist/`, `build/`, debug binaries
- **Environment**: `.env` files, config with sensitive data
- **Security**: Certificates, private keys, database files
- **IDE/Editor**: VS Code, GoLand, Vim, Emacs configurations
- **OS files**: macOS `.DS_Store`, Windows `Thumbs.db`, Linux temp files
- **Dependencies**: `vendor/`, `node_modules/`, Python `__pycache__/`
- **Testing**: Coverage reports, benchmark results, test outputs
- **Documentation**: Large API docs, generated documentation
- **Cloud/Deploy**: Docker dev files, Kubernetes secrets, Terraform state

### `.gitattributes` - Consistent File Handling
- **Line endings**: LF for Go/script files, CRLF for Windows batch files
- **Binary detection**: Proper handling of executables, images, archives
- **Language detection**: Proper Go language attribution
- **Security**: Export-ignore for sensitive files
- **Git LFS ready**: Configuration for large file tracking

### `.gitmessage` - Commit Message Template
- **Conventional Commits**: Structured commit messages
- **Types**: feat, fix, docs, style, refactor, perf, test, build, ci, chore
- **Scopes**: auth, handlers, models, config, middleware, utils
- **Format**: `<type>(<scope>): <subject>`

## 🛠️ Scripts

### `clean.sh` - Workspace Cleanup
Removes all temporary files, build artifacts, and caches:
- Build artifacts (`mbankingcore`, `migrate`, `bin/`)
- Go build/test/module cache
- Temporary files (`*.tmp`, `*.log`, `*~`, `*.bak`)
- OS generated files (`.DS_Store`, `Thumbs.db`)
- Editor temporary files (`*.swp`, `*.swo`)
- Test artifacts (`*.test`, `*.out`, coverage files)
- Debug binaries (`__debug_bin*`, `*.debug`)

### `build.sh` - Clean Build Process
Builds the application in `bin/` directory to keep workspace clean:
- Cleans previous builds
- Clears Go build cache
- Builds binary in `./bin/mbankingcore`
- Keeps root workspace free of binaries

## 🔧 VS Code Integration

### Settings Optimization
- **Search exclusions**: Git objects, build directories, large files
- **File watching**: Excludes unnecessary directories for better performance
- **Go tools**: Optimized Go language server settings
- **Git integration**: Smart commit, autofetch, sync confirmation
- **File handling**: Trim whitespace, insert final newline, format on save

### Copilot Performance
- **Large file exclusion**: Postman collections, SQL files, binaries
- **Git objects exclusion**: Prevents indexing of Git history
- **Binary exclusion**: Executables, certificates, database files

## 📈 Performance Benefits

### Before Optimization
- Git repository: 146MB with large objects
- Binary in workspace: 20MB indexed by tools
- Large JSON files: 126KB Postman collection indexed
- Complex handlers: 800+ line files slow to process

### After Optimization
- Clean workspace: Only source code tracked
- Excluded binaries: Build artifacts in separate directory
- Optimized indexing: VS Code and Copilot skip unnecessary files
- Consistent line endings: Cross-platform compatibility
- Standard commit format: Better collaboration and history

## 🎯 Best Practices

### Development Workflow
1. **Start clean**: `./clean.sh` before major changes
2. **Build properly**: `./build.sh` for consistent builds
3. **Commit regularly**: Use template with `git commit`
4. **Check status**: `git status` to verify ignored files

### File Management
- **No binaries in root**: All executables in `bin/`
- **Environment files**: Use `.env.example` template
- **Certificates**: Auto-generated in `certs/` (ignored)
- **Logs**: Temporary in `.log` files (ignored)

### Performance Tips
- **VS Code**: Restart after optimization changes
- **Git**: Large objects already excluded from new commits
- **Copilot**: Faster due to smaller indexing scope
- **Build**: Incremental builds work better with clean workspace

## 🔍 Monitoring

### Check Ignored Files
```bash
# Find files that should be ignored
find . -name "*.log" -o -name "*.tmp" -o -name ".DS_Store"

# Check Git status
git status --ignored

# Verify binary location
ls -la bin/
```

### Performance Metrics
```bash
# Repository size
du -sh .git/

# Workspace size (excluding Git)
du -sh --exclude='.git' .

# File count by type
find . -name "*.go" | wc -l
```

## 📋 Checklist

- [x] ✅ `.gitignore` optimized for Go projects
- [x] ✅ `.gitattributes` for consistent file handling
- [x] ✅ `.gitmessage` template for structured commits
- [x] ✅ `.vscodeignore` for VS Code performance
- [x] ✅ VS Code settings optimized
- [x] ✅ Build script (`build.sh`) for clean builds
- [x] ✅ Clean script (`clean.sh`) for workspace maintenance
- [x] ✅ Binary files excluded from workspace
- [x] ✅ Git commit template configured
- [x] ✅ Copilot performance optimized

## 🚨 Important Notes

1. **Binary Location**: Always build to `bin/` directory
2. **Environment Files**: Never commit `.env` with secrets
3. **Certificates**: Auto-generated, never commit private keys
4. **Large Files**: Use Git LFS for files >100MB if needed
5. **Workspace**: Keep root clean of build artifacts

This optimization ensures:
- ⚡ Faster VS Code and Copilot performance
- 🔒 Better security (no accidental commits of secrets)
- 🤝 Consistent cross-platform development
- 📦 Smaller repository size
- 🎯 Better collaboration with structured commits
