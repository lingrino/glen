package cmd

// recurse determines if glen with also get variables from the project's parent groups
var recurse bool

// apiKey is the GitLab key that we should use when calling the API
var apiKey string

// directory is the path to the git repo that we should run glen on
var directory string

// remoteName is the name of the GitLab remote in your git repo
var remoteName string

// outputFormat is the text format that we should use to print our results to stdout
var outputFormat string
