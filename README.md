![Build status](https://github.com/function61/happy-api/workflows/Build/badge.svg)
[![Download](https://img.shields.io/github/downloads/function61/happy-api/total.svg?style=for-the-badge)](https://github.com/function61/happy-api/releases)

REST API for delivering happiness - hosted on AWS Lambda.

tl;dr: put URL https://function61.com/happy
in your application to enable your users to get their daily dose of happiness.


Use case
--------

I wanted to have a "Enjoy your day!" wish at the footer of a web app I offer for my users.
I wanted the "enjoy" word to be a link that takes the user to a random picture on the
internet that brings happiness:

![](docs/example-ui.png)

**click** gets you:

![](docs/example-happiness.png)


Can I too use the URL?
----------------------

Yes! And don't be afraid to use it - I make the following promises:

- The URL is the API and it won't change, or if it will the old URL will get redirected (i.e. still work)

- The pictures will be family friendly

- The service won't have ads, or if in the long term will have ads they will be unobtrusive text-only ads.


Adding new pictures
-------------------

Generate new ID for the picture with:

```console
$ ./happy-api new
nohH
```

Add new picture under [static/images/](static/images/)

Add image attribution URL with `$ exiftool` command:

```console
$ exiftool "-artist=https://example.com/" nohH.jpg
```

Send a pull request.


How to deploy
-------------

If for some reason you want to host your own API (you could just use the public API that
we host), follow these instructions.

Deployment is easiest using our [Deployer](https://github.com/function61/deployer) tool.
You don't need it and you can upload Lambda zip and configure API gateway manually if you want.

You have to do this only for the first time:

```
$ mkdir deployments
$ deployer deployment-init happy-api "url_to_deployerspec.zip_in_GitHub_releases"
Wrote /home/joonas/deployments/happy-api/user-config.json
```

Now edit above file with your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`

Then do the actual deployment:

```
$ version="..."
$ deployer deploy happy-api "same_url_as_in_init"
```
