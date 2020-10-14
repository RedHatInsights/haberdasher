# This Dockerfile is only for testing PID1 handling. It's not useful for using
# Haberdasher in general.
FROM fedora:32

RUN dnf install -y bind-utils net-tools python3-requests && dnf clean all
COPY ./foo.py ./haberdasher /app/
WORKDIR /app
ENV PYTHONUNBUFFERED=1 HABERDASHER_STDERR_PRETTY=1
ENTRYPOINT ["./haberdasher"]
CMD ["python3", "foo.py", "--serve"]
