I was excited to see a new question I think I could help answer on the community message board, so I gave answering it a shot. I ended up getting a "thank you", yeay!

# Writing data to influxdb using the API example

Getting data into a InfluxDB is fairly straight forward. Use the API, give it a metric name, some key/value pairs, and a value. The problems come usually when a system gives you a string and having to format that string to get into a desired state.

## Process

Stand up temporary docker influxdb container

```none
docker run -d --rm --name=influxdb --log-driver=json-file -p 8086:8086 influxdb:1.5.1-alpine
```

Create database

```none
curl --include -XPOST http://localhost:8086/query --data-urlencode "q=CREATE DATABASE testflux"
```

Send data to database (_note: timestamp in influx defaults to ns_ [^1])

```none
curl -i -XPOST 'http://localhost:8086/write?db=testflux' --data-binary "metrics,host=dev-node-01,service=Linux\ CPU\ check,command=check_linux,crit=99,warn=95,max=100,min=0 value=2.87 1523830298082000000"
```

Login to the container to check out the data

```none
docker exec -it influxdb sh

/ # influx
Connected to http://localhost:8086 version 1.5.1
InfluxDB shell version: 1.5.1

> use testflux
Using database testflux

> select * from metrics
name: metrics
time                command     crit host        max min service         value warn
----                -------     ---- ----        --- --- -------         ----- ----
1523830298082000000 check_linux 99   dev-node-01 100 0   Linux CPU check 2.87  95

> show tag keys from metrics
name: metrics
tagKey
------
command
crit
host
max
min
service
warn

> show field keys from metrics
name: metrics
fieldKey fieldType
-------- ---------
value    float
```

## Conclusion

I think getting data to layout correctly when passing it into the InfluxDB write query is key. Make sure to watch out for spaces. And... most likely, everything is a tag, value is a field, and timestamps are not "needed" to be passed... if one isn't given it will use the server time of when received.

Tags are indexed.

Fields are not indexed

## References

1. Influx Data Community help - https://community.influxdata.com/t/insert-data-influxdb-prase-error/4725

[^1]: https://docs.influxdata.com/influxdb/v1.5/tools/shell/#influx-arguments
