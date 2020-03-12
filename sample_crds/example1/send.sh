#!/bin/bash
host=webhook-default.apps.repines.os.fyre.ibm.com
curl --insecure -d '{"attr1": "val1" ,"attr2":"val2"}' -H 'Content-Type: application/json' https://${host}/webhook
