#!/usr/bin/env bash

brew uninstall --force postgresql
brew doctor
brew prune
rm -rf /usr/local/var/postgres
brew install postgresql
brew services start postgresql
