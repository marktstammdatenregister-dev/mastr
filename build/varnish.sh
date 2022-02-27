#!/usr/bin/env sh
set -eu

# https://varnish-cache.org/docs/7.0/users-guide/storage-backends.html#file
# https://varnish-cache.org/docs/7.0/reference/varnishstat.html
/usr/sbin/varnishd -F -f /work/default.vcl -a http=:8080,HTTP -T none -s file,varnish.cache,2G &
/usr/bin/varnishncsa -w /dev/stdout -F '%{Host}i %h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\" \"%{Varnish:hitmiss}x\" \"%{x-location-latitude}i\" \"%{x-location-longitude}i\"' &
sh -c "while true; do sleep 30; varnishstat -1 -I 'MAIN.client_req' -I 'MAIN.cache_*' -I 'MGT.child_panic' -I 'MAIN.backend_fail' -I 'MAIN.n_lru_nuked' -I 'MAIN.s_fetch'; done"
