#System Architecture Brainstorm

### Current system:
1. Dep-server has a list of `known-versions` for each dependency in GCP
2. Dep-server polls for new versions on a timer every hour for every
   dependency. Source URIs are hard coded into the dep-server code.
3. For each new dependency/version the build and test workflow is run
4. Build is delegated to cloudfoundry/buildpacks-ci and binary-builder and may involve pulling
   from source, compiling, or processing the dependency in some way.
5. A generic smoke test is run against the compiled dependency, and if passing
   it/s uploaded to an AWS S3 bucket
6. The metadata for the dependency is gathered in an upload workflow, versions
   are converted to semver and stored in an AWS S3 bucket
7. Endpoint routing is set up through AWS Route 53 abd Cloudfront, but the
   actual dep-server runs in Google App Engine
8. If failures occur or a dependency needs to be recompiled, you must be a
   Dependencies contributor at the minimum to re-run automation

### Idea 1:
0. Workflows and actions come from a generic place, like a github-config repo
   for dependencies. They are maintained by Dependencies maintainers and
   updated by GHA in the buildpacks similarly to how we keep workflows updated
   in the buildpacks now.
1. Known versions are tracked by the buildpack, instead of dep-server. Known
   versions are easily accessible without having to go to GCP.
2. Each buildpack polls for new versions on a timer, which runs every hour (?)
3. For each new dependency/version, the buildpack has a workflow to run and grab the newest version.
   - Source URI is passed into the workflow from the buildpack
   - When possible, the dependency comes directly from source
   - If processing of any kind is needed, it happens here. It should be easily
     runnable from a local workstation as well.
   - (?) Stacks for compatibility are passed in from the buildpack
4. Any compilation/processing code comes from the buildpack repository
5. Test is run against the dependency, whether it is compiled or not.
   Processed/compiled dependencies are uploaded to an S3 bucket, accessible via buildpack maintainers (?)
6. Metadata is gathered in step 5 alongside compilation and is uploaded to the same bucket (?) (I think we can do better than this)
7. TODO: dep-server app could be a lot clearer.
