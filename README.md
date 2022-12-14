# Logging for Go with context-specific Log Level Settings

This package provides a wrapper around the [logr](https://ggithub.com/go-logr/logr)
logging system supporting a rule based approach to enable log levels
for dedicated message contexts specified at the logging location.

The rule set is configured for a logging context:

```go
    ctx := logging.New(logrLogger)
```

Any `logr.Logger` can be passed here, the level for this logger
is used as base level for the `ErrorLevel` of loggers provided
by the logging context.
If the full control should be handed over to the logging context, 
the maximum log level should be used for the sink of this logger.

If the used base level should always be 0, the base logger has to 
be set with plain mode:

```go
    ctx.SetBaseLogger(logrLogger, true)
```

Now you can add rules controlling the accepted log levels for dedicated log 
locations. First, a default log level can be set:

```go
    ctx.SetDefaultLevel(logging.InfoLevel)
```

This level restriction is used, if no rule matches a dedicated log request.

Another way to achieve the same goal is to provide a generic level rule without any
condition:

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

A realm for the actual package can be defined as local variable by using the
`Package` function:

```go
var realm = logging.Package()
```

Instead of passing `Logger`s around, now the logging `Context` is used.
It provides a method to access a logger specific for a dedicated log
request, for example, for a dedicated realm.

```go
  ctx.Logger(realm).Info("my message")
```

The provided logger offers the level specific functions, `Error`, `Warn`, `Info`, `Debug` and `Trace`.
Depending on the rule set configured for the used logging context, the level
for the given message context decides, which message to pass to the log sink of
the initial `logr.Logger`.

Alternatively a traditional `logr.Logger` for the given message context can be
obtained by using the `V` method:

```go
  ctx.V(logging.InfoLevel, realm).Info("my message")
```

The sink for this logger is configured to accept messages according to the
log level determined by th rule set of the logging context for the given
message context.

*Remark*: Returned `logr.Logger`s are always using a sink with the base level 0,
which is potentially shifted to the level of the base `logr.Logger`
used to setup the context, when forwarding to the original sink. This means
they are always directly using the log levels 0..*n*.

If no rules are configured, the default logger of the context is used
independently of the  given arguments. The given message context information is
optionally passed to the provided logger, depending on the used 
message context type.

For example, the realm is added to the logger's name.

It is also possible to provide dedicated attributes for the rule matching
process:

```go
  ctx.Logger(realm, logging.NewAttribute("test", "value")).Info("my message")
```

Such an attribute can be used as rule condition, also. This way, logging
can be enabled, for dedicated argument values of a method/function.

Both sides, the rule conditions and the message context can be a list.
For the conditions, all specified conditions must be evaluated to true, to
enable the rule. A rule is evaluated against the complete message context of
the log requests.
The default `ConditionRule` evaluates the rules against the complete log
request and a condition is *true*, if it matches at least one argument.

The rules are evaluated in the reverse order of their definition.
The first matching rule defines the finally used log level restriction and log
sink.

A `Rule` has the complete control over composing an appropriate logger.
The default condition based rule just enables the specified log level,
if all conditions match the actual log request.

For more complex conditions it is possible to compose conditions
using an `Or`, `And`, or `Not` condition.

Because `Rule` and `Condition` are interfaces, any desired behaviour
can be provided by dedicated rule and/or condition implementations.

## Default Logging Environment

This logging library provides a default logging context, it can be obtained
by

```go
  ctx := logging.DefaultContext()
```

This way it can be configured, also. It can be used for logging requests
not related to a dedicated logging context.

There is shortcut to provide a logger for a message context based on
this default context:

```go
  logging.Log(messageContext).Debug(...)
```

or

```go
  logging.Log().V(logging.DebugLevel).Info(...
```

## Configuration

It is possible to configure a logging context from a textual configuration
using `config.ConfigureWithData(ctx, bytedata)`:

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

## Nesting Contexts

Logging contents can inherit from base contexts. This way the rule set,
logger and default level settings can be reused for a sub-level context.
THis contexts then provides a new scope to define additional rules
and settings only valid for this nested context. Settings done here are not
visible to log requests evaluated against the base context.

If a nested context defines an own base logger, the rules inherited from the base
context are evaluated against this logger if evaluated for a message
context passed to the nested context (extended-self principle).

A logging context reusing the settings provided by the default logging
context can be obtained by:

```go
  ctx := logging.NewWithBase(logging.DefaultContext())
```

## Preconfigured Rules, Message Contexts and Conditions

### Rules

The base library provides the following basic rule implementations.
It is possible to define own more complex rules by implementing
the `logging.Rule` interface.

- `NewRule(level, conditions...)` a simple rule setting a log level
for a message context matching all given conditions.

### Message Contexts and Conditions

The base library already provides some ready to use conditions
and message contexts:

- `Realm`(*string*) the location context of a logging request. This could
  be some kind of denotation for a functional are or Go package. To obtain the
  package realm for some coding the function `logging.Package()` can be used.

  Used as message context, the realm name is added to the logger name for
  the log request.

- `RealmPrefix`(*string*) (only as condition) matches against a complete 
  realm tree specified by a base realm.

- `Tag`(*string*) Just some tag for a log request.

  Used as message context, the tag name is not added to the logger name for
  the log request.

- `Attribute`(*string,interface{}*) the name of an arbitrary attribute with some
  value

  Used as message context, the key/value pair is added to the log message.

