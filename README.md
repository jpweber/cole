# Cole 
> I see dead people  

Cole is a dead man switch listener. In prometheus it is common to create a _dead man switch_ which will constantly send alerts to test your entire alerting pipline. A question that comes up often is what do you have watching those _dead man switch_ alerts. Who watches the watchers, effectively.  

This is a basic implmentation of something that could watch for those deadman switch alerts, and then send alert itself if it does not receive a notification from the deadman switch within the assigned time interval.

## Status
this project is in very early stages and should not be used in production _yet_. This is Still in Work In Progress (WIP) status that does work but there are some planned features that still need to be added and things like configuration are still evolving.

## How does it work
Cole listens for http requests from prometheus alertmanager sending alerts for dream switch alert. When a message is received a timer will be started for the specified duration. If a message is not received from the deadman alert inside of that time duration, it will fire off an alert of it's own.

There is a forthcoming blog post on [jpweber.io](http://jpweber.io/blog
) on how to leverage a deadman switch alert in your prometheus monitoring and how something like Cole fits in which will provide some more detail in to the thinking of creating a tool like this.

## Supported alert integrations

* Slack
* PagerDuty
* Generic Webhook

## How to use

1) Start the cole server by any of the below defined means (bare binary, docker, etc)

2) For each DeadManSwitch that you want to check in you must generate an ID for that alert. Perform an http `GET` request to `/id` of the cole server. For example. `curl http://yourcoleaddress/id`. This will return a json payload of the following. This timerid will be part of the url you hit to check in.

    ``` json
    {
        "timerid":"bg8obqel0s1fdr02gtvg"
    }
    ```

3) Create a receiver in your alert manager config to make a call to a webhook when it recieves a DeadManSwitch alert. The wait, group and repeat intervals may need to be changed based on your needs.

    ``` yaml
    global:
     ...
    route:
     ...
        routes:
        - match:
            alertname: DeadMansSwitch
            receiver: 'cole'
            group_wait: 0s
            group_interval: 1m
            repeat_interval: 50s
    receivers:
    - name: 'cole'
    webhook_configs:
    - url: 'http://192.168.2.66:8080/ping/bg8obqel0s1fdr02gtvg'
        send_resolved: false
    ```


## Configuration

### Example using configuration file

``` 
# Example Cole configuration file

# Slack
# SenderType = "slack"
# Interval = 10
# HTTPEndpoint = "https://hooks.slack.com/services/..."
# HTTPMethod = "POST"
# SlackChannel = "#general"
# SlackUsername = "Cole - DeadManSwitch Monitor"
# SlackIcon = ":monkey_face:"


# PagerDuty
SenderType = "pagerduty"
Interval = 10
PDAPIKey = "noiD8-khbpNpgAAAAAAAAAA"
PDIntegrationKey = "5353fb993888441811111111111"
```

### Flags supported as ENV Vars

* `SENDER_TYPE`
* `INTERVAL`
* `HTTP_ENDPOINT`
* `HTTP_METHOD`
* `EMAIL_ADDR`
* `PD_KEY`
* `SLACK_CHANNEL`
* `SLACK_USERNAME`
* `SLACK_ICON`

## Example Prometheus Alert Manager config



## Run it

### With docker

``` shell
docker run -d \
-e SENDER_TYPE="slack" \
-e INTERVAL="10" \
-e HTTP_ENDPOINT="https://hooks.slack.com/services/..." \
-p 8080:8080 \
cole:0.2.0
```

### Bare binary

`./cole`

## API Endpoints

* `POST` - `/ping/<timerid>`
* `GET` - `/id`
* `GET` - `/version`

## Build locally

* clone the repo
* `dep ensure -v`
* `go build`
That is it.
