# logslice

Fast log file slicer that extracts time-range segments from large structured log files without loading them fully into memory.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

## Usage

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Pipe output to another tool:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log | grep "ERROR"
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the log file | required |
| `--from` | Start of time range (RFC3339) | required |
| `--to` | End of time range (RFC3339) | required |
| `--format` | Timestamp layout string | RFC3339 |
| `--field` | JSON field name for timestamp | `time` |

### Example

```bash
# Extract one hour of logs from a JSON-structured log file
logslice --file /var/log/app.log \
         --from "2024-01-15T14:00:00Z" \
         --to "2024-01-15T15:00:00Z" \
         --field timestamp
```

logslice uses binary search to locate the time boundaries, making it efficient even on multi-gigabyte log files.

## License

MIT © 2024 yourusername