Buffers
=======

Within Benthos is support for multiple buffering options. This document outlines
the various types of buffer along with potential use cases and configuration
examples.

## None

```yaml
buffer:
  type: none
```

Benthos works well without a buffer, this means the rate at which messages are
written will limit the rate messages are read, applying back pressure to the
input where applicable. This is perfectly acceptible for cases where Benthos is
required only as a bridge between two other message protocols and benefits such
as surge protection or persistence are not needed.

## Memory

Benthos can use RAM to buffer messages whenever the output applies back
pressure. This helps in situations where the output service is likely to reach
throughput capacity during its lifetime (during surges in data etc), or if the
output service is likely to need restarts and the input service needs to be
flushed.

Using memory is the lowest latency buffer option but will result in data loss if
the service is restarted without being fully flushed.

```yaml
buffer:
  type: memory
  memory:
    limit: 524288000
```

The `limit` field sets how many bytes the buffer can reach before blocking new
messages.

## Memory-Mapped File

Memory-Mapped files are blocks of data stored in RAM that are flushed to disk by
the host operating system. This means the service is not blocked on file IO
during the write and that messages are ready to dispatch to the output service
much sooner.

For most use cases this is the best option as it provides persistence across
service/operating system restarts and has a low latency overhead.

However, this option does not add protections around service or operating system
crashes, as there is no guarantee during runtime of exactly when the message is
flushed to disk.

```yaml
buffer:
  type: mmap_file
  mmap_file:
    directory: ""
    file_size: 262144000
    retry_period_ms: 1000
    clean_up: true
```

The `directory` field is a directory that will store all mmap buffer files,
Benthos attempts to create the entire path if it does not exist.

The `file_size` field is the maximum size each individual mmap file should be,
Benthos creates these files as it requires them and therefore this value does
not limit the capacity of the buffer. This size *must* be larger than any
messages coming through the pipeline.

The `retry_period_ms` field sets the period in milliseconds before reattempting
to open or create an mmap file after failing.

The `clean_up` field is a boolean that indicates whether Benthos should delete
old mmap files when it is finished with them. You can set this to false if you
wish to keep all content and manually clean up the directory.
