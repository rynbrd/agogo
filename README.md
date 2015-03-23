A-Go-Go: A URL Shortener
========================
A-Go-Go is a simple URL shortener written in Go. A-Go-Go initializes its
database of shortened URLs from a file stored at a particular URL. It also
hosts an endpoint which allows the user to reload the database by sending a
request to it. 

Usage
-----
A-Go-Go takes its configuration from commandline flags. The following example
retrieves its database from a database hosted on example.net named links.db and
listens for requests on the default port:

    agogo -links=http://example.net/links.db

The `-links` flag is the only required flag. By default A-Go-Go listen to port
80 on all interfaces or `0.0.0.0:http`. This can be changed with `-listen`.

The reload endpoint is located at `/_reload/`. It will respond 200 to any sort
of HTTP request. Access to the endpoint can be retricted by providing the
`-allow` flag with a host or IP address. It my be specified more than once. The
reload endpoint may be disabled entirely with the `-no-reload` flag.

Links Format
------------
The links database is a text file containing a shortened link on each line of
the file. Each line consists of two parts: the local path hosting the
redirection and the URL to redirect to. They are separated by whitespace.
Everything after the first instance of whitespace is used as the `Location`
header in the redirection. Blank lines are and lines beginning with `#` are
ignored.

License
-------
Copyright (c) 2014 Ryan Bourgeois. Licensed under BSD-Modified. See the
[LICENSE][1] file for a copy of the license.
