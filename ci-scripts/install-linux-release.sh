#!/usr/bin/env bash

set -x

sudo apt-get update
sudo apt-get install ruby ruby-dev rubygems build-essential rpm
sudo gem install --no-ri --no-rdoc fpm
