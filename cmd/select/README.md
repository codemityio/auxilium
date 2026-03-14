# `select`

## Table of contents

- [`select`](#select)
  - [Table of contents](#table-of-contents)
  - [Summary](#summary)
  - [Usage](#usage)
  - [Example](#example)

## Summary

An interactive select command.

## Usage

Coming soon...

## Example

``` shell
list=$(cat <<-EOF
one one-value-one one-value-two
two two-value-one two-value-two
EOF
)
auxilium select --label="Choose one of the following" --list="${list}"
```
