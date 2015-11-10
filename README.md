# ratelimit

ratelimit implements [Token-Bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm to help you limit transfer rate. It is rewrote from and inspired by [juju/ratelimit](https://github.com/juju/ratelimit).

## Synopsis

```go
resp, err := http.Get("http://example.com/")
if err != nil {
	// error process
}
defer resp.Body.Close()

bucket := ratelimit.NewFromRate(
	100*1024, // limit transfer rate to 10kb/s
	100*1024, // bucket capacity, tune this if needed.
	0, // allocate this much tokens a time, 0 means capacity/2
)
r := ratelimit.NewReader(resp.Body, bucket)
io.Copy(dst, r)
```

See example folder for more.

## Total transfer rate

Bucket is thread-safe. You can share same Bucket betweens readers/writers to limit the total transfer rate.

## License

LGPL v3 or later.
