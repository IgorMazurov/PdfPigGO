#!/usr/bin/env python3
import re
from collections import defaultdict

with open("port-status.md") as f:
    lines = f.readlines()

stats = defaultdict(int)

for line in lines:
    m = re.search(r'\b(NONE|COMPLETE|SKIPPED)\b', line)
    if m:
        status = m.group(1)
        if status == 'NONE':
            stats[(status, 'N/A')] += 1
        else:
            c = re.search(r'(\d+%)', line)
            stats[(status, c.group(1) if c else 'N/A')] += 1

for key, count in sorted(stats.items()):
    print(f"{key[0]:<12} {key[1]:<8} {count}")