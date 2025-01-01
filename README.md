# git-commit-ai

Proof of concept right now. This is a git extension that will look at your staged files and analyse the changes within. It will then attempt to write a good commit message for you. You will have a chance to edit it.

Currently requires an Anthropic API key, I plan to add a few other adapters. This was literally hacked together in about 2 hours so don't expect miracles.

## Instructions

Download the binary from GitHub or you can compile yourself using:

```
go build -o git-commit-ai
```

Once you have a binary you can simply put it somewhere in your `PATH`. Stage some files in git and then run:

```
git commit-ai
```

It should open in your default editor with a commit message pre-populated. You can edit it and then save and quit. The temporary commit file will be deleted and you'll be asked to confirm one more time, then it'll automatically commit with that message.
