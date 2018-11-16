# Sous Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/)
with respect to its command line interface and HTTP interface

## [Unreleased](//github.com/opentable/sous/compare/0.5.120...master)

### Fixed
* Error message when no manifest matches query on 'manifest get' and similar
  commands now lists the correct key/value pairs rather than jumbling them as
  before.
* Docker refs now always use a lower-case repo component. Previously it sometimes
  attempted to create a docker ref with upper-case chars in the repo component
  which is invalid and failed the build.

### Changed
* Client: sous jenkins cli revise format of generated Jenkinsfile
* Client: 'sous build' fails early when trying to re-build an existing version tag.
  Previously the build succeeded, and pushed a new docker image, but subsequent
  deploys did not use the new docker image, since sous uses the image digest, not
  the docker tag itself to identify images with SourceIDs.
* Client: the runmount build strategy is detected based on Dockerfile
  environment variables. Specifically, to be considered runmount, the Dockerfile
  or its parents must declare the SOUS_RUN_IMAGE_SPEC and BUILD_OUT environment variables.
  All known existing runmount containers already do this.

## [0.5.120](//github.com/opentable/sous/compare/0.5.117...0.5.120)
### Fixed
* Manifest validation now requires non-empty Owners field.
  NOTE: this will be a user experience hit, and needs to be highlighted to our users before deploy.
* Deployments to PROD clusters will use Jenkins agents with 'mesos-prod-sc'
  label (jenkins agents in PROD) to avoid QA -> PROD type of deployments.
* A race condition during whole-cluster resolutions meant the final status
  was sometimes inaccurately recorded. Real world implications of this are not
  completely clear, users are not expected to notice much difference.
* Client: sous metadata set would panic if no existing metadata was present in manifest prior to set
* Client: sous jenkins cli to use manifest metadata to generate JenkinsPipeline file
* Client: more accurate error message on 'sous manifest set' when ManifestID in YAML
  does not match that specified by flags and context.
* Client: 'artifact add' now records docker image refs using the configured default
  docker registry name. Previously it sometimes recorded local image refs which
  were not deployable.

### Changed
* Client: `SOUS_USE_SOUS_SEVER` env var must now be exactly uppercase `YES`
  in order to be considered "on". Previously any value, even empty string
  was considered "on". This change makes it easier to use in scripts.
* Client: complete overhaul of parsing singularity.json and singularity-request.json
  files. sous init with -use-otpl-deploy flag now fails unless it recognises all
  fields in the singularity.json and singularity-request.json files, and requires
  that there be a singularity-request.json file (where before it was optional).
  There is also a bunch more validation, too much to list here.

NOTE: Some features listed above as 0.5.120 were also released as 0.5.118 and 0.5.119.

## [0.5.117](//github.com/opentable/sous/compare/0.5.116...0.5.117)
### Changed
* Tests: Travis build times reduced by running sets of tests in parallel.
  This allows us to produce relese binaries, see note below in 0.5.116.


## [0.5.116](//github.com/opentable/sous/compare/0.5.115...0.5.116)
### Changed
* Client: 'sous artifact get' no longer requires -cluster flag.
* Client: 'sous artifact get' now prints artifact information (digest, type).
* Client: If `SOUS_BUILD_NOPULL=YES` then sous omits the `--pull` flag from
  Docker builds. This is mostly useful for tests. See also the `-dev` flag
  which may be more useful.

NOTE: No binaries were produced for 0.5.116 due to build timeout in Travis.

## [0.5.115](//github.com/opentable/sous/compare/0.5.114...0.5.115)
### Changed
* Client: Create a way for deploy to bypass cluster configuration for how http client
  obtains it's url for sous server.
* Server: On sous build, don't treat failure to update name cache as an error, just log.

## [0.5.114](//github.com/opentable/sous/compare/0.5.113...0.5.114)
### Changed
* Server: Default endpoint / implemented.

## [0.5.113](//github.com/opentable/sous/compare/0.5.112...0.5.113)
### Changed
* Both: Send the request URL back to the client and print to
  console on deploy.

## [0.5.112](//github.com/opentable/sous/compare/0.5.111...0.5.112)
### Changed
* Client: when receiving bad content type print the HTTP status code and
  text along with that error message.

## [0.5.111](//github.com/opentable/sous/compare/0.5.110...0.5.111)
### Fixed
* Client: builds using split container strategy were broken due to unique
  constraint violation when builder image was pushed using same SourceID as
  the runnable image. We no longer push the builder image to avoid this
  conflict. In future we may reinstate pushing builders if there is a use case
  for it.

## [0.5.110](//github.com/opentable/sous/compare/0.5.109...0.5.110)
### Fixed
* Both: Update how server errors get passed to client, if not json make more
  readable.

### Added
* Client: 'sous artifact add' command.
* Server: GET /artifact handler.

## [0.5.109](//github.com/opentable/sous/compare/0.5.104...0.5.109)
### Added
* Singularity request ID is now configurable per-deployment in the manifest.
  Changing SingularityRequestID results in the next deployment being done to
  the new SingularityRequestID, and the old request being orphaned, requiring
  manual cleanup.
* Client: New command to add artifact image to sous.

## [0.5.105](//github.com/opentable/sous/compare/0.5.102...0.5.105)
### Fixed
* Client: The client wasn't sending back the global state "defs" configuration,
  which leads to an obscure edge case that prevents the switch to Postgres
  based storage.
* Both: building with the `netcgo` build tag to ensure system DNS resolver is
  used instead of the Go native resolver.
* Server: caught a place where a potentially nil reference got a method called on it.

## [0.5.104](//github.com/opentable/sous/compare/0.5.102...0.5.104)
### Added
* Both: The client registers the digests for built artifacts to the server.
  This enhances speed and correctness of subsequent deploys.
* Developer: Adding a `sous-bootstrap` command (`go get github.com/opentable/sous/cmd/sous-bootstrap`)
  to handle initial and recovery deploys.

### Fixed
* Client: Builds are managed more reliably internally, using randomized intermediate Docker tags.
* Server: the startup failure status codes array was growing without bound.

## [0.5.102](//github.com/opentable/sous/compare/0.5.101...0.5.102)
### Added
* Server: Send user to singularity for Deployment
### Fixed
* Server: the database storage engine wasn't recording the advisory whitelists.

## [0.5.101](//github.com/opentable/sous/compare/0.5.100...0.5.101)
### Fixed
* CLI: 'sous manifest edit' was sometimes silently failing to apply changes on macOS,
  changes to how we re-read the temp file resolve this.

### Changed
* CLI: 'sous manifest edit' now uses a temp file with a .yaml extension so text
  editors are more likely to apply the right highlighting and auto formatting.
* CLI: Clearer -kind not recognised message.


## [0.5.100](//github.com/opentable/sous/compare/0.5.93...0.5.100)
### Added
* Server: state storage toggle behind a feature flag - servers can be
  configured to use the database as the source of truth.
* Manifest: per-deployment SingularityRequestID field. NOTE: This does not
  do anything yet, just gets round-tripped to Git and Postgres.

### Changed
* Client: Slack add additional channels to send via config
* Client: Now returns error/message when sous newdeploy is used to inform no longer valid.

## [0.5.93](//github.com/opentable/sous/compare/0.5.92...0.5.93)
### Changed
* Server: Postgres is now required for the server to operate.
* Both: SQLite is removed as a dependency of Sous.
* Both: Remove global logging.Log and refactor for the removal
* Server: Removed dependency on SQLite. PostgreSQL connection is now required by the server.
* Client: Refactor plumbing normalize gdm to action
* Client: sous deploy can now send a slack notification if you specify SlackHookURL and
  SlackChannel configuration in config.yaml or via environment variables
* All: when tags do not parse correctly, the error message is much clearer about why
  it was invalid, including invalid characters and their positions in the semver
  portion of the tag.

## [0.5.92](//github.com/opentable/sous/compare/0.5.91...0.5.92)
### Added
* Client: Deploy with zero instances will give a specific error message that you have zero instances
* Client: add -dev flag, at the moment just affects checking for local images
* Client: return a better message if manifestid wasn't found

## [0.5.91](//github.com/opentable/sous/compare/0.5.88...0.5.91)
### Added
* Client: Added a footer after command execution that if present, will display the request id that
  was passed in the header.
* Client: Runmount strategy, will cache maven builds using a volume mount at docker run time.

## [0.5.88](//github.com/opentable/sous/compare/0.5.87...0.5.88)
### Added
* Both: Sous commands that communicate with the server add a request ID header to trace the request.

### Changed
* Client: `sous init` now requires -kind flag which is either 'scheduled' or 'http-service'
* Client: error is returned if manifest set is sent a different source location
* All: don't log when status poller hasn't changed



## [0.5.86](//github.com/opentable/sous/compare/0.5.85...0.5.86)
### Added
* Client: `sous manifest edit` covers 90% of the manifest get/set use case, more easily.
* Server: more flexible, agile logging API.
* Server: logging API includes adding context fields to child loggers.
* All: updated to enable better running of tests in teamcity

### Changed
* Client: `sous deploy` and `sous newdeploy` now synonyms. Expect `sous newdeploy` to be removed.
* Client: Display more information in case timeout of sous newdeploy.  Also show Executor Message if failed deploy.
* Server: cleanups to logging output

## [0.5.85](//github.com/opentable/sous/compare/0.5.84...0.5.85)
### Added
* Server: Logging fields now pulled from generated logging.

## [0.5.84](//github.com/opentable/sous/compare/0.5.83...0.5.84)
### Fixed
* Server: badly formatted requests should no longer panic the HTTP client.

## [0.5.82](//github.com/opentable/sous/compare/0.5.81...0.5.82)

### Added
* Server: DB state operations are distributed over HTTP to their appropriate cluster.
* Client: add -force flag for newdeploy.  Defaults to false if not present, if true, no matter
  what GDM says, will submit to queue for deploy.

### Fixed
* Client: the flags for newdeploy were misleading. Updated flag parsing and error handling.

### Changed
* Server: Interactions with Singularity unified and simplified.
* All: Removed legacy Debug, Warn, Vomit logging infrastructure

## [0.5.78](//github.com/opentable/sous/compare/0.5.77...0.5.78)

### Added
* All: Add smoke test of the newdeploy cli
* All: Added docker run of posgtres and liquibase via task a go task executer

## [0.5.77](//github.com/opentable/sous/compare/0.5.76...0.5.77)

### Added
* Server: Single deployment rectification now waits for a reports a DeployState.
* Server: Single deployment name cache harvests - hopefully will further speed deployments.
* Server: pending state workaround for Singularity

### Fixed
* Server: /deploy-queue-item now returns correct queue position and Resolution field.
* Server: /single-deployment validation now correctly validates only after taking
  into consideration default values.
* Server: integration tests disregard failed deployments as blockers in certain cases
* Server: deployment operations are tracked by a UUID across API requests.

## [0.5.76](//github.com/opentable/sous/compare/0.5.72...0.5.76)

### Added
* CLI: sous newdeploy command: Much faster deployments. This will become
  the default 'sous deploy' after some real-world validation.
* Flaws describing problems with resource fields now get more context.
* Recording metrics for DB access (rows, time, errors)
* Server now accepts -autoresolver=false to disable to autoresolver.
* All: Added ability to add context to http requests
* Bugfixes for the newdeploy command

### Fixed
* Client: running `sous` from outside of git workspaces no longer results in
  a confusing Git error.
* Server: DB reading was omitting deployments without owners.

## [0.5.72](//github.com/opentable/sous/compare/0.5.71...0.5.72)

### Fixed
* PostgreSQL storage correctly retrieves arrays of healthcheck failure statuses
* Caught a race condition in the logging subsystem.

### Changed
* Log messages with `"@loglov3-otl": "sous-generic-v1"` field names changed:
  `fields` -> `sous-fields`, `types` -> `sous-types`, `ids` -> `sous-ids`,
  `id-values` -> `sous-id-values`.

### Added
* Server: PUT /single-deployment endpoint immediately adds a rectification
  to the queue and returns a link to monitor for completion.

## [0.5.71](//github.com/opentable/sous/compare/0.5.70...0.5.71)

### Added
* All: Updated cli, status_poller, client to start using structured logging

### Fixed
* Server: the PostgreSQL storage module successfully deduplicates proposed DB
  records generated from user input now.
* All: /all-deploy-queues returns correctly, with a somewhat different data format.

### Added
* Server: convert remainder of singularity package to generalmsg style logs

## [0.5.67](//github.com/opentable/sous/compare/0.5.66...0.5.67)

### Added
* All: Add truncated message that was supposed to go to logstash if it is found to error out from delivery to Kafka

## [0.5.66](//github.com/opentable/sous/compare/0.5.65...0.5.66)

### Added
* Server: Include connection string in db connection error log.
* Client: Add structured logging to status poller
* All: Create a new structured log that auto extracts IDs and stores in seperate fields for
  easier searching in logstash

### Changed
* All: Top-level global logger labeled "GLOBAL".
* All: Deployment builder, manifest get/set, StatusMiddleware, otplManifest now emits structured logs.
* CLI: When sous build fails due to no Dockerfile, error says exactly that.
* All: HTTP logging messages are at level "extra1" when successful

### Fixed
* All: Some formatted logs were incorrectly reporting missing values and were
  indiscriminately trying to render a single slice in the first format verb.
  Resolved so those logs format correctly.
* Server: postgres storage engine wasn't de-duplicating records as it parsed the GDM.

## [0.5.65](//github.com/opentable/sous/compare/0.5.63...0.5.65)

### Added
* Server: new endpoints /deploy-queues and /deploy-queue-item showing the
  list of all deployment queues and their lengths, and individual queue items
  respectively. /deploy-queue-item allows HTTP long-polling on the completion
  of a single rectification be providing the ?wait=true query parameter.
* Server: new endpoint /state/deployments allows cluster-specific updates to GDM.
* All: Default when testing, don't call recover when a log message fails to Deliver.
* CLI: Added timing information to report invocation message.
* All: Logging Reporter that allows allows semi flexible fields to be indexed.

## [0.5.63](//github.com/opentable/sous/compare/0.5.62...0.5.63)
### Added
* CLI: If no image is present in runspec, return a fatal flaw in build.
* Server: adding duplex state storage, to keep DB in sync until ready to switch over
* Server: Update logging to a more structured format: server, ext/singularity, resource, Generic Msg, handle_gdm, subpoller, volume, local_config, rest/client, disk_state_manager, router
* Server: A /health endpoint. For the time being, just a 200 and the running version of Sous.

### Changed
* All: error parsing repo from SourceLocation now more informative.
* Server: increased scoping of loggers - log entries should report their sources better.

### Fixed
* All: Data race in rectification queue.

## [0.5.62](//github.com/opentable/sous/compare/0.5.61...0.5.62)

### Added
* Server: If default log level is invalid, use Extreme level and print out a message.
* Logging: sous-diff-resolution logs include two new fields 'sous-diff-source-type' and
  'sous-diff-source-user', hard-coded to 'global rectifier' and 'unknown'.

### Fixed
* All: Setting a manifest with Owners field not alphabetically sorted no longer prevents updates
  to other manifests. (Owners are always stored in alphabetical order now.)
* Server: Kafka config validation now correctly checks if brokers were specified when Kafka is enabled.
* Server: Changing the Owners list now correctly causes the Singularity request to be updated.

## Internal
* Server: Storage module for Postgres uses placeholders correctly. As yet, not connected up.

## [0.5.61](//github.com/opentable/sous/compare/0.5.60...0.5.61)

### Added
* Server: Ability to load sous sibiling urls from environment variable

## [0.5.60](//github.com/opentable/sous/compare/0.5.59...0.5.60)

This release contains internal changes only; external behaviour is unmodified.
See [the set of commits in this release](//github.com/opentable/sous/compare/0.5.59...0.5.60)
for details of developer-facing changes.

## [0.5.59](//github.com/opentable/sous/compare/0.5.58...0.5.59)
### Fixed
* Server: Docker repository HTTP metrics collected and logged.
* Server: Sizes of *response* bodies logged and reported to Graphite.
* Server: Backtrace on HTTP logging excludes our libraries.

## [0.5.58](//github.com/opentable/sous/compare/0.5.57...0.5.58)

### Added
* Server: HTTP client requests to Singularity log to Kafka
* Server: HTTP client requests to the Docker registry log to Kafka
* Server: HTTP server responses log to Kafka
* Server: Additional metrics reported per request kind, per host and per kind,host pair.

## [0.5.57](//github.com/opentable/sous/compare/0.5.56...0.5.57)
### Fixed
* CLI + Server: No longer panics when printing a nil Deployment.

* Server: Profiling configuration - accepts SOUS_PROFILING environment variable correctly.

### Changed
* Developer: Server action extracted.

## [0.5.56](//github.com/opentable/sous/compare/0.5.55...0.5.56)
### Fixed
* CLI: No longer panics under normal operation (e.g. when trying to run 'sous deploy'
  outside of a git repo, or when server connection fails etc).
* CLI + Server: Logging goes to Kafka and Graphite when configured to again.

## [0.5.55](//github.com/opentable/sous/compare/0.5.54...0.5.55)
### Changed
* CLI: Quieter output for local operators. Previously many log messages were
  emitted to stderr which made CLI use difficult due to information overload.

### Fixed
* CLI: -s -q -v and -d (silent, quiet, verbose, and debug) flags now set the
  logging level appropriately. You'll now see a lot more output when using
  -d and -v flags, and only critical errors from -s, critical and console output
  from -q.

## [0.5.54](//github.com/opentable/sous/compare/0.5.53...0.5.54)
### Added
* All: a Schedule field on Manifests, which should publish to Singularity to allow for scheduled tasks.

### Fixed
* All: a bug in the Kafka subcomponent meant that the reverse sense was being applied to log entries.
  Furthermore, the severity was wrong, and the messages were being omitted.

## [0.5.53](//github.com/opentable/sous/compare/0.5.52...0.5.53)

### Fixed
* Server: the resolve complete message is properly emitted and produces stats

## [0.5.52](//github.com/opentable/sous/compare/0.5.51...0.5.52)

### Changed
* All: Logging to Kafka no longer goes through logrus.

### Added
* Client: console output related to deploys. Reports the number of instances
  intended, and reflects the successfully deployed source ID and target
  cluster.

## [0.5.51](//github.com/opentable/sous/compare/0.5.50...0.5.51)

### Added
* Server: max concurrent HTTP requests per Singularity server now configurable
  using `MaxHTTPConcurrencySingularity` in config file or
  `MAX_HTTP_CONCURRENCY_SINGULARITY` env var. Default is 10 was previously
  hard-coded to 10, so this is an opt-in change.

### Changed
* All: some logging behaviors - there may be more output than we'd like
* All: logging - capture the CLI output to Kafka, test that diff logs are generated.
* All: backtrace selection more accurate (i.e. the reported source of log entries is generally where they're actually reported)

## [0.5.50](//github.com/opentable/sous/compare/0.5.49...0.5.50)
### Fixed
* All: log entries weren't conforming to the requirements of our schemas.

## [0.5.49](//github.com/opentable/sous/compare/0.5.48...0.5.49)
### Fixed
* Server: at least one botched log message type has been caught and corrected.

### Added
* Client: more descriptive output from 'sous deploy' and 'sous update' commands.

## [0.5.48](//github.com/opentable/sous/compare/0.5.47...0.5.48)

### Fixed
* Erroneous output on stdout breaking some cli consumers.


## [0.5.47](//github.com/opentable/sous/compare/0.5.46...0.5.47)

### Fixed
Both: DI interaction causing panic on boot.

## [0.5.46](//github.com/opentable/sous/compare/0.5.44...0.5.46)

### Added
* Client: a separate status for "API requests are broken".
* Client: many new log messages during `sous update`, `sous deploy` and `sous plumbing status`.
* Client: general log message for start and end of command execution.

### Fixed
* Client: a panic would occur if the remote server wasn't available or responded with a 500.
* Attempting to fix invalid config using 'sous config' was not possible because we
  were validating config as a precondition to executing the command. We now only
  validate config for commands that rely on a valid config, so 'sous config' can be
  used to correct issues.

## [0.5.44](//github.com/opentable/sous/compare/0.5.43...0.5.44)

### Fixed

* Server: diff messages could panic when logging them if the diff didn't resolve correctly.
* All: logging panics would crash the app.
* Client: 'sous deploy' now waits for a complete resolution to take place before
  reporting failure. This avoids a race condition where earlier failures could
  be misreported as failures with the current deployment.
* Client: 'sous deploy' now bails out if no changes are detected after the present
  resolve cycle has completed, or if the latest version in the GDM does not match that
  expected. This solves an issue where deployments would appear to hang for a long time
  and eventually fail with a confusing error message, often due to conflicting updates.

## [0.5.43](//github.com/opentable/sous/compare/0.5.42...0.5.43)

### Added
* Server: deployment diffs are now logged as structured messages.
* Client: `sous init` command now has `-dryrun` flag, so you can generate a manifest
  without sending it to the server. This flag interacts with the `-flavor`,
  `-use-otpl-deploy` and `-ignore-otpl-deploy` flags as well, so you can check sous'
  intentions in all these scenarios without accidentally creating manifests you don't want.

### Fixed
* Server: Changing Startup.SkipCheck now correctly results in a re-deploy with the
  updated value.
* Client: commands 'sous deploy', 'sous manifest get' and 'sous manifest set' now receive the correct auto-detected offset
  so you no longer require the -offset flag in most cases (unless you need to override it).

## [0.5.42](//github.com/opentable/sous/compare/0.5.41...0.5.42)
### Fixed
* All: Graphite output was like `sous.sous.ci.sfautoresolver.fullcycle-duration.count`, now `sous.ci.sf.auto...`

## [0.5.41](//github.com/opentable/sous/compare/0.5.40...0.5.41)

### Fixed

* All: ot_* fields are actually populated - also instance-no
* All: metrics are scoped to env and region so that multiple sous instances don't
  clobber each other's metrics.
* All: component-id refers to the whole "sous" application; scoped loggers goes in logger-name

## [0.5.40](//github.com/opentable/sous/compare/0.5.39...0.5.40)

### Fixed

* Server: mismatches with logging schemas

## [0.5.39](//github.com/opentable/sous/compare/0.5.38...0.5.39)

### Fixed
* Server: logging configuration actually gets applied,
  so we can get graphite and kafka feeds.
* Server: Various places that log entries were tweaked to conform to ELK schema.

## [0.5.38](//github.com/opentable/sous/compare/0.5.37...0.5.38)

### Fixed
* All: restores -d and -v flags

## [0.5.37](//github.com/opentable/sous/compare/0.5.36...0.5.37)

### Fixed
* All: Cloned DI providers ("psyringe") were resulting in 2+ NameCaches, and
  uncontrolled access to the Docker registry cachce DB. A race condition led to
  errors that prevented deployment of Sous, and then blocked use of the CLI
  client. A stopgap was set up to force a NameCache to be provided early.

## [0.5.36](//github.com/opentable/sous/compare/0.5.35...0.5.36)

### Fixed
- Client: If only one otpl config found with no flavor, and a flavor is specified the found config was used.
- Client: If only one otpl config found for only one flavor, and no flavor or different flavor specified the found config was used.
- Client: `sous plumbing normalizegdm` broke the DI rules and added
  `DryrunNeither` an extra time, which led to a panic.
- Server: Initial database grooming had a race condition. Solved by ensuring NameCache is singular.

### Changed
- Client: Improve testability of default OT_ENV_FLAVOR logic and test.

## [0.5.35](//github.com/opentable/sous/compare/0.5.34...0.5.35)

### Fixed
- Client and server: various logging output is clearer and more correct.
- Client: Flavor flag wasn't being passed to otpl deploy logic.
- Client: Panic when setting OT_ENV_FLAVOR env variable if Env was unset.

## [0.5.34](//github.com/opentable/sous/compare/0.5.33...0.5.34)
### Added
- Client: When calling `sous init -flavor X` automatically add OT_ENV_FLAVOR value to the Env variables

### Fixed
- Client: `sous init -use-otpl-deploy` wasn't handling flavors properly.
- Client: `sous init -ignore-otpl-deploy` caused panic.
- Server: Logging wasn't being properly initialized and crashed on boot.

## [0.5.33](//github.com/opentable/sous/compare/0.5.32...0.5.33)
### Added
- Structured logging. Default logging will be structured (and colorful!)
  because of using logrus. Configuration should allow delivery of metrics to
  graphite and log entries to ELK.

### Fixed
- Server: Flavor metadata wasn't being pulled back into the ADS during rectify, which led to bad behaviors in deploy.

## [0.5.32](//github.com/opentable/sous/compare/0.5.31...0.5.32)
### Changed
- Developer: Refactors of filters and logging. Mostly to the good.

### Fixed
- Client: `sous build` for split containers was adding a path component,
  which broke the resulting deploy container.

## [0.5.31](//github.com/opentable/sous/compare/0.5.30...0.5.31)
### Fixed
- Client: `sous init` defaults resources correctly in the absence of other input.
## [0.5.30](//github.com/opentable/sous/compare/0.5.29...0.5.30)
### Fixed
- Client: certain commands were missing DI for a particular value. These are
  fixed, and tests added for the omission.
- Client: new values for optional fields no longer elided on `manifest set` -
  if the value was missing, the new value would be silently dropped.

## [0.5.29](//github.com/opentable/sous/compare/0.5.28...0.5.29)
### Changed
- Developer: DI system no longer used on a per-request basis.

## [0.5.28](//github.com/opentable/sous/compare/0.5.27...0.5.28)
### Fixed
- Client & Server: rework of DI to contain scope of variable assignment,
  and retain scope from CLI invocation to server.
- Client: multiple target builds weren't getting their offsets recorded correctly.
- Developer: dependency injection provider now an injected dependency of the CLI object.

## [0.5.27](//github.com/opentable/sous/compare/0.5.26...0.5.27)
### Fixed
- Client: omitting query params on update.
- Client: Failed to merge maps to null values (which came from original JSON).

## [0.5.26](//github.com/opentable/sous/compare/0.5.25...0.5.26)
### Fixed
- Client: using wrong name for the /manifest endpoint

## [0.5.25](//github.com/opentable/sous/compare/0.5.24...0.5.25)
### Fixed
- Client: simple dockerfile builds were crashing

## [0.5.24](//github.com/opentable/sous/compare/0.5.23...0.5.24)
### Fixed
- Client: the `sous manifest get` and `set` commands correctly accept a `-flavor` switch.

## [0.5.23](//github.com/opentable/sous/compare/0.5.22...0.5.23)
### Added:
- Client: Split buildpacks can now provide a list of targets,
  and produce all their build products in one `sous build`.
### Changed:
- Client: Client commands now have a "local server" available if no server is
  configurated. This is the start of the path to using HTTP client/server
  interactions for everything, as opposed to strategizing storage modules.

## [0.5.22](//github.com/opentable/sous/compare/0.5.21...0.5.22)
### Fixed
- Server: Validation checks didn't consider default values.

## [0.5.21](//github.com/opentable/sous/compare/0.5.20...0.5.21)
### Added
- Server: Startup field "SkipCheck" added - rather than omitting a
  healthcheck URI, services must set this field `true` in order to signal that
  they don't make use of a "ready" endpoint.
- Server: Full set of Singularity 0.15 HealthcheckOptions now supported.
  Previous field names on Startup retained for compatibility,
  although there may be nuanced changes in how they're interpreted
  (because Singularity no longer supports the old version.)
- Server: Logging extended to collect metrics.
- Server: Metrics exposed on an HTTP endpoint.
- Developer: LogSet extracted to new util/logging package, some refiguring of
  the types and interfaces there with an eye to pulling in a structured logger.
- Client: Command `sous plunbing normalizegdm` is a utility to round-trip the
  Sous storage format. With luck, running this command after manual changes to
  the GDM repo will correct false conflicts.
### Changed
- Server: Startup DeployConfig part fields now have cluster-based default
  values. (These should be configued in the GDM before this version is
  deployed!)
  Most notably, the CheckReadyProtocol needs a default ("HTTP" or "HTTPS")
  because Singularity validates the value on its side.

## [0.5.20](//github.com/opentable/sous/compare/0.5.19...0.5.20)
### Fixed
- Client: `sous deploy` wasn't recognizing its version if there was a prefix supplied.

## [0.5.19](//github.com/opentable/sous/compare/0.5.18...0.5.19)

Only changes in this release are related to deployment.

### Fixed
- Developer: deployment key changed to a machine account.
## [0.5.18](//github.com/opentable/sous/compare/0.5.16...0.5.18)
### Added
- Client: `sous query clusters` will enumerate all the logical clusters sous currently handles, for ease of manifest editing and deployment.

### Changed
- Developer: Makefile directives to build a Sous .deb and upload it to an Artifactory server.
- Server: /gdm now sets Etag header based on current revision of the GDM. Clients attempting
  to update with a stale Etag will have update rejected appropriately.
- Server: Updated Singularity API consumer to work with Singularity v0.15.

### Fixed
- Client: Increased buffer size in shell to handle super-long tokens.
- Client: Conflicting updates to the same manifest now fail appropriately (e.g. during 'sous deploy').
  This should mean better performance when running 'sous deploy' concurrently on the same manifest,
  as conflicting updates won't require multiple passes to resolve.


## [0.5.17]
__(extra release because of mechanical difficulties)__

## [0.5.16](//github.com/opentable/sous/compare/0.5.15...0.5.16)
- Fixed prefixed versions for update and deploy. At this point, git "package" tags should work everywhere they're used.

## [0.5.15](//github.com/opentable/sous/compare/0.5.14...0.5.15)
### Changed
- Panics during rectify print the stack trace along with the error message in the logs.
  Previously the stack trace was printed earlier in the log, making correlation
  difficult.
- Server: All read/write access to the GDM now serialised.
  This is to partially address and issue where concurrent calls to 'sous deploy'
  could result in one of them finishing in the ResolveNotIntended state.

### Fixed
- Client: 'sous build' was failing when using a semver tag with a non-numeric prefix.
  Validation logic is now shared, so 'sous build' succeeds with these tags.
- Non-destructive updates: clients won't clobber fields in the API they don't recognize. Result should be more stable, less coupled client-server relationship.

## [0.5.14](//github.com/opentable/sous/compare/0.5.13...0.5.14)

### Fixed
- A change to the Singularity API was breaking JSON unmarshaling. We now handle those errors as a "malformed" request - i.e. not Sous's to manage.

## [0.5.13](//github.com/opentable/sous/compare/0.5.12...0.5.13)

### Added
- Git tags with a non-numeric prefix and a semver suffix (e.g. 'release-1.2.3' or 'v2.3.4')
  are now considered a "semver" tag, and Sous will extract the version from them.
- Static analysis of important core data model calculations to ensure that all the components of those structures are at least "touched" during diff calculation.
- For developers only, there are 2 new build targets: `install-dev` and
  `install-brew`. These allow developers on a Mac to quickly switch between having
  a personal dev build, or the latest release from homebrew installed locally.

### Fixed
- Operations that change more than one manifest will now be rejected with an
  error. We do not believe there are any such legitimate operations, and
  there's a storage anomoly that surfaces as multiple manifests changing at
  once which we hope this will correct.
- 'sous manifest get' wrongly returned YAML with all lower-cased field names.
  Now it correctly returns YAML with upper camel-cased field names.
  Note that this does not apply to map keys, only struct fields.
- Deployment processing wasn't properly waited on, which could cause problems.

## [0.5.12](//github.com/opentable/sous/compare/0.5.11...0.5.12)
### Fixed
- Issue where deployments constantly re-deployed due to spurious Startup.Timeout diff.

## [0.5.11](//github.com/opentable/sous/compare/0.5.10...0.5.11)
### Fixed
- Singularity now accepts changes to Startup options in manifest.
- Off-by-one error with Singularity deploy IDs, fixed in 0.5.9, re-introduced in
  0.5.10. Now includes better tests surrounding edge cases.

## [0.5.10](//github.com/opentable/sous/compare/0.5.9...0.5.10)
### Fixed
- Off-by-one error with long request IDs.
- Startup information not recovered from Singularity, so not compared for deployment.

## [0.5.9](//github.com/opentable/sous/compare/0.5.8...0.5.9)
### Fixed
- Long version strings resulted in Singularity deploy IDs longer than the max
  allowed length of 49 characters. Now they are always limited to 49.

## [0.5.8](//github.com/opentable/sous/compare/0.5.7...0.5.8)
### Fixed
- Now builds and runs on Go 1.8 (one small change to URL parsing broke Sous for go 1.8).
- New Startup configuration section in manifests now correctly round-trips via 'sous
  manifest get|set' and takes part in manifest diffs.

## [0.5.7](//github.com/opentable/sous/compare/0.5.6...0.5.7)
### Changed
- Images built with Sous get a pinning tag that now includes the timestamp of
  the build, so that multiple builds on a single revision won't clobber labesls
  and make images inaccessible.

## [0.5.6](//github.com/opentable/sous/compare/0.5.5...0.5.6)
### Fixed
- Sous server was unintentionally filtering out manifests with non-empty offsets or flavors.

## [0.5.5](//github.com/opentable/sous/compare/0.5.4...0.5.5)

### Fixed
- Resolution cycles allocate much less memory, which hopefully keeps the memory headroom of Sous much smaller over time.

## [0.5.4](//github.com/opentable/sous/compare/0.5.3...0.5.4)

### Added

- Sous server now returns CORS headers so that the Sous SPA (forthcoming) can consume its data.

### Fixed

- Crashing bug on GDM updates.

## [0.5.3](//github.com/opentable/sous/compare/0.5.2...0.5.3)

### Added
- Profiling endpoints, gated with a `server` flag, or the SOUS_PROFILING env variable.

### Fixed
- Environment variable defaults from cluster definitions
  no longer elide identical variables on manifests,
  which means that common values can be added to the defaults
  without undue concern for manifest environment variables.

## [0.5.2](//github.com/opentable/sous/compare/0.5.1...0.5.2)

### Added
- Extra debug logging about how build strategies are selected.
- Startup options in manifest to set healthcheck timeout and relative
  URI path of healthcheck endpoint.

### Changed
- Singularity RequestIDs are generated with a suffix of the MD5 sum of
  pre-slug data instead of a random UUID.
- Singularity RequestIDs are shortened to no longer include FQDN or
  organization of Git repo URL.

### Fixed
- Calls to `docker build` now have a `--pull` flag so that stale cached FROM
  images don't confuse builds.

## [0.5.1](//github.com/opentable/sous/compare/0.5.0...0.5.1)

### Fixed
- Singularity RequestIDs retrieved from Singularity are reused when updating deploys,
  instead of recomputing fresh unique ones each time.

### Minor
- Added a tool called "danger" to do review of PRs.

## [0.5.0](//github.com/opentable/sous/compare/0.4.1...0.5.0)

### Added
* Split image build strategy: support for using a build image to produce artifacts to be run
  in a separate deploy image.

### Changed
* Sous detects the tasks in its purview based on metadata it sets when the task
  is created, rather than inspecting request or deploy ids.

### Fixed
* Consequent to detecting tasks based on metadata,
  Sous's requests are now compatible
  with Singularity 0.14,
  and the resulting Mesos Task IDs are suitable to use as Kafka client ids.

## [0.4.1](//github.com/opentable/sous/compare/0.4.0...0.4.1)

### Fixed
- Status for updated deploys was being reported as if they were already stable.
  The stable vs. live statuses reported by the server each now have their own
  GDM snapshot so that this determination can be made properly.

## [0.4.0](//github.com/opentable/sous/compare/0.3.0...0.4.0)

### Added
- Conflicting GDM updates now retry, up to the number of deployments in their manifest.

### Changed
- Failed deploys to Singularity are now retried until they succeed or the GDM
  changes.

## [0.3.0](//github.com/opentable/sous/compare/0.2.1...0.3.0)

### Added
- Extra metadata tagged onto the Singularity deploys.
- `sous server` now treats its configured Docker registry as canonical, so
  that, e.g. regional mirrors can be used to reduce deploy latency.

### Changed

- Digested Docker image names no longer query the registry, which should reduce
  our requests count there.

## [0.2.1](//github.com/opentable/sous/compare/0.2.0...0.2.1)

### Added

- Adds Sous related-metadata to Singularity deploys for tracking and visibility purposes.

### Fixed

- In certain conditions, Sous would report a deploy as failed before it had completed.

## [0.2.0](//github.com/opentable/sous/compare/0.1.9...0.2.0) - 2017-03-06

### Added

- 'sous deploy' now returns a nonzero exit code when tasks for a deploy fail to start
  in Singularity. This makes it more suitable for unattended contexts like CD.

### Fixed

- Source locations with empty offsets and flavors no longer confuse 'sous plumbing status'.
  Previously 'sous plumbing status' (and 'sous deploy' which depends on it) were
  failing because they matched too many deployments when treating the empty
  offset as a wildcard. They now correctly treat it as a specific value.
- 'sous build' and 'sous deploy' previously appeared to hang when running long internal
  operations, they now inform the user there will be a wait.


## [0.1.9](//github.com/opentable/sous/compare/0.1.8...0.1.9) - 2017-02-16

### Added

- 'sous init -use-otpl-deploy' now supports flavors
  defined by otpl config directories in the `<cluster>.<flavor>` format.
  If there is a single flavor defined, behaviour is like before.
  Otherwise you may supply a -flavor flag to import configurations of a particular flavor.
- config.yaml now supports a new `User` field containing `Name` and `Email`.
  If set this info is sent to the server alongside all requests,
  and is used when committing state changes (as the --author flag).
- On first run (when there is no config file),
  and when a terminal is attached
  (meaning it's likely that a user is present),
  the user is prompted to provide their name and email address,
  and the URL of their local Sous server.
- `sous deploy`
  (and `sous plumbing status`)
  now await Singularity marking the indended deployment as active before returning.

### Fixed
- Deployment filters (which are used extensively) now treat "" dirs and flavors
  as real values, rather than wildcards.

## [0.1.8](//github.com/opentable/sous/compare/v0.1.7...0.1.8) - 2017-01-17

### Added
- 'sous init' -use-otpl-config now imports owners from singularity-request.json
- 'sous update' no longer requires -tag or -repo flags
  if they can be sniffed from the git repo context.
- Docs: Installation document added at doc/install.md

### Changed

- Logging, Server: Warn when artifacts are not resolvable.
- Logging: suppress full deployment diffs in debug (-d) mode,
  only print them in verbose -v mode.
- Sous Version outputs lines less than 80 characters long.

### Fixed

- Internal error caused by reading malformed YAML manifests resolved.
- SourceLocations now return more sensible parse errors when unmarshaling from JSON.
- Resolve errors now marshalled correctly by server.
- Server /status endpoint now returns latest status from AutoResolver rather than status at boot.

## [0.1.7](//github.com/opentable/sous/compare/v0.1.6...v0.1.7) 2017-01-19

### Added

- We are now able to easily release pre-built Mac binaries.
- Documentation about Sous' intended use for driving builds.
- Change in 'sous plumbing status' to support manifests that deploy to a subset of clusters.
- 'sous deploy' now waits by default until a deploy is complete.
  This makes it much more useful in unattended CI contexts.

### Changed

- Tweaks to Makefile and build process in general.

## [0.1.6](//github.com/opentable/sous/compare/v0.1.5...v0.1.6) 2017-01-19

Not documented.

## [0.1.5](//github.com/opentable/sous/compare/v0.1.4...v0.1.5) 2017-01-19

Not documented.

## [0.1.4](//github.com/opentable/sous/compare/v0.1.3...v0.1.4) 2017-01-19

Not documented.

## Sous 0.1

Sous 0.1 adds:
- A number of new features to the CLI.
  - `sous deploy` command (alpha)
  - `-flavor` flag, and support for flavors (see below)
  - `-source` flag which can be used instead of `-repo` and `-offset`
- Automatic migrations for the Docker image name cache.
- Consistent identifier parse and print round-tripping.
- Updates to various pieces of documentation.
- Nicer Singularity request names.

## Consistency

- Changes to the schema of the local Docker image name cache database no longer require user
  intervention and re-building the entire cache from source. Now, we track the schema, and
  migrate your cache as necessary.
- SourceIDs and SourceLocations now correctly round-trip from parse to string and back again.
- SourceLocations now have a single parse method, so we always get the same results and/or errors.
- We somehow abbreviated "Actual Deployment Set" "ADC" ?! That's fixed, we're now using "ADS".

### CLI

#### New command `sous deploy` (alpha)

Is intended to provide a hook for deploying single applications from a CI context.
Right now, this command works with a local GDM, modifying it, and running rectification
locally.
Over time, we will migrate how this command works whilst maintaining its interface and
semantics.
(The intention is that eventually 'sous deploy' will work by making an API call and
allowing the server to handle updating the GDM and rectifying.)

#### New flag `-flavor`

Actually, this is more than a flag, it affects the underlying data model, and the way
we think about how deployments are grouped.

Previously, Sous enabled at most a single deployment configuration per SourceLocation
per cluster. This model covers 90% of our use cases at OpenTable, but there are
exceptions.

We added "flavor" as a core concept in Sous, which allows multiple different deployment
configurations to be defined for a single codebase (SourceLocation) in each cluster. We
don't expect this feature to be used very much, but in those cases where configuration
needs to be more granular than per cluster, you now have that option available.

All commands that accept the `-source` flag (and/or the `-repo` and `-offset` flags) now
also accept the `-flavor` flag. Flavor is intended to be a short string made of
alphanumerics and possibly a hyphen, although we are not yet validating this string.
Each `-flavor` of a given `-source` is treated as a separate application, and has its
own manifest, allowing that application to be configured globally by a single manifest,
just like any other.

To create a new flavored application, you need to `sous init` with a `-flavor` flag. E.g.:

    sous init -flavor orange

From inside a repository would initiate a flavored manifest for that repo. Further calls
to `sous update`, `sous deploy`, etc, need to also specify the flavor `orange` to
work with that manifest. You can add an arbitrary number of flavors per SourceLocation,
and using a flavored manifest does not preclude you from also using a standard manifest
with no explicit flavor.

#### New flag `-source`

The `-source` flag is intended to be a replacement for the `-repo` and
`-offset` combination, for use in development environments. Note that we do not have
any plans to remove `-repo` and `-offset` since they may still be useful, especially
in scripting environments.

Source allows you to specify your SourceLocation in a single flag.
Source also performs additional validation,
ensuring that the source you pass can be handled by Sous end-to-end.
At present, that
means the repository must be a GitHub-based, in the form:

    github.com/<user>/<repo>

If your source code is not based in the root of the repository, you can add the offset
by separating it with a comma, e.g.:

    github.com/<user>/<repo>,<offset>

Because GitHub repository paths have a fixed format that Sous understands, you can
optionally use a slash instead of a comma, so the following is equivalent:

    github.com/<user>/<repo>/<offset>

(and offset can itself contain slashes if necessary, just like before).

### Documentation

We have made various documentation improvements, but there are definitely some that are
still out of date, which we will look to resolve in the coming weeks. Improvements made
in this release include:

- Better description of networking setup for Singularity deployments.
- Update to the deployment workflow documentation.
- Some fixes to the getting started document.

### Singularity request names

Up until now, Singularity request names looked something like this:

    github.comopentablereponameclustername

Which is not a great user experience, and has a large chance of causing naming collisions.
This version of Sous changes these names to use the form:

    <SourceLocation>:<Flavor>:<ClusterName>

E.g. for a simple repo with no offset or flavor, it looks like this:

    github.com>opentable>sous::cluster-name

With an offset and flavor, it expands to something like this:

    github.com>opentable>sous,offset:flavor-name:cluster-name


### Other

There have been numerous other small tweaks and fixes, mostly at code level to make our
own lives easier. We are also conscientiously working on improving test coverage, and this
cycle hit 54%, we expect to see that rise quickly now that we fail CI when it falls. You
can track test coverage at https://codecov.io/gh/opentable/sous.

For more gory detail, check out the [full list of commits between 0.0.1 and 0.1](https://github.com/opentable/sous/compare/v0.0.1...v0.1).
