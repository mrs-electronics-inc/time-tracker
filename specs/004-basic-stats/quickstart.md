# Quickstart: Basic Stats Command

## Prerequisites
- Time tracker installed and configured
- Some time entries tracked

## Usage Examples

### View daily totals (default)
```bash
time-tracker stats
```
Output:
```
Date       Total Time
2025-09-24 02:30
2025-09-25 01:15
...
```

### View weekly totals
```bash
time-tracker stats --weekly
```
Output:
```
Week Starting  Total Time
2025-09-23      15:45
2025-09-30      12:30
...
```

### View project totals
```bash
time-tracker stats --projects
```
Output:
```
Project    Total Time
Project A  08:20
Project B  05:15
...
```

### Combined flags
```bash
time-tracker stats --weekly --projects
```
Output:
```
Project       Week Starting  Total Time
Project A     2025-09-23      04:10
Project A     2025-09-30      03:25
...
```

## Validation Steps
1. Run `time-tracker stats` and verify daily totals for past 7 days
2. Run `time-tracker stats --weekly` and verify weekly totals for past 4 weeks
3. Run `time-tracker stats --projects` and verify project groupings
4. Test with no data: should show "No data available"