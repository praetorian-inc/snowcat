# How to contribute

The entire point of us making this available is so that we can get input and
contributions, so we are really glad you're thinking about contributing. A good
place to start thinking about places we need help is our [Issues](issues) page.

Before you begin:
- Check to see if you're on the most up to date main branch
- Have you read the [code of conduct](CODE_OF_CONDUCT.md)?

If you want to reach out to us to chat, you can email the GoKart team at
<opensource@praetorian.com>.

## Testing

Make sure your code submission runs well against some test repos. Consider
testing both positive and negative use cases - that is, if you update a
signature to find more vulnerabilities, that's great. However, **please** make
sure it doesn't find vulnerabilities where there are none.

## Submitting changes

Please send a [GitHub Pull Request to Mithril](pull/new/main) with a clear list
of what you've done (read more about [pull
requests](http://help.github.com/pull-requests/)). When you send a pull request,
you'll be forever adored if you help us also add to our tests so we can make
sure this feature or bug fix stays working. Please make sure all of your commits
are atomic (one feature per commit).

Always write a clear log message for your commits. One-line messages are fine
for small changes, but bigger changes should look like this:

    $ git commit -m "A brief summary of the commit
    >
    > A paragraph describing what changed and its impact."

## Coding conventions

Start reading our code and things should be pretty clear. We optimize for
readability and simplicity, wherever possible. Clear beats out clever every time
:)

We also use [golangci-lint](https://github.com/golangci/golangci-lint
"golangci-lint") to check our own work for some common best practices. We have
included a configuration for this linter directly into the repo. Before
submission, if you could install this linter locally and run it within the
directory you're keeping Mithril, that would be greatly appreciated. If the
linter produces no output, then you did it correctly :).

(we are open to ignoring SOME linter issues, but this will need to be a
case-by-case decision.)

## Attribution

These guidelines are very loosely adapted from <https://github.com/opengovernment/opengovernment>.
