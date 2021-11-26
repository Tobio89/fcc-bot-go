# go-discord

If you DO want to try this bot, you need a folder called config, with a file called config.go, with the following code:

```go
package config

const Key = 'your bot key'
const Guild = 'your guild ID code'

```

## Code tips:

#### Event Handlers:

```go

dg.AddHandler(onReady)

func onReady...

```

The AddHandler method isn't expecting particular names, but the variables you give to the function passed in dictate what event that handler does.

#### Erasing a message:

```go
err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
    if err != nil {
        fmt.Println(err)
    }
}
```

This will erase one message.

##### Capturing a range of messages

```go
messages, _ := s.ChannelMessages(m.ChannelID, 20, msg.ID, "", "")
for _, message := range messages {
    fmt.Println(message)
}
```

Use this to capture a range of messages (20) in this case.
The empty arguments are for before or around a particular message - not that useful. Pass an empty string to skip them.
This function does not include the 'msg.ID' msg so this needs to be added after
