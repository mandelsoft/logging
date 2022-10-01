# Logging for GO with Selective Log Levels

This package provides a simple wrapper around the [logr](https://ggithub.com/go-logr/logr)
logging system supporting a rule based approach to enable log levels
for dedicated logging requests, specified at the logging location.

The rule set is configured for a logging context:

```go
    ctx := logging.New(logrLogger)
```

Any`logr.Logger` can be passed here, the level for this logger
globally filters the log levels of the provided log messages.
If the full control should be handed over to the log context, 
the maximum log level should be used for this logger.

Now you can add rules controlling the accepted log levels for dedicated log 
locations. First, a default log level can be set:

```go
    ctx.SetDefaultLevel(logging.InfoLevel)
```

This level restriction is used, if no rule matches a dedicated log request.

Another way to achieve the same goal is to provide a generic rule without any condition:

```go
    ctx.AddRule(logging.NewConditionRule(logging.InfoLevel))
```

A first rule for influencing the log level could be a realm rule.
A *Realm* represents a dedicated logical area, a good practice could be 
to use package names as realms. Realms are hierarchical consisting of
name components separated by a slash (/).

```go
    ctx.AddRule(logging.NewConditionRule(logging.DebugLevel, logging.NewRealm("github.com/mandelsoft/spiff")))
```

Alternatively `NewRealmPrefix(...)` can be used to match a complete realm hierarchy.

A realm for the actual package can be defined as local variable.

```go
var realm = logging.Package()
```

Instead of passing `Logger`s around, now the logging `Context` is used.
It provides a method to access a logger specific for a dedicated log
request, for example, for a dedicated realm.

```go
  logctx.Logger(realm).Info("my message")
```

The provided logger offers the level specific functions, `Error`, `Warn`, `Info`, `Debug` and `Trace`.
Depending on the rule set configured for the used log context, the determined level
decides, which message to pass to the log sink of the initial `logr.Logger`.

If no rules are configured, the default logger of the context is used independently of the
given arguments. The given context information is optionally passed to the
provided logger, depending on the used context type.

For example, the realm is added to the logger's name.

It is also possible to provide dedicated attributes for the mapping process:

```go
  logctx.Logger(realm, logging.NewAttribute("test", "value")).Info("my message")
```

Such an attribute can be used as rule condition, also. This way, logging
can be enabled, for dedicated argument values of a method/function.

Both sides, the rule conditions and the log context can be a list.
For the conditions, all specified conditions must be evaluated to true, to
enable the rule. A rule is evaluated against the complete log requests.
The default ` ConditionRule` evaluates the rules against the complete log
request and a condition is *true*, if it matches at least one argument.

The rules are evaluated in the reverse order of their definition.
The first matching rule defines the finally used log level restriction and log
sink.

A `Rule` has the complete control over composing an appropriate logger.
The default `ConditionRule` just enables the specified log level,
if all rules match the actual log request.

For more complex conditions it is possible to compose conditions
using an `Or`, `And`, or `Not` condition.

Because `Rule` and `Condition` are interfaces, any desired behaviour
can be provided by dedicated rule and/or condition implementations.

## Configuration

It is possible to configue a log context from a textual configuration
using `config.Configure(ctx, bytedata)`:

```yaml
defaultLevel: Info
rules:
  - rule:
      level: Debug
      conditions:
        - realm: github.com/mandelsoft/spiff
  - rule:
      level: Trace
      conditions:
        - attribute:
            name: test
            value:
               value: testvalue  # value is the *value* type, here
```

Rules might provide a deserialization by registering a type object
with `config.RegisterRuleType(name, typ)`. The factory type must implement the
interface `scheme.RuleType` and provide a value object
deserializable by yaml.

In a similar way it is possible to register deserializations for
`Condition`s. The standard condition rule supports a condition deserialization
based on those registrations.

The standard names for rules are:
 - `rule`: condition rule

The standard names for conditions are:
- `and`: AND expression for a list of sub sequent conditions
- `or`: OR expression for a list of sub sequent conditions
- `not`: negate given expression
- `realm`: name for a realm condition
- `realmprefix`: name for a realm prefix condition
- `attribute`: attribute condition given by a map with `name` and `value`.
  
The config package also offers a value deserialization using
`config.RegisterValueType`. The default value type is `value`. 
It supports an `interface{}` deserialization.

For all deserialization types flat names are reserved for
the global usage by this library. Own types should use a reverse
DNS name to avoid conflicts by different users of this logging
API.

To provide own deserialization context, an own object of type
`config.Registry` can be created using `config.NewRegistry`.
The standard registry can be obtained by `config.DefaultRegistry()`
