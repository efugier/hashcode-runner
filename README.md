# Purpose

This simple cross-platform runner was made to make our team more efficient during the **Google Hashcode 2019**.

It can run a model on multiple datasets in parallel and automatically update the submissions/scores when the output is better so that you never miss or lose a great result.

It's designed to work with a _model_ and a _scorer_ executable that respectively create and rate a solution given a dataset.

# Usage

```
./build.sh
```
```
./runner -h
```

`main.sh` and `scorer.sh` are used as dummies to test the runner.

### Test the runner

This should get you started
```
./setup.sh && touch data/A.in data/B.in data/C.in
go run runner.go utils.go --datasets=ABC
```
and make it work with the provided `model.sh` and `scorer.sh`. Experiment from there !

# Requirements

## Dataset names

Currently, a dataset name must best one character and the dataset extension must be `.in`.

**ex:** `A.in`

## Scorer and model

### Scorer
The scorer must print the score in _stdout_ (and **only the score**) and require the two following arguments:
 1. the dataset
 2. the submission to test

**ex:** `./scorer.sh ./data/A.in ./submissions-tmp/A.out.tmp`

### Model
The model must require the following arguments:
 1. the dataset
 2. the output file path

and write the solution at the given path itself.

**ex:** `./model.sh ./data/A.in ./submissions-tmp/A.out.tmp`

## Folders

* The used folders must exist, the runner will not create them. `setup.sh` can create them for you.
* If you change the submission folder to `newfolder`, `newfolder-tmp` must also exists.

## Default file structure

```
.
├── scorer.sh
├── model.sh
├── data
│  ├── A.in
│  └── B.in
├── submissions
│  ├── A.out
│  ├── A.score
│  ├── B.out
│  └── B.score
└── submissions-tmp
   ├── A.out.tmp
   └── B.out.tmp
```

The `submissions/` folder always contains the best submission for each dataset: if a new submission is better than the old one, `submissions/` will contain the new one and `submissions-tmp/` the old one. Otherwise, the new one will be in `submission-tmp/`.
