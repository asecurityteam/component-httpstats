# component-httpstats - Settings component for generating an HTTP stat client
[![GoDoc](https://godoc.org/github.com/asecurityteam/component-httpstats?status.svg)](https://godoc.org/github.com/asecurityteam/component-httpstats)

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=bugs)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=code_smells)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=coverage)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=ncloc)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=alert_status)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=security_rating)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=sqale_index)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=asecurityteam_component-httpstats&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=asecurityteam_component-httpstats)



<!-- TOC -->autoauto- [component-httpstats - Settings component for generating an HTTP stat client](#component-httpstats---settings-component-for-generating-an-http-stat-client)auto    - [Overview](#overview)auto    - [Quick Start](#quick-start)auto    - [Status](#status)auto    - [Contributing](#contributing)auto        - [Building And Testing](#building-and-testing)auto        - [License](#license)auto        - [Contributing Agreement](#contributing-agreement)autoauto<!-- /TOC -->

## Overview

This is a [`settings`](https://github.com/asecurityteam/settings) that
enables constructing an http client whose transport logs key HTTP
metrics with configurable names on every request. The resulting client
is powered by [`xstats`](https://github.com/rs/xstats) and
[`httpstats`](https://github.com/asecurityteam/httpstats).

## Quick Start

```golang
package main

import (
    "context"
    "net/http"
    "os"

    "github.com/asecurityteam/component-httpstats"
    "github.com/asecurityteam/transport"
    "github.com/asecurityteam/settings/v2"
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

## Status

This project is in incubation which means we are not yet operating this
tool in production and the interfaces are subject to change.

## Contributing

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

### License

This project is licensed under Apache 2.0. See LICENSE.txt for details.

### Contributing Agreement

Atlassian requires signing a contributor's agreement before we can accept a patch. If
you are an individual you can fill out the [individual
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d).
If you are contributing on behalf of your company then please fill out the [corporate
CLA](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b).
