#!/usr/bin/env python3

import sys
from shapely import geos, wkb

geos.WKBWriter.defaults['include_srid'] = True
for line in sys.stdin:
    wkb_in, tags = line.rstrip().split('\t')
    wkb_out = wkb.loads(bytes.fromhex(wkb_in))
    print(f"{wkb_out.wkb_hex}\t{tags}")
