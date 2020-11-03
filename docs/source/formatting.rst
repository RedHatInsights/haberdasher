Message Formatting with Haberdasher
===================================

Simply replacing the ``ENTRYPOINT`` or PID1 in a service's container with
haberdasher will provide some basic formatting turning raw messages into JSON
objects matching a basic ECS schema.

.. code-block::

  [2020-10-29 22:41:58,020] INFO: Handling signal for deleted policy Policy object (1) - invalidating associated user cache keys

.. code-block:: JSON

  {
    "ecs.version":"1.5.0",
    "@timestamp":"2020-10-29T23:00:25.628468658Z",
    "labels":{},
    "tags":[],
    "message":"INFO: Handling signal for deleted policy Policy object (1) - invalidating associated user cache keys"
  }

Moving beyond this, any JSON formatted logs will be emitted as is, under the
assumption that preformatted logs will largely fit into the ECS format. This is
where the bulk of service-side tweaking will need to happen, ensuring that any
preformatted/JSON logs fit the schema and expose the desired event data. Elastic
has developed ECS-logging libraries for several languages, including python,
javascript, java, go, and PHP. See the `Elastic github repo`_ for more information.

.. _Elastic github repo: https://github.com/elastic