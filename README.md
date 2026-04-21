# envport

> CLI tool to snapshot and restore environment variable sets across projects

## Installation

```bash
go install github.com/yourusername/envport@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envport/releases).

## Usage

**Save a snapshot of your current environment:**
```bash
envport save myproject
```

**Restore a saved snapshot:**
```bash
envport load myproject
```

**List all saved snapshots:**
```bash
envport list
```

**Remove a snapshot:**
```bash
envport delete myproject
```

**Show the contents of a snapshot:**
```bash
envport show myproject
```

Snapshots are stored in `~/.envport/` as simple JSON files, making them easy to inspect or version control.

### Example Workflow

```bash
# Working on project A — save its environment
envport save project-a

# Switch to project B
envport load project-b

# Come back to project A later
envport load project-a
```

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

## License

[MIT](LICENSE)
