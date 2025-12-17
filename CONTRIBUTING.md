# Contributing to TRIGRA

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them. Compelling examples of contributions include:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Table of Contents

- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
- [Styleguides](#styleguides)
  - [Commit Messages](#commit-messages)

## I Have a Question

> If you want to ask a question, we assume that you have read the available [Documentation](README.md).

Before you ask a question, it is best to search for existing [Issues](https://github.com/Taiwrash/trigra/issues) that might help you. In case you've found a suitable issue and still need clarification, you can write your question in this issue. It is also advisable to search the internet for answers first.

If you then still feel the need to ask a question and need clarification, we recommend the following:

- Open an [Issue](https://github.com/Taiwrash/trigra/issues/new).
- Provide as much context as you can about what you're running into.
- Provide project and platform versions (TRIGRA version, Kubernetes version, etc).

## I Want To Contribute

### Reporting Bugs

- **Use a clear and descriptive title** for the issue to identify the problem.
- **Describe the steps to reproduce the exact steps** which reproduce the issue.
- **Provide specific examples** to demonstrate the steps.
- **Describe the behavior you observed after following the steps** and point out what problem is with that behavior.
- **Explain which behavior you expected to see instead and why.**
- **Include screenshots and animated GIFs** which show you following the described steps and clearly demonstrate the problem. 

### Suggesting Enhancements

- **Use a clear and descriptive title** for the issue to identify the suggestion.
- **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
- **Provide specific examples to demonstrate the steps**.
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why.

### Your First Code Contribution

Unsure where to begin contributing? You can start by looking through these `good first issue` and `help wanted` issues:

- **Good first issues** - issues which should only require a few lines of code, and a test or two.
- **Help wanted issues** - issues which should be a bit more involved than `good first issue`.

Local development instructions:

1. Clone the repository
2. Run `go mod download`
3. Run `make build` or `go build ./...`
4. Ensure tests pass with `go test ./...`

## Styleguides

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

## License

By contributing to TRIGRA, you agree that your contributions will be licensed under its MIT License.
