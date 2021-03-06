A simple overview of keeping prometheus data for longer than the default 15 days.

## Prometheus long term retention

By default, Prometheus time series database (TSDB) is set to keep its collection of data for 15 days[^1].
This is usually enough data for metric based monitoring, but what happens when management wants monthly numbers?
We can send the Prometheus data into Influx Data's database, InfluxDB.

InfluxDB is a TSDB as well. It is based on timestamp=value pairs that are collected within a Measurement.
A Measurements is similar to a SQL table.
InfluxDB also has Series. Series are essentially every row, under a Measurement, with the same tags and fields.
These tags and fields have to have the same values.
With Prometheus, this concept does transfer one to one so bumping up against InfluxDB defaults can happen depending upon the number of unique items being collected.

## Prometheus

We can use Prometheus to look at the data with the /graph landing page. If we select a metric like "up" we can see the response is:

| element | value |
|:---|:---:|
| up{instance="localhost:9090",job="prometheus"} | 1 |

_Note: The Prometheus graph will show the latest entry for the metric selected unless we pass an amount of time._

## InfluxDB

We can interact with InfluxDB by either using "influx" at the CLI or we could use the API to make calls. I find interacting with the CLI to be proficient for messing around figuring things out while the API can be programmed... your call.

```none
# influx
```

```none
> select * from prometheus.autogen.up limit 1
name: up
time                instance       job        value
----                --------       ---        -----
1523215564367000000 localhost:9090 prometheus 1
```

## Good to knows

Bring up InfluxDB running config (_note: in container_)

```none
# influxd config
```

Show the InfluxDB stats (_note: in container, within 'influx' service_)

```none
> show stats
```

## Conclusion

Prometheus default retention period could be extended to collect data for how ever long you wish. But even Prometheus themselves don't recommend keeping data in their database for long. They do actually recommend storing it somewhere else for long term retention. I'm only appending about 18,000 metric scrapes per second right now (at work), and I'm hitting walls that I have to troubleshoot. Using InfluxDB to house long term retention and down sample data has helped out a lot. I still use Promethues for about 30 days of 100% of the collection.

I will continue to post on this topic and my findings. Sorry for the short post.

## References

1. Prometheus.io
1. Influxdata.com - (https://docs.influxdata.com/influxdb/v1.5/)

[^1]: https://prometheus.io/docs/prometheus/latest/storage/#operational-aspects
