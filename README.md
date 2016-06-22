# GoGoCal

GoGoCal (aka HoCoGoGoCal) is a Google Calendar integration written in go. It
recieves events through a Redis database and sends them to Google Calendar.

## Overview

### Process Queues

GoGoCal polls the Redis database every second for items in a set with key
`gogocal.hcpss.org:event:status:to-process`. This set contains keys of events
that need to be processed. These keys will most likely be placed there by
another integration process (Wordpress plugin, Drupal module, etc.).

GoGoCal also polls every second for items in a set with key
`gogocal.hcpss.org:event:status:to-delete`. This set contains keys of events
that need to be deleted.

Once the event is processed, it's key is placed in the set with key
`gogocal.hcpss.org:event:status:processed`. This is to make it easy for the
integration process know when an event has been processed.

If an event fails, it is placed in the set with key
`gogocal.hcpss.org:event:status:failed`. No further action is taken on the
failed events.

### Keys

All GoGoCal Redis keys are prefixed with `gogocal.hcpss.org` to prevent
collision. Keys generally follow the pattern
`<domain>:<model>:<attribute>:<value>`.

### Events

A GoGoCal Event is stored in Redis as a Hash with keys:

- source: A unique identifier for for the event source. Domain names are a good
  option. For example: events.mysite.com
- source_id: An ID for the event. When paired with the *source* this should make
  a GUID for the event. So a source_id of 72 would make the GUID
  events.mysite.com/72
- calendar: The calendar the event should be posted to. Example: joe@gmail.com
- event: The JSON encoded Google Calendar Event.
  [It's schema can be found here](https://goo.gl/fGMtP3).

## Usage

### Permission

First you need [create a service account](https://goo.gl/trXGBK) and create a
key for it. The key will be a JSON file, download and save it.

Next, you need to go to whatever Google Calendar you want to communicate with
and share it with your service account.

### Redis

GoGoCal will look for events in a Redis database, so you need to set one up.

### Parameters

`gogocal -h` prints the following, which sums it up.

```
  -a string
    	Specify the redis address. (default "redis:6379")
  -d int
    	Specify the redis database index.
  -k string
    	Specify the Google key file. (default "key.json")
  -p string
    	Specify the redis password.
  -v	Print the version.
  -version
    	Print the version.
```
