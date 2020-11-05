Haberdasher Logging via Kafka
=============================

The end goal of all of this Haberdashing is to form a log pipeline via Kafka
to a searchable log aggregator. Each service, via haberdasher, will be a Kafka
producer, sending its logs as messages to a defined Kafka logging topic. As the
consumer of that topic, we'll have Logstash instances collecting, sorting, 
partitioning, and shipping logs out to Elasticsearch/Kibana, Splunk, Cloudwatch,
etc. This way, each service can have just one configuration, "Use haberdasher to
send logs to Kafka", and the rest can all be handled by the Logstash nodes.

By default, haberdasher sends all emitted logs to stderr. From there, fulentd or
Cloudwatch or whatever configured log handler can pick them up. By changing the
``HABERDASHER_EMITTER``, ``HABERDASHER_KAFKA_BOOTSTRAP``, and
``HABERDASHER_KAFKA_TOPIC`` environment variables, haberdasher can be configured
to connect to an existing Kafka cluster and topic for log message production.
It will handle the rest.