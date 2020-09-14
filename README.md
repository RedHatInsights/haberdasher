# Haberdasher

Haberdasher is a simple command wrapper designed to consume log messages from
its wrapped command's stderr stream and retransmit them someplace else. It also
contains standard PID1 goodies for healthier container execution.

## Simple demonstration

The `foo.py` program simply ticks off integers as log messages every 2 seconds.

    $ PYTHONUNBUFFERED=1 ./haberdasher python3 foo.py
    2020/09/14 16:03:00 Initializing haberdasher.
    2020/09/14 16:03:00 Configured emitter: stdout
    Python starting
    {"ecs.version":"1.5.0","@timestamp":"2020-09-14T16:03:02.556065987-04:00","labels":{},"tags":[],"message":"0"}
    {"ecs.version":"1.5.0","@timestamp":"2020-09-14T16:03:04.558082983-04:00","labels":{},"tags":[],"message":"1"}
    {"ecs.version":"1.5.0","@timestamp":"2020-09-14T16:03:06.560023837-04:00","labels":{},"tags":[],"message":"2"}
    ^C2020/09/14 16:03:07 Signal received: interrupt
    2020/09/14 16:03:07 Sending signal to 415770
    2020/09/14 16:03:07 Trigering emitter shutdown

You can see that using the stdout emitter, it simply prints the received messages.
Since the output of `foo.py` was unstructured, each log line that Haberdasher
received is wrapped in a basic [Elastic Common Schema](https://www.elastic.co/guide/en/ecs/current/index.html)
message.

If Haberdasher receives a structured log message from its wrapped process, it
leaves it alone and retransmits it unmodified.

    $ PYTHONUNBUFFERED=1 ./haberdasher python3 foo.py --json
    2020/09/14 16:05:02 Initializing haberdasher.
    2020/09/14 16:05:02 Configured emitter: stdout
    Python starting
    {"i": 0}
    {"i": 1}
    {"i": 2}
    ^C2020/09/14 16:05:09 Signal received: interrupt
    2020/09/14 16:05:09 Sending signal to 416367
    2020/09/14 16:05:09 Trigering emitter shutdown

## Configuring Haberdasher

Haberdasher is configured entirely from environment variables.

* `HABERDASHER_EMITTER` - configures the emitter to use. `stdout` is default,
  but `kafka` is also supported.
* `HABERDASHER_TAGS` - for unstructured log lines received, Haberdasher can add
  ECS tags to the wrapped messages. This value should be a serialized JSON list.
* `HABERDASHER_LABELS` - for unstructured log lines received, Haberdasher can
  add ECS labels to the wrapped messages. This value should be a serialized
  JSON object whose values are all strings.
* `HABERDASHER_KAFKA_BOOTSTRAP` - if the `kafka` emitter is used, this is
  required and points to the bootstrap listener for your Kafka cluster
* `HABERDASHER_KAFKA_TOPIC` - if the `kafka` emitter is used, this is required
  and names the Kafka topic log messages should be written to

