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

## Configuration file
```
# Example Cole configuration file

# Slack
# SenderType = "slack"
# Interval = 10
# HTTPEndpoint = "https://hooks.slack.com/services/..."
# HTTPMethod = "POST"


# PagerDuty
SenderType = "pagerduty"
Interval = 10
PDAPIKey = "noiD8-khbpNpgAAAAAAAAAA"
PDIntegrationKey = "5353fb993888441811111111111"
```

## Simple way to see that its working
`./cole -c example.toml`


## Build locally
* clone the repo
* `dep ensure -v`
* `go build`
That is it. 
