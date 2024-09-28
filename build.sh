#!/bin/bash

echo 'Запущена сборка преокта...'
docker build -f ./docker/build.Dockerfile --network host --target builder -t diploma_builder .
docker run --rm -v .:/opt/diploma diploma_builder

