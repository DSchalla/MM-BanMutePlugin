#!/bin/bash
set -e
set -x
./package.sh;
rm ../../../mattermost/mattermost-server/plugins/com.dschalla.banmuteplugin/banmuteplugin;
cp banmuteplugin ../../../mattermost/mattermost-server/plugins/com.dschalla.banmuteplugin/banmuteplugin;