#!/usr/bin/env python3

import csv
import sys
from shapely import geos, wkb

csv.field_size_limit(2147483647)
geos.WKBWriter.defaults['include_srid'] = True
tsv_in = csv.reader(sys.stdin, delimiter='\t')
tsv_out = csv.writer(sys.stdout, delimiter='\t')
for record in tsv_in:
    wkb_in = record[0]
    wkb_out = wkb.loads(bytes.fromhex(wkb_in))
    tags = record[1]
    tsv_out.writerow([wkb_out.wkb_hex, tags])
