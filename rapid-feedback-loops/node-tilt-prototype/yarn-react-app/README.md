# React App w/ Yarn

## Getting Started

### Prerequisites
1. `tilt` CLI + k8s cluster
1. `pack` CLI

Run `tilt up` from this directory to stand up the deployment.

## Hot Reloading

Make code changes in any file within the `src` directory. A modification has
been set up for you in `src/App.js` already. 

Try adding the `react-icons` package to the project using Yarn. Run `yarn add
react-icons --modules-folder <path/to/dir/outside/of/app/dir>`. You should see
the `yarn install` process begin within the Tilt dashboard. Once the process is
complete, first uncomment the code in `src/question.js` then uncomment the
`Question` component in index.js. Finally, refresh the webpage and scroll to
the bottom of the page to see your new component!

## Limitations 

Adding packages to Node.js projects is a common part of the development
workflow. This currently does not work well with buildpacks + Tilt; running `yarn add <package>`
sometimes results in the application crashing and a full rebuild.

Ideally, upon adding a package to the project via `yarn add <package>`, a user
would need only sync their package.json and yarn.lock, without also syncing
their entire node_modules folder. A run step could then be configured to run `yarn install`
from within the container's `/workspace` directory upon modifications to package.json and/or yarn.lock.

Currently, implementing this via the Tiltfile presents some challenges:

- In order to prevent syncing the entire node_modules directory when adding
  packages (which seems to cause application failures on Tilt), care must be
  taken to place the local node_modules outside of the app directory by running `yarn
  add <package> --modules-folder <path/to/external/dir>`.

- To reinstall node_modules after syncing package.json and yarn.lock, it is
  insufficient to simply call `yarn install` inside a `run()` step since Tilt
  presumably runs `docker exec` with the command specified. This means the
  command is executed in the container outside of the launch environment
  created by the lifecycle. To remedy this, we must `run()` any commands in the
  format `/cnb/lifecycle/launcher <cmd> <args>` as the lifecycle would.

- Given the above, it is possible to run arbitrary commands in the container to
  enable dependency management workflows with Yarn. However, since watchexec is
  set up to watch the workspace, it restarts its given process on changes to
  package.json and/or yarn.lock. This creates a sort of race condition between
  the restarting watchexec process and, for example, a `yarn install` command
  which is regenerating node_modules often resulting in the application
  crashing and a full rebuild. This can be resolved by configuring watchexec to
  ignore package management artifacts (i.e. package.json, yarn.lock,
  node_modules, etc.). Tilt may still run `yarn install` on changes to
  package.json/yarn.lock, but the application process will not be restarted
  until code changes are made, avoiding any potential clashes.

## Recommendations

- When BP_LIVE_RELOAD_ENABLED is set to true, the `yarn-start` buildpack should
  set a start command which invokes whatever command it crafts from
  package.json or elsewhere via watchexec.

- The watchexec start command should be configured such that it restarts the
  given process whenever files in the workspace are changed.

- The watchexec start command should be configured such that it ignores changes
  to package.json, yarn.lock & `node_modules/` so that Yarn processes which
  modify dependency management files or web frameworks which modify the
  node_modules cache on-the-fly (e.g. React, Vue) do not trigger a restart. 


