#!/bin/bash
curl -i -X PUT -H "Content-Type: application/json; charset=UTF-8" --data-binary @in.json  "localhost:5000/compare/?comma=%3B"
