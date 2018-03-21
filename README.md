# beanstalk-tools

Tools for moving beanstalk data about

## Installation

```
go get github.com/dxw/beanstalk-tools/cmd/...
```

## Usage

To consume all items in a tube and store them in `export.json`:

```
% bst-export localhost:11300 my-tube > export.json
```

To read all items from `export.json` and put them into the tube:

```
% bst-import localhost:11300 my-tube < export.json
```

The output/input format is a stream of JSON objects.

## Licence

MIT
