# What touched
This is a small application, which will check your code and will try to fetch the features, which were touched by your changes. These results can be used during the manual testing of your product.

## How it works
You define the path where application need to check the files. It goes through all files in selected path and tries to find the `@featureType` comment where you define the type of features. Then it uses `git diff` to compare selected branches and based on the output it decides what kind of files were touched.
After that it goes through this files list and tries to find the files, where features are defined. At the end, application shows what can be potentially tested.

## How to use
In your code, please define the feature by writing of `@featureType {YOUR FEATURE NAME}` comment.
```php
<?php

/**
 * @featureType test functionality
 */
function firstFunction() {
    //Some code here
}
```

Run the command to find the touched files:
```shell script
./wt -path=/APSOLUTE/PATH/TO/YOUR/DIRECTORY -workingBranch=my-brand-new-branch
```

## Available command args
- `destinationBranch` (string)
Destination branch with which we will compare selected working branch. (default "master")
- `fileExt` (string)
The type of extension of the files which we need to check.
- `path` (string)
The type of vcs which will be used for retrieving diff information. (default ".")
- `pathToIgnoreFile` (string)
The path to file, where line-by-line written the list of paths which should be ignored. By default it's: .gitignore (default ".gitignore")
- `vcs` (string)
The type of vcs which will be used for retrieving diff information. (default "git")
- `workingBranch` (string)
Working branch which will be compared with the destination branch.

## Supported languages
Currently the application supports all languages where it's possible to define the `comments`. Eg: `#@featureType test` or `//@featureType test` or even `/* @featureType test */`

## Supported platforms
You can find this information in [`Supported OS` section of project build documentation](documentation/build.md).

## Where can be used
You can add it, for example, as step in your pipelines. So the QA engineer will see the results of application output and based on them will build the testing plan.
