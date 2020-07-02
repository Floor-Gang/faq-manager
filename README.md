# FaQ Manager
This is Floor Gang bot, this gives commands to manage an FaQ channel

## Usage
```
$ go build
$ ./faq-manager 
$ ... edit config.yml ...
$ ./faq-manager
```

## Bot Usage
 - .faq get: Get a question based on provided context
   - `.faq get <question context>`
 - .faq set: Set a question's answer
   - `.faq set <question context> NEWLINE <new answer>`
 - .faq add: Add a new question
   - `.faq add <question> NEWLINE <new answer>`
 - .faq list: Send all the currently stored FaQ embeds
   - `.faq list`
 - .faq remove: Remove a question
   - `.faq remove <question context>`
 - .faq sync: Sync the FaQ channel with the stored FaQ's
   - `.faq sync`
