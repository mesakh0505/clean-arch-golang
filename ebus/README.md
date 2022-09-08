# ebus

Ebus or event-bus used for your system event bus. If you want to achieve an event-driven programming, you can use this package. You can publish an event and also define the subscriber to that event you published. 

Please keep in mind, this is only for system event, means all event that happened in your applications you can publish between function. You can also combined with pub/sub system for reliability

If don't know how to use, just ignore it.

## Ebus in concept

If you don't understand clearly about ebus, this is how it works if we drawn in the picture. So with the ebus, no matter what is your Pub/Sub technology, you can just add it and inject it to the Ebus handler. It will make the system more scalable, and lously coupling with any pub/sub technology.
![Ebus](https://user-images.githubusercontent.com/11002383/66490820-f6270300-eadb-11e9-9f3f-7aaedcd91e94.png)

There's an article about making an ebus in Golang too here: https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997
But in that article, the logic is just too complex for beginner, because using Goroutine and channel.

You can follow that articles, if you want to try using channel and goroutine. Here the ebus is only a simple function.