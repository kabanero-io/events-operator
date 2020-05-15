#!/bin/bash
host=webhook-kabanero.apps.events.os.fyre.ibm.com
curl --insecure -d '{"attr1": "val1" ,"attr2":"val2"}' -H 'Content-Type: application/json' -H 'git-enterprise: ab.cd.ef' https://${host}/webhook
