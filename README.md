# Demun - the dmenu daemon

This listens on port 9997, and can be queried with dctl. the output of dctl can be passed to a dmenu / bemenu / rofi instance

## Wait, Why?

I run some thin clients, which are at times physically disconnected with the media I want to populate a menu with. This facilitates locality-agnostic searches for content, as well as caching the content of sometimes expensive menu creation processes. On regular computers, these tend to hold very little to no penalty; but on a very low power machine, this makes or breaks usability.

This was created in response to dsearch, available from me at https://github.com/halfwit/dsearch, to allow network transparent file selection on very slow machines such as rpis; as well to integrate cleanly with a networked plumber setup. 
