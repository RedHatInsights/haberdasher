ECS Tags and Labels
===================

The ``HABERDASHER_TAGS`` and ``HABERDASHER_LABELS`` environment variables can
be added to define additional information for the emitted log messages.

Labels
======

``HABERDASHER_LABELS`` defines a set of key-value pairs to attach to an
otherwise unstructured log message. The message itself will be added to the
``message`` field in the emitted log event, and the defined key and value to
the `ECS defined labels`_ field. Good candidates here are application/service
names, OpenShift node designators, and other identifying data.

.. code-block:: bash

  HABERDASHER_LABELS = {"service": "rbac", "hostname": $HOSTNAME}

.. code-block:: JSON

  {
      "ecs.version":"1.5.0",
      "@timestamp":"2020-09-14T16:03:02",
      "labels":{
          "service": "rbac",
          "hostname": "rbac-5d46977cb6-qz6w4"
      },
      "tags":[],
      "message":"INFO: A thing happened!"
  }

Tags
====

``HABERDASHER_TAGS`` defines a set of single keywords to attach to an emitted
log event. Much like the ``labels`` field above, this list of tags will be
added to the `ECS defined tags`_ field in the emitted message. This should be
used for identifying information that doesn't necessarily need a clarifying
'key' entry. Good candidates are environment identifiers like prod/stage/ci.

.. code-block:: bash

  HABERDASHER_TAGS = ["prod"]

.. code-block:: JSON

  {
      "ecs.version":"1.5.0",
      "@timestamp":"2020-09-14T16:03:02",
      "labels":{},
      "tags":["prod"],
      "message":"INFO: A thing happened!"
  }

.. _ECS defined labels: https://www.elastic.co/guide/en/ecs/current/ecs-base.htm

.. _ECS defined tags: https://www.elastic.co/guide/en/ecs/current/ecs-base.htm