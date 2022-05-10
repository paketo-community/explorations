List of all dependencies in the dep-server:

| Name              | Source                                                    | Compiled? And where to find  related code|
| ----------------- | --------------------------------------------------------- | ---------------------------------------- |
| bundler           | https://rubygems.org/downloads/                           | yes, binary-builder                      |
| composer          | https://getcomposer.org/download                          | no, buildpacks-ci                        |
| curl              | https://curl.se/download                                  | yes, buildpacks-ci                       |
| dotnet-aspnetcore | https://download.visualstudio.microsoft.com/download/pr   | processed, buildpacks-ci                 |
| dotnet-runtime    | https://download.visualstudio.microsoft.com/download/pr   | processed, buildpacks-ci                 |
| dotnet-sdk        | https://download.visualstudio.microsoft.com/download/pr   | processed, buildpacks-ci                 |
| go                | https://dl.google.com/go/                                 | yes, binary-builder                      |
| httpd             | http://archive.apache.org/dist/httpd                      | yes, binary-builder                      |
| icu               | https://github.com/unicode-org/icu/releases/download/     | yes, buildpacks-ci                       |
| nginx             | http://nginx.org/download/                                | yes, both locations?                     |
| node              | https://nodejs.org/dist/                                  | yes, binary-builder                      |
| php               | https://www.php.net/distributions/                        | yes, binary-builder                      |
| pip               | https://files.pythonhosted.org/packages                   | yes, with other pip deps buildpacks-ci   |
| pipenv            | https://files.pythonhosted.org/packages                   | yes, with other pip deps buildpacks-ci   |
| python            | https://www.python.org/ftp/python                         | yes, buildpacks-ci                       |
| poetry            | https://files.pythonhosted.org/packages                   | processed, buildpacks-ci                 |
| ruby              | https://cache.ruby-lang.org/pub/ruby                      | yes, binary-builder                      |
| rust              | https://static.rust-lang.org/dist/                        | yes, buildpacks-ci                       |
| tini              | https://github.com/krallin/tini/tarball                   | yes, buildpacks-ci                       |
| yarn              | https://github.com/yarnpkg/yarn/releases/download/v1.15.2 | processed, buildpacks-ci                 |
