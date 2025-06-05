#!/bin/bash

openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 \
  -subj "/CN=localhost" \
  -keyout /etc/kong/cert.key \
  -out /etc/kong/cert.crt

