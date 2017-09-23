# ananke

Ananke is a sorted key/value store implemented over an append-only
filesystem such as S3. At most one process can be writing to a
store at any time. At most one thread can be writing to a store
at any time. Readers can be arbitrarily interleaved with the writer.

Ananke is implemented as a [B-epsilon tree](http://supertech.csail.mit.edu/papers/BenderFaJa15.pdf).

## Status

This is currently a hobby project. Nothing is stable.
