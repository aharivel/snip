# snip

snip is a category-based CLI for storing and retrieving Markdown snippets.

## Install

Build locally:

```bash
make build
```

Install to your PATH (default: `/usr/local/bin`):

```bash
sudo make install
```

Install completions:

```bash
sudo make install-completions
```

## Usage

```bash
snip list
snip create ocp
snip ocp list
snip ocp show "Login to cluster"
snip ocp clip "Login to cluster"
snip ocp find "kubeconfig"
snip ocp edit
```

## Storage Format

snip stores category files under the data directory:

- Default: `~/.snip/categories/`
- Override with: `SNIP_DATA_DIR=/path/to/data`

Each category is a Markdown file named `<category>.md`. Entries are defined with
headings using the `## ` prefix. The first fenced code block in an entry is the
snippet that `snip clip` copies to the clipboard.

### Example Category File

This example doubles as a template for tests/CI.

```markdown
## Login to cluster
Use this when you need a fresh token.

~~~bash
oc login https://api.example:6443 --token=$TOKEN
~~~

## Get project list

~~~bash
oc get projects
~~~

## Fix stuck namespace
Delete finalizers and reapply.

~~~bash
oc get ns stuck -o json | jq '.spec.finalizers = []' | oc replace -f -
~~~
```

## Completion

Use the Makefile targets for completions, or generate them manually:

```bash
make completions
snip completion bash > /etc/bash_completion.d/snip
snip completion zsh > ~/.zsh/completions/_snip
snip completion fish > ~/.config/fish/completions/snip.fish
```

### Completion Hints

Completions include short descriptions so your shell can visually separate
categories, actions, and headlines. For example, Zsh shows these labels with
different colors by default.

## Editor

`snip <category> edit` opens the category file in `$EDITOR` (defaults to `vim`).

## Clipboard

snip attempts to use:

- macOS: `pbcopy`
- Linux: `wl-copy` or `xclip`

If no clipboard tool is available, `snip clip` prints the snippet to stdout.
