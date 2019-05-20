[![Build Status](https://img.shields.io/travis/function61/onni.svg?style=for-the-badge)](https://travis-ci.org/function61/onni)
[![Download](https://img.shields.io/badge/Download-bintray%20latest-green.svg?style=for-the-badge)](https://bintray.com/function61/dl/onni/_latestVersion#files)

REST API for delivering happiness.

tl;dr: put [this URL](https://29ha8gbcmc.execute-api.us-east-1.amazonaws.com/prod/happy)
to your application to enable your users to get their daily dose of happiness.

NOTE: this is a very new project and the URL will change to a prettier domain soon.


Use case
--------

I wanted to have a "Have a happy day!" wish at the footer of a web app I offer for my
customers. I wanted the "happy" word to be a link that takes the user to a random picture
on the internet that brings happiness.


Contributing
------------

TODO


How to deploy
-------------

If for some reason you want to host your own API (you could just use the public API that
we host), follow these instructions.

Deployment is easiest using our [Deployer](https://github.com/function61/deployer) tool.
You don't need it and you can upload Lambda zip and configure API gateway manually if you want.

You have to do this only for the first time:

```
$ mkdir -p deployments
$ version="..." # find this from above Bintray link
$ deployer deployment-init onni "https://dl.bintray.com/function61/dl/onni/$version/deployerspec.zip"
Wrote /home/joonas/deployments/onni/user-config.json
```

Now edit above file with your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`

Then do the actual deployment:

```
$ version="..."
$ deployer deploy onni "https://dl.bintray.com/function61/dl/onni/$version/deployerspec.zip"
```
