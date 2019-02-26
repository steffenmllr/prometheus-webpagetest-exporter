# prometheus-webpagetest-exporter [![](https://images.microbadger.com/badges/image/steffenmllr/prometheus-webpagetest-exporter.svg)](https://hub.docker.com/r/steffenmllr/prometheus-webpagetest-exporter/ "This image on Docker Hub")
> A runner that periodically tests a url with [webpagetest.org](https://www.webpagetest.org) and exposes those values to [prometheus.io](https://prometheus.io)

## Gettings Started
1. Get a Key from [webpagetest.org](https://www.webpagetest.org/getkey.php) or clone this repo and do a `docker-compose up -d`
2. Create a config and check out the [sample-config.toml](./config.toml):

```toml
# the Key from your local or remote webpagetest instance
key = "123"
# the URL to your webpagetest instance, if you run locally (and on a mac, use your real ip)
host = "http://127.0.0.1:4000"
# the Port where the service is running and the metrics are exposed
port = "3030"
# How often do you want to run the test
timer = "1h"

# An array of metrics you want to expose with prometheus
# Check out the
[[metrics]]
# Prometheus key export
key = "time_to_first_byte_ms"
# Prometheus help
help = "First Byte in ms."

# The Data from the reponse that gets collected, the array index is set as run
# Find the Reponse here: https://www.webpagetest.org/jsonResult.php?testid=190120_HQ_e97dbb371e61dd8fcf46c4feda8ddaec
data = [
    "median.firstView.TTFB",
    "median.repeatView.TTFB"
]

# An Array of sites / urls you want to tests
[[sites]]
name = "vitra home"
url = "https://www.bbc.co.uk"
# Test if you run locally with docker-compose, see https://www.webpagetest.org/getLocations.php?f=html and https://sites.google.com/a/webpagetest.org/docs/advanced-features/webpagetest-restful-apis
location = "Test.LAN"
```

3. Run the docker container (run with to access the local instance)
```
docker run -v $(pwd)/config.toml:/config.toml -it -p 3030:3030 steffenmllr/prometheus-webpagetest-exporter /config.toml
```

4. Check the output at [https://localhost:3030](https://localhost:3030)

5. Add to Prometheus to scrape

### Sample metrics

```text
# HELP page_requests_no Page Size
# TYPE page_requests_no gauge
page_requests_no{location="Test",run="1",url="https://www.bbc.co.uk"} 81
page_requests_no{location="Test",run="2",url="https://www.bbc.co.uk"} 7

# HELP page_size_bytes Page Size
# TYPE page_size_bytes gauge
page_size_bytes{location="Test",run="1",url="https://www.bbc.co.uk"} 49001
page_size_bytes{location="Test",run="2",url="https://www.bbc.co.uk"} 112527

# HELP speed_index_no Speed Index
# TYPE speed_index_no gauge
speed_index_no{location="Test",run="1",url="https://www.bbc.co.uk"} 1617
speed_index_no{location="Test",run="2",url="https://www.bbc.co.uk"} 1102

# HELP time_first_byte First Byte in ms.
# TYPE time_first_byte gauge
time_first_byte{location="Test",run="1",url="https://www.bbc.co.uk"} 138
time_first_byte{location="Test",run="2",url="https://www.bbc.co.uk"} 121

# HELP time_load_time Load time in ms.
# TYPE time_load_time gauge
time_load_time{location="Test",run="1",url="https://www.bbc.co.uk"} 1657
time_load_time{location="Test",run="2",url="https://www.bbc.co.uk"} 1163

# HELP time_start_render_ms Start render time.
# TYPE time_start_render_ms gauge
time_start_render_ms{location="Test",run="1",url="https://www.bbc.co.uk"} 500
time_start_render_ms{location="Test",run="2",url="https://www.bbc.co.uk"} 300

# HELP time_visually_complete_ms Visually Complete
# TYPE time_visually_complete_ms gauge
time_visually_complete_ms{location="Test",run="1",url="https://www.bbc.co.uk"} 2400
time_visually_complete_ms{location="Test",run="2",url="https://www.bbc.co.uk"} 1600
```
