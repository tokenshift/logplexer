# Logplexer

Combines multiple input files into a single output stream, with each line prefixed by its source.

## Use

```
$ logplexer /var/log/nginx/access.log:ACCESS /var/log/nginx/error.log:ERROR
ACCESS: 10.0.8.42 - - [06/Nov/2014:19:10:38 +0600] "GET / HTTP/1.1" 404 177 "-"
ACCESS: 10.0.8.42 - - [06/Nov/2014:19:12:14 +0600] "GET /foo HTTP/1.1" 200 4356 "-"
ERROR:  2014/11/06 19:12:14 [error] file not found
ACCESS: 10.0.8.42 - - [06/Nov/2014:19:13:24 +0600] "GET / HTTP/1.1" 200 4223 "-"
```

The arguments are of the form: `filename[:tag]`. If no tag is specified, the
filename itself will be used as the tag.
