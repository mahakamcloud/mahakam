# Mahakam CI
This directory contains configuration of pipelines and jobs executed for CI purposes. Mahakam uses [Concourse CI](https://concourse.ci).
Mahakam Concourse infrastructure is currently within private network.

## Current pipelines

Mahakam CI consists (as per today) of the following pipelines:
* `maha-pull-request` - executed for every Pull Request. runs basic syntax checks, unit tests and coverage.
* `maha-e2e-tests` - executed for open PRs with `run-e2e-tests` label.

New pipelines will be added as needed.

## Introduction to Concourse

Main way of interacting with Concourse is through its CLI, *fly*. Go to [Downloads](https://concourse.ci/downloads.html) page
to obtain the binary for your platform.

**Note:** Mahakam CI is currently within private network, and most of the commands below are accessible to the core team only.

First, you need to login to the server:
```
fly -t mahakam login -c http://10.120.1.20:8080
```

Then you can see list of pipelines and their jobs:

```
$ fly -t mahakam pipelines
name               paused  public
maha-pull-request  no      no    
maha-e2e-tests     no      no 
```

```
$ fly -t mahakam jobs -p maha-pull-request
name            paused  status     next
run-unit-tests  no      succeeded  n/a 
```

Pipelines in Concourse revolve around three main concepts: *tasks*, *jobs* and *resources* (read more about them [here](http://concourse.ci/concepts.html)).

Tasks are the smallest executable units of work. One of the most useful features of Concourse is ability to trigger a task from command line, including your local changes.