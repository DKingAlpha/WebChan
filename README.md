# WebChan
A high performance web channel focused on pure data exchanging

# Public Access
```
POST /channel/msg
GET /channel
GET /channel/from
GET /channel/from/to
DELETE /channel
```


* channel: channel id, string
* msg: message, string
* from: unix timestamp, decimal
* to: unix timestamp, decimal

# Restricted Access
Append arguments to url to post/get/delete restricted channel

## Avaliable Parameters

### key
string

1. POST with a key to post to a restricted channel, if it does not exists, create it.
2. GET/DELETE a restricted channel


### perm
Set corresponding char grant access to channel to people WITHOUT key:
* r: read
* w: write
* d: delete

By default, perm is empty, which means no one could r/w/d this channel without a key.

For example:
```
# any one could r/w this channel. but they cant delete
# if this channel exists, perm will be updated.
POST /chan123/hello?key=password&perm=rw

# any one could delete this, but no one could r/w without the key.
POST /chan456/hi?key=password&perm=d

# read ok
GET /chan456?key=password
# read failed
GET /chan456
# delete ok
DELETE /chan456
```

### time
Set time parameter to GET with timestamp in output
```
GET /chan123?time
GET /chan123?key=123&time
GET /chan123?key=123&time=1
# warn: time=0 or time=false will not disable it
GET /chan123?key=123&time=1
```


### websocket
connect to `ws://ip:port/websocket/channel_name` to subscribe new data.

## Notice

`/  ?` could not be in `channel/msg`

`&  =` could not be in `key`

