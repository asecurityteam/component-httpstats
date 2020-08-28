<a id="markdown-component-httpstats---settings-component-for-generating-an-http-stat-client" name="component-httpstats---settings-component-for-generating-an-http-stat-client"></a>
# component-httpstats - Settings component for generating an HTTP stat client
[![GoDoc](https://godoc.org/github.com/asecurityteam/component-httpstats?status.svg)](https://godoc.org/github.com/asecurityteam/component-httpstats)
[![Build Status](https://travis-ci.com/asecurityteam/component-httpstats.png?branch=master)](https://travis-ci.com/asecurityteam/component-httpstats)
[![codecov.io](https://codecov.io/github/asecurityteam/component-httpstats/coverage.svg?branch=master)](https://codecov.io/github/asecurityteam/component-httpstats?branch=master)
<!-- TOC -->

- [component-httpstats - Settings component for generating an http stat client](#component-httpstats---settings-component-for-generating-an-http-stat-client)
    - [Overview](#overview)
    - [Quick Start](#quick-start)
    - [Status](#status)
    - [Contributing](#contributing)
        - [Building And Testing](#building-and-testing)
        - [License](#license)
        - [Contributing Agreement](#contributing-agreement)

<!-- /TOC -->

<a id="markdown-overview" name="overview"></a>
## Overview

This is a [`settings`](https://github.com/asecurityteam/settings) that
enables constructing an http client whose transport logs key HTTP
metrics with configurable names on every request. The resulting client
is powered by [`xstats`](https://github.com/rs/xstats) and
[`httpstats`](https://github.com/asecurityteam/httpstats).

<a id="markdown-quick-start" name="quick-start"></a>
## Quick Start

```golang
package main

import (
    "context"
    "net/http"
    "os"

    "github.com/asecurityteam/component-httpstats"
    "github.com/asecurityteam/transport"
    "github.com/asecurityteam/settings"
)

func main() {
    ctx := context.Background()
    os.SetEnv("METRICS_BACKEND", "testdependency")
    envSource := settings.NewEnvSource(os.Environ())

    client := &http.Client(
        Transport: httpstats.New(ctx, envSource)(
            transport.New(
              transport.OptionMaxResponseHeaderBytes(4096),
              transport.OptionDisableCompression(true)
            )
        )
    )
    req, _ := http.NewRequest(http.MethodGet, "www.google.com", http.NoBody)

    // various HTTP metrics emitted and tagged with a "client_dependency"
    // value of "testdependency", among other default tag keys and values
    _, _ := client.Do(req)
}
```

<a id="markdown-status" name="status"></a>
## Status

This project is in incubation which means we are not yet operating this
tool in production and the interfaces are subject to change.

<a id="markdown-contributing" name="contributing"></a>
## Contributing

<a id="markdown-building-and-testing" name="building-and-testing"></a>
### Building And Testing

We publish a docker image called
[SDCLI](https://github.com/asecurityteam/sdcli) that bundles all of our
build dependencies. It is used by the included Makefile to help make
building and testing a bit easier. The following actions are available
through the Makefile:

-   make dep

    Install the project dependencies into a vendor directory

-   make lint

    Run our static analysis suite

-   make test

    Run unit tests and generate a coverage artifact

-   make integration

    Run integration tests and generate a coverage artifact

-   make coverage

    Report the combined coverage for unit and integration tests

<a id="markdown-license" name="license"></a>
### License

This project is licensed under Apache 2.0. See LICENSE.txt for details.

<a id="markdown-contributing-agreement" name="contributing-agreement"></a>
### Contributing Agreement

Atlassian requires signing a contributor's agreement before we can accept a patch. If
you are an individual you can fill out the [individual
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d).
If you are contributing on behalf of your company then please fill out the [corporate
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b).
