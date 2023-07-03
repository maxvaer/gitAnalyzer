<h1 align="center">
  gitAnalyzer
</h1>

<h4 align="center">A template based scanner for GitHub repositories.</h4>

## Features
* The number of repositories gitAnalyzer scans at the same time, can be set.
* Execute regular expression, console command, Bash or Python scripts.
* A crawler to fetch URLs and metadata of all public repositories.
* A Web-UI to monitor the current scan.

## Installation
gitAnalyzer requires [**Go**](https://go.dev/doc/install) for installation.  
Run the following command to install the tool:
```console
git clone https://github.com/maxvaer/gitAnalyzer.git && cd gitAnalyzer && go mod tidy && go build
```

## Usage & Templates
Take a look at the corresponding wiki pages:
* [Usage](https://github.com/maxvaer/gitAnalyzer/wiki/Usage)
* [Template](https://github.com/maxvaer/gitAnalyzer/wiki/Template)


## Metadata Database
The included crawler was used to fetch the metadata of 1.000.000+ repositories.  
To get a dump of this dataset, take a look at this repository: [GitHub-Metadata](https://github.com/maxvaer/GitHub-Metadata)  
This dataset can be used, to query a list of URLs of repositories,  
which can be scanned using gitAnalyzer.


--------

<div align="center">

gitAnalyzer is distributed under [MIT License](https://github.com/maxvaer/gitAnalyzer/blob/main/LICENSE).

</div>

