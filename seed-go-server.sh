#!/bin/sh
curl -H "Content-Type: application/json" -XPOST localhost:8080/seed --data-binary "@-"
