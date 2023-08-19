#!/usr/bin/env sh

PGPASSWORD=postgres psql -U postgres -h localhost -p 5432 -d postgres -f script/select.sql