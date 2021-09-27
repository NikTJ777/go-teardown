# go-teardown
Teardowns are named, globally-scoped, deferred function lists that teardown resources that were created in a different scope.

The primary usage of `teardowns` is in testing, although they are also useful in production code.

The prblem that `teardowns` solve is that the go `defer` statement is bound to the enclosing scope only.

Additional arrangements must be made if a resouce is created in a sub-function, and not released
until the calling function has returned.

The Go idiom is to return a cleanup func from the sub-function, which the caller should then `defer`.
This is already fragile, as the caling code may ignore the cleanup func.
And if there are additional call levels between the creation of the resource, and its release, then this becomes more difficult.

To allow independent modules or sub-systems to perform `teardown` independently, each `teardown list` is named, and is then executed independently of other lists.

In testing, it is often useful for each sub-test to have an independent `teardown list`.

A simple example is:
```
defer teardown.teardown("module-1")

// call functions in module-1 that allocate resources.
// These functions call teardown.AddTeardown("module-1", func() {...})
```
When the outer scope exits, all resources added for `teardown` in `module-1` are released.

