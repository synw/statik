# Statik

A configurable watch/autoreload development server. Watch directories for changes and
run commands

## Install

Get a binary release for Linux or clone and compile with Go

```
wget https://github.com/synw/statik/releases/download/0.1.0/statik
chmod a+x statik
```

## Watchers

Create a `statik.config.json` file and set the watchers

Syntax: directory: [*list_of_commands*]

```javascript
{
  "watch": {
    "src": [
      "yarn buildsrc"
    ],
    "lib": [
      "yarn buildlib"
    ]
  }
}
```

To run the server:

```
./statik
```

## Autoreload

Add this to the index.html to enable the autoreload feature:

```html
<script type="text/javascript">
  const ip = "localhost"; // The ip of the server: ex: 192.168.1.3 or localhost
  (function () {
    var conn = new WebSocket("ws://" + ip + ":8042/ws");
    conn.onmessage = function (evt) {
      window.location.reload();
    }
  })();
</script>
```

## Server configuration

```javascript
{
  "root": "build/web",
  "port": 8090,
  "https": false,
  "watch": {}
}
```

**root**: the root directory to serve, relative path

**port**: the port to run the server on

**https**: run the server in https mode

## Flags

Optional flags are available, they take precedence over config params if set:

**-v**: verbose mode

**-nw**: disable the watchers

**-nr**: disable the autoreload

**-nc**: do not use the config file

**-https**: run the server in https mode

**-port**: the port to run the server on

**-root**: the root directory to serve, relative path

**-certs**: help about how to generate certificates for the https mode

Example:

```
./statik -https -nr -v -root=somedir/somesubdir
```
