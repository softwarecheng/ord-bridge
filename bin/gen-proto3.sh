#!/bin/bash

# protoc --go_out=. ../indexer/pb/*.proto
protoc --go_out=paths=source_relative:. ./indexer/pb/*.proto
