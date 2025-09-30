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
+------------+------------+
|    DATE    | TOTAL TIME |
+------------+------------+
| 2025-09-24 | 02:30      |
+------------+------------+
| 2025-09-25 | 01:15      |
+------------+------------+
...
```

### View weekly totals
```bash
time-tracker stats --weekly
```
Output:
```
+---------------+------------+
| WEEK STARTING | TOTAL TIME |
+---------------+------------+
| 2025-09-23    | 15:45      |
+---------------+------------+
| 2025-09-30    | 12:30      |
+---------------+------------+
...
```

### View project totals
```bash
time-tracker stats --projects
```
Output:
```
+--------------+------------+
|   PROJECT    | TOTAL TIME |
+--------------+------------+
| Project A    | 08:20      |
+--------------+------------+
| Project B    | 05:15      |
+--------------+------------+
...
```



## Validation Steps
1. Run `time-tracker stats` and verify daily totals for past 7 days in table format
2. Run `time-tracker stats --weekly` and verify weekly totals for past 4 weeks in table format
3. Run `time-tracker stats --projects` and verify project groupings sorted by time descending in table format
4. Test `time-tracker stats --weekly --projects` shows error
5. Test with no data: should show "No data available"