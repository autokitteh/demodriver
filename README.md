# Temporal Demo Driver

```
$ make bin
$ ./bin/dd
```

Default port is 9001.

## Configuration

### Triggers path

```
DD_DRIVER_TRIGGERSPATH="examples/triggers.yaml"
```

### Slack source

```
DD_SLACKSOURCE_BOTTOKEN="xoxb-...."
DD_SLACKSOURCE_APPTOKEN="xapp-..."
```

## API

Post this form the workflow:

```
POST host:9001/api/signals
{
     "name": "signal name",
     "wid": "temporal workflow id to receive signal",
     "src": "slack|http",
     "filter": "..."
}
```

And when a matching event is received by source, DD will signsl the workflow with the trigger name and the event data.
