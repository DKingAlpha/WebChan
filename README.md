# WebChan
A high performance web channel focused on pure data exchanging

## Usage

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
Only people 

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

### Notice

`/  ?` could not be in `channel/msg`

`&  =` could not be in `key`

