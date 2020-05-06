# Git Log Issue Finder (GLIF)

This small program is used to extract a list of Jira issues from the output of a git log. This require an already checked-out
repository. 

Version: 1.4.0

## Content:
1. [Usage](#usage)
2. [diff-tags](#diff_tags_conf)
    1. [Basic Format](#Basic-format)
    1. [Examples](#Examples)
3. [Pipeline Configuration](#pipeline_configuration)
4. [Task Configuration](#task_configuration)
5. [Contact](#contact)

## <a name="usage" href="usage">Usage</a>
While this program was designed with the idea of integrating it to a *<a href="https://concourse-ci.org/" target="_blank">Concourse</a>* pipeline it can also be used a stand-alone
command-line tool.

Let's start by looking at the '--help' command:
```bash
$> glif --help
Usage of glif:
    -diff-tags string
          This parameter takes a specially formatted string: <FROM_TAG>==><TO_TAG>
                  The '==>' is litteral
                  The '<FROM_TAG>' and '<TO_TAG>' can have 2 formats:
                          1- A litteral tag name
                          2- A string with parameters (see below)
          Parameters are value declared between the '@(' and ')' litterals. Possible parameters:
                  LATEST: finds the latest value matching the string
                  LATEST-N: finds the Nth commit behind the latest value matching the string
          Examples:
          --diff-tags="v1.0.0-rc.@(LATEST-1)==>v1.0.0-rc.@(LATEST)"
    -directory string
          The directory of the git repo (default "./")
    -full-history
          Search the entire git log
    -since-latest-tag
          Search only from HEAD to the most recent tag
    -tickets string
          Comma-separated list of jira project keys

$> 
```
While the help command does not specify it, it's useful to note that every parameter can specified with one or two hyphen.

The following example assumes that you are working within the directory of your git repository.  
The most basic usage requires two (2) parameters; here's an example with the output:
```bash
$> glif --tickets="ABC,XYZ"
[ABC-001, ABC-007, XYZ-9246, ABC-045, ABC-0245, XYZ-007]
$> 
```
The above command is equivalent to this one:
```bash
$> glif --tickets="ABC,XYZ" --directory="./" --full-history
```
Now let's assume that a tag was made (1.0.0 for example). Following this tag a new feature (XYZ-999) was commited. If you run
the command but with the --since-latest-tag flags, here's the output you could expect:
```bash
$> glif --tickets="ABC,XYZ" --since-latest-tag
[XYZ-999]
$>
```
Running the command with --full-history will now give you the previous result with the added 'XYZ-999' feature.
```bash
$> glif --tickets="ABC,XYZ" --full-history
[ABC-001, ABC-007, XYZ-9246, ABC-045, ABC-0245, XYZ-007, XYZ-999]
$> 
```

If both '--full-history' and 'since-latest-tag' are specified then the '--full-history' is the one that'll take precedence.

## <a name="diff_tags_conf" href="diff_tags_conf">--diff-tags</a>
The 'diff-tags' parameters allows you to specify a very specific range of tags for glif to look into.
In the next subsection various example will be given so that developpers may have a starting point.

Each given 'diff-tags' value should then be used like this ```--diff-tags="<VALUE>"```.

### Basic format
The value given to 'diff-tags' will always be matched against the following regex:
```([a-zA-Z0-9.\-_@()/]+)==>([a-zA-Z0-9.\-_@()/]+)```.

Before the '==>' (arrow) if refered to as the 'from' and after the arrow is refered to as the 'to'. 
Therefore we have something that can be read as: FROM (tag) ==> TO (tag). 

### Examples
#### Basic example 1
This first example simply uses hard-coded values:
```
1.0.0-rc.1==>1.0.0-rc.2
```
This one will show all the matching tickets from the release candidate 1 to the release candidate 2
of version 1.0.0.

#### Basic example 2
Similarly to the first example, here we are searching between the first release candidate and the
actual release. 
```
1.0.0-rc.1==>1.0.0
```

#### Advanced example 1
Let's take the first example again but make it a little more generic. A genuine need that could arise is
to have the tickets that were worked on in the last release candidate. So the first example (the difference
between the first and second rc) covers that need when creating the tag 1.0.0-rc.2. 

But since we don't want to change the arguments of 'diff-tags' everytime we do a new release candidate, 
we can have this:

```
1.0.0-rc.@(LATEST-1)==>1.0.0-rc.@(LATEST)
```
This will lookup the tickets between the latest rc and the rc before that one (still for version 1.0.0).

#### Advanced example 2
```
@(SEMVER_LATEST)-rc.@(LATEST-1)==>@(SEMVER_LATEST)-rc.@(LATEST)
```

## <a name="pipeline_configuration" href="pipeline_configuration">Pipeline Configuration</a>

Now, here's an example of a Concourse job that uses git-log-issue-finder

```yml
jobs:
  - name: find-jira /* You can use whatever name you like */
    serial: true
    public: false
    plan: 
      - in_parallel:
        - get: <GIT_REPOSITORY_RESOURCE>
          /* Add 'passed' and/or 'trigger' configuration if needed */
      - task: git-log-issue-finder
        file: <PATH_TO_YML_TASK_CONFIGURATION>  
```

## <a name="task_configuration" href="task_configuration">Task Configuration</a>

To properly configure git-log-issue-finder, it should be done as a task and not directly as a pipeline resource. 

Here's what the task's yaml file should look like

```yml
platform: linux
image_resource:
  type: docker-image
  source:
    repository: turnscoffeeintoscripts/git-log-issue-finder
    tag: latest

params:
  TICKETS_FILTER: 'ABC,DEF,PROD,ETC'
  GIT_REPO_DIRECTORY: name-of-the-git-repo-directory
  ISSUES_DIRECTORY: name-of-the-directory-for-the-output-of-glif
  ISSUES_FILE: issues.txt
  FLAGS ''
  DIFF_TAGS: ''

inputs:
  - name: name-of-the-git-repo-directory

outputs:
  - name: name-of-the-directory-for-the-output-of-glif

run:
  path: /bin/sh
  args:
    - <PATH_TO_SHELL_SCRIPT>
```
The yaml configuration should contain three parameters ('params'). The destination (DESTINATION) parameter should contain
the name of the actual git repository folder. The tickets filter (TICKETS_FILTER) is a comma-separated list of Jira
project keys. Finally the filename (FILENAME) is the name of the file in which the result of the command will be written.
This last feature is useful if the result is needed as input of another job.

And now here's what the shell script should look like:

```bash
#!/bin/bash

set -e

# TODO
```

## <a name="contact" href="contact">Contact</a>
If you have any questions/comments please send them at: turns.coffee.into.scripts@gmail.com.

You may also submit pull-requests on github at: https://github.com/TurnsCoffeeIntoScripts/git-log-issue-finder 
to the branch 'master'.