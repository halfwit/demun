# Demun - the dmenu daemon

This listens on port 9997, and can be queried with dctl. the output of dctl can be passed to a dmenu / bemenu / rofi instance

## Wait, Why?

I run some thin clients, which are at times physically disconnected with the media I want to populate a menu with. This facilitates locality-agnostic searches for content, as well as caching the content of sometimes expensive menu creation processes. On regular computers, these tend to hold very little to no penalty; but on a very low power machine, this makes or breaks usability.

This was created in response to dsearch, available from me at https://github.com/halfwit/dsearch, to allow network transparent file selection on very slow machines such as rpis; as well to integrate cleanly with a networked plumber setup. 

## Proposed Usage

Prior to writing, it's fun to brainstorm how things will work out. Here's a bit of that.


### Dctl

Ideally, I'd like to pipe a bunch of newlines into dctl, optionally adding a host prefix to 
the menu items. 

```/bin/sh
# Add in files from our code directory, with a host prefix
find ~/code -type f -not -path '*/\.git*' | dctl -r myhost -t search add

# Add in some fancy query prefixes that my plumber knows how to handle 
printf '%s\n%s\n%s\n' "!yt" "!pl" "!g" "!ddg" | dctl -t search add
```

Then, getting them back out could be something as follows.

```/bin/sh
# list all files belonging to a specific tag
dctl -t search list | dmenu -p 'select file' | plumb -i
```

### Demun

We're wrapping github.com/takama/daemon for much of the heavier lifting. It's prescribed usage is simple for starting/stopping the daemon service

```/bin/sh

demun start - start the service
demun stop - stop the service
demun remove - remove the service
demun status - return the status of the service
demun install - create the service 

```

