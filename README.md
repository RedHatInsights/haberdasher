# Haberdasher

Haberdasher is a simple command wrapper designed to consume log messages from
its wrapped command's stderr stream and retransmit them someplace else. It also
contains standard PID1 goodies for healthier container execution.

## Simple demonstration

The `foo.py` program simply ticks off integers as log messages every 2 seconds.

    $ ./haberdasher python3 foo.py
    2020/09/14 16:03:00 Initializing haberdasher.
    2020/09/14 16:03:00 Configured emitter: stderr
    Python starting
    {"ecs.version":"1.5.0","@timestamp":"2020-09-14T16:03:02.556065987-04:00","labels":{},"tags":[],"message":"0"}
    {"ecs.version":"1.5.0","@timestamp":"2020-09-14T16:03:04.558082983-04:00","labels":{},"tags":[],"message":"1"}
    {"ecs.version":"1.5.0","@timestamp":"2020-09-14T16:03:06.560023837-04:00","labels":{},"tags":[],"message":"2"}
    ^C2020/09/14 16:03:07 Signal received: interrupt
    2020/09/14 16:03:07 Sending signal to 415770
    2020/09/14 16:03:07 Trigering emitter shutdown

You can see that using the stderr emitter, it simply prints the received messages.
Since the output of `foo.py` was unstructured, each log line that Haberdasher
received is wrapped in a basic [Elastic Common Schema](https://www.elastic.co/guide/en/ecs/current/index.html)
message.

If Haberdasher receives a structured log message from its wrapped process, it
leaves it alone and retransmits it unmodified.

    $ ./haberdasher python3 foo.py --json
    2020/09/14 16:05:02 Initializing haberdasher.
    2020/09/14 16:05:02 Configured emitter: stderr
    Python starting
    {"i": 0}
    {"i": 1}
    {"i": 2}
    ^C2020/09/14 16:05:09 Signal received: interrupt
    2020/09/14 16:05:09 Sending signal to 416367
    2020/09/14 16:05:09 Trigering emitter shutdown

## Configuring Haberdasher

Haberdasher is configured entirely from environment variables.

* `HABERDASHER_EMITTER` - configures the emitter to use. `stderr` is default,
  but `kafka` is also supported.
* `HABERDASHER_TAGS` - for unstructured log lines received, Haberdasher can add
  ECS tags to the wrapped messages. This value should be a serialized JSON list.
* `HABERDASHER_LABELS` - for unstructured log lines received, Haberdasher can
  add ECS labels to the wrapped messages. This value should be a serialized
  JSON object whose values are all strings.
* `HABERDASHER_STDERR_PRETTY` - if the `stderr` emitter is used, setting this to
  a non-empty string will result in the JSON being prettified before printing to
  stderr. This is useful in developer environments to make the messages easier
  to read.
* `HABERDASHER_KAFKA_BOOTSTRAP` - if the `kafka` emitter is used, this is
  required and points to the bootstrap listener for your Kafka cluster
* `HABERDASHER_KAFKA_TOPIC` - if the `kafka` emitter is used, this is required
  and names the Kafka topic log messages should be written to

## Adding it to your Dockerfile

To use Haberdasher in a container, you only have to make two small modifications
to your Dockerfile.

1. In a `RUN` stanza, include:

    curl -L -o /usr/bin/haberdasher https://github.com/RedHatInsights/haberdasher/releases/latest/download/haberdasher_linux_amd64 && \
    chmod 755 /usr/bin/haberdasher

2. Your `ENTRYPOINT` command should be: `["/usr/bin/haberdasher"]`

And that's it! Rebuild and you're up and running.
