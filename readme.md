# insightly

`insightly` is a command-line interface that analyzes website's accessibility by combining accessibility auditing tools (like [lighthouse](https://github.com/GoogleChrome/lighthouse), [pa11y](https://github.com/pa11y/pa11y)) with AI powered insights.

# installation

```bash
# if you don't have golang already installed, install it from https://go.dev
$ go install github.com/0xmukesh/insightly@latest
```

# usage

```bash
$ insightly COMMAND
Running command...
$ insightly (--version)
insightly version 0.0.4
$ insightly --help [COMMAND]
USAGE
  $ insightly COMMAND
...
```

# commands

- [`insightly setup`](#insightly-setup)
- [`insightly gen-ux`](#insightly-gen-ux)
- [`insightly config view`](#insightly-config-view)
- [`insightly config set`](#insightly-set)
- [`insightly config set-default`](#insightly-set-default)
- [`insightly config remove`](#insightly-remove)
- [`insightly help [COMMAND]`](#insightly-help-command)

## `insightly setup`

🔑 Setup your API keys for different LLMs and store it locally

```
USAGE
  $ insightly setup

DESCRIPTION
  Setup your API keys for different LLMs and store it locally

EXAMPLES
  $ insightly setup
```

## `insightly gen-ux`

🖨️ Generate UX reports for the given website

```
USAGE
  $ insightly gen-ux [website-url]

FLAGS:
  --llm string    Use any other LLM than your default LLM
  --save-report   Save parsed report in JSON format
  --use-ai        Use LLMs for generating a summary on how to improve the UX and accessiblity
  --use-pa11y     Use pa11y for running accessibility report

DESCRIPTION
  Generate UX reports

EXAMPLES
  $ insightly gen-ux https://example.com --use-pa11y --save-report --use-ai --llm=gemini
```

## `insighty config view`

⚙️ View configuration details

```
USAGE
  $ insightly config view

DESCRIPTION
  View configuration details

EXAMPLES
  $ insightly config view
```

## `insightly config set`

🎛️ Update configuration details

```
USAGE
  $ insightly config set

DESCRIPTION
  Update configuration details

EXAMPLES
  $ insightly config set
```

## `insightly config set-default`

🤖 Change your default LLM

```
USAGE
  $ insightly config set-default

DESCRIPTION
  Change your default LLM

EXAMPLES
  $ insightly config set-default
```

## `insightly config remove`

🗑️ Remove configuration details of a certain LLM

```
USAGE
  $ insightly config remove

DESCRIPTION
  Remove configuration details of a certain LLM

EXAMPLES
  $ insightly config remove
```
