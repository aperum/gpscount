# gpscount

The gpscount plugin collects visible and used count of satellites seen by gpsd daemon.

### Configuration:
```tom
[[inputs.gpscount]]
  # gpsd daemon listening adress to connect to
  url = "localhost:2947"
```

### Metrics:

- gpscount
  -fields:
    - used (int, count): used satellites
    - visible (int, count): seen satellites

### Example Output
```
gpscount used=5,visible=10 1593980209648811323
```
