Adding Haberdasher to A Service
===============================

Haberdasher is designed as a PID1 replacement for OpenShift containers. As the
haberdasher docs state, two changes need to be made to the service's Dockerfile
and it will handle the rest:

.. code-block:: dockerfile

  RUN curl -L -o /usr/bin/haberdasher https://github.com/RedHatInsights/haberdasher/releases/latest/download/haberdasher_linux_amd64 && chmod 755 /usr/bin/haberdasher

  ENTRYPOINT ["/usr/bin/haberdasher"]

After adding the above and rebuilding, you should see logs emitted to stderr
with haberdasher's formatting included.