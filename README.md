# Demun - the dmenu daemon

## Wait, Why?

I run some thin clients, which are at times physically disconnected with the media I want to populate a menu with. This facilitates locality-agnostic searches for content, as well as caching the content of sometimes expensive menu creation processes. On regular computers, these tend to hold very little to no penalty; but on a very low power machine, this makes or breaks usability.

A simple `find ~/code -type f` on a bigger system, which is very convenient piped into a dmenu with a fuzzy match, would take upwards of 30 seconds on an Acer Aspire E15 I was working with. Now, a menu with > 10,000 entries was built, and runs instantly every time

This was created in response to dsearch, available from me at https://github.com/halfwit/dsearch, to allow network transparent file selection on very slow machines such as rpis; as well to integrate cleanly with a networked plumber setup.

## Usage

demun allows you to pre-compute menu entries, and share listings of file resources across the network.

### Dctl

dctl [add|list]
- `-p <port>` Port to listen on
- `-r <user@host>` Dial string for remote resource, prefixed to entries if set
- `-s <addr>` Address of running demun
- `-t <tag>` Tag to use (default "path")

dctl list will return all entries related to a given tag (by default, the "path" tag)
dctl add will read from stdin, and add entries to the given tag

### Demun

demun
- `-d` Debug mode
- `-p <port>` Port to broadcast on (Default 9997)

## BUGS

- A max amount of around 5000 entries can be added per invocation of dctl add
- 
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

Demun simply wraps our db (currently just an in memory map) and provides a simple proto to modify and query it
