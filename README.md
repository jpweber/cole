# Cole 
> I see dead people  

Cole is a dead man switch listener. In prometheus it is common to create a _dead man switch_ which will constantly send alerts to test your entire alerting pipline. A question that comes up often is what do you have watching those _dead man switch_ alerts so you know when something is wrong with your alerting system. Who watches the watchers, effectively.  

This is a basic implmentation of something that could watch for those deadman switch alerts, and then send alert itself if it does not receive a notification from the deadman switch within the assigned time interval. 

## Status
this project is in very early stages and should not be used in production _yet_. So far it is a proof of concept that does work but there are some major features that still need to be added and are planned. 

### Current Limitations
These limitations are because of the current state and will be supported  

* only has awareness of one deadman alert and its corresponding duration. Meaning, you can't have this watch multiple things that could have different durations
* only reads from Command line (CLI args) and not from a config file. 
* only knows how to send alerts to http endpoints that accept a json payload. 

## How does it work
When you start cole up, one of the required flags is a time duration. Cole will start a timer for that specified duration and if doesn't recieve a message from the deadman alert inside of that time duration, it will fire off an alert of it's own. There is a forthcoming blog post on [jpweber.io](http://jpweber.io/blog
) on how to leverage a deadman switch alert in your prometheus monitoring and how something like Cole fits in which will provid some more detail in to the thinking of creating a tool like this. 

## Simple way to see that its working
`./cole -b "hello world" -e "https://reqres.in/api/users" -m POST -s "dev prometheus" -t 10`
