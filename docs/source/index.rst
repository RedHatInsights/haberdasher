Platform Logging with Haberdasher
=================================

`Haberdasher`_ is a project aimed at easing generation and collection of
structured logs from Platform services.

We've decided to adopt the `Elastic Common Schema`_ for our log formatting,
reasoning that Elastic has been in the event logging game long enough to have
some good ideas about formatting.

It handles a few different options for log configuration:

- By default, haberdasher attaches to the stderr stream of a containerized service and will wrap logged messages in basic ECS formatting
- Environment variables ``HABERDASHER_LABELS`` and ``HABERDASHER_TAGS`` can be defined to add pertinent information to otherwise unstructured messages
- Emitted messages that are already JSON formatted will be passed without alteration, assuming they conform to ECS standards
- Messages captured by Haberdasher are then sent to a defined ``HABERDASHER_EMITTER``, which defaults to ``stderr`` but can be set to ``kafka`` to send formatted log events to a defined ``HABERDASHER_KAFKA_TOPIC``

For examples of using haberdasher and ECS, see the `Insights-RBAC`_ repo which
we're using as a proof of concept.

.. toctree::
   :maxdepth: 1
   :caption: Documentation

   starting
   tags_and_labels
   formatting
   kafka

.. _Haberdasher: https://github.com/RedHatInsights/haberdasher
.. _Elastic Common Schema: https://www.elastic.co/guide/en/ecs/current/index.html
.. _Insights-RBAC: https://github.com/RedHatInsights/insights-rbac