# This Dockerfile is only for testing PID1 handling. It's not useful for using
# Haberdasher in general.
FROM fedora:32

COPY ./foo.py ./haberdasher /app/
WORKDIR /app
ENV PYTHONUNBUFFERED=1
ENTRYPOINT ["./haberdasher"]
CMD ["python3", "foo.py"]
