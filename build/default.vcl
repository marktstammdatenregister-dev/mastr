# Mostly copied from:
# https://github.com/fly-apps/fly-varnish#verifying-that-caching-works

vcl 4.0;

backend default {
  .host = "localhost";
  .port = "9090";
  .probe = {
      .url = "/Marktstammdatenregister?sql=select+%22up%22";
      .interval = 10s;
  }
}

# TODO: Consider removing all the "x-cache" stuff from here ...
sub vcl_recv {
    unset req.http.x-cache;

    # Cookies stop Varnish from caching the page. We don't use cookies, so
    # cookies are only ever set by accident.
    # https://varnish-cache.org/docs/7.0/users-guide/increasing-your-hitrate.html#cookies
    unset req.http.Cookie;

    # Respect "no-cache" pragma header.
    if (req.http.Pragma ~ "no-cache") {
        return(pass);
    }
}

sub vcl_hit {
    set req.http.x-cache = "hit";
}

sub vcl_miss {
    set req.http.x-cache = "miss";
}

sub vcl_pass {
    set req.http.x-cache = "pass";
}

sub vcl_pipe {
    set req.http.x-cache = "pipe uncacheable";
}

sub vcl_synth {
    set req.http.x-cache = "synth synth";
    set resp.http.x-cache = req.http.x-cache;
}

sub vcl_deliver {
    if (obj.uncacheable) {
        set req.http.x-cache = req.http.x-cache + " uncacheable" ;
    } else {
        set req.http.x-cache = req.http.x-cache + " cached" ;
    }
    set resp.http.x-cache = req.http.x-cache;
}
# ... to here. This stuff is 100% for debugging purposes.

# We only serve static content. Cache everything! Note that the service is
# currently restarted every weekday morning, so the maximum uptime is Friday
# morning to Monday morning (three days).
#
# Since TTL is respected by browsers too, we don't want to overdo it with the
# max-age.
#
# We set this here because Datasette does not set "Cache-control: max-age" on
# static files even when the default_cache_ttl setting is set explicitly:
# https://github.com/simonw/datasette/issues/1645
#
# TODO: Consider caching status 400, which Datasette uses for timeouts. Those
# are the most expensive queries.
#
# TODO: Consider do_gzip:
# https://varnish-cache.org/docs/7.0/users-guide/compression.html#compressing-content-if-backends-don-t
sub vcl_backend_response {
    if (beresp.status == 200 || beresp.status == 404) {
        set beresp.ttl = 1h;
    }
}
