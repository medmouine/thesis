option task = {name: "cloud-replication", every: 6h}

from(bucket: "iot-data")
    |> range(start: -6h)
    |> filter(fn: (r) => r["_measurement"] == "sensor")
    |> filter(fn: (r) => r["_field"] == "humidity" or r["_field"] == "temperature")
    |> group(columns: ["_field"])
    |> aggregateWindow(every: 10m, fn: mean, createEmpty: false)
    |> set(key: "_measurement", value: "sensors")
    |> to(bucket: "iot-datalake", host: "http://cloud-datalake.cloud", token: "influx_token", org: "cloudorg")


from(bucket: "iot-data")
  |> range(start: -24h)
  |> filter(fn: (r) => r["_measurement"] == "sensor")
  |> filter(fn: (r) => r["_field"] == "temperature" or r["_field"] == "humidity")
  |> group(columns: ["mapper_id", "_field"])
  |> timeWeightedAvg(unit: 1h)
  |> set(key: "_measurement", value: "sensors")
  |> set(key: "location", value: "dc1")
  |> duplicate(column: "_stop", as: "_time")
  |> to(bucket: "w-data", host: "http://cloud-datalake.cloud", token: "influx_token", org: "cloudorg")

