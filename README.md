
# FusionBulk

REST API for connecting FusionPBX to BulkVS M/SMS Services.


## Acknowledgements

 - [BulkVS Documentation](https://portal.bulkvs.com/api/v1.0/documentation)
 - [gorilla/mux](https://github.com/gorilla/mux)
 - [Sam Hofius' Keybase Library](samhofi.us/x/keybase/v2)
 - [CallBreezy's VoIP Services](https://callbreezy.com/)


## Features

- Extensible via REST
- FusionPBX Web Interface for sending SMS


## Badges
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)


## Deployment

Deploying this project is a bit convoluted, because it requires redirecting a path in your nginx config for your FusionPBX site:

```       
 location = /bulkvs/webhook {
                allow 52.206.134.245;
                allow 192.9.236.42;
                deny all;
                proxy_pass http://127.0.0.1:8080/bulkvs/webhook;
        }
```

1. The block above must be added to `/etc/nginx/sites-enabled/fusionpbx` and nginx must be restarted.

2. Then you must point your BulkVS Portal -> Messaging -> Messaging Webhooks to https://your-fusionpbx-site/bulkvs/webhook

3. At the moment, you may need to remove/recompile without notify.go, because I wrote that as a quick Keybase interface while I worked out the fusionPBX interface. 

With the current configuration, you can only send messages from FusionPBX. However, if you do have access to [Keybase](https://keybase.io), sending and receiving SMS works fine.
## Related

Here are some related projects

[BulkVS2Go Library](https://github.com/rudi9719/BulkVS2Go)


## Support

For support, contact [rudi](https://rudi.nightmare.haus)

